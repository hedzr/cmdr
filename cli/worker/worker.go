package worker

import (
	"context"
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/is/basics"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/conf"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

func New(cfg *cli.Config, opts ...cli.Opt) *workerS {
	w := &workerS{Config: cfg}
	w.setArgs(cfg.Args)
	bindOpts(context.TODO(), w, false, opts...)
	return w
}

func newWorker(opts ...wOpt) *workerS {
	s := &workerS{Config: cli.DefaultConfig()}
	return s.With(opts...)
}

var _ = cli.Runner((*workerS)(nil))

var onceWorker sync.Once

var uniW atomic.Pointer[workerS]

func UniqueWorker() cli.Runner {
	onceWorker.Do(func() {
		uniW.Store(newWorker())
	})
	return uniW.Load()
}

func SetUniqueWorker(s cli.Runner) {
	if w, ok := s.(*workerS); ok {
		uniW.Store(w)
	} else {
		rv := reflect.ValueOf(s)
		if rv.Kind() != reflect.Pointer {
			panic("cannot support atomic putting a type-unknown cli,Runner object")
		}
		// don't worry about this casting to force converting your pointer to *workerS,
		// because the underlying stored value is just an unsafe pointer always,
		// and it will be retrieved out as a cli.Runner later.
		uniW.Store((*workerS)(rv.UnsafePointer()))
	}
}

// taskAfterParse handles parsing error and gives a chance to inspect it.
//
// You can inspect errParsed and return it as is.
//
// You can also inspect errParsed and return nil to ignore/disable a parsing error.
type taskAfterParse func(w *workerS, ctx *parseCtx, errParsed error) (err error)

// HelpWriter needs to be compatible with [io.Writer] and [io.StringWriter].
type HelpWriter interface {
	io.Writer
	io.StringWriter
}

type workerS struct {
	*cli.Config

	tasksAfterParse []taskAfterParse
	wrHelpScreen    HelpWriter
	wrDebugScreen   HelpWriter
	// onInterpretLeadingPlusSign cli.OnInterpretLeadingPlusSign

	// app app.App

	root *cli.RootCommand
	args []string

	retCode int
	errs    errors.Error
	ready   int32 // rootCommand is set and ready for Run running.
	closed  int32 // Run has exited, and all resources released

	configFile      string
	versionSimulate string
	debugOutputFile string
	actionsMatched  cli.ActionEnum
	strictMode      bool
	strictModeLevel int
	noLoadEnv       bool

	inCompleting bool
	actions      map[cli.ActionEnum]onAction
	parsingCtx   cli.ParsedState
}

func (w *workerS) String() string {
	var sb strings.Builder
	_, _ = sb.WriteString("workerS{root: \"")
	if root := w.Root(); root != nil {
		_, _ = sb.WriteString(root.String())
	} else {
		_, _ = sb.WriteString("(root-unset)")
	}
	_, _ = sb.WriteString("\", config: ")
	if w.Config != nil {
		if b, err := json.Marshal(w.Config); err == nil {
			_, _ = sb.Write(b)
		} else {
			ctx := context.Background()
			logz.ErrorContext(ctx, "json marshalling w.Config failed", "config", w.Config, "err", err)
		}
	} else {
		sb.WriteString("(config-unset)")
	}
	_, _ = sb.WriteString("\"}")
	return sb.String()
}

func (w *workerS) Actions() (ret map[string]bool) {
	ret = make(map[string]bool)
	e := w.actionsMatched
	if e&cli.ActionShowVersion != 0 {
		ret["show-version"] = true
	}
	if e&cli.ActionShowBuiltInfo != 0 {
		ret["show-built-info"] = true
	}
	if e&cli.ActionShowHelpScreen != 0 {
		ret["show-help"] = true
	}
	if e&cli.ActionShowHelpScreenAsMan != 0 {
		ret["show-help-man"] = true
	}
	if e&cli.ActionShowTree != 0 {
		ret["show-tree"] = true
	}
	if e&cli.ActionShowDebug != 0 {
		ret["show-debug"] = true
	}
	if e&cli.ActionDefault != 0 {
		ret["default"] = true
	}
	return
}

func (w *workerS) With(opts ...wOpt) *workerS {
	for _, opt := range opts {
		opt(w)
	}
	return w
}

type onAction func(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd, args ...any) (err error)

func (w *workerS) Ready() bool {
	w.actions = map[cli.ActionEnum]onAction{
		cli.ActionShowVersion:         w.showVersion,
		cli.ActionShowBuiltInfo:       w.showBuiltInfo,
		cli.ActionShowHelpScreen:      w.showHelpScreen,
		cli.ActionShowHelpScreenAsMan: w.showHelpScreenAsMan,
		cli.ActionShowTree:            w.showTree,
		cli.ActionShowDebug:           w.showDebugScreen,
		cli.ActionShowSBOM:            w.showSBOM,
		cli.ActionDefault:             w.onDefaultAction,
	}
	return atomic.LoadInt32(&w.ready) >= 2
}

func (w *workerS) reqRootCmdReady() (yes bool) {
	yes = atomic.CompareAndSwapInt32(&w.ready, 0, 1)
	return
}

func (w *workerS) reqResourcesReady() (yes bool) {
	yes = atomic.CompareAndSwapInt32(&w.ready, 1, 2)
	return
}

func (w *workerS) Close() {
	if atomic.CompareAndSwapInt32(&w.closed, 0, 1) {
		ctx := context.Background()
		logz.DebugContext(ctx, "Runner(*workerS) closed.")
	}
}

func (w *workerS) DumpErrors(wr io.Writer) {
	if w.errs.IsEmpty() {
		return
	}
	_, _ = wr.Write([]byte(w.errs.Error()))
}

func (w *workerS) Error() (err errors.Error) {
	err = w.errs
	return
}

func (w *workerS) Recycle(errs ...error) {
	w.errs.Attach(errs...)
}

func (w *workerS) SetRoot(root *cli.RootCommand, args []string) {
	// trigger the Ready signal
	if w.reqRootCmdReady() {
		if root != nil {
			w.root = root
		}
		if args != nil {
			w.args = args
		}
	}
}

func (w *workerS) Name() string {
	if w.root != nil {
		if n := w.root.Name(); n != "" {
			return n
		}
	}
	return conf.AppName
}

func (w *workerS) Version() string {
	if v := w.versionSimulate; v != "" {
		return v
	}
	if w.root != nil {
		if v := w.root.Version; v != "" {
			return v
		}
	}
	return conf.Version
}

func (w *workerS) Root() *cli.RootCommand { return w.root }
func (w *workerS) Store() store.Store     { return w.Config.Store }

func (w *workerS) InitGlobally(ctx context.Context) {
	if w.reqResourcesReady() {
		w.initGlobalResources()
	}
}

func (w *workerS) initGlobalResources() {
	defer w.triggerGlobalResourcesInitOK()
	ctx := context.Background()
	logz.VerboseContext(ctx, "workerS.initGlobalResources")

	// to do sth...
}

func (w *workerS) triggerGlobalResourcesInitOK() {
	// to do sth...
	ctx := context.Background()
	logz.VerboseContext(ctx, "workerS.triggerGlobalResourcesInitOK")
}

func (w *workerS) attachErrors(errs ...error) { //nolint:revive,unused
	// w.errs.Attach(errs...)
	for _, err := range errs {
		w.attachError(err)
	}
}

func (w *workerS) attachError(err error) (has bool) {
	if w.errIsSignalOrNil(err) {
		return false
	}

	if has = err != nil; has {
		if w.errIsUnmatchedArg(err) {
			return false
		}
		w.errs.Attach(err)
	}
	return
}

func (w *workerS) errIsUnmatchedArg(err error) bool {
	if err == nil {
		return false
	}
	return w.UnmatchedAsError && errors.Iss(err, cli.ErrUnmatchedCommand, cli.ErrUnmatchedCommand)
}

func (w *workerS) errIsNotUnmatchedArg(err error) bool {
	if err == nil || !w.UnmatchedAsError {
		return true
	}
	return !errors.Iss(err, cli.ErrUnmatchedCommand, cli.ErrUnmatchedCommand)
}

func (w *workerS) errIsSignalOrNil(err error) bool {
	if err == nil {
		return true
	}
	return errors.Iss(err, cli.ErrShouldFallback, cli.ErrShouldStop)
}

func (w *workerS) errIsSignalFallback(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, cli.ErrShouldFallback)
}

