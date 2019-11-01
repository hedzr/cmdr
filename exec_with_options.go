/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"github.com/hedzr/cmdr/conf"
	"os"
	"path"
)

// WithXrefBuildingHooks sets the hook before and after building xref indices.
// It's replacers for AddOnBeforeXrefBuilding, and AddOnAfterXrefBuilt.
func WithXrefBuildingHooks(beforeXrefBuildingX, afterXrefBuiltX HookFunc) ExecOption {
	return func(w *ExecWorker) {
		if beforeXrefBuildingX != nil {
			w.beforeXrefBuilding = append(w.beforeXrefBuilding, beforeXrefBuildingX)
		}
		if afterXrefBuiltX != nil {
			w.afterXrefBuilt = append(w.afterXrefBuilt, afterXrefBuiltX)
		}
	}
}

// WithAutomaticEnvHooks sets the hook after building automatic environment.
func WithAutomaticEnvHooks(hook HookOptsFunc) ExecOption {
	return func(w *ExecWorker) {
		if hook != nil {
			w.afterAutomaticEnv = append(w.afterAutomaticEnv, hook)
		}
	}
}

// WithEnvVarMap adds a (envvar-name, value) map, which will be applied
// to string option value, string-slice option values, ....
// For example, you could define a key-value entry in your `<app>.yml` file:
//    app:
//      test-value: "$THIS/$APPNAME.yml"
//      home-dir: "$HOME"
// it will be expanded by mapping to OS environment and this map (WithEnvVarMap).
// That is, $THIS will be expanded to the directory path of this
// executable, $APPNAME to the app name.
// And of course, $HOME will be mapped to os home directory path.
func WithEnvVarMap(varToValue map[string]func() string) ExecOption {
	return func(w *ExecWorker) {
		if varToValue == nil {
			varToValue = make(map[string]func() string)
		}
		w.envvarToValueMap = varToValue
		testAndSetMap(w.envvarToValueMap, "THIS", func() string { return GetExecutableDir() })
		testAndSetMap(w.envvarToValueMap, "APPNAME", func() string { return conf.AppName })
		testAndSetMap(w.envvarToValueMap, "CFG_DIR", func() string { return path.Dir(GetUsedConfigFile()) })
	}
}

func testAndSetMap(m map[string]func() string, key string, value func() string) {
	if _, ok := m[key]; !ok {
		m[key] = value
	}
}

// WithEnvPrefix sets the environment variable text prefixes.
// cmdr will lookup envvars for a key.
// Default env-prefix is array ["CMDR"], ie 'CMDR_'
func WithEnvPrefix(prefix ...string) ExecOption {
	return func(w *ExecWorker) {
		w.envPrefixes = prefix
	}
}

// WithOptionsPrefix create a top-level namespace, which contains all normalized `Flag`s.
// =WithRxxtPrefix
// Default Options Prefix is array ["app"], ie 'app.xxx'
func WithOptionsPrefix(prefix ...string) ExecOption {
	return func(w *ExecWorker) {
		w.rxxtPrefixes = prefix
	}
}

// WithRxxtPrefix create a top-level namespace, which contains all normalized `Flag`s.
// cmdr will lookup envvars for a key.
// Default Options Prefix is array ["app"], ie 'app.xxx'
func WithRxxtPrefix(prefix ...string) ExecOption {
	return func(w *ExecWorker) {
		w.rxxtPrefixes = prefix
	}
}

// WithPredefinedLocations sets the environment variable text prefixes.
// cmdr will lookup envvars for a key.
// Default locations are:
//
//     []string{
//       "./ci/etc/%s/%s.yml",       // for developer
//       "/etc/%s/%s.yml",           // regular location
//       "/usr/local/etc/%s/%s.yml", // regular macOS HomeBrew location
//       "$HOME/.config/%s/%s.yml",  // per user
//       "$HOME/.%s/%s.yml",         // ext location per user
//       "$THIS/%s.yml",             // executable's directory
//       "%s.yml",                   // current directory
//     },
//
// See also InternalResetWorker
func WithPredefinedLocations(locations ...string) ExecOption {
	return func(w *ExecWorker) {
		w.predefinedLocations = locations
	}
}

// WithIgnoreWrongEnumValue will be put into `cmdrError.Ignorable`
// while wrong enumerable value found in parsing command-line
// options.
// The default is true.
//
// Main program might decide whether it's a warning or error.
//
// See also
//
// [Flag.ValidArgs]
func WithIgnoreWrongEnumValue(ignored bool) ExecOption {
	return func(w *ExecWorker) {
		w.shouldIgnoreWrongEnumValue = ignored
		// ShouldIgnoreWrongEnumValue = ignored
	}
}

