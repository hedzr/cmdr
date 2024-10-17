package cli

import (
	"context"
	"strconv"
	"strings"

	"github.com/hedzr/cmdr/v2/cli/atoa"
	"github.com/hedzr/cmdr/v2/internal/tool"
	"github.com/hedzr/evendeep/ref"
	logz "github.com/hedzr/logg/slog"
)

func (c *CmdS) Match(ctx context.Context, title string) (short bool, cc Cmd) {
	if title == "" {
		return
	}

	c.ensureXrefCommands(ctx)

	var ok bool
	if cc, ok = c.longCommands[title]; ok {
		cc.SetHitTitle(title)
		return
	}
	if cc, short = c.shortCommands[title]; short {
		cc.SetHitTitle(title)
		return
	}

	if c.onEvalSubcommandsOnce != nil || c.onEvalSubcommands != nil {
		commands := mustEnsureDynCommands(ctx, c)
		for _, cx := range commands {
			if title == cx.Name() {
				cx.SetHitTitle(title)
				cc = cx
				return
			}
			for _, ttl := range cx.AliasNames() {
				if title == ttl {
					cx.SetHitTitle(title)
					cc = cx
					return
				}
			}
		}
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

// MatchFlag try matching command title with vp.Remains, and update the relevant states.
//
// While a flag matched ok, returns vp.Matched != "" && ff != nil && err != nil
func (c *CmdS) MatchFlag(ctx context.Context, vp *FlagValuePkg) (ff *Flag, err error) { //nolint:revive
	c.ensureXrefFlags(ctx)

	var ok bool
	var matched, remains string
	if vp.Short {
		// short flag

		var cclist map[string]*Flag
		if c.onEvalFlagsOnce != nil || c.onEvalFlags != nil {
			flags := mustEnsureDynFlags(ctx, c)
			cclist = make(map[string]*Flag)
			for _, cx := range flags {
				if cx.Short != "" {
					cclist[cx.Short] = cx
				}
			}
		} else {
			cclist = c.shortFlags
		}

		if ff, ok = cclist[vp.Remains]; ok && c.testDblTilde(vp.SpecialTilde, ff) {
			vp.PartialMatched, vp.Matched, vp.Remains, ff.hitTitle, ff.hitTimes = false, vp.Remains, "", vp.Remains, ff.hitTimes+1
			return c.tryParseValue(ctx, vp, ff)
		}

		// try for compact short flags
		matched, remains, ff, err = c.partialMatchFlag(ctx, vp.Remains, vp.Short, vp.SpecialTilde, cclist)
		if vp.PartialMatched = ff != nil && err == nil; vp.PartialMatched {
			vp.Matched, vp.Remains = matched, remains
			ff, err = c.tryParseValue(ctx, vp, ff)
		}

		// try to parse headLike flag
		if vp.Matched == "" && c.headLikeFlag != nil && ref.IsNumeric(c.headLikeFlag.defaultValue) {
			var num int64
			if num, err = strconv.ParseInt(vp.Remains, 0, 64); err == nil {
				vp.Matched, vp.Remains, ff = vp.Remains, "", c.headLikeFlag
				ff.defaultValue, vp.ValueOK = int(num), true // store the parsed value
				logz.VerboseContext(ctx, "[cmdr] headLike flag matched", "flg", ff, "num", num)
			}
		}
	} else {
		var cclist map[string]*Flag
		if c.onEvalFlagsOnce != nil || c.onEvalFlags != nil {
			flags := mustEnsureDynFlags(ctx, c)
			cclist = make(map[string]*Flag)
			for _, cx := range flags {
				if cx.Long != "" {
					cclist[cx.Long] = cx
				}
				for _, t := range cx.Aliases {
					if t != "" {
						cclist[t] = cx
					}
				}
			}
		} else {
			cclist = c.longFlags
		}

		if ff, ok = cclist[vp.Remains]; ok && c.testDblTilde(vp.SpecialTilde, ff) {
			vp.PartialMatched, vp.Matched, vp.Remains, ff.hitTitle, ff.hitTimes = false, vp.Remains, "", vp.Remains, ff.hitTimes+1
			return c.tryParseValue(ctx, vp, ff)
		}

		matched, remains, ff, err = c.partialMatchFlag(ctx, vp.Remains, vp.Short, vp.SpecialTilde, cclist)
		if vp.PartialMatched = ff != nil && err == nil; vp.PartialMatched {
			vp.Matched, vp.Remains = matched, remains
			ff, err = c.tryParseValue(ctx, vp, ff)
		}
	}

	// lookup the parents, if 'ff' not matched/found
	if ff == nil && err == nil && c.owner != nil && c.owner != c {
		ff, err = c.owner.MatchFlag(ctx, vp)
		return
	}

	// when a flag matched ok
	if ff != nil && err == nil && vp.Matched != "" {
		ff.hitTitle = vp.Matched
		ff.hitTimes++
		if !vp.ValueOK {
			ff, err = c.tryParseValue(ctx, vp, ff)
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

func (c *CmdS) testDblTilde(dblTilde bool, ff *Flag) (matched bool) {
	matched = dblTilde || !ff.dblTildeOnly || (ff.dblTildeOnly && dblTilde)
	return
}

func (c *CmdS) partialMatchFlag(ctx context.Context, title string, short, dblTildeMode bool, cclist map[string]*Flag) (matched, remains string, ff *Flag, err error) { //nolint:revive
	var maxLen int
	var rightPart string

	titleOriginal := title
	if pos := strings.IndexRune(title, '='); pos >= 0 {
		rightPart = title[pos+1:]
		title = title[:pos] //nolint:revive
	}

	for k, v := range cclist {
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
		if c.OwnerIsValid() {
			if co, ok := c.owner.(*CmdS); ok {
				co.ensureXrefFlags(ctx)
				mf := co.longFlags
				if short {
					mf = co.shortFlags
				}
				matched1, remains1, ff1, err1 := co.partialMatchFlag(ctx, titleOriginal, short, dblTildeMode, mf)
				if err = err1; err != nil {
					return
				}
				if ff1 != nil && maxLen < len(matched1) {
					// if longer matched flag from parents exists, use it instead of the lastCommand's
					matched, remains, ff = matched1, remains1, ff1
				}
			}
		}
		return
	}

	if c.OwnerIsValid() {
		if co, ok := c.owner.(*CmdS); ok {
			// if no flag matched, checking the parents
			co.ensureXrefFlags(ctx)

			var cxlist map[string]*Flag
			if short {
				if co.onEvalFlagsOnce != nil || co.onEvalFlags != nil {
					flags := mustEnsureDynFlags(ctx, c.owner)
					cxlist = make(map[string]*Flag)
					for _, cx := range flags {
						if cx.Short != "" {
							cxlist[cx.Short] = cx
						}
					}
				} else {
					cxlist = co.shortFlags
				}
			} else {
				if co.onEvalFlagsOnce != nil || co.onEvalFlags != nil {
					flags := mustEnsureDynFlags(ctx, c.owner)
					cxlist = make(map[string]*Flag)
					for _, cx := range flags {
						if cx.Long != "" {
							cxlist[cx.Long] = cx
						}
						for _, t := range cx.Aliases {
							if t != "" {
								cxlist[t] = cx
							}
						}
					}
				} else {
					cxlist = co.longFlags
				}
			}

			matched, remains, ff, err = co.partialMatchFlag(ctx, titleOriginal, short, dblTildeMode, cxlist)
		}
	}
	return
}

func (c *CmdS) tryParseValue(ctx context.Context, vp *FlagValuePkg, ff *Flag) (ret *Flag, err error) {
	if ff != nil {
		ff = c.matchedForTG(ctx, ff) //nolint:revive
	}
	if ff, err = c.checkPrerequisites(ctx, vp, ff); err != nil {
		return
	}
	if ff, err = c.checkJustOnce(vp, ff); err != nil {
		return
	}

	if ff != nil && !vp.ValueOK {
		// try to parse value
		switch ff.defaultValue.(type) {
		case bool:
			ff = c.tryParseBoolValue(ctx, vp, ff) //nolint:revive
		case string:
			ff = c.tryParseStringValue(ctx, vp, ff) //nolint:revive
		case nil:
			ff = c.tryParseStringValue(ctx, vp, ff) //nolint:revive
		default:
			ff = c.tryParseOthersValue(ctx, vp, ff) //nolint:revive
		}
	}

	ret = ff
	_, err = c.checkCircuitBreak(vp, ff)
	return
}

func (c *CmdS) matchedForTG(ctx context.Context, ff *Flag) *Flag {
	// toggle group
	if co, ok := ff.owner.(*CmdS); ok {
		if co.toggles != nil {
			if m, ok := co.toggles[ff.ToggleGroup()]; ok {
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
			root := ff.Root()
			for _, fn := range ff.mutualExclusives {
				var f *Flag
				if strings.ContainsRune(fn, '.') {
					f = co.FindFlag(ctx, fn, false)
				} else {
					_, f = dottedPathToCommandOrFlagG(root, fn)
				}
				if f != nil {
					if _, ok := f.defaultValue.(bool); ok {
						f.SetDefaultValue(false)
					}
				}
			}
		}
	}
	return ff
}

func (c *CmdS) checkJustOnce(vp *FlagValuePkg, ff *Flag) (ret *Flag, err error) {
	if ff != nil && ff.justOnce {
		if ff.hitTimes > 1 {
			err = ErrFlagJustOnce.FormatWith(ff)
			return
		}
	}
	ret, _ = ff, vp
	return
}

func (c *CmdS) checkPrerequisites(ctx context.Context, vp *FlagValuePkg, ff *Flag) (ret *Flag, err error) {
	if ff != nil && len(ff.prerequisites) > 0 {
		if co, ok := c.owner.(*CmdS); ok {
			root := ff.Root()
			for _, fn := range ff.prerequisites {
				var f *Flag
				if strings.ContainsRune(fn, '.') {
					f = co.FindFlag(ctx, fn, false)
				} else {
					_, f = dottedPathToCommandOrFlagG(root, fn)
				}
				if f != nil {
					if f.hitTimes < 0 {
						err = ErrMissedPrerequisite.FormatWith(ff, f)
						return
					}
				}
			}
		}
	}
	ret, _ = ff, vp
	return
}

func (c *CmdS) checkCircuitBreak(vp *FlagValuePkg, ff *Flag) (ret *Flag, err error) {
	if ff != nil && ff.circuitBreak {
		err = ErrShouldStop
		return
	}
	ret, _ = ff, vp
	return
}

func (c *CmdS) tryParseStringValue(ctx context.Context, vp *FlagValuePkg, ff *Flag) *Flag {
	if ff.externalEditor != "" {
		if f := c.invokeExternalEditor(ctx, vp, ff); f != nil {
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

func (c *CmdS) tryParseBoolValue(ctx context.Context, vp *FlagValuePkg, ff *Flag) *Flag {
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

func (c *CmdS) tryParseOthersValue(ctx context.Context, vp *FlagValuePkg, ff *Flag) *Flag {
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

func (c *CmdS) fromString(text string, meme any) (value any) { //nolint:revive
	var err error
	value, err = atoa.Parse(text, meme)
	if err != nil {
		logz.ErrorContext(context.TODO(), "[cmdr] cannot parse text to value", "err", err, "text", text, "target-value-meme", meme)
	}
	return
}

func (c *CmdS) normalizeStringValue(sv string) string {
	return tool.StripQuotes(sv)
}

func (c *CmdS) TryOnMatched(position int, hitState *MatchState) (handled bool, err error) {
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

// MatchTitleNameFast matches a given title string without indices built.
func (c *CmdS) MatchTitleNameFast(title string) (ok bool) { //nolint:revive
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
