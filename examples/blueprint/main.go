package main

// normally tiny app

import (
	"context"
	"os"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples"
	"github.com/hedzr/cmdr/v2/examples/blueprint/cmd"
	"github.com/hedzr/cmdr/v2/examples/devmode"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

const (
	appName = "blueprint"
	desc    = `a good blueprint for you.`
	version = cmdr.Version
	author  = ``
)

func main() {
	app := cmdr.Create(appName, version, author, desc,
		cmdr.WithAutoEnvBindings(true),

		cmdr.WithTasksBeforeRun(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) {
			logz.DebugContext(ctx, "command running...", "cmd", cmd, "runner", runner, "extras", extras)
			return
		}), // cmdr.WithTasksBeforeParse(), cmdr.WithTasksBeforeRun(), cmdr.WithTasksAfterRun

		// true for debug in developing time, it'll disable onAction on each Cmd.
		// for productive mode, comment this line.
		// The envvars FORCE_DEFAULT_ACTION & FORCE_RUN can override this.
		cmdr.WithForceDefaultAction(false),

		// cmdr.WithDontGroupInHelpScreen(false), // default it's false
		cmdr.WithSortInHelpScreen(true), // default it's false
	).
		With(func(app cli.App) { logz.Debug("in dev mode?", "mode", devmode.InDevelopmentMode()) }).
		WithBuilders(
			examples.AddHeadLikeFlagWithoutCmd, // add a `--line` option, feel free to remove it.
			examples.AddToggleGroupFlags,       //
			examples.AddTypedFlags,             //
			examples.AddKilobytesFlag,          //
			examples.AddValidArgsFlag,          //
		).
		WithAdders(cmd.Commands...). // added subcommands here
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
