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

		// Write command to history for single commands
		if len(cmdPipe.Cmds) == 1 {
			cmd := cmdPipe.Cmds[0]
			repl.History.Write(fmt.Sprintf("%s %v", cmd.Command, strings.Join(cmd.Args, " ")))
		} else {
			// Write the entire pipe command to history
			var cmdStrings []string
			for _, cmd := range cmdPipe.Cmds {
				cmdStrings = append(cmdStrings, fmt.Sprintf("%s %s", cmd.Command, strings.Join(cmd.Args, " ")))
			}
			repl.History.Write(strings.Join(cmdStrings, " | "))
		}

		// Handle both single commands and pipes uniformly
		err = runner.RunPipeCmdsV2(repl, cmdPipe)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", cmdPipe.Cmds[0].Command, err)
		}
	}
}
