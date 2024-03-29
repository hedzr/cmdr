// Copyright © 2022 Hedzr Yeh.

package cli

import (
	"strings"

	"github.com/hedzr/cmdr/v2/conf"
)

func uniAddCmd(cmds []*Command, cmd *Command) []*Command { //nolint:unused
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

// backtraceCmdNames returns the sequences of a sub-command from
// top-level.
//
// - if verboseLast = false, got 'microservices.tags.list' for sub-cmd microservice/tags/list.
//
// - if verboseLast = true,  got 'microservices.tags.[ls|list|l|lst|dir]'.
//
// - at root command, it returns 'appName' or ” when verboseLast is true.
func backtraceCmdNames(cmd *BaseOpt, delimiter string, verboseLast bool) (str string) { //nolint:revive
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
	for p := cmd.owner; p != nil && p.owner != nil; {
		a = append(a, p.GetTitleName())
		p = p.owner
	}

	// reverse it
	i := 0
	j := len(a) - 1
	for i < j {
		a[i], a[j] = a[j], a[i]
		i++
		j--
	}

	str = strings.Join(a, delimiter)
	return
}

// DottedPathToCommandOrFlag searches the matched Command or Flag with the specified dotted-path.
// The searching will start from root if anyCmd is nil.
func DottedPathToCommandOrFlag(dottedPath string, anyCmd *Command) (cc *Command, ff *Flag) {
	return anyCmd.dottedPathToCommandOrFlag(dottedPath)
}

// dottedPathToCommandOrFlag searches the matched Command or Flag with the specified dotted-path.
// The searching will start from root if anyCmd is nil.
func (c *Command) dottedPathToCommandOrFlag(dottedPath string) (cmd *Command, ff *Flag) { //nolint:revive
	if c == nil {
		// anyCmd = &internalGetWorker().rootCommand.Command
		return
	}

	appName := appNameDefault
	if app := c.App(); app != nil {
		appName = app.Name()
	} else if conf.AppName != "" {
		appName = conf.AppName
	}

	if !strings.HasPrefix(dottedPath, appName) {
		dottedPath = appName + "." + dottedPath //nolint:revive
	}

	c.Walk(func(cc *Command, index, level int) {
		kp := cc.GetDottedPath()
		if !strings.HasPrefix(kp, appName) {
			kp = appName + "." + kp
		}

		if kp == dottedPath {
			cmd = cc
			return
		}

		if strings.HasPrefix(dottedPath, kp) {
			parts := strings.TrimPrefix(dottedPath, kp+".")
			if !strings.Contains(parts, ".") {
				// try matching flags in this command
				cc.ForeachFlags(func(f *Flag) (stop bool) {
					if parts == f.Long {
						cmd, ff, stop = f.owner, f, true
						if f.owner != cc {
							panic("flag's owner is not linked the proper parent.")
						}
					}
					return
				})
			}
		}
	})
	return
}
