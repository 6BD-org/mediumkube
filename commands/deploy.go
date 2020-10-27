package commands

import (
	"flag"
	"mediumkube/common"
	"mediumkube/services"
	"mediumkube/utils"

	"gopkg.in/yaml.v3"
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

	configPath := handler.flagset.String("config", "./config.yaml", "Config yaml for deployment")
	handler.flagset.Parse(args[1:])

	if Help(handler, args) {
		return
	}

	configStr := utils.ReadByte(*configPath)

	overallConfig := common.OverallConfig{}

	err := yaml.Unmarshal(configStr, &overallConfig)
	utils.CheckErr(err)

	nodeConfig := overallConfig.NodeConfig
	services.MultipassService{}.Deploy(
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
