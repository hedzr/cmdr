package worker

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hedzr/cmdr/v2/builder"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
)

func cleanApp(t *testing.T, ctx context.Context, helpScreen bool, opts ...cli.Opt) (app cli.App, ww *workerS) { //nolint:revive
	app = buildDemoApp(opts...)
	ww = postBuild(ctx, app)
	ww.InitGlobally(ctx)
	assertTrue(t, ww.Ready())

	if ww.wrHelpScreen == nil && ww.HelpScreenWriter == nil {
		ww.wrHelpScreen = &discardP{}
		if helpScreen {
			ww.wrHelpScreen = os.Stdout
		}
	}
	if ww.wrDebugScreen == nil && ww.DebugScreenWriter == nil {
		ww.wrDebugScreen = os.Stdout
	}
	ww.ForceDefaultAction = true
	ww.tasksAfterParse = []taskAfterParse{func(ctx context.Context, w *workerS, pc *parseCtx, errParsed error) (err error) { return }} //nolint:revive

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

	w := New(cfg)

	app = builder.New(w).
		Info("demo-app", "0.3.1").
		Author("hedzr")

	b := app.Cmd("jump").
		Description("jump is a demo command").
		Examples(`jump example`).
		Deprecated(`since: v0.9.1`).
		Hidden(false)
	b.Cmd("to").
		Description("to command").
		Examples(``).
		Deprecated(``).
		Hidden(false).
		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) { //nolint:revive
			return // handling command action here
		}).
		Build()

	b.Flg("full", "f").
		Default(false).
		Description("full command").
		Build()

	b.Flg("empty", "e").
		Default(false).
		Description("empty command").
		Build()
	b.Build()

	app.Flg("dry-run", "n").
		Default(false).
		Description("run all but without committing").
		Build()

	app.Flg("wet-run", "w").
		Default(false).
		Description("run all but with committing").
		Build() // no matter even if you're adding the duplicated one.

	b = app.Cmd("consul", "c").
		Description("command set for consul operations")
	b.Flg("data-center", "dc", "datacenter").
		// Description("set data-center").
		Default("dc-1").
		Build()
	b.Build()

	common.AttachServerCommand(app.Cmd("server"))

	common.AttachKvCommand(app.Cmd("kv"))

	common.AttachMsCommand(app.Cmd("ms"))

	common.AttachMoreCommandsForTest(app.Cmd("more"), true)

	b = app.Cmd("display", "da").
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

	b3 := b.Cmd("amd", "amd").
		Description("command set for AMD operations")
	b3.Flg("data-center", "dc", "datacenter").
		Default("dc-1").
		Build()
	b3.Build()

	b.Build()
	return
}

func postBuild(ctx context.Context, app cli.App) (ww *workerS) {
	if sr, ok := app.(interface{ Worker() cli.Runner }); ok {
		if ww, ok = sr.Worker().(*workerS); ok {
			if r, ok := app.(interface{ Root() *cli.RootCommand }); ok {
				// call EnsureTree without set internal flag so that we can
				// run EnsureTree again at next time (but once after Run())
				if cx, ok := r.Root().Cmd.(*cli.CmdS); ok {
					cx.EnsureTreeAlways(ctx, app, r.Root())
				}
				ww.SetRoot(r.Root(), ww.args)
			}
		}
	}
	return
}

//

//

//

func assertTrue(t testing.TB, cond bool, msg ...any) { //nolint:revive
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

func assertFalse(t testing.TB, cond bool, msg ...any) { //nolint:unused,revive
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
