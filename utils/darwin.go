//go:build darwin

package utils

import (
	"os"
	"os/exec"
)

func IsWindows() bool {
	return false
}

func FrokProcess(name string, arg ...string) (int, error) {
	// Fork the process
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}
