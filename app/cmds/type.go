package cmds

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Type struct {
	osCmds        map[string]string
	availableCmds []string
}

func InitType() *Type {
	availableCmds := []string{"exit", "type", "echo"}

	pathEnv := os.Getenv("PATH")
	pathArr := strings.Split(pathEnv, ":")

	osCmds := make(map[string]string, 0)
	for _, dirPath := range pathArr {
		files, err := os.ReadDir(dirPath)
		if err != nil {
			fmt.Errorf("error during reading dir: %v \n", err.Error())
			continue
		}

		for _, entry := range files {
			name := entry.Name()
			_, ok := osCmds[name]
			if ok {
				continue
			}
			osCmds[name] = fmt.Sprintf("%s/%s", dirPath, name)
		}
	}

	return &Type{
		availableCmds: availableCmds,
		osCmds:        osCmds,
	}
}

func (t *Type) Run(args []string) {
	searchableBin := args[0]
	has := slices.Contains(t.availableCmds, searchableBin)
	if has {
		fmt.Printf("%v is a shell builtin \n", searchableBin)
		return
	}

	val, ok := t.osCmds[searchableBin]

	if !ok {
		fmt.Printf("%v: not found \n", searchableBin)
		return
	}

	fmt.Printf("%s is %s \n", searchableBin, val)
}
