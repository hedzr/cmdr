/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log/dir"
)

// func fp00(args ...interface{}) {
//	if w := internalGetWorker().rootCommand.ow; w != nil {
//		_, _ = fmt.Fprint(w, args...)
//	} else {
//		_, _ = fmt.Printf(args...)
//	}
// }

// func fp0_(fmtStr string, args ...interface{}) {
//	if w := internalGetWorker().rootCommand.ow; w != nil {
//		_, _ = fmt.Fprintf(w, fmtStr, args...)
//	} else {
//		_, _ = fmt.Printf(fmtStr, args...)
//	}
// }

func fp(fmtStr string, args ...interface{}) {
	s := fmt.Sprintf(fmtStr, args...)
	needln := !strings.HasSuffix(s, "\n")
	_fpz(needln, s)
}

func fpK(fmtStr string, args ...interface{}) {
	s := fmt.Sprintf(fmtStr, args...)
	_fpz(false, s)
}

func _fpz(needln bool, s string) {
	var w io.Writer = os.Stdout
	if wkr := internalGetWorker(); wkr != nil {
		uniqueWorkerLock.RLock()
		defer uniqueWorkerLock.RUnlock()

		if wkr.rootCommand != nil {
			if wkr.rootCommand.ow == nil {
				return
			}
			w = wkr.rootCommand.ow
		}
		if w != nil {
			_fpzz(needln, s, w)
		}
	} else {
		_fpzz(needln, s, w)
	}
}

func _fpzz(needln bool, s string, w io.Writer) {
	_, _ = fmt.Fprintf(w, "%s", s)
	if needln {
		_, _ = fmt.Fprintln(w)
	}
}

func ffp(of io.Writer, fmtStr string, args ...interface{}) {
	s := fmt.Sprintf(fmtStr, args...)
	needln := !strings.HasSuffix(fmtStr, "\n")
	_fpzz(needln, s, of)
}

func ferr(fmtStr string, args ...interface{}) {
	var w io.Writer = os.Stderr
	if wkr := internalGetWorker(); wkr != nil {
		uniqueWorkerLock.RLock()
		defer uniqueWorkerLock.RUnlock()

		if wkr.rootCommand != nil {
			if wkr.rootCommand.oerr == nil {
				return
			}
			w = wkr.rootCommand.oerr
		}

		if w != nil {
			s := fmt.Sprintf(fmtStr, args...)
			needln := !strings.HasSuffix(fmtStr, "\n")
			_fpzz(needln, s, w)
		}
	} else {
		s := fmt.Sprintf(fmtStr, args...)
		needln := !strings.HasSuffix(fmtStr, "\n")
		_fpzz(needln, s, w)
	}
}

// fwrn print the warning message if InDebugging() is true
func fwrn(fmtStr string, args ...interface{}) {
	if InDebugging() /* || logex.GetTraceMode() */ && internalWarningMode {
		ferr(fmtStr, args...)
	}
}

var internalWarningMode bool

// flog prints information if InDebugging() is true.
//
// To enable it, building your project with '-tags=delve'.
//
// See also: https://hedzr.github.io/cmdr-docs/zh/cmdr/guide/Z81.helpers.html#func-indebugging
func flog(fmtStr string, args ...interface{}) {
	if InDebugging() /* || logex.GetTraceMode() */ {
		_, _ = fmt.Fprintf(os.Stderr, "\u001B[2m\u001B[2m"+fmtStr+"\u001B[0m\n", args...)
	}
}

// printInDevMode prints information only in developing time.
//
// If the main program has been built as a executable binary, we
// would assumed which is not in developing time.
// If GetDebugMode() is true, that's in developing time too.
func printInDevMode(fmtStr string, args ...interface{}) {
	if inDevelopingTime() {
		_, _ = fmt.Fprintf(os.Stdout, "\u001B[2m\u001B[2m"+fmtStr+"\u001B[0m\n", args...)
	}
}

