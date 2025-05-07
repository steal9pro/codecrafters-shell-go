package cmds

import (
	"fmt"
	"strings"
)

func Echo(msg []string) {
	fmt.Println(strings.Join(msg, " "))
}
