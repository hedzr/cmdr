package cmdr

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/hedzr/errors.v3"
)

// helpSystemAction is __complete command handler.
func (w *ExecWorker) helpSystemAction(cmdComplete *Command, args []string) (err error) {
	var (
		ctx = &queryShcompContext{
			w:           w,
			cmdComplete: cmdComplete,
			args:        args,

			matchedPrecededList:     make(map[string]*Command),
			matchedList:             make(map[string]*Command),
			matchedPrecededFlagList: make(map[string]*Flag),
			matchedFlagList:         make(map[string]*Flag),
			directive:               shellCompDirectiveNoFileComp,
		}
		x     strings.Builder
		keys  []string
		count int
	)

	err = ctx.lookupForHelpSystem(cmdComplete, args)
	if IsIgnorableError(err) {
		if hit := cmdComplete.GetHitStr(); hit == "help" || hit == "h" {
			if ctx.matchedCmd != nil {
				w.printHelp(ctx.matchedCmd, false)
				return
			}
		}

		defer func() {
			x.WriteString(fmt.Sprintf(":%d", ctx.directive))
			fp("%v", x.String())
			_, _ = fmt.Fprintf(os.Stderr, `%v
%v Items populated.
Args: %v
`, directivesToString(ctx.directive), count, args)
		}()

		cptLocal := getCPT()

		if ctx.matchedFlag != nil {
			keys = getSortedKeysFromFlgMap(ctx.matchedPrecededFlagList)
			for _, k := range keys {
				c := ctx.matchedPrecededFlagList[k]
				x.WriteString(fmt.Sprintf("%v\t%v\n", k, cptLocal.Translate(c.Description, BgNormal)))
				count++
			}

			keys = getSortedKeysFromFlgMap(ctx.matchedFlagList)
			for _, k := range keys {
				c := ctx.matchedFlagList[k]
				x.WriteString(fmt.Sprintf("%v\t%v\n", k, cptLocal.Translate(c.Description, BgNormal)))
				count++
			}
			return
		}

		keys = getSortedKeysFromCmdMap(ctx.matchedPrecededList)
		for _, k := range keys {
			c := ctx.matchedPrecededList[k]
			x.WriteString(fmt.Sprintf("%v\t%v\n", k, cptLocal.Translate(c.Description, BgNormal)))
			count++
		}

		keys = getSortedKeysFromCmdMap(ctx.matchedList)
		for _, k := range keys {
			c := ctx.matchedList[k]
			x.WriteString(fmt.Sprintf("%v\t%v\n", k, cptLocal.Translate(c.Description, BgNormal)))
			count++
		}
	}
	return
}

type queryShcompContext struct { //nolint:govet //meaningful order
	w           *ExecWorker
	cmdComplete *Command
	args        []string

	// returned

	matchedCmd  *Command
	matchedFlag *Flag
	exact       bool

	// dynamic updating

	matchedPrecededList, matchedList         map[string]*Command
	matchedPrecededFlagList, matchedFlagList map[string]*Flag
	directive                                int
}

func (ctx *queryShcompContext) buildMatchedSubCommandList(cmd *Command) {
	ctx.matchedCmd = cmd
	for _, c := range cmd.SubCommands {
		ctx.matchedList[c.GetTitleName()] = c
	}
}

func (ctx *queryShcompContext) winspec(s string) string {
	if s == "`\"`\"" {
		return ""
	}
	return s
}

func (ctx *queryShcompContext) lookupForHelpSystem(cmdComplete *Command, args []string) (err error) {
	var list []*Command
	cmd := &cmdComplete.root.Command
	var flg *Flag
	var exact bool
	defer func() {
		ctx.matchedCmd, ctx.matchedFlag, ctx.exact = cmd, flg, exact
	}()

	for i, ic := 0, len(args); i < ic; i++ {
		ttl := ctx.winspec(args[i])
		l := len(ttl)

		if l == 0 {
			if i == 0 {
				// root command matched
				ctx.buildMatchedSubCommandList(cmd)
				return
			} else if i == ic-1 {
				// for commandline: app __complete generate ''
				ll := len(list)
				if ll >= 1 {
					for i1 := ll - 1; i1 >= 0; i1-- {
						ctx.deleteFromMatchedList(list[i1])
					}
				}
				ctx.rebuildMatchedList(list[ll-1])
			}
			continue
		}

		if l > 0 && strings.Contains(ctx.w.switchCharset, ttl[0:1]) {
			// flags
			flg, err = ctx.lookupFlagsForHelpSystem(ttl, cmd, args, i)
			if !IsIgnorableError(err) {
				break
			}
			continue
		}

		// sub-commands
		cmd, exact, err = ctx.lookupCommandsForHelpSystem(ttl, cmd, args, i, ic)
		if !IsIgnorableError(err) {
			break
		}

		list = append(list, cmd)
	}
	return
}

func (ctx *queryShcompContext) lookupCommandsForHelpSystem(title string, parent *Command, args []string, ix, ic int) (cmdMatched *Command, exact bool, err error) {
	ctx.directive = shellCompDirectiveNoFileComp
	err = errors.NotFound
	ok := false

	for _, c := range parent.SubCommands {
		if c.VendorHidden {
			continue
		}

		if exact, ok = ctx.matchCommandTitle(c, title, ix == ic-1); ok {
			if /*ix == len(args)-1-1 && args[ix+1] == "" &&*/ exact {
				cmdMatched, err = c, nil
			} else {
				err = nil
			}
		}
	}
	return
}

func (ctx *queryShcompContext) rebuildMatchedList(cmd *Command) {
	ctx.deleteFromMatchedList(cmd)
	for _, c := range cmd.SubCommands {
		ctx.matchedList[c.GetTitleName()] = c
	}
}

