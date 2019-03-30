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
	// case "send":
	// sendCommand.Parse(os.Args[2:])
	default:
		fmt.Println(helpText)
		os.Exit(1)
	}
}

const helpText = `Usage: kubedev <subcommand>

Subcommands:

  build    Build stuff
  up       Up stuff
`
