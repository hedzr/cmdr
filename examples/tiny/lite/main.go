package main

import (
	"context"
	"io"
	"os"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/cmd"
	"github.com/hedzr/cmdr/v2/examples/dyncmd"

	"github.com/hedzr/is/dir"
	logz "github.com/hedzr/logg/slog"
	"gopkg.in/hedzr/errors.v3"
)

const (
	appName = "lite-app"
	desc    = `lite-app version of tiny app.`
	version = cmdr.Version
	author  = `The Example Authors`
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := cmdr.Create(appName, version, author, desc).
		WithAdders(cmd.Commands...).
		Build()

	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}

var Commands = []cli.CmdAdder{
	jumpCmd{},
	wrongCmd{},
	invokeCmd{},
	presetCmd{},
}

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
					cmd.Set().Set("tiny3.working", dir.GetCurrentDir())
					println()
					println("dir:", cmd.Set().WithPrefix("tiny3").MustString("working"))

					cs := cmdr.Store().WithPrefix("jump.to")
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

type wrongCmd struct{}

func (wrongCmd) Add(app cli.App) {
	app.Cmd("wrong").
		Description("a wrong command to return error for testing").
		Group("Test").
		// cmdline `FORCE_RUN=1 go run ./tiny wrong -d 8s` to verify this command to see the returned application error.
		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			dur := cmd.Store().MustDuration("duration")
			println("the duration is:", dur.String())

			ec := errors.New()
			defer ec.Defer(&err) // store the collected errors in native err and return it
			ec.Attach(io.ErrClosedPipe, errors.New("something's wrong"), os.ErrPermission)
			// see the application error by running `go run ./tiny/tiny/main.go wrong`.
			return
		}).
		Build()
}

type invokeCmd struct{}

func (invokeCmd) Add(app cli.App) {
	app.Cmd("invoke").Description(`test invoke feature`).
		With(func(b cli.CommandBuilder) {
			b.Cmd("shell").Description(`invoke shell cmd`).InvokeShell(`ls -la`).UseShell("/bin/bash").OnAction(nil).Build()
			b.Cmd("proc").Description(`invoke gui program`).InvokeProc(`say "hello, world!"`).OnAction(nil).Build()
		})
}

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
