package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		command := strings.Split(input[:len(input)-1], " ")

		switch command[0] {
		case "exit":
			os.Exit(0)
		default:
			fmt.Println(command[0] + ": command not found")
		}
	}
}
