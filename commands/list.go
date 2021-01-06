package commands

import (
	"fmt"
	"mediumkube/configurations"
	"mediumkube/services"
)

type ListHandler struct {
}

func (handler ListHandler) Handle(args []string) {
	config := configurations.Config()
	manager := services.GetNodeManager(config.Backend)
	manager.List()
}
func (handler ListHandler) Help() {
	fmt.Println("list")
}
func (handler ListHandler) Desc() string {
	return "List nodes"
}

func init() {
	name := "list"
	CMD[name] = ListHandler{}
}
