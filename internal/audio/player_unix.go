//go:build !windows

package audio

import (
	"os/exec"
	"syscall"
)

var currentCmd *exec.Cmd

func stopPlayback() {
	if currentCmd != nil && currentCmd.Process != nil {
		// Kill the process group
		syscall.Kill(-currentCmd.Process.Pid, syscall.SIGKILL)
		currentCmd = nil
	}
}
