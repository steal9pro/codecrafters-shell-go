package cmds

import (
	"fmt"
	"slices"
)

type Type struct {
	availableCmds []string
}

func InitType() *Type {
	availableCmds := []string{"exit", "type", "echo"}
	return &Type{
		availableCmds: availableCmds,
	}
}

func (t *Type) Run(args []string) {
	has := slices.Contains(t.availableCmds, args[0])
	if has {
		fmt.Printf("%v is a shell builtin \n", args[0])
		return
	}

	fmt.Printf("%v: not found \n", args[0])
}
