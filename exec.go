/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"github.com/hedzr/cmdr/conf"
	"os"
)

//

// ExecWorker is a core logic worker and holder
type ExecWorker struct {
	// beforeXrefBuildingX, afterXrefBuiltX HookFunc
	beforeXrefBuilding []HookFunc
	afterXrefBuilt     []HookFunc
	afterAutomaticEnv  []HookOptsFunc

	envPrefixes         []string
	rxxtPrefixes        []string
	predefinedLocations []string

	shouldIgnoreWrongEnumValue bool

	enableVersionCommands  bool
	enableHelpCommands     bool
	enableVerboseCommands  bool
	enableCmdrCommands     bool
	enableGenerateCommands bool

	doNotLoadingConfigFiles bool

	globalShowVersion   func()
	globalShowBuildInfo func()

	currentHelpPainter Painter

	defaultStdout *bufio.Writer
	defaultStderr *bufio.Writer

	// rootCommand the root of all commands
	rootCommand *RootCommand
	// rootOptions *Opt
	rxxtOptions *Options

	similarThreshold    float64
	noDefaultHelpScreen bool
	noColor             bool
	noEnvOverrides      bool
	strictMode          bool

	withLogex       bool
	afterArgsParsed func(cmd *Command, args []string) (err error)

	envvarToValueMap map[string]func() string
}

// ExecOption is the functional option for Exec()
type ExecOption func(w *ExecWorker)

//

var uniqueWorker *ExecWorker

func init() {
	InternalResetWorker()
}

//
//
// *******************************************
//
//

// Exec is main entry of `cmdr`.
func Exec(rootCmd *RootCommand, opts ...ExecOption) (err error) {
	w := uniqueWorker

	for _, opt := range opts {
		opt(w)
	}

	err = w.InternalExecFor(rootCmd, os.Args)
	return
}

// // InternalGetWorker is an internal helper, esp for debugging
// func InternalGetWorker() (w *ExecWorker) {
// 	w = uniqueWorker
// 	return
// }

// InternalResetWorker is an internal helper, esp for debugging
func InternalResetWorker() (w *ExecWorker) {
	uniqueWorker = &ExecWorker{
		envPrefixes:  []string{"CMDR"},
		rxxtPrefixes: []string{"app"},

		predefinedLocations: []string{
			"./ci/etc/%s/%s.yml",       // for developer
			"/etc/%s/%s.yml",           // regular location
			"/usr/local/etc/%s/%s.yml", // regular macOS HomeBrew location
			"$HOME/.config/%s/%s.yml",  // per user
			"$HOME/.%s/%s.yml",         // ext location per user
			"$THIS/%s.yml",             // executable's directory
			"%s.yml",                   // current directory
			// "./ci/etc/%s/%s.yml",
			// "/etc/%s/%s.yml",
			// "/usr/local/etc/%s/%s.yml",
			// "$HOME/.%s/%s.yml",
			// "$HOME/.config/%s/%s.yml",
		},

		shouldIgnoreWrongEnumValue: true,

		enableVersionCommands:  true,
		enableHelpCommands:     true,
		enableVerboseCommands:  true,
		enableCmdrCommands:     true,
		enableGenerateCommands: true,

		doNotLoadingConfigFiles: false,

		currentHelpPainter: new(helpPainter),

		defaultStdout: bufio.NewWriterSize(os.Stdout, 16384),
		defaultStderr: bufio.NewWriterSize(os.Stderr, 16384),

		rxxtOptions: NewOptions(),

		similarThreshold:    similarThreshold,
		noDefaultHelpScreen: false,
	}
	w = uniqueWorker
	WithEnvVarMap(nil)(w)
	return
}

// InternalExecFor is an internal helper, esp for debugging
func (w *ExecWorker) InternalExecFor(rootCmd *RootCommand, args []string) (err error) {
	var (
		pkg       = new(ptpkg)
		goCommand = &rootCmd.Command
		stop      bool
		matched   bool
	)

	if w.rootCommand == nil {
		w.setupRootCommand(rootCmd)
	}

	if len(conf.AppName) == 0 {
		conf.AppName = w.rootCommand.AppName
		conf.Version = w.rootCommand.Version
	}

	// initExitingChannelForFsWatcher()
	defer func() {
		// stop fs watcher explicitly
		stopExitingChannelForFsWatcher()

		if rootCmd.ow != nil {
			_ = rootCmd.ow.Flush()
		}
		if rootCmd.oerr != nil {
			_ = rootCmd.oerr.Flush()
		}
	}()

	err = w.preprocess(rootCmd, args)

	if err == nil {
		for pkg.i = 1; pkg.i < len(args); pkg.i++ {
			pkg.Reset()
			pkg.a = args[pkg.i]

			// --debug: long opt
			// -D:      short opt
			// -nv:     double chars short opt
			// ~~debug: long opt without opt-entry prefix.
			// ~D:      short opt without opt-entry prefix.
			// -abc:    the combined short opts
			// -nvabc, -abnvc: a,b,c,nv the four short opts, if no -n & -v defined.
			// --name=consul, --name consul, --name=consul: opt with a string, int, string slice argument
			// -nconsul, -n consul, -n=consul: opt with an argument.
			//  - -nconsul is not good format, but it could get somewhat works.
			//  - -n'consul', -n"consul" could works too.
			// -t3: opt with an argument.
			matched, stop, err = w.xxTestCmd(pkg, &goCommand, rootCmd, args)
			if e, ok := err.(*ErrorForCmdr); ok {
				ferr("%v", e)
				if !e.Ignorable {
					return
				}
			}
			if stop {
				if pkg.lastCommandHeld || (matched && pkg.flg == nil) {
					err = w.afterInternalExec(pkg, rootCmd, goCommand, args)
				}
				return
			}
		}

		err = w.afterInternalExec(pkg, rootCmd, goCommand, args)
	}

	return
}

