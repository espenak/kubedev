package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v2"
)

type TemplateVariables struct {
	RootDirectory string
	FullName      string
}

type Context struct {
	RootDirectory string `yaml:"rootDirectory"`
	ApiVersion    string `yaml:"apiVersion"`
	Name          string `yaml:"name"`
	DockerContext string `yaml:"dockerContext"`
	Verbose       bool   `yaml:"verbose"`
}

func NewContext(rootDirectory string, verbose bool) (*Context, error) {
	if _, err := os.Stat(rootDirectory); os.IsNotExist(err) {
		return nil, err
	}

	absRootDirectory, err2 := filepath.Abs(rootDirectory)
	if err2 != nil {
		return nil, err2
	}

	context := Context{}
	err3 := context.loadConfigFile(absRootDirectory)
	if err3 != nil {
		return nil, err3
	}
	context.RootDirectory = absRootDirectory
	context.Verbose = verbose

	validationError := context.clean()
	if validationError != nil {
		return nil, validationError
	}

	return &context, nil
}

func (context *Context) loadConfigFile(absRootDirectory string) error {
	configFilePath := filepath.Join(absRootDirectory, "kubedev.yml")
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err2 := yaml.Unmarshal(content, &context)
	if err2 != nil {
		return err2
	}

	return nil
}

func (context *Context) clean() error {
	validationError := func(message string) error {
		return fmt.Errorf("%s: %s", context.ConfigPath(), message)
	}
	if context.Name == "" {
		return validationError("name is required.")
	}
	if context.ApiVersion == "" {
		return validationError("apiVersion must be 'v1'.")
	}
	if context.DockerContext == "" {
		context.DockerContext = context.RootDirectory
	} else {
		dockerContext, err := filepath.Abs(filepath.Join(context.RootDirectory, context.DockerContext))
		if err != nil {
			return err
		}
		context.DockerContext = dockerContext
	}

	return nil
}

func (context Context) ConfigPath() string {
	return filepath.Join(context.RootDirectory, "kubedev.yml")
}

// func (context Context) AbsDockerContext() string {
// 	path, _ := filepath.Abs(filepath.Join(context.RootDirectory, context.DockerContext))
// 	return path
// }

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
		RootDirectory: context.RootDirectory,
		FullName:      context.FullName(),
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

func (context Context) YamlFormat() (string, error) {
	output, err := yaml.Marshal(context)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (context Context) YamlPrint() {
	output, err := context.YamlFormat()
	if err == nil {
		fmt.Println(output)
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
}
