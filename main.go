package main

import (
	"flag"
	"mediumkube/pkg/commands"
	"os"
)

func main() {

	// Setup global config
	tmpFlagSet := flag.NewFlagSet("", flag.ExitOnError)
	tmpFlagSet.Parse(os.Args[1:])
	// Handle command
	commands.RootHandler{}.Handle(tmpFlagSet.Args())

}

func init() {

}
