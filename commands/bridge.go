package commands

import (
	"fmt"
	"mediumkube/configurations"
	"mediumkube/network"
	"mediumkube/utils"
)

type BridgeHandler struct {
}

func (handler BridgeHandler) Help() {
	fmt.Println("bridge [cmd]")
	fmt.Println("commands: ")
	fmt.Println("create\tCreate bridge defined in config")
	fmt.Println("delete\tDelete bridge defined in config")
}

func (handler BridgeHandler) Desc() string {
	return "Manage bridges"
}

func (handler BridgeHandler) Handle(args []string) {
	if Help(handler, args) {
		return
	}

	if len(args) != 2 {
		handler.Help()
		return
	}

	command := args[1]
	var err error
	config := configurations.Config()
	switch command {
	case "create":
		err = network.CreateNetBridge(config.Bridge)
	case "delete":
		err = network.RemoveNetBridge(config.Bridge)
	case "show":
		network.ShowBridge(config.Bridge)
	case "up":
		network.Up(config.Bridge)
	case "down":
		network.Down(config.Bridge)
	default:
		handler.Help()
	}

	utils.CheckErr(err)

}

func init() {
	name := "bridge"
	CMD[name] = BridgeHandler{}
}
