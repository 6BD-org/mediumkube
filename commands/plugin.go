package commands

import (
	"fmt"
	"mediumkube/plugins"
)

type PluginHandler struct{}

func (handler PluginHandler) Desc() string {
	return "Plugin Operations"
}

func (handler PluginHandler) Help() {
	fmt.Println("plugin list")
	fmt.Printf("\tList all plugins\n")

	fmt.Println("plugin [exec|desc] [plugin name]")
	fmt.Printf("\tExecute or describe plugin\n")
}

func (handler PluginHandler) Handle(args []string) {
	if Help(handler, args) {
		return
	}

	plugins.Handle(args[1:])

}

func init() {
	name := "plugin"
	CMD[name] = PluginHandler{}
}
