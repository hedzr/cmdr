package cli

import (
	"io"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/store"
)

func NewConfig(opts ...Opt) *Config {
	s := DefaultConfig()
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func DefaultConfig() *Config {
	s := &Config{Store: store.NewDummyStore()}
	return s
}

// Config for cmdr system,
type Config struct {
	store.Store // default is a dummy store. create yours with store.New().

	ForceDefaultAction bool       // use builtin action for debugging if no Action specified to a command
	SortInHelpScreen   bool       // auto sort commands and flags rather than creating order
	UnmatchedAsError   bool       // unmatched command or flag as an error and threw it
	TasksBeforeParse   []Task     // globally pre-parse tasks
	TasksBeforeRun     []Task     // golbally pre-run tasks
	Loaders            []Loader   // external config loaders. use cli.WithLoader() prefer
	HelpScreenWriter   HelpWriter // redirect stdout for help screen printing
	DebugScreenWriter  HelpWriter // redirect stdout for debugging outputs
	Args               []string   // for testing
}

// Opt for cmdr system
type Opt func(s *Config)

// Runner interface for a cmdr workerS.
type Runner interface {
	// InitGlobally initialize all prerequisites, block itself until all
	// of them done and Ready signal changed. Some resources can be exceptions
	// if not required.
	InitGlobally()
	Ready() bool                 // the Runner is built and ready for Run?
	Run(opts ...Opt) (err error) // Run enter the main entry
	DumpErrors(wr io.Writer)     // prints the errors

	Error() errors.Error // return the collected errors in parsing args and invoke actions

	Store() store.Store // app settings store, config set
	Name() string       // app name
	Version() string    // app version
	Root() *RootCommand // root command

	// Actions return a state map.
	// The states can be:
	//   - show-version
	//   - show-built-info
	//   - show-help
	//   - show-help-man
	//   - show-tree
	//   - show-debug
	// These states are produced by parsing the builtin flags
	// with user's command line arguments.
	// For examples, `~~tree` causes 'show-tree' state ON,
	// `--help` causes 'show-help' state ON.
	Actions() (ret map[string]bool)
}

func WithForceDefaultAction(b bool) Opt {
	return func(s *Config) {
		s.ForceDefaultAction = b
	}
}

func WithUnmatchedAsError(b bool) Opt {
	return func(s *Config) {
		s.UnmatchedAsError = b
	}
}

func WithStore(op store.Store) Opt {
	return func(s *Config) {
		if op != nil {
			s.Store = op
		}
	}
}

func WithArgs(args ...string) Opt {
	return func(s *Config) {
		s.Args = args
	}
}

func WithHelpScreenWriter(w HelpWriter) Opt {
	return func(s *Config) {
		s.HelpScreenWriter = w
	}
}

func WithDebugScreenWriter(w HelpWriter) Opt {
	return func(s *Config) {
		s.DebugScreenWriter = w
	}
}

type HelpWriter interface {
	io.Writer
	io.StringWriter
}

func WithExternalLoaders(loaders ...Loader) Opt {
	return func(s *Config) {
		s.Loaders = loaders
	}
}

func WithTasksBeforeParse(tasks ...Task) Opt {
	return func(s *Config) {
		s.TasksBeforeParse = tasks
	}
}

func WithTasksBeforeRun(tasks ...Task) Opt {
	return func(s *Config) {
		s.TasksBeforeRun = tasks
	}
}
