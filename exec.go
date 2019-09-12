/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"fmt"
	"github.com/hedzr/cmdr/conf"
	"os"
	"strings"
)

//

//
type ExecWorker struct {
	// beforeXrefBuildingX, afterXrefBuiltX HookXrefFunc
	beforeXrefBuilding []HookXrefFunc
	afterXrefBuilt     []HookXrefFunc

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
}

// ExecOption is the functional option for Exec()
type ExecOption func(w *ExecWorker)

//

var uniqueWorker = &ExecWorker{
	envPrefixes:  []string{"CMDR"},
	rxxtPrefixes: []string{"app"},

	predefinedLocations: []string{
		"./ci/etc/%s/%s.yml",
		"/etc/%s/%s.yml",
		"/usr/local/etc/%s/%s.yml",
		os.Getenv("HOME") + "/.%s/%s.yml",
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
}

//

// WithXrefBuildingHooks sets the hook before and after building xref indices.
// It's replacers for AddOnBeforeXrefBuilding, and AddOnAfterXrefBuilt.
func WithXrefBuildingHooks(beforeXrefBuildingX, afterXrefBuiltX HookXrefFunc) func(w *ExecWorker) {
	return func(w *ExecWorker) {
		if beforeXrefBuildingX != nil {
			w.beforeXrefBuilding = append(w.beforeXrefBuilding, beforeXrefBuildingX)
		}
		if afterXrefBuiltX != nil {
			w.afterXrefBuilt = append(w.afterXrefBuilt, afterXrefBuiltX)
		}
	}
}

// WithEnvPrefix sets the environment variable text prefixes.
// cmdr will lookup envvars for a key.
func WithEnvPrefix(prefix []string) ExecOption {
	return func(w *ExecWorker) {
		w.envPrefixes = prefix
	}
}

// WithOptionsPrefix create a top-level namespace, which contains all normalized `Flag`s.
// =WithRxxtPrefix
func WithOptionsPrefix(prefix []string) ExecOption {
	return func(w *ExecWorker) {
		w.rxxtPrefixes = prefix
	}
}

// WithRxxtPrefix create a top-level namespace, which contains all normalized `Flag`s.
// cmdr will lookup envvars for a key.
func WithRxxtPrefix(prefix []string) ExecOption {
	return func(w *ExecWorker) {
		w.rxxtPrefixes = prefix
	}
}

// WithPredefinedLocations sets the environment variable text prefixes.
// cmdr will lookup envvars for a key.
func WithPredefinedLocations(locations []string) ExecOption {
	return func(w *ExecWorker) {
		w.predefinedLocations = locations
	}
}

// WithIgnoreWrongEnumValue will be put into `cmdrError.Ignorable` while wrong enumerable value found in parsing command-line options.
// main program might decide whether it's a warning or error.
// see also: [Flag.ValidArgs]
func WithIgnoreWrongEnumValue(ignored bool) ExecOption {
	return func(w *ExecWorker) {
		w.shouldIgnoreWrongEnumValue = ignored
		ShouldIgnoreWrongEnumValue = ignored
	}
}

// WithBuiltinCommands enables/disables those builtin predefined commands. Such as:
//
// 	- versionsCmds / EnableVersionCommands supports injecting the default `--version` flags and commands
// 	- helpCmds / EnableHelpCommands supports injecting the default `--help` flags and commands
// 	- verboseCmds / EnableVerboseCommands supports injecting the default `--verbose` flags and commands
// 	- generalCmdrCmds / EnableCmdrCommands support these flags: `--strict-mode`, `--no-env-overrides`
// 	- generateCmds / EnableGenerateCommands supports injecting the default `generate` commands and subcommands
//
func WithBuiltinCommands(versionsCmds, helpCmds, verboseCmds, generateCmds, generalCmdrCmds bool) ExecOption {
	return func(w *ExecWorker) {
		EnableVersionCommands = versionsCmds
		EnableHelpCommands = helpCmds
		EnableVerboseCommands = verboseCmds
		EnableCmdrCommands = generalCmdrCmds
		EnableGenerateCommands = generateCmds

		w.enableVersionCommands = versionsCmds
		w.enableHelpCommands = helpCmds
		w.enableVerboseCommands = verboseCmds
		w.enableCmdrCommands = generalCmdrCmds
		w.enableGenerateCommands = generateCmds
	}
}

// WithInternalOutputStreams sets the internal output streams for debugging
func WithInternalOutputStreams(out, err *bufio.Writer) ExecOption {
	return func(w *ExecWorker) {
		w.defaultStdout = out
		w.defaultStderr = err

		if w.defaultStdout == nil {
			w.defaultStdout = bufio.NewWriterSize(os.Stdout, 16384)
		}
		if w.defaultStderr == nil {
			w.defaultStderr = bufio.NewWriterSize(os.Stderr, 16384)
		}
	}
}

// SetCustomShowVersion supports your `ShowVersion()` instead of internal `showVersion()`
func WithCustomShowVersion(fn func()) ExecOption {
	return func(w *ExecWorker) {
		w.globalShowVersion = fn
	}
}

// SetCustomShowBuildInfo supports your `ShowBuildInfo()` instead of internal `showBuildInfo()`
func WithCustomShowBuildInfo(fn func()) ExecOption {
	return func(w *ExecWorker) {
		w.globalShowBuildInfo = fn
	}
}

// SetNoLoadConfigFiles true means no loading config files
func WithNoLoadConfigFiles(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.doNotLoadingConfigFiles = b
	}
}

