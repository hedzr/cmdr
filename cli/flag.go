package cli

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hedzr/is/term/color"
	"github.com/hedzr/store"
)

type navigator interface { //nolint:unused
	Root() *RootCommand
	Owner() *CmdS
}

func (f *Flag) Owner() *CmdS {
	if cx, ok := f.owner.(*CmdS); ok {
		return cx
	}
	return nil
}
func (f *Flag) OwnerOrParent() Backtraceable { return f.owner } // the owner of this CmdS
func (f *Flag) OwnerCmd() Cmd                { return f.owner }

// Store returns the commands subset of the application Store.
func (f *Flag) Store() store.Store {
	if f.owner != nil {
		return f.owner.Store()
	}
	return nil
}

func (f *Flag) IsToggleGroup() bool { return f.toggleGroup != "" }

func (f *Flag) ToggleGroupLeadHelpString() (lead string) { //nolint:revive
	if f.toggleGroup != "" {
		var state, b bool
		if fo := f.Owner(); fo != nil {
			if m, ok := fo.toggles[f.toggleGroup]; ok {
				if _, ok = m.Flags[f.Name()]; ok {
					if b, ok = f.defaultValue.(bool); ok {
						state = b
					}
				}
			}
			if state {
				lead = "[x] "
			} else {
				lead = "[ ] "
			}
		}
	}
	return
}

func (f *Flag) MatchedTG() (tgm *ToggleGroupMatch) {
	if f.toggleGroup != "" {
		if fo := f.Owner(); fo != nil {
			if m, ok := fo.toggles[f.toggleGroup]; ok {
				tgm = m
			}
		}
	}
	return
}

func (f *Flag) NeedParseValue() bool {
	switch f.defaultValue.(type) {
	case bool:
		return false
	default:
		return true
	}
}

type transFunc func(ss string, clr color.Color) string

func (f *Flag) DefaultValueHelpString(trans transFunc, clr, clrDefault color.Color) (hs, plain string) {
	if f.defaultValue != nil {
		if f.placeHolder != "" {
			plain = fmt.Sprintf("(Default: %s=%v)", f.placeHolder, f.defaultValue)
			hs = trans(fmt.Sprintf("(Default: <font color=%v>%s</font>=<font color=%v>%v</font>)",
				color.ToColorString(clr), f.placeHolder, color.ToColorString(clr), f.defaultValue), clrDefault)
		} else {
			plain = fmt.Sprintf("(Default: %v)", f.defaultValue)
			hs = trans(fmt.Sprintf("(Default: <font color=%v>%v</font>)",
				color.ToColorString(clr), f.defaultValue), clrDefault)
		}
	}
	return
}

func (f *Flag) ToggleGroup() string        { return f.toggleGroup }
func (f *Flag) PlaceHolder() string        { return f.placeHolder }
func (f *Flag) DefaultValue() any          { return f.defaultValue }
func (f *Flag) EnvVars() []string          { return f.envVars }
func (f *Flag) ExternalEditor() string     { return f.externalEditor }
func (f *Flag) ValidArgs() []string        { return f.validArgs }
func (f *Flag) Range() (min, max int)      { return f.min, f.max }
func (f *Flag) HeadLike() bool             { return f.headLike }
func (f *Flag) Required() bool             { return f.required }
func (f *Flag) JustOnce() bool             { return f.justOnce }
func (f *Flag) ActionStr() string          { return f.actionStr }
func (f *Flag) MutualExclusives() []string { return f.mutualExclusives }
func (f *Flag) Prerequisites() []string    { return f.prerequisites }
func (f *Flag) CircuitBreak() bool         { return f.circuitBreak }
func (f *Flag) DoubleTildeOnly() bool      { return f.dblTildeOnly }

func (f *Flag) SetToggleGroup(group string)       { f.toggleGroup = group }
func (f *Flag) SetPlaceHolder(placeHolder string) { f.placeHolder = placeHolder }
func (f *Flag) SetDefaultValue(val any)           { f.defaultValue = val }
func (f *Flag) SetEnvVars(vars ...string) {
	for _, v := range vars {
		if v != "" {
			f.envVars = append(f.envVars, v)
		}
	}
}
func (f *Flag) AppendEnvVars(vars ...string)            { f.envVars = append(f.envVars, vars...) }
func (f *Flag) SetExternalEditor(externalEditor string) { f.externalEditor = externalEditor }
func (f *Flag) SetValidArgs(validArgs ...string)        { f.validArgs = validArgs }
func (f *Flag) AppendValidArgs(validArgs ...string)     { f.validArgs = append(f.validArgs, validArgs...) }
func (f *Flag) SetRange(min, max int)                   { f.min, f.max = min, max }
func (f *Flag) SetHeadLike(headLike bool)               { f.headLike = headLike }
func (f *Flag) SetRequired(required bool)               { f.required = required }
func (f *Flag) SetJustOnce(justOnce bool)               { f.justOnce = justOnce }
func (f *Flag) SetActionStr(action string)              { f.actionStr = action }
func (f *Flag) SetMutualExclusives(ex ...string)        { f.mutualExclusives = ex }
func (f *Flag) SetPrerequisites(flags ...string)        { f.prerequisites = flags }
func (f *Flag) SetCircuitBreak(cb bool)                 { f.circuitBreak = cb }
func (f *Flag) SetDoubleTildeOnly(b bool)               { f.dblTildeOnly = b }

