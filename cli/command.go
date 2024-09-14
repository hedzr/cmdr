package cli

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/hedzr/evendeep/ref"

	"github.com/hedzr/is"
	"github.com/hedzr/is/states"
	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/cmdr/v2/cli/atoa"
	"github.com/hedzr/cmdr/v2/internal/tool"
	"github.com/hedzr/cmdr/v2/pkg/dir"
	"github.com/hedzr/cmdr/v2/pkg/exec"

	"gopkg.in/hedzr/errors.v3"
)

//
//
//

func (c *Command) IsRoot() bool { return c.root.Command == c && c.owner == nil }

func (c *Command) HasFlag(longTitle string) (f *Flag, ok bool) {
	f, ok = c.longFlags[longTitle]
	return
}

// func (c *Command) Root() *RootCommand      { return c.root }
// func (c *Command) Owner() *Command         { return c.owner }

func (c *Command) SubCommands() []*Command { return c.commands }
func (c *Command) Flags() []*Flag          { return c.flags }
func (c *Command) String() string {
	var sb strings.Builder
	_, _ = sb.WriteString("Cmd{'")
	// _, _ = sb.WriteString(c.GetTitleName())
	_, _ = sb.WriteString(c.GetDottedPath())
	_, _ = sb.WriteString("'}")
	return sb.String()
}

// TailPlaceHolder is a string at end of usage line in help screen.
//
// In help screen, a command's usage line generally has the following form:
//
//	<app-name> <sub-commands> [<options>] [<positional-args>...]
//
// The text of <positional-args> is exact TailPlaceHolder. Set TailPlaceHolder
// to "files..." might be meaningful for "load" command, looks like:
//
//	<app> yaml-formatter load [<options>] files...
func (c *Command) TailPlaceHolder() string { return strings.Join(c.tailPlaceHolders, " ") }

// RedirectTo provides the real command target for current Command.
//
// Suppose command [app build] is being redirected to [app gcc build]. There
// [app build] is a shortcut to its full commands [app gcc build].
func (c *Command) RedirectTo() (dottedPath string) { return c.redirectTo }

// GetQuotedGroupName returns the group name quoted string.
func (c *Command) GetQuotedGroupName() string {
	if strings.TrimSpace(c.group) == "" {
		return ""
	}
	i := strings.Index(c.group, ".")
	if i >= 0 {
		return fmt.Sprintf("[%v]", c.group[i+1:])
	}
	return fmt.Sprintf("[%v]", c.group)
}

// GetExpandableNamesArray returns the names array of command, includes short name and long name.
func (c *Command) GetExpandableNamesArray() []string {
	var a []string
	if len(c.Long) > 0 {
		a = append(a, c.Long)
	}
	if len(c.Short) > 0 {
		a = append(a, c.Short)
	}
	return a
}

// GetExpandableNames returns the names comma splitted string.
func (c *Command) GetExpandableNames() string {
	a := c.GetExpandableNamesArray()
	if len(a) == 1 {
		return a[0]
	} else if len(a) > 1 {
		return fmt.Sprintf("{%v}", strings.Join(a, ","))
	}
	return c.name
}

//
//

func (c *Command) AppendTailPlaceHolder(placeHolder ...string) {
	c.tailPlaceHolders = append(c.tailPlaceHolders, placeHolder...)
}

func (c *Command) SetTailPlaceHolder(placeHolders ...string) { c.tailPlaceHolders = placeHolders }
func (c *Command) SetRedirectTo(dottedPath string)           { c.redirectTo = dottedPath }
func (c *Command) SetPresetCmdLines(args ...string)          { c.presetCmdLines = args }
func (c *Command) SetInvokeProc(str string)                  { c.invokeProc = str }
func (c *Command) SetInvokeShell(str string)                 { c.invokeShell = str }
func (c *Command) SetShell(str string)                       { c.shell = str }

//
//

func (c *Command) AddSubCommand(child *Command, callbacks ...func(cc *Command)) { //nolint:revive
	if child == nil {
		return
	}

	for _, cc := range c.commands {
		if cc == child || cc.EqualTo(child) {
			return
		}
	}
	for _, cb := range callbacks {
		if cb != nil {
			cb(child)
		}
	}
	c.commands = append(c.commands, child)
	child.owner = c
	child.root = c.root
}

func (c *Command) AddFlag(child *Flag, callbacks ...func(ff *Flag)) { //nolint:revive
	if child == nil {
		return
	}

	for _, cc := range c.flags {
		if cc == child || cc.EqualTo(child) {
			return
		}
	}
	for _, cb := range callbacks {
		if cb != nil {
			cb(child)
		}
	}
	c.flags = append(c.flags, child)
	child.owner = c
	child.root = c.root
}

//
//
//

// SetOnMatched adds the onMatched handler to a command
func (c *Command) SetOnMatched(functions ...OnCommandMatchedHandler) {
	c.onMatched = append(c.onMatched, functions...)
}

