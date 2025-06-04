package worker

import (
	"github.com/hedzr/cmdr/v2/cli"
)

func rootCmdForTesting() (root *cli.RootCommand) { //nolint:unused //for test
	return nil
}

//

//

//

// func newTestRunner() cli.Runner {
// 	return &workerS{
// 		Config: &cli.Config{
// 			Store: store.New(),
// 		},
// 	}
// }

// func newTestRunner() cli.Runner {
// 	return &workerS{store: store.New()}
// }
//
// // workerS for testing only
// type workerS struct {
// 	store   store.Store
// 	retCode int
// }
//
// func (w *workerS) SetSuggestRetCode(ret int) {
// 	w.retCode = ret
// }
//
// func (w *workerS) InitGlobally(ctx context.Context) {}
// func (w *workerS) Ready() bool                      { return true }
// func (w *workerS) DumpErrors(wr io.Writer)          {}             //nolint:revive
// func (w *workerS) Error() errorsv3.Error            { return nil } //nolint:revive
// func (w *workerS) Recycle(errs ...error)            {}             //
// func (w *workerS) Store(prefix ...string) store.Store {
// 	if len(prefix) > 0 {
// 		return w.store.WithPrefix(prefix...)
// 	}
// 	return w.store
// }
// func (w *workerS) Run(ctx context.Context, opts ...Opt) (err error) { return }            //nolint:revive
// func (w *workerS) Actions() (ret map[string]bool)                   { return }            //nolint:revive
// func (w *workerS) Name() string                                     { return "for-test" } //
// func (*workerS) Version() string                                    { return "v0.0.0" }
// func (*workerS) Root() *RootCommand                                 { return nil }
// func (*workerS) Args() []string                                     { return nil }       //
// func (w *workerS) SuggestRetCode() int                              { return w.retCode } //
// func (w *workerS) ParsedState() ParsedState                         { return nil }
// func (w *workerS) LoadedSources() (results []LoadedSources)         { return }
//
// func (w *workerS) DoBuiltinAction(ctx context.Context, action ActionEnum, args ...any) (handled bool, err error) {
// 	return
// }

//

//

//

