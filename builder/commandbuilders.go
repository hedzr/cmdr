// Copyright © 2022 Hedzr Yeh.

package builder

import (
	"fmt"
	"sync/atomic"

	"github.com/hedzr/cmdr/v2/cli"
)

type buildable interface {
	Build()
}

// func NewCommandBuilder(parent *cli.CmdS, longTitle string, titles ...string) *ccb {
// 	// s := &ccb{
// 	// 	nil, parent,
// 	// 	new(cli.CmdS),
// 	// 	false, false,
// 	// }
// 	// s.Long, s.Short, s.Aliases = theTitles(longTitle, titles...)
// 	// return s
// 	return newCommandBuilderFrom(parent, nil, longTitle, titles...)
// }

func newCommandBuilderShort(b buildable, longTitle string, titles ...string) *ccb {
	return newCommandBuilderFrom(new(cli.CmdS), b, longTitle, titles...)
}

func newCommandBuilderFrom(from *cli.CmdS, b buildable, longTitle string, titles ...string) *ccb {
	s := &ccb{
		b,
		from,
		new(cli.CmdS),
		0, 0,
	}
	s.Long, s.Short, s.Aliases = theTitles(longTitle, titles...)
	return s
}

func asCommandBuilder(from *cli.CmdS, b buildable, longTitle string, titles ...string) *ccb {
	s := &ccb{
		b,
		nil,
		from,
		0, 0,
	}
	if from.OwnerIsNotNil() {
		s.parent = from.OwnerCmd().(*cli.CmdS)
	}
	s.Long, s.Short, s.Aliases = theTitles(longTitle, titles...)
	return s
}

type ccb struct {
	buildable
	parent *cli.CmdS
	*cli.CmdS
	inCmd int32
	inFlg int32
}

func (s *ccb) Build() {
	if a, ok := s.buildable.(adder); ok {
		a.addCommand(s.CmdS)
	}
	// if s.parent != nil {
	// 	s.parent.AddSubCommand(s.CmdS)
	// 	// if a, ok := s.b.(adder); ok {
	// 	// 	a.addCommand(s.CmdS)
	// 	// }
	// }

	atomic.StoreInt32(&s.inCmd, 0)
	atomic.StoreInt32(&s.inFlg, 0)
}

var _ cli.CommandBuilder = (*ccb)(nil)

func (s *ccb) With(cb func(b cli.CommandBuilder)) {
	defer s.Build()
	cb(s)
}

func (s *ccb) WithSubCmd(cb func(b cli.CommandBuilder)) {
	// if atomic.LoadInt32(&s.inCmd) != 0 {
	// 	panic("cannot call AddCmd() without Build() last Cmd()/AddCmd()")
	// }

	bc := newCommandBuilderShort(s, "new-command")
	defer bc.Build() // `Build' will add `bc'(CmdS) to s.CmdS as its SubCommand
	cb(bc)
	// atomic.AddInt32(&s.inCmd, 1)
	// return s
}

// addCommand adds a in-building Cmd into current CmdS as a child-/sub-command.
// used by adder when ccb.Build.
func (s *ccb) addCommand(child *cli.CmdS) {
	atomic.AddInt32(&s.inCmd, -1) // reset increased inCmd at AddCmd or Cmd
	s.AddSubCommand(child)
}

// addFlag adds a in-building Flg into current CmdS as its flag.
// used by adder when ccb.Build.
func (s *ccb) addFlag(child *cli.Flag) {
	atomic.AddInt32(&s.inFlg, -1)
	s.AddFlag(child)
}

func (s *ccb) NewCommandBuilder(longTitle string, titles ...string) cli.CommandBuilder {
	return s.Cmd(longTitle, titles...)
}

func (s *ccb) NewFlagBuilder(longTitle string, titles ...string) cli.FlagBuilder {
	return s.Flg(longTitle, titles...)
}

func (s *ccb) Cmd(longTitle string, titles ...string) cli.CommandBuilder {
	if atomic.LoadInt32(&s.inCmd) != 0 {
		panic("cannot call Cmd() without Build() last Cmd()/AddCmd()")
	}
	atomic.AddInt32(&s.inCmd, 1)
	return newCommandBuilderShort(s, longTitle, titles...)
}

func (s *ccb) Flg(longTitle string, titles ...string) cli.FlagBuilder {
	if atomic.LoadInt32(&s.inFlg) != 0 {
		panic("cannot call Flg() without Build() last Flg()")
	}
	atomic.AddInt32(&s.inFlg, 1)
	return newFlagBuilderShort(s, longTitle, titles...)
}

// ToggleableFlags creates a batch of toggleable flags.
//
// For example:
//
//	s.ToggleableFlags(
//	  cli.BatchToggleFlag{L: "apple", S: "a"},
//	  cli.BatchToggleFlag{L: "banana"},
//	  cli.BatchToggleFlag{L: "orange", S: "o", DV: true},
//	)
func (s *ccb) ToggleableFlags(toggleGroupName string, items ...cli.BatchToggleFlag) {
	for _, tf := range items {
		s.Flg(tf.L, tf.S).
			DefaultValue(tf.DV).
			ToggleGroup(toggleGroupName).
			Description(fmt.Sprintf("Item of toggle group %q: %q", toggleGroupName, tf.L)).
			Build()
	}
}

