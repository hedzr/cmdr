package main

import (
	"context"
	"os"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
	logz "github.com/hedzr/logg/slog"
)

const (
	appName = "external-editor"
	desc    = `a sample to show u how to use external editor flag.`
	version = cmdr.Version
	author  = `The Example Authors`
)

func main() {
	app := cmdr.Create(appName, version, author, desc).
		With(func(app cli.App) {
			app.OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
				msg := cmd.Store().MustString("message")
				println(`Hello, World.`, msg)
				return
			})
		}).
		WithBuilders(common.AddExternalEditorFlag).
		WithAdders().
		Build()

	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}
