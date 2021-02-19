package plugins

import (
	"fmt"
	"mediumkube/configurations"
	"mediumkube/network"

	"k8s.io/klog/v2"
)

type OpenPortPlugin struct {
}

func _in(port string) []string {
	iface := configurations.Config().Bridge.Name
	return []string{
		"-p", "tcp",
		"--dport", fmt.Sprintf("%v", port),
		"-i", iface,
		"-j", "ACCEPT",
	}
}

func _out(port string) []string {
	iface := configurations.Config().Bridge.Name
	return []string{
		"-p", "tcp",
		"--sport", fmt.Sprintf("%v", port),
		"-o", iface,
		"-j", "ACCEPT",
	}
}

func (plugin OpenPortPlugin) Exec(args ...string) {
	if len(args) < 1 {
		klog.Error("Invalid argument. Port must be a number")
	}
	for _, portNum := range args {
		network.InsertRuleIfNotExists("INPUT", network.IPModAPP, _in(portNum)...)
		network.InsertRuleIfNotExists("OUTPUT", network.IPModAPP, _out(portNum)...)

	}
}

func (plugin OpenPortPlugin) Desc() {
	fmt.Println("This plugin is used to open a port on mediumkube bridge")
}

func init() {
	name := "open_port"
	Plugins[name] = OpenPortPlugin{}
}
