package worker

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/hedzr/evendeep"
	"github.com/hedzr/is"
	"github.com/hedzr/is/states"
	"github.com/hedzr/is/term"
	"github.com/hedzr/is/term/color"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/logz"
	"github.com/hedzr/is/exec"
)

type discardP struct{}

func (*discardP) Write([]byte) (n int, err error) { return }

func (*discardP) WriteString(string) (n int, err error) { return }

type helpPrinter struct {
	color.Translator
	w               *workerS
	debugScreenMode bool
	debugMatches    bool
	treeMode        bool
	asManual        bool
	lastFlagGroup   string
	lastCmdGroup    string
}

const colLeftTabbedWidth = 46

type wHW struct {
	io.Writer
}

func (w *wHW) WriteString(s string) (n int, err error) { return w.Write([]byte(s)) }

func (s *helpPrinter) safeGetWriter() (wr HelpWriter) {
	wr = os.Stdout
	if s.w != nil && s.w.wrHelpScreen != nil {
		wr = s.w.wrHelpScreen
		// logz.Debug("Using helpWriter")
		// fmt.Fprintln(s.w.wrHelpScreen, "Using helpWriter / 1")
	}
	return
}

func (s *helpPrinter) Print(ctx context.Context, pc cli.ParsedState, lastCmd cli.Cmd, args ...any) { //
	wr := s.safeGetWriter()
	if len(args) > 0 {
		if t, ok := args[0].(io.Writer); ok {
			if h, ok1 := t.(HelpWriter); ok1 {
				wr = h
			} else {
				wr = &wHW{t}
			}
		}
	}

	// fmt.Fprintln(wr, "Using helpWriter : ", wr, s.w.wrHelpScreen)

	s.PrintTo(ctx, wr, pc, lastCmd, args...)
}

func (s *helpPrinter) PrintTo(ctx context.Context, wr HelpWriter, pc cli.ParsedState, lastCmd cli.Cmd, args ...any) { //
	if s.debugScreenMode {
		s.PrintDebugScreenTo(ctx, wr, pc, lastCmd)
		return
	}

	if s.Translator == nil {
		logz.Info(`[helpPrinter] getCPT`, "no-color", is.NoColorMode())
		s.Translator = color.GetCPT()
	}

	var sb strings.Builder
	tabbedW := colLeftTabbedWidth
	verboseCount := states.Env().CountOfVerbose()
	cols, rows := s.safeGetTermSize()

	var painter Painter = s
	if s.asManual {
		painter = newManPainter()
	}

	if s.treeMode {
		// ~~tree: list all commands in tree style for a overview

		grouped := true
		s.printHeader(ctx, &sb, lastCmd, pc, cols, tabbedW)
		walkCtx := &cli.WalkBackwardsCtx{}
		lastCmd.WalkGrouped(ctx, func(cc, pp cli.Cmd, ff *cli.Flag, group string, idx, level int) {
			switch {
			case ff == nil: // CmdS
				walkCtx.LastCmdGroupInc = s.printCommand(ctx, &sb, painter, &verboseCount, cc, group, idx, level, cols, tabbedW, grouped)
			default: // Flag
				s.printFlag(ctx, walkCtx, &sb, painter, &verboseCount, ff, group, idx, level, cols, tabbedW, grouped)
			}
		})
		_, _ = wr.WriteString(sb.String())
		_, _ = wr.WriteString("\n")
	} else {
		// normal help screen

		painter.printHeader(ctx, &sb, lastCmd, pc, cols, tabbedW)
		painter.printUsage(ctx, &sb, lastCmd, pc, cols, tabbedW)
		painter.printDesc(ctx, &sb, lastCmd, pc, cols, tabbedW)
		painter.printExamples(ctx, &sb, lastCmd, pc, cols, tabbedW)

		walkCtx := &cli.WalkBackwardsCtx{
			Group: !s.w.DontGroupInHelpScreen,
			Sort:  s.w.SortInHelpScreen,
		}
		lastCmd.WalkBackwardsCtx(ctx, func(ctx context.Context, pc *cli.WalkBackwardsCtx, cc cli.Cmd, ff *cli.Flag, index, groupIndex, count, level int) {
			if ff == nil {
				p := cc.OwnerCmd()
				cnt := p.CountOfCommands()
				parentIsDynamicLoading := p.IsDynamicCommandsLoading()
				isFirstItem := index == 0 && (min(cnt, count) > 0 || parentIsDynamicLoading)
				if isFirstItem {
					painter.printCommandHeading(ctx, &sb, cc, "Commands")
				} else {
					// _, _ = sb.WriteString("\nCommands[")
					// _, _ = sb.WriteString(strconv.Itoa(cnt))
					// _, _ = sb.WriteString("/")
					// _, _ = sb.WriteString(strconv.Itoa(count))
					// _, _ = sb.WriteString("]:\n")
				}
				pc.LastCmdGroupInc = s.printCommand(ctx, &sb, painter, &verboseCount, cc, cc.GroupHelpTitle(), groupIndex, 1, cols, tabbedW, walkCtx.Group)
				return
			}

			p := ff.Owner()
			cnt := p.CountOfFlags()
			parentIsDynamicLoading := p.IsDynamicFlagsLoading()
			isFirstItem := index == 0 && (min(cnt, count) > 0 || parentIsDynamicLoading)
			if isFirstItem {
				if cc.OwnerCmd() == nil {
					// _, _ = sb.WriteString("\nGlobal Flags:\n")
					painter.printFlagHeading(ctx, &sb, cc, ff, "Global Flags")
					_, _ = sb.WriteString("\n")
				} else if level == 0 {
					// _, _ = sb.WriteString("\nFlags:\n")
					painter.printFlagHeading(ctx, &sb, cc, ff, "Flags")
					_, _ = sb.WriteString("\n")
				} else if level == 1 {
					painter.printFlagHeading(ctx, &sb, cc, ff, "Parent Flags")
					_, _ = sb.WriteString("(")
					_, _ = sb.WriteString(color.ToDim(cc.String()))
					_, _ = sb.WriteString("):\n")
				} else {
					painter.printFlagHeading(ctx, &sb, cc, ff, "Grandpa Flags")
					_, _ = sb.WriteString("(")
					_, _ = sb.WriteString(color.ToDim(cc.String()))
					_, _ = sb.WriteString("):\n")
				}
			}
			s.printFlag(ctx, pc, &sb, painter, &verboseCount, ff, ff.GroupHelpTitle(), groupIndex, 1, cols, tabbedW, walkCtx.Group)
		}, walkCtx)

		painter.printTailLine(ctx, &sb, lastCmd, pc, rows, cols, tabbedW)

		_, _ = wr.WriteString(sb.String())
		_, _ = wr.WriteString("\n")
	}

	logz.VerboseContext(ctx, "tty cols", "cols", cols, "rows", rows, "tree-mode", s.treeMode, "show-tree", s.w.Actions())

	if !s.debugMatches {
		return
	}

	sb.Reset()
	s.printDebugMatches(ctx, &sb, wr, pc)
}

