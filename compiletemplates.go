package main

import (
	"log"
)

func CompileTemplates(args []string) {
	cliFlags := NewCliFlags(true)
	cliFlags.UsageHeader = "kubedev compile-templates <context directory>"
	cliFlags.Parse(args)
	context := cliFlags.NewContext()

	if err := context.BuildTemplates(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Templates built into: %#v", context.BuiltTemplatesDirectory())
}
