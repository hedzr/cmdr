// Copyright © 2022 Hedzr Yeh.

package cmdr

import (
	"context"
	"os"

	"github.com/hedzr/is/basics"
	"github.com/hedzr/store"

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
func App() cli.Runner { return worker.UniqueWorker() }

func AppName() string            { return App().Name() }            // the app's name
func AppVersion() string         { return App().Version() }         // the app's version
func AppDescription() string     { return App().Root().Desc() }     // the app's short description
func AppDescriptionLong() string { return App().Root().DescLong() } // the app's long description

// CmdLines returns the whole command-line as space-separated slice.
func CmdLines() []string { return worker.UniqueWorker().Args() }

func Parsed() bool                   { return App().ParsedState() != nil }            // is parsed ok?
func ParsedLastCmd() cli.Cmd         { return App().ParsedState().LastCmd() }         // the parsed last command
func ParsedCommands() []cli.Cmd      { return App().ParsedState().MatchedCommands() } // the parsed commands
func ParsedPositionalArgs() []string { return App().ParsedState().PositionalArgs() }  // the rest positional args

func LoadedSources() []cli.LoadedSources { return App().LoadedSources() } // the loaded config files or other sources

// Set returns the KVStore associated with current App().
func Set() store.Store { return App().Store() }

// Store returns the child Store at location 'app.cmd'.
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
//	cs := cmdr.Store().WithPrefix("server.start")
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
// A: Running `app ~~tree`, `app -v ~~tree` or `app ~~tree -vvv` can get
// a list of subcommands tree, and with those builtin hidden commands,
// and with those vendor hidden commands. In this case, `-vvv` dumps the
// hidden commands and vendor-hidden commands.
func Store() store.Store { return Set().WithPrefix(cli.CommandsStoreKey) }

// RemoveOrderedPrefix removes '[a-z0-9]+\.' at front of string.
func RemoveOrderedPrefix(s string) string {
	return cli.RemoveOrderedPrefix(s)
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

func WithForceDefaultAction(b bool) cli.Opt {
	return func(s *cli.Config) {
		s.ForceDefaultAction = b
	}
}

func WithUnmatchedAsError(b bool) cli.Opt {
	return func(s *cli.Config) {
		s.UnmatchedAsError = b
	}
}

// WithStore gives a user-defined Store as initial, or by default
// cmdr makes a dummy Store internally.
//
// So you must have a new Store to be transferred into cmdr if
// you want integrating cmdr and fully-functional Store. Like this,
//
//		app := prepareApp()
//		if err := app.Run(
//			cmdr.WithStore(store.New()),        // create a standard Store instead of internal dummyStore
//			// cmdr.WithExternalLoaders(
//			// 	local.NewConfigFileLoader(),    // import "github.com/hedzr/cmdr-loaders/local" to get in advanced external loading features
//			// 	local.NewEnvVarLoader(),
//			// ),
//			cmdr.WithForceDefaultAction(false), // true for debug in developing time
//		); err != nil {
//			logz.ErrorContext(ctx, "Application Error:", "err", err)
//		}
//
//	 func prepareApp() cli.App {
//			app = cmdr.New().                   // the minimal app is `cmdr.New()`
//				Info("tiny-app", "0.3.1").
//				Author("example.com Authors")
//		}
func WithStore(conf store.Store) cli.Opt {
	return func(s *cli.Config) {
		s.Store = conf
	}
}

// WithExternalLoaders sets the loaders of external sources, which will be loaded
// at cmdr's preparing time (xref-building time).
//
// The orders could be referred as:
//
// - constructing cmdr commands system (by your prepareApp)
// - entering cmdr.Run
// - cmdr preparing stage
//   - build commands and flags xref
//   - load and apply envvars if matched
//   - load external sources
//   - post preparing time
//
// - cmdr parsing stage
// - cmdr invoking stage
// - cmdr cleanup stage
//
// Using our loaders repo is a good idea: https://github.com/hedzr/cmdr-loaders
//
// Typical app:
//
//	app = cmdr.New(opts...).
//		Info("tiny-app", "0.3.1").
//		Author("The Example Authors") // .Description(``).Header(``).Footer(``)
//		cmdr.WithStore(store.New()), // use an option store explicitly, or a dummy store by default
//
//		cmdr.WithExternalLoaders(
//			local.NewConfigFileLoader(), // import "github.com/hedzr/cmdr-loaders/local" to get in advanced external loading features
//			local.NewEnvVarLoader(),
//		),
//	)
//	if err := app.Run(ctx); err != nil {
//		logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
//		os.Exit(app.SuggestRetCode())
//	}
func WithExternalLoaders(loaders ...cli.Loader) cli.Opt {
	return func(s *cli.Config) {
		s.Loaders = append(s.Loaders, loaders...)
	}
}

// WithTasksBeforeParse installs callbacks before parsing stage.
//
// The internal stages are: initial -> preload + xref -> parse -> run/invoke -> post-actions.
func WithTasksBeforeParse(tasks ...cli.Task) cli.Opt {
	return func(s *cli.Config) {
		s.TasksBeforeParse = append(s.TasksBeforeParse, tasks...)
	}
}

// WithTasksParsed installs callbacks after parsing stage.
//
// The internal stages are: initial -> preload + xref -> parse -> run/invoke -> post-actions.
//
// Another way is disabling cmdr default executing/run/invoke stage by WithDontExecuteAction(true).
func WithTasksParsed(tasks ...cli.Task) cli.Opt {
	return func(s *cli.Config) {
		s.TasksParsed = append(s.TasksParsed, tasks...)
	}
}

// WithTasksBeforeRun installs callbacks before run/invoke stage.
//
// The internal stages are: initial -> preload + xref -> parse -> run/invoke -> post-actions.
//
// The internal stages and user-defined tasks are:
//   - initial
//   - preload & xref
//   - <tasksBeforeParse>
//   - parse
//   - <tasksParsed>
//   - <tasksBeforeRun> ( = tasksAfterParse )
//   - exec (run/invoke)
//   - <tasksAfterRun>
//   - <tasksPostCleanup>
//   - basics.closers...Close()
func WithTasksBeforeRun(tasks ...cli.Task) cli.Opt {
	return func(s *cli.Config) {
		s.TasksBeforeRun = append(s.TasksBeforeRun, tasks...)
	}
}

// WithTasksAfterRun installs callbacks after run/invoke stage.
//
// The internal stages are: initial -> preload + xref -> parse -> run/invoke -> post-actions.
func WithTasksAfterRun(tasks ...cli.Task) cli.Opt {
	return func(s *cli.Config) {
		s.TasksAfterRun = append(s.TasksAfterRun, tasks...)
	}
}

// WithTasksPostCleanup install callbacks at cmdr ending.
//
// See the stagings order introduce at [WithTasksBeforeRun].
//
// See also WithTasksSetupPeripherals, WithPeripherals.
func WithTasksPostCleanup(tasks ...cli.Task) cli.Opt {
	return func(s *cli.Config) {
		s.TasksPostCleanup = append(s.TasksPostCleanup, tasks...)
	}
}

// WithTasksSetupPeripherals gives a special chance to setup
// your server's peripherals (such as database, redis, message
// queues, or others).
//
// For these server peripherals, a better practices would be
// initializing them with WithTasksSetupPeripherals and
// destroying them at WithTasksAfterRun.
//
// Another recommendation is implementing your server peripherals
// as a [basics.Closeable] component, and register it with
// basics.RegisterPeripherals(), and that's it. These objects
// will be destroyed at cmdr ends (later than WithTasksAfterRun
// but it's okay).
//
//	import "github.com/hedzr/is/basics"
//	type Obj struct{}
//	func (o *Obj) Init(context.Context) *Obj { return o } // initialize itself
//	func (o *Obj) Close(){...}                            // destroy itself
//	app := cmdr.New(cmdr.WithTasksSetupPeripherals(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) {
//	    obj := new(Obj)
//	    basics.RegisterPeripheral(obj.Init(ctx))          // initialize obj at first, and register it to basics.Closers for auto-shutting-down
//	    return
//	}))
//	...
func WithTasksSetupPeripherals(tasks ...cli.Task) cli.Opt {
	return func(s *cli.Config) {
		s.TasksBeforeRun = append(s.TasksBeforeRun, tasks...)
	}
}

// WithPeripherals is a better way to register your server peripherals
// than WithTasksSetupPeripherals. But the 'better' is less.
//
//	import "github.com/hedzr/is/basics"
//	type Obj struct{}
//	func (o *Obj) Init(context.Context) *Obj { return o } // initialize itself
//	func (o *Obj) Close(){...}                            // destory itself
//	obj := new(Obj)                                       // initialize obj at first
//	ctx := context.Background()                           //
//	app := cmdr.New(cmdr.WithPeripherals(obj.Init(ctx))   // and register it to basics.Closers for auto-shutting-down
//	...
func WithPeripherals(peripherals ...basics.Peripheral) cli.Opt {
	return func(s *cli.Config) {
		basics.RegisterPeripheral(peripherals...)
	}
}

func WithSortInHelpScreen(b bool) cli.Opt {
	return func(s *cli.Config) {
		s.SortInHelpScreen = b
	}
}

func WithDontGroupInHelpScreen(b bool) cli.Opt {
	return func(s *cli.Config) {
		s.DontGroupInHelpScreen = b
	}
}

// WithDontExecuteAction prevents internal exec stage which will
// invoke the matched command's [cli.Cmd.OnAction] handler.
//
// If [cli.Config.DontExecuteAction] is true, cmdr works like
// classical golang stdlib 'flag', which will stop after parsed
// without any further actions.
//
// cmdr itself is a parsing-and-executing processor. We will
// launch a matched command's handlers by default.
func WithDontExecuteAction(b bool) cli.Opt {
	return func(s *cli.Config) {
		s.DontExecuteAction = b
	}
}

// WithAutoEnvBindings enables the feature which can auto-bind env-vars
// to flags default value.
//
// For example, APP_JUMP_TO_FULL=1 -> Cmd{"jump.to"}.Flag{"full"}.default-value = true.
// In this case, `APP` is default prefix so the env-var is different
// than other normal OS env-vars (like HOME, etc.).
//
// You may specify another prefix optionally. For instance, prefix
// "CT" will cause CT_JUMP_TO_FULL=1 binding to
// Cmd{"jump.to"}.Flag{"full"}.default-value = true.
//
// Also you can specify the prefix with multiple section just
// like "CT_ACCOUNT_SERVICE", it will be treated as a normal
// plain string and concatted with the rest parts, so
// "CT_ACCOUNT_SERVICE_JUMP_TO_FULL=1" will be bound in.
func WithAutoEnvBindings(b bool, prefix ...string) cli.Opt {
	return func(s *cli.Config) {
		s.AutoEnv = b
		for _, p := range prefix {
			s.AutoEnvPrefix = p
		}
	}
}

// WithConfig allows you passing a [*cli.Config] object directly.
func WithConfig(conf *cli.Config) cli.Opt {
	return func(s *cli.Config) {
		if conf == nil {
			*s = cli.Config{}
		} else {
			*s = *conf
		}
	}
}

// func WithOnInterpretLeadingPlusSign(cb func()) cli.Opt {
// 	return func(s *cli.Config) {
// 		s.onInterpretLeadingPlusSign = cb
// 	}
// }