func (s *helpPrinter) PrintDebugScreenTo(ctx context.Context, wr HelpWriter, pc cli.ParsedState, lastCmd cli.Cmd) {
	if s.Translator == nil {
		s.Translator = color.GetCPT()
	}

	var sb strings.Builder

	// tabbedW := colLeftTabbedWidth
	// verboseCount := states.Env().CountOfVerbose()
	// cols, rows := s.safeGetTermSize()

	text := s.w.Store().Dump()
	_, _ = sb.WriteString("\nStore:\n")
	_, _ = sb.WriteString(text)
	_, _ = sb.WriteString("\n")

	_, _ = wr.WriteString(sb.String())
	// _, _ = wr.WriteString("\n")

	sb.Reset()
	s.printEnv(ctx, &sb, wr, pc)

	sb.Reset()
	s.printRaw(ctx, &sb, wr, pc)

	sb.Reset()
	s.printMore(ctx, &sb, wr, pc)

	sb.Reset()
	s.printDebugMatches(ctx, &sb, wr, pc)
}

func (s *helpPrinter) safeGetTermSize() (cols, rows int) {
	cols, rows = term.GetTtySize()
	if cols == 0 || rows == 0 {
		const virtualTtyWidthOrHeight = 4096
		cols, rows = virtualTtyWidthOrHeight, virtualTtyWidthOrHeight
	} else {
		for _, en := range []string{"COLS", "COLUMNS", "TERM_COLS", "TERM_COLUMNS", "CMDR_COLS"} {
			if C := os.Getenv(en); C != "" {
				if cols64, err := strconv.ParseInt(C, 10, 64); err == nil && cols64 > 0 {
					cols = int(cols64)
					break
				}
			}
		}
	}
	return
}

type Painter interface {
	printHeader(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int)
	printUsage(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int)
	printDesc(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int)
	printExamples(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int)
	printNotes(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int)
	printTailLine(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, rows, cols, tabbedW int)

	// printEnv(ctx context.Context, sb *strings.Builder, wr HelpWriter, pc cli.ParsedState)
	// printRaw(ctx context.Context, sb *strings.Builder, wr HelpWriter, pc cli.ParsedState)
	// printMore(ctx context.Context, sb *strings.Builder, wr HelpWriter, pc cli.ParsedState)
	// printDebugMatches(ctx context.Context, sb *strings.Builder, wr HelpWriter, pc cli.ParsedState)

	// printCommand(ctx context.Context, sb *strings.Builder, verboseCount *int, cc cli.Cmd, group string, idx, level, cols, tabbedW int, grouped bool)

	printCommandHeading(ctx context.Context, sb *strings.Builder, cc cli.Cmd, title string)
	printCommandGroupTitle(ctx context.Context, sb *strings.Builder, group string, indent int)
	printCommandRow(ctx context.Context, sb *strings.Builder,
		cc cli.Cmd,
		indentSpaces, left, right, dep, depPlain string,
		cols, tabbedW int, deprecated, dim bool)

	// printFlag(ctx context.Context, sb *strings.Builder, verboseCount *int, ff *cli.Flag, group string, idx, level, cols, tabbedW int, grouped bool)

	printFlagHeading(ctx context.Context, sb *strings.Builder, cc cli.Cmd, ff *cli.Flag, title string)
	printFlagGroupTitle(ctx context.Context, sb *strings.Builder, group string, indent int)
	printFlagRow(ctx context.Context, sb *strings.Builder,
		ff *cli.Flag,
		indentSpaces, left, right, tg, def, defPlain, dep, depPlain, env, envPlain string,
		cols, tabbedW int, deprecated, dim bool)
}

var _ Painter = (*helpPrinter)(nil)

