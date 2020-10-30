package commands

import (
	"flag"
	"fmt"
	"log"
	"mediumkube/common"
	"mediumkube/services"
	"mediumkube/utils"

	"gopkg.in/yaml.v3"
)

type InitHandler struct {
	flagSet *flag.FlagSet
}

func (handler InitHandler) Help() {

	fmt.Println("init [nodeName] kwargs")
	handler.flagSet.Usage()
}

func (InitHandler) Desc() string {
	return "Initialize a k8s cluster"
}

func initCmd(args []common.Arg) []string {
	cmd := []string{"kubeadm", "init"}
	for _, arg := range args {
		cmd = append(cmd, "--"+arg.Key)
		cmd = append(cmd, arg.Value)
	}
	return cmd
}

func (handler InitHandler) Handle(args []string) {

	configPath := handler.flagSet.String("config", "./config.yaml", "Config file that defines kubernetes init args")
	node := handler.flagSet.String("node", "node1", "Node to be inited")
	handler.flagSet.Parse(args[1:])

	if len(args) < 2 {
		fmt.Println("Insufficient arguments.")
		handler.Help()
	}

	if Help(handler, args) {
		return
	}

	configStr := utils.ReadByte(*configPath)

	overallConfig := common.OverallConfig{}

	err := yaml.Unmarshal(configStr, &overallConfig)
	utils.CheckErr(err)

	kubeInitArgs := overallConfig.KubeInit.Args
	cmd := initCmd(kubeInitArgs)
	services.MultipassService{}.Exec(*node, cmd, true)
	// TODO: Add post-command to enable kubectl

	log.Printf("Doing post-init configurations")
	services.MultipassService{}.ExecScript(*node, "./k8s/scripts/post-init.sh", false)

}

func init() {
	name := "init"

	handler := InitHandler{
		flagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}

	CMD[name] = handler
}
