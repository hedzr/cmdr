/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import "os"

var(
	// EnableVersionCommands supports injecting the default `--version` flags and commands
	// Deprecated from v1.5.0
	EnableVersionCommands = true
	// EnableHelpCommands supports injecting the default `--help` flags and commands
	// Deprecated from v1.5.0
	EnableHelpCommands = true
	// EnableVerboseCommands supports injecting the default `--verbose` flags and commands
	// Deprecated from v1.5.0
	EnableVerboseCommands = true
	// EnableCmdrCommands support these flags: `--strict-mode`, `--no-env-overrides`
	// Deprecated from v1.5.0
	EnableCmdrCommands = true
	// EnableGenerateCommands supports injecting the default `generate` commands and subcommands
	// Deprecated from v1.5.0
	EnableGenerateCommands = true

	// EnvPrefix attaches a prefix to key to retrieve the option value.
	// Deprecated from v1.5.0
	EnvPrefix = []string{"CMDR"}

	// ShouldIgnoreWrongEnumValue will be put into `cmdrError.Ignorable` while wrong enumerable value found in parsing command-line options.
	// main program might decide whether it's a warning or error.
	// see also: [Flag.ValidArgs]
	// Deprecated from v1.5.0
	ShouldIgnoreWrongEnumValue = false
)

// AddOnBeforeXrefBuilding add hook func
// Deprecated from v1.5.0
func AddOnBeforeXrefBuilding(cb HookXrefFunc) {
	uniqueWorker.AddOnBeforeXrefBuilding(cb)
}

// AddOnAfterXrefBuilt add hook func
// Deprecated from v1.5.0
func AddOnAfterXrefBuilt(cb HookXrefFunc) {
	uniqueWorker.AddOnAfterXrefBuilt(cb)
}

// ExecWith is main entry of `cmdr`.
// Deprecated from v1.5.0
func ExecWith(rootCmd *RootCommand, beforeXrefBuildingX, afterXrefBuiltX HookXrefFunc) (err error) {
	w := uniqueWorker
	if beforeXrefBuildingX != nil {
		w.beforeXrefBuilding = append(w.beforeXrefBuilding, beforeXrefBuildingX)
	}
	if afterXrefBuiltX != nil {
		w.afterXrefBuilt = append(w.afterXrefBuilt, afterXrefBuiltX)
	}
	err = w.InternalExecFor(rootCmd, os.Args)
	return
}
