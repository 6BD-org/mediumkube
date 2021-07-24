package handlers

import (
	"flag"
	"fmt"
	"log"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/services"
)

type ResetHandler struct {
	flagSet *flag.FlagSet
}

func (handler ResetHandler) Help() {
	fmt.Println("Remove generated K8s files from a node")
	handler.flagSet.Usage()
}

func (ResetHandler) Desc() string {
	return "Reset a k8s node"
}

func (handler ResetHandler) Handle(args []string) {

	handler.flagSet.Parse(args[1:])
	fmt.Println(handler.flagSet.Args())

	if Help(handler, args) {
		return
	}

	node := args[1]

	cmd := []string{"kubeadm", "reset"}

	services.GetNodeManager(configurations.Config().Backend).AttachAndExec(node, cmd, true)

	log.Println("Removing custom kube config")
	config := configurations.Config()
	rmConfigCmd := []string{"rm", config.VMKubeConfigDir}
	services.GetNodeManager(configurations.Config().Backend).AttachAndExec(node, rmConfigCmd, true)

}

func init() {
	name := "reset"
	CMD[name] = ResetHandler{
		flagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}
}
