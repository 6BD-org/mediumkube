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
	handler.flagset.Usage()
}

func (handler TransferHandler) Handle(args []string) {

	recursive := handler.flagset.Bool("r", false, "If true, directory is copied recursively")
	handler.flagset.Parse(args[1:])
	if Help(handler, args) {
		return
	}

	if len(args) < 3 {
		handler.Help()
		return
	}

	args = handler.flagset.Args()

	fmt.Println(args, *recursive)
	if *recursive {
		services.GetNodeManager(configurations.Config().Backend).TransferR(args[0], args[1])
	} else {
		services.GetNodeManager(configurations.Config().Backend).Transfer(args[0], args[1])
	}

}

func init() {
	name := "transfer"
	CMD[name] = TransferHandler{
		flagset: flag.NewFlagSet(name, flag.ExitOnError),
	}
}
