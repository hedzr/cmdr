package cmd

import (
	"context"
	"io"
	"os"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2/cli"
)

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