func (s *helpPrinter) translate(pc cli.ParsedState, pattern string, fg color.Color) string {
	return s.Translate(pc.Translate(os.ExpandEnv(pattern)), fg)
}

func (s *helpPrinter) printHeader(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int) {
	if cc.Root() == nil {
		if cc.OwnerCmd() != cc {
			logz.Fatal("unsatisfied link to root: cc.Root() is invalid", "cc", cc)
		}
		return
	}
	header := cc.Root().Header()
	line := exec.StripLeftTabs(os.ExpandEnv(header))
	_, _ = sb.WriteString(s.translate(pc, line, color.FgDefault))
	_, _ = sb.WriteString("\n")
	_, _, _ = pc, cols, tabbedW
	_ = ctx
}

func (s *helpPrinter) printTailLine(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, rows, cols, tabbedW int) {
	footer := strings.TrimSpace(cc.Root().Footer())
	if footer != "" {
		if strings.Contains(footer, "{{") && strings.Contains(footer, "}}") {
			var sb strings.Builder
			if tpl, err := template.New("footer").Parse(footer); err != nil {
				logz.FatalContext(ctx, "failed to parse footer template", "err", err, "footer", footer)
				return
			} else if err = tpl.Execute(&sb, struct{ Cols, Rows, Tabstop int }{
				Cols:    cols,
				Rows:    rows,
				Tabstop: tabbedW,
			}); err != nil {
				logz.FatalContext(ctx, "failed to execute footer template", "err", err, "footer", footer)
				return
			}
			footer = sb.String()
		}

		_, _ = sb.WriteString("\n")
		str := exec.StripLeftTabs(os.ExpandEnv(footer))
		// line := fmt.Sprintf("<dim>%s</dim>", str)
		_, _ = sb.WriteString(color.ToDim("%v", s.translate(pc, str, color.FgDefault)))
		if !strings.HasSuffix(footer, "\n") {
			_, _ = sb.WriteString("\n")
		}
		// if s.w.actionsMatched&actionShowTree != 0 {
		// 	_, _ = sb.WriteString("~~tree\n")
		// }
	}
	_, _, _ = pc, cols, tabbedW
	_ = ctx
}

func (s *helpPrinter) printUsage(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int) {
	// cc.App() could have a nil value while cc is a dynamic command.
	// But cc.Root() is always available and point to the proper target.
	appName := cc.Root().App().Name()
	titles := cc.GetCommandTitles()
	tail := "[files...]"
	if tph := cc.TailPlaceHolder(); tph != "" {
		tail = tph
	}
	line := fmt.Sprintf("$ <kbd>%s</kbd> %s [Options...]%s\n", appName, titles, tail)
	_, _ = sb.WriteString("\nUsage:\n\n  ")
	// _, _ = sb.WriteString("\n")
	_, _ = sb.WriteString(s.translate(pc, line, color.FgDefault))
	_, _, _ = pc, cols, tabbedW
	_ = ctx
}

func (s *helpPrinter) printDesc(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int) {
	desc := cc.DescLong()
	if desc != "" {
		_, _ = sb.WriteString("\nDescription:\n\n")
		desc = exec.StripLeftTabs(os.ExpandEnv(desc))
		line := color.ToDim("%v", s.translate(pc, desc, color.FgDefault))
		line = exec.LeftPad(line, 2)
		_, _ = sb.WriteString(line)
	}
	_, _, _ = pc, cols, tabbedW
	_ = ctx
}

func (s *helpPrinter) removeLastEmptyLines(lines []string) (lastLine int) {
	for lastLine = len(lines) - 1; lastLine >= 0; lastLine-- {
		if strings.TrimSpace(lines[lastLine]) == "" {
			continue
		}
		break
	}
	return
}

func (s *helpPrinter) printExamples(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int) {
	examples := cc.Examples()
	if examples != "" {
		_, _ = sb.WriteString("\nExamples:\n\n")
		str := exec.StripLeftTabs(os.ExpandEnv(examples))

		lines := strings.Split(str, "\n")
		for i, ln := range lines {
			if strings.HasPrefix(ln, "$ ") {
				lines[i] = "$ <font color=\"green\">" + ln[2:] + "</color>"
			} else if ln != "" {
				lines[i] = "<dim>" + ln + "</dim>"
			}
		}

		lastLine := s.removeLastEmptyLines(lines)
		str = strings.Join(lines[0:lastLine+1], "\n")
		line := s.translate(pc, str, color.FgLightGray)
		line = exec.LeftPad(line, 2)

		_, _ = sb.WriteString(line)
	}
	_, _, _ = pc, cols, tabbedW

	s.printNotes(ctx, sb, cc, pc, cols, tabbedW)
}

