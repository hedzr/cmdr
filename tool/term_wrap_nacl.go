// +build nacl

// Copyright © 2020 Hedzr Yeh.

package tool

// ReadPassword reads the password from stdin with safe protection
func ReadPassword() (text string, err error) {
	return randomStringPure(9), nil
}
