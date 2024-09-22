package main

import (
	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/dir"
)

func main() {
	app := prepareApp()

	// // simple run the parser of app and trigger the matched command's action
	// _ = app.Run(
	// 	cmdr.WithForceDefaultAction(false), // true for debug in developing time
	// )

	if err := app.Run(
		cmdr.WithStore(store.New()), // use a option store, if not specified by store.New(), a dummy store allocated

		// cmdr.WithExternalLoaders(
		// 	local.NewConfigFileLoader(),
		// 	local.NewEnvVarLoader(),
		// ),
		cmdr.WithForceDefaultAction(true), // true for debug in developing time
	); err != nil {
		logz.Error("Application Error:", "err", err)
	}
}

func prepareApp() (app cli.App) {
	app = cmdr.New().
		Info("tiny-app", "0.3.1").
		Author("hedzr")

	app.Flg("no-default").
		Description("disable force default action").
		OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) {
			app.Store().Set("app.force-default-action", false)
			return
		}).
		// Group(cli.UnsortedGroup).
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
				OnAction(func(cmd *cli.Command, args []string) (err error) {
					app.Store().Set("app.demo.working", dir.GetCurrentDir())
					println()
					println(dir.GetCurrentDir())
					println()
					println(app.Store().Dump())
					return // handling command action here
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
		Group(cli.UnsortedGroup).
		Build()

	app.Flg("wet-run", "w").
		Default(false).
		Description("run all but with committing").
		// Group(cli.UnsortedGroup).
		Build() // no matter even if you're adding the duplicated one.

	return
}
