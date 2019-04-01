package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type DockerDirectory struct {
	context       Context
	DirectoryPath string
}

func NewDockerDirectory(context Context, directoryPath string) (*DockerDirectory, error) {
	dockerFilePath := filepath.Join(directoryPath, "Dockerfile")
	if _, err := os.Stat(dockerFilePath); os.IsNotExist(err) {
		return nil, err
	}
	dockerDirectory := DockerDirectory{context, directoryPath}
	return &dockerDirectory, nil
}

func (dockerDirectory DockerDirectory) Name() string {
	return filepath.Base(dockerDirectory.DirectoryPath)
}

func (dockerDirectory DockerDirectory) FilePath() string {
	return filepath.Join(
		dockerDirectory.DirectoryPath,
		"Dockerfile")
}

func (dockerDirectory DockerDirectory) ImageName() string {
	return fmt.Sprintf("%s/%s", dockerDirectory.context.FullName(), dockerDirectory.Name())
}

func makeDockerBuildVariableName(prefix, suffix string) string {
	re := regexp.MustCompile(`[^a-zA-Z_]`)
	fullVariableName := prefix + "_" + suffix
	return strings.ToUpper(re.ReplaceAllString(fullVariableName, "_"))
}

func (dockerDirectory DockerDirectory) BuildCommandArgs() []string {
	commandArgs := []string{"build", dockerDirectory.context.DockerContext, "-f", dockerDirectory.FilePath(), "-t", dockerDirectory.ImageName()}
	for key, value := range dockerDirectory.context.Vars {
		buildArg := fmt.Sprintf("%v=%v", makeDockerBuildVariableName("vars", key), value)
		commandArgs = append(commandArgs, "--build-arg")
		commandArgs = append(commandArgs, buildArg)
	}
	for key, value := range dockerDirectory.context.UserConfig {
		buildArg := fmt.Sprintf("%v=%v", makeDockerBuildVariableName("userConfig", key), value)
		commandArgs = append(commandArgs, "--build-arg")
		commandArgs = append(commandArgs, buildArg)
	}
	return commandArgs
}

func (dockerDirectory DockerDirectory) Build() error {
	buildCommandArgs := dockerDirectory.BuildCommandArgs()
	log.Printf("Building docker image %v", dockerDirectory.ImageName())
	log.Printf("... running: \"docker\" executable with args %v", buildCommandArgs)
	cmd := exec.Command("docker", buildCommandArgs...)
	if dockerDirectory.context.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err == nil {
		log.Printf("... docker image %v built successfully.", dockerDirectory.ImageName())
		return err
	} else {
		log.Printf("... docker image %v build FAILED.", dockerDirectory.ImageName())
		log.Println(err)
	}
	return nil
}
