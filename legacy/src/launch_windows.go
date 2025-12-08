//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

// setupProcessAttributes configures the command to run detached on Windows
func setupProcessAttributes(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | 0x00000008, // CREATE_NO_WINDOW
	}
}
