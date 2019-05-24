/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

type (
	manPainter struct {
		writer io.Writer
		// buffer bufio.Writer
	}
)

func newManPainter() *manPainter {
	return &manPainter{
		writer: new(bytes.Buffer),
		// buffer:bufio.NewWriterSize(),
	}
}

func (s *manPainter) Results() (res []byte) {
	if bb, ok := s.writer.(*bytes.Buffer); ok {
		res = bb.Bytes()
	}
	return
}

func (s *manPainter) Flush() {
	// if bb, ok := s.writer.(*bytes.Buffer); ok {
	// 	_, _ = fmt.Fprintf(os.Stdout, "%v\n", bb.String())
	// }
}

func (s *manPainter) Printf(fmtStr string, args ...interface{}) {
	str := fmt.Sprintf(fmtStr, args...)
	str = strings.ReplaceAll(str, "-", `\-`)
	str = strings.ReplaceAll(str, "`cmdr`", `\fBcmdr\fP`)
	_, _ = s.writer.Write([]byte(str))
}

type manHdrData struct {
	RootCommand
	TimeMY      string
	ManExamples string
}

func (s *manPainter) FpPrintHeader(command *Command) {
	root := command.root
	a := &manHdrData{
		*root,
		time.Now().Format("Jan 2006"),
		manExamples(root.Examples, root),
	}

	s.Printf("%v", tplApply(`
.pc
.nh
.TH {{.AppName}} 1 "{{.TimeMY}}" "{{.Version}}" "Tool with cmdr"
Auto generated by hedzr/cmdr

.SH NAME
.PP
{{.AppName}} v{{.Version}} - {{.Copyright}}

`, a))

	if command.IsRoot() {
		s.Printf("%v", tplApply(`

.SH SYNOPSIS
.PP
\fB{{.AppName}} generate manual [flags]\fP

.SH DESCRIPTIONS
.PP
{{.LongDescription}}

.SH EXAMPLES

{{.ManExamples}}

.\" .SH TIPS
.\" 
.\" 	NAME, SYNOPSIS, CONFIGURATION, DESCRIPTION, OPTIONS, EXIT STATUS, RETURN VALUE, ERRORS, 
.\" 	ENVIRONMENT, FILES, VERSIONS, CONFORMING TO, NOTES, BUGS, EXAMPLE, AUTHORS, and SEE ALSO.

`, a))
	}
}

func (s *manPainter) FpPrintHelpTailLine(command *Command) {
	root := command.root
	s.Printf(`
.SH SEE ALSO
.PP
\fB%v(1)\fP

.SH HISTORY
.PP
%v Auto generated by hedzr/cmdr
`, root.AppName, time.Now().Format("02-Jan-2006")) // , time.RFC822Z
}

func (s *manPainter) FpUsagesTitle(command *Command, title string) {
	if !command.IsRoot() {
		s.Printf("\n.SH %s\n", "SYNOPSIS")
	}
	// s.Printf("\n\x1b[%dm\x1b[%dm%s\x1b[0m", bgNormal, darkColor, title)
	// fp("  [\x1b[%dm\x1b[%dm%s\x1b[0m]", bgDim, darkColor, normalize(group))
}

func (s *manPainter) FpUsagesLine(command *Command, fmt, appName, cmdList, cmdsTitle, tailPlaceHolder string) {
	if !command.IsRoot() {
		s.Printf(".PP\n\\fB%s\\fP %v%s%s [Options] [Parent/Global Options]"+fmt+"\n\n", appName, cmdList, cmdsTitle, tailPlaceHolder)
	}
}

func (s *manPainter) FpDescTitle(command *Command, title string) {
	if !command.IsRoot() {
		if len(command.LongDescription) > 0 {
			s.Printf("\n.SH %s\n", title)
		} else if len(command.Description) > 0 {
			s.Printf("\n.SH %s\n", title)
		}
	}
}

func (s *manPainter) FpDescLine(command *Command) {
	if !command.IsRoot() {
		if len(command.LongDescription) > 0 {
			s.Printf(".PP\n%v\n", manBr(command.LongDescription))
		} else if len(command.Description) > 0 {
			s.Printf(".PP\n%v\n", command.Description)
		}
	}
}

func (s *manPainter) FpExamplesTitle(command *Command, title string) {
	if !command.IsRoot() {
		if len(command.Examples) > 0 && command.HasParent() {
			s.Printf("\n.SH %s\n", title)
		}
	}
}

func (s *manPainter) FpExamplesLine(command *Command) {
	if !command.IsRoot() {
		if len(command.Examples) > 0 && command.HasParent() {
			s.Printf("%v\n", manExamples(command.Examples, command.root))
		}
	}
}

func (s *manPainter) FpCommandsTitle(command *Command) {
	var title string
	title = "COMMANDS AND SUB-COMMANDS"
	// if command.HasParent() {
	// 	title = "Commands"
	// } else {
	// 	title = "Sub-Commands"
	// }
	s.Printf("\n.SH %s\n", title)
}

func (s *manPainter) FpCommandsGroupTitle(group string) {
	if group != UnsortedGroup {
		// fp("  [%s]:", normalize(group))
		s.Printf(".SS \"%s\"\n", StripOrderPrefix(group))
	} else {
		s.Printf(".SS \"%s\"\n", "General")
	}
}

func (s *manPainter) FpCommandsLine(command *Command) {
	if !command.Hidden {
		// s.Printf("  %-48s%v", command.GetTitleNames(), command.Description)
		// s.Printf("\n\x1b[%dm\x1b[%dm%s\x1b[0m", bgNormal, darkColor, title)
		// s.Printf("  [\x1b[%dm\x1b[%dm%s\x1b[0m]", bgDim, darkColor, normalize(group))
		s.Printf(".TP\n.BI %s\n%s\n", manWs(command.GetTitleNames()), command.Description)
	}
}

func (s *manPainter) FpFlagsTitle(command *Command, flag *Flag, title string) {
	s.Printf("\n.SH %s\n", "OPTIONS")
}

func (s *manPainter) FpFlagsGroupTitle(group string) {
	if group != UnsortedGroup {
		s.Printf(".SS \"%s\"\n", StripOrderPrefix(group))
	} else {
		s.Printf(".SS \"%s\"\n", "General")
	}
}

func (s *manPainter) FpFlagsLine(command *Command, flag *Flag, defValStr string) {
	s.Printf(".TP\n.BI %s\n%s\n%s\n", manWs(flag.GetTitleFlagNames()), flag.Description, defValStr)
}

//
//
//
