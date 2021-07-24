package handlers

import (
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/services"
)

type PurgeHandler struct {
}

func (handler PurgeHandler) Handle(args []string) {
	if Help(handler, args) {
		return
	}

	if len(args) < 2 {
		handler.Help()
	}
	config := configurations.Config()
	nodeManager := services.GetNodeManager(config.Backend)
	for _, node := range args[1:] {
		nodeManager.Purge(node)
	}
}
func (handler PurgeHandler) Help() {

}
func (handler PurgeHandler) Desc() string {
	return "Purge a node"
}

func init() {
	name := "purge"
	CMD[name] = PurgeHandler{}
}
