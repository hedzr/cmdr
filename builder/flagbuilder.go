// Copyright © 2022 Hedzr Yeh.

package builder

import (
	"github.com/hedzr/cmdr/v2/cli"
)

// func NewFlagBuilder(parent *cli.CmdS, defaultValue any, longTitle string, titles ...string) *ffb {
// 	// s := &ffb{
// 	// 	nil, parent,
// 	// 	new(cli.Flag),
// 	// }
// 	// s.Long, s.Short, s.Aliases = theTitles(longTitle, titles...)
// 	// return s
// 	return newFlagBuilderFrom(parent, nil, defaultValue, longTitle, titles...)
// }

func newFlagBuilderShort(b buildable, longTitle string, titles ...string) *ffb {
	return newFlagBuilderFrom(nil, b, nil, longTitle, titles...)
}

func newFlagBuilderFrom(parent *cli.CmdS, b buildable, defaultValue any, longTitle string, titles ...string) *ffb {
	s := &ffb{
		b, parent,
		new(cli.Flag),
	}
	s.Flag.SetDefaultValue(defaultValue)
	s.Long, s.Short, s.Aliases = theTitles(longTitle, titles...)
	return s
}

type ffb struct {
	buildable
	parent *cli.CmdS
	*cli.Flag
}

func (s *ffb) With(cb func(b cli.FlagBuilder)) {
	defer s.Build()
	cb(s)
}

func (s *ffb) Build() {
	if a, ok := s.buildable.(adder); ok {
		a.addFlag(s.Flag)
	}
	// if s.parent != nil {
	// 	s.parent.AddFlag(s.Flag)
	// 	// if a, ok := s.b.(adder); ok {
	// 	// 	a.addFlag(s.Flag)
	// 	// }
	// }
}

func (s *ffb) SetApp(app buildable) { s.buildable = app }

func (s *ffb) Titles(longTitle string, titles ...string) cli.FlagBuilder {
	s.Long, s.Short, s.Aliases = theTitles(longTitle, titles...)
	return s
}

func (s *ffb) Default(defaultValue any) cli.FlagBuilder {
	s.Flag.SetDefaultValue(defaultValue)
	return s
}

func (s *ffb) ExtraShorts(shorts ...string) cli.FlagBuilder {
	s.Flag.SetShorts(shorts...)
	return s
}

func (s *ffb) Description(description string, longDescription ...string) cli.FlagBuilder {
	s.Flag.SetDescription(description, longDescription...)
	return s
}

func (s *ffb) Examples(examples string) cli.FlagBuilder {
	s.Flag.SetExamples(examples)
	return s
}

func (s *ffb) Group(group string) cli.FlagBuilder {
	s.Flag.SetGroup(group)
	return s
}

func (s *ffb) Deprecated(deprecated string) cli.FlagBuilder {
	s.Flag.SetDeprecated(deprecated)
	return s
}

func (s *ffb) Hidden(hidden bool, vendorHidden ...bool) cli.FlagBuilder {
	s.Flag.SetHidden(hidden, vendorHidden...)
	return s
}

func (s *ffb) ToggleGroup(group string) cli.FlagBuilder {
	s.Flag.SetToggleGroup(group)
	return s
}

func (s *ffb) PlaceHolder(placeHolder string) cli.FlagBuilder {
	s.Flag.SetPlaceHolder(placeHolder)
	return s
}

func (s *ffb) DefaultValue(val any) cli.FlagBuilder {
	s.Flag.SetDefaultValue(val)
	return s
}

func (s *ffb) EnvVars(vars ...string) cli.FlagBuilder {
	s.Flag.SetEnvVars(vars...)
	return s
}

func (s *ffb) AppendEnvVars(vars ...string) cli.FlagBuilder {
	s.Flag.AppendEnvVars(vars...)
	return s
}

func (s *ffb) ExternalEditor(externalEditor string) cli.FlagBuilder {
	s.Flag.SetExternalEditor(externalEditor)
	return s
}

func (s *ffb) ValidArgs(validArgs ...string) cli.FlagBuilder {
	s.Flag.SetValidArgs(validArgs...)
	return s
}

func (s *ffb) AppendValidArgs(validArgs ...string) cli.FlagBuilder {
	s.Flag.AppendValidArgs(validArgs...)
	return s
}

func (s *ffb) Range(min, max int) cli.FlagBuilder {
	s.Flag.SetRange(min, max)
	return s
}

func (s *ffb) HeadLike(headLike bool, bounds ...int) cli.FlagBuilder { //nolint:revive
	s.Flag.SetHeadLike(headLike)
	return s
}

func (s *ffb) Required(required bool) cli.FlagBuilder {
	s.Flag.SetRequired(required)
	return s
}

func (s *ffb) CompJustOnce(justOnce bool) cli.FlagBuilder {
	s.Flag.SetJustOnce(justOnce)
	return s
}

func (s *ffb) CompActionStr(action string) cli.FlagBuilder {
	s.Flag.SetActionStr(action)
	return s
}

func (s *ffb) CompMutualExclusives(ex ...string) cli.FlagBuilder {
	s.Flag.SetMutualExclusives(ex...)
	return s
}

func (s *ffb) CompPrerequisites(flags ...string) cli.FlagBuilder {
	s.Flag.SetPrerequisites(flags...)
	return s
}

func (s *ffb) CompCircuitBreak(cb bool) cli.FlagBuilder {
	s.Flag.SetCircuitBreak(cb)
	return s
}

func (s *ffb) DoubleTildeOnly(b bool) cli.FlagBuilder {
	s.Flag.SetDoubleTildeOnly(b)
	return s
}

func (s *ffb) OnParseValue(handler cli.OnParseValueHandler) cli.FlagBuilder {
	s.Flag.SetOnParseValueHandler(handler)
	return s
}

func (s *ffb) OnMatched(handler cli.OnMatchedHandler) cli.FlagBuilder {
	s.Flag.SetOnMatchedHandler(handler)
	return s
}

func (s *ffb) OnChanging(handler cli.OnChangingHandler) cli.FlagBuilder {
	s.Flag.SetOnChangingHandler(handler)
	return s
}

func (s *ffb) OnChanged(handler cli.OnChangedHandler) cli.FlagBuilder {
	s.Flag.SetOnChangedHandler(handler)
	return s
}

func (s *ffb) OnSet(handler cli.OnSetHandler) cli.FlagBuilder {
	s.Flag.SetOnSetHandler(handler)
	return s
}
