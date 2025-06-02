// Copyright © 2022 Hedzr Yeh.

package cli

import (
	"context"
	"fmt"
	"time"
)

func (c *RootCommand) SetApp(app App) *RootCommand {
	c.app = app
	c.SetRoot(c) // c.root = c
	if a, ok := app.(interface {
		SetRootCommand(command *RootCommand) App
	}); ok {
		a.SetRootCommand(c)
	}
	return c
}

func (c *RootCommand) RedirectToSet() map[string]map[Cmd][]*CmdS {
	return c.redirectCmds
}

func (c *RootCommand) NewCmd(longTitle string) *CmdS {
	cc := &CmdS{
		BaseOpt: BaseOpt{
			Long:  longTitle,
			owner: c.Cmd,
			root:  c,
		},
	}
	if x, ok := c.Cmd.(interface {
		AddSubCommand(child *CmdS, callbacks ...func(cc *CmdS))
	}); ok {
		x.AddSubCommand(cc)
	}
	return cc
}

func (c *RootCommand) NewFlg(longTitle string) *Flag {
	cc := &Flag{
		BaseOpt: BaseOpt{
			Long:  longTitle,
			owner: c.Root(),
			root:  c.Root(),
		},
	}
	if x, ok := c.Cmd.(interface {
		AddFlag(ff *Flag, callbacks ...func(ff *Flag))
	}); ok {
		x.AddFlag(cc)
	}
	return cc
}

// Attach attaches new root command on it
func (c *RootCommand) Attach(newRootCommand Cmd) {
	c.Cmd = newRootCommand
	newRootCommand.SetOwnerCmd(nil)
	newRootCommand.SetRoot(c)

	// update 'root' fields in any commands/flags
	// hist := make(map[Cmd]bool)
	ctx := context.TODO()
	c.Walk(ctx, func(cc Cmd, index, level int) {
		if cx, ok := cc.(interface{ SetRoot(command *RootCommand) }); ok {
			cx.SetRoot(c)
		}
		cc.ForeachFlags(ctx, func(f *Flag) (stop bool) { f.root = c; return })
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
	ver := c.Version
	if ver != "" && ver[0] == 'v' {
		ver = ver[1:]
	} else if ver == "" {
		ver = "?"
	}
	return fmt.Sprintf("%v v%v ~ %v by %v ~ All Rights Reserved.",
		c.AppName, ver, c.Copyright, c.Author)
}

func (c *RootCommand) Footer() string {
	if c.FooterLine != "" {
		return c.FooterLine
	}
	return defaultTailLine
}

const (
	defaultTailLine = `
Type '-h'/'-?' or '--help' to get this help screen ({{.Cols}}x{{.Rows}}/{{.Tabstop}}).
More: '-D'/'--debug', '-V'/'--version', '-#'/'--build-info', '--no-color'...`
)
