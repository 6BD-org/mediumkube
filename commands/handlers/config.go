package handlers

import (
	"fmt"
	"mediumkube/pkg/utils"
	"os/exec"

	"k8s.io/klog/v2"
)

type ConfigHandler struct {
}

var (
	textBackends = []string{"vim", "vi"}
	backendMap   = map[string]func(str string) *exec.Cmd{}
)

func (h ConfigHandler) Desc() string {
	return "Config mediumkube"
}

func (h ConfigHandler) Help() {
	fmt.Println("mediumkube config. Supported text backends are: ")
	for _, tb := range textBackends {
		fmt.Printf("\t- %s\n", tb)
	}
}

func (h ConfigHandler) Handle(args []string) {
	for _, tb := range textBackends {
		cmdFunc, ok := backendMap[tb]
		if !utils.BinaryExists(tb) {
			continue
		}
		if !ok {
			klog.Error("backend not registered", tb)
		}
		cmd := cmdFunc("/etc/mediumkube/config.yaml")
		utils.AttachAndExec(cmd)
		return
	}
	klog.Error("No supported backends")
	h.Help()
}

func init() {
	name := "config"
	CMD[name] = ConfigHandler{}

	backendMap["vim"] = func(str string) *exec.Cmd {
		cmd := exec.Command("vim", str)
		return cmd
	}

	backendMap["vi"] = func(str string) *exec.Cmd {
		cmd := exec.Command("vi", str)
		return cmd
	}
}