// SetCurrentHelpPainter allows to change the behavior and facade of help screen.
func WithHelpPainter(painter Painter) ExecOption {
	return func(w *ExecWorker) {
		w.currentHelpPainter = painter
	}
}

// AddOnConfigLoadedListener add an functor on config loaded and merged
func WithConfigLoadedListener(c ConfigReloaded) ExecOption {
	return func(w *ExecWorker) {
		AddOnConfigLoadedListener(c)
	}
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

// InternalExecFor is an internal helper, esp for debugging
func (w *ExecWorker) InternalExecFor(rootCmd *RootCommand, args []string) (err error) {
	var (
		pkg       = new(ptpkg)
		goCommand = &rootCmd.Command
		stop      bool
		matched   bool
		// helpFlag = rootCmd.allFlags[UnsortedGroup]["help"]
	)

	// for deprecated variables
	//
	w.shouldIgnoreWrongEnumValue = ShouldIgnoreWrongEnumValue
	//
	w.enableVersionCommands = EnableVersionCommands
	w.enableHelpCommands = EnableHelpCommands
	w.enableVerboseCommands = EnableVerboseCommands
	w.enableCmdrCommands = EnableCmdrCommands
	w.enableGenerateCommands = EnableGenerateCommands
	//
	w.envPrefixes = EnvPrefix
	w.rxxtPrefixes = RxxtPrefix

	if rootCommand == nil {
		w.setRootCommand(rootCmd)
	}

	defer func() {
		_ = rootCmd.ow.Flush()
		_ = rootCmd.oerr.Flush()
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
		err = rxxtOptions.buildAutomaticEnv(rootCmd)
	}

	if err == nil {
		for _, x := range w.afterXrefBuilt {
			x(rootCmd, args)
		}
	}
	return
}

func (w *ExecWorker) afterInternalExec(pkg *ptpkg, rootCmd *RootCommand, goCommand *Command, args []string) (err error) {
	if !pkg.needHelp {
		pkg.needHelp = GetBoolP(w.getPrefix(), "help")
	}

	if !pkg.needHelp && len(pkg.unknownCmds) == 0 && len(pkg.unknownFlags) == 0 {
		if goCommand.Action != nil {
			args := w.getArgs(pkg, args)

			if goCommand != &rootCmd.Command {
				if rootCmd.PostAction != nil {
					defer rootCmd.PostAction(goCommand, args)
				}
				if rootCmd.PreAction != nil {
					if err = rootCmd.PreAction(goCommand, args); err == ErrShouldBeStopException {
						return nil
					}
				}
			}

			if goCommand.PostAction != nil {
				defer goCommand.PostAction(goCommand, args)
			}

			if err = goCommand.Action(goCommand, args); err == ErrShouldBeStopException {
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

	w.printHelp(goCommand, pkg.needFlagsHelp)
	return
}

func (w *ExecWorker) buildXref(rootCmd *RootCommand) (err error) {
	// build xref for root command and its all sub-commands and flags
	// and build the default values
	w.buildRootCrossRefs(rootCmd)

	if !w.doNotLoadingConfigFiles {
		// pre-detects for `--config xxx`, `--config=xxx`, `--configxxx`
		if err = w.parsePredefinedLocation(); err != nil {
			return
		}

		// and now, loading the external configuration files
		err = w.loadFromPredefinedLocation(rootCmd)

		if len(w.envPrefixes) > 0 {
			EnvPrefix = w.envPrefixes
		}
		w.envPrefixes = EnvPrefix
		envPrefix := strings.Split(GetStringR("env-prefix"), ".")
		if len(envPrefix) > 0 {
			w.envPrefixes = envPrefix
		}
	}
	return
}

func (w *ExecWorker) parsePredefinedLocation() (err error) {
	// pre-detects for `--config xxx`, `--config=xxx`, `--configxxx`
	if ix, str, yes := partialContains(os.Args, "--config"); yes {
		var location string
		if i := strings.Index(str, "="); i > 0 {
			location = str[i+1:]
		} else if len(str) > 8 {
			location = str[8:]
		} else if ix+1 < len(os.Args) {
			location = os.Args[ix+1]
		}

		location = trimQuotes(location)

		if len(location) > 0 && FileExists(location) {
			if yes, err = IsDirectory(location); yes {
				if FileExists(location + "/conf.d") {
					SetPredefinedLocations([]string{location + "/%s.yml"})
				} else {
					SetPredefinedLocations([]string{location + "/%s/%s.yml"})
				}
			} else if yes, err = IsRegularFile(location); yes {
				SetPredefinedLocations([]string{location})
			}
		}
	}
	return
}

func (w *ExecWorker) loadFromPredefinedLocation(rootCmd *RootCommand) (err error) {
	// and now, loading the external configuration files
	for _, s := range getExpandedPredefinedLocations() {
		fn := s
		switch strings.Count(fn, "%s") {
		case 2:
			fn = fmt.Sprintf(s, rootCmd.AppName, rootCmd.AppName)
		case 1:
			fn = fmt.Sprintf(s, rootCmd.AppName)
		}

		if FileExists(fn) {
			err = rxxtOptions.LoadConfigFile(fn)
			if err != nil {
				return
			}
			conf.CfgFile = fn
			break
		}
	}
	return
}

// AddOnBeforeXrefBuilding add hook func
func (w *ExecWorker) AddOnBeforeXrefBuilding(cb HookXrefFunc) {
	w.beforeXrefBuilding = append(w.beforeXrefBuilding, cb)
}

// AddOnAfterXrefBuilt add hook func
func (w *ExecWorker) AddOnAfterXrefBuilt(cb HookXrefFunc) {
	w.afterXrefBuilt = append(w.afterXrefBuilt, cb)
}

func (w *ExecWorker) setRootCommand(rootCmd *RootCommand) {
	rootCommand = rootCmd

	rootCommand.ow = w.defaultStdout
	rootCommand.oerr = w.defaultStderr
}

func (w *ExecWorker) getPrefix() string {
	return strings.Join(w.rxxtPrefixes, ".")
}

func (w *ExecWorker) getArgs(pkg *ptpkg, args []string) []string {
	var a []string
	if pkg.i+1 < len(args) {
		a = args[pkg.i+1:]
	}
	return a
}