func (s *helpPrinter) printNotes(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int) {
	if root := cc.Root(); root.Cmd == cc && root.RedirectTo() != "" {
		_, _ = sb.WriteString("\nNotes:\n\n")

		str := exec.StripLeftTabs(fmt.Sprintf(`<i>Root Command was been redirected to Subcommand</i>: "<b>%s</b>"`, root.RedirectTo()))
		line := color.ToDim("%v", s.translate(pc, str, color.FgDefault))
		line = exec.LeftPad(line, 2)
		_, _ = sb.WriteString(line)

		if m := root.RedirectToSet(); m != nil {
			_, _ = sb.WriteString("\n")
			_, _ = sb.WriteString("  All Redirected Commands: \n\n")
			for k, v := range m {
				for to, froms := range v {
					for _, from := range froms {
						_, _ = sb.WriteString("    ")
						_, _ = sb.WriteString(s.translate(pc, fmt.Sprintf("<dim>%s --(<b>%s</b>)-> %s</dim>\n", from, k, to), color.FgDefault))
					}
				}
			}
			// _, _ = sb.WriteString("\n")
		}
	} else if rt := cc.RedirectTo(); rt != "" {
		_, _ = sb.WriteString("\nNotes:\n\n")
		str := exec.StripLeftTabs(fmt.Sprintf(`<i>This Command was been redirected to</i>: "<b>%s</b>"`, rt))
		line := color.ToDim("%v", s.translate(pc, str, color.FgDefault))
		line = exec.LeftPad(line, 2)
		_, _ = sb.WriteString(line)
	} else if cc1, ok := cc.(*cli.CmdS); ok {
		for k, v := range root.RedirectToSet() {
			if froms, ok := v[cc1]; ok {
				_, _ = sb.WriteString("\nNotes:\n\n")
				_, _ = sb.WriteString("  These commands redirect to here:\n\n")
				for _, from := range froms {
					_, _ = sb.WriteString("    ")
					_, _ = sb.WriteString(s.translate(pc, fmt.Sprintf("<dim>%s --(<b>%s</b>)-> Me</dim>\n", from, k), color.FgDefault))
				}
			}
		}
	}
	_, _, _ = pc, cols, tabbedW
	_ = ctx
}

func (s *helpPrinter) printEnv(ctx context.Context, sb *strings.Builder, wr HelpWriter, pc cli.ParsedState) {
	if found := pc.HasFlag("env", func(ff *cli.Flag, state *cli.MatchState) bool {
		return state.DblTilde && state.HitTimes > 0
	}); !found {
		return
	}

	_, _ = sb.WriteString("\nEnvironments:\n")
	var keys, extras []string
	b, m := false, map[string]string{}
	for _, line := range os.Environ() {
		a := strings.Split(line, "=")
		if a[0] == "_" {
			b = true
		}
		if b {
			extras = append(extras, a[0])
		} else {
			keys = append(keys, a[0])
		}
		m[a[0]] = a[1]
	}

	slices.Sort(keys)
	for _, key := range keys {
		_, _ = sb.WriteString("  ")
		_, _ = sb.WriteString(key)
		_, _ = sb.WriteString(" = ")
		_, _ = sb.WriteString(color.ToDim(m[key]))
		_, _ = sb.WriteString("\n")
	}

	for _, key := range extras {
		_, _ = sb.WriteString("  ")
		_, _ = sb.WriteString(key)
		_, _ = sb.WriteString(" = ")
		_, _ = sb.WriteString(color.ToDim(m[key]))
		_, _ = sb.WriteString("\n")
	}

	_, _ = wr.WriteString(sb.String())
	_, _ = wr.WriteString("\n")
	_ = ctx
}

func (s *helpPrinter) printRaw(ctx context.Context, sb *strings.Builder, wr HelpWriter, pc cli.ParsedState) {
	if found := pc.HasFlag("raw", func(ff *cli.Flag, state *cli.MatchState) bool {
		return state.DblTilde && state.HitTimes > 0
	}); !found {
		return
	}

	_, _ = sb.WriteString("\nRaw:\n")
	_, _ = wr.WriteString(sb.String())
	_, _ = wr.WriteString("\n")
	_ = ctx
}

func (s *helpPrinter) printMore(ctx context.Context, sb *strings.Builder, wr HelpWriter, pc cli.ParsedState) {
	if found := pc.HasFlag("more", func(ff *cli.Flag, state *cli.MatchState) bool {
		return state.DblTilde && state.HitTimes > 0
	}); !found {
		return
	}

	_, _ = sb.WriteString("\nMore:\n")
	_, _ = wr.WriteString(sb.String())
	_, _ = wr.WriteString("\n")
	_ = ctx
}

func (s *helpPrinter) printDebugMatches(ctx context.Context, sb *strings.Builder, wr HelpWriter, pc cli.ParsedState) {
	if len(pc.MatchedCommands()) > 0 {
		_, _ = sb.WriteString("\nMatched commands:\n")
		for i, cc := range pc.MatchedCommands() {
			_, _ = sb.WriteString(s.Translate(fmt.Sprintf("  - %d. <code>%s</code> | %v\n", i+1, cc.HitTitle(), cc), color.FgDefault))
		}
	}
	if len(pc.MatchedFlags()) > 0 {
		_, _ = sb.WriteString("\nMatched flags:\n")
		i := 0
		for ff, st := range pc.MatchedFlags() {
			i++
			short, tilde := "", ""
			if st.Short {
				short = "short"
			}
			if st.DblTilde {
				if st.Short {
					tilde = ",TILDE"
				} else {
					tilde = "TILDE"
				}
			}
			_, _ = sb.WriteString(s.translate(pc,
				fmt.Sprintf(
					"  - %d. <code>%s</code> <dim>(+%v)</dim> %v <dim>/%v%v/</dim> | <dim>[owner: %v]</dim>\n",
					i, ff.GetHitStr(), ff.GetTriggeredTimes(), ff, short, tilde, ff.Owner().String()),
				color.FgDefault))
		}
	}

	if s.w != nil {
		_, _ = sb.WriteString("\nACTIONS:\n")
		_, _ = sb.WriteString(s.w.actionsMatched.String())
		_, _ = sb.WriteString("\n")
	}

	if sb.Len() > 0 {
		if s.w != nil && s.w.wrDebugScreen != nil {
			wr = s.w.wrDebugScreen
		}
		_, _ = wr.WriteString(sb.String())
		// _, _ = wr.WriteString("\n")
	}
	_ = ctx
}

