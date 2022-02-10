/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"fmt"
	"github.com/hedzr/cmdr/tool"
	"golang.org/x/net/html"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// some refs:
// - github.com/labstack/gommon/color
//

const (
	defaultTimestampFormat = time.RFC3339

	// https://en.wikipedia.org/wiki/ANSI_escape_code
	// https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97

	// FgBlack terminal color code
	FgBlack = 30
	// FgRed terminal color code
	FgRed = 31
	// FgGreen terminal color code
	FgGreen = 32
	// FgYellow terminal color code
	FgYellow = 33
	// FgBlue terminal color code
	FgBlue = 34
	// FgMagenta terminal color code
	FgMagenta = 35
	// FgCyan terminal color code
	FgCyan = 36
	// FgLightGray terminal color code
	FgLightGray = 37
	// FgDarkGray terminal color code
	FgDarkGray = 90
	// FgLightRed terminal color code
	FgLightRed = 91
	// FgLightGreen terminal color code
	FgLightGreen = 92
	// FgLightYellow terminal color code
	FgLightYellow = 93
	// FgLightBlue terminal color code
	FgLightBlue = 94
	// FgLightMagenta terminal color code
	FgLightMagenta = 95
	// FgLightCyan terminal color code
	FgLightCyan = 96
	// FgWhite terminal color code
	FgWhite = 97

	// BgNormal terminal color code
	BgNormal = 0
	// BgBoldOrBright terminal color code
	BgBoldOrBright = 1
	// BgDim terminal color code
	BgDim = 2
	// BgItalic terminal color code
	BgItalic = 3
	// BgUnderline terminal color code
	BgUnderline = 4
	// BgUlink terminal color code
	BgUlink = 5
	// BgInverse _
	BgInverse = 7
	// BgHidden terminal color code
	BgHidden = 8
	// BgStrikeout terminal color code
	BgStrikeout = 9

	// DarkColor terminal color code
	DarkColor = FgLightGray
)

var (
	// onceColorPrintTranslator sync.Once
	cpt   colorPrintTranslator
	cptNC = colorPrintTranslator{noColorMode: true}
)

type colorPrintTranslator struct {
	noColorMode bool // strip color code simply
}

func (c *colorPrintTranslator) Translate(s string, initialFg int) string {
	return c.TranslateTo(s, initialFg)
}

func (c *colorPrintTranslator) resetColors(sb strings.Builder, states []int) func() {
	return func() {
		var st string
		st = "\x1b[0m"
		sb.WriteString(st)
		if len(states) > 0 {
			st = fmt.Sprintf("\x1b[%dm", states[len(states)-1])
			sb.WriteString(st)
		}
	}
}

func (c *colorPrintTranslator) colorize(sb strings.Builder, states []int, walker func(node *html.Node, level int)) func(node *html.Node, clr int, representation string, level int) {
	return func(node *html.Node, clr int, representation string, level int) {
		if representation != "" {
			sb.WriteString(fmt.Sprintf("\x1b[%sm", representation))
		} else {
			sb.WriteString(fmt.Sprintf("\x1b[%dm", clr))
		}
		states = append(states, clr)
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walker(child, level+1)
		}
		states = states[0 : len(states)-1]
		c.resetColors(sb, states)()
	}
}

func (c *colorPrintTranslator) TranslateTo(s string, initialState int) string {
	if c.noColorMode {
		return c._ss(s)
	}

	node, err := html.Parse(bufio.NewReader(strings.NewReader(s)))
	if err != nil {
		return c._sz(s)
	}

	return c.translateTo(node, s, initialState)
}

