// Copyright © 2020 Hedzr Yeh.

// +build go1.12

package cmdr

import "strings"

func replaceAll(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}
