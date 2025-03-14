package main

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples"
	"github.com/hedzr/is/term/color"
	logz "github.com/hedzr/logg/slog"
)

const (
	appName = "required"
	desc    = `a sample to show u what error raised by a missed required flag.`
	version = cmdr.Version
	author  = `The Example Authors`
)

func main() {
	app := cmdr.Create(appName, version, author, desc).
		With(func(app cli.App) {
			app.OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
				msg := cmd.Store().MustString("required")
				println(`Hello, World.`, msg)
				return
			})
		}).
		WithBuilders(examples.AddRequiredFlag).
		WithAdders().
		Build()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Run(ctx); err != nil {
		if errors.Is(err, cli.ErrRequiredFlag) {
			fmt.Printf(color.StripLeftTabsC(`
				The <b>REQUIRED</b> Flag not present:
				
				<font color="red">%s</red>
				`),
				err)
			os.Exit(1)
		}
		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
		os.Exit(app.SuggestRetCode())
	} else if rc := app.SuggestRetCode(); rc != 0 {
		os.Exit(rc)
	}
}
