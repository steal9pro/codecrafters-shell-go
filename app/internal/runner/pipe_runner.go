package runner

import (
	"github.com/codecrafters-io/shell-starter-go/app/cmds"
	"github.com/codecrafters-io/shell-starter-go/app/internal/reader"
)

func RunPipeCmds(repl *cmds.Repl, cmdPipe *reader.CmdsPipe) error {
	if cmdPipe == nil {
		return ErrInvalidCommand
	}

	for _, cmd := range cmdPipe.Cmds {
		err := RunSingleCmd(repl, cmd)

		if err != nil {
			return err
		}
	}

	return nil
}
