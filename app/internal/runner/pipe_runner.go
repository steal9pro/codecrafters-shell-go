package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/codecrafters-io/shell-starter-go/app/cmds"
	"github.com/codecrafters-io/shell-starter-go/app/internal/output"
	"github.com/codecrafters-io/shell-starter-go/app/internal/reader"
)

type OutputStreamWriter struct {
	output output.Output
}

func (osw *OutputStreamWriter) Write(p []byte) (n int, err error) {
	osw.output.WriteStream(bytes.NewReader(p))
	return len(p), nil
}

type PipeRunner struct {
	repl     *cmds.Repl
	commands []*reader.Cmd
	ctx      context.Context
	cancel   context.CancelFunc
	pipes    []*io.PipeReader
	writers  []*io.PipeWriter
	wg       sync.WaitGroup
	errChan  chan error
}

func NewPipeRunner(repl *cmds.Repl, cmdPipe *reader.CmdsPipe) *PipeRunner {
	ctx, cancel := context.WithCancel(context.Background())

	numPipes := len(cmdPipe.Cmds) - 1
	pipes := make([]*io.PipeReader, numPipes)
	writers := make([]*io.PipeWriter, numPipes)

	for i := 0; i < numPipes; i++ {
		pipes[i], writers[i] = io.Pipe()
	}

	return &PipeRunner{
		repl:     repl,
		commands: cmdPipe.Cmds,
		ctx:      ctx,
		cancel:   cancel,
		pipes:    pipes,
		writers:  writers,
		errChan:  make(chan error, len(cmdPipe.Cmds)),
	}
}

func RunPipeCmdsV2(repl *cmds.Repl, cmdPipe *reader.CmdsPipe) error {
	if cmdPipe == nil {
		return ErrInvalidCommand
	}

	if len(cmdPipe.Cmds) == 0 {
		return ErrEmptyCommand
	}

	if len(cmdPipe.Cmds) == 1 {
		return RunSingleCmd(repl, cmdPipe.Cmds[0])
	}

	runner := NewPipeRunner(repl, cmdPipe)
	defer runner.cleanup()

	return runner.execute()
}

func (pr *PipeRunner) execute() error {
	for i, cmd := range pr.commands {
		pr.wg.Add(1)
		go pr.runCommand(i, cmd)
	}

	// Monitor for completion and handle cleanup
	go pr.monitor()

	pr.wg.Wait()
	close(pr.errChan)

	// Return first error if any
	for err := range pr.errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (pr *PipeRunner) runCommand(index int, cmd *reader.Cmd) {
	defer pr.wg.Done()

	select {
	case <-pr.ctx.Done():
		return
	default:
	}

	var err error
	if isBuiltinCommandV2(cmd.Command) {
		err = pr.runBuiltinCommand(index, cmd)
	} else {
		err = pr.runExternalCommand(index, cmd)
	}

	if err != nil {
		pr.errChan <- err
		pr.cancel() // Cancel other commands on error
	}
}

func (pr *PipeRunner) runBuiltinCommand(index int, cmd *reader.Cmd) error {
	var cmdOutput output.Output = pr.repl.GetOutput()

	// If not the last command, redirect to pipe
	if index < len(pr.writers) {
		cmdOutput = &output.PipeOutput{Writer: pr.writers[index]}
		defer pr.writers[index].Close()
	}

	// If not the first command, consume input from previous pipe to prevent hanging
	if index > 0 {
		go func() {
			io.Copy(io.Discard, pr.pipes[index-1])
		}()
	}

	// Create modified repl for this command
	cmdRepl := pr.createCommandRepl(cmdOutput)

	switch cmd.Command {
	case "echo":
		cmds.Echo(cmdRepl, cmd.Args)
	case "history":
		cmdRepl.History.Run(cmd.Args)
	case "type":
		exe := cmds.NewCmd(cmdRepl, cmd.Command)
		exe.Run(cmd.Args)
	case "pwd":
		cmdRepl.Pwd()
	case "cd":
		if len(cmd.Args) > 0 {
			cmdRepl.Cd(cmd.Args[0])
		}
	case "exit":
		cmdRepl.History.Close()
		os.Exit(0)
	}

	return nil
}

func (pr *PipeRunner) runExternalCommand(index int, cmd *reader.Cmd) error {
	_, ok := pr.repl.CmdExist(cmd.Command)
	if !ok {
		return fmt.Errorf("%s: command not found", cmd.Command)
	}

	execCmd := exec.CommandContext(pr.ctx, cmd.Command, cmd.Args...)

	if index > 0 {
		execCmd.Stdin = pr.pipes[index-1]
	}

	if index < len(pr.writers) {
		// Not the last command - pipe to next
		execCmd.Stdout = pr.writers[index]
		defer pr.writers[index].Close()
	} else {
		// Last command - create pipe to terminal
		stdout, err := execCmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to create stdout pipe: %v", err)
		}

		go func() {
			osw := &OutputStreamWriter{output: pr.repl.GetOutput()}
			io.Copy(osw, stdout)
		}()
	}

	// Set up stderr (always goes to terminal)
	stderr, err := execCmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	go func() {
		pr.repl.GetErrorOutput().WriteStream(stderr)
	}()

	// Start and wait for command
	if err := execCmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- execCmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					if status.Signal() == syscall.SIGPIPE {
						return nil // SIGPIPE is normal
					}
				}
			}
			return err
		}
		return nil
	case <-pr.ctx.Done():
		if execCmd.Process != nil {
			execCmd.Process.Kill()
		}
		<-done // Wait for process to actually exit
		return nil
	}
}

func (pr *PipeRunner) monitor() {
	// Give commands a moment to start
	time.Sleep(10 * time.Millisecond)

	// Monitor for early completion
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-pr.ctx.Done():
			return
		case <-ticker.C:
			// Check if any downstream processes have exited
			continue
		}
	}
}

func (pr *PipeRunner) cleanup() {
	pr.cancel()

	for _, writer := range pr.writers {
		if writer != nil {
			writer.Close()
		}
	}

	// Give processes time to clean up
	time.Sleep(50 * time.Millisecond)
}

func (pr *PipeRunner) createCommandRepl(cmdOutput output.Output) *cmds.Repl {
	newRepl := &cmds.Repl{}
	*newRepl = *pr.repl

	newRepl.SetOutput(cmdOutput)
	return newRepl
}

func isBuiltinCommandV2(command string) bool {
	builtins := map[string]bool{
		"echo":    true,
		"history": true,
		"type":    true,
		"pwd":     true,
		"cd":      true,
		"exit":    true,
	}
	return builtins[command]
}
