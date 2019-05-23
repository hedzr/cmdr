/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

type (
	Painter interface {
		// printHeader()
		// printHelpUsages(command)
		// printHelpDescription(command)
		// printHelpExamples(command)
		// printHelpSection(command, justFlags)
		// printHelpTailLine(command)
		fp(fmtStr string, args ...interface{})

		fpUsagesTitle(title string)
		fpUsagesLine(fmt, appName, cmdList, cmdsTitle, tailPlaceHolder string)
		fpDescTitle(title string)
		fpDescLine(desc string)
		fpExamplesTitle(title string)
		fpExamplesLine(examples string)
		fpCommandsTitle(command *Command)
		fpCommandsGroupTitle(group string)
		fpCommandsLine(command *Command)
		fpFlagsTitle(title string)
		fpFlagsGroupTitle(group string)
		fpFlagsLine(flag *Flag, defValStr string)
	}
)
