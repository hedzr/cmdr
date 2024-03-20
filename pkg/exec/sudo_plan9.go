//go:build plan9
// +build plan9

// Copyright © 2020 Hedzr Yeh.

package exec

import (
	"os/exec"
)

// IsExitError checks the error object
func IsExitError(err error) (int, bool) {
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode(), true
	}

	return 0, false
}
