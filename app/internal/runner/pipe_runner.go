package runner

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/codecrafters-io/shell-starter-go/app/cmds"
	"github.com/codecrafters-io/shell-starter-go/app/internal/input"
	"github.com/codecrafters-io/shell-starter-go/app/internal/output"
	"github.com/codecrafters-io/shell-starter-go/app/internal/reader"
)

func RunPipeCmds(repl *cmds.Repl, cmdPipe *reader.CmdsPipe) error {
	if cmdPipe == nil {
		return ErrInvalidCommand
	}

	if len(cmdPipe.Cmds) == 0 {
		return ErrEmptyCommand
	}

	if len(cmdPipe.Cmds) == 1 {
		return RunSingleCmd(repl, cmdPipe.Cmds[0])
	}

	// Create channels for communication between commands
	channels := make([]*output.ChannelOutput, len(cmdPipe.Cmds)-1)
	for i := 0; i < len(cmdPipe.Cmds)-1; i++ {
		channels[i] = output.NewChannelOutput()
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(cmdPipe.Cmds))

	// Execute each command in the pipe
	for i, cmd := range cmdPipe.Cmds {
		wg.Add(1)
		go func(index int, command *reader.Cmd) {
			defer wg.Done()

			// Create a new repl instance for this command
			cmdRepl := createPipeRepl(repl, index, channels)

			err := runPipeCommand(cmdRepl, command, index, channels)
			fmt.Println("complete cmd ", index)
			if err != nil {
				errChan <- err
			}
		}(i, cmd)
	}

	// Wait for all commands to complete
	go func() {
		wg.Wait()
		close(errChan)

		// Close all channels
		for _, ch := range channels {
			ch.Close()
		}
	}()

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func createPipeRepl(originalRepl *cmds.Repl, cmdIndex int, channels []*output.ChannelOutput) *cmds.Repl {
	// Create a copy of the repl with modified input/output
	newRepl := &cmds.Repl{}
	*newRepl = *originalRepl

	// Set output for this command
	if cmdIndex < len(channels) {
		// Not the last command, redirect output to channel
		newRepl.RedirectStdOutToChannel(channels[cmdIndex])
	} else {
		// Last command, use original output
		newRepl.ResetOutput()
	}

	return newRepl
}

func runPipeCommand(repl *cmds.Repl, cmdStruct *reader.Cmd, cmdIndex int, channels []*output.ChannelOutput) error {
	if cmdStruct == nil {
		return ErrInvalidCommand
	}

	if cmdStruct.Command == "" {
		return ErrEmptyCommand
	}

	args := cmdStruct.Args

	switch cmdStruct.Command {
	case "echo":
		cmds.Echo(repl, args)
	case "history":
		repl.History.Run(args)
	case "type":
		exe := cmds.NewCmd(repl, cmdStruct.Command)
		exe.Run(args)
	case "pwd":
		repl.Pwd()
	case "cd":
		if len(args) > 0 {
			repl.Cd(args[0])
		}
	case "exit":
		repl.History.Close()
		os.Exit(0)
	default:
		path, ok := repl.CmdExist(cmdStruct.Command)
		if !ok {
			return fmt.Errorf("%s: command not found", cmdStruct.Command)
		}
		return runOSCmdInPipe(repl, path, cmdStruct.Command, args, cmdIndex, channels)
	}

	return nil
}

func runOSCmdInPipe(repl *cmds.Repl, path, name string, args []string, cmdIndex int, channels []*output.ChannelOutput) error {
	cmd := exec.Command(name, args...)

	// Set up input from previous command if not the first
	if cmdIndex > 0 {
		fmt.Printf("command index %v \n", cmdIndex)
		prevChannel := channels[cmdIndex-1].GetChannel()
		channelInput := input.NewChannelInput(prevChannel)
		cmd.Stdin = channelInput
	}

	// Set up output
	if cmdIndex < len(channels) {
		fmt.Println("not the last cmd ", cmdIndex)
		// Not the last command, pipe output to next command
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("error creating stdout pipe: %v", err)
		}

		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("error creating stderr pipe: %v", err)
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("error starting command: %v", err)
		}

		// Handle stdout and stderr
		go func() {
			repl.GetChannelOutput().WriteStream(stdoutPipe, false)
		}()

		go func() {
			repl.PrintErrorStream(stderrPipe)
		}()

		// Wait for command to complete
		return cmd.Wait()
	} else {
		fmt.Println("last cmd")
		// Last command, output to terminal
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("error creating stdout pipe: %v", err)
		}

		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("error creating stderr pipe: %v", err)
		}

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("error starting command: %v", err)
		}

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			repl.GetOutput().WriteStream(stdoutPipe, false)
		}()

		go func() {
			defer wg.Done()
			repl.GetErrorOutput().WriteStream(stderrPipe, true)
		}()

		wg.Wait()
		return cmd.Wait()
	}
}
