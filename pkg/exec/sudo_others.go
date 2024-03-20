//go:build !plan9 && !nacl
// +build !plan9,!nacl

// Copyright © 2020 Hedzr Yeh.

package exec

import (
	"os/exec"
	"syscall"
)

// IsExitError checks the error object
func IsExitError(err error) (int, bool) {
	if ee, ok := err.(*exec.ExitError); ok {
		if status, ok := ee.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus(), true
		}
	}

	return 0, false
}
