package handlers

import (
	"fmt"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/services"
)

type StartHandler struct {
}

func (handler StartHandler) Handle(args []string) {

	config := configurations.Config()

	if Help(handler, args) {
		return
	}

	manager := services.GetNodeManager(config.Backend)
	if len(args) < 2 {
		handler.Help()
		return
	}

	for _, node := range args[1:] {
		manager.Start(node)
	}
}
func (handler StartHandler) Help() {
	fmt.Println("start node1 node2 node3")
}

func (handler StartHandler) Desc() string {
	return "Start virtual machines"
}

func init() {
	name := "start"
	CMD[name] = StartHandler{}
}
