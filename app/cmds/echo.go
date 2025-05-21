package cmds

import (
	"fmt"
	"strings"
)

func Echo(repl *Repl, msg []string) {
	repl.Print(fmt.Sprintf(strings.Join(msg, " ")))
}