// GetDottedNamePath return the dotted key path of this flag
// in the options store.
func (f *Flag) GetDottedNamePath() string {
	if f.owner != nil {
		if op := f.owner.GetDottedPath(); op != "" {
			return f.owner.GetDottedPath() + "." + f.GetTitleName()
		}
	}
	return f.GetTitleName()
}

func (f *Flag) TryOnParseValue(index int, hitCaption, hitValue string, args []string) ( //nolint:revive
	handled bool, newVal any, remainsPartInHitValue string, err error,
) {
	if f.onParseValue != nil {
		handled = true
		newVal, remainsPartInHitValue, err = f.onParseValue(f, index, hitCaption, hitValue, args)
		if !f.errIsSignalFallback(err) {
			err, handled = nil, false
		}
	}
	return
}

func (f *Flag) TryOnMatched(position int, hitState *MatchState) (handled bool, err error) {
	if f.onMatched != nil {
		handled = true
		err = f.onMatched(f, position, hitState)
		if !f.errIsSignalFallback(err) {
			err, handled = nil, false
		}
	}
	return
}

func (f *Flag) TryOnChanging(oldVal, newVal any) (handled bool, err error) {
	if f.onChanging != nil {
		handled = true
		err = f.onChanging(f, oldVal, newVal)
		if !f.errIsSignalFallback(err) {
			err, handled = nil, false
		}
	}
	return
}

func (f *Flag) TryOnChanged(oldVal, newVal any) {
	if f.onChanged != nil {
		f.onChanged(f, oldVal, newVal)
	}
}

func (f *Flag) TryOnSet(oldVal, newVal any) {
	if f.onSet != nil {
		f.onSet(f, oldVal, newVal)
	}
}

func (f *Flag) SetOnParseValueHandler(handler OnParseValueHandler) {
	f.onParseValue = handler
}

func (f *Flag) SetOnMatchedHandler(handler OnMatchedHandler) {
	f.onMatched = handler
}

func (f *Flag) SetOnChangingHandler(handler OnChangingHandler) {
	f.onChanging = handler
}

func (f *Flag) SetOnChangedHandler(handler OnChangedHandler) {
	f.onChanged = handler
}

func (f *Flag) SetOnSetHandler(handler OnSetHandler) {
	f.onSet = handler
}

func (f *Flag) SetNegatable(b bool, items ...string) {
	f.negatable = b
	f.negItems = items
}

// Negatable flag supports auto-orefixing by `--no-`.
//
// For a flag named as 'warning`, both `--warning` and
// `--no-warning` are available in cmdline.
//
// While items slice are supplied, it would create a
// `-W`-style nesatable flag. It's a gcc `-W` style
// flag, for example, `-Wunused-variable` in gcc is raising
// a warning for unused variables, `-Wno-unused-variable`
// is disabling warning for unused variable.
//
// So, our `-W`-style can be building with:
//
//	parent.Flg("warnings", "W").
//	  Description("gcc-style negatable flag: <code>-Wunused-variable</code> and -Wno-unused-variable", "").
//	  Group("Negatable").
//	  Negatable(true, "unused-variable", "unused-parameter",
//	    "unused-function", "unused-but-set-variable",
//	    "unused-private-field", "unused-label").
//	  Default(false).
//	  Build()
//
// For a negatable flag `--warning`, extracting the final
// flag values from `cmd.Store().MustBool("warning")` and
// `cmd.Store().MustBool("no-warning")`.
//
// For a `-W`-style nesatable flag `--warnings`/`-W` in
// above sample code, extracting the final values from
// `cmd.Store().MustBool("warnings.unused-variable")`
// and so on. And you can also extract a selected items
// set from `cmd.Store().MustStringSlice("warnings.selected")`.
//
// See also [Flag.NegatableItems]
func (f *Flag) Negatable() bool { return f.negatable }

// NegatableItems provides the alternative items for a `-W`-style
// negatable flag.
//
// See also [Flag.Negatable]
func (f *Flag) NegatableItems() []string { return f.negItems }