func printleftpad(sb *strings.Builder, cond bool, tabbedW int) {
	if cond {
		_, _ = sb.WriteString("\n")
		_, _ = sb.WriteString(strings.Repeat(" ", tabbedW))
	}
}

func trans(ss string, translator color.Translator, clr color.Color, deprecated bool) string {
	ss = translator.Translate(strings.TrimSpace(ss), clr)
	if deprecated {
		ss = term.StripEscapes(ss)
	}
	return ss
}

func (s *helpPrinter) printCommandHeading(ctx context.Context, sb *strings.Builder, cc cli.Cmd, title string) {
	_, _ = sb.WriteString(fmt.Sprintf("\n%s:\n\n", title))
}

func (s *helpPrinter) printFlagHeading(ctx context.Context, sb *strings.Builder, cc cli.Cmd, ff *cli.Flag, title string) {
	_, _ = sb.WriteString(fmt.Sprintf("\n%s:\n", title))
}

func (s *helpPrinter) printCommandGroupTitle(ctx context.Context, sb *strings.Builder, group string, indent int) {
	_, _ = sb.WriteString(strings.Repeat("  ", indent))

	noColor := is.NoColorMode()
	colorful := !noColor
	if colorful {
		s.WriteColor(sb, CurrentGroupTitleColor)
		s.WriteBgColor(sb, CurrentGroupTitleBgColor)
	}
	_, _ = sb.WriteString("[")
	_, _ = sb.WriteString(group)
	_, _ = sb.WriteString("]")
	if colorful {
		s.Reset(sb)
	}
	_, _ = sb.WriteString("\n")
}

func (s *helpPrinter) printCommandRow(ctx context.Context, sb *strings.Builder,
	cc cli.Cmd,
	indentSpaces, left, right, dep, depPlain string,
	cols, tabbedW int, deprecated, dim bool,
) {
	_, _ = sb.WriteString(indentSpaces)

	noColor := is.NoColorMode()
	colorful := !noColor

	if colorful {
		if dim {
			s.WriteBgColor(sb, color.BgDim)
		}
		if deprecated {
			s.WriteBgColor(sb, color.BgStrikeout)
			s.WriteColor(sb, CurrentDescColor)
		} else {
			s.WriteColor(sb, CurrentTitleColor)
		}
	}
	_, _ = sb.WriteString(left)

	var split bool
	var printed int
	if right != "" {
		if colorful {
			s.WriteColor(sb, CurrentDescColor)
		}

		rCols := cols - tabbedW
		l := len(right) + len(depPlain)
		if l >= rCols {
			var prt string
			var ix int
			for len(right) > rCols {
				prt, right = right[:rCols], right[rCols:]
				printleftpad(sb, ix > 0, tabbedW)
				_, _ = sb.WriteString(trans(prt, s, CurrentDescColor, deprecated))
				ix++
			}
			if right != "" {
				if ix > 0 {
					printleftpad(sb, ix > 0, tabbedW)
				} else {
					split, printed = true, len(right)
				}
				_, _ = sb.WriteString(trans(right, s, CurrentDescColor, deprecated))
			}
		} else {
			_, _ = sb.WriteString(trans(right, s, CurrentDescColor, deprecated))
		}
		// sb.WriteString(trans(right, CurrentDescColor))
	} else {
		if colorful {
			s.WriteColor(sb, CurrentDescColor)
			_, _ = sb.WriteString(trans("<i>(no desc)</i>", s, CurrentDescColor, deprecated))
		} else {
			_, _ = sb.WriteString("(no desc)")
		}
	}

	if dep != "" {
		if split {
			printleftpad(sb, split, tabbedW)
			split = false
		}
		if printed >= 0 {
			_, _ = sb.WriteString(" ")
			_, _ = sb.WriteString(dep)
		}
		logz.VerboseContext(ctx, "[watching] split flag", "split", split)
	}

	if colorful {
		s.Reset(sb) // reset fg/bg colors by color Translator
	}
}

