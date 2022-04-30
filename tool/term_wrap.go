//go:build (darwin || dragonfly || freebsd || linux || netbsd || openbsd || aix || arm_linux || solaris) && !nacl && !plan9
// +build darwin dragonfly freebsd linux netbsd openbsd aix arm_linux solaris
// +build !nacl
// +build !plan9

// Copyright © 2020 Hedzr Yeh.

package tool

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/crypto/ssh/terminal"
)

// NOTE:
//   SA1019: package golang.org/x/crypto/ssh/terminal is deprecated: this package moved to golang.org/x/term.
// Here we keep old reference for backward-compatibility to go1.11 (even lower)

//

// ReadPassword reads the password from stdin with safe protection
func ReadPassword() (text string, err error) {
	var bytePassword []byte
	if bytePassword, err = terminal.ReadPassword(syscall.Stdin); err == nil {
		fmt.Println() // it's necessary to add a new line after user's input
		text = string(bytePassword)
	} else {
		fmt.Println() // it's necessary to add a new line after user's input
	}
	return
}

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) {
	var sz struct {
		rows, cols, xPixels, yPixels uint16
	}
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, //nolint:dogsled //like it
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&sz)))
	cols, rows = int(sz.cols), int(sz.rows)
	return
}
