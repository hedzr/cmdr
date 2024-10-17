// Copyright © 2022 Hedzr Yeh.

package builder

import (
	"os"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/conf"
	"github.com/hedzr/cmdr/v2/pkg/dir"
)

func New(w cli.Runner) cli.App {
	b := &appS{
		Runner: w,
		root:   newDefaultRoot(),
		args:   os.Args,
	}
	if x, ok := w.(interface{ Args() []string }); ok {
		if args := x.Args(); args != nil {
			b.args = args
		}
	}
	return b
}

func newDefaultRoot() *cli.RootCommand {
	if conf.AppName == "" {
		conf.AppName = dir.Basename(dir.GetExecutablePath())
	}
	if conf.Version == "" {
		conf.Version = "v0.0.1-debug"
	}

	root := new(cli.RootCommand)
	root.Cmd = new(cli.CmdS)
	root.AppName = conf.AppName
	root.Version = conf.Version
	return root
}

type setRoot interface {
	SetRoot(root *cli.RootCommand, args []string)
}

type adder interface {
	addCommand(child *cli.CmdS)
	addFlag(child *cli.Flag)
}
