package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
	"github.com/hedzr/cmdr/v2/examples/devmode"
	logz "github.com/hedzr/logg/slog"
)

const (
	appName = "negatable"
	desc    = `a sample to demo negatable flag.`
	version = cmdr.Version
	author  = ``
)

func main() {
	app := cmdr.Create(appName, version, author, desc).
		With(func(app cli.App) { logz.Debug("in dev mode?", "mode", devmode.InDevelopmentMode()) }).
		WithBuilders(common.AddNegatableFlag).
		WithAdders().
		// override the onAction defined in common.AddNegatableFlag()
		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			// logz.PrintlnContext(ctx, cmd.Set().Dump())

			// v := cmd.Store().MustBool("warning")
			// w := cmd.Store().MustBool("no-warning")
			// println("warning flag: ", v)
			// println("no-warning flag: ", w)

			// wf := cmd.FlagBy("warning")
			// if wf != nil && wf.LeadingPlusSign() {
			// 	println("warning flag with leading plus sign: +w FOUND.")
			// }
			// return

			println("Dumping Store ----------\n", cmdr.Set().Dump())

			cmd.App().DoBuiltinAction(ctx, cli.ActionDefault)

			cs := cmd.Store()

			v := cs.MustBool("warning")
			w := cs.MustBool("no-warning")
			println("warning toggle group: ", cs.MustString("warning"))
			println("  warning flag: ", v)
			println("  no-warning flag: ", w)

			wf := cmd.FlagBy("warning")
			if wf != nil && wf.LeadingPlusSign() {
				println()
				println("NOTABLE: a flag with leading plus sign: +w FOUND.")
			} else {
				logz.WarnContext(ctx, "cannot found flag 'warning'")
			}

			sw1 := cs.MustBool("warnings.unused-variable")
			sw2 := cs.MustBool("warnings.unused-parameter")
			sw3 := cs.MustBool("warnings.unused-function")
			sw4 := cs.MustBool("warnings.unused-but-set-variable")
			sw5 := cs.MustBool("warnings.unused-private-field")
			sw6 := cs.MustBool("warnings.unused-label")
			fmt.Printf(`

--warnings, -W:
    > TG: %q
    > TG.selected: %q
	unused-variable:		%v, (no-): %v
	unused-parameter:		%v, (no-): %v
	unused-function:		%v, (no-): %v
	unused-but-set-variable:	%v, (no-): %v
	unused-private-field:		%v, (no-): %v
	unused-label:			%v, (no-): %v
`,
				cs.MustString("warnings"),
				cs.MustStringSlice("warnings.selected"),
				sw1, cs.MustBool("warnings.no-unused-variable"),
				sw2, cs.MustBool("warnings.no-unused-parameter"),
				sw3, cs.MustBool("warnings.no-unused-function"),
				sw4, cs.MustBool("warnings.no-unused-but-set-variable"),
				sw5, cs.MustBool("warnings.no-unused-private-field"),
				sw6, cs.MustBool("warnings.no-unused-label"),
			)

			return
		}).
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
