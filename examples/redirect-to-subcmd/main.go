package main

import (
	"context"
	"os"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/cmd"
	logz "github.com/hedzr/logg/slog"
)

const (
	appName = "redirect-to-subcmd"
	desc    = `a sample to show u how to redirect from root cmd to
	a specific subcommand:
	when end-user type "$APP",
	the executor will be guided to "$APP wrong".`
	version = cmdr.Version
	author  = ``
)

func main() {
	app := cmdr.Create(appName, version, author, desc).With(func(app cli.App) {
		// redirect root commands into "wrong" subcmd for testing.
		//
		// For a dad command such as "server" command, it
		// would translate `app start|stop` -> `app server start|stop`.
		app.WithRootCommand(func(root *cli.RootCommand) {
			root.SetRedirectTo("wrong")
		})
	}).WithAdders(cmd.Commands...).Build()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}
