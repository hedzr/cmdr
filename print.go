/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/conf"
	"regexp"
	"sort"
	"strings"
	"time"
)

func fp(fmtStr string, args ...interface{}) {
	_, _ = fmt.Fprintf(rootCommand.ow, fmtStr+"\n", args...)
}

func ferr(fmtStr string, args ...interface{}) {
	_, _ = fmt.Fprintf(rootCommand.oerr, fmtStr+"\n", args...)
}

func printHeader() {
	if len(rootCommand.Header) == 0 {
		fp("%v by %v - v%v", rootCommand.Copyright, rootCommand.Author, rootCommand.Version)
	} else {
		fp("%v", rootCommand.Header)
	}
}

func printHelp(command *Command, justFlags bool) {
	if GetInt("app.help-zsh") > 0 {
		printHelpZsh(command, justFlags)

	} else if GetBool("app.help-bash") {
		// TODO for bash
		printHelpZsh(command, justFlags)

	} else {

		printHeader()

		printHelpUsages(command)
		printHelpDescription(command)
		printHelpExamples(command)

		printHelpSection(command, justFlags)

		printHelpTailLine(command)

	}

	if rxxtOptions.GetBool("debug") {
		// "  [\x1b[2m\x1b[%dm%s\x1b[0m]"
		fp("\n\x1b[2m\x1b[%dmDUMP:\n\n%v\x1b[0m\n", darkColor, rxxtOptions.DumpAsString())
	}
}

func printHelpZsh(command *Command, justFlags bool) {
	if command == nil {
		command = &rootCommand.Command
	}

	printHelpZshCommands(command, justFlags)
}

func printHelpZshCommands(command *Command, justFlags bool) {
	if !justFlags {
		var x strings.Builder
		x.WriteString(fmt.Sprintf("%d: :((", GetInt("app.help-zsh")))
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

func printHelpUsages(command *Command) {
	if len(rootCommand.Header) == 0 {
		fp("\nUsages: ")

		ttl := "[Commands] "
		if command.owner != nil {
			if len(command.SubCommands) == 0 {
				ttl = ""
			} else {
				ttl = "[Sub-Commands] "
			}
		}

		cmds := strings.ReplaceAll(backtraceCmdNames(command), ".", " ")
		if len(cmds) > 0 {
			cmds += " "
		}
		fp("    %s %v%s%s [Options] [Parent/Global Options]", rootCommand.Name, cmds, ttl, command.TailPlaceHolder)
	}
}

func printHelpDescription(command *Command) {
	if len(command.Description) > 0 {
		fp("\nDescription: \n    %v", command.Description)
	}
}

func printHelpExamples(command *Command) {
	if len(command.Examples) > 0 {
		fp("%v", command.Examples)
	}
}

func printHelpTailLine(command *Command) {
	fp("\nType '-h' or '--help' to get command help screen.")
}

func printHelpSection(command *Command, justFlags bool) {
	if !justFlags {
		printHelpCommandSection(command, justFlags)
	}
	printHelpFlagSections(command, justFlags)
}

func printHelpCommandSection(command *Command, justFlags bool) {
	count := 0
	for _, items := range command.allCmds {
		count += len(items)
	}

	if count > 0 {
		if command.owner == nil {
			fp("\nCommands:")
		} else {
			fp("\nSub-Commands:")
		}
		k0 := make([]string, 0)
		for k := range command.allCmds {
			if k != UnsortedGroup {
				k0 = append(k0, k)
			}
		}
		sort.Strings(k0)
		// k0 = append(k0, UnsortedGroup)
		k0 = append([]string{UnsortedGroup}, k0...)

		for _, group := range k0 {
			groups := command.allCmds[group]
			if len(groups) > 0 {
				if group != UnsortedGroup {
					// fp("  [%s]:", normalize(group))
					fp("  [\x1b[2m\x1b[%dm%s\x1b[0m]", darkColor, normalize(group))
				}

				k1 := make([]string, 0)
				for k := range groups {
					k1 = append(k1, k)
				}
				sort.Strings(k1)

				for _, nm := range k1 {
					cmd := groups[nm]
					if !cmd.Hidden {
						fp("  %-48s%v", cmd.GetTitleNames(), cmd.Description)
					}
				}
			}
		}
	}
}

func printHelpFlagSections(command *Command, justFlags bool) {
	sectionName := "Options"

GO_PRINT_FLAGS:
	count := 0
	for _, items := range command.allFlags {
		count += len(items)
	}

	if count > 0 {
		fp("\n%v:", sectionName)
		k2 := make([]string, 0)
		for k := range command.allFlags {
			if k != UnsortedGroup {
				k2 = append(k2, k)
			}
		}
		sort.Strings(k2)
		k2 = append([]string{UnsortedGroup}, k2...)

		for _, group := range k2 {
			groups := command.allFlags[group]
			if len(groups) > 0 {
				if group != UnsortedGroup {
					// // echo -e "Normal \e[2mDim"
					// _, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m\x1b[2m\x1b[%dm[%04d]\x1b[0m%-44s \x1b[2m\x1b[%dm%s\x1b[0m ",
					// 	levelColor, levelText, darkColor, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message, darkColor, caller)
					fp("  [\x1b[2m\x1b[%dm%s\x1b[0m]", darkColor, normalize(group))
				}

				k3 := make([]string, 0)
				for k := range groups {
					k3 = append(k3, k)
				}
				sort.Strings(k3)

				for _, nm := range k3 {
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
						fp("  %-48s%v%s", flg.GetTitleFlagNames(), flg.Description, defValStr)
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

func normalize(s string) string {
	if xre.MatchString(s) {
		s = s[strings.Index(s, ".")+1:]
	}
	return s
}

// SetCustomShowVersion supports your `ShowVersion()` instead of internal `showVersion()`
func SetCustomShowVersion(fn func()) {
	globalShowVersion = fn
}

// SetCustomShowBuildInfo supports your `ShowBuildInfo()` instead of internal `showBuildInfo()`
func SetCustomShowBuildInfo(fn func()) {
	globalShowBuildInfo = fn
}

func showVersion() {
	if globalShowVersion != nil {
		globalShowVersion()
		return
	}

	fp("v%v", conf.Version)
	fp("%v", conf.AppName)
	fp("%v", conf.Buildstamp)
	fp("%v", conf.Githash)
}

func showBuildInfo() {
	if globalShowBuildInfo != nil {
		globalShowBuildInfo()
		return
	}

	printHeader()
	// buildTime
	fp("Build Timestamp: %v. Githash: %v", conf.Buildstamp, conf.Githash)
}

const (
	defaultTimestampFormat = time.RFC3339

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97

	boldOrBright = 1
	dim          = 2
	underline    = 4
	blink        = 5
	hidden       = 8

	darkColor = lightGray
)

var (
	xre = regexp.MustCompile(`^[0-9A-Za-z]+\.(.+)$`)
)
