package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/cmds"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		rawCmd := strings.Split(input[:len(input)-1], " ")
		command := rawCmd[0]
		args := rawCmd[1:]

		repl := cmds.InitRepl()

		switch command {
		case "echo":
			cmds.Echo(args)
		case "type":
			exe := cmds.NewCmd(repl, "type")
			exe.Run(args)
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
