// Copyright © 2022 Hedzr Yeh.

package cli

import (
	"context"
	"strings"

	"github.com/hedzr/cmdr/v2/conf"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

func uniAddCmd(cmds []*CmdS, cmd *CmdS) []*CmdS { //nolint:unused
	for _, f := range cmds {
		if f == cmd || f.EqualTo(cmd) {
			return cmds
		}
	}
	return append(cmds, cmd)
}

func uniAddFlg(flags []*Flag, flg *Flag) []*Flag { //nolint:unused
	for _, f := range flags {
		if f == flg || f.EqualTo(flg) {
			return flags
		}
	}
	return append(flags, flg)
}

func uniAddStr(a []string, s string) []string {
	for _, f := range a {
		if f == s {
			return a
		}
	}
	return append(a, s)
}

func uniAddStrS(a []string, ss ...string) []string { //nolint:revive
	for _, s := range ss {
		found := false
		for _, f := range a {
			if f == s {
				found = true
				break
			}
		}
		if !found {
			a = append(a, s) //nolint:revive
		}
	}
	return a
}

type BacktraceableMin interface {
	GetDottedPath() string
	GetTitleName() string
	GetTitleNamesArray() []string
	GetTitleNames(maxWidth ...int) (title, rest string) // return the joint string of short,full,aliases names

	// OwnerOrParent return parent owner cmdS as a BacktraceableMin.
	//
	// A typical backtracing loop is:
	//
	//	var cc BacktraceableMin = cmd
	//	for cc.OwnerIsNotNil() {
	//	  // do sth with cc.(*CmdS)
	//	  cc = cc.OwnerOrParent()
	//	}
	OwnerOrParent() BacktraceableMin
	OwnerIsNil() bool
	OwnerIsNotNil() bool
	OwnerIsRoot() bool
	IsRoot() bool
	Root() *RootCommand
	App() App
}

type Backtraceable interface {
	BacktraceableMin

	Walk(ctx context.Context, cb WalkCB)
	WalkFast(ctx context.Context, cb WalkFastCB) (stop bool)
	ForeachFlags(context.Context, func(f *Flag) (stop bool)) (stop bool)
}

func DottedPath(cmd BacktraceableMin) string {
	return backtraceCmdNamesG(cmd, ".", false)
}

func DottedPathWith(cmd BacktraceableMin, delimiter string, verboseLast bool) string {
	return backtraceCmdNamesG(cmd, delimiter, verboseLast)
}

func backtraceCmdNamesG(cmd BacktraceableMin, delimiter string, verboseLast bool) (str string) { //nolint:revive
	var a []string
	if verboseLast {
		va := cmd.GetTitleNamesArray()
		if len(va) > 0 {
			vas := strings.Join(va, "|")
			a = append(a, "["+vas+"]")
		}
	} else {
		a = append(a, cmd.GetTitleName())
	}

	history := make(map[BacktraceableMin]struct{})
	p := cmd
	for p != nil && p.OwnerIsNotNil() {
		p = p.OwnerOrParent()
		if p != nil { // p has parent
			if _, ok := history[p]; !ok {
				history[p] = struct{}{}
				n := p.GetTitleName()
				a = append(a, n)
			} else {
				logz.Warn("[cmdr][backtraceCmdNamesG][ILLEGAL] the owners has cycle ref", "working-p", p.(Cmd).LongTitle(), "from-cmd", cmd)
			}
		}
	}

	// reverse it
	i := 0
	j := len(a) - 1
	for i < j {
		a[i], a[j] = a[j], a[i]
		i++
		j--
	}

	for i = 0; i < len(a); i++ {
		if a[i] != "" {
			break
		}
	}

	str = strings.Join(a[i:], delimiter)
	return
}

func dottedPathToCommandOrFlagG(c Backtraceable, dottedPath string) (cmd Backtraceable, ff *Flag) { //nolint:revive
	if c == nil {
		// c = internalGetWorker().rootCommand.Cmd
		return
	}

	appName := appNameDefault
	if app := c.App(); app != nil {
		appName = app.Name()
	} else if conf.AppName != "" {
		appName = conf.AppName
	}
	if appName != appNameDefault {
		appName = ""
	}

	if !strings.HasPrefix(dottedPath, appName) {
		dottedPath = appName + "." + dottedPath //nolint:revive
	}

	ctx := context.TODO()
	c.WalkFast(ctx, func(cc Cmd, index, level int) (stop bool) {
		var cx Backtraceable
		var ok bool
		if cx, ok = cc.(Backtraceable); !ok {
			return
		}

		kp := cc.GetDottedPath()
		if !strings.HasPrefix(kp, appName) {
			kp = appName + "." + kp
		}

		if stop = kp == dottedPath; stop {
			cmd = cx
			return
		}

		if strings.HasPrefix(dottedPath, kp) {
			parts := strings.TrimPrefix(dottedPath, kp+".")
			if !strings.Contains(parts, ".") {
				// try matching flags in this command
				stop = cx.ForeachFlags(ctx, func(f *Flag) (stop bool) {
					if parts == f.Long {
						cmd, ff, stop = f.OwnerOrParent(), f, true
						if f.OwnerOrParent() != cx {
							panic("flag's owner is not linked the proper parent.")
						}
					}
					return
				})
			}
		}
		return
	})
	return
}

func (c *CmdS) dottedPathToCommandOrFlag(dottedPath string) (cmd Backtraceable, ff *Flag) { //nolint:revive
	if c == nil {
		// anyCmd = &internalGetWorker().rootCommand.CmdS
		return
	}
	return dottedPathToCommandOrFlagG(c, dottedPath)
}

// DottedPathToCommandOrFlag1 searches the matched CmdS or Flag with the specified dotted-path.
// anyCmd is the starting of this searching.
func DottedPathToCommandOrFlag1(dottedPath string, anyCmd Backtraceable) (cc Backtraceable, ff *Flag) {
	return dottedPathToCommandOrFlagG(anyCmd, dottedPath)
}

// // backtraceCmdNames returns the sequences of a sub-command from
// // top-level.
// //
// // - if verboseLast = false, got 'microservices.tags.list' for sub-cmd microservice/tags/list.
// //
// // - if verboseLast = true,  got 'microservices.tags.[ls|list|l|lst|dir]'.
// //
// // - at root command, it returns 'appName' or ” when verboseLast is true.
// func backtraceCmdNames(cmd *BaseOpt, delimiter string, verboseLast bool) (str string) { //nolint:revive
// 	var a []string
// 	if verboseLast {
// 		va := cmd.GetTitleNamesArray()
// 		if len(va) > 0 {
// 			vas := strings.Join(va, "|")
// 			a = append(a, "["+vas+"]")
// 		}
// 	} else {
// 		a = append(a, cmd.GetTitleName())
// 	}
// 	for p := cmd.owner; p != nil && p.owner != nil; {
// 		a = append(a, p.GetTitleName())
// 		p = p.owner
// 	}
//
// 	// reverse it
// 	i := 0
// 	j := len(a) - 1
// 	for i < j {
// 		a[i], a[j] = a[j], a[i]
// 		i++
// 		j--
// 	}
//
// 	str = strings.Join(a, delimiter)
// 	return
// }
//
// // DottedPathToCommandOrFlag searches the matched CmdS or Flag with the specified dotted-path.
// // The searching will start from root if anyCmd is nil.
// func DottedPathToCommandOrFlag(dottedPath string, anyCmd *CmdS) (cc *CmdS, ff *Flag) {
// 	return anyCmd.dottedPathToCommandOrFlag(dottedPath)
// }
//
// // dottedPathToCommandOrFlag searches the matched CmdS or Flag with the specified dotted-path.
// // The searching will start from root if anyCmd is nil.
// func (c *CmdS) dottedPathToCommandOrFlag(dottedPath string) (cmd *CmdS, ff *Flag) { //nolint:revive
// 	if c == nil {
// 		// anyCmd = &internalGetWorker().rootCommand.CmdS
// 		return
// 	}
//
// 	appName := appNameDefault
// 	if app := c.App(); app != nil {
// 		appName = app.Name()
// 	} else if conf.AppName != "" {
// 		appName = conf.AppName
// 	}
//
// 	if !strings.HasPrefix(dottedPath, appName) {
// 		dottedPath = appName + "." + dottedPath //nolint:revive
// 	}
//
// 	ctx := context.TODO()
// 	c.Walk(ctx, func(cc Cmd, index, level int) {
// 		kp := cc.GetDottedPath()
// 		if !strings.HasPrefix(kp, appName) {
// 			kp = appName + "." + kp
// 		}
//
// 		if kp == dottedPath {
// 			cmd = cc
// 			return
// 		}
//
// 		if strings.HasPrefix(dottedPath, kp) {
// 			parts := strings.TrimPrefix(dottedPath, kp+".")
// 			if !strings.Contains(parts, ".") {
// 				// try matching flags in this command
// 				cc.ForeachFlags(func(f *Flag) (stop bool) {
// 					if parts == f.Long {
// 						cmd, ff, stop = f.owner, f, true
// 						if f.owner != cc {
// 							panic("flag's owner is not linked the proper parent.")
// 						}
// 					}
// 					return
// 				})
// 			}
// 		}
// 	})
// 	return
// }
