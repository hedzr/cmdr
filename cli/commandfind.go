package cli

import (
	"context"
)

// FindSubCommand find sub-command with `longName` from `cmd`.
//
// If wide is true, FindSubCommand will try to match
// long, short and aliases titles,
// or it matches only long title.
func (c *Command) FindSubCommand(ctx context.Context, longName string, wide bool) (res *Command) {
	// return FindSubCommand(longName, c)
	if wide {
		if r, ok := c.longCommands[longName]; ok {
			return r
		}
		if r, ok := c.shortCommands[longName]; ok {
			return r
		}
	}
	commands := mustEnsureDynCommands(ctx, c)
	for _, cx := range commands {
		if longName == cx.Long {
			return cx
		}
	}
	return
}

// FindSubCommandRecursive find sub-command with `longName` from `cmd` recursively
func (c *Command) FindSubCommandRecursive(ctx context.Context, longName string, wide bool) (res *Command) { //nolint:revive
	// return FindSubCommandRecursive(longName, c)
	commands := mustEnsureDynCommands(ctx, c)
	res = c.findSubCommandIn(ctx, c, commands, longName, wide)
	return
}

func (c *Command) findSubCommandIn(ctx context.Context, cc *Command, children []*Command, longName string, wide bool) (res *Command) { //nolint:revive
	if wide {
		if r, ok := c.longCommands[longName]; ok {
			return r
		}
		if r, ok := c.shortCommands[longName]; ok {
			return r
		}
	}
	for _, cx := range children {
		if longName == cx.Long {
			return cx
		}
		cclist := mustEnsureDynCommands(ctx, cx)
		if len(cclist) > 0 {
			if res = cx.findSubCommandIn(ctx, cx, cclist, longName, wide); res != nil {
				return
			}
		}
	}
	return
}

// FindFlag find flag with `longName` from `cmd`
func (c *Command) FindFlag(ctx context.Context, longName string, wide bool) (res *Flag) {
	// return FindFlag(longName, c)
	if wide {
		if r, ok := c.longFlags[longName]; ok {
			return r
		}
		if r, ok := c.shortFlags[longName]; ok {
			return r
		}
	}
	flags := mustEnsureDynFlags(ctx, c)
	for _, cx := range flags {
		if longName == cx.Long {
			return cx
		}
	}
	return
}

// FindFlagRecursive find flag with `longName` from `cmd` recursively
func (c *Command) FindFlagRecursive(ctx context.Context, longName string, wide bool) (res *Flag) {
	commands := mustEnsureDynCommands(ctx, c)
	return c.findFlagIn(ctx, c, commands, longName, wide)
}

func (c *Command) findFlagIn(ctx context.Context, cc *Command, children []*Command, longName string, wide bool) (res *Flag) {
	// return FindFlagRecursive(longName, c)
	if wide {
		if r, ok := cc.longFlags[longName]; ok {
			return r
		}
		if r, ok := cc.shortFlags[longName]; ok {
			return r
		}
	}
	flags := mustEnsureDynFlags(ctx, cc)
	for _, cx := range flags {
		if longName == cx.Long {
			return cx
		}
	}

	commands := mustEnsureDynCommands(ctx, c)
	for _, cx := range commands {
		// if len(cx.SubCommands) > 0 {
		if res = cx.findFlagIn(ctx, cx, commands, longName, false); res != nil {
			return
		}
		// }
	}
	return
}

func (c *Command) FindFlagBackwards(ctx context.Context, longName string) (res *Flag) {
	commands := mustEnsureDynCommands(ctx, c)
	return c.findFlagBackwardsIn(ctx, c, commands, longName)
}

func (c *Command) findFlagBackwardsIn(ctx context.Context, cc *Command, children []*Command, longName string) (res *Flag) {
	for _, cx := range c.flags {
		if longName == cx.Long {
			res = cx
			return
		}
	}
	if c.owner != nil && c.owner != c {
		commands := mustEnsureDynCommands(ctx, c.owner)
		res = c.owner.findFlagBackwardsIn(ctx, c.owner, commands, longName)
	}
	return
}
