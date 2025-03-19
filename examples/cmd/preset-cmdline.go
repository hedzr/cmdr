package cmd

import (
	"context"

	"github.com/hedzr/cmdr/v2/cli"
)

type presetCmd struct{}

func (presetCmd) Add(app cli.App) {
	app.Cmd("preset", "p").
		Description("preset command to inject into user input").
		With(func(b cli.CommandBuilder) {
			b.Flg("preset", "p").
				Default(false).
				Description("preset arg").
				Build()
			b.Cmd("cmd", "c").Description("inject `-pv` into user input cmdline").
				PresetCmdLines(`-pv`).
				OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
					_, err = app.DoBuiltinAction(ctx, cli.ActionDefault)
					return
				}).Build()
		})
}
