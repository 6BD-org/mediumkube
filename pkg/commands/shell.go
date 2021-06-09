package commands

import (
	"fmt"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/services"
)

// ShellHandler attaches to the bash of a node
type ShellHandler struct {
}

func (handler ShellHandler) Desc() string {
	return "Attach to the shell of a node"
}

func (handler ShellHandler) Help() {
	fmt.Println("shell [nodeToAttach]")
}

func (handler ShellHandler) Handle(args []string) {

	if Help(handler, args) {
		return
	}

	if len(args) < 2 {
		handler.Help()
		return
	}

	services.GetNodeManager(configurations.Config().Backend).Shell(args[1])
}

func init() {
	name := "shell"
	CMD[name] = ShellHandler{}
}