func (w *workerS) setArgs(args []string)     { w.args = args }
func (w *workerS) Args() (args []string)     { return w.args }
func (w *workerS) SuggestRetCode() int       { return w.retCode } //
func (w *workerS) SetSuggestRetCode(ret int) { w.retCode = ret }
func (w *workerS) ParsedState() cli.ParsedState {
	if w != nil {
		return w.parsingCtx
	}
	return nil
}

func bindOpts[Opt cli.Opt](ctx context.Context, w *workerS, installAsUnique bool, opts ...Opt) {
	for _, opt := range opts {
		opt(w.Config)
	}

	if w.HelpScreenWriter != nil {
		w.wrHelpScreen = w.HelpScreenWriter
	}
	if w.DebugScreenWriter != nil {
		w.wrDebugScreen = w.DebugScreenWriter
	}

	// if cx, ok := w.root.Cmd.(*cli.CmdS); ok {
	// 	cx.EnsureTree(ctx, w, w.root)
	// }

	// update args with w.Config.Args
	if len(w.Config.Args) > 0 {
		w.args = w.Config.Args
	}

	if installAsUnique {
		if app := UniqueWorker(); app != w {
			SetUniqueWorker(w)
		}
	}
}

func (w *workerS) Run(ctx context.Context, opts ...cli.Opt) (err error) {
	bindOpts(ctx, w, true, opts...)

	// shutdown basics.Closers for the registered Peripheral, Closers.
	// See also: basics.RegisterPeripheral, basics.RegisterClosable,
	// basics.RegisterCloseFns, basics.RegisterCloseFn, and
	// basics.RegisterClosers
	defer basics.Close()

	w.errs = errors.New(w.root.AppName)
	defer w.errs.Defer(&err)

	if err = w.preProcess(ctx); err != nil {
		w.attachErrors(err)
		return
	}

	dummy := func() bool { return true }
	pc := &parseCtx{argsPtr: &w.args, root: w.root, forceDefaultAction: w.ForceDefaultAction}
	defer func() { w.attachError(w.postProcess(ctx, pc)) }()
	if w.invokeTasks(ctx, pc, w.errs, w.Config.TasksBeforeParse...) ||
		w.attachError(w.parse(ctx, pc)) ||
		w.invokeTasks(ctx, pc, w.errs, w.Config.TasksParsed...) ||
		w.invokeTasks(ctx, pc, w.errs, w.Config.TasksBeforeRun...) ||
		w.attachError(w.exec(ctx, pc)) ||
		w.invokeTasks(ctx, pc, w.errs, w.Config.TasksAfterRun...) ||
		w.invokeTasks(ctx, pc, w.errs, w.Config.TasksPostCleanup...) ||
		dummy() {
		// any errors occurred
		return
	}

	return
}

// invokeTasks returns true to identify there are some tasks handled and errors occurred.
// it returns false if no errors occurred or no any tasks handled.
func (w *workerS) invokeTasks(ctx context.Context, pc *parseCtx, errs errors.Error, tasks ...cli.Task) (ret bool) {
	for _, tsk := range tasks {
		if tsk != nil {
			if err := tsk(ctx, pc.LastCmd(), w, pc, pc.PositionalArgs()); err != nil {
				ret = true
				errs.Attach(err)
			}
		}
	}
	return
}
