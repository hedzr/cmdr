/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import (
	"errors"
	"github.com/hedzr/cmdr"
	"log"
	"os"
	"os/exec"
)

func runInstaller(cmd *cmdr.Command, args []string) (err error) {
	if !isRoot() {
		log.Fatal("This program must be run as root! (sudo)")
	}

	return
}

func runUninstaller(cmd *cmdr.Command, args []string) (err error) {
	if !isRoot() {
		log.Fatal("This program must be run as root! (sudo)")
	}

	return
}

func isRoot() bool {
	return os.Getuid() == 0
}

// Exec executes a command setting both standard input, output and error.
func Exec(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

// ExecSudo executes a command under "sudo".
func ExecSudo(cmd string, args ...string) error {
	return Exec("sudo", append([]string{cmd}, args...)...)
}

var ErrNoRoot = errors.New("MUST have administrator privileges")