func (f *Flag) SetLeadingPlusSign(b bool) {
	f.leadingPlusSign = b
}

func (f *Flag) LeadingPlusSign() bool { return f.leadingPlusSign }

//
//

// func (f *Flag) Root() *RootCommand { return c.root }
// func (f *Flag) Owner() *CmdS    { return c.owner }

// GetTriggeredTimes returns the matched times
func (f *Flag) GetTriggeredTimes() int    { return f.hitTimes }
func (f *Flag) GetTriggeredTitle() string { return f.hitTitle }

// GetDescZsh temp
func (f *Flag) GetDescZsh() (desc string) {
	desc = f.description
	if desc == "" {
		desc = EraseAnyWSs(f.GetTitleZshFlagName())
	}
	// desc = replaceAll(desc, " ", "\\ ")
	desc = reSQ.ReplaceAllString(desc, `*$1*`)
	desc = reBQ.ReplaceAllString(desc, `**$1**`)
	desc = reSQnp.ReplaceAllString(desc, "''")
	desc = reBQnp.ReplaceAllString(desc, "\\`")
	desc = strings.ReplaceAll(desc, ":", "\\:")
	desc = strings.ReplaceAll(desc, "[", "\\[")
	desc = strings.ReplaceAll(desc, "]", "\\]")
	return
}

// GetTitleFlagNames temp
func (f *Flag) GetTitleFlagNames() (title, rest string) {
	return f.GetTitleFlagNamesBy(",")
}

// GetTitleZshFlagName temp
func (f *Flag) GetTitleZshFlagName() (str string) {
	if len(f.Long) > 0 {
		str += "--" + f.Long
	} else if len(f.Short) > 0 {
		str += "-" + f.Short
	}
	return
}

// GetTitleZshFlagShortName temp
func (f *Flag) GetTitleZshFlagShortName() (str string) {
	if len(f.Short) > 0 {
		str += "-" + f.Short
	} else if len(f.Long) > 0 {
		str += "--" + f.Long
	}
	return
}

// GetTitleZshNamesBy temp
func (f *Flag) GetTitleZshNamesBy(delimiter string, allowPrefix, quoted bool) (str string) {
	return f.GetTitleZshNamesExtBy(delimiter, allowPrefix, quoted, true, true)
}

// GetTitleZshNamesExtBy temp
func (f *Flag) GetTitleZshNamesExtBy(delimiter string, allowPrefix, quoted, shortTitleOnly, longTitleOnly bool) (str string) { //nolint:revive
	// quote := false
	prefix, suffix := "", ""
	if _, ok := f.defaultValue.(bool); !ok {
		suffix = "="
		// } else if _, ok := s.DefaultValue.(bool); ok {
		//	suffix = "-"
	}
	if allowPrefix && !f.justOnce {
		quoted, prefix = true, "*" //nolint:revive
	}
	if !longTitleOnly && len(f.Short) > 0 {
		if quoted {
			str += "'" + prefix + "-" + f.Short + suffix + "'"
		} else {
			str += prefix + "-" + f.Short + suffix
		}
		if shortTitleOnly {
			return
		}
	}
	if len(f.Long) > 0 {
		if str != "" {
			str += delimiter
		}
		if quoted {
			str += "'" + prefix + "--" + f.Long + suffix + "'"
		} else {
			str += prefix + "--" + f.Long + suffix
		}
	}
	return
}

// GetTitleZshFlagNamesArray temp
func (f *Flag) GetTitleZshFlagNamesArray() (ary []string) {
	if len(f.Short) == 1 || len(f.Short) == 2 {
		if len(f.placeHolder) > 0 {
			ary = append(ary, "-"+f.Short+"=") // +s.PlaceHolder)
		} else {
			ary = append(ary, "-"+f.Short)
		}
	}
	if len(f.Long) > 0 {
		if len(f.placeHolder) > 0 {
			ary = append(ary, "--"+f.Long+"=") // +s.PlaceHolder)
		} else {
			ary = append(ary, "--"+f.Long)
		}
	}
	return
}

// GetTitleFlagNamesBy temp
//
// A title with prefix '_'/'__' will be hidden from the
// result array.
func (f *Flag) GetTitleFlagNamesBy(delimiter string, maxW ...int) (title, rest string) {
	return f.GetTitleFlagNamesByMax(delimiter, len(f.Short), maxW...)
}