// // appS for testing only
// type appS struct {
// 	Runner
// 	root  *RootCommand
// 	args  []string
// 	inCmd bool
// 	inFlg bool
// }
//
// func (s *appS) GetRunner() Runner { return s.Runner }
//
// func (s *appS) NewCommandBuilder(longTitle string, titles ...string) CommandBuilder {
// 	return s.Cmd(longTitle, titles...)
// }
//
// func (s *appS) NewFlagBuilder(longTitle string, titles ...string) FlagBuilder {
// 	return s.Flg(longTitle, titles...)
// }
//
// func (s *appS) Cmd(longTitle string, titles ...string) CommandBuilder { //nolint:revive
// 	s.inCmd = true
// 	// return newCommandBuilder(s, longTitle, titles...)
// 	return nil
// }
//
// func (s *appS) With(cb func(app App)) { //nolint:revive
// 	cb(s)
// }
//
// func (s *appS) WithOpts(opts ...Opt) App {
// 	return s
// }
//
// func (s *appS) Flg(longTitle string, titles ...string) FlagBuilder { //nolint:revive
// 	s.inFlg = true
// 	// return newFlagBuilder(s, longTitle, titles...)
// 	return nil
// }
//
// func (s *appS) AddCmd(f func(b CommandBuilder)) App { //nolint:revive
// 	// b := newCommandBuilder(s, "")
// 	// defer b.Build()
// 	// cb(b)
// 	return s
// }
//
// func (s *appS) AddFlg(cb func(b FlagBuilder)) App { //nolint:revive
// 	// b := newFlagBuilder(s, "")
// 	// defer b.Build()
// 	// cb(b)
// 	return s
// }
//
// func (s *appS) NewCmdFrom(from *CmdS, cb func(b CommandBuilder)) App { //nolint:revive
// 	// b := newCommandBuilderFrom(from, s, "")
// 	// defer b.Build()
// 	// cb(b)
// 	return s
// }
//
// func (s *appS) NewFlgFrom(from *CmdS, defaultValue any, cb func(b FlagBuilder)) App { //nolint:revive
// 	// b := newFlagBuilderFrom(from, s, defaultValue, "")
// 	// defer b.Build()
// 	// cb(b)
// 	return s
// }
//
// func (s *appS) ToggleableFlags(toggleGroupName string, items ...BatchToggleFlag) {}
//
// func (s *appS) RootBuilder(cb func(b CommandBuilder)) App { return s }
//
// func (s *appS) addCommand(child *CmdS) { //nolint:unused
// 	s.inCmd = false
// 	if cx, ok := s.root.Cmd.(*CmdS); ok {
// 		cx.AddSubCommand(child)
// 	}
// }
//
// func (s *appS) addFlag(child *Flag) { //nolint:unused
// 	s.inFlg = false
// 	if cx, ok := s.root.Cmd.(*CmdS); ok {
// 		cx.AddFlag(child)
// 	}
// }
//
// func (s *appS) Info(name, version string, desc ...string) App {
// 	s.ensureNewApp()
// 	if s.root.AppName == "" {
// 		s.root.AppName = name
// 	}
// 	if s.root.Version == "" {
// 		s.root.Version = version
// 	}
// 	if cx, ok := s.root.Cmd.(*CmdS); ok {
// 		cx.SetDescription("", desc...)
// 	}
// 	return s
// }
//
// func (s *appS) Examples(examples ...string) App {
// 	s.ensureNewApp()
// 	if cx, ok := s.root.Cmd.(*CmdS); ok {
// 		cx.SetExamples(examples...)
// 	}
// 	return s
// }
//
// func (s *appS) Copyright(copyright string) App {
// 	s.ensureNewApp()
// 	s.root.Copyright = copyright
// 	return s
// }
//
// func (s *appS) Author(author string) App {
// 	s.ensureNewApp()
// 	s.root.Author = author
// 	return s
// }
//
// func (s *appS) Description(desc string) App {
// 	s.ensureNewApp()
// 	s.root.SetDesc(desc)
// 	return s
// }
//
// func (s *appS) Header(headerLine string) App {
// 	s.ensureNewApp()
// 	s.root.HeaderLine = headerLine
// 	return s
// }
//
// func (s *appS) Footer(footerLine string) App {
// 	s.ensureNewApp()
// 	s.root.FooterLine = footerLine
// 	return s
// }
//
// func (s *appS) OnAction(handler OnInvokeHandler) App {
// 	return s
// }
//
// func (s *appS) SetRootCommand(root *RootCommand) App {
// 	s.root = root
// 	return s
// }
//
// func (s *appS) WithRootCommand(cb func(root *RootCommand)) App {
// 	cb(s.root)
// 	return s
// }
//
// func (s *appS) RootCommand() *RootCommand { return s.root }
//
// func (s *appS) Name() string       { return s.root.AppName }
// func (s *appS) Version() string    { return s.root.Version }
// func (s *appS) Worker() Runner     { return s.Runner }
// func (s *appS) Root() *RootCommand { return s.root }
// func (s *appS) Args() []string     { return s.args }
//
// func (s *appS) ensureNewApp() App { //nolint:unparam
// 	if s.root == nil {
// 		s.root = &RootCommand{
// 			AppName: conf.AppName,
// 			Version: conf.Version,
// 			app:     s,
// 			// Copyright:  "",
// 			// Author:     "",
// 			// HeaderLine: "",
// 			// FooterLine: "",
// 			// CmdS:    nil,
// 		}
// 	}
// 	if s.root.Cmd == nil {
// 		s.root.Cmd = new(CmdS)
// 		s.root.Cmd.SetName(s.root.AppName)
// 	}
// 	return s
// }
//
// func (s *appS) Build() {
// 	type setRoot interface {
// 		SetRoot(root *RootCommand, args []string)
// 	}
// 	if sr, ok := s.Runner.(setRoot); ok {
// 		ctx := context.Background()
// 		if cx, ok := s.root.Cmd.(*CmdS); ok {
// 			cx.EnsureTree(ctx, s, s.root)
// 		}
// 		sr.SetRoot(s.root, s.args)
// 	}
// }
//
// func (s *appS) Run(ctx context.Context, opts ...Opt) (err error) {
// 	if s.inCmd {
// 		return errors.New("a NewCommandBuilder()/Cmd() call needs ending with Build()")
// 	}
// 	if s.inFlg {
// 		return errors.New("a NewFlagBuilder()/Flg() call needs ending with Build()")
// 	}
//
// 	if s.root == nil || s.root.Cmd == nil {
// 		return errors.New("the RootCommand hasn't been built")
// 	}
//
// 	s.Build() // set rootCommand into worker
//
// 	s.Runner.InitGlobally(ctx) // let worker starts initializations
//
// 	if !s.Runner.Ready() {
// 		return errors.New("the RootCommand hasn't been built, or Init() failed. Has builder.App.Build() called? ")
// 	}
//
// 	err = s.Runner.Run(ctx, opts...)
//
// 	// if err != nil {
// 	// 	s.Runner.DumpErrors(os.Stderr)
// 	// }
//
// 	return
// }