func (w *ExecWorker) printHelp(command *Command, justFlags bool) {
	if len(w.afterHelpScreen) > 0 {
		defer func() {
			for _, c := range w.afterHelpScreen {
				c(w, w.currentHelpPainter, command, justFlags)
			}
		}()
	}
	for _, c := range w.beforeHelpScreen {
		c(w, w.currentHelpPainter, command, justFlags)
	}

	initTabStop(defaultTabStop)

	// fluent -nc --help-zsh 1   # print _arguments array for sub-commands level 1
	// fluent -nc --help-zsh 2   # print _arguments array for sub-commands level 2
	// fluent Ø ms tags, ...
	switch {
	case GetIntR("help-zsh") > 0:
		w.printHelpZsh(command, justFlags)
	case GetBoolR("help-bash"):
		w.printHelpZsh(command, justFlags)
	default:
		// NOTE: checking `~~debug`
		if w.rxxtOptions.GetBoolEx("debug", false) {
			w.paintTildeDebugCommand(w.rxxtOptions.GetBoolEx("value-type"))
		} else {
			w.paintFromCommand(w.currentHelpPainter, command, justFlags)
		}
	}

	if w.currentHelpPainter != nil {
		w.currentHelpPainter.Results()
		w.currentHelpPainter.Reset()

		w.paintFromCommand(nil, command, false) // for gocov testing
	}
}

// DebugOutputTildeInfo prints debug information like `~~debug`
func DebugOutputTildeInfo(showType bool) {
	internalGetWorker().paintTildeDebugCommand(showType)
}

// paintTildeDebugCommand for `~~debug`
func (w *ExecWorker) paintTildeDebugCommand(showType bool) {
	of := os.Stdout
	if fn := GetStringR("debug-output"); fn != "" {
		var err error
		if of, err = os.Create(fn); err == nil {
			defer func() {
				if of != nil {
					_ = of.Close()
				}
			}()
		}
	}
	if GetNoColorMode() {
		ffp(of, "\nDUMP:\n\n%v\n", w.rxxtOptions.DumpAsString(showType))
	} else {
		// "  [\x1b[2m\x1b[%dm%s\x1b[0m]"
		ffp(of, "\n\x1b[2m\x1b[%dmDUMP:\n\n%v\x1b[0m\n", DarkColor, w.rxxtOptions.DumpAsString(showType))

		if w.rxxtOptions.GetBoolEx("env") {
			ffp(of, "---- ENV: ")
			for _, s := range os.Environ() {
				s2 := strings.Split(s, "=")
				ffp(of, "  - %s = \x1b[2m\x1b[%dm%s\x1b[0m", s2[0], DarkColor, s2[1])
			}
		}
		// if w.rxxtOptions.GetBoolEx("yaml") {
		// 	ffp(of, "---- YAML Tree: ")
		// 	if str, err := AsYamlExt(); err != nil {
		// 		log.Errorf("BAD: %v", err)
		// 	} else {
		// 		ffp(of, "\x1b[2m\x1b[%dm%s\x1b[0m", DarkColor, str)
		// 	}
		// }
		if w.rxxtOptions.GetBoolEx("more") {
			ffp(of, "---- INFO: ")
			ffp(of, "Exec: \x1b[2m\x1b[%dm%s\x1b[0m, %s", DarkColor, dir.GetExecutablePath(), dir.GetExecutableDir())
		}
	}
}

func (w *ExecWorker) paintFromCommand(p Painter, command *Command, justFlags bool) {
	if p == nil {
		return
	}

	w.printHeader(p, command)

	w.printHelpUsages(p, command)
	w.printHelpDescription(p, command)
	w.printHelpExamples(p, command)
	w.printHelpSection(p, command, justFlags)

	w.printHelpTailLine(p, command)

	p.Flush()
}

func (w *ExecWorker) printHeader(p Painter, command *Command) {
	p.FpPrintHeader(command)
}

func (w *ExecWorker) printHelpTailLine(p Painter, command *Command) {
	p.FpPrintHelpTailLine(command)
}

