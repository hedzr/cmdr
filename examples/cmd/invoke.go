package cmd

import (
	"github.com/hedzr/cmdr/v2/cli"
)

type invokeCmd struct{}

func (invokeCmd) Add(app cli.App) {
	app.Cmd("invoke").Description(`test invoke feature`).
		With(func(b cli.CommandBuilder) {
			b.Cmd("shell").Description(`invoke shell cmd`).InvokeShell(`ls -la`).UseShell("/bin/bash").OnAction(nil).Build()
			b.Cmd("proc").Description(`invoke gui program`).InvokeProc(`say "hello, world!"`).OnAction(nil).Build()
		})
}
