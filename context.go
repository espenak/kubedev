package main

import (
  "text/template"
  "path/filepath"
  "fmt"
  "os"
)

type Context struct {
  RootDirectory string
  Verbose bool
}

func NewContext (rootDirectory string, verbose bool) *Context {
  context := Context{
    RootDirectory: rootDirectory,
    Verbose: verbose,
  }
  context.LoadConfig()
  return &context
}

func (context *Context) LoadConfig() {
  fmt.Println("TODO: Load config")
}

func (context Context) TemplatesDirectory() string {
  return filepath.Join(context.RootDirectory, "templates")
}

func (context Context) BuiltTemplatesDirectory() string {
  return filepath.Join(context.RootDirectory, "_kubedev_built_templates")
}

func (context Context) ParseTemplates () (*template.Template, error) {
  return template.ParseGlob(
    filepath.Join(context.TemplatesDirectory(), "*"))
}

func (context Context) BuildTemplate(parsedTemplates *template.Template, templateName string) error {
  outPath := filepath.Join(context.BuiltTemplatesDirectory(), templateName)
  outFile, err := os.Create(outPath)
  if err != nil {
    return err
  }
  defer outFile.Close()
  parsedTemplates.ExecuteTemplate(outFile, templateName, "Hello world")
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
