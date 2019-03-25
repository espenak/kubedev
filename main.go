package main

import (
    "fmt"
    "os"
)

func main() {
  var subcommand = os.Args[1]
  switch subcommand {
  case "build":
    Build(os.Args[2:])
  // case "send":
  	// sendCommand.Parse(os.Args[2:])
  default:
  	fmt.Println(helpText)
  	os.Exit(2)
  }
  // fmt.Printf("Hello world!! %v %t\n", subcommand, verbose)
}


// func printUsage(cliFlags *flag.FlagSet) {
// 	fmt.Fprintf(cliFlags.Output(), helpTextHeader)
//   cliFlags.PrintDefaults()
//   fmt.Fprintf(cliFlags.Output(), helpTextFooter)
// }

const helpText = `Usage: kubedev <subcommand>

Subcommands:

  build    Build stuff
  up       Up stuff
`
