package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/cmds"
	"github.com/codecrafters-io/shell-starter-go/app/internal/args"
	"github.com/codecrafters-io/shell-starter-go/app/internal/output"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		command, args := args.ParseArgs()

		repl := cmds.InitRepl()
		redirectStdout, redirectStdErr, fileName := output.ParseRedirectIfPresent(args)

		if redirectStdout {
			repl.RedirectStdOutToFile(fileName)
			args = args[0 : len(args)-2]
		}

		if redirectStdErr {
			repl.RedirectStdErrToFile(fileName)
			args = args[0 : len(args)-2]
		}

		switch command {
		case "echo":
			cmds.Echo(repl, args)
		case "type":
			exe := cmds.NewCmd(repl, "type")
			exe.Run(args)
		case "pwd":
			repl.Pwd()
		case "cd":
			repl.Cd(args[0])
		case "exit":
			os.Exit(0)
		default:
			_, ok := repl.CmdExist(command)
			if !ok {
				repl.PrintError(fmt.Sprintf(command + ": command not found"))
				continue
			}
			cmds.RunOSCmd(repl, command, args)
		}
	}
}
