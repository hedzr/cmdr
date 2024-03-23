// Copyright © 2022 Hedzr Yeh.

package cmdr

import (
	"os"

	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2/builder"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/worker"
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
//	app.AddCmd(func(b config.CommandBuilder) {
//	    b.Titles("jump").
//	        Description("jump command").
//	        Examples(``).
//	        Deprecated(``).
//	        Hidden(false).
//	        AddCmd(func(b config.CommandBuilder) {
//	            b.Titles("to").
//	                Description("to command").
//	                Examples(``).
//	                Deprecated(``).
//	                Hidden(false).
//	                OnAction(func(cmd *obj.Command, args []string) (err error) {
//	                    return // handling command action here
//	                }).
//	                Build()
//	        }).
//	        Build()
//	}).AddFlg(func(b config.FlagBuilder) {
//	    b.Titles("dry-run", "n").Default(false).Build()
//	})
//
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
//	    logz.Error("Application Error:", "err", err)
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
	cfg := cli.NewConfig(opts...)
	w := worker.New(cfg)
	return builder.New(w)
}

// App returns a light version of builder.Runner (a.k.a. *worker.Worker).
//
// Generally it's a unique instance in one system.
//
// It's available once New() / Exec() called, else nil.
func App() cli.Runner { return worker.UniqueWorker() }

// Exec starts a new cmdr app (parsing cmdline args based on the given rootCmd)
// from scratch.
//
// It's a reserved API for back-compatible with cmdr v1.
//
// It'll be removed completely at the recently future version.
//
// Deprecated since 2.1
func Exec(rootCmd *cli.RootCommand, opts ...cli.Opt) (err error) {
	// if is.InDebugging() {
	// 	_ = exec.Run("/bin/false")
	// 	// cabin.Version()
	// 	// cpcn.Out()
	// }

	app := New(opts...).WithRootCommand(rootCmd)
	err = app.Run()
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

func WithStore(conf store.Store) cli.Opt {
	return func(s *cli.Config) {
		s.Store = conf
	}
}

func WithExternalLoaders(loaders ...cli.Loader) cli.Opt {
	return func(s *cli.Config) {
		s.Loaders = loaders
	}
}

func WithTasksBeforeParse(tasks ...cli.Task) cli.Opt {
	return func(s *cli.Config) {
		s.TasksBeforeParse = tasks
	}
}

func WithTasksBeforeRun(tasks ...cli.Task) cli.Opt {
	return func(s *cli.Config) {
		s.TasksBeforeRun = tasks
	}
}
