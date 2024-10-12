package worker

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"

	"gopkg.in/hedzr/errors.v3"

	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/conf"
)

func New(cfg *cli.Config, opts ...cli.Opt) *workerS {
	w := &workerS{Config: cfg}
	w.setArgs(cfg.Args)
	bindOpts(w, false, opts...)
	return w
}

func newWorker(opts ...wOpt) *workerS {
	s := &workerS{Config: cli.DefaultConfig()}
	for _, opt := range opts {
		opt(s)
	}
	return s
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

	// app app.App

	root *cli.RootCommand
	args []string

	errs   errors.Error
	ready  int32 // rootCommand is set and ready for Run running.
	closed int32 // Run has exited, and all resources released

	configFile      string
	versionSimulate string
	debugOutputFile string
	actionsMatched  actionEnum
	strictMode      bool
	strictModeLevel int
	noLoadEnv       bool

	inCompleting bool
	actions      map[actionEnum]onAction
	parsingCtx   *parseCtx
}

func (w *workerS) String() string {
	var sb strings.Builder
	_, _ = sb.WriteString("workerS{root: \"")
	_, _ = sb.WriteString(w.Root().String())
	_, _ = sb.WriteString("\", config: ")
	if b, err := json.Marshal(w.Config); err == nil {
		_, _ = sb.Write(b)
	} else {
		logz.Error("json marshalling w.Config failed", "config", w.Config, "err", err)
	}
	_, _ = sb.WriteString("\"}")
	return sb.String()
}

type actionEnum int

const (
	actionShowVersion actionEnum = 1 << iota
	actionShowBuiltInfo
	actionShowHelpScreen
	actionShowHelpScreenAsMan
	actionShowTree  // Tree. `~~tree` | show all of commands (& flags) as a tree
	actionShowDebug // Debug. `~~debug` | show debug information for debugging cmdr internal states
	actionShowDebugEnv
	actionShowDebugMore
	actionShowDebugRaw
	actionShowDebugValueType
	actionShowSBOM
	// actionShortMode
	// actionDblTildeMode
)

func (e actionEnum) String() string {
	var sb strings.Builder
	if e&actionShowVersion != 0 {
		_, _ = sb.WriteString("- ShowVersion\n")
	}
	if e&actionShowBuiltInfo != 0 {
		_, _ = sb.WriteString("- ShowBuiltInfo\n")
	}
	if e&actionShowHelpScreen != 0 {
		_, _ = sb.WriteString("- ShowHelpScreen\n")
	}
	if e&actionShowHelpScreenAsMan != 0 {
		_, _ = sb.WriteString("- ShowHelpScreenAsMan\n")
	}
	if e&actionShowTree != 0 {
		_, _ = sb.WriteString("- ShowTree\n")
	}
	if e&actionShowDebug != 0 {
		_, _ = sb.WriteString("- ShowDebug\n")
	}
	return sb.String()
}

func (w *workerS) Actions() (ret map[string]bool) {
	ret = make(map[string]bool)
	e := w.actionsMatched
	if e&actionShowVersion != 0 {
		ret["show-version"] = true
	}
	if e&actionShowBuiltInfo != 0 {
		ret["show-built-info"] = true
	}
	if e&actionShowHelpScreen != 0 {
		ret["show-help"] = true
	}
	if e&actionShowHelpScreenAsMan != 0 {
		ret["show-help-man"] = true
	}
	if e&actionShowTree != 0 {
		ret["show-tree"] = true
	}
	if e&actionShowDebug != 0 {
		ret["show-debug"] = true
	}
	return
}

func (w *workerS) Ready() bool {
	w.actions = map[actionEnum]onAction{
		actionShowVersion:         w.showVersion,
		actionShowBuiltInfo:       w.showBuiltInfo,
		actionShowHelpScreen:      w.showHelpScreen,
		actionShowHelpScreenAsMan: w.showHelpScreenAsMan,
		actionShowTree:            w.showTree,
		actionShowDebug:           w.showDebugScreen,
		actionShowSBOM:            w.showSBOM,
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
		logz.Debug("[cmdr] Runner(*workerS) closed.")
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

func (w *workerS) InitGlobally() {
	if w.reqResourcesReady() {
		w.initGlobalResources()
	}
}

func (w *workerS) initGlobalResources() {
	defer w.triggerGlobalResourcesInitOK()
	logz.Debug("[cmdr] workerS.initGlobalResources")
}

func (w *workerS) triggerGlobalResourcesInitOK() {
	logz.Debug("[cmdr] workerS.triggerGlobalResourcesInitOK")
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

func (w *workerS) setArgs(args []string) { w.args = args }
func (w *workerS) Args() (args []string) { return w.args }

func bindOpts[Opt cli.Opt](w *workerS, installAsUnique bool, opts ...Opt) {
	for _, opt := range opts {
		opt(w.Config)
	}

	if w.HelpScreenWriter == nil {
		w.wrHelpScreen = w.HelpScreenWriter
	}
	if w.DebugScreenWriter == nil {
		w.wrDebugScreen = w.DebugScreenWriter
	}

	if installAsUnique {
		if app := UniqueWorker(); app != w {
			SetUniqueWorker(w)
		}
	}
}

func (w *workerS) Run(opts ...cli.Opt) (err error) {
	bindOpts(w, true, opts...)

	w.errs = errors.New(w.root.AppName)
	defer w.errs.Defer(&err)

	ctx := parseCtx{root: w.root, forceDefaultAction: w.ForceDefaultAction}
	dummy := func() bool { return true }

	if w.attachError(w.preProcess()) {
		return
	}
	defer func() { w.attachError(w.postProcess(&ctx)) }()

	if w.invokeTasks(&ctx, w.errs, w.Config.TasksBeforeParse...) ||
		w.attachError(w.parse(&ctx)) ||
		w.invokeTasks(&ctx, w.errs, w.Config.TasksBeforeRun...) ||
		w.attachError(w.exec(&ctx)) ||
		w.invokeTasks(&ctx, w.errs, w.Config.TasksAfterRun...) ||
		dummy() {
		// any errors occurred
		return
	}

	return
}

func (w *workerS) invokeTasks(ctx *parseCtx, errs errors.Error, tasks ...cli.Task) (ret bool) {
	for _, tsk := range tasks {
		if tsk != nil {
			if err := tsk(ctx.LastCmd(), w, ctx, ctx.PositionalArgs()); err != nil {
				ret = true
				errs.Attach(err)
			}
		}
	}
	_ = ctx
	return
}
