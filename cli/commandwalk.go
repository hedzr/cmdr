package cli

//
// ForeachXXX and WalkXXX
//

import (
	"fmt"
	"slices"

	logz "github.com/hedzr/logg/slog"
)

// ForeachSubCommands is another way to Walk on all commands.
func (c *Command) ForeachSubCommands(cb func(cc *Command) (stop bool)) (stop bool) {
	for _, cc := range c.commands {
		if cc != nil && cb != nil {
			if stop = cb(cc); stop {
				break
			}
		}
	}

	if c.onEvalSubcommandsOnce != nil && !c.onEvalSubcommandsOnce.invoked {
	}
	if c.onEvalSubcommands != nil {
	}

	return
}

// ForeachFlags is another way to WalkEverything on all flags.
func (c *Command) ForeachFlags(cb func(f *Flag) (stop bool)) (stop bool) {
	for _, cc := range c.flags {
		if cc != nil && cb != nil {
			if stop = cb(cc); stop {
				break
			}
		}
	}

	if c.onEvalFlagsOnce != nil && !c.onEvalFlagsOnce.invoked {
	}
	if c.onEvalFlags != nil {
	}

	return
}

// ForeachGroupedSubCommands loops for all grouped commands.
//
// This function works proper except EnsureXref called.
// EnsureXref will link the whole command tree and build all
// internal indexes and maps.
func (c *Command) ForeachGroupedSubCommands(cb func(group string, items []*Command)) {
	c.ensureXrefGroup()
	for group, items := range c.allCommands {
		cb(group, items.A)
	}
}

// ForeachGroupedFlags loops for all grouped flags.
//
// This function works proper except EnsureXref called.
// EnsureXref will link the whole command tree and build all
// internal indexes and maps.
func (c *Command) ForeachGroupedFlags(cb func(group string, items []*Flag)) {
	c.ensureXrefGroup()
	for group, items := range c.allFlags {
		cb(group, items.A)
	}
}

// WalkBackwardsCtx used by WalkBackwards
type WalkBackwardsCtx struct {
	Group  bool
	Sort   bool
	hist   map[*Command]bool
	histff map[*Flag]bool
}

// WalkBackwardsCB is a callback functor used by WalkBackwards.
//
// cc and ff is the looping command and/or one of its flags. ff == nil means
// that the looping item is a command in this turn.
//
// index and groupIndex is 0-based.
// if no-sort and no-group specified, groupIndex will be -1.
//
// count is count of items in a command/flag group.
//
// level is how many times the nested flags or commands backwards to
// the root command.
type WalkBackwardsCB func(ctx *WalkBackwardsCtx, cc *Command, ff *Flag, index, groupIndex, count, level int)

type WalkCB func(cc *Command, index, level int)

type WalkGroupedCB func(cc, pp *Command, ff *Flag, group string, idx, level int)

type WalkEverythingCB func(cc, pp *Command, ff *Flag, cmdIndex, flgIndex, level int)

//

//

//

// WalkBackwards is a simple way to loop for all commands.
//
// It's specially used by Help Screen Printer, and provides
// a visiting with no-sorting but grouping turn.
//
// If you wanna control the grouping and sorting way, feel
// free about WalkBackwardsCtx.
func (c *Command) WalkBackwards(cb WalkBackwardsCB) {
	ctx := &WalkBackwardsCtx{
		Group: true,
		Sort:  false,
	}
	c.WalkBackwardsCtx(cb, ctx)
}

// WalkBackwardsCtx is a simple way to loop for all commands.
//
// It provides the visiting method with user-specified
// sorting and group way.
//
// In grouping modes, the unsorted commands and flags are
// always put at first position.
//
// In no-sort and no-group mode, the commands and flags are
// accessed with the insertion turn.
func (c *Command) WalkBackwardsCtx(cb WalkBackwardsCB, ctx *WalkBackwardsCtx) {
	if ctx.hist == nil {
		ctx.hist = make(map[*Command]bool)
	}
	if ctx.histff == nil {
		ctx.histff = make(map[*Flag]bool)
	}

	if ctx.Group {
		c.walkBackwardsImplGrouping(ctx.Sort, ctx, c, 0, cb)
		return
	}

	c.walkBackwardsImplNoGrouping(ctx.Sort, ctx, c, 0, cb)
}

