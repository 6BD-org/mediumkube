package commands

import (
	"fmt"
	"mediumkube/utils"
	"os/exec"
)

// RootHandler root command handler
type RootHandler struct {
}

// CMD Sub commands shoud register itself to this map
var CMD = make(map[string]Handler)

// Help help
func (handler RootHandler) Help() {
	fmt.Print("Available commands:\n\r")
	for k := range CMD {
		fmt.Printf("  -%v: %v\n\r", k, CMD[k].Desc())
	}
}

// Handle handle
func (handler RootHandler) Handle(args []string) {

	multipassDelegated := []string{"list", "delete", "purge"}

	if len(args) < 2 {
		fmt.Printf("%v\n", "Insufficient arguments\n")
		handler.Help()
	} else {

		// Multipass commands are compatible
		for _, v := range multipassDelegated {
			if args[1] == v {
				cmd := exec.Command("multipass", args[1:]...)
				utils.ExecWithStdio(cmd)
				return
			}
		}

		switch args[1] {
		case "render":
			CMD["render"].Handle(args[1:])
		case "deploy":
			CMD["deploy"].Handle(args[1:])
		case "init":
			CMD["init"].Handle(args[1:])
		case "reset":
			CMD["reset"].Handle(args[1:])
		case "help":
			handler.Help()
		default:
			fmt.Printf("%v\n", "Invalid Command")
			handler.Help()
		}
	}

}
