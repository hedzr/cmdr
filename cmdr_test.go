package cmdr_test

import (
	"context"
	"io"
	"os"
	"testing"

	"gopkg.in/hedzr/errors.v3"

	cmdr "github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/cmd"
	"github.com/hedzr/cmdr/v2/examples/devmode"
	"github.com/hedzr/cmdr/v2/examples/dyncmd"
	"github.com/hedzr/cmdr/v2/pkg/logz"
	"github.com/hedzr/is/dir"
	"github.com/hedzr/store"
)

func ExampleCreate() {
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

func ExampleCreate_buildFromAStructValue() {
	var Root struct {
		Remove struct {
			Force     bool `help:"Force removal."`
			Recursive bool `help:"Recursively remove files."`

			Paths []string `arg:"" name:"path" help:"Paths to remove." type:"path"`
		} `title:"remove" shorts:"rm" cmd:"" help:"Remove files."`

		List struct {
			Paths []string `arg:"" optional:"" name:"path" help:"Paths to list." type:"path"`
		} `title:"list" shorts:"ls" cmd:"" help:"List paths."`
	}

	app := cmdr.Create(appName, version, author, desc).
		// WithAdders(cmd.Commands...).
		BuildFrom(&Root)

	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}

func ExampleCreator_BuildFrom() {
	var Root struct {
		Remove struct {
			Force     bool `help:"Force removal."`
			Recursive bool `help:"Recursively remove files."`

			Paths []string `arg:"" name:"path" help:"Paths to remove." type:"path"`
		} `title:"remove" shorts:"rm" cmd:"" help:"Remove files."`

		List struct {
			Paths []string `arg:"" optional:"" name:"path" help:"Paths to list." type:"path"`
		} `title:"list" shorts:"ls" cmd:"" help:"List paths."`
	}

	app := cmdr.Create(appName, version, author, desc).
		// WithAdders(cmd.Commands...).
		BuildFrom(&Root)

	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}

const (
	appName = "concise"
	desc    = `concise version of tiny app.`
	version = cmdr.Version
	author  = `The Example Authors`
)

type TinyApp struct{}

func (TinyApp) Traditional()            {}
func (TinyApp) ViaCreate()              {}
func (TinyApp) ViaCreateViaFromStruct() {}

func ExampleTinyApp_ViaCreateViaFromStruct() {
	var Root struct {
		Remove struct {
			Force     bool `help:"Force removal."`
			Recursive bool `help:"Recursively remove files."`

			Paths []string `arg:"" name:"path" help:"Paths to remove." type:"path"`
		} `title:"remove" shorts:"rm" cmd:"" help:"Remove files."`

		List struct {
			Paths []string `arg:"" optional:"" name:"path" help:"Paths to list." type:"path"`
		} `title:"list" shorts:"ls" cmd:"" help:"List paths."`
	}

	app := cmdr.Create(appName, version, author, desc).
		// WithAdders(cmd.Commands...).
		BuildFrom(&Root)

	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}

func ExampleTinyApp_ViaCreate() { ExampleCreate() }

func ExampleTinyApp_Traditional() {
	ctx := context.Background()

	app := prepareApp(
		cmdr.WithStore(store.New()), // use an option store explicitly, or a dummy store by default

		// cmdr.WithExternalLoaders(
		// 	local.NewConfigFileLoader(), // import "github.com/hedzr/cmdr-loaders/local" to get in advanced external loading features
		// 	local.NewEnvVarLoader(),
		// ),

		cmdr.WithTasksBeforeRun(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) {
			logz.DebugContext(ctx, "command running...", "cmd", cmd, "runner", runner, "extras", extras)
			return
		}), // cmdr.WithTasksBeforeParse(), cmdr.WithTasksBeforeRun(), cmdr.WithTasksAfterRun

		// true for debug in developing time, it'll disable onAction on each Cmd.
		// for productive mode, comment this line.
		// The envvars FORCE_DEFAULT_ACTION & FORCE_RUN can override this.
		cmdr.WithForceDefaultAction(true),

		cmdr.WithSortInHelpScreen(true),       // default it's false
		cmdr.WithDontGroupInHelpScreen(false), // default it's false

		cmdr.WithAutoEnvBindings(true),
	)

	logz.Debug("in dev mode?", "mode", devmode.InDevelopmentMode())

	// // simple run the parser of app and trigger the matched command's action
	// _ = app.Run(
	// 	cmdr.WithForceDefaultAction(false), // true for debug in developing time
	// )

	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	}
}

