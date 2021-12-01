// Copyright © 2020 Hedzr Yeh.

package tool

import (
	"regexp"
	"strings"
)

// StripQuotes strips single or double quotes around a string.
func StripQuotes(s string) string {
	return trimQuotes(s)
}

func trimQuotes(s string) string {
	if s[0] == '\'' {
		if s[len(s)-1] == '\'' {
			return s[1 : len(s)-1]
		}
		return s[1:]

	} else if s[0] == '"' {
		if s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
		return s[1:]

	} else if s[len(s)-1] == '\'' {
		return s[0 : len(s)-1]
	} else if s[len(s)-1] == '"' {
		return s[0 : len(s)-1]
	}
	return s
}

// func eraseMultiWSs(s string) string {
// 	return reSimp.ReplaceAllString(s, " ")
// }

// EraseAnyWSs eats any whitespaces inside the giving string s.
func EraseAnyWSs(s string) string {
	return reSimpSimp.ReplaceAllString(s, "")
}

// var reSimp = regexp.MustCompile(`[ \t][ \t]+`)
var reSimpSimp = regexp.MustCompile(`[ \t]+`)

// EscapeCompletionTitle escapes ';' character for zsh completion system
func EscapeCompletionTitle(title string) string {
	ret := strings.ReplaceAll(title, "'", "\"")
	ret = strings.ReplaceAll(ret, ":", "\\:")
	return ret
}