func (c *Command) CanInvoke() bool { return c.onInvoke != nil }

// SetPostActions adds the post-action to a command
func (c *Command) SetPostActions(functions ...OnPostInvokeHandler) {
	c.postActions = append(c.postActions, functions...)
}

// SetPreActions adds the pre-action to a command
func (c *Command) SetPreActions(functions ...OnPreInvokeHandler) {
	c.preActions = append(c.preActions, functions...)
}

// SetAction adds the onInvoke action to a command
func (c *Command) SetAction(fn OnInvokeHandler) {
	c.onInvoke = fn
}

func (c *Command) HasOnAction() bool { return c.onInvoke != nil }

func (c *Command) Invoke(args []string) (err error) {
	var deferActions func(errInvoked error)
	if deferActions, err = c.RunPreActions(c, args); err != nil {
		return
	}
	defer func() { deferActions(err) }() // err must be delayed caught here

	if c.onInvoke != nil {
		err = c.onInvoke(c, args)
	}
	return
}

func (c *Command) RunPreActions(cmd *Command, args []string) (deferAction func(errInvoked error), err error) { //nolint:revive
	ec := errors.New("[PRE-INVOKE]")
	defer ec.Defer(&err)
	if c.root.Command != c {
		for _, a := range c.root.preActions {
			if a != nil {
				ec.Attach(a(cmd, args))
			}
		}
	}

	for _, a := range c.preActions {
		if a != nil {
			ec.Attach(a(cmd, args))
		}
	}

	if !ec.IsEmpty() {
		deferAction = func(errInvoked error) {}
		return
	}

	deferAction = c.getDeferAction(cmd, args)
	return
}

func (c *Command) getDeferAction(cmd *Command, args []string) func(errInvoked error) { //nolint:revive
	return func(errInvoked error) {
		ecp := errors.New("[POST-INVOKE]")
		if !errors.Iss(errInvoked, ErrShouldStop, ErrShouldFallback) { // condition is true if errInvoked is nil
			ecp.Attach(errInvoked) // no matter, attaching a nil error is no further effect
		}

		for _, a := range c.postActions {
			if a != nil {
				ecp.Attach(a(cmd, args, errInvoked))
			}
		}
		if c.root.Command != c {
			for _, a := range c.root.postActions {
				if a != nil {
					ecp.Attach(a(cmd, args, errInvoked))
				}
			}
		}

		if !ecp.IsEmpty() {
			logz.Fatal("Error(s) occurred when running post-actions:", "error", ecp.Error())
		}
	}
}

// func (c *Command) RunPreActions(cmd *Command, args []string) (deferAction func(errInvoked error), err error) {
// 	var ec = errors.New("[PRE-INVOKE]")
// 	defer ec.Defer(&err)
// 	if c.preInvoke != nil {
// 		ec.Attach(c.preInvoke(cmd, args))
// 	}
//
// 	deferAction = func(errInvoked error) {
// 		var ecp = errors.New("[POST-INVOKE]")
// 		if c.postInvoke != nil {
// 			ecp.Attach(c.postInvoke(cmd, args, errInvoked))
// 		}
// 		if !ecp.IsEmpty() {
// 			logz.Fatalf("Error(s) occurred when running post-actions: %v", ecp)
// 		}
// 		return
// 	}
// 	return
// }

//
//
//

// EnsureTree associates owner and app between all subCommands and app/runner/rootCommand.
// EnsureTree links all commands as a tree (make root and owner linked).
func (c *Command) EnsureTree(app App, root *RootCommand) {
	root.app = app // link RootCommand.app to app
	root.name = app.Name()
	c.ensureTreeR(app, root)
}

// ensureTreeR link Command.owner to its parent, and Command.root to root.
// ensureTreeR links all commands as a tree (make root and owner linked).
func (c *Command) ensureTreeR(app App, root *RootCommand) { //nolint:unparam,revive
	c.WalkEverything(func(cc, pp *Command, ff *Flag, cmdIndex, flgIndex, level int) {
		cc.owner, cc.root, _ = pp, root, app
		if ff != nil {
			ff.owner, ff.root = cc, root
		}
	})
}

// EnsureXref builds the internal indexes and maps.
//
// Called by worker.Worker in preparing time (preProcess).
//
// ForeachSubCommands, ForeachFlags, ForeachGroupedSubCommands, and
// ForeachGroupedFlags needs EnsureXref called.
func (c *Command) EnsureXref(cb ...func(cc *Command, index, level int)) {
	c.Walk(func(cc *Command, index, level int) {
		cc.ensureXrefCommands()
		cc.ensureXrefFlags()
		cc.ensureXrefGroup()
		for _, fn := range cb {
			fn(cc, index, level)
		}
	})
}

