package commands

import (
	"flag"
	"mediumkube/common"
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

}

func (handler DeployHandler) Handle(args []string) {
	configPath := handler.flagset.String("config", "./config.yaml", "Config yaml for deployment")
	handler.flagset.Parse(args)
	configStr := utils.ReadByte(*configPath)

	overallConfig := common.OverallConfig{}

	err := yaml.Unmarshal(configStr, &overallConfig)
	utils.CheckErr(err)

	nodeConfig := overallConfig.NodeConfig

}

func init() {
	var name = "deploy"
	CMD[name] = DeployHandler{
		flagset: flag.NewFlagSet(name, flag.ExitOnError),
	}
}
