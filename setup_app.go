package cmdr

import (
	"os"

	"github.com/hedzr/cmdr/v2/builder"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/cli/worker"
	"github.com/hedzr/cmdr/v2/pkg/logz"
	"github.com/hedzr/store"
)

// Create provides a concise interface to create an cli app easily.
func Create(appName, version, author, desc string, opts ...cli.Opt) Creator {
	s := &cs{}
	return s.Create(appName, version, author, desc, opts...)
}

type Creator interface {
	WithOpts(opts ...cli.Opt) Creator
	// WithAdders receives a couple of cli.CmdAdder adders
	// which can initialize a command standalone.
	WithAdders(adders ...cli.CmdAdder) Creator
	// WithBuilder receives a couple of callbacks of cli.CommandBuilder
	// which can initialize command and flag.
	WithBuilders(builders ...func(b cli.CommandBuilder)) Creator
	// With a callback, you can initialize commands and
	// flags by app object directly.
	With(cb func(app cli.App)) Creator

	// Build creates the final app object and stop the
	// building sequence of a builder pattern.
	Build() (app cli.App)
}

type cs struct {
	app cli.App
}

func (s *cs) Create(appName, version, author, desc string, opts ...cli.Opt) Creator {
	// A cmdr app will close all peripherals in basics.Closers() at exiting.
	// So you could always register the objects which wanna be cleanup at
	// app terminating, by [basics.RegisterPeripheral(...)].
	// See also: https://github.com/hedzr/is/blob/master/basics/ and

	_ = os.Setenv("CMDR_VERSION", Version)
	logz.Verbose("setup env-var at earlier time", "CMDR_VERSION", Version)
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

func (s *cs) Build() (app cli.App) {
	return s.app
}
