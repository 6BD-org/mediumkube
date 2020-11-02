package commands

import (
	"flag"
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
	handler.flagset.Usage()

}

func (handler DeployHandler) Handle(args []string) {

	handler.flagset.Parse(args[1:])

	if Help(handler, args) {
		return
	}

	overallConfig := configurations.Config()

	nodeConfig := overallConfig.NodeConfig
	services.GetMultipassService().Deploy(
		overallConfig.NodeNum,
		nodeConfig.CPU,
		nodeConfig.MEM,
		nodeConfig.DISK,
		overallConfig.Image,
		overallConfig.CloudInit,
	)
}

func init() {
	var name = "deploy"
	CMD[name] = DeployHandler{
		flagset: flag.NewFlagSet(name, flag.ExitOnError),
	}
}
