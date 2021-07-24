package handlers

import (
	"fmt"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/services"
)

type StopHandler struct {
}

func (handler StopHandler) Handle(args []string) {

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
		manager.Stop(node)
	}
}
func (handler StopHandler) Help() {
	fmt.Println("stop node1 node2 node3")
}

func (handler StopHandler) Desc() string {
	return "Stop virtual machines"
}

func init() {
	name := "stop"
	CMD[name] = StopHandler{}
}
