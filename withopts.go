package cmdr

import (
	"context"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/is/basics"
	"github.com/hedzr/store"
)

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
func WithStore(conf store.Store, topLevelPrefix ...string) cli.Opt {
	return func(s *cli.Config) {
		prefix := cli.DefaultStoreKeyPrefix
		for _, pre := range topLevelPrefix {
			prefix = pre
		}
		s.Store = conf.WithPrefix(prefix)
	}
}

func WithRawStore(conf store.Store) cli.Opt {
	return func(s *cli.Config) {
		s.Store = conf
	}
}

// WithExternalLoaders appends the loaders of external sources, which will be loaded
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

// WithExternalLoadersReplaced sets the loaders of external sources, which will be loaded
// at cmdr's preparing time (xref-building time).
func WithExternalLoadersReplaced(loaders ...cli.Loader) cli.Opt {
	return func(s *cli.Config) {
		s.Loaders = loaders
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
//
//	type Obj struct{}
//	func (o *Obj) Open(context.Context) error { return nil } // initialize itself
//	func (o *Obj) Close(){...}                               // destroy itself
//
//	ctx := context.Background()                              //
//	app := cmdr.New(cmdr.WithPeripherals(&Obj{})             // and register it to basics.Closers for auto-shutting-down
//	...
//
// If a peripheral implements `Open(ctx context.Context) error`, it
// will be initialized before running a hit subcommand.
func WithPeripherals(peripherals PeripheralMap) cli.Opt {
	return func(s *cli.Config) {
		s.Store.Set(cli.PeripheralsStoreKey, peripherals)
		for _, peripheral := range peripherals {
			basics.RegisterPeripheral(peripheral)
			if p, ok := peripheral.(interface {
				Open(ctx context.Context) error
			}); ok {
				WithTasksSetupPeripherals(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) {
					err = p.Open(ctx)
					return
				})(s)
			}
		}
	}
}

type PeripheralMap map[string]basics.Peripheral

func Peripheral(name string) basics.Peripheral {
	var pm = Set().MustGet(cli.PeripheralsStoreKey)
	if m, ok := pm.(map[string]basics.Peripheral); ok {
		return m[name]
	}
	return nil
}

func PeripheralT[T basics.Peripheral](name string) T {
	var pm = Set().MustGet(cli.PeripheralsStoreKey)
	if m, ok := pm.(map[string]basics.Peripheral); ok {
		return m[name].(T)
	}
	var t T
	return t
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

func WithOnInterpretLeadingPlusSign(cb cli.OnInterpretLeadingPlusSign) cli.Opt {
	return func(s *cli.Config) {
		s.OnInterpretLeadingPlusSign = cb
	}
}