func (c *Command) ensureXrefCommands() { //nolint:revive
	if c.longCommands == nil {
		c.longCommands = make(map[string]*Command)
		for _, cc := range c.commands {
			for _, ss := range cc.GetLongTitleNamesArray() {
				c.longCommands[ss] = cc
			}
		}
	}
	if c.shortCommands == nil {
		c.shortCommands = make(map[string]*Command)
		for _, cc := range c.commands {
			for _, ss := range cc.GetShortTitleNamesArray() {
				c.shortCommands[ss] = cc
			}
		}
	}

	// if c.allCommands == nil {
	// 	c.allCommands = make(map[string]map[string]*Command)
	// 	for _, cc := range c.commands {
	// 		if cc.Short != "" {
	// 			c.shortCommands[cc.Short] = cc
	// 		}
	// 	}
	// }
}

func (c *Command) ensureXrefFlags() { //nolint:revive
	if c.longFlags == nil {
		c.longFlags = make(map[string]*Flag)
		for _, ff := range c.flags {
			c.ensureToggleGroups(ff)
			ff.ensureXref()
			if ff.headLike {
				if ff.owner.headLikeFlag != nil && ff.owner.headLikeFlag != ff {
					logz.Warn("too much head-like flags", "last-head-like-flag", ff.owner.headLikeFlag, "this-one", ff)
				}
				ff.owner.headLikeFlag = ff
			}
			for _, ss := range ff.GetLongTitleNamesArray() {
				c.longFlags[ss] = ff
			}
		}
	}
	if c.shortFlags == nil {
		c.shortFlags = make(map[string]*Flag)
		for _, ff := range c.flags {
			c.ensureToggleGroups(ff)
			ff.ensureXref()
			if ff.headLike {
				if ff.owner.headLikeFlag != nil && ff.owner.headLikeFlag != ff {
					logz.Warn("too much head-like flags", "last-head-like-flag", ff.owner.headLikeFlag, "this-one", ff)
				}
				ff.owner.headLikeFlag = ff
			}
			for _, ss := range ff.GetShortTitleNamesArray() {
				c.shortFlags[ss] = ff
			}
		}
	}
}

func (c *Command) ensureToggleGroups(ff *Flag) {
	if tg := ff.ToggleGroup(); tg != "" {
		if c.toggles == nil {
			c.toggles = make(map[string]*ToggleGroupMatch)
		}
		if c.toggles[tg] == nil {
			c.toggles[tg] = &ToggleGroupMatch{Flags: make(map[string]*Flag)}
		}
		c.toggles[tg].Flags[ff.Title()] = ff
		if ff.group == "" {
			ff.group = tg
		}
	}
}

func (c *Command) ensureXrefGroup() { //nolint:revive
	if c.allCommands == nil {
		c.allCommands = make(map[string]*CmdSlice)
		for _, cc := range c.commands {
			cc.ensureXrefCommands()
			cc.ensureXrefFlags()
			group := cc.SafeGroup()
			if m, ok := c.allCommands[group]; ok {
				m.A = append(m.A, cc)
			} else {
				c.allCommands[group] = &CmdSlice{A: []*Command{cc}}
			}
		}
	}
	if c.allFlags == nil {
		c.allFlags = make(map[string]*FlgSlice)
		for _, cc := range c.flags {
			group := cc.SafeGroup()
			if m, ok := c.allFlags[group]; ok {
				m.A = append(m.A, cc)
			} else {
				c.allFlags[group] = &FlgSlice{A: []*Flag{cc}}
			}
		}
	}
}

//

//

func (c *Command) Match(title string) (short bool, cc *Command) {
	c.ensureXrefCommands()

	if title == "" {
		return
	}

	var ok bool
	if cc, ok = c.longCommands[title]; ok {
		cc.hitTitle = title
		cc.hitTimes++
		return
	}
	if cc, short = c.shortCommands[title]; short {
		cc.hitTitle = title
		cc.hitTimes++
		return
	}
	return
}

type FlagValuePkg struct {
	Args    []string
	AteArgs int

	SpecialTilde bool
	Short        bool

	Matched string
	Remains string

	PartialMatched bool
	Flags          []*Flag // matched flags, reserved.
	ValueOK        bool
	Value          any
}

// NewFVP gets a new FlagValuePkg done.
// A FlagValuePkg is a internal structure for tracing the flag's matching and parsing.
func NewFVP(args []string, remains string, short, plusSign, dblTilde bool) (vp *FlagValuePkg) {
	vp = &FlagValuePkg{
		Args:         args,
		Short:        short,
		SpecialTilde: dblTilde,
		Remains:      remains,
	}
	if plusSign {
		vp.Short, vp.ValueOK, vp.Value = true, true, true
	}
	return
}

func (s *FlagValuePkg) Reset() {
	s.Matched, s.ValueOK, s.Value, s.Flags, s.PartialMatched = "", false, nil, nil, false
}

func (c *Command) testDblTilde(dblTilde bool, ff *Flag) (matched bool) {
	matched = dblTilde || !ff.dblTildeOnly || (ff.dblTildeOnly && dblTilde)
	return
}