func (s *helpPrinter) printCommand(ctx context.Context, sb *strings.Builder,
	painter Painter,
	verboseCount *int, cc cli.Cmd,
	group string, idx, level, cols, tabbedW int, grouped bool,
) (groupedInc int) {
	if (cc.HiddenBR() && *verboseCount < 1) || (cc.VendorHiddenBR() && *verboseCount < 3) {
		return
	}

	_ = idx
	if grouped {
		if grp := cc.GroupHelpTitle(); grp != "" {
			if grp != s.lastCmdGroup {
				s.lastCmdGroup = grp
				if group != "" {
					painter.printCommandGroupTitle(ctx, sb, group, level+groupedInc)
					// _, _ = sb.WriteString(strings.Repeat("  ", level+groupedInc))
					// s.WriteColor(sb, CurrentGroupTitleColor)
					// s.WriteBgColor(sb, CurrentGroupTitleBgColor)
					// _, _ = sb.WriteString("[")
					// _, _ = sb.WriteString(group)
					// _, _ = sb.WriteString("]")
					// s.Reset(sb)
					// _, _ = sb.WriteString("\n")
					groupedInc++
				}
			} else {
				groupedInc++
			}
		}
	}

	indentSpaces := strings.Repeat("  ", level+groupedInc)
	// _, _ = sb.WriteString(indentSpaces)

	w := tabbedW - (level+groupedInc)*2
	ttl, restTitles := cc.GetTitleNames(tabbedW - (level+groupedInc)*2)
	if ttl == "" {
		ttl = "(not-specified)"
	}

	if !grouped && group != "" {
		ttl += " /[" + group + "]"
	}

	dim := (cc.HiddenBR() && *verboseCount > 0) || (cc.VendorHiddenBR() && *verboseCount >= 3)
	deprecated := cc.Deprecated() != ""
	// trans := func(ss string, clr color.Color) string {
	// 	ss = s.Translate(strings.TrimSpace(ss), clr)
	// 	if deprecated {
	// 		ss = term.StripEscapes(ss)
	// 	}
	// 	return ss
	// }

	if w >= len(ttl) {
		w -= len(ttl)
	}

	if root := cc.Root(); root != nil && root.RedirectTo() == cc.Name() {
		var ss strings.Builder
		s.Translator.HighlightFast(&ss, ttl)
		s.Translator.DimFast(&ss, " <- (root)")
		if w >= 10 {
			w -= 10
		}
		ttl = ss.String()
	}

	left, right := fmt.Sprintf("%s%s", ttl, strings.Repeat(" ", w)), cc.Desc()
	dep, depPlain := cc.DeprecatedHelpString(func(ss string, clr color.Color) string {
		return trans(ss, s, clr, deprecated)
	}, CurrentDeprecatedColor, CurrentDescColor)

	painter.printCommandRow(ctx, sb, cc, indentSpaces, left, right, dep, depPlain, cols, tabbedW, deprecated, dim)
	_, _ = sb.WriteString("\n")

	if restTitles != "" {
		_, _ = sb.WriteString(indentSpaces)
		_, _ = sb.WriteString("   = ")

		noColor := is.NoColorMode()
		colorful := !noColor
		if colorful {
			s.WriteBgColor(sb, color.BgItalic)
		}
		_, _ = sb.WriteString(restTitles)
		if colorful {
			s.Reset(sb)
		}
		_, _ = sb.WriteString("\n")
	}

	return
}

func (s *helpPrinter) printFlagGroupTitle(ctx context.Context, sb *strings.Builder, group string, indent int) {
	_, _ = sb.WriteString(strings.Repeat("  ", indent))
	noColor := is.NoColorMode()
	colorful := !noColor
	if colorful {
		s.WriteColor(sb, CurrentGroupTitleColor)
		s.WriteBgColor(sb, CurrentGroupTitleBgColor)
	}
	_, _ = sb.WriteString("[")
	_, _ = sb.WriteString(group)
	_, _ = sb.WriteString("]")
	if colorful {
		s.Reset(sb)
	}
	_, _ = sb.WriteString("\n")
}

