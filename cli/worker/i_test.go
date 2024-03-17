package worker

import (
	"fmt"
	"os"
	"testing"

	"github.com/hedzr/cmdr/v2/builder"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples"
)

func cleanApp(t *testing.T, helpScreen bool) (app cli.App, ww *workerS) { //nolint:revive
	app = buildDemoApp()
	ww = postBuild(app)
	ww.InitGlobally()
	assertTrue(t, ww.Ready())

	ww.wrHelpScreen = &discardP{}
	if helpScreen {
		ww.wrHelpScreen = os.Stdout
	}
	ww.wrDebugScreen = os.Stdout
	ww.ForceDefaultAction = true
	ww.tasksAfterParse = []taskAfterParse{func(w *workerS, ctx *parseCtx, errParsed error) (err error) { return }}

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

func buildDemoApp() (app cli.App) { //nolint:revive
	// cfg := cli.New(cli.WithStore(store.New()))

	cfg := cli.NewConfig()

	w := New(cfg)

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
					OnAction(func(cmd *cli.Command, args []string) (err error) {
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

func postBuild(app cli.App) (ww *workerS) {
	if sr, ok := app.(interface{ Worker() cli.Runner }); ok {
		if ww, ok = sr.Worker().(*workerS); ok {
			if r, ok := app.(interface{ Root() *cli.RootCommand }); ok {
				r.Root().EnsureTree(app, r.Root())
				ww.SetRoot(r.Root(), ww.args)
			}
		}
	}
	return
}

//

//

//

func assertTrue(t testing.TB, cond bool, msg ...any) {
	if cond {
		return
	}

	var mesg string
	if len(msg) > 0 {
		if format, ok := msg[0].(string); ok {
			mesg = fmt.Sprintf(format, msg[1:]...)
		} else {
			mesg = fmt.Sprint(msg...)
		}
	}

	t.Fatalf("assertTrue failed: %s", mesg)
}

func assertFalse(t testing.TB, cond bool, msg ...any) {
	if !cond {
		return
	}

	var mesg string
	if len(msg) > 0 {
		if format, ok := msg[0].(string); ok {
			mesg = fmt.Sprintf(format, msg[1:]...)
		} else {
			mesg = fmt.Sprint(msg...)
		}
	}

	t.Fatalf("assertFalse failed: %s", mesg)
}