// MatchFlag try matching command title with vp.Remains, and update the relevant states.
//
// While a flag matched ok, returns vp.Matched != "" && ff != nil && err != nil
func (c *Command) MatchFlag(vp *FlagValuePkg) (ff *Flag, err error) { //nolint:revive
	c.ensureXrefFlags()

	var ok bool
	var matched, remains string
	if vp.Short { // short flag
		if ff, ok = c.shortFlags[vp.Remains]; ok && c.testDblTilde(vp.SpecialTilde, ff) {
			vp.PartialMatched, vp.Matched, vp.Remains, ff.hitTitle, ff.hitTimes = false, vp.Remains, "", vp.Remains, ff.hitTimes+1
			return c.tryParseValue(vp, ff)
		}

		// try for compact short flags
		matched, remains, ff, err = c.partialMatchFlag(vp.Remains, vp.Short, vp.SpecialTilde, c.shortFlags)
		if vp.PartialMatched = ff != nil && err == nil; vp.PartialMatched {
			vp.Matched, vp.Remains = matched, remains
			ff, err = c.tryParseValue(vp, ff)
		}

		// try to parse headLike flag
		if vp.Matched == "" && c.headLikeFlag != nil && ref.IsNumeric(c.headLikeFlag.defaultValue) {
			var num int64
			if num, err = strconv.ParseInt(vp.Remains, 0, 64); err == nil {
				vp.Matched, vp.Remains, ff = vp.Remains, "", c.headLikeFlag
				ff.defaultValue, vp.ValueOK = int(num), true // store the parsed value
				logz.Verbose("headLike flag matched", "flg", ff, "num", num)
			}
		}
	} else {
		if ff, ok = c.longFlags[vp.Remains]; ok && c.testDblTilde(vp.SpecialTilde, ff) {
			vp.PartialMatched, vp.Matched, vp.Remains, ff.hitTitle, ff.hitTimes = false, vp.Remains, "", vp.Remains, ff.hitTimes+1
			return c.tryParseValue(vp, ff)
		}
		matched, remains, ff, err = c.partialMatchFlag(vp.Remains, vp.Short, vp.SpecialTilde, c.longFlags)
		if vp.PartialMatched = ff != nil && err == nil; vp.PartialMatched {
			vp.Matched, vp.Remains = matched, remains
			ff, err = c.tryParseValue(vp, ff)
		}
	}

	// lookup the parents, if 'ff' not matched/found
	if ff == nil && err == nil && c.owner != nil && c.owner != c {
		ff, err = c.owner.MatchFlag(vp)
		return
	}

	// when a flag matched ok
	if ff != nil && err == nil && vp.Matched != "" {
		ff.hitTitle = vp.Matched
		ff.hitTimes++
		if !vp.ValueOK {
			ff, err = c.tryParseValue(vp, ff)
			// // tryParseValue ...
			// if vp.PartialMatched {
			// 	//
			// } else {
			// 	//
			// }
		}
	}
	return
}

func (c *Command) partialMatchFlag(title string, short, dblTildeMode bool, mFlags map[string]*Flag) (matched, remains string, ff *Flag, err error) { //nolint:revive
	var maxLen int
	var rightPart string

	titleOriginal := title
	if pos := strings.IndexRune(title, '='); pos >= 0 {
		rightPart = title[pos+1:]
		title = title[:pos] //nolint:revive
	}

	for k, v := range mFlags {
		if strings.HasPrefix(title, k) {
			if maxLen < len(k) {
				if c.testDblTilde(dblTildeMode, v) {
					// keep the longest matched flag here
					maxLen, matched, remains, ff = len(k), k, title[len(k):], v
					if remains == "" && rightPart != "" {
						remains = rightPart
					}
				}
			}
		}
	}

	if maxLen > 0 {
		// if any flag matched, checking the parents for looking up the longer ones
		if c.owner != nil && c.owner != c {
			c.owner.ensureXrefFlags()
			mf := c.owner.longFlags
			if short {
				mf = c.owner.shortFlags
			}
			matched1, remains1, ff1, err1 := c.owner.partialMatchFlag(titleOriginal, short, dblTildeMode, mf)
			if err = err1; err != nil {
				return
			}
			if ff1 != nil && maxLen < len(matched1) {
				// if longer matched flag from parents exists, use it instead of the lastCommand's
				matched, remains, ff = matched1, remains1, ff1
			}
		}
		return
	}

	if c.owner != nil && c.owner != c {
		// if no flag matched, checking the parents
		c.owner.ensureXrefFlags()
		mf := c.owner.longFlags
		if short {
			mf = c.owner.shortFlags
		}
		matched, remains, ff, err = c.owner.partialMatchFlag(titleOriginal, short, dblTildeMode, mf)
	}
	return
}