func (w *ExecWorker) printHelpZsh(command *Command, justFlags bool) {
	if command == nil {
		command = &w.rootCommand.Command
	}

	w.printHelpZshCommands(command, justFlags)
}

func (w *ExecWorker) printHelpZshCommands(command *Command, justFlags bool) {
	if !justFlags {
		var x strings.Builder
		for _, cx := range command.SubCommands {
			for _, n := range cx.GetExpandableNamesArray() {
				x.WriteString(fmt.Sprintf("%v:'%v' \n", n, tool.EscapeCompletionTitle(cx.Description)))
			}
		}
		// x.WriteString(fmt.Sprintf("%d: :((", GetIntP(w.getPrefix(), "help-zsh")))
		// for _, cx := range command.SubCommands {
		//	for _, n := range cx.GetExpandableNamesArray() {
		//		x.WriteString(fmt.Sprintf(`%v:'%v' `, n, cx.Description))
		//	}
		//
		//	// fp(`  %-25s  %v%v`, cx.GetName(), cx.GetQuotedGroupName(), cx.Description)
		//
		//	// fp(`%v:%v`, cx.GetExpandableNames(), cx.Description)
		//	// printHelpZshCommands(cx)
		// }
		// x.WriteString("))")
		directive := 4
		x.WriteString(fmt.Sprintf(":%d\n", directive))
		fp("%v", x.String())
	} else {
		for _, flg := range command.Flags {
			// fp(`  %-25s  %v`,
			// 	// "--help", //
			// 	// flg.GetTitleZshFlagNames(" "),
			// 	flg.GetTitleZshFlagName(), flg.GetDescZsh())
			for _, ff := range flg.GetTitleZshFlagNamesArray() {
				// fp(`  %-25s  %v`, ff, flg.GetDescZsh())
				fp(`%s[%v]`, ff, flg.GetDescZsh())
				// fp(`%s[%v]:%v:`, ff, flg.GetDescZsh(), flg.DefaultValuePlaceholder)
			}
		}
		fp(`(: -)--help[Print usage]`)
		// fp(`  %-25s  %v`, "--help", "Print Usage")
	}
}

func (w *ExecWorker) printHelpUsages(p Painter, command *Command) {
	if w.rootCommand.Header == "" || !command.IsRoot() {
		p.FpUsagesTitle(command, "Usages")

		ttl := "[Commands] "
		if command.owner != nil {
			if len(command.SubCommands) == 0 {
				ttl = ""
			} else {
				ttl = "[Sub-Commands] "
			}
		}

		cmds := replaceAll(backtraceCmdNames(command, true), ".", " ")
		if len(cmds) > 0 {
			cmds += " "
		}

		ph := getph1st(command)
		p.FpUsagesLine(command, "", w.rootCommand.Name, cmds, ttl, ph)
	}
}

func getph1st(c *Command) string {
	if c.TailPlaceHolder != "" {
		return c.TailPlaceHolder
	}
	if len(c.TailPlaceHolders) > 0 {
		return c.TailPlaceHolders[0]
	}
	return ""
}

func getphtail(c *Command) string {
	if len(c.TailPlaceHolders) > 0 {
		return strings.Join(c.TailPlaceHolders[1:], " ")
	}
	return ""
}

func (w *ExecWorker) printHelpDescription(p Painter, command *Command) {
	if len(command.Description) > 0 {
		p.FpDescTitle(command, "Description")
		p.FpDescLine(command)
		// fp("\nDescription: \n    %v", command.Description)
	}
}

func (w *ExecWorker) printHelpExamples(p Painter, command *Command) {
	if len(command.Examples) > 0 {
		p.FpExamplesTitle(command, "Examples")
		p.FpExamplesLine(command)
		// fp("%v", command.Examples)
	}
}

func findMaxL(s1 []aSection, maxL int) int {
	for _, s := range s1 {
		if s.maxL > maxL {
			maxL = s.maxL
		}
	}
	return maxL
}

