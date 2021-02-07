package commands

import (
	"flag"
	"fmt"
	"mediumkube/configurations"
	"mediumkube/services"
)

type TransferHandler struct {
	flagset *flag.FlagSet
}

func (handler TransferHandler) Desc() string {
	return "Transfer a file from host to node"
}

func (handler TransferHandler) Help() {
	fmt.Println("transfer file1 node1:file2")
}

func (handler TransferHandler) Handle(args []string) {
	if Help(handler, args) {
		return
	}

	if len(args) < 3 {
		handler.Help()
		return
	}

	services.GetNodeManager(configurations.Config().Backend).Transfer(args[1], args[2])

}

func init() {
	name := "transfer"
	CMD[name] = TransferHandler{
		flagset: flag.NewFlagSet(name, flag.ExitOnError),
	}
}
