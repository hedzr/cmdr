/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"github.com/hedzr/log/closers"
	"github.com/hedzr/logex"
	"gopkg.in/hedzr/errors.v2"
	"os"
	"strings"
)

// Exec is main entry of `cmdr`.
//
//
func Exec(rootCmd *RootCommand, opts ...ExecOption) (err error) {
	defer func() {
		// stop fs watcher explicitly
		stopExitingChannelForFsWatcher()

		closers.Close()
	}()

	w := internalGetWorker()

	for _, opt := range opts {
		opt(w)
	}

	_, err = w.InternalExecFor(rootCmd, os.Args)
	if err == ErrShouldBeStopException {
		err = nil
	}
	return
}

// InternalExecFor is an internal helper, esp for debugging
func (w *ExecWorker) InternalExecFor(rootCmd *RootCommand, args []string) (last *Command, err error) {
	var pkg = new(ptpkg)

	if w.rootCommand == nil {
		w.setupRootCommand(rootCmd)
	}
	w.preparePtPkg(pkg)

	// initExitingChannelForFsWatcher()
	defer w.postExecFor(rootCmd, pkg)

	err = w.preprocess(rootCmd, args)
	if err == nil {
		w.rxxtOptions.setToAppendMode()
		last, err = w.internalExecFor(pkg, rootCmd, args)
	}
	return
}

func (w *ExecWorker) postExecFor(rootCmd *RootCommand, pkg *ptpkg) {
	// stop fs watcher explicitly
	// stopExitingChannelForFsWatcher()

	if rootCmd.ow != nil {
		_ = rootCmd.ow.Flush()
	}
	if rootCmd.oerr != nil {
		_ = rootCmd.oerr.Flush()
	}

	w.lastPkg = pkg

	if w.writeBackAlterConfigs {
		w.rxxtOptions.Flush()
	}
}

func (w *ExecWorker) preprocess(rootCmd *RootCommand, args []string) (err error) {
	flog("--> preprocess")
	for _, x := range w.beforeXrefBuilding {
		if x != nil {
			x(rootCmd, args)
		}
	}

	err = w.buildXref(rootCmd, args)

	if err == nil {
		flog("--> preprocess / rxxtOptions.buildAutomaticEnv()")
		err = w.rxxtOptions.buildAutomaticEnv(rootCmd)
	}

	flog("--> preprocess / rxxtOptions.setCB(onOptionMergingSet)")
	w.rxxtOptions.setCB(w.onOptionMergingSet, w.onOptionSet)

	if err == nil {
		flog("--> preprocess / afterXrefBuilt()")
		for _, x := range w.afterXrefBuilt {
			x(rootCmd, args)
		}
	}

	flog("--> preprocess / END: trace=%v/logex:%v, debug=%v/logex:%v, inDebugging:%v",
		GetTraceMode(), logex.GetTraceMode(), GetDebugMode(), logex.GetDebugMode(),
		logex.InDebugging())
	return
}

//func (w *ExecWorker) internalExecForV2(pkg *ptpkg, rootCmd *RootCommand, args []string) (last *Command, err error) {
//	var (
//		goCommand    = &rootCmd.Command
//		stopF, stopC bool
//		matched      bool
//	)
//
//	flog("--> process...")
//	for pkg.i = 1; pkg.i < len(args); pkg.i++ {
//		// if pkg.ResetAnd(args[pkg.i]) == 0 {
//		// 	continue
//		// }
//		lr := pkg.ResetAnd(args[pkg.i])
//		flog("--> parsing %q (idx=%v, len=%v) | pkg.lastCommandHeld=%v", pkg.a, pkg.i, lr, pkg.lastCommandHeld)
//
//		matched, stopC, stopF, err = w.xxTestCmd(pkg, &goCommand, rootCmd, &args)
//		if err != nil {
//			var e *ErrorForCmdr
//			if errors.As(err, &e) {
//				ferr("%v", e)
//				if !e.Ignorable {
//					return
//				}
//			}
//		}
//		if !matched {
//			pkg.remainArgs = append(pkg.remainArgs, pkg.a)
//		}
//		if stopF {
//			//if pkg.lastCommandHeld || (matched && pkg.flg == nil) {
//			//	pkg.remainArgs = append(pkg.remainArgs, pkg.a)
//			//}
//			break
//		}
//		if stopC && !matched {
//			break
//		}
//	}
//
//	last = goCommand
//	pkg.remainArgs = append(pkg.remainArgs, args[pkg.i:]...)
//	err = w.afterInternalExec(pkg, rootCmd, goCommand, args, stopC || pkg.lastCommandHeld)
//
//	return
//}

