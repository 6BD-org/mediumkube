package commands

import (
	"fmt"
	"mediumkube/configurations"
	"mediumkube/utils"
	"os/exec"
)

// RootHandler root command handler
type RootHandler struct {
}

// CMD Sub commands shoud register itself to this map
var CMD = make(map[string]Handler)
var multipassDelegated = []string{"list", "delete", "purge", "exec", "shell", "info", "launch", "start", "find"}

// Help help
func (handler RootHandler) Help() {
	fmt.Print("Available commands:\n\r")
	for k := range CMD {
		fmt.Printf("  -%v: %v\n\r", k, CMD[k].Desc())
	}

	fmt.Println("\n\rMultipass compatible commands:")
	for _, cmd := range multipassDelegated {
		fmt.Printf(" | %v", cmd)
	}
	fmt.Print("\n\r")
}

// Handle handle
func (handler RootHandler) Handle(args []string) {

	if len(args) < 2 {
		fmt.Printf("%v\n", "Insufficient arguments\n")
		handler.Help()
	} else {
		config := configurations.Config()

		if config.Backend == "multipass" {
			// Multipass commands are compatible
			for _, v := range multipassDelegated {
				if args[1] == v {
					cmd := exec.Command("multipass", args[1:]...)
					utils.AttachAndExec(cmd)
					return
				}
			}
		}

		switch args[1] {
		case "list":
			CMD["list"].Handle(args[1:])
		case "render":
			CMD["render"].Handle(args[1:])
		case "deploy":
			CMD["deploy"].Handle(args[1:])
		case "init":
			CMD["init"].Handle(args[1:])
		case "reset":
			CMD["reset"].Handle(args[1:])
		case "join":
			CMD["join"].Handle(args[1:])
		case "apply":
			CMD["apply"].Handle(args[1:])
		case "purge":
			CMD["purge"].Handle(args[1:])
		case "start":
			CMD["start"].Handle(args[1:])
		case "stop":
			CMD["stop"].Handle(args[1:])
		case "help":
			handler.Help()
		default:
			fmt.Printf("%v\n", "Invalid Command")
			handler.Help()
		}
	}

}