func (w *ExecWorker) xxTestCmd(pkg *ptpkg, goCommand **Command, rootCmd *RootCommand, args []string) (matched, stop bool, err error) {
	if len(pkg.a) > 0 && (pkg.a[0] == '-' || pkg.a[0] == '/' || pkg.a[0] == '~') {
		if len(pkg.a) == 1 {
			pkg.needHelp = true
			pkg.needFlagsHelp = true
			return
		}

		// flag
		if stop, err = w.flagsPrepare(pkg, goCommand, args); stop || err != nil {
			return
		}
		if pkg.flg != nil && pkg.found {
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
		matched, stop, err = w.flagsMatching(pkg, cc, goCommand, args)

	} else {
		// testing the next command, but the last one has already been the end of command series.
		if pkg.lastCommandHeld {
			pkg.i--
			stop = true
			return
		}

		// or, keep going on...
		// if matched, stop, err = cmdMatching(pkg, goCommand, args); stop || err != nil {
		// 	return
		// }
		matched, stop, err = w.cmdMatching(pkg, goCommand, args)
	}
	return
}

func (w *ExecWorker) preprocess(rootCmd *RootCommand, args []string) (err error) {
	for _, x := range w.beforeXrefBuilding {
		x(rootCmd, args)
	}

	err = w.buildXref(rootCmd)

	if err == nil {
		err = w.rxxtOptions.buildAutomaticEnv(rootCmd)
	}

	if err == nil {
		for _, x := range w.afterXrefBuilt {
			x(rootCmd, args)
		}
	}
	return
}

func (w *ExecWorker) afterInternalExec(pkg *ptpkg, rootCmd *RootCommand, goCommand *Command, args []string) (err error) {
	w.checkState(pkg)

	if !pkg.needHelp && len(pkg.unknownCmds) == 0 && len(pkg.unknownFlags) == 0 {
		if goCommand.Action != nil {
			args := w.getArgs(pkg, args)

			if goCommand != &rootCmd.Command {
				if err = w.beforeInvokeCommand(rootCmd, goCommand, args); err == ErrShouldBeStopException {
					return nil
				}
			}

			if err = w.invokeCommand(rootCmd, goCommand, args); err == ErrShouldBeStopException {
				return nil
			}

			return
		}
	}

	// if GetIntP(getPrefix(), "help-zsh") > 0 || GetBoolP(getPrefix(), "help-bash") {
	// 	if len(goCommand.SubCommands) == 0 && !pkg.needFlagsHelp {
	// 		// pkg.needFlagsHelp = true
	// 	}
	// }

	if w.noDefaultHelpScreen == false {
		w.printHelp(goCommand, pkg.needFlagsHelp)
	}
	return
}

func (w *ExecWorker) checkState(pkg *ptpkg) {
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
}

func (w *ExecWorker) beforeInvokeCommand(rootCmd *RootCommand, goCommand *Command, args []string) (err error) {
	if rootCmd.PostAction != nil {
		defer rootCmd.PostAction(goCommand, args)
	}

	if w.withLogex {
		if err = w.initWithLogex(goCommand, args); err == ErrShouldBeStopException {
			return
		}
	}

	if w.afterArgsParsed != nil {
		if err = w.afterArgsParsed(goCommand, args); err == ErrShouldBeStopException {
			return
		}
	}

	if rootCmd.PreAction != nil {
		if err = rootCmd.PreAction(goCommand, args); err == ErrShouldBeStopException {
			return
		}
	}
	return
}

func (w *ExecWorker) invokeCommand(rootCmd *RootCommand, goCommand *Command, args []string) (err error) {
	if goCommand.PostAction != nil {
		defer goCommand.PostAction(goCommand, args)
	}

	if err = goCommand.Action(goCommand, args); err == ErrShouldBeStopException {
		return
	}
	return
}
