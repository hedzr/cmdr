/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/tool"
	"strconv"
	"strings"
)

type (
	helpPainter struct {
	}
)

func (s *helpPainter) Reset() {
}

func (s *helpPainter) Flush() {
}

func (s *helpPainter) Results() (res []byte) {
	return
}

func (s *helpPainter) Printf(fmtStr string, args ...interface{}) {
	_, _ = fmt.Fprintf(internalGetWorker().rootCommand.ow, fmtStr+"\n", args...)
}

func (s *helpPainter) FpPrintHeader(command *Command) {
	if len(command.root.Header) == 0 {
		s.Printf("%v by %v - v%v", command.root.Copyright, command.root.Author, command.root.Version)
	} else {
		s.Printf("%v", command.root.Header)
	}
}

func (s *helpPainter) FpPrintHelpTailLine(command *Command) {
	if internalGetWorker().enableHelpCommands {
		if GetNoColorMode() {
			s.Printf(fmtTailLineNC, internalGetWorker().helpTailLine)
		} else {
			s.Printf(fmtTailLine, CurrentGroupTitleColor, internalGetWorker().helpTailLine)
		}
	}
}

func (s *helpPainter) FpUsagesTitle(command *Command, title string) {
	s.Printf("\n%s:", title)
	// s.Printf("\n\x1b[%dm\x1b[%dm%s\x1b[0m", BgNormal, DarkColor, title)
	// fp("  [\x1b[%dm\x1b[%dm%s\x1b[0m]", BgDim, DarkColor, StripOrderPrefix(group))
}

func (s *helpPainter) FpUsagesLine(command *Command, fmt, appName, cmdList, cmdsTitle, tailPlaceHolder string) {
	if strings.HasPrefix(cmdList, appName) {
		appName = ""
	} else {
		cmdList = " " + cmdList
	}
	if len(tailPlaceHolder) > 0 {
		tailPlaceHolder = command.TailPlaceHolder
	} else {
		tailPlaceHolder = "[tail args...]"
	}
	s.Printf("    %s%v%s%s [Options] [Parent/Global Options]"+fmt, appName, cmdList, cmdsTitle, tailPlaceHolder)
}

func (s *helpPainter) FpDescTitle(command *Command, title string) {
	s.Printf("\n%s:", title)
}

func (s *helpPainter) FpDescLine(command *Command) {
	s.Printf("    %v", command.Description)
}

func (s *helpPainter) FpExamplesTitle(command *Command, title string) {
	s.Printf("\n%s:", title)
}

func (s *helpPainter) FpExamplesLine(command *Command) {
	str := tplApply(command.Examples, command.root)
	for _, line := range strings.Split(str, "\n") {
		s.Printf("    %v", line)
	}
}

func (s *helpPainter) FpCommandsTitle(command *Command) {
	var title string
	if command.owner == nil {
		title = "Commands"
	} else {
		title = "Sub-Commands"
	}
	s.Printf("\n%s:", title)
}

func (s *helpPainter) FpCommandsGroupTitle(group string) {
	if group != UnsortedGroup {
		if GetNoColorMode() {
			s.Printf(fmtCmdGroupTitleNC, tool.StripOrderPrefix(group))
		} else {
			s.Printf(fmtCmdGroupTitle, CurrentGroupTitleColor, tool.StripOrderPrefix(group))
		}
	}
}

func (s *helpPainter) FpCommandsLine(command *Command) {
	if !command.Hidden {
		if len(command.Deprecated) > 0 {
			if GetNoColorMode() {
				s.Printf(fmtCmdlineDepNC, command.GetTitleNames(), command.Description, command.Deprecated)
			} else {
				s.Printf(fmtCmdlineDep, BgNormal, CurrentDescColor, command.GetTitleNames(), command.Description, command.Deprecated)
			}
		} else {
			if GetNoColorMode() {
				s.Printf(fmtCmdlineNC, command.GetTitleNames(), command.Description)
			} else {
				// s.Printf("  %-48s%v", command.GetTitleNames(), command.Description)
				// s.Printf("\n\x1b[%dm\x1b[%dm%s\x1b[0m", BgNormal, DarkColor, title)
				// s.Printf("  [\x1b[%dm\x1b[%dm%s\x1b[0m]", BgDim, DarkColor, StripOrderPrefix(group))
				s.Printf(fmtCmdline, command.GetTitleNames(), BgNormal, CurrentDescColor, command.Description)
			}
		}
	}
}

// func (s *helpPainter) FpFlagsSssTitle(flag *Flag) {
// 	var title string
// 	if flag.owner == nil {
// 		title = "Commands"
// 	} else {
// 		title = "Sub-Commands"
// 	}
// 	s.Printf("\n%s:", title)
// }

func (s *helpPainter) FpFlagsTitle(command *Command, flag *Flag, title string) {
	s.Printf("\n%s:", title)
}

