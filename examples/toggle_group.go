package examples

import (
	"context"
	"fmt"

	"github.com/hedzr/cmdr/v2/cli"
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
func AddToggleGroupCommand(parent cli.CommandBuilder) { //nolint:revive
	// toggle-group-test - without a default choice

	parent.Cmd("tg-test", "tg").
		Description("toggable group, with a default choice", "tg test new features,\nverbose long descriptions here.").
		Group("Toggleable Group").
		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			fmt.Printf("*** Got fruit (toggle group): %q\n", cmd.Store().MustString(toggleGroupName))
			fmt.Printf("         multiple selections: %q\n", cmd.Store().MustStringSlice(toggleGroupName+".selected"))

			println()
			// fmt.Printf("> STDIN MODE: %v \n", cmd.Set().MustBool("mx-test.stdin"))
			// fmt.Println()
			return
		}).
		With(func(cb cli.CommandBuilder) {
			cb.Flg("apple", "").
				Default(false).
				Description("the test text.", "").
				ToggleGroup(toggleGroupName).
				Build()
			cb.Flg("banana", "").
				Default(false).
				Description("the test text.", "").
				ToggleGroup(toggleGroupName).
				Build()
			cb.Flg("orange", "").
				Default(true).
				Description("the test text.", "").
				ToggleGroup(toggleGroupName).
				Build()
		})

	// tg2 - with a default choice

	parent.Cmd("tg-test2", "tg2", "toggle-group-test2").
		Description("toggable group 2, without default choice", "tg2 test new features,\nverbose long descriptions here.").
		Group("Toggleable Group").
		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
			fmt.Printf("*** Got fruit (toggle group): %q\n", cmd.Store().MustString(toggleGroupName))
			fmt.Printf("         multiple selections: %q\n", cmd.Store().MustStringSlice(toggleGroupName+".selected"))

			println()
			// fmt.Printf("> STDIN MODE: %v \n", cmd.Set().MustBool("mx-test.stdin"))
			// fmt.Println()
			return
		}).
		With(func(cb cli.CommandBuilder) {
			cb.Flg("apple", "a").
				Default(true).
				Description("the test text.", "").
				ToggleGroup(toggleGroupName).
				Build()
			cb.Flg("banana", "b").
				Default(false).
				Description("the test text.", "").
				ToggleGroup(toggleGroupName).
				Build()
			cb.Flg("orange", "o").
				Default(false).
				Description("the test text.", "").
				ToggleGroup(toggleGroupName).
				Build()
		})
}

func AddToggleGroupFlags(c cli.CommandBuilder) { //nolint:revive
	c.Flg("apple", "a").
		Default(false).
		Description("the test text.", "").
		ToggleGroup(toggleGroupName).
		Group("Toggleable Group").
		Build()

	c.Flg("banana", "b").
		Default(false).
		Description("the test text.", "").
		ToggleGroup(toggleGroupName).
		Group("Toggleable Group").
		Build()

	c.Flg("orange", "o").
		Default(true).
		Description("the test text.", "").
		ToggleGroup(toggleGroupName).
		Group("Toggleable Group").
		Build()
}

const toggleGroupName = "fruit"
