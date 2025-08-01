package builder

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync/atomic"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/conf"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

var _ cli.App = (*appS)(nil)

type appS struct {
	cli.Runner
	opts  []cli.Opt
	root  *cli.RootCommand
	args  []string
	inCmd int32
	inFlg int32
}

func (s *appS) GetRunner() cli.Runner { return s.Runner }

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

	err = s.Runner.Run(ctx, append(s.opts, opts...)...)

	// if err != nil {
	// 	s.Runner.DumpErrors(os.Stderr)
	// }

	if err == nil && !s.Runner.Error().IsEmpty() {
		err = s.Runner.Error()
	}
	return
}

func (s *appS) Name() string                    { return s.root.AppName }
func (s *appS) Version() string                 { return s.root.Version }
func (s *appS) Worker() cli.Runner              { return s.Runner }
func (s *appS) Root() *cli.RootCommand          { return s.root }
func (s *appS) Args() []string                  { return s.args }
func (s *appS) SetCancelFunc(cancelFunc func()) { s.Runner.SetCancelFunc(cancelFunc) }
func (s *appS) CancelFunc() func()              { return s.Runner.CancelFunc() }

func (s *appS) Build() {
	if sr, ok := s.Runner.(setRoot); ok {
		ctx := context.Background()
		logz.VerboseContext(ctx, "builder.appS.Build() - setRoot")
		if cx, ok := s.root.Cmd.(*cli.CmdS); ok {
			// first time to link cmd.root and cmd.owner fields
			cx.EnsureTreeAlways(ctx, s, s.root)
		}
		sr.SetRoot(s.root, s.args)
	}
}

func (s *appS) With(cb func(app cli.App)) {
	defer s.Build()
	cb(s)
}

func (s *appS) WithOpts(opts ...cli.Opt) cli.App {
	// s.opts = append(s.opts, opts...)
	for _, opt := range opts {
		// NOTE that `cli.WithArgs()` returns a Opt object,
		// which will be used to compare with `opt` here.
		if funcPtrSame[cli.Opt](opt, cli.WithArgs()) {
			var cfg cli.Config
			opt(&cfg)
			s.args = cfg.Args
		} else {
			s.opts = append(s.opts, opt)
		}
	}
	return s
}

func funcPtrSame[T any](fn1, fn2 T) bool {
	sf1 := reflect.ValueOf(fn1)
	sf2 := reflect.ValueOf(fn2)
	if sf1.Kind() != reflect.Func || sf2.Kind() != reflect.Func {
		return false
	}
	return sf1.Pointer() == sf2.Pointer()
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

	// build & publish a release:
	//   go build -ldflags="-X 'github.com/hedzr/cmdr/v2/conf.Version=${git-version}' -X 'github.com/hedzr/cmdr/v2/conf.AppName=me' -X 'github.com/your-name/your-cli/cli/consts.Version=-'"
	//
	// local run & debug (set your local consts.version to "-" to keep conf.Version unchanged):
	//   go run -ldflags="-X 'github.com/hedzr/cmdr/v2/conf.Version=${git-version}' -X 'github.com/hedzr/cmdr/v2/conf.AppName=me' -X 'github.com/your-name/your-cli/cli/consts.Version=-'"
	//
	// If you're running without ldflags, your local `const.Version` would be passed into cmdr as the final version number.
	// You can always set `const.Version` to "-" to tell cmdr ignore it. If so, cmdr loads `conf.Version` as the final version number.
	if strings.Trim(strings.TrimSpace(version), "-_.+") != "" {
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

func (s *appS) NewCommandBuilder(longTitle string, titles ...string) cli.CommandBuilder {
	return s.Cmd(longTitle, titles...)
}

func (s *appS) NewFlagBuilder(longTitle string, titles ...string) cli.FlagBuilder {
	return s.Flg(longTitle, titles...)
}

func (s *appS) FromStruct(structValue any, opts ...cli.StructBuilderOpt) cli.StructBuilder {
	atomic.AddInt32(&s.inCmd, 1)
	return newStructBuilderShort(s, structValue, opts...)
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

func (s *appS) ToggleableFlags(toggleGroupName string, items ...cli.BatchToggleFlag) {
	from := s.root.Cmd.(*cli.CmdS)
	b := asCommandBuilder(from, s, "")
	b.ToggleableFlags(toggleGroupName, items...)
	b.Build()
}

func (s *appS) RootBuilder(cb func(b cli.CommandBuilder)) cli.App {
	from := s.root.Cmd.(*cli.CmdS)
	atomic.AddInt32(&s.inCmd, 1)
	b := asCommandBuilder(from, s, "")
	cb(b)
	b.Build()
	// atomic.AddInt32(&s.inCmd, 1)
	return s
}

func (s *appS) OnAction(handler cli.OnInvokeHandler) cli.App {
	s.RootBuilder(func(b cli.CommandBuilder) {
		b.OnAction(handler)
	})
	return s
}

// func (s *ccb) OnAction(handler cli.OnInvokeHandler) cli.CommandBuilder {
// 	s.SetAction(handler)
// 	return s
// }
//
// func (s *ccb) OnPreAction(handlers ...cli.OnPreInvokeHandler) cli.CommandBuilder {
// 	s.SetPreActions(handlers...)
// 	return s
// }
//
// func (s *ccb) OnPostAction(handlers ...cli.OnPostInvokeHandler) cli.CommandBuilder {
// 	s.SetPostActions(handlers...)
// 	return s
// }
//
// func (s *ccb) OnMatched(handler cli.OnCommandMatchedHandler) cli.CommandBuilder {
// 	s.SetOnMatched(handler)
// 	return s
// }

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
	if s.root != nil && s.root.Cmd != nil {
		if rc, ok := s.root.Cmd.(*cli.CmdS); ok {
			if rc == child {
				atomic.AddInt32(&s.inCmd, -1)
				return
			}

			if child != nil && isAssumedAsRootCmd(child.Long) {
				s.root.Cmd = child
				atomic.CompareAndSwapInt32(&s.inCmd, 1, 0)
				return
			}
		}
	}

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
