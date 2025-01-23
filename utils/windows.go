//go:build windows

package utils

import (
	"os"
	"os/exec"
	"syscall"
)

func IsWindows() bool {
	return true
}

func FrokProcess(name string, arg ...string) (int, error) {
	// Fork the process
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	err := cmd.Start()
	if err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}