func (c *Command) tryParseValue(vp *FlagValuePkg, ff *Flag) (ret *Flag, err error) {
	if ff != nil {
		ff = c.matchedForTG(ff) //nolint:revive
	}
	if ff, err = c.checkPrerequisites(vp, ff); err != nil {
		return
	}
	if ff, err = c.checkJustOnce(vp, ff); err != nil {
		return
	}

	if ff != nil && !vp.ValueOK {
		// try to parse value
		switch ff.defaultValue.(type) {
		case bool:
			ff = c.tryParseBoolValue(vp, ff) //nolint:revive
		case string:
			ff = c.tryParseStringValue(vp, ff) //nolint:revive
		case nil:
			ff = c.tryParseStringValue(vp, ff) //nolint:revive
		default:
			ff = c.tryParseOthersValue(vp, ff) //nolint:revive
		}
	}

	ret = ff

	if _, err = c.checkCircuitBreak(vp, ff); err != nil {
		return
	}
	return
}

func (c *Command) matchedForTG(ff *Flag) *Flag {
	// toggle group
	if ff.owner.toggles != nil {
		if m, ok := ff.owner.toggles[ff.ToggleGroup()]; ok {
			if f, ok := m.Flags[ff.Name()]; ok {
				for _, v := range m.Flags {
					v.SetDefaultValue(false)
				}
				f.SetDefaultValue(true)
				m.Matched, m.MatchedTitle = f, f.Name()
			}
		}
	}
	// mutual exclusives
	if len(ff.mutualExclusives) > 0 {
		for _, fn := range ff.mutualExclusives {
			var f *Flag
			if strings.ContainsRune(fn, '.') {
				f = ff.owner.FindFlag(fn, false)
			} else {
				_, f = ff.Root().dottedPathToCommandOrFlag(fn)
			}
			if f != nil {
				if _, ok := f.defaultValue.(bool); ok {
					f.SetDefaultValue(false)
				}
			}
		}
	}
	return ff
}

func (c *Command) checkJustOnce(vp *FlagValuePkg, ff *Flag) (ret *Flag, err error) {
	if ff != nil && ff.justOnce {
		if ff.hitTimes > 1 {
			err = ErrFlagJustOnce.FormatWith(ff)
			return
		}
	}
	ret, _ = ff, vp
	return
}

func (c *Command) checkPrerequisites(vp *FlagValuePkg, ff *Flag) (ret *Flag, err error) {
	if ff != nil && len(ff.prerequisites) > 0 {
		for _, fn := range ff.prerequisites {
			var f *Flag
			if strings.ContainsRune(fn, '.') {
				f = ff.owner.FindFlag(fn, false)
			} else {
				_, f = ff.Root().dottedPathToCommandOrFlag(fn)
			}
			if f != nil {
				if f.hitTimes < 0 {
					err = ErrMissedPrerequisite.FormatWith(ff, f)
					return
				}
			}
		}
	}
	ret, _ = ff, vp
	return
}

func (c *Command) checkCircuitBreak(vp *FlagValuePkg, ff *Flag) (ret *Flag, err error) {
	if ff != nil && ff.circuitBreak {
		err = ErrShouldStop
		return
	}
	ret, _ = ff, vp
	return
}

var (
	ErrMissedPrerequisite = errors.New("Flag %q needs %q was set at first") // flag need a prerequisite flag exists.
	ErrFlagJustOnce       = errors.New("Flag %q MUST BE set once only")     // flag cannot be set more than one time.
)

func (c *Command) tryParseStringValue(vp *FlagValuePkg, ff *Flag) *Flag {
	if ff.externalEditor != "" {
		if f := c.invokeExternalEditor(vp, ff); f != nil {
			return f
		}
	}

	if vp.Remains != "" {
		vp.ValueOK, vp.Value, vp.Remains = true, c.normalizeStringValue(vp.Remains), ""
	} else if vp.AteArgs < len(vp.Args) {
		vp.ValueOK, vp.Value, vp.AteArgs = true, c.normalizeStringValue(vp.Args[vp.AteArgs]), vp.AteArgs+1
	} else {
		vp.ValueOK, vp.Value = true, ""
	}
	ff.defaultValue = vp.Value
	return ff
}

func (c *Command) tryParseBoolValue(vp *FlagValuePkg, ff *Flag) *Flag {
	if len(vp.Remains) > 0 {
		switch ch := vp.Remains[0]; ch {
		case '+':
			vp.Value, vp.ValueOK = true, true
			vp.Remains = vp.Remains[1:]
			ff.defaultValue = vp.Value
		case '-':
			vp.Value, vp.ValueOK = false, true
			vp.Remains = vp.Remains[1:]
			ff.defaultValue = vp.Value
		default:
			vp.Value, vp.ValueOK = true, true
			ff.defaultValue = vp.Value
		}
	} else if !vp.ValueOK {
		vp.Value, vp.ValueOK = true, true
		ff.defaultValue = vp.Value
	} else {
		ff.defaultValue = vp.Value
	}
	return ff
}

