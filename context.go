package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type TemplateVariables struct {
	ContextDirectory string
	Name             string
}

type Context struct {
	RootDirectory    string
	AbsRootDirectory string
	Verbose          bool
	Name             string
	Config           *Config
}

func NewContext(rootDirectory string, verbose bool) (*Context, error) {
	if _, err := os.Stat(rootDirectory); os.IsNotExist(err) {
		return nil, err
	}

	absRootDirectory, err2 := filepath.Abs(rootDirectory)
	if err2 != nil {
		return nil, err2
	}

	context := Context{
		RootDirectory:    rootDirectory,
		AbsRootDirectory: absRootDirectory,
		Verbose:          verbose,
		Name:             "tullball",
	}

	var err3 error
	context.Config, err3 = LoadConfigFromFile(context.ConfigPath())
	if err3 != nil {
		return nil, err3
	}

	return &context, nil
}

func (context Context) ConfigPath() string {
	return filepath.Join(context.RootDirectory, "kubedev.yml")
}

func (context Context) DockerContext() string {
	path, _ := filepath.Abs(filepath.Join(context.RootDirectory, context.Config.DockerContext))
	return path
}

func (context Context) TemplatesDirectory() string {
	return filepath.Join(context.RootDirectory, "templates")
}

func (context Context) BuiltTemplatesDirectory() string {
	return filepath.Join(context.RootDirectory, "_kubedev_built_templates")
}

func (context Context) ParseTemplates() (*template.Template, error) {
	return template.ParseGlob(
		filepath.Join(context.TemplatesDirectory(), "*"))
}

func (context Context) FullName() string {
	return fmt.Sprintf("kubedev-%s", context.Name)
}

func (context Context) TemplateVariables() *TemplateVariables {
	return &TemplateVariables{
		ContextDirectory: context.AbsRootDirectory,
		Name:             context.FullName(),
	}
}

func (context Context) BuildTemplate(parsedTemplates *template.Template, templateName string) error {
	outPath := filepath.Join(context.BuiltTemplatesDirectory(), templateName)
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	parsedTemplates.ExecuteTemplate(outFile, templateName, context.TemplateVariables())
	outFile.Sync()
	return nil
}

func (context Context) BuildTemplates() error {
	os.MkdirAll(context.BuiltTemplatesDirectory(), os.ModePerm)
	parsedTemplates, err := context.ParseTemplates()
	if err != nil {
		return err
	}
	for _, tpl := range parsedTemplates.Templates() {
		err := context.BuildTemplate(parsedTemplates, tpl.Name())
		if err != nil {
			return err
		}
	}
	return nil
}