func findMaxL2(s2 []aGroupedSections, maxL int) int {
	for _, s1 := range s2 {
		for _, s := range s1.sections {
			if s.maxL > maxL {
				maxL = s.maxL
			}
		}
	}
	return maxL
}

func findMaxR(s1 []aSection, maxR int) int {
	for _, s := range s1 {
		if s.maxR > maxR {
			maxR = s.maxR
		}
	}
	return maxR
}

func findMaxR2(s2 []aGroupedSections, maxR int) int {
	for _, s1 := range s2 {
		for _, s := range s1.sections {
			if s.maxR > maxR {
				maxR = s.maxR
			}
		}
	}
	return maxR
}

func getTextPiece(str string, start, want int) (text, ending string) {
	var sb, tried strings.Builder
	src := []rune(str[start:])
	var tryEscape, tryAnsiColor bool
	var tryPos int
	type controls struct {
		seq    string
		pos    int
		esclen int //nolint:structcheck,unused //keep it
	}
	var escapeSeqs []controls
	for _, c := range src { //nolint:gosimple,staticcheck //no problem
		if c == '\x1b' {
			tryEscape, tryAnsiColor = true, false
			tryPos = len([]rune(sb.String()))
			tried.Reset()
			tried.WriteRune(c)
			continue
		}

		if tryEscape {
			if tryAnsiColor {
				if unicode.IsDigit(c) || c == ';' || c == ',' {
					tried.WriteRune(c)
					continue
				}
				if c == 'm' {
					tried.WriteRune(c)
					tryEscape, escapeSeqs = false, append(escapeSeqs, controls{pos: tryPos, seq: tried.String()})
					continue
				}
			} else if c == '[' {
				tried.WriteRune(c)
				tryAnsiColor = true
				continue
			}
			sb.WriteString(tried.String())
		}
		if len([]rune(sb.String())) >= want {
			break
		}
		sb.WriteRune(c)
	}

	var out strings.Builder
	outs := []rune(sb.String())
	var last int
	for _, cc := range escapeSeqs {
		out.WriteString(string(outs[last:cc.pos]))
		out.WriteString(cc.seq)
		last, ending = cc.pos, cc.seq
	}
	out.WriteString(string(outs[last:]))
	text = out.String()
	return
}

func (w *ExecWorker) prCommands(p Painter, command *Command, s1 []aSection, maxL, cols int) {
	if len(s1) > 0 {
		p.FpCommandsTitle(command)
		for _, s := range s1 {
			p.FpCommandsGroupTitle(s.title)
			fmtStrL, fmtStrR, fmtStrMR := fmt.Sprintf("%%-%dv", maxL+2), "%v\u001B[0m\n", fmt.Sprintf("%%%dv\x1b[%dm%%v%%v\u001B[0m\n", maxL+2, CurrentDeprecatedColor)
			for i, l := range s.bufLL {
				str := l.String()
				if tool.IsTtyEscaped(str) {
					strCleaned := tool.StripEscapes(str)
					fl := fmt.Sprintf("%%-%dv", maxL+2+(len(str)-len(strCleaned)))
					p.Print(fl, l.String())
				} else {
					p.Print(fmtStrL, l.String())
				}

				str = s.bufLR[i].String()
				// if len(str) > cols {
				ww := maxL + 2
				s2w := cols - ww
				if s2w < len(str) && !InTesting() {
					firstPiece, ending := getTextPiece(str, 0, s2w)
					p.Print(fmtStrR, firstPiece)
					for ix := len(firstPiece); ix < len(str); {
						rs, end := getTextPiece(str, ix, s2w)
						p.Print(fmtStrMR, " ", ending, rs)
						ending = end
						ix += len(rs)
					}
					// p.Print("ww, s2w, cols = %v, %v, %v\n", ww, s2w, cols)
				} else {
					p.Print(fmtStrR, str)
				}
			}
		}
	}
}