func (c *Command) tryParseOthersValue(vp *FlagValuePkg, ff *Flag) *Flag {
	if vp.Remains != "" {
		vp.ValueOK, vp.Value, vp.Remains = true, c.fromString(vp.Remains, ff.defaultValue), ""
	} else {
		vp.ValueOK, vp.Value, vp.AteArgs = true, c.fromString(vp.Args[vp.AteArgs], ff.defaultValue), vp.AteArgs+1
	}
	if ref.IsSlice(vp.Value) {
		if ff.hitTimes == 0 {
			ff.defaultValue = vp.Value
		} else {
			ff.defaultValue = ref.SliceMerge(ff.defaultValue, vp.Value)
		}
	} else {
		ff.defaultValue = vp.Value
	}
	return ff
}

func (c *Command) fromString(text string, meme any) (value any) { //nolint:revive
	var err error
	value, err = atoa.Parse(text, meme)
	if err != nil {
		logz.ErrorContext(context.TODO(), "cannot parse text to value", "err", err, "text", text, "target-value-meme", meme)
	}
	return
}

func (c *Command) normalizeStringValue(sv string) string {
	return tool.StripQuotes(sv)
}

//

func (c *Command) invokeExternalEditor(vp *FlagValuePkg, ff *Flag) *Flag {
	if vp.Remains != "" {
		arg := c.normalizeStringValue(vp.Remains)
		vp.ValueOK, vp.Value, vp.Remains = true, arg, ""
		ff.defaultValue = arg
		return ff
	}
	if vp.AteArgs < len(vp.Args) {
		arg := c.normalizeStringValue(vp.Args[vp.AteArgs])
		if !strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "~") {
			vp.ValueOK, vp.Value, vp.AteArgs = true, arg, vp.AteArgs+1
			ff.defaultValue = arg
			return ff
		}
	}

	logz.Debug("external editor", "ex-editor", ff.externalEditor)
	if cmd := os.Getenv(ff.externalEditor); cmd != "" {
		f, err := os.CreateTemp(os.TempDir(), "message*.tmp")
		if err != nil {
			logz.Error("cannot create temporary file for flag", "flag", ff)
			return nil
		}
		file := f.Name()
		f.Close()
		cmdS := tool.SplitCommandString(cmd)
		cmdS = append(cmdS, file)
		defer func(dst string) {
			err = dir.DeleteFile(dst)
			if err != nil {
				logz.Error("cannot delete temporary file for flag", "flag", ff)
			}
		}(file)

		logz.Debug("invoke external editor", "ex-editor", ff.externalEditor, "cmd", cmdS)
		if is.DebuggerAttached() {
			vp.ValueOK, vp.Value = true, "<<stdoutTextForDebugging>>"
			logz.Warn("use debug text", "flag", ff, "text", vp.Value)
			return ff
		}

		err = exec.CallSliceQuiet([]string{"which", cmdS[0]}, func(retCode int, stdoutText string) {
			if retCode == 0 {
				cmdS[0] = strings.TrimSpace(strings.TrimSuffix(stdoutText, "\n"))
				logz.Debug("got external editor real-path", "cmd", cmdS)
			}
		})

		if err != nil {
			logz.Error("cannot invoke which Command", "flag", ff, "cmd", cmdS)
			return nil
		}

		var content []byte
		content, err = tool.LaunchEditor(cmdS[0], func() string { return cmdS[1] })
		if err != nil {
			logz.Error("Error on launching cmd", "err", err, "cmd", cmdS)
			return nil
		}

		// f, err = os.Open(file)
		// if err != nil {
		// 	logz.Error("cannot open temporary file for reading content", "file", file, "flag", ff, "cmd", cmdS)
		// 	return nil
		// }
		// defer f.Close()
		// vp.ValueOK, vp.Value = true, dir.MustReadAll(f)

		vp.ValueOK, vp.Value = true, string(content)
		ff.defaultValue = string(content)
		// logz.Debug("invoked external editor", "ex-editor", ff.externalEditor, "text", string(content))
		return ff
	}
	logz.Warn("Unknown External Editor for flag.", "ex-editor", ff.externalEditor, "flag", ff)
	return nil
}

//

func (c *Command) TryOnMatched(position int, hitState *MatchState) (handled bool, err error) {
	if c.onMatched != nil {
		handled = true
		for _, m := range c.onMatched {
			err = m(c, position, hitState)
			if !c.errIsSignalFallback(err) {
				err, handled = nil, false
			}
		}
	}
	return
}

//

// EqualTo compares with another one based on its titles
func (c *Command) EqualTo(rh *Command) (ok bool) {
	if c == nil {
		return rh == nil
	}
	if rh == nil {
		return false
	}
	return c.GetTitleName() == rh.GetTitleName()
}

