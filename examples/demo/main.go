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
	appName = "demo-app"
	desc    = ``
	version = cmdr.Version
	author  = ``
)

func main() {
	app := cmdr.Create(appName, version, author, desc).With(func(app cli.App) {

	}).WithAdders(cmd.Commands...).Build()

	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}