func (w *ExecWorker) prFlags(p Painter, command *Command, s2 []aGroupedSections, maxL, cols int) {
	for _, s1 := range s2 {
		if len(s1.sections) > 0 {
			p.FpFlagsTitle(command, nil, s1.title)
			for _, s := range s1.sections {
				// p.FpCommandsGroupTitle(s.title)
				p.FpFlagsGroupTitle(s.title, s.isToggleGroup)

				// fmtStr := fmt.Sprintf("%%-%dv%%v\n", maxL+2)
				// for i, l := range s.bufLL {
				//	p.Print(fmtStr, l.String(), s.bufLR[i].String())
				// }

				fmtStrL, fmtStrR, fmtStrMR := fmt.Sprintf("%%-%dv", maxL+2), "%v\u001B[0m\n", fmt.Sprintf("%%%dv\x1b[%dm%%v%%v\u001B[0m\n", maxL+2, CurrentDeprecatedColor)
				for i, l := range s.bufLL {
					str := l.String()
					if tool.IsTtyEscaped(str) {
						strCleaned := tool.StripEscapes(str)
						fl := fmt.Sprintf("%%-%dv", maxL+2+(len(str)-len(strCleaned)))
						p.Print(fl, l.String())
					} else {
						p.Print(fmtStrL, l.String())
					}

					str = s.bufLR[i].String()
					// if len(str) > cols {
					ww := maxL + 2
					s2w := cols - ww
					if s2w < len(str) && !InTesting() {
						firstPiece, ending := getTextPiece(str, 0, s2w)
						p.Print(fmtStrR, firstPiece)
						for ix := len(firstPiece); ix < len(str); {
							rs, end := getTextPiece(str, ix, s2w)
							p.Print(fmtStrMR, " ", ending, rs)
							ending = end
							ix += len(rs)
						}
						// p.Print("ww, s2w, cols = %v, %v, %v\n", ww, s2w, cols)
					} else {
						p.Print(fmtStrR, str)
					}
				}
			}
		}
	}
}

func (w *ExecWorker) printHelpSection(p Painter, command *Command, justFlags bool) {
	var (
		s1         []aSection
		s2         []aGroupedSections
		maxL, maxR int
	)

	if !justFlags {
		s1 = printHelpCommandSection(p, command, justFlags)
	}
	s2 = printHelpFlagSections(p, command, justFlags)

	maxL = findMaxL2(s2, findMaxL(s1, 0))

	cols, _ := tool.GetTtySize()
	if cols <= 0 || cols > 512 {
		// fmt.Printf("\n\ncols = %v, maxL = %v\n\n\n", cols, maxL)
		maxR = findMaxR2(s2, findMaxR(s1, 0))
		cols = maxL + maxR + 2
		if cols < 80 {
			cols = 80
		}
		// fmt.Printf("\n\ncols = %v, maxL = %v\n\n\n", cols, maxL)
	}
	w.prCommands(p, command, s1, maxL, cols)
	w.prFlags(p, command, s2, maxL, cols)

	// return
}

func getSortedKeysFromCmdGroupedMap(m map[string]map[string]*Command) (k0 []string) {
	k0 = make([]string, 0)
	for k := range m {
		if k != UnsortedGroup {
			k0 = append(k0, k)
		}
	}
	sort.Strings(k0)
	// k0 = append(k0, UnsortedGroup)
	k0 = append([]string{UnsortedGroup}, k0...)
	return
}

func getSortedKeysFromCmdMap(groups map[string]*Command) (k1 []string) {
	k1 = make([]string, 0)
	for k := range groups {
		k1 = append(k1, k)
	}
	sort.Strings(k1)
	return
}

// func getSortedKeysFromBaseOptMap(groups map[string]*BaseOpt) (k1 []string) {
//	k1 = make([]string, 0)
//	for k := range groups {
//		k1 = append(k1, k)
//	}
//	sort.Strings(k1)
//	return
// }

type aSection struct {
	title         string
	bufLL, bufLR  []bytes.Buffer
	maxL, maxR    int
	isToggleGroup bool
}

type aGroupedSections struct {
	title    string
	sections []aSection
}