func (s *ccb) AddCmd(cb func(b cli.CommandBuilder)) cli.CommandBuilder {
	// if atomic.LoadInt32(&s.inCmd) != 0 {
	// 	panic("cannot call AddCmd() without Build() last Cmd()/AddCmd()")
	// }

	bc := newCommandBuilderShort(s, "new-command")
	defer bc.Build() // `Build' will add `bc'(CmdS) to s.CmdS as its SubCommand
	cb(bc)
	atomic.AddInt32(&s.inCmd, 1)
	return s
}

func (s *ccb) AddFlg(cb func(b cli.FlagBuilder)) cli.CommandBuilder {
	// if atomic.LoadInt32(&s.inFlg) != 0 {
	// 	panic("cannot call AddFlg() without Build() last Flg()/AddFlg()")
	// }

	bc := newFlagBuilderShort(s, "new-flag")
	defer bc.Build() // `Build' will add `bc'(Flag) to s.CmdS as its Flag
	// atomic.AddInt32(&s.inFlg, 1)
	// defer func() { atomic.AddInt32(&s.inFlg, -1) }()
	cb(bc)
	return s
}

//

//

//

func (s *ccb) Titles(longTitle string, titles ...string) cli.CommandBuilder {
	s.Long, s.Short, s.Aliases = theTitles(longTitle, titles...)
	return s
}

func (s *ccb) ExtraShorts(shorts ...string) cli.CommandBuilder {
	s.SetShorts(shorts...)
	return s
}

func (s *ccb) Description(description string, longDescription ...string) cli.CommandBuilder {
	s.SetDescription(description, longDescription...)
	return s
}

func (s *ccb) Examples(examples string) cli.CommandBuilder {
	s.SetExamples(examples)
	return s
}

func (s *ccb) Group(group string) cli.CommandBuilder {
	s.SetGroup(group)
	return s
}

func (s *ccb) Deprecated(deprecated string) cli.CommandBuilder {
	s.SetDeprecated(deprecated)
	return s
}

func (s *ccb) Hidden(hidden bool, vendorHidden ...bool) cli.CommandBuilder {
	s.SetHidden(hidden, vendorHidden...)
	return s
}

func (s *ccb) TailPlaceHolders(placeHolders ...string) cli.CommandBuilder {
	s.SetTailPlaceHolder(placeHolders...)
	return s
}

func (s *ccb) RedirectTo(dottedPath string) cli.CommandBuilder {
	s.SetRedirectTo(dottedPath)
	return s
}

// OnAction sets the onAction handler.
//
// a call to `OnAction(nil)` will set the underlying onAction handlet empty.
func (s *ccb) OnAction(handler cli.OnInvokeHandler) cli.CommandBuilder {
	s.SetAction(handler)
	return s
}

func (s *ccb) OnPreAction(handlers ...cli.OnPreInvokeHandler) cli.CommandBuilder {
	s.SetPreActions(handlers...)
	return s
}

func (s *ccb) OnPostAction(handlers ...cli.OnPostInvokeHandler) cli.CommandBuilder {
	s.SetPostActions(handlers...)
	return s
}

func (s *ccb) OnMatched(handler cli.OnCommandMatchedHandler) cli.CommandBuilder {
	s.SetOnMatched(handler)
	return s
}

func (s *ccb) OnEvaluateSubCommandsFromConfig(path ...string) cli.CommandBuilder {
	s.SetOnEvaluateSubCommandsFromConfig(path...)
	return s
}

func (s *ccb) OnEvaluateSubCommands(handler cli.OnEvaluateSubCommands) cli.CommandBuilder {
	s.SetOnEvaluateSubCommands(handler)
	return s
}

func (s *ccb) OnEvaluateSubCommandsOnce(handler cli.OnEvaluateSubCommands) cli.CommandBuilder {
	s.SetOnEvaluateSubCommandsOnce(handler)
	return s
}

func (s *ccb) OnEvaluateFlags(handler cli.OnEvaluateFlags) cli.CommandBuilder {
	s.SetOnEvaluateFlags(handler)
	return s
}

func (s *ccb) OnEvaluateFlagsOnce(handler cli.OnEvaluateFlags) cli.CommandBuilder {
	s.SetOnEvaluateFlagsOnce(handler)
	return s
}

func (s *ccb) PresetCmdLines(args ...string) cli.CommandBuilder {
	s.SetPresetCmdLines(args...)
	return s
}

func (s *ccb) IgnoreUnmatched(ignore ...bool) cli.CommandBuilder {
	i := true
	for _, v := range ignore {
		i = v
	}
	s.SetIgnoreUnmatched(i)
	return s
}

func (s *ccb) PassThruNow(enterPassThruModeRightNow ...bool) cli.CommandBuilder {
	p := true
	for _, v := range enterPassThruModeRightNow {
		p = v
	}
	s.SetPassThruNow(p)
	return s
}

func (s *ccb) InvokeProc(executablePath string) cli.CommandBuilder {
	s.SetInvokeProc(executablePath)
	return s
}

func (s *ccb) InvokeShell(commandLine string) cli.CommandBuilder {
	s.SetInvokeShell(commandLine)
	return s
}

func (s *ccb) UseShell(shellPath string) cli.CommandBuilder {
	s.SetShell(shellPath)
	return s
}

func theTitles(longTitle string, titles ...string) (lt, st string, aliases []string) {
	lt = longTitle
	switch len(titles) {
	// case 2:
	// 	st = titles[1]
	// 	fallthrough
	case 0:
		// do nothing
	case 1:
		fallthrough //nolint:gocritic
	default:
		st = titles[0]
		for i := 1; i < len(titles); i++ {
			aliases = append(aliases, titles[i])
		}
	}
	return
}