func prepareApp(opts ...cli.Opt) (app cli.App) {
	app = cmdr.New(opts...).
		Info("tiny-app", "0.3.1").
		Author("The Example Authors") // .Description(``).Header(``).Footer(``)

	// another way to disable `cmdr.WithForceDefaultAction(true)` is using
	// env-var FORCE_RUN=1 (builtin already).
	app.Flg("no-default").
		Description("disable force default action").
		OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) {
			if b, ok := hitState.Value.(bool); ok {
				f.Set().Set("app.force-default-action", b) // disable/enable the final state about 'force default action'
			}
			return
		}).
		Build()

	app.Cmd("jump").
		Description("jump command").
		Examples(`jump example`). // {{.AppName}}, {{.AppVersion}}, {{.DadCommands}}, {{.Commands}}, ...
		Deprecated(`v1.1.0`).
		// Group(cli.UnsortedGroup).
		Hidden(false, false).
		OnEvaluateSubCommands(dyncmd.OnEvalJumpSubCommands).
		With(func(b cli.CommandBuilder) {
			b.Cmd("to").
				Description("to command").
				Examples(``).
				Deprecated(`v0.1.1`).
				// Group(cli.UnsortedGroup).
				Hidden(false).
				OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
					// cmd.Set() == cmdr.Set(), cmd.Store() == cmdr.Store()
					cmd.Set().Set("app.demo.working", dir.GetCurrentDir())
					println()
					println(cmd.Set().WithPrefix("app.demo").MustString("working"))

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
		Build() // no matter even if you're adding the duplicated one.

	app.Cmd("wrong").
		Description("a wrong command to return error for testing").
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
		With(func(b cli.CommandBuilder) {
			b.Flg("duration", "d").
				Default("5s").
				Description("a duration var").
				Build()
		})
	return
}

func TestDottedPathToCommandOrFlag(t *testing.T) {
	ctx := context.Background()
	app := cmdr.Create("app", "v1", `author`, `desc`,
		cli.WithArgs("test-app", "--debug"),
	).OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
		cc, ff := cmdr.DottedPathToCommandOrFlag("generate.shell.zsh", nil)
		if cc == nil || cc.GetTitleName() != "shell" || ff.Title() != "zsh" {
			t.Fail()
		}
		cc, ff = cmdr.DottedPathToCommandOrFlag("generate.doc", nil)
		if ff != nil || cc.GetTitleName() != "doc" {
			t.Fail()
		}
		return
	}).
		Build()
	err := app.Run(ctx)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestStoreGetSectionFrom(t *testing.T) {
	ctx := context.Background()
	app := cmdr.Create("test-app", "v1", `author`, `desc`,
		cli.WithArgs("test-app", "--debug"),
	).OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
		cs := cmdr.Store()
		b := cs.MustBool("debug")
		println("Dump: \n", cs.Dump())
		println("debug flag: ", b)
		if !b {
			t.Fail()
		}
		println(cmdr.Set().Dump())

		type manS struct {
			Dir  string
			Type int
		}
		type genS struct {
			Manual manS
		}
		var v genS
		set := cmdr.Set(cli.CommandsStoreKey)
		err = set.GetSectionFrom("generate", &v)
		if err != nil {
			t.Fail()
		}
		if v.Manual.Type != 1 {
			t.Fail()
		}
		return
	}).
		Build()
	err := app.Run(ctx)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestJumpTo(t *testing.T) {
	ctx := context.Background()
	app := cmdr.Create("test-app", "v1", `author`, `desc`,
		cli.WithArgs("test-app", "jump", "to", "--debug"),
	).
		WithAdders(cmd.Commands[0]).
		Build()
	err := app.Run(ctx)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestTo(t *testing.T) {
	ctx := context.Background()
	app := cmdr.Create("app", "v1", `author`, `desc`,
		cli.WithArgs("test-app", "--debug"),
	).OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
		b := cmdr.Store().MustBool("debug")
		println("debug flag: ", b)
		if !b {
			t.Fail()
		}
		println(cmdr.Set().Dump())

		type manS struct {
			Dir  string
			Type int
		}
		type genS struct {
			Manual manS
		}
		var v genS
		err = cmdr.To("cmd.generate", &v)
		if err != nil {
			t.Fail()
		}
		if v.Manual.Type != 1 {
			t.Fail()
		}
		return
	}).
		Build()
	err := app.Run(ctx)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestExecNoRoot(t *testing.T) {
	if err := cmdr.Exec(nil); !errors.Iss(err, cli.ErrEmptyRootCommand) {
		t.Errorf("Error: %v", err)
	}
}

// func TestExecSimple(t *testing.T) {
// 	if err := cmdr.Exec(testdata.BuildCommands(true)); !errors.Iss(err, cli.ErrEmptyRootCommand) {
// 		t.Errorf("Error: %v", err)
// 	}
// }

func TestGetSet(t *testing.T) {
	ctx := context.Background()
	app := cmdr.Create("app", "v1", `author`, `desc`,
		cli.WithArgs("test-app", "--debug"),
	).
		With(func(app cli.App) {
			app.OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
				b := cmdr.Store().MustBool("debug")
				println("debug flag: ", b)
				if !b {
					t.Fail()
				}
				// println(cmdr.Set().Dump())
				return
			})
		}).
		Build()
	err := app.Run(ctx)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