// MatchTitleNameFast matches a given title string without indices built.
func (c *Command) MatchTitleNameFast(title string) (ok bool) { //nolint:revive
	if title == "" {
		return
	}

	ok = c.Long == title || c.Short == title
	if !ok {
		for _, t := range c.Aliases {
			if ok = t == title; ok {
				break
			}
		}
	}
	return
}

// FindSubCommand find sub-command with `longName` from `cmd`.
//
// If wide is true, FindSubCommand try to match  long, short and aliases titles,
// If wide is false, only long title matched.
func (c *Command) FindSubCommand(longName string, wide bool) (res *Command) {
	// return FindSubCommand(longName, c)
	for _, cx := range c.commands {
		if wide {
			for k, v := range c.longCommands {
				if k == longName {
					res = v
					return
				}
			}
			for k, v := range c.shortCommands {
				if k == longName {
					res = v
					return
				}
			}
		} else {
			if longName == cx.Long {
				res = cx
				return
			}
		}
	}
	return
}

// FindSubCommandRecursive find sub-command with `longName` from `cmd` recursively
func (c *Command) FindSubCommandRecursive(longName string, wide bool) (res *Command) { //nolint:revive
	// return FindSubCommandRecursive(longName, c)
	for _, cx := range c.commands {
		if wide {
			for k, v := range c.longCommands {
				if k == longName {
					res = v
					return
				}
			}
			for k, v := range c.shortCommands {
				if k == longName {
					res = v
					return
				}
			}
		} else {
			if longName == cx.Long {
				res = cx
				return
			}
		}
	}

	for _, cx := range c.commands {
		if len(cx.commands) > 0 {
			if res = cx.FindSubCommandRecursive(longName, wide); res != nil {
				return
			}
		}
	}
	return
}

// FindFlag find flag with `longName` from `cmd`
func (c *Command) FindFlag(longName string, wide bool) (res *Flag) {
	// return FindFlag(longName, c)
	for _, cx := range c.flags {
		if wide {
			for k, v := range c.longFlags {
				if k == longName {
					res = v
					return
				}
			}
			for k, v := range c.shortFlags {
				if k == longName {
					res = v
					return
				}
			}
		} else {
			if longName == cx.Long {
				res = cx
				return
			}
		}
	}
	return
}

// FindFlagRecursive find flag with `longName` from `cmd` recursively
func (c *Command) FindFlagRecursive(longName string, wide bool) (res *Flag) {
	// return FindFlagRecursive(longName, c)
	for _, cx := range c.flags {
		if wide {
			for k, v := range c.longFlags {
				if k == longName {
					res = v
					return
				}
			}
			for k, v := range c.shortFlags {
				if k == longName {
					res = v
					return
				}
			}
		} else {
			if longName == cx.Long {
				res = cx
				return
			}
		}
	}

	for _, cx := range c.commands {
		// if len(cx.SubCommands) > 0 {
		if res = cx.FindFlagRecursive(longName, false); res != nil {
			return
		}
		// }
	}
	return
}

func (c *Command) FindFlagBackwards(longName string) (res *Flag) {
	for _, cx := range c.flags {
		if longName == cx.Long {
			res = cx
			return
		}
	}
	if c.owner != nil && c.owner != c {
		res = c.owner.FindFlagBackwards(longName)
	}
	return
}

//
//
//

func (c *Command) GetGroupedCommands(group string) (commands []*Command) {
	c.ensureXrefGroup()
	commands = c.allCommands[group].A
	return
}

func (c *Command) GetGroupedFlags(group string) (flags []*Flag) {
	c.ensureXrefGroup()
	flags = c.allFlags[group].A
	return
}

