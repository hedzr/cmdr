package builder

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/conf"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

type appS struct {
	cli.Runner
	root  *cli.RootCommand
	args  []string
	inCmd int32
	inFlg int32
}

func (s *appS) Run(ctx context.Context, opts ...cli.Opt) (err error) {
	if atomic.LoadInt32(&s.inCmd) != 0 {
		return errors.New("app/rootCmd: a Cmd() call needs ending with Build()")
	}
	if atomic.LoadInt32(&s.inFlg) != 0 {
		return errors.New("a NewFlagBuilder()/Flg() call needs ending with Build()")
	}

	if s.root == nil || s.root.Cmd == nil {
		return cli.ErrEmptyRootCommand
	}

	// Any structural errors causes panic rather than
	// returning an error object directly.
	// s.Build() will transfer s.root and s.args into
	// runner object (s.Runner) which shall support
	// `setRoot' interface.
	s.Build() // set rootCommand into worker

	s.Runner.InitGlobally(ctx) // let worker starts initializations

	if !s.Runner.Ready() {
		return cli.ErrCommandsNotReady
	}

	err = s.Runner.Run(ctx, opts...)

	// if err != nil {
	// 	s.Runner.DumpErrors(os.Stderr)
	// }

	if err == nil && !s.Runner.Error().IsEmpty() {
		err = s.Runner.Error()
	}
	return
}

func (s *appS) Name() string           { return s.root.AppName }
func (s *appS) Version() string        { return s.root.Version }
func (s *appS) Worker() cli.Runner     { return s.Runner }
func (s *appS) Root() *cli.RootCommand { return s.root }
func (s *appS) Args() []string         { return s.args }

func (s *appS) Build() {
	if sr, ok := s.Runner.(setRoot); ok {
		ctx := context.Background()
		logz.VerboseContext(ctx, "builder.appS.Build() - setRoot")
		if cx, ok := s.root.Cmd.(*cli.CmdS); ok {
			cx.EnsureTree(ctx, s, s.root)
		}
		sr.SetRoot(s.root, s.args)
	}
}

func (s *appS) ensureNewApp() cli.App { //nolint:unparam
	if s.root == nil {
		s.root = &cli.RootCommand{
			AppName: conf.AppName,
			Version: conf.Version,
			// Copyright:  "",
			// Author:     "",
			// HeaderLine: "",
			// FooterLine: "",
			// CmdS:    nil,
		}
	}
	if s.root.Cmd == nil {
		s.root.Cmd = new(cli.CmdS)
		s.root.Cmd.SetName(s.root.AppName)
	}
	return s
}

// func (s *appS) Store() store.Store {
// 	return s.Runner.Store()
// }

func (s *appS) Info(name, version string, desc ...string) cli.App {
	s.ensureNewApp()
	if name != "" {
		s.root.AppName = name
		if name != conf.AppName {
			conf.AppName = name
		}
	}
	if version != "" {
		s.root.Version = version
		if version != conf.Version {
			conf.Version = version
		}
	}
	if cx, ok := s.root.Cmd.(*cli.CmdS); ok {
		cx.SetDescription("", desc...)
	}
	return s
}

func (s *appS) Examples(examples ...string) cli.App {
	s.ensureNewApp()
	if cx, ok := s.root.Cmd.(*cli.CmdS); ok {
		cx.SetExamples(examples...)
	}
	return s
}

func (s *appS) Copyright(copyright string) cli.App {
	s.ensureNewApp()
	s.root.Copyright = copyright
	return s
}

func (s *appS) Author(author string) cli.App {
	s.ensureNewApp()
	s.root.Author = author
	return s
}

func (s *appS) Description(desc string) cli.App {
	s.ensureNewApp()
	s.root.SetDesc(desc)
	return s
}

func (s *appS) Header(headerLine string) cli.App {
	s.ensureNewApp()
	s.root.HeaderLine = headerLine
	return s
}

func (s *appS) Footer(footerLine string) cli.App {
	s.ensureNewApp()
	s.root.FooterLine = footerLine
	return s
}

func (s *appS) SetRootCommand(root *cli.RootCommand) cli.App {
	s.root = root
	return s
}

func (s *appS) WithRootCommand(cb func(root *cli.RootCommand)) cli.App {
	cb(s.root)
	return s
}

func (s *appS) RootCommand() *cli.RootCommand { return s.root }

func (s *appS) With(cb func(app cli.App)) {
	defer s.Build()
	cb(s)
}

func (s *appS) NewCommandBuilder(longTitle string, titles ...string) cli.CommandBuilder {
	return s.Cmd(longTitle, titles...)
}

func (s *appS) NewFlagBuilder(longTitle string, titles ...string) cli.FlagBuilder {
	return s.Flg(longTitle, titles...)
}

func (s *appS) Cmd(longTitle string, titles ...string) cli.CommandBuilder {
	atomic.AddInt32(&s.inCmd, 1)
	return newCommandBuilderShort(s, longTitle, titles...)
}

func (s *appS) Flg(longTitle string, titles ...string) cli.FlagBuilder {
	atomic.AddInt32(&s.inFlg, 1)
	return newFlagBuilderShort(s, longTitle, titles...)
}

func (s *appS) NewCmdFrom(from *cli.CmdS, cb func(b cli.CommandBuilder)) cli.App {
	b := newCommandBuilderFrom(from, s, "")
	defer b.Build()
	cb(b)
	return s
}

func (s *appS) NewFlgFrom(from *cli.CmdS, defaultValue any, cb func(b cli.FlagBuilder)) cli.App {
	b := newFlagBuilderFrom(from, s, defaultValue, "")
	defer b.Build()
	cb(b)
	return s
}

func (s *appS) AddCmd(cb func(b cli.CommandBuilder)) cli.App {
	b := newCommandBuilderShort(s, "")
	defer b.Build()
	cb(b)
	return s
}

func (s *appS) AddFlg(cb func(b cli.FlagBuilder)) cli.App {
	b := newFlagBuilderShort(s, "")
	defer b.Build()
	cb(b)
	return s
}

func (s *appS) addCommand(child *cli.CmdS) {
	atomic.AddInt32(&s.inCmd, -1)
	if s.root == nil {
		s.root = &cli.RootCommand{Cmd: child}
	} else {
		if cx, ok := s.root.Cmd.(*cli.CmdS); ok {
			cx.AddSubCommand(child)
		}
	}
}

func (s *appS) addFlag(child *cli.Flag) {
	atomic.AddInt32(&s.inFlg, -1)
	if cx, ok := s.root.Cmd.(*cli.CmdS); ok {
		cx.AddFlag(child)
	}
}
