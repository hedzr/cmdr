package loaders

import (
	"context"
	"os"
	"path"
	"strings"

	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/hcl"
	"github.com/hedzr/store/codecs/hjson"
	"github.com/hedzr/store/codecs/json"
	"github.com/hedzr/store/codecs/nestext"
	"github.com/hedzr/store/codecs/toml"
	"github.com/hedzr/store/codecs/yaml"
	"github.com/hedzr/store/providers/file"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/dir"
)

const confSubFolderName = "conf.d"

func NewConfigFileLoader(opts ...Opt) *conffileloader {
	s := &conffileloader{confDFolderName: confSubFolderName}
	s.initOnce()
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	return s
}

type Opt func(s *conffileloader)

func WithFolderMap(m map[string][]*Item) Opt {
	return func(s *conffileloader) {
		s.folderMap = m
	}
}

func WithConfDFolderName(name string) Opt {
	return func(s *conffileloader) {
		s.confDFolderName = name
	}
}

type conffileloader struct {
	folderMap       map[string][]*Item
	suffixCodecMap  map[string]func() store.Codec
	confDFolderName string
}

type Item struct {
	// In a Folder, we try to stat() '$APP.yaml' or with another suffix.
	// But if Dot is true, '.$APP.yaml' will be stat() and loaded.
	Folder           string
	Dot              bool // prefix '.' to the filename?
	Recursive        bool // following 'conf.d' subdirectory?
	Watch            bool // enable watching routine?
	WriteBack        bool // write-back to "alternative config" file?
	hit              bool // this item is valid and the config file loaded?
	writeBackHandler writeBackHandler
}

type writeBackHandler interface {
	Save(ctx context.Context) error
}

func (w *conffileloader) Save(ctx context.Context) (err error) {
	for _, class := range []string{"primary", "secondary", "alternative"} {
		for _, str := range w.folderMap[class] {
			if str.hit && str.WriteBack && str.writeBackHandler != nil {
				err = str.writeBackHandler.Save(ctx)
			}
		}
	}
	return
}

func (w *conffileloader) Load(app cli.App) (err error) {
	// var conf = app.Store()

	for _, class := range []string{"primary", "secondary", "alternative"} {
		for _, it := range w.folderMap[class] {
			folderEx := os.ExpandEnv(it.Folder)
			logz.Verbose("loading config files from Folder", "class", class, "Folder", it.Folder, "Folder-expanded", folderEx)
			if !dir.FileExists(folderEx) {
				continue
			}

			var found bool
			found, err = w.loadAppConfig(folderEx, it, app)

			if root := path.Join(folderEx, w.confDFolderName); it.Recursive && found && dir.FileExists(root) {
				found, err = w.loadSubDir(root, app)
			}
		}
	}

	// logz.Verbose("Store.Dump")
	// logz.Verbose(conf.Dump())
	return
}

func (w *conffileloader) LoadFile(filename string, app cli.App) (err error) {
	return w.loadConfigFile(filename, path.Ext(filename), &Item{Watch: true, WriteBack: false}, app)
}

func (w *conffileloader) loadAppConfig(folderExpanded string, it *Item, app cli.App) (found bool, err error) {
	rootCmd := app.RootCommand()

	if file, _ := dir.IsRegularFile(folderExpanded); file {
		err = w.loadConfigFile(folderExpanded, path.Ext(folderExpanded), it, app)
		if err == nil {
			found = true
			logz.Verbose("config file loaded", "file", folderExpanded)
		}
		return
	}

	err = dir.ForFileMax(folderExpanded, 0, 1,
		func(depth int, dirName string, fi os.DirEntry) (stop bool, err error) {
			baseName, ext, appName := dir.Basename(fi.Name()), dir.Ext(fi.Name()), rootCmd.AppName
			if it.Dot {
				appName = "." + appName
			}
			if baseName != appName {
				return
			}

			logz.Verbose("loading config file", "dir", dirName, "file", fi.Name())
			err = w.loadConfigFile(path.Join(dirName, fi.Name()), ext, it, app)
			if err == nil {
				logz.Verbose("config file loaded", "file", path.Join(dirName, fi.Name()))
				found, stop = true, true
			}
			return
		})
	return
}

