// Copyright © 2022 Hedzr Yeh.

package cmdr

import (
	"context"
	"os"
	"strings"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/store"
	"github.com/hedzr/store/radix"

	"github.com/hedzr/cmdr/v2/builder"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/worker"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

// func NewOpt[T any](defaultValue ...T) config.Opt {
// 	return nil
// }

// New starts a new cmdr app.
//
// With the returned builder.App, you may build root and sub-commands fluently.
//
//	app := cmdr.New().
//	    Info("demo-app", "0.3.1").
//	    Author("hedzr")
//	app.Cmd("jump").
//		Description("jump command").
//		Examples(`jump example`).
//		Deprecated(`jump is a demo command`).
//		With(func(b cli.CommandBuilder) {
//			b.Hidden(false)
//			b.Cmd("to").
//				Description("to command").
//				Examples(``).
//				Deprecated(`v0.1.1`).
//				Hidden(false).
//				OnAction(func(cmd *cli.CmdS, args []string) (err error) {
//					main1()
//					return // handling command action here
//				}).
//				With(func(b cli.CommandBuilder) {
//					b.Flg("full", "f").
//						Default(false).
//						Description("full command").
//						Build()
//				})
//		})
//	app.Flg("dry-run", "n").
//	    Default(false).
//	    Build() // no matter even if you're adding the duplicated one.
//
//	// // simple run the parser of app and trigger the matched command's action
//	// _ = app.Run(
//	//     cmdr.WithForceDefaultAction(false), // true for debug in developing time
//	// )
//
//	if err := app.Run(
//	    cmdr.WithForceDefaultAction(false), // true for debug in developing time
//	); err != nil {
//	    logz.ErrorContext(ctx, "Application Error:", "err", err)
//	}
//
// After the root command and all its children are built, use app.[config.App.Run]
// to parse end-user's command-line arguments, and invoke the bound
// action on the hit subcommand.
//
// It is not necessary to attach an action onto a parent command, because
// its subcommands are the main characters - but you still can do that.
func New(opts ...cli.Opt) cli.App {
	_ = os.Setenv("CMDR_VERSION", Version)
	logz.Verbose("setup env-var at earlier time", "CMDR_VERSION", Version)
	cfg := cli.NewConfig(opts...)
	w := worker.New(cfg)
	return builder.New(w)
}

// App returns a light version of builder.Runner (a.k.a. *worker.Worker).
//
// Generally it's a unique instance in one system.
//
// It's available once New() / Exec() called, else nil.
//
// App returns a cli.Runner instance, which is different with builder.App.
//
// In most cases, App() return the exact app object (a &workerS{} instance).
// But it' not the real worker if you're requesting shared app instance.
// A shared app instance must be made by New() and valued context:
//
//	ctx := context.Background()
//	ctx = context.WithValue(ctx, "shared.cmdr.app", true)
//
//	app := cmdr.Create(appName, version, author, desc).
//		WithAdders(cmd.Commands...).
//		OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
//			fmt.Printf("app.name = %s\n", cmdr.App().Name())
//			fmt.Printf("app.unique = %v\n", cmdr.App()) // return an uncertain app object
//			app := cmd.Root().App()                     // this is the real app object associate with current RootCommand
//			fmt.Printf("app = %v\n", app)
//			return
//		}).
//		Build()
//
//	if err := app.Run(ctx); err != nil {
//		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
//		os.Exit(app.SuggestRetCode())
//	} else if rc := app.SuggestRetCode(); rc != 0 {
//		os.Exit(rc)
//	}
//
// What effect words with a shared app object?
//
// Suppose you have a standard app running, and also you're requesting
// to build a new child app instance for some special purpose, thought
// you could declare it as a shared instance so that the main app never
// be replaced with the new one.
func App() cli.Runner { return worker.UniqueWorker() }

func AppName() string            { return App().Name() }            // the app's name
func AppVersion() string         { return App().Version() }         // the app's version
func AppDescription() string     { return App().Root().Desc() }     // the app's short description
func AppDescriptionLong() string { return App().Root().DescLong() } // the app's long description

// CmdLines returns the whole command-line as space-separated slice.
func CmdLines() []string { return App().Args() }

// Error returns all errors occurred with a leading parent.
//
// Suppose a long term action raised many errors in runtime,
// these errors will be collected and bundled as one package
// so that they could be traced at the app terminating time.
func Error() errors.Error { return App().Error() }

// Recycle collects your errors into a container.
// You can retrieve the error container with Error() later.
func Recycle(errs ...error) { App().Recycle(errs...) }

// Parsed identify cmdr.v2 ended the command-line arguments
// parsing task.
func Parsed() bool                   { return App().ParsedState() != nil }            // is parsed ok?
func ParsedLastCmd() cli.Cmd         { return App().ParsedState().LastCmd() }         // the parsed last command
func ParsedCommands() []cli.Cmd      { return App().ParsedState().MatchedCommands() } // the parsed commands
func ParsedPositionalArgs() []string { return App().ParsedState().PositionalArgs() }  // the rest positional args
func ParsedState() cli.ParsedState   { return App().ParsedState() }                   // return the parsed state

// LoadedSources records all external sources loaded ok.
//
// Each type LoadedSources is a map which collect the children
// objects if a external source loaded them successfully.
//
// LoadedSources() is an array to represent all loaders.
func LoadedSources() []cli.LoadedSources { return App().LoadedSources() } // the loaded config files or other sources

// Store returns the child Store tree at location 'app.cmd'.
//
// By default, cmdr maintains all command-line subcommands and flags
// as a child tree in the associated Set ("store") internally.
//
// You can check out the flags state by querying in this child store.
//
// For example, we have a command 'server'->'start' and its
// flag 'foreground', therefore we can query the flag what if
// it was given by user's 'app server start --foreground':
//
//	fore := cmdr.Store().MustBool("server.start.foreground", false)
//	if fore {
//	   runRealMain()
//	} else {
//	   service.Start("start", runRealMain) // start the real main as a service
//	}
//
//	// second form:
//	cs := cmdr.Store("server.start")
//	fore := cs.MustBool("foreground")
//	port := cs.MustInt("port", 7893)
//
// Q: How to inspect the internal Store()?
//
// A: Running `app [any subcommands] [any options] ~~debug` will dump
// the internal Store() tree.
//
// Q: Can I list all subcommands?
//
// A: A end-user can run `app ~~tree` in the shell to list them.
// Also, `app -v ~~tree` or `app ~~tree -vvv` can get a list
// of subcommands tree, and with those builtin hidden
// commands, and with those vendor hidden commands. In
// this case, `-vvv` dumps the hidden commands and
// vendor-hidden commands.
//
// Since v2.1.16, the passing prefix parameters will be
// joint as a dottedPath with dot char.
// So `Set("a", "b", "c")` is equivelant with `Set("a.b.c")`.
func Store(prefix ...string) store.Store {
	switch len(prefix) {
	case 0:
		return Set(cli.CommandsStoreKey)
	case 1:
		return Set(cli.CommandsStoreKey, prefix[0])
	default:
		return Set(append([]string{cli.CommandsStoreKey}, prefix...)...)
	}
}

// Set returns the `app` subtree as a KVStore associated
// with current App().
//
//	conf := cmdr.Set()
//	cmdStore := conf.WithPrefix(cli.CommandsStoreKey)
//	assert(cmdrStore == cmdr.Store())
//
// Set() can be used for accessing the app settings.
//
// cmdr will load and merge all found external config
// files and sources into Set/Store.
//
// So, if you have a config file at `/etc/<app>/<app>.toml`
// with the following content,
//
//	[logging]
//	file = "/var/log/app/stdout.log"
//
// The file will be loaded to `app.logging`. Now you
// can access it with this,
//
//	println(cmdr.Set().MustString("logging.file"))
//
// NOTE:
//
// To enable external source loading mechanism, you may
// need to use `hedzr/cmdr-loader`. Please surf it
// for more detail.
//
// You can also create your own external source if absent.
//
// Since v2.1.16, the passing prefix parameters will be
// joint as a dottedPath with dot char.
// So `Set("a", "b", "c")` is equivelant with `Set("a.b.c")`.
func Set(prefix ...string) store.Store {
	if len(prefix) == 0 {
		return App().Store()
	}
	var pre = strings.Join(prefix, ".")
	if pre == "" {
		return App().Store() // .WithPrefix(cli.DefaultStoreKeyPrefix)
	}
	return App().Store().WithPrefix(pre)
}

const DefaultStorePrefix = cli.DefaultStoreKeyPrefix
const CommandsStoreKey = cli.CommandsStoreKey

// RemoveOrderedPrefix removes '[a-z0-9]+\.' at front of string.
func RemoveOrderedPrefix(s string) string {
	return cli.RemoveOrderedPrefix(s)
}

// DottedPathToCommandOrFlag searches the matched CmdS or
// Flag with the specified dotted-path.
//
// anyCmd is the starting of this searching.
// Give it a nil if you have no idea.
func DottedPathToCommandOrFlag(dottedPath string, anyCmd cli.Backtraceable) (cc cli.Backtraceable, ff *cli.Flag) {
	if anyCmd == nil {
		anyCmd = App().Root()
	}
	return cli.DottedPathToCommandOrFlag1(dottedPath, anyCmd)
}

// To finds a given path and loads the subtree into
// 'holder', typically 'holder' could be a struct.
//
// For yaml input
//
//	app:
//	  server:
//	    sites:
//	      - name: default
//	        addr: ":7999"
//	        location: ~/Downloads/w/docs
//
// The following codes can load it into sitesS struct:
//
//	var sites sitesS
//	err = cmdr.To("server.sites", &sites)
//
//	type sitesS struct{ Sites []siteS }
//
//	type siteS struct {
//	  Name        string
//	  Addr        string
//	  Location    string
//	}
//
// In this above case, 'store' loaded yaml and built it
// into memory, and extract 'server.sites' into 'sitesS'.
// Since 'server.sites' is a yaml array, it was loaded
// as a store entry and holds a slice value, so GetSectionFrom
// extract it to sitesS.Sites field.
//
// The optional MOpt operators could be:
//   - WithKeepPrefix
//   - WithFilter
func To[T any](path string, holder *T, opts ...radix.MOpt[any]) (err error) {
	return App().Store().To(path, holder, opts...)
}

// Exec starts a new cmdr app (parsing cmdline args based on the given rootCmd)
// from scratch.
//
// It's a reserved API for back-compatible with cmdr v1.
//
// It'll be removed completely at the recently future version.
//
// Deprecated since 2.1 by app.Run()
func Exec(rootCmd *cli.RootCommand, opts ...cli.Opt) (err error) {
	// if is.InDebugging() {
	// 	_ = exec.Run("/bin/false")
	// 	// cabin.Version()
	// 	// cpcn.Out()
	// }

	app := New(opts...).SetRootCommand(rootCmd)
	err = app.Run(context.Background())
	return
}
