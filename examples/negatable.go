package examples

import (
	"context"

	"github.com/hedzr/cmdr/v2/cli"
)

func AddNegatableFlag(parent cli.CommandBuilder)           { AddNegatableFlagImpl(parent) }
func AddNegatableFlagWithoutCmd(parent cli.CommandBuilder) { AddNegatableFlagImpl(parent, true) }

func AddNegatableFlagImpl(parent cli.CommandBuilder, dontHandlingParentCmd ...bool) { //nolint:revive
	parent.Flg("warning", "w").
		Description("negatable flag: <code>--no-warning</code> is available", "").
		Group("Negatable").
		Negatable(true).
		LeadingPlusSign(true).
		Default(false).
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
			v := cmd.Store().MustBool("warning")
			w := cmd.Store().MustBool("no-warning")
			println("warning flag: ", v)
			println("no-warning flag: ", w)
			wf := cmd.FlagBy("warning")
			if wf != nil && wf.LeadingPlusSign() {
				println("warning flag with leading plus sign: +w FOUND.")
			}
			return
		})
	}

	parent.Examples(`Try to use negatable flag,
	  
	  $ $APP --no-warning
	    <code>cmd.Store().MustBool("no-warning")</code> will be 'true'.
	  $ $APP --warning
	    <code>cmd.Store().MustBool("warning")</code> will be 'true'.
	`)
}
