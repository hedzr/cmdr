package cli

//
// ForeachXXX and WalkXXX
//

import (
	"context"
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
type WalkBackwardsCB func(ctx context.Context, pc *WalkBackwardsCtx, cc *Command, ff *Flag, index, groupIndex, count, level int)

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
func (c *Command) WalkBackwards(ctx context.Context, cb WalkBackwardsCB) {
	pc := &WalkBackwardsCtx{
		Group: true,
		Sort:  false,
	}
	c.WalkBackwardsCtx(ctx, cb, pc)
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
func (c *Command) WalkBackwardsCtx(ctx context.Context, cb WalkBackwardsCB, pc *WalkBackwardsCtx) {
	if pc.hist == nil {
		pc.hist = make(map[*Command]bool)
	}
	if pc.histff == nil {
		pc.histff = make(map[*Flag]bool)
	}

	if pc.Group {
		c.walkBackwardsImplGrouping(ctx, pc.Sort, pc, c, 0, cb)
		return
	}

	c.walkBackwardsImplNoGrouping(ctx, pc.Sort, pc, c, 0, cb)
}

// walkBackwardsImplGrouped _
// passed
func (c *Command) walkBackwardsImplGrouping(ctx context.Context, sort bool, pc *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	if cb != nil {
		if level == 0 {
			// walk the first level command's subcommands.
			gsAndPrintCommands(ctx, sort, pc, cmd, level, cb)
		}
		gsAndPrintFlags(ctx, sort, pc, cmd, level, cb)
	}

	if cmd.owner != nil && cmd.owner != cmd {
		c.walkBackwardsImplGrouping(ctx, sort, pc, cmd.owner, level+1, cb)
	}
}

// walkBackwardsImplSorted _
// passed
func (c *Command) walkBackwardsImplNoGrouping(ctx context.Context, sort bool, pc *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	if cb != nil {
		if level == 0 {
			sortAndPrintCommands(ctx, sort, pc, cmd, level, cb)
		}
		sortAndPrintFlags(ctx, sort, pc, cmd, level, cb)
	}

	if cmd.owner != nil && cmd.owner != cmd {
		c.walkBackwardsImplNoGrouping(ctx, sort, pc, cmd.owner, level+1, cb)
	}
}

func gsAndPrintCommands(ctx context.Context, sort bool, pc *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	// walk the first level command's subcommands.
	commands := mustEnsureDynCommands(ctx, cmd)
	count := len(cmd.commands)
	m := make(map[string]map[string]*Command)
	for i, cc := range commands {
		if _, ok := pc.hist[cc]; !ok {
			pc.hist[cc] = true
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
			cb(ctx, pc, cc, nil, ii, gi, count, level)
			gi++
			ii++
		}
	}
}

func gsAndPrintFlags(ctx context.Context, sort bool, pc *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	flags := mustEnsureDynFlags(ctx, cmd)
	count := len(flags)
	m := make(map[string]map[string]*Flag)
	for i, ff := range flags {
		if _, ok := pc.histff[ff]; !ok {
			pc.histff[ff] = true
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
			cb(ctx, pc, cmd, ff, ii, gi, count, level)
			gi++
			ii++
		}
	}
}

func sortAndPrintCommands(ctx context.Context, sort bool, pc *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	commands := mustEnsureDynCommands(ctx, cmd)
	count := len(commands)
	var m map[string]*Command
	if sort {
		m = make(map[string]*Command)
	}
	for i, cc := range commands {
		if _, ok := pc.hist[cc]; !ok {
			pc.hist[cc] = true
			if sort {
				m[cc.Name()] = cc
			} else {
				cb(ctx, pc, cc, nil, i, -1, count, level)
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
			cb(ctx, pc, cc, nil, i, -1, count, level)
		}
	}
}

func sortAndPrintFlags(ctx context.Context, sort bool, pc *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	flags := mustEnsureDynFlags(ctx, cmd)
	count := len(flags)
	var m map[string]*Flag
	if sort {
		m = make(map[string]*Flag)
	}
	for i, ff := range flags {
		if _, ok := pc.histff[ff]; !ok {
			pc.histff[ff] = true
			if sort {
				m[ff.Name()] = ff
			} else {
				cb(ctx, pc, cmd, ff, i, -1, count, level)
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
			cb(ctx, pc, cmd, ff, i, -1, count, level)
		}
	}
}

// Walk is a simple way to loop for all commands in original order.
func (c *Command) Walk(ctx context.Context, cb WalkCB) {
	hist := make(map[*Command]bool)
	c.walkImpl(ctx, hist, c, 0, cb)
}

func (c *Command) walkImpl(ctx context.Context, hist map[*Command]bool, cmd *Command, level int, cb WalkCB) {
	if cb != nil {
		cb(cmd, 0, level)
	}

	commands := mustEnsureDynCommands(ctx, cmd)
	for _, cc := range commands {
		if cc != nil {
			if _, ok := hist[cc]; !ok {
				hist[cc] = true
				c.walkImpl(ctx, hist, cc, level+1, cb)
			} else {
				logz.Warn("[cmdr] loop ref found", "dad", cmd, "cc", cc)
			}
		}
	}
}

// WalkGrouped loops for all commands and its flags with grouped order.
func (c *Command) WalkGrouped(ctx context.Context, cb WalkGroupedCB) {
	hist := make(map[*Command]bool)
	c.walkGroupedImpl(ctx, hist, c, nil, 0, 0, cb)
}

func (c *Command) walkGroupedImpl(ctx context.Context, hist map[*Command]bool, dad, grandpa *Command, cmdIdx, level int, cb WalkGroupedCB) { //nolint:revive
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
				c.walkGroupedImpl(ctx, hist, cc, dad, i, level+1, cb)
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
func (c *Command) WalkEverything(ctx context.Context, cb WalkEverythingCB) {
	hist := make(map[*Command]bool)
	c.walkEx(ctx, hist, c, nil, 0, 0, cb)
}

func (c *Command) walkEx(ctx context.Context, hist map[*Command]bool, dad, grandpa *Command, level, cmdIndex int, cb WalkEverythingCB) { //nolint:revive
	cb(dad, grandpa, nil, cmdIndex, 0, level)

	commands := mustEnsureDynCommands(ctx, dad)
	for i, cc := range commands {
		if cc != nil {
			if _, ok := hist[cc]; !ok {
				hist[cc] = true
				// cb(cc,nil, i, 0, level)
				c.walkEx(ctx, hist, cc, dad, level+1, i, cb)
			} else {
				logz.Warn("[cmdr] loop ref found", "dad", dad, "grandpa", grandpa, "cc", cc)
			}
		}
	}

	flags := mustEnsureDynFlags(ctx, dad)
	for i, ff := range flags {
		if ff != nil {
			cb(dad, grandpa, ff, cmdIndex, i, level)
		}
	}
}

//

func mustEnsureDynCommands(ctx context.Context, cmd *Command) (commands []*Command) {
	commands = cmd.commands
	if cclist, err := ensureDynCommands(ctx, cmd); err != nil {
		logz.Error("cannot evaluate dynamic-commands", "err", err)
	} else {
		commands = append(commands, cclist...)
	}
	return
}

func mustEnsureDynFlags(ctx context.Context, cmd *Command) (flags []*Flag) {
	flags = cmd.flags
	if cclist, err := ensureDynFlags(ctx, cmd); err != nil {
		logz.Error("cannot evaluate dynamic-flags", "err", err)
	} else {
		flags = append(flags, cclist...)
	}
	return
}

func ensureDynCommands(ctx context.Context, cmd *Command) (list []*Command, err error) {
	var c BaseOptI

	if cmd.onEvalSubcommandsOnce != nil {
		if !cmd.onEvalSubcommandsOnce.invoked {
			var iter EvalIterator
			if iter, err = cmd.onEvalSubcommandsOnce.cb(ctx, cmd); err != nil {
				return
			}

			hasNext := true
			for hasNext {
				if c, hasNext, err = iter(ctx); err != nil {
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
		if iter, err = cmd.onEvalSubcommands.cb(ctx, cmd); err != nil || iter == nil {
			return
		}

		hasNext := true
		for hasNext {
			if c, hasNext, err = iter(ctx); err != nil {
				return
			}
			if cc, ok := c.(*Command); ok {
				list = append(list, cc)
			}
		}
	}
	return
}

func ensureDynFlags(ctx context.Context, cmd *Command) (list []*Flag, err error) {
	var c BaseOptI

	if cmd.onEvalFlagsOnce != nil {
		if !cmd.onEvalFlagsOnce.invoked {
			var iter EvalIterator
			if iter, err = cmd.onEvalFlagsOnce.cb(ctx, cmd); err != nil {
				return
			}

			hasNext := true
			for hasNext {
				if c, hasNext, err = iter(ctx); err != nil {
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
		if iter, err = cmd.onEvalFlags.cb(ctx, cmd); err != nil || iter == nil {
			return
		}

		hasNext := true
		for hasNext {
			if c, hasNext, err = iter(ctx); err != nil {
				return
			}
			if cc, ok := c.(*Flag); ok {
				list = append(list, cc)
			}
		}
	}
	return
}
