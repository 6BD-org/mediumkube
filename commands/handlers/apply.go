package handlers

import (
	"fmt"
	"mediumkube/pkg/services"
)

type ApplyHandler struct {
}

func (handler ApplyHandler) Help() {
	fmt.Println("apply [filePath]")
}

func (handler ApplyHandler) Desc() string {
	return "Apply K8s resources to cluster"
}

func (handler ApplyHandler) Handle(args []string) {
	if Help(handler, args) {
		return
	}

	if len(args) < 2 {
		handler.Help()
		panic("Too few arguments")
	}

	// TODO, get k8s client and apply changes
	file := args[1]
	services.GetK8sService().Apply(file)

}

func init() {
	name := "apply"
	CMD[name] = ApplyHandler{}
}
