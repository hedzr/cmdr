package cli

import (
	"context"
)

// FindSubCommand find sub-command with `longName` from `cmd`.
//
// If wide is true, FindSubCommand will try to match
// long, short and aliases titles,
// or it matches only long title.
func (c *CmdS) FindSubCommand(ctx context.Context, longName string, wide bool) (res Cmd) {
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
		if longName == cx.Name() {
			return cx
		}
	}
	return
}

// FindSubCommandRecursive find sub-command with `longName` from `cmd` recursively
func (c *CmdS) FindSubCommandRecursive(ctx context.Context, longName string, wide bool) (res Cmd) { //nolint:revive
	// return FindSubCommandRecursive(longName, c)
	commands := mustEnsureDynCommands(ctx, c)
	res = c.findSubCommandIn(ctx, c, commands, longName, wide)
	return
}

func (c *CmdS) findSubCommandIn(ctx context.Context, cc Cmd, children []Cmd, longName string, wide bool) (res Cmd) { //nolint:revive
	if wide {
		if r, ok := c.longCommands[longName]; ok {
			return r
		}
		if r, ok := c.shortCommands[longName]; ok {
			return r
		}
	}
	for _, cx := range children {
		if longName == cx.Name() {
			return cx
		}
		cclist := mustEnsureDynCommands(ctx, cx)
		if len(cclist) > 0 {
			if k, ok := cx.(CmdPriv); ok {
				if res = k.findSubCommandIn(ctx, cx, cclist, longName, wide); res != nil {
					return
				}
			}
		}
	}
	return
}

// FindFlag find flag with `longName` from `cmd`
func (c *CmdS) FindFlag(ctx context.Context, longName string, wide bool) (res *Flag) {
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
func (c *CmdS) FindFlagRecursive(ctx context.Context, longName string, wide bool) (res *Flag) {
	commands := mustEnsureDynCommands(ctx, c)
	return c.findFlagIn(ctx, c, commands, longName, wide)
}

func (c *CmdS) findFlagIn(ctx context.Context, cc Cmd, children []Cmd, longName string, wide bool) (res *Flag) {
	// return FindFlagRecursive(longName, c)

	// TODO Cmd.longFlags
	// if wide {
	// 	if r, ok := cc.longFlags[longName]; ok {
	// 		return r
	// 	}
	// 	if r, ok := cc.shortFlags[longName]; ok {
	// 		return r
	// 	}
	// }
	flags := mustEnsureDynFlags(ctx, cc)
	for _, cx := range flags {
		if longName == cx.Long {
			return cx
		}
	}

	commands := mustEnsureDynCommands(ctx, c)
	for _, cx := range commands {
		if k, ok := cx.(CmdPriv); ok {
			if res = k.findFlagIn(ctx, cx, commands, longName, false); res != nil {
				return
			}
		}
	}
	return
}

func (c *CmdS) FindFlagBackwards(ctx context.Context, longName string) (res *Flag) {
	commands := mustEnsureDynCommands(ctx, c)
	return c.findFlagBackwardsIn(ctx, c, commands, longName)
}

func (c *CmdS) findFlagBackwardsIn(ctx context.Context, cc Cmd, children []Cmd, longName string) (res *Flag) {
	for _, cx := range c.flags {
		if longName == cx.Long {
			res = cx
			return
		}
	}
	if pp := c.owner; pp != nil && pp != c {
		commands := mustEnsureDynCommands(ctx, pp)
		if pf, ok := pp.(CmdPriv); ok {
			res = pf.findFlagBackwardsIn(ctx, pp, commands, longName)
		}
	}
	_, _ = cc, children
	return
}
