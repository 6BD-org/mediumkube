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
	configDir := tmpFlagSet.String("config", "/etc/mediumkube/config.yaml", "Configuration file")
	tmpFlagSet.Parse(os.Args[1:])
	configurations.InitConfig(*configDir)
	// Handle command
	commands.RootHandler{}.Handle(tmpFlagSet.Args())

}

func init() {

}
