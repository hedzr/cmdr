// +build linux,windows,freebsd,darwin,aix,arm_linux,netbsd,openbsd,plan9,solaris

/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func readPassword() (text string, err error) {
	var bytePassword []byte
	if bytePassword, err = terminal.ReadPassword(int(syscall.Stdin)); err == nil {
		fmt.Println() // it's necessary to add a new line after user's input
		text = string(bytePassword)
	} else {
		fmt.Println() // it's necessary to add a new line after user's input
	}
	return
}
