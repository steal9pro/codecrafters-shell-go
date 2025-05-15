package cmds

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Cmd interface {
	Run(args []string)
}

type Repl struct {
	osCmds map[string]string
}

func InitRepl() *Repl {
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

	return &Repl{
		osCmds: osCmds,
	}
}

func (r *Repl) CmdExist(cmdName string) (string, bool) {
	path, ok := r.osCmds[cmdName]

	return path, ok
}

func (r *Repl) Pwd() {
	absPath, err := os.Getwd()
	if err != nil {
		log.Printf("error during running: %v", err.Error())
		return
	}

	fmt.Println(absPath)
}

func (r *Repl) Cd(path string) {
	err := os.Chdir(path)
	if err != nil {
		fmt.Printf("%s: %s: %s \n", "cd", path, "No such file or directory")
	}
}

func NewCmd(repl *Repl, name string) Cmd {
	switch name {
	case "type":
		return InitType(repl)
	}
	return nil
}

func RunOSCmd(name string, args []string) {
	cmd := exec.Command(name, args...)
	byteResp, err := cmd.Output()
	if err != nil {
		log.Printf("error during running: %v", err.Error())
	}

	fmt.Printf("%s", string(byteResp))
}
