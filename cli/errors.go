package cli

import (
	errorsv3 "gopkg.in/hedzr/errors.v3"
)

var (
	// ErrUnmatchedCommand means Unmatched command found. It's just a state, not a real error, see [cli.Config.UnmatchedAsError]
	ErrUnmatchedCommand = errorsv3.New("UNKNOWN Command FOUND: %q | cmd=%v")
	// ErrUnmatchedFlag means Unmatched flag found. It's just a state, not a real error, see [cli.Config.UnmatchedAsError]
	ErrUnmatchedFlag = errorsv3.New("UNKNOWN Flag FOUND: %q | cmd=%v")
	// ErrRequiredFlag means required flag must be set explicitly
	ErrRequiredFlag = errorsv3.New("Flag %q is REQUIRED | cmd=%v")
)
