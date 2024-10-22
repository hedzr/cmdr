package cli

import (
	"context"
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
	store.Store `json:"store,omitempty"` // default is a dummy store. create yours with store.New().

	ForceDefaultAction    bool              `json:"force_default_action,omitempty"`    // use builtin action for debugging if no Action specified to a command
	DontGroupInHelpScreen bool              `json:"no_group_in_help_screen,omitempty"` // group commands and flags by its group-name
	SortInHelpScreen      bool              `json:"sort_in_help_screen,omitempty"`     // auto sort commands and flags rather than creating order
	UnmatchedAsError      bool              `json:"unmatched_as_error,omitempty"`      // unmatched command or flag as an error and threw it
	TasksAfterXref        []Task            `json:"-"`                                 // while command linked and xref'd, it's time to insert user-defined commands dynamically.
	TasksAfterLoader      []Task            `json:"-"`                                 // while external loaders loaded.
	TasksBeforeParse      []Task            `json:"-"`                                 // globally pre-parse tasks
	TasksBeforeRun        []Task            `json:"-"`                                 // globally pre-run tasks, it's also used as TasksAfterParsed
	TasksAfterRun         []Task            `json:"-"`                                 // globally post-run tasks
	Loaders               []Loader          `json:"-"`                                 // external loaders. use cli.WithLoader() prefer
	HelpScreenWriter      HelpWriter        `json:"help_screen_writer,omitempty"`      // redirect stdout for help screen printing
	DebugScreenWriter     HelpWriter        `json:"debug_screen_writer,omitempty"`     // redirect stdout for debugging outputs
	Args                  []string          `json:"args,omitempty"`                    // for testing
	Env                   map[string]string `json:"env,omitempty"`                     // inject env var & values
}

// Opt for cmdr system
type Opt func(s *Config)

// Runner interface for a cmdr workerS.
type Runner interface {
	// InitGlobally initialize all prerequisites, block itself until all
	// of them done and Ready signal changed. Some resources can be exceptions
	// if not required.
	InitGlobally(ctx context.Context)                 // ctx context.Context)
	Ready() bool                                      // the Runner is built and ready for Run?
	Run(ctx context.Context, opts ...Opt) (err error) // Run enter the main entry
	DumpErrors(wr io.Writer)                          // prints the errors

	Error() errors.Error // return the collected errors in parsing args and invoke actions

	Store() store.Store // app settings store, config set
	Name() string       // app name
	Version() string    // app version
	Root() *RootCommand // root command
	Args() []string     // command-line

	SuggestRetCode() int       // os process return code
	SetSuggestRetCode(ret int) // update ret code (0-255) from onAction, onTask, ...
	ParsedState() ParsedState  // the parsed states

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

type ParsedState interface {
	NoCandidateChildCommands() bool
	LastCmd() Cmd
	MatchedCommands() []Cmd
	MatchedFlags() map[*Flag]*MatchState
	PositionalArgs() []string

	CommandMatchedState(c Cmd) (ms *MatchState)
	FlagMatchedState(f *Flag) (ms *MatchState)

	HasCmd(longTitle string, validator func(cc Cmd, state *MatchState) bool) (found bool)
	HasFlag(longTitle string, validator func(ff *Flag, state *MatchState) bool) (found bool)

	Translate(pattern string) (result string)

	DadCommandsText() string
	CommandsText() string
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
