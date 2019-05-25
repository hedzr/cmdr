/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

type (
	Painter interface {
		Printf(fmtStr string, args ...interface{})

		FpPrintHeader(command *Command)
		FpPrintHelpTailLine(command *Command)

		FpUsagesTitle(command *Command, title string)
		FpUsagesLine(command *Command, fmt, appName, cmdList, cmdsTitle, tailPlaceHolder string)
		FpDescTitle(command *Command, title string)
		FpDescLine(command *Command)
		FpExamplesTitle(command *Command, title string)
		FpExamplesLine(command *Command)

		FpCommandsTitle(command *Command)
		FpCommandsGroupTitle(group string)
		FpCommandsLine(command *Command)
		FpFlagsTitle(command *Command, flag *Flag, title string)
		FpFlagsGroupTitle(group string)
		FpFlagsLine(command *Command, flag *Flag, defValStr string)

		Flush()

		Results() []byte

		// clear any internal states and reset itself
		Reset()
	}
)
