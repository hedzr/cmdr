// +build darwin dragonfly freebsd linux netbsd openbsd windows aix arm_linux plan9 solaris
// +build !nacl

// Copyright © 2020 Hedzr Yeh.

package tool

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

// ReadPassword reads the password from stdin with safe protection
func ReadPassword() (text string, err error) {
	var bytePassword []byte
	if bytePassword, err = terminal.ReadPassword(int(syscall.Stdin)); err == nil {
		fmt.Println() // it's necessary to add a new line after user's input
		text = string(bytePassword)
	} else {
		fmt.Println() // it's necessary to add a new line after user's input
	}
	return
}
