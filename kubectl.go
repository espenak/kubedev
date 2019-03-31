package main

import (
	"os"
	"os/exec"
	"strings"
)

type KubeCtl struct {
	KubectlExectutablePath string
}

func (kubeCtl KubeCtl) GetContexts() ([]string, error) {
	cmd := exec.Command(kubeCtl.KubectlExectutablePath, "config", "view", "-o", "jsonpath={.contexts[*].name}")
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(string(output), " "), nil
}

func (kubeCtl KubeCtl) GetCurrentContext() (string, error) {
	cmd := exec.Command(kubeCtl.KubectlExectutablePath, "config", "current-context")
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
