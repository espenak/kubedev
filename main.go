package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(helpText)
		os.Exit(1)
	}
	var subcommand = os.Args[1]
	switch subcommand {
	case "build":
		Build(os.Args[2:])
	case "compile-templates":
		CompileTemplates(os.Args[2:])
	default:
		fmt.Println(helpText)
		os.Exit(1)
	}
}

const helpText = `Usage: kubedev <subcommand>

Subcommands:

  build               Build the docker images.
  compile-templates   Just compile the templates - do not build anything.
  up                  Bring the development environment up. 
    				  Shortcut for several of the other commands.
`
