package cmd

import (
	"github.com/hedzr/cmdr/v2/cli"
)

type invokeCmd struct{}

func (invokeCmd) Add(app cli.App) {
	app.Cmd("invoke").
		With(func(b cli.CommandBuilder) {
			b.Cmd("shell").InvokeShell(`ls -la`).UseShell("/bin/bash").OnAction(nil).Build()
			b.Cmd("proc").InvokeProc(`say "hello, world!"`).OnAction(nil).Build()
		})
}
