/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/conf"
	"os"
	"sort"
	"strings"
)

func fp(fmtStr string, args ...interface{}) {
	_, _ = fmt.Fprintf(uniqueWorker.rootCommand.ow, fmtStr+"\n", args...)
}

func ferr(fmtStr string, args ...interface{}) {
	_, _ = fmt.Fprintf(uniqueWorker.rootCommand.oerr, fmtStr+"\n", args...)
}

func (w *ExecWorker) printHelp(command *Command, justFlags bool) {
	initTabStop(tabStop)

	if GetIntR("help-zsh") > 0 {
		w.printHelpZsh(command, justFlags)
	} else if GetBoolR("help-bash") {
		// TODO for bash
		w.printHelpZsh(command, justFlags)
	} else {
		w.paintFromCommand(w.currentHelpPainter, command, justFlags)
	}

	// NOTE: checking `~~debug`
	if w.rxxtOptions.GetBoolEx("debug", false) {
		w.paintTildeDebugCommand()
	}
	if w.currentHelpPainter != nil {
		w.currentHelpPainter.Results()
		w.currentHelpPainter.Reset()

		w.paintFromCommand(nil, command, false) // for gocov testing
	}
}

// paintTildeDebugCommand for `~~debug`
func (w *ExecWorker) paintTildeDebugCommand() {
	if GetNoColorMode() {
		fp("\nDUMP:\n\n%v\n", w.rxxtOptions.DumpAsString())
	} else {
		// "  [\x1b[2m\x1b[%dm%s\x1b[0m]"
		fp("\n\x1b[2m\x1b[%dmDUMP:\n\n%v\x1b[0m\n", DarkColor, w.rxxtOptions.DumpAsString())

		if w.rxxtOptions.GetBoolEx("env") {
			fp("---- ENV: ")
			for _, s := range os.Environ() {
				s2 := strings.Split(s, "=")
				fp("  - %s = \x1b[2m\x1b[%dm%s\x1b[0m", s2[0], DarkColor, s2[1])
			}
		}
		if w.rxxtOptions.GetBoolEx("more") {
			fp("---- INFO: ")
			fp("Exec: \x1b[2m\x1b[%dm%s\x1b[0m, %s", DarkColor, GetExcutablePath(), GetExecutableDir())
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
		x.WriteString(fmt.Sprintf("%d: :((", GetIntP(w.getPrefix(), "help-zsh")))
		for _, cx := range command.SubCommands {
			for _, n := range cx.GetExpandableNamesArray() {
				x.WriteString(fmt.Sprintf(`%v:'%v' `, n, cx.Description))
			}

			// fp(`  %-25s  %v%v`, cx.GetName(), cx.GetQuotedGroupName(), cx.Description)

			// fp(`%v:%v`, cx.GetExpandableNames(), cx.Description)
			// printHelpZshCommands(cx)
		}
		x.WriteString("))")
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
	if len(w.rootCommand.Header) == 0 || !command.IsRoot() {
		p.FpUsagesTitle(command, "Usages")

		ttl := "[Commands] "
		if command.owner != nil {
			if len(command.SubCommands) == 0 {
				ttl = ""
			} else {
				ttl = "[Sub-Commands] "
			}
		}

		cmds := strings.ReplaceAll(uniqueWorker.backtraceCmdNames(command), ".", " ")
		if len(cmds) > 0 {
			cmds += " "
		}

		p.FpUsagesLine(command, "", w.rootCommand.Name, cmds, ttl, command.TailPlaceHolder)
	}
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

func (w *ExecWorker) printHelpSection(p Painter, command *Command, justFlags bool) {
	if !justFlags {
		printHelpCommandSection(p, command, justFlags)
	}
	printHelpFlagSections(p, command, justFlags)
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

func printHelpCommandSection(p Painter, command *Command, justFlags bool) {
	count := 0
	for _, items := range command.allCmds {
		count += len(items)
	}

	if count > 0 {
		p.FpCommandsTitle(command)

		k0 := getSortedKeysFromCmdGroupedMap(command.allCmds)
		for _, group := range k0 {
			groups := command.allCmds[group]
			if len(groups) > 0 {
				p.FpCommandsGroupTitle(group)
				for _, nm := range getSortedKeysFromCmdMap(groups) {
					p.FpCommandsLine(groups[nm])
				}
			}
		}
	}
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

func printHelpFlagSections(p Painter, command *Command, justFlags bool) {
	sectionName := "Options"

GO_PRINT_FLAGS:
	count := 0
	for _, items := range command.allFlags {
		count += len(items)
	}

	if count > 0 {
		p.FpFlagsTitle(command, nil, sectionName)
		k2 := getSortedKeysFromFlgGroupedMap(command.allFlags)
		for _, group := range k2 {
			groups := command.allFlags[group]
			if len(groups) > 0 {
				p.FpFlagsGroupTitle(group)
				for _, nm := range getSortedKeysFromFlgMap(groups) {
					flg := groups[nm]
					if !flg.Hidden {
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
						p.FpFlagsLine(command, flg, defValStr)
						// fp("  %-48s%v%s", flg.GetTitleFlagNames(), flg.Description, defValStr)
					}
				}
			}
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
		goto GO_PRINT_FLAGS
	}

}

func (w *ExecWorker) showVersion() {
	if w.globalShowVersion != nil {
		w.globalShowVersion()
		return
	}

	fp(`v%v
%v
%v
%v
%v`, conf.Version, conf.AppName, conf.Buildstamp, conf.Githash, conf.GoVersion)
}

func (w *ExecWorker) showBuildInfo() {
	if w.globalShowBuildInfo != nil {
		w.globalShowBuildInfo()
		return
	}

	w.printHeader(w.currentHelpPainter, &w.rootCommand.Command)
	// buildTime
	fp(`
       Built by: %v
Build Timestamp: %v
        Githash: %v`, conf.GoVersion, conf.Buildstamp, conf.Githash)
}
