package commands

import (
	"flag"
	"fmt"
	"mediumkube/common"
	"mediumkube/configurations"
	"mediumkube/services"
)

type DeployHandler struct {
	flagset *flag.FlagSet
}

func (handler DeployHandler) Desc() string {
	return "Deploy a new K8s cluster"
}

func (handler DeployHandler) Help() {
	fmt.Println("deploy node1 node2...")
	fmt.Println("Absence of node indicates deploy all defined nodes")

}

func (handler DeployHandler) Handle(args []string) {

	handler.flagset.Parse(args[1:])

	if Help(handler, args) {
		return
	}
	overallConfig := configurations.Config()
	var nodes []common.NodeConfig
	nodeNames := args[1:]

	if len(nodeNames) == 0 {
		nodes = overallConfig.NodeConfig
	} else {
		nodes = make([]common.NodeConfig, 0)
		nodeMap := make(map[string]common.NodeConfig)
		for _, n := range overallConfig.NodeConfig {
			nodeMap[n.Name] = n
		}
		for _, name := range nodeNames {
			node, ok := nodeMap[name]
			if !ok {
				panic(fmt.Sprintf("Node node defined: %v", name))
			}
			nodes = append(nodes, node)
		}
	}

	services.GetMultipassService().Deploy(
		nodes,
		overallConfig.CloudInit,
		overallConfig.Image,
	)
}

func init() {
	var name = "deploy"
	CMD[name] = DeployHandler{
		flagset: flag.NewFlagSet(name, flag.ExitOnError),
	}
}
