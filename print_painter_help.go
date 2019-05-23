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

func (s *helpPainter) fp(fmtStr string, args ...interface{}) {
	_, _ = fmt.Fprintf(rootCommand.ow, fmtStr+"\n", args...)
}

func (s *helpPainter) fpUsagesTitle(title string) {
	s.fp("\n%s:", title)
	// s.fp("\n\x1b[%dm\x1b[%dm%s\x1b[0m", bgNormal, darkColor, title)
	// fp("  [\x1b[%dm\x1b[%dm%s\x1b[0m]", bgDim, darkColor, normalize(group))
}

func (s *helpPainter) fpUsagesLine(fmt, appName, cmdList, cmdsTitle, tailPlaceHolder string) {
	s.fp("    %s %v%s%s [Options] [Parent/Global Options]"+fmt, appName, cmdList, cmdsTitle, tailPlaceHolder)
}

func (s *helpPainter) fpDescTitle(title string) {
	s.fp("\n%s:", title)
}

func (s *helpPainter) fpDescLine(desc string) {
	s.fp("    %v", desc)
}

func (s *helpPainter) fpExamplesTitle(title string) {
	s.fp("\n%s:", title)
}

func (s *helpPainter) fpExamplesLine(examples string) {
	for _, line := range strings.Split(examples, "\n") {
		s.fp("    %v", line)
	}
}

func (s *helpPainter) fpCommandsTitle(command *Command) {
	var title string
	if command.owner == nil {
		title = "Commands"
	} else {
		title = "Sub-Commands"
	}
	s.fp("\n%s:", title)
}

func (s *helpPainter) fpCommandsGroupTitle(group string) {
	if group != UnsortedGroup {
		// fp("  [%s]:", normalize(group))
		s.fp("  [\x1b[2m\x1b[%dm%s\x1b[0m]", currentGroupTitleColor, normalize(group))
	}
}

func (s *helpPainter) fpCommandsLine(command *Command) {
	if !command.Hidden {
		// s.fp("  %-48s%v", command.GetTitleNames(), command.Description)
		// s.fp("\n\x1b[%dm\x1b[%dm%s\x1b[0m", bgNormal, darkColor, title)
		// s.fp("  [\x1b[%dm\x1b[%dm%s\x1b[0m]", bgDim, darkColor, normalize(group))
		s.fp("  %-48s\x1b[%dm\x1b[%dm%s\x1b[0m", command.GetTitleNames(), bgDim, currentDescColor, command.Description)
	}
}

func (s *helpPainter) fpFlagsSssTitle(flag *Flag) {
	var title string
	if flag.owner == nil {
		title = "Commands"
	} else {
		title = "Sub-Commands"
	}
	s.fp("\n%s:", title)
}

func (s *helpPainter) fpFlagsTitle(title string) {
	s.fp("\n%s:", title)
}

func (s *helpPainter) fpFlagsGroupTitle(group string) {
	if group != UnsortedGroup {
		// fp("  [%s]:", normalize(group))
		// // echo -e "Normal \e[2mDim"
		// _, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m\x1b[2m\x1b[%dm[%04d]\x1b[0m%-44s \x1b[2m\x1b[%dm%s\x1b[0m ",
		// 	levelColor, levelText, darkColor, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message, darkColor, caller)
		s.fp("  [\x1b[2m\x1b[%dm%s\x1b[0m]", currentGroupTitleColor, normalize(group))
	}
}

func (s *helpPainter) fpFlagsLine(flg *Flag, defValStr string) {
	s.fp("  %-48s\x1b[%dm\x1b[%dm%s\x1b[%dm\x1b[%dm%s\x1b[0m",
		flg.GetTitleFlagNames(), bgNormal, currentDescColor, flg.Description,
		bgItalic, currentDefaultValueColor, defValStr)
}
