package commands

import (
	"flag"
	"fmt"
	"log"
	"mediumkube/common"
	"mediumkube/configurations"
	"mediumkube/services"
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

	node := handler.flagSet.String("node", "node1", "Node to be inited")
	handler.flagSet.Parse(args[1:])

	if len(args) < 2 {
		fmt.Println("Insufficient arguments.")
		handler.Help()
	}

	if Help(handler, args) {
		return
	}

	overallConfig := configurations.Config()

	kubeInitArgs := overallConfig.KubeInit.Args
	cmd := initCmd(kubeInitArgs)
	services.GetMultipassService().Exec(*node, cmd, true)
	// TODO: Add post-command to enable kubectl

	log.Printf("Doing post-init configurations")
	services.GetMultipassService().ExecScript(*node, "./k8s/scripts/post-init.sh", false)

	// Transfer kube-config
	log.Println("Transferring kube config")
	services.GetMultipassService().Transfer(fmt.Sprintf("%v:%v", node, overallConfig.VMKubeConfigDir), kubeConfigPath(overallConfig))

}

func init() {
	name := "init"

	handler := InitHandler{
		flagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}

	CMD[name] = handler
}
