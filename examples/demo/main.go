package main

import (
	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/loaders"
	"github.com/hedzr/cmdr/v2/pkg/dir"
)

func main() {
	app := prepareApp()

	// // simple run the parser of app and trigger the matched command's action
	// _ = app.Run(
	// 	cmdr.WithForceDefaultAction(false), // true for debug in developing time
	// )

	if err := app.Run(
		cmdr.WithStore(store.New()),
		cmdr.WithExternalLoaders(
			loaders.NewConfigFileLoader(),
			loaders.NewEnvVarLoader(),
		),
		cmdr.WithForceDefaultAction(true), // true for debug in developing time
	); err != nil {
		logz.Error("Application Error:", "err", err)
	}
}

func prepareApp() (app cli.App) {
	app = cmdr.New().
		Info("demo-app", "0.3.1").
		Author("hedzr")
	app.AddFlg(func(b cli.FlagBuilder) {
		b.Titles("no-default").
			Description("disable force default action").
			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) {
				app.Store().Set("app.force-default-action", false)
				return
			})
	})
	app.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("jump").
			Description("jump command").
			Examples(`jump example`).
			Deprecated(`jump is a demo command`).
			Hidden(false)

		b.AddCmd(func(b cli.CommandBuilder) {
			b.Titles("to").
				Description("to command").
				Examples(``).
				Deprecated(`v0.1.1`).
				Hidden(false).
				OnAction(func(cmd *cli.Command, args []string) (err error) {
					app.Store().Set("app.demo.working", dir.GetCurrentDir())
					println()
					println(dir.GetCurrentDir())
					println()
					println(app.Store().Dump())
					return // handling command action here
				})
			b.AddFlg(func(b cli.FlagBuilder) {
				b.Default(false).
					Titles("full", "f").
					Description("full command").
					Build()
			})
		})
	})

	app.AddFlg(func(b cli.FlagBuilder) {
		b.Titles("dry-run", "n").
			Default(false).
			Description("run all but without committing")
	})

	app.Flg("wet-run", "w").
		Default(false).
		Description("run all but with committing").
		Build() // no matter even if you're adding the duplicated one.
	return
}
