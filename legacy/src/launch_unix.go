//go:build !windows

package main

import (
	"os/exec"
)

// setupProcessAttributes configures the command to run detached on Unix systems
func setupProcessAttributes(cmd *exec.Cmd) {
	// On Unix systems, don't redirect output to prevent blocking
	cmd.Stdout = nil
	cmd.Stderr = nil
}
