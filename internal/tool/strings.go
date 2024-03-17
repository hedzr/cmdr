package tool

import "strings"

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
func SplitCommandString(s string, quoteChars ...rune) []string {
	var qc rune = '"'
	var m = map[rune]bool{qc: true}
	for _, q := range quoteChars {
		qc, m[q] = q, true //nolint:ineffassign
	}

	quoted, ch := false, rune(0)

	a := strings.FieldsFunc(s, func(r rune) bool {
		if ch == 0 {
			if _, ok := m[r]; ok {
				quoted, ch = !quoted, r
			}
		} else if ch == r {
			quoted, ch = !quoted, r
		}
		return !quoted && r == ' '
	})

	var b []string
	for _, s := range a {
		b = append(b, trimQuotes(s))
	}

	return b
}
