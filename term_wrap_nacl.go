// +build nacl

/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

func readPassword() (text string, err error) {
	return randomStringPure(9), nil
}
