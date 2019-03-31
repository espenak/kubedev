package main

import (
	"io/ioutil"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type UserSettings struct {
	KubectlExecutable string            `yaml:"kubeCtlExecutable"`
	Config            map[string]string `yaml:"config"`
}

func (userSettings *UserSettings) readFromFile(settingsFilePath string) error {
	if content, err := ioutil.ReadFile(settingsFilePath); err != nil {
		return err
	} else {
		if err := yaml.Unmarshal(content, &userSettings); err != nil {
			return err
		}
	}
	userSettings.clean()
	return nil
}

func (userSettings *UserSettings) userDirectoryFilePath() (string, error) {
	if u, err := user.Current(); err != nil {
		return "", err
	} else {
		return filepath.Join(u.HomeDir, ".kubedev.usersettings.yml"), nil
	}
}

func (userSettings *UserSettings) readFromUserDirectory() error {
	if settingsFilePath, err := userSettings.userDirectoryFilePath(); err != nil {
		return err
	} else {
		return userSettings.readFromFile(settingsFilePath)
	}
}

func (userSettings *UserSettings) clean() {
	if userSettings.KubectlExecutable == "" {
		userSettings.KubectlExecutable = "kubectl"
	}
	if userSettings.Config == nil {
		userSettings.Config = make(map[string]string)
	}
}

func (userSettings UserSettings) ToYaml() ([]byte, error) {
	if output, err := yaml.Marshal(userSettings); err != nil {
		return nil, err
	} else {
		return output, nil
	}
}

func (userSettings UserSettings) ToYamlString() string {
	yamlBytes, _ := userSettings.ToYaml()
	return string(yamlBytes)
}

func (userSettings *UserSettings) Update(otherUserSettings UserSettings) {
	userSettings.KubectlExecutable = otherUserSettings.KubectlExecutable
	for configKey, configValue := range otherUserSettings.Config {
		userSettings.Config[configKey] = configValue
	}
}
