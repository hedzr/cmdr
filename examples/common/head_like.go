package common

import (
	"context"

	"github.com/hedzr/cmdr/v2/cli"
)

func AddHeadLikeFlag(parent cli.CommandBuilder)           { AddHeadLikeFlagImpl(parent) }
func AddHeadLikeFlagWithoutCmd(parent cli.CommandBuilder) { AddHeadLikeFlagImpl(parent, true) }

func AddHeadLikeFlagImpl(parent cli.CommandBuilder, dontHandlingParentCmd ...bool) { //nolint:revive
	parent.Flg("lines", "l").
		Description("`head -1` like", "").
		Group("Head Like").
		HeadLike(true).
		Default(1).
		// Required(true).
		Build()

	var noHpc bool
	for _, p := range dontHandlingParentCmd {
		if p {
			noHpc = true
		}
	}
	if !noHpc {
		// give root command an action to handle it
		parent.OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			lines := cmd.Store().MustInt("lines")
			println("using lines: ", lines)
			return
		})
	}

	parent.Examples(`Try to use head-like app,
	  
	  $ $APP -567
	    this command request 567 lines just like "$APP --lines 567"
	`)
}
