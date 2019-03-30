package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/atrox/homedir"
	"gopkg.in/yaml.v2"
)

type Context struct {
	RootDirectory string            `yaml:"rootDirectory"`
	ApiVersion    string            `yaml:"apiVersion"`
	Name          string            `yaml:"name"`
	DockerContext string            `yaml:"dockerContext"`
	Verbose       bool              `yaml:"verbose"`
	Paths         map[string]string `yaml:"paths"`
	Vars          map[string]string `yaml:"vars"`
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

	if context.Paths == nil {
		context.Paths = make(map[string]string)
	} else {
		for pathKey, rawPath := range context.Paths {
			cleanedPath, err := homedir.Expand(rawPath)
			if err != nil {
				return err
			}
			if !filepath.IsAbs(cleanedPath) {
				cleanedPath, err = filepath.Abs(filepath.Join(context.RootDirectory, cleanedPath))
				if err != nil {
					return err
				}
			}
			context.Paths[pathKey] = cleanedPath
		}
	}

	if context.Vars == nil {
		context.Vars = make(map[string]string)
	}
	// for k, v := range m {
	//   fmt.Println("k:", k, "v:", v)
	// }

	return nil
}

func (context Context) ConfigPath() string {
	return filepath.Join(context.RootDirectory, "kubedev.yml")
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
	return &TemplateVariables{context}
}

func (context Context) BuildTemplate(parsedTemplates *template.Template, templateName string) error {
	outPath := filepath.Join(context.BuiltTemplatesDirectory(), templateName)
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	compileError := parsedTemplates.ExecuteTemplate(outFile, templateName, context.TemplateVariables())
	if compileError != nil {
		return err
	}
	outFile.Sync()
	return err
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

func (context Context) GetAllDockerDirectories() ([]DockerDirectory, error) {
	var dockerDirectories []DockerDirectory
	dockerImagesDirectory := filepath.Join(context.RootDirectory, "dockerimages")
	fileInfo, readDirErr := ioutil.ReadDir(dockerImagesDirectory)
	if readDirErr != nil {
		return nil, readDirErr
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			dockerImageDirectory := filepath.Join(dockerImagesDirectory, file.Name())
			dockerDirectory, err := NewDockerDirectory(context, dockerImageDirectory)
			if err == nil {
				dockerDirectories = append(dockerDirectories, *dockerDirectory)
			} else {
				log.Printf("WARNING: %v", err)
			}
		}
	}
	return dockerDirectories, nil
}

func (context Context) BuildAllDockerImages() error {
	dockerDirectories, err := context.GetAllDockerDirectories()
	if err != nil {
		return err
	}
	for _, dockerDirectory := range dockerDirectories {
		err = dockerDirectory.Build()
		if err != nil {
			return err
		}
	}
	return nil
}
