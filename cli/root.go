// Copyright © 2022 Hedzr Yeh.

package cli

import (
	"fmt"
	"time"
)

// Attach attaches new root command on it
func (c *RootCommand) Attach(newRootCommand *Command) {
	c.Command = newRootCommand
	newRootCommand.owner = nil
	newRootCommand.root = c

	// update 'root' fields in any commands/flags
	hist := make(map[*Command]bool)
	c.walkImpl(hist, c.Command, 0, func(cc *Command, index, level int) {
		cc.root = c
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
