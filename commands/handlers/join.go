package handlers

import (
	"flag"
	"fmt"
	"log"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/services"
	"mediumkube/pkg/utils"
)

type JoinHandler struct {
	flagSet *flag.FlagSet
}

func (handler JoinHandler) Help() {
	fmt.Println("join [nodeToJoin] [masterNode]")
}

func (handler JoinHandler) Desc() string {
	return "Join a node to Kubernetes cluster"
}

func (handler JoinHandler) Handle(args []string) {

	if Help(handler, args) {
		return
	}

	if len(args) < 3 {
		handler.Help()
		return
	}
	nodeToJoin := args[1]
	master := args[2]

	log.Printf("Joining node %v to %v...\n\r", nodeToJoin, master)

	// Step 1. On master node, execute `kubeamd token create --print-join-command`
	// 	to obtain join command

	// Step 2. Execute that command on node to join

	mpSvc := services.GetNodeManager(configurations.Config().Backend)

	tokenCmd := []string{"kubeadm", "token", "create", "--print-join-command"}

	joinCommandStr := mpSvc.Exec(master, tokenCmd, true)
	joinCmd := utils.SplitCmd(joinCommandStr)

	log.Printf("Obtained join command: %v", joinCommandStr)

	mpSvc.AttachAndExec(nodeToJoin, joinCmd, true)
}

func init() {
	name := "join"
	CMD[name] = JoinHandler{
		flagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}
}