func (w *conffileloader) loadConfigFile(filename, ext string, it *Item, app cli.App) (err error) {
	logz.Verbose("loading config file", "file", filename)
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}
	if codec, ok := w.suffixCodecMap[ext]; ok {
		var wr writeBackHandler
		wr, err = app.Store().Load(context.Background(),
			// store.WithStorePrefix("app.yaml"),
			// store.WithPosition("app"),
			store.WithCodec(codec()),
			store.WithProvider(file.New(filename,
				file.WithWatchEnabled(it.Watch),
				file.WithWriteBackEnabled(it.WriteBack))),
		)
		if err == nil && it.WriteBack {
			it.writeBackHandler = wr
			it.hit = true
		}
	}
	return
}

func (w *conffileloader) loadSubDir(root string, app cli.App) (found bool, err error) {
	err = dir.ForFile(root,
		func(depth int, dirName string, fi os.DirEntry) (stop bool, err error) {
			ext := dir.Ext(fi.Name())
			if strings.HasPrefix(ext, ".") {
				ext = ext[1:]
			}
			if codec, ok := w.suffixCodecMap[ext]; ok {
				filename := path.Join(dirName, fi.Name())
				_, err = app.Store().Load(context.Background(),
					// store.WithStorePrefix("app.yaml"),
					// store.WithPosition("app"),
					store.WithCodec(codec()),
					store.WithProvider(file.New(filename)),
				)
				if err == nil {
					logz.Verbose("conf.d file loaded", "file", filename)
					found, stop = true, true
				}
			}
			return
		})
	return
}

func (w *conffileloader) SetAlternativeConfigFile(file string) {
	w.folderMap["alternative"] = []*Item{{Folder: file, Watch: true}}
}

func (w *conffileloader) initOnce() {
	if w.folderMap == nil {
		w.folderMap = map[string][]*Item{
			// Primary configs, which define the baseline of app config, are generally
			// bundled with application release.
			// App installer will dispatch primary config files to the standard directory
			// position. It's `/etc/$APP/` on linux, or `/usr/loca/etc/$app` on macOS by
			// Homebrew.
			// For debugging easier in developing, we also check `./ci/etc/$app`.
			"primary": {
				{Folder: "/etc/$APP", Recursive: true, Watch: true},
				{Folder: "/usr/local/etc/$APP", Recursive: true, Watch: true},
				{Folder: "./ci/etc/$APP", Recursive: true, Watch: true},
			},
			// Secondary configs, which may make some patches on the baseline if necessary.
			// On linux and macOS, it can be `~/.$app` or `~/.config/$app` (`XDG_CONFIG_DIR`).
			"secondary": {
				{Folder: "$HOME/.$APP", Recursive: true, Watch: true},
				{Folder: "$CONFIG_DIR/$APP", Recursive: true, Watch: true},
				{Folder: "./ci/config/$APP", Recursive: true, Watch: true},
			},
			// Alternative config, which is live config, can be read and written.
			// Application, such as cmdr-based, reads primary config on startup, and
			// patches it with secondary config, and updates these configs with
			// alternative config finally.
			// At application terminating, the changes can be written back to alternative
			// config.
			"alternative": {{Folder: ".", Dot: true, Recursive: true, Watch: true, WriteBack: true}},
		}
	}
	if w.suffixCodecMap == nil {
		w.suffixCodecMap = map[string]func() store.Codec{
			"yaml":       func() store.Codec { return yaml.New() },
			"yml":        func() store.Codec { return yaml.New() },
			"json":       func() store.Codec { return json.New() },
			"hjson":      func() store.Codec { return hjson.New() },
			"toml":       func() store.Codec { return toml.New() },
			"hcl":        func() store.Codec { return hcl.New() },
			"nestedtext": func() store.Codec { return nestext.New() },
			"txt":        func() store.Codec { return nestext.New() },
			"conf":       func() store.Codec { return nestext.New() },
			"":           func() store.Codec { return nestext.New() },
		}
	}
}
