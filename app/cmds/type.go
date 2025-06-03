package cmds

import (
	"fmt"
	"slices"
)

var AvailableCmds = []string{"exit", "type", "echo", "pwd", "cd", "history"}

type Type struct {
	repl          *Repl
	availableCmds []string
}

func InitType(repl *Repl) *Type {
	return &Type{
		repl:          repl,
		availableCmds: AvailableCmds,
	}
}

func (t *Type) Run(args []string) {
	if len(args) == 0 {
		t.repl.PrintError("type: missing operand")
		return
	}

	searchableBin := args[0]
	has := slices.Contains(t.availableCmds, searchableBin)
	if has {
		t.repl.Print(fmt.Sprintf("%v is a shell builtin\n", searchableBin))
		return
	}

	path, ok := t.repl.CmdExist(searchableBin)

	if !ok {
		t.repl.Print(fmt.Sprintf("%v: not found\n", searchableBin))
		return
	}

	t.repl.Print(fmt.Sprintf("%s is %s\n", searchableBin, path))
}
