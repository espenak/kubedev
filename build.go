package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func Build(args []string) {
	var verbose bool
	var contextRootDirectory string
	cliFlags := flag.NewFlagSet("kubedev", flag.ExitOnError)
	cliFlags.BoolVar(&verbose, "v", false, "Verbose mode")
	cliFlags.StringVar(&contextRootDirectory, "d", "", "Kubedev context directory")
	cliFlags.Usage = func() { printBuildUsage(cliFlags) }

	if len(args) < 1 {
		printBuildUsage(cliFlags)
		os.Exit(1)
	}

	if err := cliFlags.Parse(args); err != nil {
		cliFlags.Usage()
		os.Exit(1)
	}

	// fmt.Printf("build! %t\n", verbose)
	context, err := NewContext(contextRootDirectory, verbose)
	if err != nil {
		log.Fatal(err)
	}

	if verbose {
		fmt.Println("** Config: **")
		context.YamlPrint()
		fmt.Println("")
	}

	if context.BuildTemplates() != nil {
		log.Fatal(err)
	}
	context.BuildAllDockerImages()
}

func printBuildUsage(flagSet *flag.FlagSet) {
	fmt.Fprintf(flagSet.Output(), buildHelpText)
	flagSet.PrintDefaults()
}

const buildHelpText = `Usage: kubedev build <context directory>

Options:
`
