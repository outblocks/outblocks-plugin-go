//go:build !windows

package command

import (
	"os"
	"os/exec"
	"syscall"
)

func CmdSignal(cmd *exec.Cmd, signal os.Signal) error {
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		err = syscall.Kill(-pgid, signal.(syscall.Signal)) //nolint:errcheck
		_ = cmd.Process.Release()

		return err
	}

	return cmd.Process.Signal(signal)
}

func defaultSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setpgid: true,
	}
}
