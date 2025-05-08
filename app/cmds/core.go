package cmds

type Cmd interface {
	Run(args []string)
}

func NewCmd(name string) Cmd {
	switch name {
	case "type":
		return InitType()
	}
	return nil
}
