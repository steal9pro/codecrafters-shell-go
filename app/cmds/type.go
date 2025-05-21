package cmds

import (
	"fmt"
	"slices"
)

type Type struct {
	repl          *Repl
	availableCmds []string
}

func InitType(repl *Repl) *Type {
	availableCmds := []string{"exit", "type", "echo", "pwd", "cd"}

	return &Type{
		repl:          repl,
		availableCmds: availableCmds,
	}
}

func (t *Type) Run(args []string) {
	searchableBin := args[0]
	has := slices.Contains(t.availableCmds, searchableBin)
	if has {
		t.repl.Print(fmt.Sprintf("%v is a shell builtin", searchableBin))
		return
	}

	path, ok := t.repl.CmdExist(searchableBin)

	if !ok {
		t.repl.Print(fmt.Sprintf("%v: not found", searchableBin))
		return
	}

	t.repl.Print(fmt.Sprintf("%s is %s", searchableBin, path))
}
