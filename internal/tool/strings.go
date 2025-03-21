package tool

import (
	"regexp"

	"github.com/hedzr/is/exec"
)

// StripQuotes strips first and last quote char (double quote or single quote).
func StripQuotes(s string) string { return trimQuotes(s) }

// TrimQuotes strips first and last quote char (double quote or single quote).
func TrimQuotes(s string) string { return trimQuotes(s) }

// func trimQuotes(s string) string {
// 	if len(s) >= 2 {
// 		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
// 			return s[1 : len(s)-1]
// 		}
// 	}
// 	return s
// }

func trimQuotes(s string) string {
	switch {
	case s[0] == '\'':
		if s[len(s)-1] == '\'' {
			return s[1 : len(s)-1]
		}
		return s[1:]
	case s[0] == '"':
		if s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
		return s[1:]
	case s[len(s)-1] == '\'':
		return s[0 : len(s)-1]
	case s[len(s)-1] == '"':
		return s[0 : len(s)-1]
	}
	return s
}

// SplitCommandString allows split command-line by quote
// characters (default is double-quote).
//
// In: `bash -c 'echo hello world!'`
// Out: []string{ "bash", "-c", "echo hello world!"}
//
// For example:
//
//	in := `bash -c 'echo hello world!'`
//	out := SplitCommandString(in, '\'', '"')
//	println(out)   // will got: []string{ "bash", "-c", "echo hello world!"}
func SplitCommandString(s string, quoteChars ...rune) []string {
	return exec.SplitCommandString(s, quoteChars...)
}

// StripOrderPrefix strips the prefix string fragment for sorting order.
// see also: Command.Group, Flag.Group, ...
// An order prefix is a dotted string with multiple alphabet and digit. Such as:
// "zzzz.", "0001.", "700.", "A1." ...
func StripOrderPrefix(s string) string {
	a := xre.FindStringSubmatch(s)
	return a[2]
	// if xre.MatchString(s) {
	//	s = s[strings.Index(s, ".")+1:]
	// }
	// return s
}

// HasOrderPrefix tests whether an order prefix is present or not.
// An order prefix is a dotted string with multiple alphabet and digit. Such as:
// "zzzz.", "0001.", "700.", "A1." ...
func HasOrderPrefix(s string) bool {
	return xre.MatchString(s)
}

var xre = regexp.MustCompile(`^([0-9A-Za-z]+[.])?(.+)$`)
