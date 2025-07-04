package main

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/is/term/color"
	logz "github.com/hedzr/logg/slog"
)

const (
	appName = "valid-args"
	desc    = `a sample to show u how to using enum values in option.`
	version = cmdr.Version
	author  = ``
)

func main() {
	app := cmdr.Create(appName, version, author, desc).
		With(func(app cli.App) {}).
		WithBuilders(validArgsCommand).
		WithAdders().
		Build()

	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		if errors.Is(err, cli.ErrValidArgs) {
			fmt.Printf(color.StripLeftTabsC(`
				The input is not in valid list:
				
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

func validArgsCommand(parent cli.CommandBuilder) { //nolint:revive
	parent.Flg("enum", "e").
		Description("valid args option", "").
		Group("Test").
		ValidArgs("apple", "banana", "orange").
		Default("").
		Build()

	// give root command an action to handle it
	parent.OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
		v := cmd.Store().MustString("enum")
		println("enum value is: ", v)
		return
	})

	parent.Examples(`Try to use value-args app,
	  
	  $ $APP -e apple
	    this command works ok.
	  $ $APP -e mongo
	    can't work because valid args are: apple, banana, and orange.
	  $ $APP --help
	    check out the valid values in help screen.
	`)
}
