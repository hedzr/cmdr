package worker

import (
	"context"
	"os"
	"path"
	"path/filepath"

	"github.com/hedzr/cmdr/v2/internal/tool"
	"github.com/hedzr/cmdr/v2/pkg/dir"
	"github.com/hedzr/is"
	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/atoa"
)

func (w *workerS) preProcess() (err error) {
	// w.Config.Store.Load() // from external providers: 1. consul, 2. local,

	w.preEnvSet() // setup envvars: APP, APP_NAME, etc.

	if err = w.linkCommands(w.root); err != nil {
		return
	}

	if err = w.loadLoaders(); err != nil {
		return
	}

	// w.Config.Store.Load() // from external providers: 3. env

	w.postEnvLoad()

	if is.VerboseBuild() || is.DebugBuild() { // `-tags delve` or `-tags verbose` used in building?
		logz.Verbose("Dumping Store ...")
		logz.Verbose(w.Store().Dump())
	}
	return
}

func (w *workerS) preEnvSet() {
	// NOTE 'CMDR_VERSION' has been setup.

	_ = os.Setenv("APP", w.Name())
	_ = os.Setenv("APPNAME", w.Name())
	_ = os.Setenv("APP_NAME", w.Name())
	_ = os.Setenv("APP_VER", w.Version())
	_ = os.Setenv("APP_VERSION", w.Version())

	logz.Verbose("appName", "appName", w.Name(), "appVer", w.Version())

	xdgPrefer := true
	if w.Store().Has("app.force-xdg-dir-prefer") {
		xdgPrefer = w.Store().MustBool("app.force-xdg-dir-prefer", false)
		logz.Verbose("reset force-xdg-dir-prefer from store value", "app.force-xdg-dir-prefer", xdgPrefer)
	}

	home := tool.HomeDir()
	if os.Getenv("HOME") == "" {
		_ = os.Setenv("HOME", home)
	} else {
		home = os.Getenv("HOME")
	}

	dir := ""
	if os.Getenv("CONFIG_DIR") == "" { // XDG_CONFIG_DIR
		if xdgPrefer {
			dir = filepath.Join(home, ".config", w.root.AppName)
		} else {
			dir, _ = os.UserConfigDir()
		}
		_ = os.Setenv("CONFIG_DIR", dir)
	}
	if os.Getenv("CACHE_DIR") == "" {
		if xdgPrefer {
			dir = filepath.Join(home, ".cache", w.root.AppName)
		} else {
			dir, _ = os.UserCacheDir()
		}
		_ = os.Setenv("CACHE_DIR", dir)
	}
	if os.Getenv("COMMON_SHARE_DIR") == "" {
		dir = filepath.Join("/usr", "local", "share", w.root.AppName)
		_ = os.Setenv("COMMON_SHARE_DIR", dir)
	}
	if os.Getenv("LOCAL_SHARE_DIR") == "" {
		dir = filepath.Join(home, ".local", "share", w.root.AppName)
		_ = os.Setenv("LOCAL_SHARE_DIR", dir)
	}
	if os.Getenv("TEMP_DIR") == "" {
		dir := os.TempDir()
		_ = os.Setenv("TEMP_DIR", dir)
	}
}

func (w *workerS) postEnvLoad() {
	// logz.Info("env", "FORCE_DEFAULT_ACTION", os.Getenv("FORCE_DEFAULT_ACTION"))
	if w.Store().Has("app.force-default-action") {
		w.ForceDefaultAction = w.Store().MustBool("app.force-default-action", false)
		logz.Verbose("reset forceDefaultAction from store value", "ForceDefaultAction", w.ForceDefaultAction)
	}
	if tool.ToBool(os.Getenv("FORCE_DEFAULT_ACTION"), false) {
		w.ForceDefaultAction = true
		logz.Info("set ForceDefaultAction true", "ForceDefaultAction", w.ForceDefaultAction)
	}
}

func (w *workerS) linkCommands(root *cli.RootCommand) (err error) {
	if err = w.addBuiltinCommands(root); err == nil {
		root.EnsureTree(root.App(), root) // link the added builtin commands too
		if err = w.xrefCommands(root); err == nil {
			if err = w.commandsToStore(root); err == nil {
				logz.Verbose("*RootCommand linked itself")
			}
		}
	}
	return
}

func (w *workerS) commandsToStore(root *cli.RootCommand) (err error) {
	if root == nil || w.Store() == nil {
		return
	}

	conf := w.Store().WithPrefixReplaced("app.cmd")

	if conf == nil {
		return
	}

	err = w.commandsToStoreChild(root, conf)
	return
}

func (w *workerS) commandsToStoreChild(root *cli.RootCommand, conf store.Store) (err error) { //nolint:revive
	fromString := func(text string, meme any) (value any) { //nolint:revive
		var err error
		value, err = atoa.Parse(text, meme)
		if err != nil {
			value = text
		}
		return
	}

	conf.WithinLoading(func() {
		root.WalkEverything(func(cc, pp *cli.Command, ff *cli.Flag, cmdIndex, flgIndex, level int) { //nolint:revive
			if ff != nil {
				if evs := ff.EnvVars(); len(evs) > 0 {
					for _, ev := range evs {
						if v, has := os.LookupEnv(ev); has {
							data := fromString(v, ff.DefaultValue())
							ff.SetDefaultValue(data)
						}
					}
				}
				if conf != nil {
					conf.Set(ff.GetDottedPath(), ff.DefaultValue())
				}
			}
		})
	})
	return
}

// loadLoaders try to load the external loaders, for loading the config files.
func (w *workerS) loadLoaders() (err error) {
	w.precheckLoaders()

	// By default, we try loading `$(pwd)/.appName.json'.
	// The main reason is the feature doesn't take new dependence to another 3rd-party lib.
	// For cmdr/v2, we restrict to go builtins, google, and ours libraries.
	// And, ours libraries will not import any others except go builtins and google's.
	if len(w.Config.Loaders) == 0 {
		appDir := dir.GetCurrentDir()
		appName := w.Name()
		jsonLoader := &jsonLoaderS{}
		jsonLoader.filename = path.Join(appDir, "."+appName+".json")
		logz.Debug("use internal tiny json loader", "filename", jsonLoader.filename)
		w.Config.Loaders = append(w.Config.Loaders, jsonLoader)
		return
	}

	for _, loader := range w.Config.Loaders {
		if loader != nil {
			if err = loader.Load(w.root.App()); err != nil {
				break
			}
		}
	}
	return
}

func (w *workerS) precheckLoaders() {
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
			logz.Warn("Config file has been ignored", "config-file", w.configFile)
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

func (w *workerS) xrefCommands(cmd *cli.RootCommand, cb ...func(cc *cli.Command, index, level int)) (err error) { //nolint:unparam
	cmd.EnsureXref(cb...)
	return
}

// func (w *workerS) walk(cmd *cli.Command, cb func(cc *cli.Command, ff *cli.Flag)) {
// 	w.walkR(cmd, cb, 0)
// }
//
// func (w *workerS) walkR(cmd *cli.Command, cb func(cc *cli.Command, ff *cli.Flag), level int) {
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

func (w *workerS) postProcess(ctx *parseCtx) (err error) {
	_ = ctx
	return
}