// walkBackwardsImplGrouped _
// passed
func (c *Command) walkBackwardsImplGrouping(sort bool, ctx *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	if cb != nil {
		if level == 0 {
			// walk the first level command's subcommands.
			gsAndPrintCommands(sort, ctx, cmd, level, cb)
		}
		gsAndPrintFlags(sort, ctx, cmd, level, cb)
	}

	if cmd.owner != nil && cmd.owner != cmd {
		c.walkBackwardsImplGrouping(sort, ctx, cmd.owner, level+1, cb)
	}
}

// walkBackwardsImplSorted _
// passed
func (c *Command) walkBackwardsImplNoGrouping(sort bool, ctx *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	if cb != nil {
		if level == 0 {
			sortAndPrintCommands(sort, ctx, cmd, level, cb)
		}
		sortAndPrintFlags(sort, ctx, cmd, level, cb)
	}

	if cmd.owner != nil && cmd.owner != cmd {
		c.walkBackwardsImplNoGrouping(sort, ctx, cmd.owner, level+1, cb)
	}
}

func gsAndPrintCommands(sort bool, ctx *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	// walk the first level command's subcommands.
	commands := mustEnsureDynCommands(cmd)
	count := len(cmd.commands)
	m := make(map[string]map[string]*Command)
	for i, cc := range commands {
		if _, ok := ctx.hist[cc]; !ok {
			ctx.hist[cc] = true
			if _, ok = m[cc.SafeGroup()]; !ok {
				m[cc.SafeGroup()] = make(map[string]*Command)
			}
			if sort {
				m[cc.SafeGroup()][cc.Name()] = cc
			} else {
				name := fmt.Sprintf("%04d.%v", i, cc.Name())
				m[cc.SafeGroup()][name] = cc
			}
		}
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	ii := 0
	for _, grp := range keys {
		grpItems := m[grp]
		keysItems := make([]string, 0, len(grpItems))
		for k := range grpItems {
			kk := k
			keysItems = append(keysItems, kk)
		}
		slices.Sort(keysItems)
		gi := 0
		for _, it := range keysItems {
			cc := grpItems[it]
			cb(ctx, cc, nil, ii, gi, count, level)
			gi++
			ii++
		}
	}
}

func gsAndPrintFlags(sort bool, ctx *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	flags := mustEnsureDynFlags(cmd)
	count := len(flags)
	m := make(map[string]map[string]*Flag)
	for i, ff := range flags {
		if _, ok := ctx.histff[ff]; !ok {
			ctx.histff[ff] = true
			if _, ok = m[ff.SafeGroup()]; !ok {
				m[ff.SafeGroup()] = make(map[string]*Flag)
			}
			if sort {
				m[ff.SafeGroup()][ff.Name()] = ff
			} else {
				name := fmt.Sprintf("%04d.%v", i, ff.Name())
				m[ff.SafeGroup()][name] = ff
			}
		}
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	ii := 0
	for _, grp := range keys {
		grpItems := m[grp]
		keysItems := make([]string, 0, len(grpItems))
		for k := range grpItems {
			kk := k
			keysItems = append(keysItems, kk)
		}
		slices.Sort(keysItems)
		gi := 0
		for _, it := range keysItems {
			ff := grpItems[it]
			cb(ctx, cmd, ff, ii, gi, count, level)
			gi++
			ii++
		}
	}
}

func sortAndPrintCommands(sort bool, ctx *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	commands := mustEnsureDynCommands(cmd)
	count := len(commands)
	var m map[string]*Command
	if sort {
		m = make(map[string]*Command)
	}
	for i, cc := range commands {
		if _, ok := ctx.hist[cc]; !ok {
			ctx.hist[cc] = true
			if sort {
				m[cc.Name()] = cc
			} else {
				cb(ctx, cc, nil, i, -1, count, level)
			}
		}
	}
	if sort {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for i, k := range keys {
			cc := m[k]
			cb(ctx, cc, nil, i, -1, count, level)
		}
	}
}

func sortAndPrintFlags(sort bool, ctx *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	flags := mustEnsureDynFlags(cmd)
	count := len(flags)
	var m map[string]*Flag
	if sort {
		m = make(map[string]*Flag)
	}
	for i, ff := range flags {
		if _, ok := ctx.histff[ff]; !ok {
			ctx.histff[ff] = true
			if sort {
				m[ff.Name()] = ff
			} else {
				cb(ctx, cmd, ff, i, -1, count, level)
			}
		}
	}
	if sort {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for i, k := range keys {
			ff := m[k]
			cb(ctx, cmd, ff, i, -1, count, level)
		}
	}
}

// Walk is a simple way to loop for all commands in original order.
func (c *Command) Walk(cb WalkCB) {
	hist := make(map[*Command]bool)
	c.walkImpl(hist, c, 0, cb)
}

func (c *Command) walkImpl(hist map[*Command]bool, cmd *Command, level int, cb WalkCB) {
	if cb != nil {
		cb(cmd, 0, level)
	}

	commands := mustEnsureDynCommands(cmd)
	for _, cc := range commands {
		if cc != nil {
			if _, ok := hist[cc]; !ok {
				hist[cc] = true
				c.walkImpl(hist, cc, level+1, cb)
			} else {
				logz.Warn("[cmdr] loop ref found", "dad", cmd, "cc", cc)
			}
		}
	}
}

// WalkGrouped loops for all commands and its flags with grouped order.
func (c *Command) WalkGrouped(cb WalkGroupedCB) {
	hist := make(map[*Command]bool)
	c.walkGroupedImpl(hist, c, nil, 0, 0, cb)
}

func (c *Command) walkGroupedImpl(hist map[*Command]bool, dad, grandpa *Command, cmdIdx, level int, cb WalkGroupedCB) { //nolint:revive
	cb(dad, grandpa, nil, dad.GroupHelpTitle(), cmdIdx, level)

	// todo need ensure dynamic commands (and flags)

	grpKeys := make([]string, 0)
	for gg := range dad.allCommands {
		grpKeys = append(grpKeys, gg)
	}
	slices.Sort(grpKeys)

	for _, gg := range grpKeys {
		for i, cc := range dad.allCommands[gg].A {
			if _, ok := hist[cc]; !ok {
				hist[cc] = true
				// cb(cc,nil, i, 0, level)
				c.walkGroupedImpl(hist, cc, dad, i, level+1, cb)
			} else {
				logz.Warn("[cmdr] loop ref found", "dad", dad, "grandpa", grandpa, "cc", cc)
			}
		}
	}

	// for i, cc := range dad.commands {
	// 	if cc != nil {
	// 		if _, ok := hist[cc]; !ok {
	// 			hist[cc] = true
	// 			// cb(cc,nil, i, 0, level)
	// 			c.walkGrouped(hist, cc, dad, level+1, cb)
	// 		} else {
	// 			logz.Warn("[cmdr] loop ref found", "dad", dad, "grandpa", grandpa, "cc", cc)
	// 		}
	// 	}
	// }

	grpKeys = make([]string, 0)
	for gg := range dad.allFlags {
		grpKeys = append(grpKeys, gg)
	}
	slices.Sort(grpKeys)

	for _, gg := range grpKeys {
		for i, ff := range dad.allFlags[gg].A {
			cb(dad, grandpa, ff, ff.GroupHelpTitle(), i, level)
		}
	}

	// for i, ff := range dad.flags {
	// 	if ff != nil {
	// 		cb(dad, grandpa, ff, i, level)
	// 	}
	// }
}

// WalkEverything loops for all commands and its flags.
func (c *Command) WalkEverything(cb WalkEverythingCB) {
	hist := make(map[*Command]bool)
	c.walkEx(hist, c, nil, 0, 0, cb)
}

func (c *Command) walkEx(hist map[*Command]bool, dad, grandpa *Command, level, cmdIndex int, cb WalkEverythingCB) { //nolint:revive
	cb(dad, grandpa, nil, cmdIndex, 0, level)

	commands := mustEnsureDynCommands(dad)
	for i, cc := range commands {
		if cc != nil {
			if _, ok := hist[cc]; !ok {
				hist[cc] = true
				// cb(cc,nil, i, 0, level)
				c.walkEx(hist, cc, dad, level+1, i, cb)
			} else {
				logz.Warn("[cmdr] loop ref found", "dad", dad, "grandpa", grandpa, "cc", cc)
			}
		}
	}

	flags := mustEnsureDynFlags(dad)
	for i, ff := range flags {
		if ff != nil {
			cb(dad, grandpa, ff, cmdIndex, i, level)
		}
	}
}

//

func mustEnsureDynCommands(cmd *Command) (commands []*Command) {
	commands = cmd.commands
	if cclist, err := ensureDynCommands(cmd); err != nil {
		logz.Error("cannot evaluate dynamic-commands", "err", err)
	} else {
		commands = append(commands, cclist...)
	}
	return
}

func mustEnsureDynFlags(cmd *Command) (flags []*Flag) {
	flags = cmd.flags
	if cclist, err := ensureDynFlags(cmd); err != nil {
		logz.Error("cannot evaluate dynamic-flags", "err", err)
	} else {
		flags = append(flags, cclist...)
	}
	return
}

func ensureDynCommands(cmd *Command) (list []*Command, err error) {
	ctx := 0
	var c BaseOptI

	if cmd.onEvalSubcommandsOnce != nil {
		if !cmd.onEvalSubcommandsOnce.invoked {
			var iter EvalIterator
			if iter, err = cmd.onEvalSubcommandsOnce.cb(ctx, cmd); err != nil {
				return
			}

			hasNext := true
			for hasNext {
				if c, hasNext, err = iter(); err != nil {
					return
				}
				if cc, ok := c.(*Command); ok {
					list = append(list, cc)
				}
			}
			cmd.onEvalSubcommandsOnce.invoked = true
			cmd.onEvalSubcommandsOnce.commands = list
		} else {
			list = cmd.onEvalSubcommandsOnce.commands
		}
	}

	if cmd.onEvalSubcommands != nil {
		var iter EvalIterator
		if iter, err = cmd.onEvalSubcommands.cb(ctx, cmd); err != nil {
			return
		}

		hasNext := true
		for hasNext {
			if c, hasNext, err = iter(); err != nil {
				return
			}
			if cc, ok := c.(*Command); ok {
				list = append(list, cc)
			}
		}
	}
	return
}

func ensureDynFlags(cmd *Command) (list []*Flag, err error) {
	ctx := 0
	var c BaseOptI

	if cmd.onEvalFlagsOnce != nil {
		if !cmd.onEvalFlagsOnce.invoked {
			var iter EvalIterator
			if iter, err = cmd.onEvalFlagsOnce.cb(ctx, cmd); err != nil {
				return
			}

			hasNext := true
			for hasNext {
				if c, hasNext, err = iter(); err != nil {
					return
				}
				if cc, ok := c.(*Flag); ok {
					list = append(list, cc)
				}
			}
			cmd.onEvalFlagsOnce.invoked = true
			cmd.onEvalFlagsOnce.flags = list
		} else {
			list = cmd.onEvalFlagsOnce.flags
		}
	}

	if cmd.onEvalFlags != nil {
		var iter EvalIterator
		if iter, err = cmd.onEvalFlags.cb(ctx, cmd); err != nil {
			return
		}

		hasNext := true
		for hasNext {
			if c, hasNext, err = iter(); err != nil {
				return
			}
			if cc, ok := c.(*Flag); ok {
				list = append(list, cc)
			}
		}
	}
	return
}
