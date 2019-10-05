/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"os"
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
//
// Main program might decide whether it's a warning or error.
//
// See also
//
// [Flag.ValidArgs]
func WithIgnoreWrongEnumValue(ignored bool) ExecOption {
	return func(w *ExecWorker) {
		w.shouldIgnoreWrongEnumValue = ignored
		ShouldIgnoreWrongEnumValue = ignored
	}
}

// WithBuiltinCommands enables/disables those builtin predefined commands. Such as:
//
//  - versionsCmds / EnableVersionCommands supports injecting the default `--version` flags and commands
//  - helpCmds / EnableHelpCommands supports injecting the default `--help` flags and commands
//  - verboseCmds / EnableVerboseCommands supports injecting the default `--verbose` flags and commands
//  - generalCmdrCmds / EnableCmdrCommands support these flags: `--strict-mode`, `--no-env-overrides`
//  - generateCmds / EnableGenerateCommands supports injecting the default `generate` commands and subcommands
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
func WithHelpTabStop(tabStop int) ExecOption {
	return func(w *ExecWorker) {
		initTabStop(tabStop)
	}
}
