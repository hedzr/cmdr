// Copyright © 2020 Hedzr Yeh.

package cmdr_test

import (
	"bufio"
	"github.com/hedzr/cmdr"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestGetTextPieces(t *testing.T) {
	for _, tt := range []string{
		`Options:\x1b[2m\1b[37mMisc\1b]0m
  [[2m[37mMisc[0m]
      --config=[Locations of config files]  [0m[90mload config files from where you specified[3m[90m (default [Locations of config files]=)[0m
  -q, --quiet                               [0m[90mNo more screen output.[3m[90m [env: QUITE] (default=true)[0m
  -v, --verbose                             [0m[90mShow this help screen[3m[90m [env: VERBOSE] (default=false)[0m
[2m[37m
`,
	} {
		_ = cmdr.GetTextPiecesForTest(tt, 0, 1000)
	}
}

func TestParseHtml(t *testing.T) {
	source := `
	load <code>config</code> files from where you specified
	`

	node, err := html.Parse(bufio.NewReader(strings.NewReader(source)))
	if err != nil {
		t.Error(err)
	}

	var sb strings.Builder
	var walker func(node *html.Node, level int)
	walker = func(node *html.Node, level int) {
		switch node.Type {
		case html.DocumentNode, html.DoctypeNode, html.CommentNode:
		case html.ErrorNode:
		case html.ElementNode:
			switch node.Data {
			case "html", "head", "body":
			case "code":
				sb.WriteString("\x1b[1m")
				for child := node.FirstChild; child != nil; child = child.NextSibling {
					walker(child, level+1)
				}
				sb.WriteString("\x1b[0m")
				return
			default:
				// Logger.Debugf("%v, %v, lvl #%d\n", node.Type, node.Data, level)
				// sb.WriteString(node.Data)
			}
		case html.TextNode:
			// Logger.Debugf("%v, %v, lvl #%d\n", node.Type, node.Data, level)
			sb.WriteString(node.Data)
			return
		default:
			// sb.WriteString(node.Data)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walker(child, level+1)
		}
	}

	walker(node, 0)

	t.Logf("%v", sb.String())
}

func TestCPT(t *testing.T) {
	source := `
	load <code>config</code> files from where <font color="green"><b>you</b></font> specified
	<del>scan</del> <u>folder</u> and save <i>result</i> to <code>bgo.yml</code>, as <mark>project settings</mark>
	`

	str := cmdr.Cpt().Translate(source, 4)
	t.Logf("\x1b[4m%v\x1b[0m", str)
}

func TestCPTNC(t *testing.T) {
	source := `
	load <code>config</code> files from where <font color="green"><b>you</b></font> specified
	<del>scan</del> <u>folder</u> and save <i>result</i> to <code>bgo.yml</code>, as <mark>project settings</mark>
	`

	str := cmdr.CptNC().Translate(source, 0)
	t.Logf("%v", str)
}

func TestStripLeftTabs(t *testing.T) {
	source := `
		// load <code>config</code> files from where you specified
			if tool.IsTtyEscaped(s) {
		    clean := tool.StripEscapes(s)
			return c.stripHTMLTags(clean)
		}
	`
	expected := `
// load config files from where you specified
	if tool.IsTtyEscaped(s) {
    clean := tool.StripEscapes(s)
	return c.stripHTMLTags(clean)
}
`
	expected2 := `
// load <code>config</code> files from where you specified
	if tool.IsTtyEscaped(s) {
    clean := tool.StripEscapes(s)
	return c.stripHTMLTags(clean)
}
`
	sz := cmdr.StripLeftTabs(source)
	if sz != expected {
		t.Errorf("unexpect result\n%v", sz)
	}

	sz = cmdr.StripLeftTabsOnly(source)
	if sz != expected2 {
		t.Errorf("unexpect result\n%v", sz)
	}
}

func TestStripHtmlTags(t *testing.T) {
	source := `
		// load <code>config</code> files from where you specified
			if tool.IsTtyEscaped(s) {
			clean := tool.StripEscapes(s)
			return c.stripHTMLTags(clean)
		}
		<del>scan</del> <u>folder</u> and save <i>result</i> to <code>bgo.yml</code>, as <mark>project settings</mark>
	`
	expected := `
		// load config files from where you specified
			if tool.IsTtyEscaped(s) {
			clean := tool.StripEscapes(s)
			return c.stripHTMLTags(clean)
		}
		scan folder and save result to bgo.yml, as project settings
	`
	sz := cmdr.StripHTMLTags(source)
	if sz != expected {
		t.Errorf("unexpect result\n%v", sz)
	}
}
