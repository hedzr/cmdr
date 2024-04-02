// Copyright © 2022 Hedzr Yeh.

package builder

import (
	"os"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/conf"
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
	root := new(cli.RootCommand)
	root.Command = new(cli.Command)
	root.AppName = conf.AppName
	root.Version = conf.Version
	return root
}

type setRoot interface {
	SetRoot(root *cli.RootCommand, args []string)
}

type adder interface {
	addCommand(child *cli.Command)
	addFlag(child *cli.Flag)
}
