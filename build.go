package main

import (
    "fmt"
    "flag"
    "os"
)

func Build (args []string) {
  verbose := false
  cliFlags := flag.NewFlagSet("kubedev", flag.ExitOnError)
  cliFlags.BoolVar(&verbose, "verbose", false, "Verbose mode")
  cliFlags.Usage = func() { printUsage(cliFlags) }

  if len(args) < 1 {
    printUsage(cliFlags)
    os.Exit(1)
  }
  contextRootDirectory := args[0]

  if err := cliFlags.Parse(args[1:]); err != nil {
		cliFlags.Usage()
		os.Exit(1)
	}

  // fmt.Printf("build! %t\n", verbose)
  context, err := NewContext(contextRootDirectory, verbose)
  if (err != nil) {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }

  if context.BuildTemplates() != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}

func printUsage (flagSet *flag.FlagSet) {
  fmt.Fprintf(flagSet.Output(), buildHelpText)
  flagSet.PrintDefaults()
}

const buildHelpText = `Usage: kubedev build <context directory>

Options:
`