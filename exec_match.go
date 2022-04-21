/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"gopkg.in/hedzr/errors.v3"
	"os"
	"strings"
)

func (w *ExecWorker) switchCharMatching(pkg *ptpkg, goCommand **Command, args []string) (matched, stop bool, err error) {
	// pkg.needHelp = true
	// pkg.needFlagsHelp = true
	ra := (args)[pkg.i:]
	if len(ra) > 1 {
		ra = ra[1:]
	}

	stop = true
	pkg.lastCommandHeld = false
	pkg.needHelp = false
	pkg.needFlagsHelp = false
	if !w.noUseOnSwitchCharHitHandler && !w.inCompleting {
		if w.onSwitchCharHitHandler != nil {
			err = w.onSwitchCharHitHandler(*goCommand, pkg.a, ra)
		} else {
			err = defaultOnSwitchCharHit(*goCommand, pkg.a, ra)
		}
	} else {
		pkg.remainArgs = append(pkg.remainArgs, pkg.a)
	}
	return
}

func (w *ExecWorker) cmdMatching(pkg *ptpkg, goCommand **Command, args []string) (matched, stop bool, err error) {
	cc := *goCommand

	if cmd, ok := cc.plainCmds[pkg.a]; ok {
		cmd.strHit = pkg.a
		*goCommand = cmd
		matched = true
		flog("    -> command %q hit (a=%q, idx=%v)...", cmd.GetTitleName(), pkg.a, pkg.i)
		stop, err = w.cmdMatched(pkg, *goCommand, args)
		return
	}

	if len(cc.SubCommands) == 0 { // (*goCommand).Action != nil &&
		// the args remained are files, not sub-commands.
		pkg.i--
		pkg.lastCommandHeld = true
		pkg.iLastCommand = pkg.i
		return
	}

	if cc == &w.rootCommand.Command && pkg.aliasCommand != nil {
		matched, stop, err = w.cmdMatching(pkg, &pkg.aliasCommand, args)
		return
	}

	if w.treatUnknownCommandAsArgs {
		pkg.lastCommandHeld, stop = true, true
		return
	}

	flog("    . adding unknown command %q", pkg.a)
	pkg.unknownCmds = append(pkg.unknownCmds, pkg.a)
	unknownCommand(pkg, cc, args)
	return
}

func (w *ExecWorker) cmdMatched(pkg *ptpkg, goCommand *Command, args []string) (stop bool, err error) {
	pkg.iLastCommand = pkg.i

	w.hitCommands = append(w.hitCommands, goCommand)
	// goCommand.strHit = pkg.fn
	goCommand.times++

	if len(goCommand.SubCommands) == 0 { // (*goCommand).Action != nil &&
		// the args remained are files, not sub-commands.
		pkg.lastCommandHeld, stop = true, true
	}

	if goCommand.onMatched != nil {
		err = goCommand.onMatched(goCommand, args)
	}
	return
}

func (w *ExecWorker) flagsPrepare(pkg *ptpkg, goCommand **Command, args []string) (stop bool, err error) {
	if len(pkg.a) > 1 {
		if strings.Contains(w.switchCharset, pkg.a[1:2]) { // '--', '~~', '//'
			if len(pkg.a) == 2 {
				// disableParser = true // '--': ignore the following args // PassThrough hit!
				stop = true
				pkg.lastCommandHeld = false
				pkg.needHelp = false
				pkg.needFlagsHelp = false
				if w.inCompleting {
					pkg.remainArgs = append(pkg.remainArgs, pkg.a)
					return
				}
				ra := args[pkg.i:]
				if len(ra) > 0 {
					ra = ra[1:]
				}
				if w.onPassThruCharHitHandler != nil {
					err = w.onPassThruCharHitHandler(*goCommand, pkg.a, ra)
				} else {
					err = defaultOnPassThruCharHit(*goCommand, pkg.a, ra)
				}
				return
			}

			// long flag
			pkg.matchLongFlag(goCommand)
			return
		}

		// short flag
		var ok bool
		if ok, err = pkg.doExtractingHeadLikeFlag(goCommand, args); ok {
			return
		}
		pkg.doParseSuffix()
		pkg.matchShortFlag(goCommand)
	}
	return
}

func (w *ExecWorker) flagsMatching(pkg *ptpkg, cc *Command, goCommand **Command, args []string) (matched, stop bool, err error) {
	var upLevel bool
goUp:
	pkg.found = false
	if pkg.short {
		a := "-" + pkg.fn + pkg.savedFn
		flog("    .  . matching short flag for %q", a)
		if i := pkg.matchLongestShortFlag(cc, a, 1); i >= 0 {
			pkg.fn, pkg.savedFn = a[1:i], a[i:]
			pkg.flg, matched = cc.plainShortFlags[pkg.fn]
		}
	} else {
		flog("    .  . matching long flag for --%v", pkg.fn)
		matched = pkg.matchLongFlagsRecursively(cc, pkg.fn, 0) >= 0
	}

	if matched {
		if upLevel, stop, err = w.flagsCheckAndDoMatched(pkg, goCommand, args); upLevel {
			goto goUp
		}
		return
	}

	if cc.owner != nil {
		// match the flag within parent's flags set.
		cc = cc.owner
		goto goUp
	}

	if !pkg.assigned && pkg.short {
		// try matching 2-chars short opt
		if len(pkg.savedFn) > 0 {
			fnf := pkg.fn + pkg.savedFn
			pkg.fn, pkg.savedFn = fnf[0:2], fnf[2:]
			*goCommand = pkg.savedGoCommand
			if (*goCommand).owner != nil {
				goto goUp
			}
		}
	}

	if pkg.i == len(args)-1 && w.inCompleting {
		stop = true
		pkg.remainArgs = append(pkg.remainArgs, pkg.a)
		return
	}

	pkg.unknownFlags = append(pkg.unknownFlags, pkg.a)
	unknownFlag(pkg, *goCommand, args)

	return
}

