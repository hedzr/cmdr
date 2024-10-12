package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/hedzr/is"
	"github.com/hedzr/is/states"
	logz "github.com/hedzr/logg/slog"

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

func (c *Command) SetOnEvaluateSubCommands(handler OnEvaluateSubCommands) {
	c.onEvalSubcommands = &struct{ cb OnEvaluateSubCommands }{cb: handler}
}

func (c *Command) SetOnEvaluateSubCommandsOnce(handler OnEvaluateSubCommands) {
	c.onEvalSubcommandsOnce = &struct {
		cb       OnEvaluateSubCommands
		invoked  bool
		commands []*Command
	}{cb: handler}
}

func (c *Command) SetOnEvaluateFlags(handler OnEvaluateFlags) {
	c.onEvalFlags = &struct{ cb OnEvaluateFlags }{cb: handler}
}

func (c *Command) SetOnEvaluateFlagsOnce(handler OnEvaluateFlags) {
	c.onEvalFlagsOnce = &struct {
		cb      OnEvaluateFlags
		invoked bool
		flags   []*Flag
	}{cb: handler}
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
		logz.Verbose("[cmdr] cmd.RunPreActions failed", "err", err)
		return
	}
	defer func() { deferActions(err) }() // err must be delayed caught here

	logz.Verbose("[cmdr] cmd.Invoke()", "onInvoke", c.onInvoke)
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

		// if !errors.Iss(errInvoked, ErrShouldStop, ErrShouldFallback) { // condition is true if errInvoked is nil
		// 	ecp.Attach(errInvoked) // no matter, attaching a nil error is no further effect
		// }

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
			logz.Panic("[cmdr] Error(s) occurred when running post-actions:", "error", ecp.Error())
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
// 			logz.Fatalf("[cmdr] Error(s) occurred when running post-actions: %v", ecp)
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
					logz.Warn("[cmdr] too much head-like flags", "last-head-like-flag", ff.owner.headLikeFlag, "this-one", ff)
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
					logz.Warn("[cmdr] too much head-like flags", "last-head-like-flag", ff.owner.headLikeFlag, "this-one", ff)
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

	logz.Debug("[cmdr] external editor", "ex-editor", ff.externalEditor)
	if cmd := os.Getenv(ff.externalEditor); cmd != "" {
		file := tool.TempFileName("message*.tmp", "message001.tmp", c.App().Name())
		cmdS := tool.SplitCommandString(cmd)
		cmdS = append(cmdS, file)
		defer func(dst string) {
			if err := dir.DeleteFile(dst); err != nil {
				logz.Error("[cmdr] cannot delete temporary file for flag", "flag", ff)
			}
		}(file)

		logz.Debug("[cmdr] invoke external editor", "ex-editor", ff.externalEditor, "cmd", cmdS)
		if is.DebuggerAttached() {
			vp.ValueOK, vp.Value = true, "<<stdoutTextForDebugging>>"
			logz.Warn("[cmdr] use debug text", "flag", ff, "text", vp.Value)
			return ff
		}

		if err := exec.CallSliceQuiet([]string{"which", cmdS[0]}, func(retCode int, stdoutText string) {
			if retCode == 0 {
				cmdS[0] = strings.TrimSpace(strings.TrimSuffix(stdoutText, "\n"))
				logz.Debug("[cmdr] got external editor real-path", "cmd", cmdS)
			}
		}); err != nil {
			logz.Error("[cmdr] cannot invoke which Command", "flag", ff, "cmd", cmdS)
			return nil
		}

		var content []byte
		var err error
		content, err = tool.LaunchEditorWithGetter(cmdS[0], func() string { return cmdS[1] }, false)
		if err != nil {
			logz.Error("[cmdr] Error on launching cmd", "err", err, "cmd", cmdS)
			return nil
		}

		// content, err = tool.LaunchEditorWith(cmdS[0], cmdS[1])
		// if err != nil {
		// 	logz.Error("[cmdr] Error on launching cmd", "err", err, "cmd", cmdS)
		// 	return nil
		// }
		//
		// content, err = tool.LaunchEditor(cmdS[0])
		// if err != nil {
		// 	logz.Error("[cmdr] Error on launching cmd", "err", err, "cmd", cmdS)
		// 	return nil
		// }

		// f, err = os.Open(file)
		// if err != nil {
		// 	logz.Error("[cmdr] cannot open temporary file for reading content", "file", file, "flag", ff, "cmd", cmdS)
		// 	return nil
		// }
		// defer f.Close()
		// vp.ValueOK, vp.Value = true, dir.MustReadAll(f)

		vp.ValueOK, vp.Value = true, string(content)
		ff.defaultValue = string(content)
		// logz.Debug("[cmdr] invoked external editor", "ex-editor", ff.externalEditor, "text", string(content))
		return ff
	}
	logz.Warn("[cmdr] Unknown External Editor for flag.", "ex-editor", ff.externalEditor, "flag", ff)
	return nil
}

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
