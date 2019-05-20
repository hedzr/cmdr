/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import "regexp"

func eraseMultiWSs(s string) string {
	return reSimp.ReplaceAllString(s, " ")
}

func eraseAnyWSs(s string) string {
	return reSimpSimp.ReplaceAllString(s, "")
}

var reSimp = regexp.MustCompile(`[ \t][ \t]+`)
var reSimpSimp = regexp.MustCompile(`[ \t]+`)
