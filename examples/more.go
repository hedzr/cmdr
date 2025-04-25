package examples

import (
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
)

func AttachServerCommand(parent cli.CommandBuilder) {
	common.AttachServerCommand(parent)
}

// AttachKvCommand adds 'kv' command to a builder (eg: app.App)
//
// Example:
//
//	app := cli.New().Info(...)
//	app.AddCmd(func(b cli.CommandBuilder) {
//	  examples.AttachKvCommand(b)
//	})
//	// Or:
//	examples.AttachKvCommand(app.NewCommandBuilder())
func AttachKvCommand(parent cli.CommandBuilder) {
	common.AttachKvCommand(parent)
}

func AttachMsCommand(parent cli.CommandBuilder) {
	common.AttachMsCommand(parent)
}

func AttachMoreCommandsForTest(parent cli.CommandBuilder, moreAndMore bool) {
	common.AttachMoreCommandsForTest(parent, moreAndMore)
}

func AddMxCommand(parent cli.CommandBuilder) {
	common.AddMxCommand(parent)
}

func AddXyPrintCommand(parent cli.CommandBuilder) {
	common.AddXyPrintCommand(parent)
}

func AddPanicTestCommand(parent cli.CommandBuilder) {
	common.AddPanicTestCommand(parent)
}

func AddTtySizeTestCommand(parent cli.CommandBuilder) {
	common.AddTtySizeTestCommand(parent)
}
