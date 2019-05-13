/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import "regexp"

func simp(s string) string {
	return reSimp.ReplaceAllString(s, " ")
}

func simpsimp(s string) string {
	return reSimpSimp.ReplaceAllString(s, "")
}

var reSimp = regexp.MustCompile(`[ \t][ \t]+`)
var reSimpSimp = regexp.MustCompile(`[ \t]+`)
