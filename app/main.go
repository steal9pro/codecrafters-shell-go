package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/cmds"
	"github.com/codecrafters-io/shell-starter-go/app/internal/output"
	"github.com/codecrafters-io/shell-starter-go/app/internal/reader"
)

func main() {
	repl := cmds.InitRepl()
	defer repl.History.Close()

	for {
		repl.ResetOutput()
		streamReader := reader.NewStreamReader(repl.GetTrieNode())

		fmt.Fprint(os.Stdout, "$ ")

		command, args, err := streamReader.ReadCommand()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading command: %v\n", err)
			continue
		}

		if command == "" {
			continue
		}

		repl.History.Write(fmt.Sprintf("%s %v", command, strings.Join(args, " ")))

		redirectStdout, redirectStdErr, appendStdout, appendStdErr, fileName := output.ParseRedirectIfPresent(args)

		if redirectStdout || appendStdout {
			repl.RedirectStdOutToFile(fileName, appendStdout)
			args = args[0 : len(args)-2]
		}

		if redirectStdErr || appendStdErr {
			repl.RedirectStdErrToFile(fileName, appendStdErr)
			args = args[0 : len(args)-2]
		}

		switch command {
		case "echo":
			cmds.Echo(repl, args)
		case "history":
			repl.History.Run(args)
		case "type":
			exe := cmds.NewCmd(repl, command)
			exe.Run(args)
		case "pwd":
			repl.Pwd()
		case "cd":
			repl.Cd(args[0])
		case "exit":
			repl.History.Close()
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
