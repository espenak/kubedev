package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ConfigFilePath string `yaml:"configFilePath"`
	RootDirectory  string `yaml:"rootDirectory"`
	ApiVersion     string `yaml:"apiVersion"`
	Name           string `yaml:"name"`
	DockerContext  string `yaml:"dockerContext"`
}

func LoadConfigFromFile(configFilePath string) (*Config, error) {
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	config := Config{}
	err2 := yaml.Unmarshal(content, &config)
	if err2 != nil {
		return nil, err2
	}

	configFilePath, err3 := filepath.Abs(configFilePath)
	if err3 != nil {
		return nil, err3
	}
	config.ConfigFilePath = configFilePath
	config.RootDirectory = filepath.Dir(config.ConfigFilePath)

	return &config, nil
}

func (config Config) YamlFormat() (string, error) {
	output, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (config Config) YamlPrint() {
	output, err := config.YamlFormat()
	if err == nil {
		fmt.Println(output)
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
}
