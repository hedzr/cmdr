package loaders

import (
	"os"
	"testing"

	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2/builder"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/worker"
	"github.com/hedzr/cmdr/v2/examples"
)

type discardP struct{}

func (*discardP) Write([]byte) (n int, err error) { return }

func (*discardP) WriteString(string) (n int, err error) { return }

func TestWorkerS_Pre(t *testing.T) {
	app, ww := cleanApp(t,
		cli.WithStore(store.New()),
		cli.WithExternalLoaders(NewConfigFileLoader(), NewEnvVarLoader()),
		cli.WithArgs("--debug"),
		cli.WithHelpScreenWriter(&discardP{}),
	)

	// app := buildDemoApp()
	// ww := postBuild(app)
	// ww.InitGlobally()
	// assert.EqualTrue(t, ww.Ready())
	//
	// ww.ForceDefaultAction = true
	// ww.wrHelpScreen = &discardP{}
	// ww.wrDebugScreen = os.Stdout
	// ww.wrHelpScreen = os.Stdout

	// ww.setArgs([]string{"--debug"})
	// ww.Config.Store = store.New()
	// ww.Loaders = []cli.Loader{NewConfigFileLoader(), NewEnvVarLoader()}
	_ = app

	err := ww.Run(
		cli.WithTasksBeforeParse(func(root *cli.RootCommand, runner cli.Runner) (err error) {
			root.SelfAssert()
			t.Logf("root.SelfAssert() passed. runner = %v", runner)
			return
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func cleanApp(t *testing.T, opts ...cli.Opt) (app cli.App, ww cli.Runner) { //nolint:revive
	opts = append(opts,
		// cli.WithHelpScreenWriter(os.Stdout),
		cli.WithDebugScreenWriter(os.Stdout),
		cli.WithForceDefaultAction(true),
		cli.WithTasksBeforeParse(func(root *cli.RootCommand, runner cli.Runner) (err error) {
			_, _ = root, runner
			return
		}),
	)
	app = buildDemoApp(opts...)
	ww = postBuild(app)
	if r, ok := ww.(interface{ InitGlobally() }); ok {
		r.InitGlobally()
	}
	if !ww.Ready() {
		t.Fatalf("not ready")
	}

	// assertTrue(t, ww.Ready())
	// ww.wrHelpScreen = &discardP{}
	// if helpScreen {
	// 	ww.wrHelpScreen = os.Stdout
	// }
	// ww.wrDebugScreen = os.Stdout
	// ww.ForceDefaultAction = true
	// ww.tasksAfterParse = []taskAfterParse{func(w *workerS, ctx *parseCtx, errParsed error) (err error) { return }} //nolint:revive

	// ww.setArgs([]string{"--debug"})
	// err := ww.Run(withTasksBeforeParse(func(root *cli.RootCommand, runner cli.Runner) (err error) {
	// 	root.SelfAssert()
	// 	t.Logf("root.SelfAssert() passed.")
	// 	return
	// }))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// ww.TasksBeforeParse = nil
	return
}

func buildDemoApp(opts ...cli.Opt) (app cli.App) { //nolint:revive
	cfg := cli.NewConfig(opts...)

	// cfg := cli.New(cli.WithStore(store.New()))

	w := worker.New(cfg)

	app = builder.New(w).
		Info("demo-app", "0.3.1").
		Author("hedzr")

	app.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("jump").
			Description("jump is a demo command").
			Examples(`jump example`).
			Deprecated(`since: v0.9.1`).
			Hidden(false).
			AddCmd(func(b cli.CommandBuilder) {
				b.Titles("to").
					Description("to command").
					Examples(``).
					Deprecated(``).
					Hidden(false).
					OnAction(func(cmd *cli.Command, args []string) (err error) { //nolint:revive
						return // handling command action here
					})

				b.AddFlg(func(b cli.FlagBuilder) {
					b.Default(false).
						Titles("full", "f").
						Description("full command")
				})
			})
		b.AddFlg(func(b cli.FlagBuilder) {
			b.Default(false).
				Titles("empty", "e").
				Description("empty command")
		})
	}).AddFlg(func(b cli.FlagBuilder) {
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
		examples.AttachMoreCommandsForTest(b, true)
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

func postBuild(app cli.App, args ...string) (ww cli.Runner) {
	if sr, ok := app.(interface{ Worker() cli.Runner }); ok {
		ww = sr.Worker()
		if ww1, ok := sr.Worker().(interface {
			SetRoot(root *cli.RootCommand, args []string)
		}); ok {
			if r, ok := app.(interface{ Root() *cli.RootCommand }); ok {
				r.Root().EnsureTree(app, r.Root())
				ww1.SetRoot(r.Root(), app.Args())
			}
		}
	}
	return
}
