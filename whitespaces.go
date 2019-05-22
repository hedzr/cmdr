/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import "regexp"

func trimQuotes(s string) string {
	if s[0] == '\'' {
		if s[len(s)-1] == '\'' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func eraseMultiWSs(s string) string {
	return reSimp.ReplaceAllString(s, " ")
}

func eraseAnyWSs(s string) string {
	return reSimpSimp.ReplaceAllString(s, "")
}

var reSimp = regexp.MustCompile(`[ \t][ \t]+`)
var reSimpSimp = regexp.MustCompile(`[ \t]+`)
