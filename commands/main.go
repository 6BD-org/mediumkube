package main

import (
	"flag"
	"mediumkube/commands/handlers"
	"os"
)

func main() {

	// Setup global config
	tmpFlagSet := flag.NewFlagSet("", flag.ExitOnError)
	tmpFlagSet.Parse(os.Args[1:])
	// Handle command
	handlers.RootHandler{}.Handle(tmpFlagSet.Args())

}

func init() {

}
