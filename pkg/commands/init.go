package commands

import (
	"flag"
	"fmt"
	"log"
	"mediumkube/pkg/common"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/k8s"
	"mediumkube/pkg/plugins"
	"mediumkube/pkg/services"
	"mediumkube/pkg/utils"
	"os"
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
	services.GetNodeManager(overallConfig.Backend).Exec(*node, cmd, true)
	// TODO: Add post-command to enable kubectl

	log.Printf("Doing post-init configurations")
	plugins.Plugins["kube_post_init"].Exec(*node)
	// Transfer kube-config
	log.Printf("Mkdir %v\n", utils.GetFileDir(k8s.KubeConfigPath(overallConfig)))

	// We need execute permission here for cascade mkdir
	err := os.MkdirAll(utils.GetFileDir(k8s.KubeConfigPath(overallConfig)), os.FileMode(0777))
	utils.CheckErr(err)
	log.Println("Transferring kube config")
	services.GetNodeManager(overallConfig.Backend).Transfer(fmt.Sprintf("%v:%v", *node, overallConfig.VMKubeConfigDir), k8s.KubeConfigPath(overallConfig))

}

func init() {
	name := "init"

	handler := InitHandler{
		flagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}

	CMD[name] = handler
}
