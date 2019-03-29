package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func Build(args []string) {
	verbose := false
	cliFlags := flag.NewFlagSet("kubedev", flag.ExitOnError)
	cliFlags.BoolVar(&verbose, "v", false, "Verbose mode")
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
	if err != nil {
		log.Fatal(err)
	}

	if context.BuildTemplates() != nil {
		log.Fatal(err)
	}

	if verbose {
		fmt.Println("** Config: **")
		context.YamlPrint()
		fmt.Println("")
	}
}

func printUsage(flagSet *flag.FlagSet) {
	fmt.Fprintf(flagSet.Output(), buildHelpText)
	flagSet.PrintDefaults()
}

const buildHelpText = `Usage: kubedev build <context directory>

Options:
`
