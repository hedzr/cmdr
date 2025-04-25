package worker

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hedzr/evendeep"
	"github.com/hedzr/is"
	"github.com/hedzr/store"
	"github.com/hedzr/store/radix"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/atoa"
	"github.com/hedzr/cmdr/v2/internal/tool"
	"github.com/hedzr/cmdr/v2/pkg/logz"
	"github.com/hedzr/is/dir"
)

func (w *workerS) preProcess(ctx context.Context) (err error) {
	// w.Config.Store.Load() // from external providers: 1. consul, 2. local,
	logz.VerboseContext(ctx, "pre-processing...")
	dummyParseCtx := parseCtx{root: w.root, forceDefaultAction: w.ForceDefaultAction}

	w.preEnvSet(ctx) // setup envvars: APP, APP_NAME, etc.

	var aliasMap map[string]*cli.CmdS
	if aliasMap, err = w.linkCommands(ctx, w.root); err != nil {
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

	if err = w.postLinkCommands(ctx, w.root, aliasMap); err != nil {
		return
	}

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

func (w *workerS) linkCommands(ctx context.Context, root *cli.RootCommand) (aliasMap map[string]*cli.CmdS, err error) {
	if err = w.addBuiltinCommands(root); err == nil {
		if cx, ok := root.Cmd.(*cli.CmdS); ok {
			// second time to link cmd.root and cmd.owner fields, include builtins commands and flags now.
			cx.EnsureTree(ctx, root.App(), root)
			if err = w.xrefCommands(ctx, root, func(cc cli.Cmd, index, level int) {
				if x := cc.OnEvaluateSubCommandsFromConfig(); x != "" {
					if aliasMap == nil {
						aliasMap = make(map[string]*cli.CmdS)
					}
					if c, ok1 := cc.(*cli.CmdS); ok1 {
						aliasMap[x] = c
					}
				}
			}); err == nil {
				if err = w.commandsToStore(ctx, root); err == nil {
					logz.VerboseContext(ctx, "linkCommands() - *RootCommand linked itself")
				}
			}
		}
	}
	return
}
func (w *workerS) postLinkCommands(ctx context.Context, root *cli.RootCommand, aliasMap map[string]*cli.CmdS) (err error) {
	for k, cc := range aliasMap {
		conf := w.Store().WithPrefix(k)
		from := conf.Prefix() + "."
		cvt := evendeep.Cvt{}
		conf.Walk(from, func(path, fragment string, node radix.Node[any]) {
			logz.VerboseContext(ctx, "post-link-command", "path", path, "fragment", fragment, "value", node.Data())
			if path != from {
				line := strings.TrimSpace(cvt.String(node.Data()))
				key := path[len(from):]
				if len(line) < 1 {
					return
				}

				c := &cli.CmdS{
					BaseOpt: cli.BaseOpt{
						Long: fmt.Sprintf("alias-%s", key),
					},
				}
				c.SetName(key)
				c.SetDesc(fmt.Sprintf("%s: [%s]/%s", cc.LongTitle(), k, key))
				c.SetGroup("Alias")

				switch prefix := line[0]; prefix {
				case '!':
					c.SetInvokeProc(strings.TrimSpace(line[1:]))
					logz.VerboseContext(ctx, "invoke proc cmd", "target", c.InvokeProc())
				case '>':
					c.SetRedirectTo(strings.ReplaceAll(strings.TrimSpace(line[1:]), " ", "."))
					logz.VerboseContext(ctx, "redirectTo cmd", "target", c.RedirectTo())
				default:
					c.SetInvokeShell(strings.TrimSpace(line))
					logz.VerboseContext(ctx, "invoke shell cmd", "target", c.InvokeShell())
				}

				err = cc.AddSubCommand(c)
			}
		})
	}
	_ = root
	return
}

func (w *workerS) commandsToStore(ctx context.Context, root *cli.RootCommand) (err error) {
	if root == nil || w.Store() == nil {
		return
	}

	// get a new Store with prefix "app.cmd"
	conf := w.Store().WithPrefix(cli.CommandsStoreKey)
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
	_ = root
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
					key := ff.GetDottedPath()
					_, _ = conf.Set(key, ff.DefaultValue())
					logz.DebugContext(ctx, "linking a flag", "ff", ff, "key", key, "TG", ff.ToggleGroup())
					if ff.Negatable() {
						if items := ff.NegatableItems(); len(items) > 0 {
							// in this branch, all following codes are no effects,
							// so its could be removed safely.
							// but we kept them because it is a reference specially
							// for negatable items stlye.

							// a := strings.Split(ff.LongTitle(), ".")
							// title := fmt.Sprintf("%s.no-%s", a[0], a[1])
							title := ff.LongTitle()
							nf := ff.Owner().FlagBy(title)
							if nf == nil {
								logz.WarnContext(ctx, "FlagBy() return nil", "querying title", title)
							}
							// logz.DebugContext(ctx, "linking a negatable item", "ff", nf, "title", title)

							// key := nf.GetDottedNamePath()
							// key1 := nf.GetDottedPath()
							// _, _ = conf.Set(key1, nf.DefaultValue())
						} else {
							// in fact, this branch would never be reached.
							// because a normal negatable flag will be reset
							// once its shadowed flag created.

							title := fmt.Sprintf("no-%s", ff.LongTitle())
							nf := ff.Owner().FlagBy(title)
							key1 := nf.GetDottedPath()
							_, _ = conf.Set(key1, nf.DefaultValue())
							logz.DebugContext(ctx, "linking a shadowed negatable flag", "ff", nf, "title", title, "key.1", key, "TG", nf.ToggleGroup())
							// _, _ = conf.Set(nf.GetDottedPath(), nf.DefaultValue())
						}
					}
				}
			} else if x, ok := cc.(interface{ TransferIntoStore(store.Store, bool) }); ok {
				x.TransferIntoStore(conf, false)
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
		jsonLoader1 := &jsonLoaderS{filename: path.Join(appDir, "."+appName+".json")}
		jsonLoader2 := &jsonLoaderS{filename: path.Join(appDir, appName+".json")}
		logz.DebugContext(ctx, "use internal tiny json loader", "filename", jsonLoader1.filename)
		w.Config.Loaders = append(w.Config.Loaders, jsonLoader1, jsonLoader2)
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
