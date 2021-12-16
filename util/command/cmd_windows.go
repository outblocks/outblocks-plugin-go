//go:build windows

package command

import (
	"os"
	"os/exec"
	"syscall"
)

func CmdSignal(cmd *exec.Cmd, signal os.Signal) error {
	return cmd.Process.Signal(signal)
}

func defaultSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
