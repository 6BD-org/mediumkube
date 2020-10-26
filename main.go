package main

import (
	"mediumkube/commands"
	"os"
)

func main() {
	commands.RootHandler{}.Handle(os.Args)
}

func init() {

}
