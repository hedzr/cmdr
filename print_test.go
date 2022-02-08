// Copyright © 2020 Hedzr Yeh.

package cmdr

import (
	"bufio"
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
		_ = getTextPiece(tt, 0, 1000)
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

	walker(node, 0)

	t.Logf("%v", sb.String())
}

func TestCPT(t *testing.T) {
	source := `
	load <code>config</code> files from where you specified
	<del>scan</del> <u>folder</u> and save <i>result</i> to <code>bgo.yml</code>, as <mark>project settings</mark>
	`

	str := cpt.Translate(source, 4)
	t.Logf("\x1b[4m%v\x1b[0m", str)
}

func TestCPTNC(t *testing.T) {
	source := `
	load <code>config</code> files from where you specified
	<del>scan</del> <u>folder</u> and save <i>result</i> to <code>bgo.yml</code>, as <mark>project settings</mark>
	`

	str := cptNC.Translate(source, 0)
	t.Logf("%v", str)
}
