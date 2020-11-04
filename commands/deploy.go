package commands

import (
	"flag"
	"mediumkube/configurations"
	"mediumkube/services"
	"os"

	"github.com/wylswz/logflog/flogger"
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

	go func() {
		os.RemoveAll(logPath(overallConfig))
		os.MkdirAll(logPath(overallConfig), 0777)
		flogger.FLog([]string{logPath(overallConfig)})
	}()

	mount := make(map[string]string)
	mount[logPath(overallConfig)] = overallConfig.VMLogDir

	nodeConfig := overallConfig.NodeConfig
	services.GetMultipassService().Deploy(
		overallConfig.NodeNum,
		nodeConfig.CPU,
		nodeConfig.MEM,
		nodeConfig.DISK,
		overallConfig.Image,
		overallConfig.CloudInit,
		mount,
	)
}

func init() {
	var name = "deploy"
	CMD[name] = DeployHandler{
		flagset: flag.NewFlagSet(name, flag.ExitOnError),
	}
}
