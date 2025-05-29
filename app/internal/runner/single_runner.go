package runner

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/cmds"
	"github.com/codecrafters-io/shell-starter-go/app/internal/output"
	"github.com/codecrafters-io/shell-starter-go/app/internal/reader"
)

var ErrCommandNotFound = fmt.Errorf("command not found")
var ErrInvalidCommand = fmt.Errorf("invalid command")
var ErrEmptyCommand = fmt.Errorf("empty command")

func RunSingleCmd(repl *cmds.Repl, cmdStruct *reader.Cmd) error {
	if cmdStruct == nil {
		return ErrInvalidCommand
	}

	if cmdStruct.Command == "" {
		return ErrEmptyCommand
	}

	args := cmdStruct.Args

	redirectStdout, redirectStdErr, appendStdout, appendStdErr, fileName := output.ParseRedirectIfPresent(args)

	if redirectStdout || appendStdout {
		repl.RedirectStdOutToFile(fileName, appendStdout)
		args = args[0 : len(args)-2]
	}

	if redirectStdErr || appendStdErr {
		repl.RedirectStdErrToFile(fileName, appendStdErr)
		args = args[0 : len(args)-2]
	}

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
		repl.Cd(args[0])
	case "exit":
		repl.History.Close()
		os.Exit(0)
	default:
		_, ok := repl.CmdExist(cmdStruct.Command)
		if !ok {
			repl.PrintError(fmt.Sprintf(cmdStruct.Command + ": command not found"))
			return ErrCommandNotFound
		}
		cmds.RunOSCmd(repl, cmdStruct.Command, args)
	}

	return nil
}
