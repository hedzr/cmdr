package main

// normally tiny app

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/dir"
	"github.com/hedzr/is"
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

		cmdr.WithTasksBeforeRun(func(ctx context.Context, cmd *cli.Command, runner cli.Runner, extras ...any) (err error) {
			logz.Debug("command running...", "cmd", cmd, "runner", runner, "extras", extras)
			return
		}), // cmdr.WithTasksBeforeParse(), cmdr.WithTasksBeforeRun(), cmdr.WithTasksAfterRun

		// true for debug in developing time, it'll disable onAction on each Cmd.
		// for productive mode, comment this line.
		cmdr.WithForceDefaultAction(true),

		cmdr.WithSortInHelpScreen(true),       // default it's false
		cmdr.WithDontGroupInHelpScreen(false), // default it's false
	)

	// // simple run the parser of app and trigger the matched command's action
	// _ = app.Run(
	// 	cmdr.WithForceDefaultAction(false), // true for debug in developing time
	// )

	if err := app.Run(ctx); err != nil {
		logz.Error("Application Error:", "err", err)
		os.Exit(app.SuggestRetCode())
	}
}

func onEvalJumpSubCommands(ctx context.Context, c *cli.Command) (it cli.EvalIterator, err error) {
	return
}

func prepareApp(opts ...cli.Opt) (app cli.App) {
	app = cmdr.New(opts...).
		Info("tiny-app", "0.3.1").
		Author("hedzr")

	app.Flg("no-default").
		Description("disable force default action").
		OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) {
			if b, ok := hitState.Value.(bool); ok {
				app.Store().Set("app.force-default-action", b) // disable/enable the final state about 'force default action'
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
		OnEvaluateSubCommands(onEvalJumpSubCommands).
		With(func(b cli.CommandBuilder) {
			b.Cmd("to").
				Description("to command").
				Examples(``).
				Deprecated(`v0.1.1`).
				// Group(cli.UnsortedGroup).
				Hidden(false).
				OnAction(func(ctx context.Context, cmd *cli.Command, args []string) (err error) {
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
		Group(cli.UnsortedGroup).
		Build()

	app.Flg("wet-run", "w").
		Default(false).
		Description("run all but with committing").
		// Group(cli.UnsortedGroup).
		Build() // no matter even if you're adding the duplicated one.

	app.Cmd("wrong").
		Description("a wrong command to return error for testing").
		OnAction(func(ctx context.Context, cmd *cli.Command, args []string) (err error) {
			ec := errors.New()
			defer ec.Defer(&err) // store the collected errors in native err and return it
			ec.Attach(io.ErrClosedPipe, errors.New("something's wrong"), os.ErrPermission)
			return // handling command action here
		}).
		With(func(b cli.CommandBuilder) {
			b.Flg("full", "f").
				Default(false).
				Description("full command").
				// Group(cli.UnsortedGroup).
				Build()
		})
	return
}

var onceDev sync.Once
var devMode bool

func init() {
	// onceDev is a redundant operation, but we still keep it to
	// fit for defensive programming style.
	onceDev.Do(func() {
		log.SetFlags(log.LstdFlags | log.Lmsgprefix | log.LUTC | log.Lshortfile | log.Lmicroseconds)
		log.SetPrefix("")

		if dir.FileExists(".dev-mode") {
			devMode = true
		} else if dir.FileExists("go.mod") {
			data, err := os.ReadFile("go.mod")
			if err != nil {
				return
			}
			content := string(data)

			// dev := true
			if strings.Contains(content, "github.com/hedzr/cmdr/v2/pkg/") {
				devMode = false
			}

			// I am tiny-app in cmdr/v2, I will be launched in dev-mode always
			if strings.Contains(content, "module github.com/hedzr/cmdr/v2") {
				devMode = true
			}
		}

		if devMode {
			is.SetDebugMode(true)
			logz.Debug(".dev-mode file detected, entering Debug Mode...")
		}

		if is.DebugBuild() {
			is.SetDebugMode(true)
		}

		if is.VerboseBuild() {
			is.SetVerboseMode(true)
		}
	})
}
