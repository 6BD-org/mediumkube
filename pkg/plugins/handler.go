package plugins

import (
	"fmt"

	"k8s.io/klog/v2"
)

func handleList() {
	for k := range Plugins {
		fmt.Println("-", k)
	}
}

func handleDesc(pluginName string) {
	plugin, ok := Plugins[pluginName]
	if !ok {
		klog.Error("plugin not found")
		return
	}

	plugin.Desc()
}

func handleExec(pluginName string, args ...string) {
	plugin, ok := Plugins[pluginName]
	if !ok {
		klog.Error("plugin not found")
		return
	}

	plugin.Exec(args...)
}

// Handle plugin commands
func Handle(args []string) {
	if len(args) < 1 {
		klog.Error("Invalid arguments")
		return
	}
	if args[0] == "list" {
		handleList()
	}

	if args[0] == "desc" {
		if len(args) < 2 {
			klog.Error("Invalid arguments")
			return
		}
		handleDesc(args[1])
	}

	if args[0] == "exec" {
		if len(args) < 2 {
			klog.Error("Invalid arguments")
			return
		}
		if len(args) == 2 {
			handleExec(args[1])
			return
		}
		handleExec(args[1], args[2:]...)
	}
}
