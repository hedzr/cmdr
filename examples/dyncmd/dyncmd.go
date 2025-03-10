package dyncmd

import (
	"context"

	"github.com/hedzr/cmdr/v2/cli"
)

// OnEvalJumpSubCommands querys shell scripts in EXT directory
// (typically it is `/usr/local/lib/<app-name>/ext/`) and build
// as subcommands dynamically.
//
// In this demo app, we looks up `./ci/pkg/usr.local.lib.large-app/ext`
// with hard-code.
//
// EXT directory: see the [cmdr.UsrLibDir()] for its location.
func OnEvalJumpSubCommands(ctx context.Context, c cli.Cmd) (it cli.EvalIterator, err error) {
	return onEvalJumpSubCommands(ctx, c)
}
