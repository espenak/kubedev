package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CliFlags struct {
	FlagSet              *flag.FlagSet
	UsageHeader          string
	UsageFooter          string
	Verbose              bool
	ContextRootDirectory string
}

func NewCliFlags(withContextRootDirectory bool) *CliFlags {
	cliFlags := CliFlags{}
	cliFlags.FlagSet = flag.NewFlagSet("kubedev", flag.ExitOnError)
	cliFlags.FlagSet.BoolVar(&cliFlags.Verbose, "v", false, "Verbose mode")
	if withContextRootDirectory {
		cliFlags.FlagSet.StringVar(&cliFlags.ContextRootDirectory, "d", ".", "Kubedev context directory")
	}
	cliFlags.FlagSet.Usage = func() { cliFlags.Usage() }
	return &cliFlags
}

func (cliFlags *CliFlags) Usage() {
	fmt.Fprintf(cliFlags.FlagSet.Output(), "Usage: ")
	fmt.Fprintf(cliFlags.FlagSet.Output(), cliFlags.UsageHeader)
	fmt.Fprintf(cliFlags.FlagSet.Output(), "\n\nOptions:\n")
	cliFlags.FlagSet.PrintDefaults()
	fmt.Fprintf(cliFlags.FlagSet.Output(), cliFlags.UsageFooter)
}

func (cliFlags *CliFlags) Parse(args []string) {
	if err := cliFlags.FlagSet.Parse(args); err != nil {
		cliFlags.Usage()
		os.Exit(1)
	}
}

func (cliFlags *CliFlags) NewContext() *Context {
	context, err := NewContext(cliFlags.ContextRootDirectory, cliFlags.Verbose)
	if err != nil {
		log.Fatal(err)
	}
	return context
}
