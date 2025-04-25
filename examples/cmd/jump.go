package cmd

import (
	"context"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/dyncmd"
	"github.com/hedzr/is/dir"
	logz "github.com/hedzr/logg/slog"
)

type jumpCmd struct{}

func (jumpCmd) Add(app cli.App) {
	app.Cmd("jump").
		Description("jump command").
		Examples(`jump example`). // {{.AppName}}, {{.AppVersion}}, {{.DadCommands}}, {{.Commands}}, ...
		Deprecated(`v1.1.0`).
		Group("Test").
		// Group(cli.UnsortedGroup).
		// Hidden(false).
		OnEvaluateSubCommands(dyncmd.OnEvalJumpSubCommands).
		OnEvaluateSubCommandsFromConfig().
		// Both With(cb) and Build() to end a building sequence
		With(func(b cli.CommandBuilder) {
			b.Cmd("to").
				Description("to command").
				Examples(``).
				Deprecated(`v0.1.1`).
				OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
					// cmd.Set() == cmdr.Set(), cmd.Store() == cmdr.Store()
					_, _ = cmd.Set().Set("tiny3.working", dir.GetCurrentDir())
					println()
					println("dir:", cmd.Set().WithPrefix("tiny3").MustString("working"))

					cs := cmd.Set().WithPrefix("jump.to")
					if cs.MustBool("full") {
						println()
						println(cmd.Set().Dump())
					}
					cs2 := cmd.Store()
					if cs2.MustBool("full") != cs.MustBool("full") {
						logz.Panic("a bug found")
					}
					app.SetSuggestRetCode(1) // ret code must be in 0-255
					return                   // handling command action here
				}).
				With(func(b cli.CommandBuilder) {
					b.Flg("full", "f").
						Default(false).
						Description("full command").
						Build()
				})
		})
}