// WithBuiltinCommands enables/disables those builtin predefined commands. Such as:
//
//  - versionsCmds / EnableVersionCommands supports injecting the default `--version` flags and commands
//  - helpCmds / EnableHelpCommands supports injecting the default `--help` flags and commands
//  - verboseCmds / EnableVerboseCommands supports injecting the default `--verbose` flags and commands
//  - generalCmdrCmds / EnableCmdrCommands support these flags: `--strict-mode`, `--no-env-overrides`, and `--no-color`
//  - generateCmds / EnableGenerateCommands supports injecting the default `generate` commands and sub-commands
//
func WithBuiltinCommands(versionsCmds, helpCmds, verboseCmds, generateCmds, generalCmdrCmds bool) ExecOption {
	return func(w *ExecWorker) {
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

// WithCustomShowVersion supports your `ShowVersion()` instead of internal `showVersion()`
func WithCustomShowVersion(fn func()) ExecOption {
	return func(w *ExecWorker) {
		w.globalShowVersion = fn
	}
}

// WithCustomShowBuildInfo supports your `ShowBuildInfo()` instead of internal `showBuildInfo()`
func WithCustomShowBuildInfo(fn func()) ExecOption {
	return func(w *ExecWorker) {
		w.globalShowBuildInfo = fn
	}
}

// WithNoDefaultHelpScreen true to disable printing help screen if without any arguments
func WithNoDefaultHelpScreen(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.noDefaultHelpScreen = b
	}
}

// WithNoLoadConfigFiles true means no loading config files
func WithNoLoadConfigFiles(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.doNotLoadingConfigFiles = b
	}
}

// WithHelpPainter allows to change the behavior and facade of help screen.
func WithHelpPainter(painter Painter) ExecOption {
	return func(w *ExecWorker) {
		w.currentHelpPainter = painter
	}
}

// WithConfigLoadedListener add an functor on config loaded and merged
func WithConfigLoadedListener(c ConfigReloaded) ExecOption {
	return func(w *ExecWorker) {
		AddOnConfigLoadedListener(c)
	}
}

// WithHelpTabStop sets the tab-stop position in the help screen
// Default tabstop is 48
func WithHelpTabStop(tabStop int) ExecOption {
	return func(w *ExecWorker) {
		initTabStop(tabStop)
	}
}

// WithUnknownOptionHandler enables your customized wrong command/flag processor.
// internal processor supports smart suggestions for those wrong commands and flags.
func WithUnknownOptionHandler(handler UnknownOptionHandler) ExecOption {
	return func(w *ExecWorker) {
		unknownOptionHandler = handler
	}
}

// WithSimilarThreshold defines a threshold for command/option similar detector.
// Default threshold is 0.6666666666666666.
// See also JaroWinklerDistance
func WithSimilarThreshold(similiarThreshold float64) ExecOption {
	return func(w *ExecWorker) {
		w.similarThreshold = similiarThreshold
	}
}

// WithNoColor make console outputs plain and without ANSI escape colors
//
// Since v1.6.2+
func WithNoColor(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.noColor = b
	}
}

// WithNoEnvOverrides enables the internal no-env-overrides mode
//
// Since v1.6.2+
//
// In this mode, cmdr do NOT find and transfer equivalent envvar
// value into cmdr options store.
func WithNoEnvOverrides(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.noEnvOverrides = b
	}
}

// WithStrictMode enables the internal strict mode
//
// Since v1.6.2+
//
// In this mode, any warnings will be treat as an error and cause app
// fatal exit.
//
// In normal mode, these cases are assumed as warnings:
// - flag name not found
// - command or sub-command name not found
// - value extracting failed
// - ...
func WithStrictMode(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.strictMode = b
	}
}

// WithAfterArgsParsed sets a callback point after command-line args parsed by cmdr internal exec().
// Your callback func will be invoked before invoking the matched command `cmd`.
// At this time, all command-line args parsed and a command found.
//
// If program was launched with empty or wrong arguments, your callback func won't triggered.
//
// When empty argument or `--help` found, cmdr will display help screen. To customize it
// see also cmdr.WithCustomShowVersion and cmdr.WithCustomShowBuildInfo.
//
// When any wrong/warn arguments found, cmdr will display some tip message. To customize it
// see also cmdr.WithUnknownOptionHandler.
//
func WithAfterArgsParsed(hookFunc func(cmd *Command, args []string) (err error)) ExecOption {
	return func(w *ExecWorker) {
		w.afterArgsParsed = hookFunc
	}
}
