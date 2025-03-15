package cli

import (
	"errors"

	errorsv3 "gopkg.in/hedzr/errors.v3"
)

var (
	ErrShouldFallback   = errors.New("fallback")                                                                                      // cmdr should fall back to the default internal impl, such as defaultAction, ....
	ErrShouldStop       = errors.New("stop")                                                                                          // cmdr should end the parsing loop ASAP, instead old ErrShouldBeStopException.
	ErrEmptyRootCommand = errors.New("the RootCommand hasn't been built")                                                             // obs
	ErrCommandsNotReady = errors.New("the RootCommand hasn't been built, or InitGlobally() failed. Has builder.App.Build() called? ") // obs

	// ErrUnmatchedCommand means Unmatched command found. It's just a state, not a real error, see [cli.Config.UnmatchedAsError]
	ErrUnmatchedCommand = errorsv3.New("UNKNOWN CmdS FOUND: %q | cmd=%v")
	// ErrUnmatchedFlag means Unmatched flag found. It's just a state, not a real error, see [cli.Config.UnmatchedAsError]
	ErrUnmatchedFlag = errorsv3.New("UNKNOWN Flag FOUND: %q | cmd=%v")
	// ErrRequiredFlag means required flag must be set explicitly
	ErrRequiredFlag = errorsv3.New("Flag %q is REQUIRED | cmd=%v")
	ErrValidArgs    = errorsv3.New("Flag %q expects a valid input is in list: %v | cmd=%v")

	ErrMissedPrerequisite = errorsv3.New("Flag %q needs %q was set at first") // flag need a prerequisite flag exists.
	ErrFlagJustOnce       = errorsv3.New("Flag %q MUST BE set once only")     // flag cannot be set more than one time.
)
