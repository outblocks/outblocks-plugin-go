package util

import (
	"os"
	"os/exec"
	"runtime"
)

func NewCmdAsUser(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", command)
	}

	shell, ok := os.LookupEnv("SHELL")
	if !ok {
		shell = "sh"
	}

	uidStr, ok := os.LookupEnv("SUDO_UID")
	if ok {
		return exec.Command("sudo", "-E", "-u", "#"+uidStr, shell, "-c", command)
	}

	return exec.Command(shell, "-c", command)
}