func countOfCommandsItems(p Painter, command *Command, justFlags bool) (count int) {
	for _, items := range command.allCmds {
		for _, c := range items {
			if c.VendorHidden {
				continue
			}
			if !c.Hidden || getFlagHitCount(command, "verbose") > 1 {
				count++
			}
		}
	}
	return
}

func lpgrp(p Painter, group string, g map[string]*Command) (sections []aSection) {
	var section aSection
	section.title = group // [nm].GetTitleName()
	for _, nm := range getSortedKeysFromCmdMap(g) {
		bufL, bufR := p.FpCommandsLine(g[nm])
		if bufL.Len() > 0 && bufR.Len() > 0 {
			section.bufLL, section.bufLR = append(section.bufLL, bufL), append(section.bufLR, bufR)
			if section.maxL < bufL.Len() {
				section.maxL = bufL.Len()
			}
			if section.maxR < bufR.Len() {
				section.maxR = bufR.Len()
			}
		}
	}
	if section.maxL > 0 {
		sections = append(sections, section)
	}
	return
}

func printHelpCommandSection(p Painter, command *Command, justFlags bool) (sections []aSection) {
	count := countOfCommandsItems(p, command, justFlags)
	if count > 0 {
		k0 := getSortedKeysFromCmdGroupedMap(command.allCmds)
		for _, group := range k0 {
			g := command.allCmds[group]
			if len(g) > 0 {
				sections = append(sections, lpgrp(p, group, g)...)
				// var section aSection
				// section.title = group // [nm].GetTitleName()
				// for _, nm := range getSortedKeysFromCmdMap(g) {
				// 	bufL, bufR := p.FpCommandsLine(g[nm])
				// 	if bufL.Len() > 0 && bufR.Len() > 0 {
				// 		section.bufLL, section.bufLR = append(section.bufLL, bufL), append(section.bufLR, bufR)
				// 		if section.maxL < bufL.Len() {
				// 			section.maxL = bufL.Len()
				// 		}
				// 		if section.maxR < bufR.Len() {
				// 			section.maxR = bufR.Len()
				// 		}
				// 	}
				// }
				// if section.maxL > 0 {
				// 	sections = append(sections, section)
				// }
			}
		}
	}
	return
}

func getSortedKeysFromFlgGroupedMap(m map[string]map[string]*Flag) (k2 []string) {
	k2 = make([]string, 0)
	for k := range m {
		if k != UnsortedGroup {
			k2 = append(k2, k)
		}
	}
	sort.Strings(k2)
	k2 = append([]string{UnsortedGroup}, k2...)
	return
}

func getSortedKeysFromFlgMap(groups map[string]*Flag) (k3 []string) {
	k3 = make([]string, 0)
	for k := range groups {
		k3 = append(k3, k)
	}
	sort.Strings(k3)
	return
}

func findMaxShortLength(groups map[string]*Flag) (maxShort int) {
	for _, flg := range groups {
		// flg := groups[nm]
		if /*!flg.Hidden &&*/ maxShort < len(flg.Short) {
			maxShort = len(flg.Short)
		}
	}
	return
}

func countOfFlagsItems(p Painter, command *Command, justFlags bool) (count int) {
	for _, items := range command.allFlags {
		for _, c := range items {
			if c.VendorHidden {
				continue
			}
			if !c.Hidden || getFlagHitCount(command, "verbose") > 1 {
				count++
			}
		}
	}
	return
}