func (w *ExecWorker) shouldTerminate(err error) (shouldTerminate bool) {
	if err != nil {
		var e *ErrorForCmdr
		if errors.As(err, &e) {
			ferr("%v", e)
			if !e.Ignorable {
				shouldTerminate = true
			}
		} else {
			shouldTerminate = true
		}
	}
	return
}

func (w *ExecWorker) internalExecFor(pkg *ptpkg, rootCmd *RootCommand, args []string) (last *Command, err error) {
	var (
		goCommand    = &rootCmd.Command
		stopF, stopC bool
		matched      bool
	)

	flog("--> process...")
	for pkg.i = 1; pkg.i < len(args); pkg.i++ {
		// if pkg.ResetAnd(args[pkg.i]) == 0 {
		// 	continue
		// }
		lr := pkg.ResetAnd(args[pkg.i])
		flog("--> parsing %q (idx=%v, len=%v) | pkg.lastCommandHeld=%v", pkg.a, pkg.i, lr, pkg.lastCommandHeld)

		// --debug:        long opt
		// -D:             short opt
		// -nv:            double chars short opt, more chars are supported
		// ~~debug:        long opt without opt-entry prefix.
		// ~D:             short opt without opt-entry prefix.
		// -abc:           the combined short opts
		// -nvabc, -abnvc: a,b,c,nv the four short opts, if no -n & -v defined.
		// --name=consul, --name consul, --nameconsul: opt with a string, int, string slice argument
		// -nconsul, -n consul, -n=consul: opt with an argument.
		//  - -nconsul is not good format, but it could get somewhat works.
		//  - -n'consul', -n"consul" could works too.
		// -t3: opt with an argument.

		matched, stopC, stopF, err = w.xxTestCmd(pkg, &goCommand, rootCmd, &args)
		if w.shouldTerminate(err) {
			return
		}
		if stopF {
			if (matched && pkg.flg == nil) || (!matched && (pkg.lastCommandHeld || len((*goCommand).SubCommands) == 0)) {
				err = w.afterInternalExec(pkg, rootCmd, goCommand, args, stopC || pkg.lastCommandHeld)
			}
			return
		}
		if stopC && !matched && !pkg.lastCommandHeld {
			break
		}
	}

	last = goCommand
	err = w.afterInternalExec(pkg, rootCmd, goCommand, args, stopC || pkg.lastCommandHeld)
	return
}

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) xxTestCmd(pkg *ptpkg, goCommand **Command, rootCmd *RootCommand, args *[]string) (matched, stopC, stopF bool, err error) {

	if len(pkg.a) > 0 && strings.Contains(w.switchCharset, pkg.a[0:1]) { // pkg.a[0] == '/' ||
		if len(pkg.a) == 1 {
			matched, stopF, err = w.switchCharMatching(pkg, goCommand, *args)
			return
		}
		return w.xxTestCmdFlags(pkg, goCommand, rootCmd, args)
	}

	// testing the next command, but the last one has already been the end of command series.
	if pkg.lastCommandHeld {
		// if pkg.i == len(args) {	pkg.i-- }
		stopC, matched = true, true
		pkg.remainArgs = append(pkg.remainArgs, pkg.a)
		return
	}

	// or, keep going on...
	// if matched, stop, err = cmdMatching(pkg, goCommand, args); stop || err != nil {
	// 	return
	// }
	matched, stopC, err = w.cmdMatching(pkg, goCommand, *args)
	if matched && len((*goCommand).presetCmdLines) > 0 && (*goCommand).Invoke != "" {
		w.updateArgs(pkg, goCommand, rootCmd, args)
	}

	return
}

func (w *ExecWorker) xxTestCmdFlags(pkg *ptpkg, goCommand **Command, rootCmd *RootCommand, args *[]string) (matched, stopC, stopF bool, err error) {
	// flag
	if stopF, err = w.flagsPrepare(pkg, goCommand, *args); stopF || err != nil {
		return
	}
	if pkg.flg != nil && pkg.found { // if headLike flag matched
		matched = true
		return
	}

	// fn + val
	// fn: short,
	// fn: long
	// fn: short||val: such as '-t3'
	// fn: long=val, long='val', long="val", long val, long 'val', long "val"
	// fn: longval, long'val', long"val"

	pkg.savedGoCommand = *goCommand
	cc := *goCommand
	// if matched, stop, err = flagsMatching(pkg, cc, goCommand, args); stop || err != nil {
	// 	return
	// }
	flog("    -> matching flag for %q", pkg.a)

	if pkg.aliasCommand != nil && *goCommand == &rootCmd.Command {
		matched, stopF, err = w.flagsMatching(pkg, pkg.aliasCommand, &pkg.aliasCommand, *args)
		if matched || stopF || err != nil {
			*goCommand = pkg.flg.owner
			return
		}
	}

	matched, stopF, err = w.flagsMatching(pkg, cc, goCommand, *args)
	return
}

