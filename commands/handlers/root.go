package handlers

import (
	"fmt"
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

	fmt.Print("\n\r")
}

// Handle handle
func (handler RootHandler) Handle(args []string) {

	if len(args) < 1 {
		fmt.Printf("%v\n", "Insufficient arguments\n")
		handler.Help()
	} else {

		switch args[0] {
		case "help":
			handler.Help()
		default:
			cmd, ok := CMD[args[0]]
			if !ok {
				fmt.Printf("%v\n", "Invalid Command")
				handler.Help()
			}
			cmd.Handle(args)
		}
	}

}
