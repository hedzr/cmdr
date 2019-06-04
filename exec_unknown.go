/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

var (
	unknownOptionHandler func(isFlag bool, title string, cmd *Command, args []string)
)

// SetUnknownOptionHandler enables your customized wrong command/flag processor.
// internal processor supports smart suggestions for those wrong commands and flags.
func SetUnknownOptionHandler(handler func(isFlag bool, title string, cmd *Command, args []string)) {
	unknownOptionHandler = handler
}

func unknownCommand(pkg *ptpkg, cmd *Command, args []string) {
	ferr("\n\x1b[%dmUnknown command:\x1b[0m %v", BgBoldOrBright, pkg.a)
	if unknownOptionHandler != nil {
		unknownOptionHandler(false, pkg.a, cmd, args)
	} else {
		unknownCommandDetector(pkg, cmd, args)
	}
}

func unknownFlag(pkg *ptpkg, cmd *Command, args []string) {
	ferr("\n\x1b[%dmUnknown flag:\x1b[0m %v", BgBoldOrBright, pkg.a)
	if unknownOptionHandler != nil && !pkg.short {
		unknownOptionHandler(true, pkg.a, cmd, args)
	} else {
		unknownFlagDetector(pkg, cmd, args)
	}
}

func unknownCommandDetector(pkg *ptpkg, cmd *Command, args []string) {
	sndSrc := soundex(pkg.a)
	ever := false
	for k := range cmd.plainCmds {
		snd := soundex(k)
		if sndSrc == snd {
			ferr("  - do you mean: %v", k)
			ever = true
			// } else {
			// 	ferr("  . %v -> %v: --%v -> %v", pkg.a, sndSrc, k, snd)
		}
	}
	if !ever && cmd.HasParent() {
		unknownCommandDetector(pkg, cmd.GetOwner(), args)
	}
}

func unknownFlagDetector(pkg *ptpkg, cmd *Command, args []string) {
	sndSrc := soundex(pkg.a)
	if !pkg.short {
		ever := false
		for k := range cmd.plainLongFlags {
			snd := soundex(k)
			if sndSrc == snd {
				ferr("  - do you mean: --%v", k)
				ever = true
				// } else {
				// 	ferr("  . %v -> %v: --%v -> %v", pkg.a, sndSrc, k, snd)
			}
		}
		if !ever && cmd.HasParent() {
			unknownFlagDetector(pkg, cmd.GetOwner(), args)
		}
	}
}
