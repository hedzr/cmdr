package worker

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/hedzr/is"
	"github.com/hedzr/is/states"
	"github.com/hedzr/is/term"
	"github.com/hedzr/is/term/color"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/exec"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

type discardP struct{}

func (*discardP) Write([]byte) (n int, err error) { return }

func (*discardP) WriteString(string) (n int, err error) { return }

type helpPrinter struct {
	color.Translator
	debugScreenMode bool
	debugMatches    bool
	treeMode        bool
	w               *workerS
}

const colLeftTabbedWidth = 56

func (s *helpPrinter) Print(ctx context.Context, pc cli.ParsedState, lastCmd cli.Cmd) { //nolint:revive //
	wr := s.safeGetWriter()
	s.PrintTo(ctx, wr, pc, lastCmd)
}

func (s *helpPrinter) PrintTo(ctx context.Context, wr HelpWriter, pc cli.ParsedState, lastCmd cli.Cmd) { //nolint:revive //
	if s.debugScreenMode {
		s.PrintDebugScreenTo(ctx, wr, pc, lastCmd)
		return
	}

	if s.Translator == nil {
		s.Translator = color.GetCPT()
	}

	var sb strings.Builder
	tabbedW := colLeftTabbedWidth
	verboseCount := states.Env().CountOfVerbose()
	cols, rows := s.safeGetTermSize()

	if s.treeMode {
		// ~~tree: list all commands in tree style for a overview

		grouped := true
		s.printHeader(ctx, &sb, lastCmd, pc, cols, tabbedW)
		lastCmd.WalkGrouped(ctx, func(cc, pp cli.Cmd, ff *cli.Flag, group string, idx, level int) { //nolint:revive
			switch {
			case ff == nil: // CmdS
				s.printCommand(ctx, &sb, &verboseCount, cc, group, idx, level, cols, tabbedW, grouped)
			default: // Flag
				s.printFlag(ctx, &sb, &verboseCount, ff, group, idx, level, cols, tabbedW, grouped)
			}
		})
		_, _ = wr.WriteString(sb.String())
		_, _ = wr.WriteString("\n")
	} else {
		// normal help screen

		s.printHeader(ctx, &sb, lastCmd, pc, cols, tabbedW)
		s.printUsage(ctx, &sb, lastCmd, pc, cols, tabbedW)
		s.printDesc(ctx, &sb, lastCmd, pc, cols, tabbedW)
		s.printExamples(ctx, &sb, lastCmd, pc, cols, tabbedW)

		walkCtx := &cli.WalkBackwardsCtx{
			Group: !s.w.DontGroupInHelpScreen,
			Sort:  s.w.SortInHelpScreen,
		}
		lastCmd.WalkBackwardsCtx(ctx, func(ctx context.Context, pc *cli.WalkBackwardsCtx, cc cli.Cmd, ff *cli.Flag, index, groupIndex, count, level int) {
			if ff == nil {
				cnt := cc.OwnerCmd().CountOfCommands()
				if index == 0 && min(cnt, count) > 0 {
					_, _ = sb.WriteString("\nCommands:\n")
				} else { //nolint:revive,staticcheck
					// _, _ = sb.WriteString("\nCommands[")
					// _, _ = sb.WriteString(strconv.Itoa(cnt))
					// _, _ = sb.WriteString("/")
					// _, _ = sb.WriteString(strconv.Itoa(count))
					// _, _ = sb.WriteString("]:\n")
				}
				s.printCommand(ctx, &sb, &verboseCount, cc, cc.GroupHelpTitle(), groupIndex, 1, cols, tabbedW, walkCtx.Group)
				return
			}

			cnt := ff.Owner().CountOfFlags()
			if index == 0 && min(cnt, count) > 0 {
				if cc.OwnerCmd() == nil {
					_, _ = sb.WriteString("\nGlobal Flags:\n")
				} else if level == 0 {
					_, _ = sb.WriteString("\nFlags:\n")
				} else if level == 1 {
					_, _ = sb.WriteString("\nParent Flags (")
					_, _ = sb.WriteString(color.ToDim(cc.String()))
					_, _ = sb.WriteString("):\n")
				} else {
					_, _ = sb.WriteString("\nGrandpa Flags (")
					_, _ = sb.WriteString(color.ToDim(cc.String()))
					_, _ = sb.WriteString("):\n")
				}
			}
			s.printFlag(ctx, &sb, &verboseCount, ff, ff.GroupHelpTitle(), groupIndex, 1, cols, tabbedW, walkCtx.Group)
		}, walkCtx)

		s.printTailLine(ctx, &sb, lastCmd, pc, cols, tabbedW)

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

func (s *helpPrinter) PrintDebugScreenTo(ctx context.Context, wr HelpWriter, pc cli.ParsedState, lastCmd cli.Cmd) { //nolint:revive //
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

func (s *helpPrinter) safeGetWriter() (wr HelpWriter) {
	wr = os.Stdout
	if s.w != nil && s.w.wrHelpScreen != nil {
		wr = s.w.wrHelpScreen
	}
	return
}

func (s *helpPrinter) safeGetTermSize() (cols, rows int) {
	cols, rows = term.GetTtySize()
	if cols == 0 || rows == 0 {
		const virtualTtyWidthOrHeight = 4096
		cols, rows = virtualTtyWidthOrHeight, virtualTtyWidthOrHeight
	}
	return
}

func (s *helpPrinter) printHeader(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int) { //nolint:revive,unparam
	// app, root := cc.App(), cc.Root()
	// _ = app
	if cc.Root() == nil {
		logz.Fatal("unsatisfied link to root: cc.Root() is invalid", "cc", cc)
	}
	header := cc.Root().Header()
	line := exec.StripLeftTabs(header)
	_, _ = sb.WriteString(s.Translate(line, color.FgDefault))
	_, _ = sb.WriteString("\n")
	_, _, _ = pc, cols, tabbedW
}

func (s *helpPrinter) printUsage(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int) { //nolint:revive,unparam
	// app, root := cc.App(), cc.Root()
	// _ = app
	appName := cc.App().Name()
	titles := cc.GetCommandTitles()
	tail := "[files...]"
	if tph := cc.TailPlaceHolder(); tph != "" {
		tail = tph
	}
	line := fmt.Sprintf("$ <kbd>%s</kbd> %s [Options...]%s\n", appName, titles, tail)
	_, _ = sb.WriteString("\nUsage:\n\n  ")
	// _, _ = sb.WriteString("\n")
	_, _ = sb.WriteString(s.Translate(line, color.FgDefault))
	_, _, _ = pc, cols, tabbedW
}

func (s *helpPrinter) printDesc(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int) { //nolint:revive,unparam
	desc := cc.DescLong()
	if desc != "" {
		_, _ = sb.WriteString("\nDescription:\n\n")
		desc = exec.StripLeftTabs(desc)
		desc = pc.Translate(desc)
		line := s.Translate(desc, color.FgDefault)
		line = exec.LeftPad(line, 2)
		_, _ = sb.WriteString(line)
	}
	_, _, _ = pc, cols, tabbedW
}

func (s *helpPrinter) printExamples(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int) { //nolint:revive,unparam
	examples := cc.Examples()
	if examples != "" {
		_, _ = sb.WriteString("\nExamples:\n\n")
		str := exec.StripLeftTabs(examples)
		str = pc.Translate(str)
		line := s.Translate(str, color.FgDefault)
		line = exec.LeftPad(line, 2)
		_, _ = sb.WriteString(line)
	}
	_, _, _ = pc, cols, tabbedW
}

func (s *helpPrinter) printTailLine(ctx context.Context, sb *strings.Builder, cc cli.Cmd, pc cli.ParsedState, cols, tabbedW int) { //nolint:revive,unparam
	footer := strings.TrimSpace(cc.Root().Footer())
	if footer != "" {
		_, _ = sb.WriteString("\n")
		str := exec.StripLeftTabs(footer)
		line := fmt.Sprintf("<dim>%s</dim>", str)
		_, _ = sb.WriteString(s.Translate(line, color.FgDefault))
		if !strings.HasSuffix(footer, "\n") {
			_, _ = sb.WriteString("\n")
		}
		// if s.w.actionsMatched&actionShowTree != 0 {
		// 	_, _ = sb.WriteString("~~tree\n")
		// }
	}
	_, _, _ = pc, cols, tabbedW
}

func (s *helpPrinter) printEnv(ctx context.Context, sb *strings.Builder, wr HelpWriter, pc cli.ParsedState) {
	if found := pc.HasFlag("env", func(ff *cli.Flag, state *cli.MatchState) bool { //nolint:revive
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
}

func (s *helpPrinter) printRaw(ctx context.Context, sb *strings.Builder, wr HelpWriter, pc cli.ParsedState) {
	if found := pc.HasFlag("raw", func(ff *cli.Flag, state *cli.MatchState) bool { //nolint:revive
		return state.DblTilde && state.HitTimes > 0
	}); !found {
		return
	}

	_, _ = sb.WriteString("\nRaw:\n")
	_, _ = wr.WriteString(sb.String())
	_, _ = wr.WriteString("\n")
}

func (s *helpPrinter) printMore(ctx context.Context, sb *strings.Builder, wr HelpWriter, pc cli.ParsedState) {
	if found := pc.HasFlag("more", func(ff *cli.Flag, state *cli.MatchState) bool { //nolint:revive
		return state.DblTilde && state.HitTimes > 0
	}); !found {
		return
	}

	_, _ = sb.WriteString("\nMore:\n")
	_, _ = wr.WriteString(sb.String())
	_, _ = wr.WriteString("\n")
}

func (s *helpPrinter) printDebugMatches(ctx context.Context, sb *strings.Builder, wr HelpWriter, pc cli.ParsedState) { //nolint:revive
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
			_, _ = sb.WriteString(s.Translate(fmt.Sprintf(
				"  - %d. <code>%s</code> <dim>(+%v)</dim> %v <dim>/%v%v/</dim> | <dim>[owner: %v]</dim>\n",
				i, ff.GetHitStr(), ff.GetTriggeredTimes(), ff, short, tilde, ff.Owner().String()), color.FgDefault))
		}
	}

	if s.w != nil {
		_, _ = sb.WriteString("\nACTIONS:\n")
		_, _ = sb.WriteString(s.w.actionsMatched.String())
		_, _ = sb.WriteString("\n")
	}

	if sb.Len() > 0 {
		if s.w != nil && s.w.wrDebugScreen != nil {
			wr = s.w.wrDebugScreen //nolint:revive
		}
		_, _ = wr.WriteString(sb.String())
		// _, _ = wr.WriteString("\n")
	}
}

func (s *helpPrinter) printCommand(ctx context.Context, sb *strings.Builder, verboseCount *int, cc cli.Cmd, group string, idx, level, cols, tabbedW int, grouped bool) { //nolint:revive
	if (cc.Hidden() && *verboseCount < 1) || (cc.VendorHidden() && *verboseCount < 3) { //nolint:revive
		return
	}

	groupedInc := 0
	if grouped && group != "" {
		if idx == 0 {
			_, _ = sb.WriteString(strings.Repeat("  ", level))
			s.WriteColor(sb, CurrentGroupTitleColor)
			s.WriteBgColor(sb, CurrentGroupTitleBgColor)
			_, _ = sb.WriteString("[")
			_, _ = sb.WriteString(group)
			_, _ = sb.WriteString("]")
			s.Reset(sb)
			_, _ = sb.WriteString("\n")
		}
		groupedInc++
	}

	_, _ = sb.WriteString(strings.Repeat("  ", level+groupedInc))
	w := tabbedW - (level+groupedInc)*2
	ttl := cc.GetTitleNames()
	if ttl == "" {
		ttl = "(empty)"
	}

	if !grouped && group != "" {
		ttl += " /[" + group + "]"
	}

	deprecated := cc.Deprecated() != ""
	trans := func(ss string, clr color.Color) string {
		ss = s.Translate(strings.TrimSpace(ss), clr)
		if deprecated {
			ss = term.StripEscapes(ss)
		}
		return ss
	}

	left, right := fmt.Sprintf("%-"+strconv.Itoa(w)+"s", ttl), cc.Desc()
	dep, depPlain := cc.DeprecatedHelpString(trans, CurrentDeprecatedColor, CurrentDescColor)

	if (cc.Hidden() && *verboseCount > 0) || (cc.VendorHidden() && *verboseCount >= 3) { //nolint:revive
		s.WriteBgColor(sb, color.BgDim)
	}

	if deprecated {
		s.WriteBgColor(sb, color.BgStrikeout)
		s.WriteColor(sb, CurrentDescColor)
	} else {
		s.WriteColor(sb, CurrentTitleColor)
	}
	_, _ = sb.WriteString(left)

	var split bool
	var printed int
	printleftpad := func(cond bool) {
		if cond {
			_, _ = sb.WriteString("\n")
			_, _ = sb.WriteString(strings.Repeat(" ", tabbedW))
		}
	}

	// s.DimFast(&sb, s.Translate(right, color.BgNormal))
	// s.ColoredFast(&sb, CurrentDescColor, s.Translate(right, CurrentDescColor))
	if right != "" {
		s.WriteColor(sb, CurrentDescColor)

		rCols := cols - tabbedW
		l := len(right) + len(depPlain)
		if l >= rCols {
			var prt string
			var ix int
			for len(right) > rCols {
				prt, right = right[:rCols], right[rCols:]
				printleftpad(ix > 0)
				_, _ = sb.WriteString(trans(prt, CurrentDescColor))
				ix++
			}
			if right != "" {
				if ix > 0 {
					printleftpad(ix > 0)
				} else {
					split, printed = true, len(right)
				}
				_, _ = sb.WriteString(trans(right, CurrentDescColor))
			}
		} else {
			_, _ = sb.WriteString(trans(right, CurrentDescColor))
		}
		// sb.WriteString(trans(right, CurrentDescColor))
	} else {
		s.WriteColor(sb, CurrentDescColor)
		_, _ = sb.WriteString(trans("<i>desc</i>", CurrentDescColor))
	}

	if dep != "" {
		if split {
			printleftpad(split)
			split = false
		}
		if printed >= 0 {
			_, _ = sb.WriteString(" ")
			_, _ = sb.WriteString(dep)
		}
		logz.VerboseContext(ctx, "[watching] split flag", "split", split)
	}

	s.Reset(sb) // reset fg/bg colors by color Translator
	_, _ = sb.WriteString("\n")
}

func (s *helpPrinter) printFlag(ctx context.Context, sb *strings.Builder, verboseCount *int, ff *cli.Flag, group string, idx, level, cols, tabbedW int, grouped bool) { //nolint:revive
	if (ff.Hidden() && *verboseCount < 1) || (ff.VendorHidden() && *verboseCount < 3) { //nolint:revive
		return
	}

	groupedInc := 0
	if grouped && group != "" {
		if idx == 0 {
			_, _ = sb.WriteString(strings.Repeat("  ", level+groupedInc))
			s.WriteColor(sb, CurrentGroupTitleColor)
			s.WriteBgColor(sb, CurrentGroupTitleBgColor)
			_, _ = sb.WriteString("[")
			_, _ = sb.WriteString(group)
			_, _ = sb.WriteString("]")
			s.Reset(sb)
			_, _ = sb.WriteString("\n")
		}
		groupedInc++
	}

	_, _ = sb.WriteString(strings.Repeat("  ", level+groupedInc))
	// ttl := strings.Join(ff.GetTitleZshFlagNamesArray(), ",")
	ttl := ff.GetTitleFlagNamesBy(",")
	w := tabbedW - (level+groupedInc)*2 // - len(ttl)

	deprecated := ff.Deprecated() != ""
	trans := func(ss string, clr color.Color) string {
		ss = s.Translate(strings.TrimSpace(ss), clr)
		if deprecated {
			ss = term.StripEscapes(ss)
		}
		return ss
	}

	left, right := fmt.Sprintf("%-"+strconv.Itoa(w)+"s", ttl), ff.Desc()
	tg := ff.ToggleGroupLeadHelpString()
	def, defPlain := ff.DefaultValueHelpString(trans, CurrentDefaultValueColor, CurrentDescColor)
	dep, depPlain := ff.DeprecatedHelpString(trans, CurrentDeprecatedColor, CurrentDescColor)

	if (ff.Hidden() && *verboseCount > 0) || (ff.VendorHidden() && *verboseCount >= 3) { //nolint:revive
		s.WriteBgColor(sb, color.BgDim)
	}

	if deprecated {
		s.WriteBgColor(sb, color.BgStrikeout)
		s.WriteColor(sb, CurrentDescColor)
	} else {
		s.WriteColor(sb, CurrentTitleColor)
	}
	_, _ = sb.WriteString(left)

	// s.DimFast(&sb, s.Translate(right, color.BgNormal))
	if tg != "" {
		// s.ColoredFast(&sb, CurrentFlagTitleColor, tg)
		s.WriteColor(sb, CurrentFlagTitleColor)
		_, _ = sb.WriteString(tg)
	}
	var split bool
	var printed int
	printleftpad := func(cond bool) {
		if cond {
			_, _ = sb.WriteString("\n")
			_, _ = sb.WriteString(strings.Repeat(" ", tabbedW))
		}
	}
	rCols := cols - tabbedW
	if right != "" {
		s.WriteColor(sb, CurrentDescColor)

		_, l, l1st := len(right), len(right)+len(defPlain)+len(depPlain), len(tg)
		// aa := []string{}
		if l+l1st >= rCols {
			var prt string
			var ix int
			for len(right)+l1st >= rCols {
				prt, right = right[:rCols-l1st], right[rCols-l1st:]
				printleftpad(ix > 0)
				// aa = append(aa, prt)
				_, _ = sb.WriteString(trans(prt, CurrentDescColor))
				ix++
				l1st = 0
			}
			if right != "" {
				str := trans(right, CurrentDescColor)
				if ix > 0 {
					printleftpad(ix > 0)
				} else {
					split, printed = true, len(is.StripEscapes(str))
				}
				_, _ = sb.WriteString(str)
			}
		} else {
			_, _ = sb.WriteString(trans(right, CurrentDescColor))
		}

		// if ff.Long == "addr" {
		// 	sb.WriteString(fmt.Sprintf(" / l=%d, l1st=%d, len(right)/def=%d/%d, rCols=%d", l, l1st, len(right), len(defPlain), rCols))
		// }
		// sb.WriteString(fmt.Sprintf(" / l=%d/%d, l1st=%d, len(right)=%d, rCols=%d", l, lr, l1st, len(right), rCols))
		// sb.WriteString(trans(right, CurrentDescColor))
	} else {
		s.WriteColor(sb, CurrentDescColor)
		_, _ = sb.WriteString(trans("<i>desc</i>", CurrentDescColor))
		printed += 4
	}

	if def != "" && printed >= 0 {
		if split {
			deflen := len(is.StripEscapes(def))
			printed += deflen
			if printed >= rCols {
				printleftpad(split)
			}
		}
		if sb.String()[sb.Len()-1] != ' ' {
			_, _ = sb.WriteString(" ")
		}
		_, _ = sb.WriteString(def)
	}

	if dep != "" {
		if split {
			deplen := len(is.StripEscapes(dep))
			printed += deplen
			if printed >= rCols {
				printleftpad(split)
			}
		}
		if sb.String()[sb.Len()-1] != ' ' {
			_, _ = sb.WriteString(" ")
		}
		_, _ = sb.WriteString(dep)
		logz.VerboseContext(ctx, "split flag is", "split", split)
	}
	// s.ColoredFast(&sb, CurrentDefaultValueColor, def)
	// s.ColoredFast(&sb, CurrentDeprecatedColor, dep)
	// sb.WriteString(s.Translate(right, color.BgDefault))
	s.Reset(sb)
	_, _ = sb.WriteString("\n")
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
