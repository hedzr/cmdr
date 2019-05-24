/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"strings"
)

type (
	helpPainter struct {
	}
)

func (s *helpPainter) Flush() {
}

func (s *helpPainter) Printf(fmtStr string, args ...interface{}) {
	_, _ = fmt.Fprintf(rootCommand.ow, fmtStr+"\n", args...)
}

func (s *helpPainter) FpPrintHeader(command *Command) {
	if len(command.root.Header) == 0 {
		s.Printf("%v by %v - v%v", command.root.Copyright, command.root.Author, command.root.Version)
	} else {
		s.Printf("%v", command.root.Header)
	}
}

func (s *helpPainter) FpPrintHelpTailLine(command *Command) {
	s.Printf("\nType '-h' or '--help' to get command help screen.")
}

func (s *helpPainter) FpUsagesTitle(command *Command, title string) {
	s.Printf("\n%s:", title)
	// s.Printf("\n\x1b[%dm\x1b[%dm%s\x1b[0m", BgNormal, DarkColor, title)
	// fp("  [\x1b[%dm\x1b[%dm%s\x1b[0m]", BgDim, DarkColor, StripOrderPrefix(group))
}

func (s *helpPainter) FpUsagesLine(command *Command, fmt, appName, cmdList, cmdsTitle, tailPlaceHolder string) {
	s.Printf("    %s %v%s%s [Options] [Parent/Global Options]"+fmt, appName, cmdList, cmdsTitle, tailPlaceHolder)
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
		// fp("  [%s]:", StripOrderPrefix(group))
		s.Printf("  [\x1b[2m\x1b[%dm%s\x1b[0m]", CurrentGroupTitleColor, StripOrderPrefix(group))
	}
}

func (s *helpPainter) FpCommandsLine(command *Command) {
	if !command.Hidden {
		// s.Printf("  %-48s%v", command.GetTitleNames(), command.Description)
		// s.Printf("\n\x1b[%dm\x1b[%dm%s\x1b[0m", BgNormal, DarkColor, title)
		// s.Printf("  [\x1b[%dm\x1b[%dm%s\x1b[0m]", BgDim, DarkColor, StripOrderPrefix(group))
		s.Printf("  %-48s\x1b[%dm\x1b[%dm%s\x1b[0m", command.GetTitleNames(), BgNormal, CurrentDescColor, command.Description)
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
		// fp("  [%s]:", StripOrderPrefix(group))
		// // echo -e "Normal \e[2mDim"
		// _, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m\x1b[2m\x1b[%dm[%04d]\x1b[0m%-44s \x1b[2m\x1b[%dm%s\x1b[0m ",
		// 	levelColor, levelText, DarkColor, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message, DarkColor, caller)
		s.Printf("  [\x1b[2m\x1b[%dm%s\x1b[0m]", CurrentGroupTitleColor, StripOrderPrefix(group))
	}
}

func (s *helpPainter) FpFlagsLine(command *Command, flg *Flag, defValStr string) {
	s.Printf("  %-48s\x1b[%dm\x1b[%dm%s\x1b[%dm\x1b[%dm%s\x1b[0m",
		flg.GetTitleFlagNames(), BgNormal, CurrentDescColor, flg.Description,
		BgItalic, CurrentDefaultValueColor, defValStr)
}