// GetTitleFlagNamesByMax temp
//
// A title with prefix '_'/'__' will be hidden from the
// result array.
func (f *Flag) GetTitleFlagNamesByMax(delimiter string, maxShort int, maxW ...int) (title, rest string) {
	var sb, sbr strings.Builder

	if f.Short == "" {
		// if no flag.Short,
		_, _ = sb.WriteString(strings.Repeat(" ", maxShort))
	} else {
		_, _ = sb.WriteRune('-')
		_, _ = sb.WriteString(f.Short)
		_, _ = sb.WriteString(delimiter)
		if len(f.Short) < maxShort {
			_, _ = sb.WriteString(strings.Repeat(" ", maxShort-len(f.Short)))
		}
	}

	if f.Short == "" {
		_, _ = sb.WriteRune(' ')
		_, _ = sb.WriteRune(' ')
	}
	_, _ = sb.WriteRune(' ')
	if f.dblTildeOnly {
		_, _ = sb.WriteString("~~")
	} else {
		_, _ = sb.WriteString("--")
	}
	_, _ = sb.WriteString(f.Long)
	if len(f.placeHolder) > 0 {
		// str += fmt.Sprintf("=\x1b[2m\x1b[%dm%s\x1b[0m", DarkColor, s.PlaceHolder)
		_, _ = sb.WriteString(fmt.Sprintf("=%s", f.placeHolder))
	}

	si := len(f.Aliases)
	ll := sb.Len()
	maxWidth := 0
	for _, w := range maxW {
		maxWidth = max(maxWidth, w)
	}
	for i, sz := range f.Aliases {
		if sz == "" || strings.HasPrefix(sz, "_") {
			continue
		}

		if i > si {
			_, _ = sbr.WriteString(delimiter)
			_, _ = sbr.WriteString("--")
			_, _ = sbr.WriteString(sz)
			continue
		}
		if maxWidth > 0 && ll+len(delimiter)+2+len(sz) >= maxWidth {
			_, _ = sbr.WriteString("--")
			_, _ = sbr.WriteString(sz)
			si = i
			continue
		}

		ll += len(delimiter) + 2 + len(sz)
		_, _ = sb.WriteString(delimiter)
		_, _ = sb.WriteString("--")
		_, _ = sb.WriteString(sz)
	}
	title, rest = sb.String(), sbr.String()
	return
}

func (f *Flag) String() string {
	var sb strings.Builder
	_, _ = sb.WriteString("Flg{'")
	_, _ = sb.WriteString(f.GetDottedNamePath())
	_, _ = sb.WriteString("'}")
	return sb.String()
}

// EqualTo _
func (f *Flag) EqualTo(rh *Flag) (ok bool) {
	if f == nil {
		return rh == nil
	}
	if rh == nil {
		return false
	}
	return f.GetTitleName() == rh.GetTitleName()
}

// Delete removes myself from the command owner.
func (f *Flag) Delete() {
	if f == nil || f.owner == nil {
		return
	}

	if fo := f.Owner(); fo != nil {
		for i, cc := range fo.flags {
			if f == cc {
				fo.flags = append(fo.flags[0:i], fo.flags[i+1:]...)
				return
			}
		}
	}
}

func (f *Flag) ensureXref() {
	//
}

func (f *Flag) EnvVarsHelpString(trans func(ss string, clr color.Color) string, clr, clrDefault color.Color) (hs, plain string) {
	if len(f.envVars) != 0 {
		var envVars []string
		for _, v := range f.envVars {
			if v != "" {
				envVars = append(envVars, v)
			}
		}
		dep := strings.Join(envVars, ",")
		plain = fmt.Sprintf("[Env: %s]", dep)
		hs = trans(fmt.Sprintf("[Env: <font color=%v>%s</font>]", color.ToColorString(clr), dep), clrDefault)
	}
	return
}

func (f *Flag) Clone() any {
	return &Flag{
		BaseOpt: *(f.BaseOpt.Clone().(*BaseOpt)),

		toggleGroup:  f.toggleGroup,
		placeHolder:  f.placeHolder,
		defaultValue: f.defaultValue,
		envVars:      slices.Clone(f.envVars),

		externalEditor: f.externalEditor,
		validArgs:      slices.Clone(f.validArgs),
		min:            f.min,
		max:            f.max,
		headLike:       f.headLike,
		required:       f.required,

		onParseValue: f.onParseValue,
		onMatched:    f.onMatched,
		onChanging:   f.onChanging,
		onChanged:    f.onChanged,
		onSet:        f.onSet,

		actionStr:        f.actionStr,
		mutualExclusives: slices.Clone(f.mutualExclusives),
		prerequisites:    slices.Clone(f.prerequisites),
		justOnce:         f.justOnce,
		circuitBreak:     f.circuitBreak,
		dblTildeOnly:     f.dblTildeOnly,
		negatable:        f.negatable,
		negItems:         slices.Clone(f.negItems),

		leadingPlusSign: f.leadingPlusSign,
	}
}
