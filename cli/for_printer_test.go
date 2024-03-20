package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/hedzr/is/term/color"
)

type helpPrinter struct {
	color.Translator
	debugMatches bool     //nolint:unused
	w            *workerS //nolint:unused
}

func (s *helpPrinter) Print(ctx any, lastCmd *Command) { //nolint:unparam,revive
	if s.Translator == nil {
		s.Translator = color.GetCPT()
	}

	const currentDescColor = color.FgDarkGray
	type HelpWriter interface {
		io.Writer
		io.StringWriter
	}

	var sb strings.Builder
	var wr HelpWriter = os.Stdout
	// if s.w != nil && s.w.wrHelpScreen != nil {
	// 	wr = s.w.wrHelpScreen
	// }

	lastCmd.WalkEverything(func(cc, pp *Command, ff *Flag, cmdIndex, flgIndex, level int) {
		switch {
		case ff == nil && level > 0:
			_, _ = sb.WriteString(strings.Repeat("  ", level))
			w := 56 - level*2
			ttl := cc.GetTitleNames()
			if ttl == "" {
				ttl = "(empty)"
			}

			left, right := fmt.Sprintf("%-"+strconv.Itoa(w)+"s", ttl), fmt.Sprintf("%s\n", cc.Desc())
			_, _ = sb.WriteString(left)
			// s.DimFast(&sb, s.Translate(right, color.BgNormal))
			s.ColoredFast(&sb, currentDescColor, s.Translate(right, currentDescColor))

		case ff != nil:
			_, _ = sb.WriteString(strings.Repeat("  ", level+1))
			// ttl := strings.Join(ff.GetTitleZshFlagNamesArray(), ",")
			ttl := ff.GetTitleFlagNamesBy(",")
			w := 56 - (level+1)*2 // - len(ttl)

			left, right := fmt.Sprintf("%-"+strconv.Itoa(w)+"s", ttl), fmt.Sprintf("%s\n", ff.Desc())
			s.ColoredFast(&sb, color.FgGreen, left)
			// s.DimFast(&sb, s.Translate(right, color.BgNormal))
			s.ColoredFast(&sb, currentDescColor, s.Translate(right, currentDescColor))
			// sb.WriteString(s.Translate(right, color.BgDefault))
		}
	})
	_, _ = wr.WriteString(sb.String())
	_, _ = wr.WriteString("\n")
}
