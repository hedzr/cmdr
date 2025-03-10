package main

import (
	"context"
	"os"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/examples/cmd"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

const (
	appName = "concise"
	desc    = `concise version of tiny app.`
	version = cmdr.Version
	author  = `The Example Authors`
)

func main() {
	ctx := context.Background()

	app := cmdr.Create(appName, version, author, desc).
		WithAdders(cmd.Commands...).
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