// ForeachSubCommands is another way to Walk on all commands.
func (c *Command) ForeachSubCommands(cb func(cc *Command) (stop bool)) (stop bool) {
	for _, cc := range c.commands {
		if cc != nil && cb != nil {
			if stop = cb(cc); stop {
				break
			}
		}
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

func (c *Command) CountOfCommands() int {
	vc := states.Env().CountOfVerbose()
	cnt := 0
	for _, cc := range c.commands {
		if cc.vendorHidden {
			if vc > 2 {
				cnt++
			}
		} else if cc.hidden {
			if vc > 0 {
				cnt++
			}
		} else {
			cnt++
		}
	}
	return cnt
}

func (c *Command) CountOfFlags() int {
	vc := states.Env().CountOfVerbose()
	cnt := 0
	for _, cc := range c.flags {
		// if (vc > 0 && cc.Hidden()) || (vc > 2 && cc.VendorHidden()) || (!cc.hidden && !cc.vendorHidden) {
		// 	cnt++
		// }
		if cc.vendorHidden {
			if vc > 2 {
				cnt++
			}
		} else if cc.hidden {
			if vc > 0 {
				cnt++
			}
		} else {
			cnt++
		}
	}
	return cnt
}

// WalkBackwards is a simple way to loop for all commands.
func (c *Command) WalkBackwards(cb WalkBackwardsCB) {
	ctx := &WalkBackwardsCtx{
		Grouped: true,
		hist:    make(map[*Command]bool),
		histff:  make(map[*Flag]bool),
	}
	c.walkBackwardsImpl(ctx, c, 0, cb)
}
func (c *Command) WalkBackwardsCtx(cb WalkBackwardsCB, ctx *WalkBackwardsCtx) {
	if ctx.hist == nil {
		ctx.hist = make(map[*Command]bool)
	}
	if ctx.histff == nil {
		ctx.histff = make(map[*Flag]bool)
	}
	c.walkBackwardsImpl(ctx, c, 0, cb)
}

type WalkBackwardsCtx struct {
	Grouped bool
	hist    map[*Command]bool
	histff  map[*Flag]bool
}

// WalkBackwardsCB is a callback functor used by WalkBackwards.
//
// cc and ff is the looping command and/or one of its flags. ff == nil means
// that the looping item is a command in this turn.
//
// index and groupIndex is 0-based.
//
// count is count of items in a command/flag group.
//
// level is how many times the nested flags or commands backwards to
// the root command.
type WalkBackwardsCB func(cc *Command, ff *Flag, index, groupIndex, count, level int)

func (c *Command) walkBackwardsImpl(ctx *WalkBackwardsCtx, cmd *Command, level int, cb WalkBackwardsCB) {
	if cb != nil {
		if level == 0 {
			count := len(cmd.commands)
			lastG, gi := "", 0
			for i, cc := range cmd.commands {
				if _, ok := ctx.hist[cc]; !ok {
					ctx.hist[cc] = true
					if ctx.Grouped {
						if g := cc.GroupHelpTitle(); g != lastG {
							lastG, gi = g, 0
						}
					}
					cb(cc, nil, i, gi, count, level)
					gi++
				}
			}
		}

		count := len(cmd.flags)
		lastG, gi := "", 0
		for i, ff := range cmd.flags {
			if _, ok := ctx.histff[ff]; !ok {
				ctx.histff[ff] = true
				if ctx.Grouped {
					if g := ff.GroupHelpTitle(); g != lastG {
						lastG, gi = g, 0
					}
				}
				cb(cmd, ff, i, gi, count, level)
				gi++
			}
		}
	}

	if cmd.owner != nil && cmd.owner != cmd {
		c.walkBackwardsImpl(ctx, cmd.owner, level+1, cb)
	}
}

// Walk is a simple way to loop for all commands.
func (c *Command) Walk(cb WalkCB) {
	hist := make(map[*Command]bool)
	c.walkImpl(hist, c, 0, cb)
}

type WalkCB func(cc *Command, index, level int)

func (c *Command) walkImpl(hist map[*Command]bool, cmd *Command, level int, cb WalkCB) {
	if cb != nil {
		cb(cmd, 0, level)
	}

	for _, cc := range cmd.commands {
		if cc != nil {
			if _, ok := hist[cc]; !ok {
				hist[cc] = true
				c.walkImpl(hist, cc, level+1, cb)
			} else {
				logz.Warn("loop ref found", "dad", cmd, "cc", cc)
			}
		}
	}
}

// WalkGrouped loops for all commands and its flags with grouped order.
func (c *Command) WalkGrouped(cb WalkGroupedCB) {
	hist := make(map[*Command]bool)
	c.walkGroupedImpl(hist, c, nil, 0, 0, cb)
}

type WalkGroupedCB func(cc, pp *Command, ff *Flag, group string, idx, level int)

func (c *Command) walkGroupedImpl(hist map[*Command]bool, dad, grandpa *Command, cmdIdx, level int, cb WalkGroupedCB) { //nolint:revive
	cb(dad, grandpa, nil, dad.GroupHelpTitle(), cmdIdx, level)

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
				logz.Warn("loop ref found", "dad", dad, "grandpa", grandpa, "cc", cc)
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
	// 			logz.Warn("loop ref found", "dad", dad, "grandpa", grandpa, "cc", cc)
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

type WalkEverythingCB func(cc, pp *Command, ff *Flag, cmdIndex, flgIndex, level int)

func (c *Command) walkEx(hist map[*Command]bool, dad, grandpa *Command, level, cmdIndex int, cb WalkEverythingCB) { //nolint:revive
	cb(dad, grandpa, nil, cmdIndex, 0, level)

	for i, cc := range dad.commands {
		if cc != nil {
			if _, ok := hist[cc]; !ok {
				hist[cc] = true
				// cb(cc,nil, i, 0, level)
				c.walkEx(hist, cc, dad, level+1, i, cb)
			} else {
				logz.Warn("loop ref found", "dad", dad, "grandpa", grandpa, "cc", cc)
			}
		}
	}

	for i, ff := range dad.flags {
		if ff != nil {
			cb(dad, grandpa, ff, cmdIndex, i, level)
		}
	}
}
