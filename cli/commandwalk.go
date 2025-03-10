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
func (c *CmdS) ForeachSubCommands(ctx context.Context, cb func(cc *CmdS) (stop bool)) (stop bool) {
	for _, cc := range c.commands {
		if cc != nil && cb != nil {
			if stop = cb(cc); stop {
				break
			}
		}
	}

	if c.onEvalSubcommandsOnce != nil && !c.onEvalSubcommandsOnce.invoked {
		// todo onEvalSubcommandsOnce
	}
	if c.onEvalSubcommands != nil {
		// todo onEvalSubcommands
	}

	return
}

// ForeachFlags is another way to WalkEverything on all flags.
func (c *CmdS) ForeachFlags(ctx context.Context, cb func(f *Flag) (stop bool)) (stop bool) {
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
func (c *CmdS) ForeachGroupedSubCommands(ctx context.Context, cb func(group string, items []*CmdS)) {
	c.ensureXrefGroup(ctx)
	for group, items := range c.allCommands {
		cb(group, items.A)
	}
}

// ForeachGroupedFlags loops for all grouped flags.
//
// This function works proper except EnsureXref called.
// EnsureXref will link the whole command tree and build all
// internal indexes and maps.
func (c *CmdS) ForeachGroupedFlags(ctx context.Context, cb func(group string, items []*Flag)) {
	c.ensureXrefGroup(ctx)
	for group, items := range c.allFlags {
		cb(group, items.A)
	}
}

// WalkBackwardsCtx used by WalkBackwards
type WalkBackwardsCtx struct {
	Group  bool
	Sort   bool
	hist   map[Cmd]bool
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
type WalkBackwardsCB func(ctx context.Context, pc *WalkBackwardsCtx, cc Cmd, ff *Flag, index, groupIndex, count, level int)

type WalkCB func(cc Cmd, index, level int)

type WalkFastCB func(cc Cmd, index, level int) (stop bool)

type WalkGroupedCB func(cc, pp Cmd, ff *Flag, group string, idx, level int)

type WalkEverythingCB func(cc, pp Cmd, ff *Flag, cmdIndex, flgIndex, level int)

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
func (c *CmdS) WalkBackwards(ctx context.Context, cb WalkBackwardsCB) {
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
func (c *CmdS) WalkBackwardsCtx(ctx context.Context, cb WalkBackwardsCB, pc *WalkBackwardsCtx) {
	if pc.hist == nil {
		pc.hist = make(map[Cmd]bool)
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
func (c *CmdS) walkBackwardsImplGrouping(ctx context.Context, sort bool, pc *WalkBackwardsCtx, cmd Cmd, level int, cb WalkBackwardsCB) {
	if cb != nil {
		if level == 0 {
			// walk the first level command's subcommands.
			gsAndPrintCommands(ctx, sort, pc, cmd, level, cb)
		}
		gsAndPrintFlags(ctx, sort, pc, cmd, level, cb)
	}

	if cmd.OwnerIsValid() {
		c.walkBackwardsImplGrouping(ctx, sort, pc, cmd.OwnerCmd(), level+1, cb)
	}
}

// walkBackwardsImplSorted _
// passed
func (c *CmdS) walkBackwardsImplNoGrouping(ctx context.Context, sort bool, pc *WalkBackwardsCtx, cmd Cmd, level int, cb WalkBackwardsCB) {
	if cb != nil {
		if level == 0 {
			sortAndPrintCommands(ctx, sort, pc, cmd, level, cb)
		}
		sortAndPrintFlags(ctx, sort, pc, cmd, level, cb)
	}

	if cmd.OwnerIsValid() {
		c.walkBackwardsImplNoGrouping(ctx, sort, pc, cmd.OwnerCmd(), level+1, cb)
	}
}

func gsAndPrintCommands(ctx context.Context, sort bool, pc *WalkBackwardsCtx, cmd Cmd, level int, cb WalkBackwardsCB) {
	// walk the first level command's subcommands.
	commands := mustEnsureDynCommands(ctx, cmd)
	count := len(cmd.SubCommands())
	m := make(map[string]map[string]Cmd)
	for i, cc := range commands {
		if _, ok := pc.hist[cc]; !ok {
			pc.hist[cc] = true
			if _, ok = m[cc.SafeGroup()]; !ok {
				m[cc.SafeGroup()] = make(map[string]Cmd)
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

func gsAndPrintFlags(ctx context.Context, sort bool, pc *WalkBackwardsCtx, cmd Cmd, level int, cb WalkBackwardsCB) {
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

func sortAndPrintCommands(ctx context.Context, sort bool, pc *WalkBackwardsCtx, cmd Cmd, level int, cb WalkBackwardsCB) {
	commands := mustEnsureDynCommands(ctx, cmd)
	count := len(commands)
	var m map[string]Cmd
	if sort {
		m = make(map[string]Cmd)
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

func sortAndPrintFlags(ctx context.Context, sort bool, pc *WalkBackwardsCtx, cmd Cmd, level int, cb WalkBackwardsCB) {
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

// WalkFast is a simple way to loop for all commands in original order.
func (c *CmdS) WalkFast(ctx context.Context, cb WalkFastCB) (stop bool) {
	hist := make(map[Cmd]bool)
	return c.walkFastImpl(ctx, hist, c, 0, cb)
}

func (c *CmdS) walkFastImpl(ctx context.Context, hist map[Cmd]bool, cmd Cmd, level int, cb WalkFastCB) (stop bool) {
	if cb != nil {
		cb(cmd, 0, level)
	}

	commands := mustEnsureDynCommands(ctx, cmd)
	for _, cc := range commands {
		if cc != nil {
			if _, ok := hist[cc]; !ok {
				hist[cc] = true
				if stop = c.walkFastImpl(ctx, hist, cc, level+1, cb); stop {
					return
				}
			} else {
				logz.WarnContext(ctx, "[cmdr] loop ref found", "dad", cmd, "cc", cc)
			}
		}
	}
	return
}

// Walk is a simple way to loop for all commands in original order.
func (c *CmdS) Walk(ctx context.Context, cb WalkCB) {
	hist := make(map[Cmd]bool)
	c.walkImpl(ctx, hist, c, 0, cb)
}

func (c *CmdS) walkImpl(ctx context.Context, hist map[Cmd]bool, cmd Cmd, level int, cb WalkCB) {
	if cb != nil {
		cb(cmd, 0, level)
	}

	// for _, cc := range cmd.SubCommands() {
	// 	if cc != nil {
	// 		if _, ok := hist[cc]; !ok {
	// 			hist[cc] = true
	// 			c.walkImpl(ctx, hist, cc, level+1, cb)
	// 		} else {
	// 			logz.WarnContext(ctx, "[cmdr] loop ref found", "dad", cmd, "cc", cc)
	// 		}
	// 	}
	// }

	commands := mustEnsureDynCommands(ctx, cmd)
	for _, cc := range commands {
		if cc != nil {
			if _, ok := hist[cc]; !ok {
				hist[cc] = true
				c.walkImpl(ctx, hist, cc, level+1, cb)
			} else {
				logz.WarnContext(ctx, "[cmdr] loop ref found", "dad", cmd, "cc", cc)
			}
		}
	}
}

// WalkGrouped loops for all commands and its flags with grouped order.
func (c *CmdS) WalkGrouped(ctx context.Context, cb WalkGroupedCB) {
	hist := make(map[Cmd]bool)
	c.walkGroupedImpl(ctx, hist, c, nil, 0, 0, cb)
}

// AllGroupKeys collects group keys and returns them.
// Setting chooseFlag to true to grab the owned flags group
// names; setting sort to true to return a sorted slice.
func (c *CmdS) AllGroupKeys(chooseFlag, sort bool) []string {
	grpKeys := make([]string, 0)
	if chooseFlag {
		for gg := range c.allFlags {
			grpKeys = append(grpKeys, gg)
		}
	} else {
		for gg := range c.allCommands {
			grpKeys = append(grpKeys, gg)
		}
	}
	if sort {
		slices.Sort(grpKeys)
	}
	return grpKeys
}

// CommandsInGroup return all commands in a given group key.
func (c *CmdS) CommandsInGroup(groupTitle string) (list []Cmd) {
	if c.allCommands != nil {
		for _, a := range c.allCommands[groupTitle].A {
			list = append(list, a)
		}
	}
	return
}

// FlagsInGroup return all flags in a given group key.
func (c *CmdS) FlagsInGroup(groupTitle string) (list []*Flag) {
	if c.allFlags != nil {
		for _, a := range c.allFlags[groupTitle].A {
			list = append(list, a)
		}
	}
	return
}

func (c *CmdS) walkGroupedImpl(ctx context.Context, hist map[Cmd]bool, dad, grandpa Cmd, cmdIdx, level int, cb WalkGroupedCB) { //nolint:revive
	cb(dad, grandpa, nil, dad.GroupHelpTitle(), cmdIdx, level)

	// todo need ensure dynamic commands (and flags)

	grpKeys := dad.AllGroupKeys(false, true)
	for _, gg := range grpKeys {
		for i, cc := range dad.CommandsInGroup(gg) {
			if _, ok := hist[cc]; !ok {
				hist[cc] = true
				// cb(cc,nil, i, 0, level)
				c.walkGroupedImpl(ctx, hist, cc, dad, i, level+1, cb)
			} else {
				logz.WarnContext(ctx, "[cmdr] loop ref found", "dad", dad, "grandpa", grandpa, "cc", cc)
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
	// 			logz.WarnContext(ctx, "[cmdr] loop ref found", "dad", dad, "grandpa", grandpa, "cc", cc)
	// 		}
	// 	}
	// }

	grpKeys = dad.AllGroupKeys(true, true)
	for _, gg := range grpKeys {
		for i, ff := range dad.FlagsInGroup(gg) {
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
func (c *CmdS) WalkEverything(ctx context.Context, cb WalkEverythingCB) {
	hist := make(map[Cmd]bool)
	c.walkEx(ctx, hist, c, nil, 0, 0, cb)
}

func (c *CmdS) walkEx(ctx context.Context, hist map[Cmd]bool, dad, grandpa Cmd, level, cmdIndex int, cb WalkEverythingCB) { //nolint:revive
	cb(dad, grandpa, nil, cmdIndex, 0, level)

	// for i, cc := range dad.SubCommands() {
	// 	if cc != nil {
	// 		if _, ok := hist[cc]; !ok {
	// 			hist[cc] = true
	// 			// cb(cc,nil, i, 0, level)
	// 			c.walkEx(ctx, hist, cc, dad, level+1, i, cb)
	// 		} else {
	// 			logz.WarnContext(ctx, "[cmdr] loop ref found", "dad", dad, "grandpa", grandpa, "cc", cc)
	// 		}
	// 	}
	// }

	commands := mustEnsureDynCommands(ctx, dad)
	for i, cc := range commands {
		if cc != nil {
			if _, ok := hist[cc]; !ok {
				hist[cc] = true
				// cb(cc,nil, i, 0, level)
				c.walkEx(ctx, hist, cc, dad, level+1, i, cb)
			} else {
				logz.WarnContext(ctx, "[cmdr] loop ref found", "dad", dad, "grandpa", grandpa, "cc", cc)
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

func mustEnsureDynCommands(ctx context.Context, cmd Cmd) (commands []Cmd) {
	if cclist, err := ensureDynCommands(ctx, cmd); err != nil {
		logz.ErrorContext(ctx, "[cmdr] cannot evaluate dynamic-commands", "err", err)
		for _, cc := range cmd.OnEvalSubcommandsOnceCache() {
			commands = append(commands, cc)
		}
	} else {
		commands = append(commands, cclist...)
	}
	return
}

func mustEnsureDynFlags(ctx context.Context, cmd Cmd) (flags []*Flag) {
	flags = cmd.Flags()
	if cclist, err := ensureDynFlags(ctx, cmd); err != nil {
		logz.ErrorContext(ctx, "[cmdr] cannot evaluate dynamic-flags", "err", err)
	} else {
		flags = append(flags, cclist...)
	}
	return
}

func ensureDynCommands(ctx context.Context, cmd Cmd) (list []Cmd, err error) {
	var c Cmd

	for _, cc := range cmd.SubCommands() {
		list = append(list, cc)
	}

	if cb := cmd.OnEvalSubcommands(); cb != nil {
		logz.VerboseContext(ctx, "[cmdr] checking dynamic commands (always)", "cmd", cmd)

		var iter EvalIterator
		if iter, err = cb(ctx, cmd); err != nil || iter == nil {
			return
		}

		hasNext := true
		for hasNext {
			if c, hasNext, err = iter(ctx); err != nil {
				return
			}
			list = append(list, c)
		}
	}

	if cmd.OnEvalSubcommandsOnce() != nil {
		logz.InfoContext(ctx, "[cmdr] checking dynamic commands (once)", "cmd", cmd)
		var lst []Cmd
		if !cmd.OnEvalSubcommandsOnceInvoked() {
			var iter EvalIterator
			if iter, err = cmd.OnEvalSubcommandsOnce()(ctx, cmd); err != nil {
				return
			}

			hasNext := true
			for hasNext {
				if c, hasNext, err = iter(ctx); err != nil {
					return
				}
				lst = append(lst, c)
			}
			cmd.OnEvalSubcommandsOnceSetCache(lst)
		} else {
			lst = cmd.OnEvalSubcommandsOnceCache()
		}
		list = append(list, lst...)
	}

	return
}

func ensureDynFlags(ctx context.Context, cmd Cmd) (list []*Flag, err error) {
	var f *Flag

	if cb := cmd.OnEvalFlagsOnce(); cb != nil {
		logz.InfoContext(ctx, "[cmdr] checking dynamic flags (once)", "cmd", cmd)
		if !cmd.OnEvalFlagsOnceInvoked() {
			var iter EvalFlagIterator
			if iter, err = cb(ctx, cmd); err != nil {
				return
			}

			hasNext := true
			for hasNext {
				if f, hasNext, err = iter(ctx); err != nil {
					return
				}
				list = append(list, f)
			}
			cmd.OnEvalFlagsOnceSetCache(list)
			// cmd.onEvalFlagsOnce.invoked = true
			// cmd.onEvalFlagsOnce.flags = list
		} else {
			// list = cmd.onEvalFlagsOnce.flags
			list = cmd.OnEvalFlagsOnceCache()
		}
	}

	if cb := cmd.OnEvalFlags(); cb != nil {
		logz.InfoContext(ctx, "[cmdr] checking dynamic flags (always)", "cmd", cmd)
		var iter EvalFlagIterator
		if iter, err = cb(ctx, cmd); err != nil || iter == nil {
			return
		}

		hasNext := true
		for hasNext {
			if f, hasNext, err = iter(ctx); err != nil {
				return
			}
			list = append(list, f)
		}
	}
	return
}
