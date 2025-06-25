package cmdr

import (
	"github.com/hedzr/cmdr/v2/builder"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/worker"
	"github.com/hedzr/store"
)

// Create provides a concise interface to create an cli app easily.
//
// See also [New].
//
// A tiny sample app can be:
//
//	package main
//
//	import (
//		"context"
//		"os"
//
//		"github.com/hedzr/cmdr/v2"
//		"github.com/hedzr/cmdr/v2/examples/cmd"
//		"github.com/hedzr/cmdr/v2/pkg/logz"
//	)
//
//	const (
//		appName = "concise"
//		desc    = `concise version of tiny app.`
//		version = cmdr.Version
//		author  = `The Example Authors`
//	)
//
//	func main() {
//		app := cmdr.Create(appName, version, author, desc).
//			WithAdders(cmd.Commands...).
//			Build()
//
//		ctx, cancel := context.WithCancel(context.Background())
//		defer cancel()
//
//		if err := app.Run(ctx); err != nil {
//			logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
//			os.Exit(app.SuggestRetCode())
//		} else if rc := app.SuggestRetCode(); rc != 0 {
//			os.Exit(rc)
//		}
//	}
func Create(appName, version, author, desc string, opts ...cli.Opt) Creator {
	s := &cs{}
	return s.Create(appName, version, author, desc, opts...)
}

type Creator interface {
	WithOpts(opts ...cli.Opt) Creator
	// WithAdders receives a couple of cli.CmdAdder adders
	// which can initialize a command standalone.
	WithAdders(adders ...cli.CmdAdder) Creator
	// WithBuilders receives a couple of callbacks of cli.CommandBuilder
	// which can initialize command and flag.
	WithBuilders(builders ...func(b cli.CommandBuilder)) Creator
	// With a callback, you can initialize commands and
	// flags by app object directly.
	With(cb func(app cli.App)) Creator

	// OnAction apply the root-level onAction handler to the app object.
	OnAction(handler cli.OnInvokeHandler) Creator

	// BuildFrom builds command system from a given struct-value and
	// creates the final app object right now.
	BuildFrom(structValue any, opts ...cli.StructBuilderOpt) (app cli.App)

	// Build creates the final app object and stop the
	// building sequence of a builder pattern.
	Build() (app cli.App)
}

type cs struct {
	app cli.App
}

func (s *cs) OnAction(handler cli.OnInvokeHandler) Creator {
	s.app.OnAction(handler)
	return s
}

func (s *cs) Create(appName, version, author, desc string, opts ...cli.Opt) Creator {
	// A cmdr app will close all peripherals in basics.Closers() at exiting.
	// So you could always register the objects which wanna be cleanup at
	// app terminating, by [basics.RegisterPeripheral(...)].
	// See also: https://github.com/hedzr/is/blob/master/basics/ and

	earlierInitForNew()
	cfg := cli.NewConfig(append([]cli.Opt{WithStore(store.New())}, opts...)...)
	w := worker.New(cfg)
	s.app = builder.New(w)

	// s.app = New(
	// 	// use an option store explicitly, or a dummy store by default
	// 	append([]cli.Opt{WithStore(store.New())}, opts...)...,
	// )

	s.app.
		Info(appName, version, desc).
		Author(author) // .Description(``).Header(``).Footer(``)
	return s
}

func (s *cs) WithOpts(opts ...cli.Opt) Creator {
	s.app.WithOpts(opts...)
	return s
}

func (s *cs) WithBuilders(builders ...func(b cli.CommandBuilder)) Creator {
	for _, cb := range builders {
		s.app.RootBuilder(cb)
	}
	return s
}

func (s *cs) With(cb func(app cli.App)) Creator {
	if cb != nil {
		cb(s.app)
	}
	return s
}

func (s *cs) WithAdders(adders ...cli.CmdAdder) Creator {
	for _, adder := range adders {
		adder.Add(s.app)
	}
	return s
}

func (s *cs) BuildFrom(structValue any, opts ...cli.StructBuilderOpt) (app cli.App) {
	s.app.
		FromStruct(structValue, opts...).
		Build()
	return s.app
}

func (s *cs) Build() (app cli.App) {
	return s.app
}
