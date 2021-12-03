// Copyright © 2020 Hedzr Yeh.

package tool

import (
	"reflect"
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

// ReverseAnySlice reverse any slice/array
func ReverseAnySlice(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// ReverseStringSlice reverse a string slice
func ReverseStringSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