func (w *ExecWorker) updateArgs(pkg *ptpkg, goCommand **Command, rootCmd *RootCommand, args *[]string) {
	*args = append((*args)[0:pkg.i+1], append((*goCommand).presetCmdLines, (*args)[pkg.i+1:]...)...)
	cmdPath := (*goCommand).Invoke
	if cmd, matched := w.locateCommand(cmdPath, *goCommand); matched {
		*goCommand = cmd
	}
}

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) afterInternalExec(pkg *ptpkg, rootCmd *RootCommand, goCommand *Command, args []string, stopC bool) (err error) {

	flog("--> afterInternalExec: trace=%v/logex:%v, debug=%v/logex:%v, indebugging:%v",
		GetTraceMode(), logex.GetTraceMode(), GetDebugMode(), logex.GetDebugMode(),
		logex.InDebugging())

	w.checkStates(pkg)

	if !pkg.needHelp && len(pkg.unknownCmds) == 0 && len(pkg.unknownFlags) == 0 {
		if goCommand.Action == nil && goCommand == &rootCmd.Command && pkg.aliasCommand != nil {
			goCommand = pkg.aliasCommand
		}
		action := goCommand.Action
		if goCommand.Action == nil && assumeDefaultAction {
			action = defaultAction
		}
		if action != nil {
			rArgs := w.getRemainArgs(pkg, args)
			err = w.doInvokeCommand(rootCmd, action, goCommand, rArgs)
			return
		}
	}

	if w.noDefaultHelpScreen == false {
		rArgs := w.getRemainArgs(pkg, args)
		err = w.doInvokeHelpScreen(pkg, rootCmd, goCommand, rArgs)
	}
	return
}

func (w *ExecWorker) doInvokeHelpScreen(pkg *ptpkg, rootCmd *RootCommand, goCommand *Command, remainArgs []string) (err error) {

	defer w.deferRunPostActionOfRootLevel(rootCmd, goCommand, remainArgs)()
	err = w.runPreActionOfRootLevel(rootCmd, goCommand, remainArgs)
	if err == nil {
		w.printHelp(goCommand, pkg.needFlagsHelp)
	}
	return

}

func (w *ExecWorker) doInvokeCommand(rootCmd *RootCommand, action Handler, goCommand *Command, remainArgs []string) (err error) {
	if goCommand != &rootCmd.Command {
		if w.noCommandAction {
			return
		}

		defer w.deferRunPostActionOfRootLevel(rootCmd, goCommand, remainArgs)()

		if err = w.checkArgs(rootCmd, goCommand, remainArgs); err != nil {
			return
		}

		if err = w.runPreActionOfRootLevel(rootCmd, goCommand, remainArgs); err != nil {
			return
		}
	}

	if err = w.invokeCommand(rootCmd, action, goCommand, remainArgs); err == ErrShouldBeStopException {
		return nil
	}

	return
}

func (w *ExecWorker) deferRunPostActionOfRootLevel(rootCmd *RootCommand, goCommand *Command, remainArgs []string) (deferFn func()) {
	postActions := w.gatherPostActions(rootCmd)
	if len(postActions) > 0 {
		deferFn = func() {
			for _, fn := range postActions {
				if fn != nil {
					fn(goCommand, remainArgs)
				}
			}
		}
	} else {
		deferFn = func() {}
	}
	return
}

func (w *ExecWorker) runPreActionOfRootLevel(rootCmd *RootCommand, goCommand *Command, remainArgs []string) (err error) {
	var preActions = w.gatherPreActions(rootCmd)
	c := errors.NewContainer("cannot invoke preActions")
	for _, fn := range preActions {
		if fn != nil {
			switch e := fn(goCommand, remainArgs); {
			case e == ErrShouldBeStopException:
				return e
			case e != nil:
				c.Attach(e)
			}
		}
	}
	err = c.Error()
	return
}

func (w *ExecWorker) gatherPreActions(rootCmd *RootCommand) (preActions []Handler) {
	preActions = append(preActions, w.afterArgsParsed, rootCmd.PreAction)
	preActions = append(preActions, rootCmd.PreActions...)
	return
}

