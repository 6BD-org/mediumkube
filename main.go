package main

import (
	"flag"
	"mediumkube/commands"
	"mediumkube/configurations"
	"os"
)

func main() {

	// Setup global config
	tmpFlagSet := flag.NewFlagSet("", flag.ExitOnError)
	configDir := tmpFlagSet.String("config", "./config.yaml", "Configuration file")
	tmpFlagSet.Parse(os.Args)
	configurations.InitConfig(*configDir)

	// Handle command
	commands.RootHandler{}.Handle(os.Args)

}

func init() {

}
