package worker

import (
	"context"
	"os"
	"path"
	"path/filepath"

	"github.com/hedzr/is"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/atoa"
	"github.com/hedzr/cmdr/v2/internal/tool"
	"github.com/hedzr/cmdr/v2/pkg/dir"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

func (w *workerS) preProcess(ctx context.Context) (err error) {
	// w.Config.Store.Load() // from external providers: 1. consul, 2. local,
	logz.VerboseContext(ctx, "pre-processing...")
	dummyParseCtx := parseCtx{root: w.root, forceDefaultAction: w.ForceDefaultAction}

	w.preEnvSet(ctx) // setup envvars: APP, APP_NAME, etc.

	if err = w.linkCommands(ctx, w.root); err != nil {
		return
	}

	if w.invokeTasks(ctx, &dummyParseCtx, w.errs, w.Config.TasksAfterXref...) {
		return
	}

	if err = w.loadLoaders(ctx); err != nil {
		return
	}

	if w.invokeTasks(ctx, &dummyParseCtx, w.errs, w.Config.TasksAfterLoader...) {
		return
	}

	// w.Config.Store.Load() // from external providers: 3. env

	w.postEnvLoad(ctx)

	if is.VerboseBuild() || is.DebugBuild() { // `-tags delve` or `-tags verbose` used in building?
		logz.VerboseContext(ctx, "Dumping Store ...")
		logz.VerboseContext(ctx, w.Store().Dump())
	}
	return
}

func (w *workerS) preEnvSet(ctx context.Context) {
	// NOTE 'CMDR_VERSION' has been setup.

	if w.Env != nil {
		for k, v := range w.Env {
			_ = os.Setenv(k, v)
		}
	}

	appName := w.Name()

	_ = os.Setenv("APP", appName)
	_ = os.Setenv("APPNAME", appName)
	_ = os.Setenv("APP_NAME", appName)
	_ = os.Setenv("APP_VER", w.Version())
	_ = os.Setenv("APP_VERSION", w.Version())
	_ = os.Setenv("EXE", dir.GetExecutablePath())
	_ = os.Setenv("EXE_DIR", dir.GetExecutableDir())

	logz.VerboseContext(ctx, "preEnvSet()", "appName", appName, "appVer", w.Version())

	// xdgPrefer := true
	// if w.Store().Has("app.force-xdg-dir-prefer") {
	// 	xdgPrefer = w.Store().MustBool("app.force-xdg-dir-prefer", false)
	// 	logz.VerboseContext(ctx, "reset force-xdg-dir-prefer from store value", "app.force-xdg-dir-prefer", xdgPrefer)
	// }

	home := tool.HomeDir()
	if os.Getenv("HOME") == "" {
		_ = os.Setenv("HOME", home)
	} else {
		home = os.Getenv("HOME")
	}

	dir := tool.ConfigDir(appName)
	if os.Getenv("CONFIG_DIR") == "" { // XDG_CONFIG_DIR
		_ = os.Setenv("CONFIG_DIR", dir)
	}

	dir = tool.CacheDir(appName)
	if os.Getenv("CACHE_DIR") == "" {
		_ = os.Setenv("CACHE_DIR", dir)
	}

	if os.Getenv("COMMON_SHARE_DIR") == "" {
		dir = filepath.Join("/usr", "local", "share", w.root.AppName)
		_ = os.Setenv("COMMON_SHARE_DIR", dir)
	}

	dir = tool.DataDir(appName)
	if os.Getenv("LOCAL_SHARE_DIR") == "" {
		_ = os.Setenv("LOCAL_SHARE_DIR", dir)
	}
	if os.Getenv("DATA_DIR") == "" {
		_ = os.Setenv("DATA_DIR", dir)
	}

	tmpdir := tool.TempDir(appName)
	if os.Getenv("TEMP_DIR") == "" {
		_ = os.Setenv("TEMP_DIR", tmpdir)
	}
}

func (w *workerS) postEnvLoad(ctx context.Context) {
	logz.VerboseContext(ctx, "postEnvLoad()", "FORCE_DEFAULT_ACTION", os.Getenv("FORCE_DEFAULT_ACTION"))
	if w.Store().Has("app.force-default-action") {
		w.ForceDefaultAction = w.Store().MustBool("app.force-default-action", false)
		logz.VerboseContext(ctx, "postEnvLoad() - reset forceDefaultAction from store value", "ForceDefaultAction", w.ForceDefaultAction)
	}
	if is.ToBool(os.Getenv("FORCE_DEFAULT_ACTION"), false) {
		w.ForceDefaultAction = true
		logz.InfoContext(ctx, "postEnvLoad() - set ForceDefaultAction true", "ForceDefaultAction", w.ForceDefaultAction)
	}
	if w.ForceDefaultAction {
		if envForceRun := is.ToBool(os.Getenv("FORCE_RUN")); envForceRun {
			w.ForceDefaultAction = false
			logz.InfoContext(ctx, "postEnvLoad() - reset forceDefaultAction since FORCE_RUN defined", "ForceDefaultAction", w.ForceDefaultAction)
		}
	}
}

func (w *workerS) linkCommands(ctx context.Context, root *cli.RootCommand) (err error) {
	if err = w.addBuiltinCommands(root); err == nil {
		if cx, ok := root.Cmd.(*cli.CmdS); ok {
			cx.EnsureTree(ctx, root.App(), root) // link the added builtin commands too
			if err = w.xrefCommands(ctx, root); err == nil {
				if err = w.commandsToStore(ctx, root); err == nil {
					logz.VerboseContext(ctx, "linkCommands() - *RootCommand linked itself")
				}
			}
		}
	}
	return
}

func (w *workerS) commandsToStore(ctx context.Context, root *cli.RootCommand) (err error) {
	if root == nil || w.Store() == nil {
		return
	}

	// get a new Store with prefix "app.cmd"
	conf := w.Store().WithPrefixReplaced(cli.CommandsStoreKey)
	if conf == nil {
		return
	}

	err = w.commandsToStoreR(ctx, root, conf)
	return
}

func (w *workerS) commandsToStoreR(ctx context.Context, root *cli.RootCommand, conf store.Store) (err error) { //nolint:revive
	fromString := func(text string, meme any) (value any) { //nolint:revive
		var err error
		value, err = atoa.Parse(text, meme)
		if err != nil {
			value = text
		}
		return
	}
	if w.Config.AutoEnvPrefix == "" {
		w.Config.AutoEnvPrefix = "APP"
	}
	worker := func(cx cli.Cmd) {
		// lookup all flags...
		//    and bind the value with its envvars field;
		//    also bind the auto-binding env vars;
		cx.WalkEverything(ctx, func(cc, pp cli.Cmd, ff *cli.Flag, cmdIndex, flgIndex, level int) { //nolint:revive
			if ff != nil {
				if evs := ff.EnvVars(); len(evs) > 0 {
					for _, ev := range evs {
						if v, has := os.LookupEnv(ev); has {
							data := fromString(v, ff.DefaultValue())
							ff.SetDefaultValue(data)
						}
					}
				}
				if w.Config.AutoEnv {
					ev := ff.GetAutoEnvVarName(w.Config.AutoEnvPrefix, true)
					if v, has := os.LookupEnv(ev); has {
						data := fromString(v, ff.DefaultValue())
						ff.SetDefaultValue(data)
					}
				}
				if conf != nil {
					conf.Set(ff.GetDottedPath(), ff.DefaultValue())
				}
			}
		})
	}
	if cx, ok := w.root.Cmd.(*cli.CmdS); ok && cx != nil && conf != nil {
		// using store.WithinLoading to disable onSet callbacks and reset internal modified states.
		conf.WithinLoading(func() { worker(cx) })
	} else if cx := w.root.Cmd; cx != nil {
		worker(cx)
	}
	return
}

// loadLoaders try to load the external loaders, for loading the config files.
func (w *workerS) loadLoaders(ctx context.Context) (err error) {
	w.precheckLoaders(ctx)

	// By default, we try loading `$(pwd)/.appName.json' if there
	// is no any loaders specified.
	//
	// The main reason is the feature doesn't take new dependence
	// to another 3rd-party lib.
	//
	// For cmdr/v2, we restrict to go builtins, google, and ours
	// libraries. And, ours libraries will not import any others
	// except go builtins and google's.
	if len(w.Config.Loaders) == 0 {
		appDir := dir.GetCurrentDir()
		appName := w.Name()
		jsonLoader := &jsonLoaderS{}
		jsonLoader.filename = path.Join(appDir, "."+appName+".json")
		logz.DebugContext(ctx, "use internal tiny json loader", "filename", jsonLoader.filename)
		w.Config.Loaders = append(w.Config.Loaders, jsonLoader)
	}

	for _, loader := range w.Config.Loaders {
		if loader != nil {
			if err = loader.Load(ctx, w.root.App()); err != nil {
				if _, ok := loader.(*jsonLoaderS); !ok {
					break
				}
				err = nil
			}
		}
	}
	return
}

func (w *workerS) LoadedSources() (results []cli.LoadedSources) {
	for _, loader := range w.Config.Loaders {
		if loader != nil {
			if q, ok := loader.(cli.QueryLoadedSources); ok {
				results = append(results, q.LoadedSources())
			}
		}
	}
	return
}

func (w *workerS) precheckLoaders(ctx context.Context) {
	if w.configFile != "" {
		found := false
		for _, loader := range w.Config.Loaders {
			// try calling (*conffileloader).SetAlternativeConfigFile(file)
			if x, ok := loader.(interface{ SetAlternativeConfigFile(file string) }); ok {
				x.SetAlternativeConfigFile(w.configFile)
				found = true
				break
			}
		}
		if !found {
			logz.WarnContext(ctx, "Config file has been ignored", "config-file", w.configFile)
		}
	}
}

// writeBackToLoaders implements write-back mechanism:
// At the end of app terminated, the modified Store entries will be written back to "alternative config".
func (w *workerS) writeBackToLoaders(ctx context.Context) (err error) {
	for _, loader := range w.Config.Loaders {
		if loader != nil {
			// see also (*conffileloader).Save(ctx) and file provider, and (*loadS).Save() and trySave()
			if x, ok := loader.(store.Writeable); ok && x != nil {
				err = x.Save(ctx)
				if err != nil {
					break
				}
			}
		}
	}
	return
}

func (w *workerS) xrefCommands(ctx context.Context, root *cli.RootCommand, cb ...func(cc cli.Cmd, index, level int)) (err error) { //nolint:unparam
	if cx, ok := root.Cmd.(*cli.CmdS); ok {
		cx.EnsureXref(ctx, cb...)
	}
	return
}

// func (w *workerS) walk(cmd *cli.CmdS, cb func(cc *cli.CmdS, ff *cli.Flag)) {
// 	w.walkR(cmd, cb, 0)
// }
//
// func (w *workerS) walkR(cmd *cli.CmdS, cb func(cc *cli.CmdS, ff *cli.Flag), level int) {
// 	cb(cmd, nil)
//
// 	for _, ff := range cmd.Flags() {
// 		cb(cmd, ff)
// 	}
//
// 	for _, child := range cmd.SubCommands() {
// 		w.walkR(child, cb, level+1)
// 	}
// }

func (w *workerS) postProcess(ctx context.Context, pc *parseCtx) (err error) {
	logz.VerboseContext(ctx, "post-processing...")
	_, _ = pc, ctx
	return
}
