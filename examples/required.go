package examples

import (
	"github.com/hedzr/cmdr/v2/cli"
)

// AddRequiredFlag _
//
// Sample code,
//
//	app.RootBuilder(examples.AddRequiredFlag)
//
// Or,
//
//	app := cmdr.Create(appName, version, author, desc).
//		WithBuilders(examples.AddRequiredFlag).
//		Build()
//	app.Run(context.TODO())
func AddRequiredFlag(c cli.CommandBuilder) { //nolint:revive
	c.Flg("required", "r").
		Default("").
		Required(true).
		Description("the required text string wanted.", "").
		Group("Test").
		Build()
}
