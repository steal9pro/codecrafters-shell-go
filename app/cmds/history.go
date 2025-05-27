package cmds

import "fmt"

type History struct {
	repl *Repl
}

func InitHistory(repl *Repl) *History {
	return &History{
		repl: repl,
	}
}

func (h *History) Run(args []string) {
	fmt.Println("history")
}
