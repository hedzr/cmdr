package cli

import (
	"errors"
)

var (
	ErrShouldFallback   = errors.New("fallback")                                                                                      // cmdr should fall back to the default internal impl, such as defaultAction, ....
	ErrShouldStop       = errors.New("stop")                                                                                          // cmdr should end the parsing loop ASAP, instead old ErrShouldBeStopException.
	ErrEmptyRootCommand = errors.New("the RootCommand hasn't been built")                                                             // obs
	ErrCommandsNotReady = errors.New("the RootCommand hasn't been built, or InitGlobally() failed. Has builder.App.Build() called? ") // obs
)
