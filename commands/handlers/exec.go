package handlers

import (
	"flag"
	"fmt"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/services"

	"k8s.io/klog/v2"
)

// ExecHandler provides ability to execute commands on nodes
// For multipass backend, this is already supported by multipass
// command redirecting
type ExecHandler struct {
	flagset *flag.FlagSet
}

func (handler ExecHandler) Help() {
	fmt.Printf("exec nodeName cmd arg1 arg2 ... --option1 val1 ...")
}

func (handler ExecHandler) Desc() string {
	return "Execute a command on a node"
}

func (handler ExecHandler) Handle(args []string) {
	if len(args) < 3 {
		klog.Error("Invalid arguments")
		handler.Help()
	}

	services.GetDomainManager(configurations.Config().Backend).Exec(args[1], args[2:], true)
}

func init() {
	name := "exec"
	CMD[name] = ExecHandler{
		flagset: flag.NewFlagSet(name, flag.ExitOnError),
	}

}
