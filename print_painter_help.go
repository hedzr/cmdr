/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
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
	_, _ = fmt.Fprintf(uniqueWorker.rootCommand.ow, fmtStr+"\n", args...)
}

func (s *helpPainter) FpPrintHeader(command *Command) {
	if len(command.root.Header) == 0 {
		s.Printf("%v by %v - v%v", command.root.Copyright, command.root.Author, command.root.Version)
	} else {
		s.Printf("%v", command.root.Header)
	}
}

func (s *helpPainter) FpPrintHelpTailLine(command *Command) {
	if EnableHelpCommands {
		s.Printf("\nType '-h' or '--help' to get command help screen.")
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
	s.Printf("    %s%v%s%s [Options] [Parent/Global Options] [tail args...]"+fmt, appName, cmdList, cmdsTitle, tailPlaceHolder)
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
		if GetBoolR("no-color") {
			s.Printf(fmtCmdGroupTitleNC, StripOrderPrefix(group))
		} else {
			s.Printf(fmtCmdGroupTitle, CurrentGroupTitleColor, StripOrderPrefix(group))
		}
	}
}

func (s *helpPainter) FpCommandsLine(command *Command) {
	if !command.Hidden {
		if len(command.Deprecated) > 0 {
			if GetBoolR("no-color") {
				s.Printf(fmtCmdlineDepNC, command.GetTitleNames(), command.Description, command.Deprecated)
			} else {
				s.Printf(fmtCmdlineDep, BgNormal, CurrentDescColor, command.GetTitleNames(), command.Description, command.Deprecated)
			}
		} else {
			if GetBoolR("no-color") {
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
		if GetBoolR("no-color") {
			s.Printf(fmtGroupTitleNC, StripOrderPrefix(group))
		} else {
			// fp("  [%s]:", StripOrderPrefix(group))
			// // echo -e "Normal \e[2mDim"
			// _, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m\x1b[2m\x1b[%dm[%04d]\x1b[0m%-48s \x1b[2m\x1b[%dm%s\x1b[0m ",
			// 	levelColor, levelText, DarkColor, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message, DarkColor, caller)
			s.Printf(fmtGroupTitle, CurrentGroupTitleColor, StripOrderPrefix(group))
		}
	}
}

func (s *helpPainter) FpFlagsLine(command *Command, flg *Flag, defValStr string) {
	if len(flg.ValidArgs) > 0 {
		defValStr = fmt.Sprintf("%v, in %v", defValStr, flg.ValidArgs)
	}
	if flg.Min >= 0 && flg.Max > 0 {
		defValStr = fmt.Sprintf("%v, in [%v..%v]", defValStr, flg.Min, flg.Max)
	}
	if len(flg.Deprecated) > 0 {
		if GetBoolR("no-color") {
			s.Printf(fmtFlagsDepNC, // "  %-48s%s%s [deprecated since %v]",
				flg.GetTitleFlagNames(), flg.Description, defValStr, flg.Deprecated)
		} else {
			s.Printf(fmtFlagsDep, // "  \x1b[%dm\x1b[%dm%-48s%s\x1b[%dm\x1b[%dm%s\x1b[0m [deprecated since %v]",
				BgNormal, CurrentDescColor, flg.GetTitleFlagNames(), flg.Description,
				BgItalic, CurrentDefaultValueColor, defValStr, flg.Deprecated)
		}
	} else {
		if GetBoolR("no-color") {
			s.Printf(fmtFlagsNC, flg.GetTitleFlagNames(), flg.Description, defValStr)
		} else {
			s.Printf(fmtFlags, // "  %-48s\x1b[%dm\x1b[%dm%s\x1b[%dm\x1b[%dm%s\x1b[0m",
				flg.GetTitleFlagNames(), BgNormal, CurrentDescColor, flg.Description,
				BgItalic, CurrentDefaultValueColor, defValStr)
		}
	}
}

func initTabStop(ts int) {
	tabStop = ts

	var s = strconv.Itoa(tabStop)

	fmtCmdGroupTitle = "  [\x1b[2m\x1b[%dm%s\x1b[0m]"
	fmtCmdGroupTitleNC = "  [%s]"

	fmtCmdline = "  %-" + s + "s\x1b[%dm\x1b[%dm%s\x1b[0m"
	fmtCmdlineDep = "  \x1b[%dm\x1b[%dm%-" + s + "s%s\x1b[0m [deprecated since %v]"
	fmtCmdlineNC = "  %-" + s + "s%s"
	fmtCmdlineDepNC = "  %-" + s + "s%s [deprecated since %v]"

	fmtGroupTitle = "  [\x1b[2m\x1b[%dm%s\x1b[0m]"
	fmtGroupTitleNC = "  [%s]"

	fmtFlagsDep = "  \x1b[%dm\x1b[%dm%-" + s + "s%s\x1b[%dm\x1b[%dm%s\x1b[0m [deprecated since %v]"
	fmtFlags = "  %-" + s + "s\x1b[%dm\x1b[%dm%s\x1b[%dm\x1b[%dm%s\x1b[0m"
	fmtFlagsDepNC = "  %-" + s + "s%s%s [deprecated since %v]"
	fmtFlagsNC = "  %-" + s + "s%s%s"
}

var (
	tabStop                                                  = 48
	fmtCmdGroupTitle, fmtCmdGroupTitleNC                     string
	fmtCmdline, fmtCmdlineDep, fmtCmdlineNC, fmtCmdlineDepNC string
	fmtGroupTitle, fmtGroupTitleNC                           string
	fmtFlags, fmtFlagsDep, fmtFlagsNC, fmtFlagsDepNC         string
)
