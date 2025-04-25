package examples

import (
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
)

// AddToggleGroupCommand _
//
// Sample code,
//
//	app.RootBuilder(examples.AddToggleGroupCommand)
//
// Or,
//
//	app := cmdr.Create(appName, version, author, desc).
//		WithBuilders(examples.AddToggleGroupCommand).
//		Build()
//	app.Run(context.TODO())
func AddToggleGroupCommand(parent cli.CommandBuilder) {
	common.AddToggleGroupCommand(parent)
}

func AddToggleGroupFlags(c cli.CommandBuilder) { //nolint:revive
	common.AddToggleGroupFlags(c)
}
