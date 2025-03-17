package main

import (
	"context"
	"os"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples"
	"github.com/hedzr/cmdr/v2/examples/devmode"
	logz "github.com/hedzr/logg/slog"
)

const (
	appName = "negatable"
	desc    = `a sample to demo negatable flag.`
	version = cmdr.Version
	author  = ``
)

func main() {
	app := cmdr.Create(appName, version, author, desc).
		With(func(app cli.App) { logz.Debug("in dev mode?", "mode", devmode.InDevelopmentMode()) }).
		WithBuilders(examples.AddNegatableFlag).
		WithAdders().
		// override the onAction defined in examples.AddNegatableFlag()
		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			// logz.PrintlnContext(ctx, cmd.Set().Dump())

			v := cmd.Store().MustBool("warning")
			w := cmd.Store().MustBool("no-warning")
			println("warning flag: ", v)
			println("no-warning flag: ", w)
			return
		}).
		Build()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}
