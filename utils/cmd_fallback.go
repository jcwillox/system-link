//go:build !windows

package utils

import "os/exec"

func makeCmdHidden(cmd *exec.Cmd) {
	return
}
