// Copyright © 2022 Hedzr Yeh.

package cli

import (
	"context"
	"fmt"
	"time"
)

func (c *RootCommand) SetApp(app App) *RootCommand {
	c.app = app
	c.root = c
	if a, ok := app.(interface {
		SetRootCommand(command *RootCommand) App
	}); ok {
		a.SetRootCommand(c)
	}
	return c
}

func (c *RootCommand) NewCmd(longTitle string) *Command {
	cc := &Command{
		BaseOpt: BaseOpt{
			Long:  longTitle,
			owner: c.Command,
			root:  c,
		},
	}
	c.AddSubCommand(cc)
	return cc
}

func (c *RootCommand) NewFlg(longTitle string) *Flag {
	cc := &Flag{
		BaseOpt: BaseOpt{
			Long:  longTitle,
			owner: c.root.Command,
			root:  c.root,
		},
	}
	c.AddFlag(cc)
	return cc
}

// Attach attaches new root command on it
func (c *RootCommand) Attach(newRootCommand *Command) {
	c.Command = newRootCommand
	newRootCommand.owner = nil
	newRootCommand.root = c

	// update 'root' fields in any commands/flags
	hist := make(map[BaseOptI]bool)
	ctx := context.TODO()
	c.walkImpl(ctx, hist, c.Command, 0, func(cc BaseOptI, index, level int) {
		if cx, ok := cc.(interface{ SetRoot(command *RootCommand) }); ok {
			cx.SetRoot(c)
		}
		cc.ForeachFlags(func(f *Flag) (stop bool) { f.root = c; return })
	})

	// c.attachBuiltinCommands()
}

//
//
//

// AppendPostActions adds the global post-action to cmdr system
func (c *RootCommand) AppendPostActions(functions ...OnPostInvokeHandler) {
	c.postActions = append(c.postActions, functions...)
}

// AppendPreActions adds the global pre-action to cmdr system
func (c *RootCommand) AppendPreActions(functions ...OnPreInvokeHandler) {
	c.preActions = append(c.preActions, functions...)
}

func (c *RootCommand) AppDescription() string     { return c.Desc() }
func (c *RootCommand) AppLongDescription() string { return c.DescLong() }

func (c *RootCommand) Header() string {
	if c.HeaderLine != "" {
		return c.HeaderLine
	}

	if c.Copyright == "" {
		// c.Copyright = fmt.Sprintf("Copyright (c) %v", time.Now().Year())
		c.Copyright = fmt.Sprintf("Copyright © %v", time.Now().Year())
	}
	if c.Author == "" {
		c.Author = fmt.Sprintf("%v Authors", c.AppName)
	}
	return fmt.Sprintf("%v v%v ~ %v by %v ~ All Rights Reserved.",
		c.AppName, c.AppVersion(), c.Copyright, c.Author)
}

func (c *RootCommand) Footer() string {
	if c.FooterLine != "" {
		return c.FooterLine
	}
	return defaultTailLine
}

const (
	defaultTailLine = `
Type '-h'/'-?' or '--help' to get command help screen. 
More: '-D'/'--debug'['--env'|'--raw'|'--more'], '-V'/'--version', '-#'/'--build-info', '--no-color', '--strict-mode', '--no-env-overrides'...`
)
