package main

import (
	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples"
	"github.com/hedzr/cmdr/v2/loaders"
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
		cmdr.WithForceDefaultAction(false), // true for debug in developing time
	); err != nil {
		logz.Error("Application Error:", "err", err)
	}
}

func prepareApp() (app cli.App) {
	app = cmdr.New().
		Info("large-app", "0.3.1").
		Author("hedzr")

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

	app.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("consul", "c").
			Description("command set for consul operations")

		b.Flg("data-center", "dc", "datacenter").
			// Description("set data-center").
			Default("dc-1").
			Build()
	})

	app.AddCmd(func(b cli.CommandBuilder) {
		examples.AttachServerCommand(b)
	})

	app.AddCmd(func(b cli.CommandBuilder) {
		examples.AttachKvCommand(b)
	})

	app.AddCmd(func(b cli.CommandBuilder) {
		examples.AttachMsCommand(b)
	})

	app.AddCmd(func(b cli.CommandBuilder) {
		examples.AttachMoreCommandsForTest(b, false)
	})

	app.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("display", "da").
			Description("command set for display adapter operations")

		b1 := b.Cmd("voodoo", "vd").
			Description("command set for voodoo operations")
		b1.Flg("data-center", "dc", "datacenter").
			Default("dc-1").
			Build()
		b1.Build()

		b2 := b.Cmd("nvidia", "nv").
			Description("command set for nvidia operations")
		b2.Flg("data-center", "dc", "datacenter").
			Default("dc-1").
			Build()
		b2.Build()

		b.AddCmd(func(b cli.CommandBuilder) {
			b.Titles("amd", "amd").
				Description("command set for AMD operations")
			b.Flg("data-center", "dc", "datacenter").
				Default("dc-1").
				Build()
		})
	})

	return
}
