package main

// Simplest tiny app

import (
	"context"
	"io"
	"os"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/dir"
	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"
)

func main() {
	ctx := context.Background()
	app := prepareApp(
		cmdr.WithStore(store.New()), // use a option store, if not specified by store.New(), a dummy store allocated

		// cmdr.WithExternalLoaders(
		// 	local.NewConfigFileLoader(),      // import "github.com/hedzr/cmdr-loaders/local" to get in advanced external loading features
		// 	local.NewEnvVarLoader(),
		// ),

		cmdr.WithTasksBeforeRun(func(ctx context.Context, cmd cli.BaseOptI, runner cli.Runner, extras ...any) (err error) {
			logz.DebugContext(ctx, "command running...", "cmd", cmd, "runner", runner, "extras", extras)
			return
		}),

		// true for debug in developing time, it'll disable onAction on each Cmd.
		// for productive mode, comment this line.
		cmdr.WithForceDefaultAction(true),
	)
	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err)
		os.Exit(app.SuggestRetCode())
	}
}

func prepareApp(opts ...cli.Opt) (app cli.App) {
	app = cmdr.New(opts...).
		Info("tiny-app", "0.3.1").
		Author("hedzr")

	app.Flg("no-default").
		Description("disable force default action").
		OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) {
			if b, ok := hitState.Value.(bool); ok {
				// disable/enable the final state about 'force default action'
				app.Store().Set("app.force-default-action", b)
			}
			return
		}).
		Build()

	app.Cmd("jump").
		Description("jump command").
		Examples(`jump example`).
		Deprecated(`v1.1.0`).
		// Group(cli.UnsortedGroup).
		Hidden(false).
		With(func(b cli.CommandBuilder) {
			b.Cmd("to").
				Description("to command").
				Examples(``).
				Deprecated(`v0.1.1`).
				// Group(cli.UnsortedGroup).
				Hidden(false).
				OnAction(func(ctx context.Context, cmd cli.BaseOptI, args []string) (err error) {
					app.Store().Set("app.demo.working", dir.GetCurrentDir())
					println()
					println(dir.GetCurrentDir())
					println()
					println(app.Store().Dump())
					app.SetSuggestRetCode(1) // ret code must be in 0-255
					return                   // handling command action here
				}).
				With(func(b cli.CommandBuilder) {
					b.Flg("full", "f").
						Default(false).
						Description("full command").
						// Group(cli.UnsortedGroup).
						Build()
				})
		})

	app.Flg("dry-run", "n").
		Default(false).
		Description("run all but without committing").
		Build()

	app.Flg("wet-run", "w").
		Default(false).
		Description("run all but with committing").
		Build() // no matter even if you're adding the duplicated one.

	app.Cmd("wrong").
		Description("a wrong command to return error for testing").
		OnAction(func(ctx context.Context, cmd cli.BaseOptI, args []string) (err error) {
			ec := errors.New()
			defer ec.Defer(&err) // store the collected errors in native err and return it
			ec.Attach(io.ErrClosedPipe, errors.New("something's wrong"), os.ErrPermission)
			return
		}).
		With(func(b cli.CommandBuilder) {
			b.Flg("full", "f").
				Default(false).
				Description("full command").
				Build()
		})
	return
}