func (w *ExecWorker) flagsCheckAndDoMatched(pkg *ptpkg, goCommand **Command, args []string) (upLevel, stop bool, err error) {
	if err = w.checkFlagCanBeHere(pkg); err == nil {
		if upLevel, stop, err = w.flagsMatched(pkg, *goCommand, args); stop || err != nil {
			return
		}
	}
	return
}

func (w *ExecWorker) flagsMatched(pkg *ptpkg, goCommand *Command, args []string) (upLevel, stop bool, err error) {
	w.hitFlags = append(w.hitFlags, pkg.flg)
	pkg.flg.times++
	pkg.flg.strHit = pkg.fn

	if err = pkg.tryExtractingValue(args); err != nil {
		var perr *ErrorForCmdr
		if errors.As(err, &perr) && !perr.Ignorable {
			stop = true
			return // matched ok, and value extracted ok, no further parsing needs
		}
	}

	if pkg.found {
		// if !GetBoolP(getPrefix(), "quiet") {
		// 	logrus.Debugf("-- flag '%v' hit, go ahead...", pkg.flg.GetTitleName())
		// }
		stop, err = w.handleHandlersForFlag(pkg, goCommand, args)
		if stop || err != nil {
			return
		}

		// if pkg.flg.Action != nil {
		//	if err = pkg.flg.Action(goCommand, w.tmpGetRemainArgs(pkg, args)); err == ErrShouldBeStopException {
		//		stop = true
		//		err = nil
		//		return
		//	} else if err != nil {
		//		return
		//	}
		// }

		if isBool(pkg.flg.DefaultValue) || isNil1(pkg.flg.DefaultValue) {
			flog("    .  . [tryToggleGroup] %q = %v", pkg.fn, pkg.val)
			pkg.tryToggleGroup()
		}

		if !pkg.assigned {
			if len(pkg.savedFn) > 0 && len(pkg.savedVal) == 0 {
				pkg.fn = pkg.savedFn[0:1]
				pkg.savedFn = pkg.savedFn[1:]
				// goto GO_UP
				upLevel = true
			}
		} else {
			flog("    .  . [value assigned] %q = %v", pkg.fn, pkg.val)
		}
	}
	return
}

func (w *ExecWorker) handleHandlersForFlag(pkg *ptpkg, goCommand *Command, args []string) (stop bool, err error) {
	handleIt := func(handler Handler) (stop bool, err error) {
		if handler != nil {
			remainArgs := w.tmpGetRemainArgs(pkg, args)
			if err = handler(goCommand, remainArgs); IsIgnorableError(err) && err != nil {
				stop, err = true, nil
			}
		}
		return
	}
	stop, err = handleIt(pkg.flg.onMatched)
	if stop || err != nil {
		return
	}
	stop, err = handleIt(pkg.flg.Action)
	return
}

func (w *ExecWorker) checkFlagCanBeHere(pkg *ptpkg) (err error) {
	if err = w.checkPrerequisites(pkg); err != nil {
		return
	}
	if err = w.checkDblTildeStatus(pkg); err != nil {
		return
	}
	return
}

func (w *ExecWorker) checkDblTildeStatus(pkg *ptpkg) (err error) {
	if pkg.flg.dblTildeOnly {
		if pkg.a[:2] != "~~" {
			err = errors.New("a Flag '~~%v' request double tilde prefix only", pkg.flg.GetTitleName())
		}
	}
	return
}

func (w *ExecWorker) checkPrerequisites(pkg *ptpkg) (err error) {
	if len(pkg.flg.prerequisites) > 0 {
		for _, longTitleOrDottedPath := range pkg.flg.prerequisites {
			var cc *Command
			if strings.Contains(longTitleOrDottedPath, ",") {
				cc = dottedPathToCommand(longTitleOrDottedPath, nil)
			}
			if cc != nil {
				if err = w.checkPrerequisitesForCmd(cc, pkg); err != nil {
					return
				}
			}
		}
	}
	return
}

func (w *ExecWorker) checkPrerequisitesForCmd(cc *Command, pkg *ptpkg) (err error) {
	for _, f := range cc.Flags {
		if f.times == 0 {
			err = errors.New("the matching Flag '-%v' needs prerequisites are present, but '-%v' missed",
				pkg.flg.GetTitleName(),
				f.GetTitleName())
			return
		}
	}
	return
}

// MatchAndTest _
func MatchAndTest(inputCommandlineWithoutArg0 string, opts ...ExecOption) (last *Command, err error) {
	return matchForTest(inputCommandlineWithoutArg0, opts...)
}

// match try parsing the input command-line, the result is the last hit *Command.
func matchForTest(inputCommandlineWithoutArg0 string, opts ...ExecOption) (last *Command, err error) {
	saved := internalGetWorker()
	savedUnknownOptionHandler := unknownOptionHandler
	rootCmd := saved.rootCommand

	defer func() {
		uniqueWorkerLock.Lock()
		uniqueWorker = saved
		unknownOptionHandler = savedUnknownOptionHandler
		uniqueWorkerLock.Unlock()
	}()

	// get a new ExecWorker instance
	w := internalResetWorkerNoLock()

	for _, opt := range opts {
		opt(w)
	}

	w.noDefaultHelpScreen = true
	w.noUnknownCmdTip = true
	w.noCommandAction = true
	unknownOptionHandler = emptyUnknownOptionHandler

	line := os.Args[0] + " " + inputCommandlineWithoutArg0
	last, err = w.InternalExecFor(rootCmd, strings.Split(line, " "))
	return
}
