package main

import (
	"context"
	"os"

	"github.com/hedzr/cmdr/v2"
	logz "github.com/hedzr/logg/slog"
)

const (
	appName = "struct"
	desc    = `struct buidler version of tiny app.`
	version = cmdr.Version
	author  = `The Example Authors`
)

func main() {
	var Root struct {
		Remove struct {
			Full struct {
				NoForce bool `desc:"DON'T Force removal."`
			} `desc:"remove full of files"`

			Force     bool `help:"Force removal."`
			Recursive bool `help:"Recursively remove files."`

			Paths []string `arg:"" name:"path" help:"Paths to remove." type:"path"`
		} `title:"remove" shorts:"rm" cmd:"" help:"Remove files."`

		List struct {
			Paths []string `name:"path" help:"Paths to list." type:"path"`
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
