// Copyright © 2022 Hedzr Yeh.

package builder

import (
	"sync/atomic"

	"github.com/hedzr/cmdr/v2/cli"
)

type buildable interface {
	Build()
}

// func NewCommandBuilder(parent *cli.Command, longTitle string, titles ...string) *ccb {
// 	// s := &ccb{
// 	// 	nil, parent,
// 	// 	new(cli.Command),
// 	// 	false, false,
// 	// }
// 	// s.Long, s.Short, s.Aliases = theTitles(longTitle, titles...)
// 	// return s
// 	return newCommandBuilderFrom(parent, nil, longTitle, titles...)
// }

func newCommandBuilderShort(b buildable, longTitle string, titles ...string) *ccb {
	return newCommandBuilderFrom(new(cli.Command), b, longTitle, titles...)
}

func newCommandBuilderFrom(from *cli.Command, b buildable, longTitle string, titles ...string) *ccb {
	s := &ccb{
		b, from,
		new(cli.Command),
		0, 0,
	}
	s.Long, s.Short, s.Aliases = theTitles(longTitle, titles...)
	return s
}

type ccb struct {
	buildable
	parent *cli.Command
	*cli.Command
	inCmd int32
	inFlg int32
}

func (s *ccb) Build() {
	if a, ok := s.buildable.(adder); ok {
		a.addCommand(s.Command)
	}
	// if s.parent != nil {
	// 	s.parent.AddSubCommand(s.Command)
	// 	// if a, ok := s.b.(adder); ok {
	// 	// 	a.addCommand(s.Command)
	// 	// }
	// }

	atomic.StoreInt32(&s.inCmd, 0)
	atomic.StoreInt32(&s.inFlg, 0)
}

// addCommand adds a in-building Cmd into current Command as a child-/sub-command.
// used by adder when ccb.Build.
func (s *ccb) addCommand(child *cli.Command) {
	atomic.AddInt32(&s.inCmd, -1) // reset increased inCmd at AddCmd or Cmd
	s.AddSubCommand(child)
}

// addFlag adds a in-building Flg into current Command as its flag.
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

func (s *ccb) AddCmd(cb func(b cli.CommandBuilder)) cli.CommandBuilder {
	// if atomic.LoadInt32(&s.inCmd) != 0 {
	// 	panic("cannot call AddCmd() without Build() last Cmd()/AddCmd()")
	// }

	bc := newCommandBuilderShort(s, "new-command")
	defer bc.Build() // `Build' will add `bc'(Command) to s.Command as its SubCommand
	cb(bc)
	atomic.AddInt32(&s.inCmd, 1)
	return s
}

func (s *ccb) AddFlg(cb func(b cli.FlagBuilder)) cli.CommandBuilder {
	// if atomic.LoadInt32(&s.inFlg) != 0 {
	// 	panic("cannot call AddFlg() without Build() last Flg()/AddFlg()")
	// }

	bc := newFlagBuilderShort(s, "new-flag")
	defer bc.Build() // `Build' will add `bc'(Flag) to s.Command as its Flag
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

func (s *ccb) PresetCmdLines(args ...string) cli.CommandBuilder {
	s.SetPresetCmdLines(args...)
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