func (s *helpPrinter) printFlagRow(ctx context.Context, sb *strings.Builder,
	ff *cli.Flag,
	indentSpaces, left, right, tg, def, defPlain, dep, depPlain, env, envPlain string,
	cols, tabbedW int, deprecated, dim bool,
) {
	_, _ = sb.WriteString(indentSpaces)

	noColor := is.NoColorMode()
	colorful := !noColor

	if ff.Required() {
		_, _ = sb.WriteString("* ")
	}

	if ff.Short == "" {
		sb.WriteRune(' ')
	}

	if colorful {
		if dim {
			s.WriteBgColor(sb, color.BgDim)
		}
		if deprecated {
			s.WriteBgColor(sb, color.BgStrikeout)
			s.WriteColor(sb, CurrentDescColor)
		} else {
			s.WriteColor(sb, CurrentTitleColor)
		}
	}
	_, _ = sb.WriteString(left)

	// s.DimFast(&sb, s.Translate(right, color.BgNormal))
	if tg != "" {
		// s.ColoredFast(&sb, CurrentFlagTitleColor, tg)
		if colorful {
			s.WriteColor(sb, CurrentFlagTitleColor)
		}
		_, _ = sb.WriteString(tg)
	}

	var split bool
	var printed int
	// printleftpad := func(cond bool) {
	// 	if cond {
	// 		_, _ = sb.WriteString("\n")
	// 		_, _ = sb.WriteString(strings.Repeat(" ", tabbedW))
	// 	}
	// }
	rCols := cols - tabbedW
	if right != "" {
		if colorful {
			s.WriteColor(sb, CurrentDescColor)
		}

		_, l, l1st := len(right), len(right)+len(defPlain)+len(depPlain)+len(envPlain), len(tg)
		// aa := []string{}
		if l+l1st >= rCols {
			var prt string
			var ix int
			for len(right)+l1st >= rCols {
				prt, right = right[:rCols-l1st], right[rCols-l1st:]
				printleftpad(sb, ix > 0, tabbedW)
				// aa = append(aa, prt)
				_, _ = sb.WriteString(trans(prt, s, CurrentDescColor, deprecated))
				ix++
				l1st = 0
			}
			if right != "" {
				str := trans(right, s, CurrentDescColor, deprecated)
				if ix > 0 {
					printleftpad(sb, ix > 0, tabbedW)
				} else {
					split, printed = true, len(is.StripEscapes(str))
				}
				_, _ = sb.WriteString(str)
			}
		} else {
			if colorful {
				_, _ = sb.WriteString(trans(right, s, CurrentDescColor, deprecated))
			} else {
				_, _ = sb.WriteString(right)
			}
		}

		// if ff.Long == "addr" {
		// 	sb.WriteString(fmt.Sprintf(" / l=%d, l1st=%d, len(right)/def=%d/%d, rCols=%d", l, l1st, len(right), len(defPlain), rCols))
		// }
		// sb.WriteString(fmt.Sprintf(" / l=%d/%d, l1st=%d, len(right)=%d, rCols=%d", l, lr, l1st, len(right), rCols))
		// sb.WriteString(trans(right, CurrentDescColor))
	} else {
		if colorful {
			s.WriteColor(sb, CurrentDescColor)
			_, _ = sb.WriteString(trans("<i>(no desc)</i>", s, CurrentDescColor, deprecated))
		} else {
			_, _ = sb.WriteString("(no desc)")
		}
		printed += 4
	}

	if env != "" && printed >= 0 {
		if split {
			envlen := len(envPlain)
			printed += envlen
			if printed >= rCols {
				printleftpad(sb, split, tabbedW)
				printed = envlen
			}
		}
		if sb.String()[sb.Len()-1] != ' ' {
			_, _ = sb.WriteRune(' ')
		}
		_, _ = sb.WriteString(env)
	}

	if def != "" && printed >= 0 {
		if split {
			deflen := len(defPlain) // len(is.StripEscapes(def))
			printed += deflen
			if printed >= rCols {
				printleftpad(sb, split, tabbedW)
				printed = deflen
			}
		}
		if sb.String()[sb.Len()-1] != ' ' {
			_, _ = sb.WriteRune(' ')
		}
		_, _ = sb.WriteString(def)
	}

	if ff.Required() {
		str := "<kbd>REQUIRED</kbd>"
		esc := s.Translate(str, CurrentFlagTitleColor)
		if split {
			deflen := len(str)
			printed += deflen
			if printed >= rCols {
				printleftpad(sb, split, tabbedW)
				printed = deflen
			}
		}
		if sb.String()[sb.Len()-1] != ' ' {
			_, _ = sb.WriteRune(' ')
		}
		_, _ = sb.WriteString(esc)
	}

	if dep != "" {
		if split {
			deplen := len(depPlain) // len(is.StripEscapes(dep))
			printed += deplen
			if printed >= rCols {
				printleftpad(sb, split, tabbedW)
				printed = deplen
			}
		}
		if sb.String()[sb.Len()-1] != ' ' {
			_, _ = sb.WriteRune(' ')
		}
		_, _ = sb.WriteString(dep)
		logz.VerboseContext(ctx, "split flag is", "split", split)
	}

	if va := ff.ValidArgs(); len(va) > 0 {
		cvt := evendeep.Cvt{}
		str := cvt.String(va)
		sl, w := len(str), cols-tabbedW
		split = true
		inc := 0
		if colorful {
			s.WriteBgColor(sb, color.BgDim)
			s.WriteBgColor(sb, color.BgItalic)
		}
		for printed = 0; printed < sl; {
			// _, _ = sb.WriteRune('\n')
			printleftpad(sb, split, tabbedW)
			// n, _ := sb.WriteString("Valid Args: ")
			end := printed + w - inc
			if end > sl {
				end = sl
			}
			sub := str[printed:end]
			// s.Dim(sb, func(out io.Writer) {
			// s.Italic(sb, func(out io.Writer) {
			// 	if inc == 1 {
			// 		_, _ = out.Write([]byte{' '})
			// 	}
			// 	_, _ = out.Write([]byte(sub))
			// })
			// })

			if inc == 1 {
				_, _ = sb.Write([]byte{' '})
			}
			_, _ = sb.Write([]byte(sub))

			printed += len(sub)
			inc = 1
		}
	}

	// s.ColoredFast(&sb, CurrentDefaultValueColor, def)
	// s.ColoredFast(&sb, CurrentDeprecatedColor, dep)
	// sb.WriteString(s.Translate(right, color.BgDefault))
	if colorful {
		s.Reset(sb)
	}
}

