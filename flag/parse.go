/*
 * Copyright © 2019 Hedzr Yeh.
 */

package flag

import (
	"github.com/hedzr/cmdr"
	"log"
	"os"
)

func init() {
	pfRootCmd = cmdr.Root("", "")
}

// Parse parses the command-line flags from os.Args[1:]. Must be called
// after all flags are defined and before flags are accessed by the program.
func Parse() {
	// Ignore errors; CommandLine is set for ExitOnError.
	// CommandLine.Parse(os.Args[1:])
	if err := cmdr.Exec(pfRootCmd.RootCommand(), cmdr.WithNoDefaultHelpScreen(true)); err != nil {
		log.Fatal(err)
	}
	parsed = true
}

// Parsed reports whether the command-line flags have been parsed.
func Parsed() bool {
	// return CommandLine.Parsed()
	return parsed
}

// Args returns the non-flag command-line arguments.
func Args() []string { return os.Args[1:] }

// TreatAsLongOpt treat name as long option name or short.
func TreatAsLongOpt(b bool) bool {
	treatAsLongOpt = b
	return b
}
