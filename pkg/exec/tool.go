package exec

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"unicode/utf8"
)

// LookPath searches for an executable named file in the
// directories named by the PATH environment variable.
// If file contains a slash, it is tried directly and the PATH is not consulted.
// The result may be an absolute path or a path relative to the current directory.
func LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func toStringSimple(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

// StripLeftTabs removes the padding tabs at left margin.
// The least tab chars will be erased at left side of lines, and the
// tab chars beyond the least at left will be kept.
func StripLeftTabs(s string) string { return stripLeftTabs(s) }

func stripLeftTabs(s string) string {
	var lines []string
	var tabs int = 1000
	var emptyLines []int
	var sb strings.Builder
	var line int
	var noLastLF bool = !strings.HasSuffix(s, "\n")

	scanner := bufio.NewScanner(bufio.NewReader(strings.NewReader(s)))
	for scanner.Scan() {
		str := scanner.Text()
		i, n, allTabs := 0, len(str), true
		for ; i < n; i++ {
			if str[i] != '\t' {
				allTabs = false
				if tabs > i && i > 0 {
					tabs = i
					break
				}
			}
		}
		if i == n && allTabs {
			emptyLines = append(emptyLines, line)
		}
		lines = append(lines, str)
		line++
	}

	pad := strings.Repeat("\t", tabs)
	for i, str := range lines {
		if strings.HasPrefix(str, pad) {
			sb.WriteString(str[tabs:])
		} else if inIntSlice(i, emptyLines) {
		} else {
			sb.WriteString(str)
		}
		if noLastLF && i == len(lines)-1 {
			break
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func inIntSlice(i int, slice []int) bool {
	for _, n := range slice {
		if n == i {
			return true
		}
	}
	return false
}

// LeftPad inserts spaces at beginning of each line
func LeftPad(s string, pad int) string { return leftPad(s, pad) }

func leftPad(s string, pad int) string {
	if pad <= 0 {
		return s
	}

	var sb strings.Builder
	padStr := strings.Repeat(" ", pad)
	scanner := bufio.NewScanner(bufio.NewReader(strings.NewReader(s)))
	for scanner.Scan() {
		sb.WriteString(padStr)
		sb.WriteString(scanner.Text())
		sb.WriteRune('\n')
	}
	return sb.String()
}

// StripHtmlTags aggressively strips HTML tags from a string.
// It will only keep anything between `>` and `<`.
func StripHtmlTags(s string) string { return stripHtmlTags(s) }

const (
	htmlTagStart = 60 // Unicode `<`
	htmlTagEnd   = 62 // Unicode `>`
)

// Aggressively strips HTML tags from a string.
// It will only keep anything between `>` and `<`.
func stripHtmlTags(s string) string {
	// Setup a string builder and allocate enough memory for the new string.
	var builder strings.Builder
	builder.Grow(len(s) + utf8.UTFMax)

	in := false // True if we are inside an HTML tag.
	start := 0  // The index of the previous start tag character `<`
	end := 0    // The index of the previous end tag character `>`

	for i, c := range s {
		// If this is the last character and we are not in an HTML tag, save it.
		if (i+1) == len(s) && end >= start {
			builder.WriteString(s[end:])
		}

		// Keep going if the character is not `<` or `>`
		if c != htmlTagStart && c != htmlTagEnd {
			continue
		}

		if c == htmlTagStart {
			// Only update the start if we are not in a tag.
			// This make sure we strip out `<<br>` not just `<br>`
			if !in {
				start = i
			}
			in = true

			// Write the valid string between the close and start of the two tags.
			builder.WriteString(s[end:start])
			continue
		}
		// else c == htmlTagEnd
		in = false
		end = i + 1
	}
	s = builder.String()
	return s
}

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
