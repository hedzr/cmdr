/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

func cmdMatching(pkg *ptpkg, goCommand **Command, args []string) (stop bool, err error) {
	// command, files
	if cmd, ok := (*goCommand).plainCmds[pkg.a]; ok {
		cmd.strHit = pkg.a
		*goCommand = cmd
		// logrus.Debugf("-- command '%v' hit, go ahead...", cmd.GetTitleName())
		stop, err = cmdMatched(pkg, *goCommand, args)
	} else {
		if (*goCommand).Action != nil && len((*goCommand).SubCommands) == 0 {
			// the args remained are files, not sub-commands.
			pkg.i--
			stop = true
			return
		}

		pkg.unknownCmds = append(pkg.unknownCmds, pkg.a)
		unknownCommand(pkg, *goCommand, args)
	}
	return
}

func cmdMatched(pkg *ptpkg, goCommand *Command, args []string) (stop bool, err error) {
	if goCommand.PreAction != nil {
		if err = goCommand.PreAction(goCommand, getArgs(pkg, args)); err == ErrShouldBeStopException {
			return false, nil
		}
	}
	return
}

func flagsPrepare(pkg *ptpkg, goCommand **Command, args []string) (stop bool, err error) {
	if len(pkg.a) > 1 && (pkg.a[1] == '-' || pkg.a[1] == '~') {
		if len(pkg.a) == 2 {
			// disableParser = true // '--': ignore the following args
			stop = true
			return
		}

		// long flag
		pkg.fn = pkg.a[2:]
		findValueAttached(pkg, &pkg.fn)
	} else {
		pkg.suffix = pkg.a[len(pkg.a)-1]
		if pkg.suffix == '+' || pkg.suffix == '-' {
			pkg.a = pkg.a[0 : len(pkg.a)-1]
		} else {
			pkg.suffix = 0
		}

		if i := matchShortFlag(pkg, *goCommand, pkg.a); i >= 0 {
			pkg.fn = pkg.a[1:i]
			pkg.savedFn = pkg.a[i:]
		} else {
			pkg.fn = pkg.a[1:2]
			pkg.savedFn = pkg.a[2:]
		}
		pkg.short = true
		findValueAttached(pkg, &pkg.savedFn)
	}
	return
}

func flagsMatching(pkg *ptpkg, cc *Command, goCommand **Command, args []string) (stop bool, err error) {
	var ok, upLevel bool
GO_UP:
	pkg.found = false
	if pkg.short {
		pkg.flg, ok = cc.plainShortFlags[pkg.fn]
	} else {
		ok = matchForLongFlags(pkg.fn, cc, pkg)
	}

	if ok {
		if upLevel, stop, err = flagsMatched(pkg, *goCommand, args); stop || err != nil {
			return
		}
		if upLevel {
			goto GO_UP
		}
	} else {
		if cc.owner != nil {
			// match the flag within parent's flags set.
			cc = cc.owner
			goto GO_UP
		}
		if !pkg.assigned && pkg.short {
			// try matching 2-chars short opt
			if len(pkg.savedFn) > 0 {
				fnf := pkg.fn + pkg.savedFn
				pkg.fn = fnf[0:2]
				pkg.savedFn = fnf[2:]
				*goCommand = pkg.savedGoCommand
				goto GO_UP
			}
		}

		pkg.unknownFlags = append(pkg.unknownFlags, pkg.a)
		unknownFlag(pkg, *goCommand, args)
	}
	return
}

func flagsMatched(pkg *ptpkg, goCommand *Command, args []string) (upLevel, stop bool, err error) {
	if err = tryExtractingValue(pkg, args); err != nil {
		stop = true
		return
	}

	if pkg.found {
		// if !GetBoolP(getPrefix(), "quiet") {
		// 	logrus.Debugf("-- flag '%v' hit, go ahead...", pkg.flg.GetTitleName())
		// }
		if pkg.flg.Action != nil {
			if err = pkg.flg.Action(goCommand, getArgs(pkg, args)); err == ErrShouldBeStopException {
				stop = true
				err = nil
				return
			}
		}
		if isBool(pkg.flg.DefaultValue) || isNil1(pkg.flg.DefaultValue) {
			toggleGroup(pkg)
		}

		if !pkg.assigned {
			if len(pkg.savedFn) > 0 && len(pkg.savedVal) == 0 {
				pkg.fn = pkg.savedFn[0:1]
				pkg.savedFn = pkg.savedFn[1:]
				// goto GO_UP
				upLevel = true
			}
		}
	}
	return
}

func matchForLongFlags(fn string, cc *Command, pkg *ptpkg) (ok bool) {
	var ln = len(fn)
	for ; ln > 1; ln-- {
		fn = pkg.fn[0:ln]
		pkg.flg, ok = cc.plainLongFlags[fn]
		if ok {
			if ln < len(pkg.fn) {
				pkg.val = pkg.fn[ln:]
				pkg.fn = fn
				pkg.assigned = true
			}
			break
		}
	}
	return
}
