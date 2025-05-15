package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/cmds"
	"github.com/codecrafters-io/shell-starter-go/app/internal/args"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		command, args := args.ParseArgs()

		repl := cmds.InitRepl()

		switch command {
		case "echo":
			cmds.Echo(args)
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
				fmt.Println(command + ": command not found")
				continue
			}
			cmds.RunOSCmd(command, args)
		}
	}
}