func (s *helpPainter) FpFlagsGroupTitle(group string) {
	if group != UnsortedGroup {
		if GetNoColorMode() {
			s.Printf(fmtGroupTitleNC, tool.StripOrderPrefix(group))
		} else {
			// fp("  [%s]:", StripOrderPrefix(group))
			// // echo -e "Normal \e[2mDim"
			// _, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m\x1b[2m\x1b[%dm[%04d]\x1b[0m%-48s \x1b[2m\x1b[%dm%s\x1b[0m ",
			// 	levelColor, levelText, DarkColor, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message, DarkColor, caller)
			s.Printf(fmtGroupTitle, CurrentGroupTitleColor, tool.StripOrderPrefix(group))
		}
	}
}

func (s *helpPainter) FpFlagsLine(command *Command, flg *Flag, maxShort int, defValStr string) {
	if len(flg.ValidArgs) > 0 {
		defValStr = fmt.Sprintf("%v, in %v", defValStr, flg.ValidArgs)
	}
	if flg.Min >= 0 && flg.Max > 0 {
		defValStr = fmt.Sprintf("%v, in [%v..%v]", defValStr, flg.Min, flg.Max)
	}
	var envKeys string
	if len(flg.EnvVars) > 0 {
		var sb strings.Builder
		for _, k := range flg.EnvVars {
			if len(strings.TrimSpace(k)) > 0 {
				sb.WriteString(strings.TrimSpace(k))
				sb.WriteRune(',')
			}
		}
		if sb.Len() > 0 {
			envKeys = fmt.Sprintf(" [env: %v]", strings.TrimRight(sb.String(), ","))
		}
	}
	if len(flg.Deprecated) > 0 {
		if GetNoColorMode() {
			s.Printf(fmtFlagsDepNC, // "  %-48s%s%s [deprecated since %v]",
				flg.GetTitleFlagNamesByMax(",", maxShort), flg.Description, envKeys, defValStr, flg.Deprecated)
		} else {
			s.Printf(fmtFlagsDep, // "  \x1b[%dm\x1b[%dm%-48s%s\x1b[%dm\x1b[%dm%s\x1b[0m [deprecated since %v]",
				BgNormal, CurrentDescColor, flg.GetTitleFlagNamesByMax(",", maxShort), flg.Description,
				BgItalic, CurrentDefaultValueColor, envKeys, defValStr, flg.Deprecated)
		}
	} else {
		if GetNoColorMode() {
			s.Printf(fmtFlagsNC, flg.GetTitleFlagNamesByMax(",", maxShort), flg.Description, envKeys, defValStr)
		} else {
			s.Printf(fmtFlags, // "  %-48s\x1b[%dm\x1b[%dm%s\x1b[%dm\x1b[%dm%s\x1b[0m",
				flg.GetTitleFlagNamesByMax(",", maxShort), BgNormal, CurrentDescColor, flg.Description,
				BgItalic, CurrentDefaultValueColor, envKeys, defValStr)
		}
	}
}

func initTabStop(ts int) {
	defaultTabStop = ts

	var s = strconv.Itoa(defaultTabStop)

	fmtCmdGroupTitle = "  [\x1b[2m\x1b[%dm%s\x1b[0m]"
	fmtCmdGroupTitleNC = "  [%s]"

	fmtCmdline = "  %-" + s + "s\x1b[%dm\x1b[%dm%s\x1b[0m"
	fmtCmdlineDep = "  \x1b[%dm\x1b[%dm%-" + s + "s%s\x1b[0m [deprecated since %v]"
	fmtCmdlineNC = "  %-" + s + "s%s"
	fmtCmdlineDepNC = "  %-" + s + "s%s [deprecated since %v]"

	fmtGroupTitle = "  [\x1b[2m\x1b[%dm%s\x1b[0m]"
	fmtGroupTitleNC = "  [%s]"

	fmtFlagsDep = "  \x1b[%dm\x1b[%dm%-" + s + "s%s\x1b[%dm\x1b[%dm%v%s\x1b[0m [deprecated since %v]"
	fmtFlags = "  %-" + s + "s\x1b[%dm\x1b[%dm%s\x1b[%dm\x1b[%dm%v%s\x1b[0m"
	fmtFlagsDepNC = "  %-" + s + "s%s%v%s [deprecated since %v]"
	fmtFlagsNC = "  %-" + s + "s%s%v%s"

	fmtTailLine = "\x1b[2m\x1b[%dm%s\x1b[0m"
	fmtTailLineNC = "%s"
}

var (
	defaultTabStop                                           = 48
	fmtCmdGroupTitle, fmtCmdGroupTitleNC                     string
	fmtCmdline, fmtCmdlineDep, fmtCmdlineNC, fmtCmdlineDepNC string
	fmtGroupTitle, fmtGroupTitleNC                           string
	fmtFlags, fmtFlagsDep, fmtFlagsNC, fmtFlagsDepNC         string
	fmtTailLine, fmtTailLineNC                               string
)

const defaultTailLine = `
Type '-h'/'-?' or '--help' to get command help screen. 
More: '-D'/'--debug'['--env'|'--raw'|'--more'], '-V'/'--version', '-#'/'--build-info', '--no-color', '--strict-mode', '--no-env-overrides'...`
