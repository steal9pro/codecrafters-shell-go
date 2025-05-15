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
	availableCmds := []string{"exit", "type", "echo", "pwd"}

	return &Type{
		repl:          repl,
		availableCmds: availableCmds,
	}
}

func (t *Type) Run(args []string) {
	searchableBin := args[0]
	has := slices.Contains(t.availableCmds, searchableBin)
	if has {
		fmt.Printf("%v is a shell builtin \n", searchableBin)
		return
	}

	path, ok := t.repl.CmdExist(searchableBin)

	if !ok {
		fmt.Printf("%v: not found \n", searchableBin)
		return
	}

	fmt.Printf("%s is %s \n", searchableBin, path)
}