func (c *colorPrintTranslator) translateTo(root *html.Node, s string, initialState int) string {
	var states = []int{initialState}
	var sb strings.Builder
	var walker func(node *html.Node, level int)
	colorize := c.colorize(sb, states, walker)
	nilfn := func(node *html.Node, level int) {}
	colorizeIt := func(clr int) func(node *html.Node, level int) {
		return func(node *html.Node, level int) {
			colorize(node, clr, "", level)
		}
	}
	m := map[string]func(node *html.Node, level int){
		"html": nilfn, "head": nilfn, "body": nilfn,
		"b": colorizeIt(BgBoldOrBright), "strong": colorizeIt(BgBoldOrBright), "em": colorizeIt(BgBoldOrBright),
		"i": colorizeIt(BgItalic), "cite": colorizeIt(BgItalic),
		"u":    colorizeIt(BgUnderline),
		"mark": colorizeIt(BgInverse),
		"del":  colorizeIt(BgStrikeout),
	}
	walker = func(node *html.Node, level int) {
		switch node.Type {
		case html.DocumentNode, html.DoctypeNode, html.CommentNode:
		case html.ErrorNode:
		case html.ElementNode:
			if fn, ok := m[node.Data]; ok {
				fn(node, level)
				return
			}

			switch node.Data {
			//case "html", "head", "body":
			//	// do nothing
			//case "b", "strong", "em":
			//	colorize(node, BgBoldOrBright, "", level)
			//	return
			//case "i", "cite":
			//	colorize(node, BgItalic, "", level)
			//	return
			//case "u":
			//	colorize(node, BgUnderline, "", level)
			//	return
			//case "mark":
			//	colorize(node, BgInverse, "", level)
			//	return
			//case "del":
			//	colorize(node, BgStrikeout, "", level)
			//	return
			case "font":
				for _, a := range node.Attr {
					if a.Key == "color" {
						clr := c.toColorInt(a.Val)
						colorize(node, clr, "", level)
						return
					}
				}
			case "kbd", "code":
				colorize(node, 51, "51;1", level)
				return
			default:
				//Logger.Debugf("%v, %v, lvl #%d\n", node.Type, node.Data, level)
				//sb.WriteString(node.Data)
			}
		case html.TextNode:
			//Logger.Debugf("%v, %v, lvl #%d\n", node.Type, node.Data, level)
			sb.WriteString(node.Data)
			return
		default:
			//sb.WriteString(node.Data)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walker(child, level+1)
		}
	}
	walker(root, 0)
	return sb.String()
}

func (c *colorPrintTranslator) _sz(s string) string {
	return s
}

func (c *colorPrintTranslator) _ss(s string) string {
	if tool.IsTtyEscaped(s) {
		clean := tool.StripEscapes(s)
		return c.stripHTMLTags(clean)
	}
	return c.stripHTMLTags(s)
}

var onceoptCM sync.Once
var cptCM map[string]int

func (c *colorPrintTranslator) toColorInt(s string) int {
	onceoptCM.Do(func() {
		cptCM = map[string]int{
			"black":     FgBlack,
			"red":       FgRed,
			"green":     FgGreen,
			"yellow":    FgYellow,
			"blue":      FgBlue,
			"magenta":   FgMagenta,
			"cyan":      FgCyan,
			"lightgray": FgLightGray, "light-gray": FgLightGray,
			"darkgray": FgDarkGray, "dark-gray": FgDarkGray,
			"lightred": FgLightRed, "light-red": FgLightRed,
			"lightgreen": FgLightGreen, "light-green": FgLightGreen,
			"lightyellow": FgLightYellow, "light-yellow": FgLightYellow,
			"lightblue": FgLightBlue, "light-blue": FgLightBlue,
			"lightmagenta": FgLightMagenta, "light-magenta": FgLightMagenta,
			"lightcyan": FgLightCyan, "light-cyan": FgLightCyan,
			"white": FgWhite,
		}
	})
	if i, ok := cptCM[strings.ToLower(s)]; ok {
		return i
	}
	return 0
}

func (c *colorPrintTranslator) stripLeftTabs(s string) string {
	r := c.stripLeftTabsOnly(s)
	return c.Translate(r, 0)
}

func (c *colorPrintTranslator) stripLeftTabsOnly(s string) string {

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

// StripLeftTabs strips the least left side tab chars from lines.
// StripLeftTabs strips html tags too.
func StripLeftTabs(s string) string { return cptNC.stripLeftTabs(s) }

// StripLeftTabsOnly strips the least left side tab chars from lines.
func StripLeftTabsOnly(s string) string { return cptNC.stripLeftTabsOnly(s) }

// StripHTMLTags aggressively strips HTML tags from a string.
// It will only keep anything between `>` and `<`.
func StripHTMLTags(s string) string { return cptNC.stripHTMLTags(s) }

const (
	htmlTagStart = 60 // Unicode `<`
	htmlTagEnd   = 62 // Unicode `>`
)

// Aggressively strips HTML tags from a string.
// It will only keep anything between `>` and `<`.
func (c *colorPrintTranslator) stripHTMLTags(s string) string {
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
