package cmds

import (
	"fmt"
	"slices"
)

// var AvailableCmds = []string{"exit", "type", "echo", "pwd", "cd", "xyz_baz", "xyz_baz_foo", "xyz_baz_foo_quz"}
var AvailableCmds = []string{"exit", "type", "echo", "pwd", "cd"}

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