func (ctx *queryShcompContext) deleteFromMatchedList(cmd *Command) {
	// if _, ok := ctx.matchedList[cmd.GetTitleName()]; ok {
	delete(ctx.matchedList, cmd.GetTitleName())
	// }
}

func (ctx *queryShcompContext) matchCommandTitle(c *Command, titleChecking string, fuzzy bool) (exact, ok bool) {
	ok = ctx.visitCommandTitles(c, false, func(c *Command, title string) (stopNow bool) {
		if exact, ok = ctx.doMatchCommandTitle(c, title, titleChecking, fuzzy); ok {
			stopNow = true
		}
		return
	})
	return
}

func (ctx *queryShcompContext) visitCommandTitles(c *Command, justFullTitle bool, fn func(c *Command, title string) (stopNow bool)) (ok bool) {
	if c.Full != "" && fn(c, c.Full) {
		return true
	}
	if justFullTitle {
		return
	}
	if c.Short != "" && fn(c, c.Short) {
		return true
	}
	for _, t := range c.Aliases {
		if t != "" && fn(c, t) {
			return true
		}
	}
	return
}

func (ctx *queryShcompContext) doMatchCommandTitle(c *Command, title, titleChecking string, fuzzy bool) (exact, ok bool) {
	if title == titleChecking {
		ctx.matchedList[c.GetTitleName()] = c
		exact, ok = true, true
		return
	}
	if fuzzy && strings.HasPrefix(title, titleChecking) {
		ctx.matchedPrecededList[c.GetTitleName()] = c
		ok = true
	} else if !noPartialMatching && fuzzy && strings.Contains(title, titleChecking) {
		ctx.matchedPrecededList[c.GetTitleName()] = c
		ok = true
	}
	return
}

func (ctx *queryShcompContext) lookupFlagsForHelpSystem(titleChecking string, parent *Command, args []string, ix int) (flgMatched *Flag, err error) {
	ctx.directive = shellCompDirectiveNoFileComp
	err = errors.NotFound

	sw1 := len(titleChecking) > 0 && strings.ContainsAny(titleChecking[0:1], ctx.w.switchCharset)
	sw2 := len(titleChecking) > 1 && strings.ContainsAny(titleChecking[1:2], ctx.w.switchCharset)

goUp:
	for _, c := range parent.Flags {
		if c.VendorHidden {
			continue
		}

		if _, ok := ctx.matchFlagTitle(c, titleChecking, sw1, sw2); ok {
			// if /*ix == len(args)-1-1 && args[ix+1] == "" &&*/ exact {
			//	flgMatched, err = c, nil
			// } else {
			//	err = nil
			// }
			flgMatched, err = c, nil
		}
	}
	if parent.owner != nil && parent.owner != parent {
		parent = parent.owner
		goto goUp
	}
	return
}

func (ctx *queryShcompContext) matchFlagTitle(c *Flag, titleChecking string, sw1, sw2 bool) (exact, ok bool) {
	if len(titleChecking) == 1 && sw1 {
		ctx.matchedPrecededFlagList["--"+c.GetTitleName()] = c
		ok = true
		return
	}

	if len(titleChecking) == 2 && sw1 && sw2 {
		ctx.matchedPrecededFlagList["--"+c.GetTitleName()] = c
		ok = true
		return
	}

	if sw1 && sw2 {
		exact, ok = ctx.matchLongFlagTitle(c, titleChecking)
	}

	if !ok && c.Short != "" && sw1 {
		exact, ok = ctx.matchShortFlagTitle(c, titleChecking)
	}

	return
}

func (ctx *queryShcompContext) matchShortFlagTitle(c *Flag, titleChecking string) (exact, ok bool) {
	t := titleChecking[1:]
	if c.Short == t {
		ctx.matchedFlagList["-"+c.Short] = c
		exact, ok = true, true
	} else if strings.HasPrefix(c.Short, t) {
		ctx.matchedPrecededFlagList["-"+c.Short] = c
		ok = true
	}
	return
}

func (ctx *queryShcompContext) matchLongFlagTitle(c *Flag, titleChecking string) (exact, ok bool) {
	t := titleChecking[2:]
	switch {
	case c.Full == t:
		ctx.matchedFlagList["--"+c.GetTitleName()] = c
		exact, ok = true, true
	case strings.HasPrefix(c.Full, t):
		ctx.matchedPrecededFlagList["--"+c.GetTitleName()] = c
		ok = true
	default:
		for _, st := range c.Aliases {
			if st == t {
				ctx.matchedFlagList["--"+c.GetTitleName()] = c
				exact, ok = true, true
				break
			} else if strings.HasPrefix(st, t) {
				ctx.matchedPrecededFlagList["--"+c.GetTitleName()] = c
				ok = true
			}
		}
	}
	return
}

var noPartialMatching = true

const (
	shellCompDirectiveError         = 1
	shellCompDirectiveNoSpace       = 2
	shellCompDirectiveNoFileComp    = 4
	shellCompDirectiveFilterFileExt = 8
	shellCompDirectiveFilterDirs    = 16
)

func directivesToString(d int) string {
	// var sb strings.Builder
	if d&shellCompDirectiveError != 0 {
		return "ShellCompDirectiveError"
	}
	if d&shellCompDirectiveNoSpace != 0 {
		return "ShellCompDirectiveNoSpace"
	}
	if d&shellCompDirectiveNoFileComp != 0 {
		return "ShellCompDirectiveNoFileComp"
	}
	if d&shellCompDirectiveFilterFileExt != 0 {
		return "ShellCompDirectiveFilterFileExt"
	}
	if d&shellCompDirectiveFilterDirs != 0 {
		return "ShellCompDirectiveFilterDirs"
	}
	return ""
}