func (w *ExecWorker) gatherPostActions(rootCmd *RootCommand) (postActions []Invoker) {
	postActions = append(rootCmd.PostActions, rootCmd.PostAction)
	return
}

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) checkArgs(rootCmd *RootCommand, goCommand *Command, remainArgs []string) (err error) {

	if w.logexInitialFunctor != nil {
		err = w.logexInitialFunctor(goCommand, remainArgs)
		// ; err == ErrShouldBeStopException {
	}

	if err == nil {
		err = w.checkRequiredArgs(goCommand, remainArgs)
	}

	//if err == nil {
	//	err = w.checkPrerequisiteArgs(goCommand, remainArgs)
	//}

	if err == nil && w.afterArgsParsed != nil {
		err = w.afterArgsParsed(goCommand, remainArgs)
		//; err == ErrShouldBeStopException {
	}

	return
}

////goland:noinspection GoUnusedParameter
//func (w *ExecWorker) checkPrerequisiteArgs(goCommand *Command, remainArgs []string) (err error) {
//	c := errors.NewContainer("obligatory flag missed")
//	return
//}

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) checkRequiredArgs(goCommand *Command, remainArgs []string) (err error) {
	c := errors.NewContainer("required flag missed")

	cmd := goCommand
UP:
	for gn, gv := range cmd.allFlags {
		for fn, fv := range gv {
			if fv.Required && fv.times < 1 {
				c.Attach(errors.New("\n    The required flag %q in group %q missed", fn, gn))
			}
		}
	}
	if cmd.HasParent() {
		cmd = cmd.owner
		goto UP
	}

	err = c.Error()
	return
}

func (w *ExecWorker) checkStates(pkg *ptpkg) {
	if !pkg.needHelp {
		pkg.needHelp = GetBoolP(w.getPrefix(), "help")
	}

	if w.noColor {
		Set("no-color", true)
	}

	if w.noEnvOverrides {
		Set("no-env-overrides", true)
	}

	if w.strictMode {
		Set("strict-mode", true)
	}

	SetDebugMode(GetDebugMode()) // send state to hedzr/log
	SetTraceMode(GetTraceMode()) // send state to hedzr/log
	if GetDebugMode() {
		w.setLoggerLevel(DebugLevel)
	}
	if GetTraceMode() {
		w.setLoggerLevel(TraceLevel)
	}
}

// func (w *ExecWorker) beforeInvokeCommand(rootCmd *RootCommand, goCommand *Command, args []string) (err error) {
// 	if rootCmd.PostAction != nil {
// 		defer rootCmd.PostAction(goCommand, args)
// 	}
//
// 	if w.logexInitialFunctor != nil {
// 		if err = w.logexInitialFunctor(goCommand, args); err == ErrShouldBeStopException {
// 			return
// 		}
// 	}
//
// 	if w.afterArgsParsed != nil {
// 		if err = w.afterArgsParsed(goCommand, args); err == ErrShouldBeStopException {
// 			return
// 		}
// 	}
//
// 	if rootCmd.PreAction != nil {
// 		if err = rootCmd.PreAction(goCommand, args); err == ErrShouldBeStopException {
// 			return
// 		}
// 	}
// 	return
// }

//goland:noinspection GoUnusedParameter
func (w *ExecWorker) invokeCommand(rootCmd *RootCommand, action Handler, goCommand *Command, remainArgs []string) (err error) {
	if unhandledErrorHandler != nil {
		defer func() {
			// fmt.Println("defer caller")
			if ex := recover(); ex != nil {
				// debug.PrintStack()
				// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
				// dumpStacks()

				// https://stackoverflow.com/questions/52103182/how-to-get-the-stacktrace-of-a-panic-and-store-as-a-variable
				// http://hustcat.github.io/dive-into-stack-defer-panic-recover-in-go/
				// fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))

				// fmt.Printf("recover success. error: %v", ex)
				unhandledErrorHandler(ex)
				if e, ok := ex.(error); ok {
					err = e
				} else {
					err = errors.New("unexpected unknown error handled")
				}
			}
		}()
	}

	if goCommand.PostAction != nil {
		defer goCommand.PostAction(goCommand, remainArgs)
	}

	if goCommand.PreAction != nil {
		err = goCommand.PreAction(goCommand, remainArgs)
		// err != ErrShouldBeStopException
	}

	if err == nil {
		err = action(goCommand, remainArgs)
	}
	return
}

// func dumpStacks() {
// 	buf := make([]byte, 16384)
// 	buf = buf[:runtime.Stack(buf, true)]
// 	fmt.Printf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===\n", buf)
// }
