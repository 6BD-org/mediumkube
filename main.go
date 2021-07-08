package main

import (
	"flag"
	"mediumkube/pkg/commands"
	"mediumkube/pkg/configurations"
	"os"
)

func main() {

	// Setup global config
	tmpFlagSet := flag.NewFlagSet("", flag.ExitOnError)
	tmpFlagSet.Parse(os.Args[1:])
	configurations.InitConfig()
	// Handle command
	commands.RootHandler{}.Handle(tmpFlagSet.Args())

}

func init() {

}