func printHelpFlagSectionsChild(p Painter, command *Command, groups map[string]*Flag, groupTitle string) (section aSection) {
	// p.FpFlagsGroupTitle(group)
	section.title = groupTitle
	k3 := getSortedKeysFromFlgMap(groups)
	maxShort := findMaxShortLength(groups)
	for _, nm := range k3 {
		flg := groups[nm]
		if flg.VendorHidden || (flg.Hidden && getFlagHitCount(command, "verbose") <= 1) {
			continue
		}

		defValStr := ""
		if flg.DefaultValue != nil {
			if ss, ok := flg.DefaultValue.(string); ok && len(ss) > 0 {
				if len(flg.DefaultValuePlaceholder) > 0 {
					defValStr = fmt.Sprintf(" (default %v='%s')", flg.DefaultValuePlaceholder, ss)
				} else {
					defValStr = fmt.Sprintf(" (default='%s')", ss)
				}
			} else {
				if len(flg.DefaultValuePlaceholder) > 0 {
					defValStr = fmt.Sprintf(" (default %v=%v)", flg.DefaultValuePlaceholder, flg.DefaultValue)
				} else {
					defValStr = fmt.Sprintf(" (default=%v)", flg.DefaultValue)
				}
			}
		}
		bufL, bufR := p.FpFlagsLine(command, flg, maxShort, defValStr)
		section.bufLL, section.bufLR = append(section.bufLL, bufL), append(section.bufLR, bufR)
		if section.maxL < bufL.Len() {
			section.maxL = bufL.Len()
		}
		if section.maxR < bufR.Len() {
			section.maxR = bufR.Len()
		}
		if flg.ToggleGroup != "" {
			section.isToggleGroup = true
		}
		// fp("  %-48s%v%s", flg.GetTitleFlagNames(), flg.Description, defValStr)
	}
	return
}

func printHelpFlagSections(p Painter, command *Command, justFlags bool) (aGroupedSectionsList []aGroupedSections) {
	sectionName := "Options"

GoPrintFlags:
	count := countOfFlagsItems(p, command, justFlags)
	if count > 0 {
		var gs aGroupedSections
		k2 := getSortedKeysFromFlgGroupedMap(command.allFlags)
		for _, group := range k2 {
			groups := command.allFlags[group]
			if len(groups) > 0 {
				section := printHelpFlagSectionsChild(p, command, groups, group)
				if section.maxL > 0 {
					gs.sections = append(gs.sections, section)
				}
			}
		}
		if len(gs.sections) > 0 {
			gs.title = sectionName
			aGroupedSectionsList = append(aGroupedSectionsList, gs)
		}
	}

	if command.owner != nil {
		command = command.owner
		// sectionName = "Parent/Global Options"
		if command.owner == nil {
			sectionName = "Global Options"
		} else {
			sectionName = fmt.Sprintf("Parent (`%v`) Options", command.GetTitleName())
		}
		goto GoPrintFlags
	}

	return
}

func (w *ExecWorker) showVersion(cmd *Command) {
	if w.globalShowVersion != nil {
		w.globalShowVersion()
		return
	}

	fp(`v%v
%v
%v
%v
%v
%v
%v
%v`,
		strings.TrimLeft(conf.Version, "v"), conf.AppName,
		conf.Buildstamp, conf.Githash,
		conf.GoVersion, conf.GitSummary,
		conf.Serial, conf.SerialString,
	)
}

func (w *ExecWorker) showBuildInfo(cmd *Command) {
	if w.globalShowBuildInfo != nil {
		w.globalShowBuildInfo()
		return
	}

	initTabStop(defaultTabStop)

	w.printHeader(w.currentHelpPainter, &w.rootCommand.Command)

	if conf.GoVersion != "" {
		fp(`           Built by: %v`, conf.GoVersion)
	}

	if conf.BuilderComments != "" {
		fp(`   Builder Comments: %v`, conf.BuilderComments)
	}

	if conf.Serial != "" {
		fp(`       Built Serial: %v`, conf.Serial)
	}
	if conf.SerialString != "" {
		fp(`Built Serial String: %v`, conf.SerialString)
	}

	ts := conf.Buildstamp
	if ts == "" {
		ts = time.Now().UTC().Format("")
	}
	dt, err := time.Parse("", ts)
	if err == nil {
		ts = dt.Format("")
	}
	fp(`    Built Timestamp: %v`, ts)

	fp(`
         Git Commit: %v
        Git Summary: %v
    Git Description: %v
`, conf.Githash, conf.GitSummary, conf.GitDesc)

	w.currentHelpPainter.FpPrintHelpTailLine(cmd)
}
