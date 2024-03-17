package worker

import (
	"context"
	"os"

	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/atoa"
	"github.com/hedzr/cmdr/v2/internal/tool"
	"github.com/hedzr/store"
)

func (w *workerS) preProcess() (err error) {
	// w.Config.Store.Load() // from external providers: 1. consul, 2. local,

	w.preEnvSet()

	if err = w.linkCommands(w.root); err != nil {
		return
	}

	if err = w.loadLoaders(); err != nil {
		return
	}

	// w.Config.Store.Load() // from external providers: 3. env

	// logz.Info("env", "FORCE_DEFAULT_ACTION", os.Getenv("FORCE_DEFAULT_ACTION"))
	if tool.ToBool(os.Getenv("FORCE_DEFAULT_ACTION"), false) {
		w.ForceDefaultAction = true
		logz.Info("ForceDefaultAction set true.", "FORCE_DEFAULT_ACTION", os.Getenv("FORCE_DEFAULT_ACTION"))
	}

	logz.Verbose("Store.Dump")
	logz.Verbose(w.Store().Dump())

	return
}

func (w *workerS) preEnvSet() {
	_ = os.Setenv("APP", w.root.AppName)
	_ = os.Setenv("APPNAME", w.root.AppName)
	_ = os.Setenv("APP_NAME", w.root.AppName)
	_ = os.Setenv("APP_VER", w.root.AppVersion())
	_ = os.Setenv("APP_VERSION", w.root.AppVersion())

	logz.Verbose("appName", "appName", w.root.AppName, w.root.AppVersion())

	if os.Getenv("HOME") == "" {
		home, _ := os.UserHomeDir()
		_ = os.Setenv("HOME", home)
	}
	if os.Getenv("CACHE_DIR") == "" {
		dir, _ := os.UserCacheDir()
		_ = os.Setenv("CACHE_DIR", dir)
	}
	if os.Getenv("CONFIG_DIR") == "" {
		dir, _ := os.UserConfigDir()
		_ = os.Setenv("CONFIG_DIR", dir)
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
	if w.Store() == nil || root == nil {
		return
	}

	conf := w.Store().WithPrefixReplaced("app.cmd")

	if conf == nil {
		return
	}

	fromString := func(text string, meme any) (value any) { //nolint:revive
		var err error
		value, err = atoa.Parse(text, meme)
		if err != nil {
			value = text
		}
		return
	}

	conf.WithinLoading(func() {
		root.WalkEverything(func(cc, pp *cli.Command, ff *cli.Flag, cmdIndex, flgIndex, level int) {
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

	for _, loader := range w.Config.Loaders {
		if loader != nil {
			err = loader.Load(w.root.App())
			if err != nil {
				break
			}
		}
	}
	return
}

// writeBackToLoaders implements write-back mechanism:
// At the end of app terminated, the modified Store entries will be written back to "alternative config".
func (w *workerS) writeBackToLoaders(ctx context.Context) (err error) {
	for _, loader := range w.Config.Loaders {
		if loader != nil {
			// see also (*conffileloader).Save(ctx) and file provider, and (*loadS).Save() and trySave()
			if x, ok := loader.(store.Writeable); ok {
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

func (w *workerS) postProcess() (err error) {
	return
}
