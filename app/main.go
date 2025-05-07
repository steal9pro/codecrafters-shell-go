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

		switch command {
		case "echo":
			cmds.Echo(args)
		case "exit":
			os.Exit(0)
		default:
			fmt.Println(command + ": command not found")
		}
	}
}
