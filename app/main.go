package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/cmds"
	"github.com/codecrafters-io/shell-starter-go/app/internal/reader"
	"github.com/codecrafters-io/shell-starter-go/app/internal/runner"
)

func main() {
	repl := cmds.InitRepl()
	defer repl.History.Close()

	for {
		repl.ResetOutput()
		streamReader := reader.NewStreamReader(repl.GetTrieNode(), repl.History)

		fmt.Fprint(os.Stdout, "$ ")

		cmdPipe, err := streamReader.ReadCommand()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading command: %v\n", err)
			continue
		}

		if cmdPipe == nil {
			continue
		}

		if len(cmdPipe.Cmds) == 0 {
			continue
		}

		if len(cmdPipe.Cmds) == 1 {
			cmd := cmdPipe.Cmds[0]
			repl.History.Write(fmt.Sprintf("%s %v", cmd.Command, strings.Join(cmd.Args, " ")))

			err := runner.RunSingleCmd(repl, cmd)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
			}
		}

		if len(cmdPipe.Cmds) > 1 {
			err := runner.RunPipeCmds(repl, cmdPipe)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
			}
		}
	}
}