func (s *helpPrinter) printFlag(ctx context.Context, pc *cli.WalkBackwardsCtx,
	sb *strings.Builder,
	painter Painter,
	verboseCount *int, ff *cli.Flag, group string,
	idx, level, cols, tabbedW int, grouped bool,
) {
	if (ff.HiddenBR() && *verboseCount < 1) || (ff.VendorHiddenBR() && *verboseCount < 3) {
		return
	}

	groupedInc := pc.LastCmdGroupInc
	if s.treeMode {
		groupedInc++
	}

	ofs := 0
	_ = idx
	if grouped && group != "" {
		if grp := ff.GroupHelpTitle(); grp != s.lastFlagGroup {
			s.lastFlagGroup = grp
			painter.printFlagGroupTitle(ctx, sb, group, level+groupedInc)
			// _, _ = sb.WriteString(strings.Repeat("  ", level+groupedInc))
			// s.WriteColor(sb, CurrentGroupTitleColor)
			// s.WriteBgColor(sb, CurrentGroupTitleBgColor)
			// _, _ = sb.WriteString("[")
			// _, _ = sb.WriteString(group)
			// _, _ = sb.WriteString("]")
			// s.Reset(sb)
			// _, _ = sb.WriteString("\n")
		}
		groupedInc++
		if ff.Required() {
			ofs = -1
		}
	} else if grouped && !ff.OwnerIsNotNil() { // don't apply on a detached flag item
		groupedInc++
		if ff.Required() {
			ofs = -1
		}
	}

	indentSpaces := strings.Repeat("  ", level+groupedInc+ofs)
	// ttl := strings.Join(ff.GetTitleZshFlagNamesArray(), ",")
	ttl, restTitles := ff.GetTitleFlagNamesBy(",", tabbedW-len(indentSpaces))
	w := tabbedW - (level+groupedInc)*2 // - len(ttl)

	// _, _ = sb.WriteString(indentSpaces)
	//
	// if ff.Required() {
	// 	_, _ = sb.WriteString("* ")
	// }
	//
	// if ff.Short == "" {
	// 	sb.WriteRune(' ')
	// 	w--
	// }

	if ff.Short == "" {
		w--
	}

	dim := (ff.HiddenBR() && *verboseCount > 0) || (ff.VendorHiddenBR() && *verboseCount >= 3)
	deprecated := ff.Deprecated() != ""
	// trans := func(ss string, clr color.Color) string {
	// 	ss = s.Translate(strings.TrimSpace(ss), clr)
	// 	if deprecated {
	// 		ss = term.StripEscapes(ss)
	// 	}
	// 	return ss
	// }

	// left, right := fmt.Sprintf("%-"+strconv.Itoa(w)+"s", ttl), ff.Desc()
	if w >= len(ttl) {
		w -= len(ttl)
	}
	left, right := fmt.Sprintf("%s%s", ttl, strings.Repeat(" ", w)), ff.Desc()
	tg := ff.ToggleGroupLeadHelpString()
	trans1 := func(ss string, clr color.Color) string { return trans(ss, s, clr, deprecated) }
	def, defPlain := ff.DefaultValueHelpString(trans1, CurrentDefaultValueColor, CurrentDescColor)
	dep, depPlain := ff.DeprecatedHelpString(trans1, CurrentDeprecatedColor, CurrentDescColor)
	env, envPlain := ff.EnvVarsHelpString(trans1, CurrentEnvVarsColor, CurrentDescColor)

	painter.printFlagRow(ctx, sb, ff, indentSpaces, left, right, tg, def, defPlain, dep, depPlain, env, envPlain, cols, tabbedW, deprecated, dim)
	_, _ = sb.WriteString("\n")

	if restTitles != "" {
		_, _ = sb.WriteString(indentSpaces)
		_, _ = sb.WriteString("    = ")

		noColor := is.NoColorMode()
		colorful := !noColor
		if colorful {
			s.WriteBgColor(sb, color.BgItalic)
		}
		_, _ = sb.WriteString(restTitles)
		if colorful {
			s.Reset(sb)
		}
		_, _ = sb.WriteString("\n")
	}

	if ff.HeadLike() && !s.asManual {
		_, _ = sb.WriteString(indentSpaces)
		_, _ = sb.WriteString("    ")
		if ff.Required() {
			_, _ = sb.WriteString("  ")
		}
		row := fmt.Sprintf("-<i>number</i> = --%s=<i>number</i>\n", ff.Title())
		esc := s.Translate(row, CurrentFlagTitleColor)
		_, _ = sb.WriteString(esc)
	}
}

var (
	//
	// doNotLoadingConfigFiles = false

	// // rootCommand the root of all commands
	// rootCommand *RootCommand
	// // rootOptions *Opt
	// rxxtOptions = newOptions()

	// usedConfigFile
	// usedConfigFile            string
	// usedConfigSubDir          string
	// configFiles               []string
	// onConfigReloadedFunctions map[ConfigReloaded]bool
	//
	// predefinedLocations = []string{
	// 	"./ci/etc/%s/%s.yml",
	// 	"/etc/%s/%s.yml",
	// 	"/usr/local/etc/%s/%s.yml",
	// 	os.Getenv("HOME") + "/.%s/%s.yml",
	// }

	//
	// defaultStdout = bufio.NewWriterSize(os.Stdout, 16384)
	// defaultStderr = bufio.NewWriterSize(os.Stderr, 16384)

	//
	// currentHelpPainter Painter

	CurrentTitleColor     = color.FgDefault
	CurrentFlagTitleColor = color.FgGreen

	// CurrentHiddenColor the print color for left part of a hidden opt
	CurrentHiddenColor = color.FgDarkGray

	// CurrentDeprecatedColor the print color for deprecated opt line
	CurrentDeprecatedColor = color.FgDarkGray

	CurrentEnvVarsColor = color.FgLightGray

	// CurrentDescColor the print color for description line
	CurrentDescColor = color.FgDarkGray
	// CurrentDefaultValueColor the print color for default value line
	CurrentDefaultValueColor = color.FgCyan
	// CurrentGroupTitleColor the print color for titles
	CurrentGroupTitleColor   = color.FgWhite
	CurrentGroupTitleBgColor = color.BgDim

	// globalShowVersion   func()
	// globalShowBuildInfo func()

	// beforeXrefBuilding []HookFunc
	// afterXrefBuilt     []HookFunc

	// getEditor sets callback to get editor program
	// getEditor func() (string, error)

	// defaultStringMetric = tool.JaroWinklerDistance(tool.JWWithThreshold(similarThreshold))
)
