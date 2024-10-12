package cli

// FindSubCommand find sub-command with `longName` from `cmd`.
//
// If wide is true, FindSubCommand will try to match
// long, short and aliases titles,
// or it matches only long title.
func (c *Command) FindSubCommand(longName string, wide bool) (res *Command) {
	// return FindSubCommand(longName, c)
	if wide {
		if r, ok := c.longCommands[longName]; ok {
			return r
		}
		if r, ok := c.shortCommands[longName]; ok {
			return r
		}
	}
	commands := mustEnsureDynCommands(c)
	for _, cx := range commands {
		if longName == cx.Long {
			return cx
		}
	}
	return
}

// FindSubCommandRecursive find sub-command with `longName` from `cmd` recursively
func (c *Command) FindSubCommandRecursive(longName string, wide bool) (res *Command) { //nolint:revive
	// return FindSubCommandRecursive(longName, c)
	commands := mustEnsureDynCommands(c)
	res = c.findSubCommandIn(c, commands, longName, wide)
	return
}

func (c *Command) findSubCommandIn(cc *Command, children []*Command, longName string, wide bool) (res *Command) { //nolint:revive
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
		cclist := mustEnsureDynCommands(cx)
		if len(cclist) > 0 {
			if res = cx.findSubCommandIn(cx, cclist, longName, wide); res != nil {
				return
			}
		}
	}
	return
}

// FindFlag find flag with `longName` from `cmd`
func (c *Command) FindFlag(longName string, wide bool) (res *Flag) {
	// return FindFlag(longName, c)
	if wide {
		if r, ok := c.longFlags[longName]; ok {
			return r
		}
		if r, ok := c.shortFlags[longName]; ok {
			return r
		}
	}
	flags := mustEnsureDynFlags(c)
	for _, cx := range flags {
		if longName == cx.Long {
			return cx
		}
	}
	return
}

// FindFlagRecursive find flag with `longName` from `cmd` recursively
func (c *Command) FindFlagRecursive(longName string, wide bool) (res *Flag) {
	commands := mustEnsureDynCommands(c)
	return c.findFlagIn(c, commands, longName, wide)
}

func (c *Command) findFlagIn(cc *Command, children []*Command, longName string, wide bool) (res *Flag) {
	// return FindFlagRecursive(longName, c)
	if wide {
		if r, ok := cc.longFlags[longName]; ok {
			return r
		}
		if r, ok := cc.shortFlags[longName]; ok {
			return r
		}
	}
	flags := mustEnsureDynFlags(cc)
	for _, cx := range flags {
		if longName == cx.Long {
			return cx
		}
	}

	commands := mustEnsureDynCommands(c)
	for _, cx := range commands {
		// if len(cx.SubCommands) > 0 {
		if res = cx.findFlagIn(cx, commands, longName, false); res != nil {
			return
		}
		// }
	}
	return
}

func (c *Command) FindFlagBackwards(longName string) (res *Flag) {
	commands := mustEnsureDynCommands(c)
	return c.findFlagBackwardsIn(c, commands, longName)
}

func (c *Command) findFlagBackwardsIn(cc *Command, children []*Command, longName string) (res *Flag) {
	for _, cx := range c.flags {
		if longName == cx.Long {
			res = cx
			return
		}
	}
	if c.owner != nil && c.owner != c {
		commands := mustEnsureDynCommands(c.owner)
		res = c.owner.findFlagBackwardsIn(c.owner, commands, longName)
	}
	return
}
