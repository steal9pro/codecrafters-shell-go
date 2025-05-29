package cmds

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/internal/autocompletition"
	"github.com/codecrafters-io/shell-starter-go/app/internal/output"
)

type Cmd interface {
	Run(args []string)
}

type Repl struct {
	osCmds        map[string]string
	output        output.Output
	errorOutput   output.Output
	channelOutput *output.ChannelOutput
	trieNode      *autocompletition.TrieNode
	History       *History
}

func InitRepl() *Repl {
	pathEnv := os.Getenv("PATH")
	pathArr := strings.Split(pathEnv, ":")

	osCmds := make(map[string]string, 20)
	osCmdsArray := make([]string, 20)

	for _, dirPath := range pathArr {
		files, err := os.ReadDir(dirPath)
		if err != nil {
			fmt.Errorf("error during reading dir: %v \n", err.Error())
			continue
		}

		for _, entry := range files {
			name := entry.Name()
			_, ok := osCmds[name]
			if ok {
				continue
			}
			osCmds[name] = fmt.Sprintf("%s/%s", dirPath, name)
			osCmdsArray = append(osCmdsArray, name)
		}
	}

	wholeCmdsArray := append(osCmdsArray, AvailableCmds...)

	// Initialize the TrieNode for autocomplete
	rootNode := autocompletition.InitTrieNode()
	rootNode.LoadWordsToTrie(wholeCmdsArray)

	return &Repl{
		osCmds:        osCmds,
		output:        output.NewOutput(),
		errorOutput:   output.NewOutput(),
		channelOutput: nil,
		trieNode:      rootNode,
		History:       InitHistory(),
	}
}

func (r *Repl) ResetOutput() {
	r.output = output.NewOutput()
	r.errorOutput = output.NewOutput()
	r.channelOutput = nil
}

func (r *Repl) RedirectStdOutToChannel(channelOutput *output.ChannelOutput) {
	r.channelOutput = channelOutput
	r.output = channelOutput
}

func (r *Repl) GetChannelOutput() *output.ChannelOutput {
	return r.channelOutput
}

func (r *Repl) GetOutput() output.Output {
	return r.output
}

func (r *Repl) GetErrorOutput() output.Output {
	return r.errorOutput
}

func (r *Repl) PrintErrorStream(reader io.Reader) {
	r.errorOutput.WriteStream(reader, true)
}

func (r *Repl) ShowCmds() {
	r.Print("Available commands:")
	// r.trieNode.GetAllWords("e")
	r.trieNode.Display(0)
}

func (r *Repl) RedirectStdOutToFile(fileName string, append bool) {
	r.output = output.NewFileOutput(fileName, append)
}

func (r *Repl) RedirectStdErrToFile(fileName string, append bool) {
	r.errorOutput = output.NewFileOutput(fileName, append)
}

func (r *Repl) PrintError(msg string) {
	r.errorOutput.PrintError(fmt.Sprintf("%s", msg))
}

func (r *Repl) Print(msg string) {
	r.output.Print(msg)
}

func (r *Repl) CmdExist(cmdName string) (string, bool) {
	path, ok := r.osCmds[cmdName]

	return path, ok
}

func (r *Repl) Pwd() {
	absPath, err := os.Getwd()
	if err != nil {
		r.PrintError(fmt.Sprintf("error during running: %v", err.Error()))
		return
	}

	r.Print(absPath)
}

func (r *Repl) Cd(path string) {
	if strings.Contains(path, "~") {
		homePath := os.Getenv("HOME")
		path = strings.Replace(path, "~", homePath, 1)
	}

	err := os.Chdir(path)
	if err != nil {
		r.PrintError(fmt.Sprintf("%s: %s: %s", "cd", path, "No such file or directory"))
	}
}

func (r *Repl) GetTrieNode() *autocompletition.TrieNode {
	return r.trieNode
}

func NewCmd(repl *Repl, name string) Cmd {
	switch name {
	case "type":
		return InitType(repl)
	}
	return nil
}

func RunOSCmd(repl *Repl, name string, args []string) {
	cmd := exec.Command(name, args...)

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		repl.PrintError(fmt.Sprintf("error creating stdout pipe: %v", err.Error()))
		return
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		repl.PrintError(fmt.Sprintf("error creating stderr pipe: %v", err.Error()))
		return
	}

	// Start the command before reading from pipes
	if err := cmd.Start(); err != nil {
		repl.PrintError(fmt.Sprintf("error starting command: %v", err.Error()))
		return
	}

	// Use goroutines to handle stdout and stderr streams
	stdoutDone := make(chan bool)
	stderrDone := make(chan bool)

	go func() {
		repl.output.WriteStream(stdoutPipe, false)
		stdoutDone <- true
	}()

	go func() {
		repl.errorOutput.WriteStream(stderrPipe, true)
		stderrDone <- true
	}()

	// Wait for both streams to complete
	<-stdoutDone
	<-stderrDone

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Command failed with non-zero exit code
			// repl.PrintError(fmt.Sprintf("command exited with status %d", exitErr.ExitCode()))
		} else {
			repl.PrintError(fmt.Sprintf("error waiting for command: %v", err.Error()))
		}
	}
}
