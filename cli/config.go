package cli

import (
	"context"
	"io"
	"strings"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/store"
)

const DefaultStoreKeyPrefix = "app"
const CommandsStoreKey = "cmd"
const PeripheralsStoreKey = "peripherals"

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

	ForceDefaultAction         bool                       `json:"force_default_action,omitempty"`    // use builtin action for debugging if no Action specified to a command
	DontGroupInHelpScreen      bool                       `json:"no_group_in_help_screen,omitempty"` // group commands and flags by its group-name
	DontExecuteAction          bool                       `json:"no_execute_action,omitempty"`       // just parsing, without executing [cli.Cmd.OnAction]
	SortInHelpScreen           bool                       `json:"sort_in_help_screen,omitempty"`     // auto sort commands and flags rather than creating order
	UnmatchedAsError           bool                       `json:"unmatched_as_error,omitempty"`      // unmatched command or flag as an error and threw it
	TasksAfterXref             []Task                     `json:"-"`                                 // while command linked and xref'd, it's time to insert user-defined commands dynamically.
	TasksAfterLoader           []Task                     `json:"-"`                                 // while external loaders loaded.
	TasksBeforeParse           []Task                     `json:"-"`                                 // globally pre-parse tasks
	TasksParsed                []Task                     `json:"-"`                                 // globally post-parse tasks
	TasksBeforeRun             []Task                     `json:"-"`                                 // globally pre-run tasks, it's also used as TasksAfterParsed
	TasksAfterRun              []Task                     `json:"-"`                                 // globally post-run tasks
	TasksPostCleanup           []Task                     `json:"-"`                                 // globally post-run tasks, specially for cleanup actions
	Loaders                    []Loader                   `json:"-"`                                 // external loaders. use cli.WithLoader() prefer
	HelpScreenWriter           HelpWriter                 `json:"help_screen_writer,omitempty"`      // redirect stdout for help screen printing
	DebugScreenWriter          HelpWriter                 `json:"debug_screen_writer,omitempty"`     // redirect stdout for debugging outputs
	Args                       []string                   `json:"args,omitempty"`                    // for testing
	Env                        map[string]string          `json:"env,omitempty"`                     // inject env var & values
	AutoEnv                    bool                       `json:"auto_env,omitempty"`                // enable envvars auto-binding?
	AutoEnvPrefix              string                     `json:"auto_env_prefix,omitempty"`         // envvars auto-binding prefix, bind them to corresponding flags
	OnInterpretLeadingPlusSign OnInterpretLeadingPlusSign `json:"-"`                                 // parsing '+shortFlag`
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

	Error() errors.Error   // return the collected errors in parsing args and invoke actions
	Recycle(errs ...error) // collect the errs object and return the bundled error to main()

	Store() store.Store // app settings store, config set
	Name() string       // app name
	Version() string    // app version
	Root() *RootCommand // root command
	Args() []string     // command-line

	SuggestRetCode() int                      // os process return code
	SetSuggestRetCode(ret int)                // update ret code (0-255) from onAction, onTask, ...
	ParsedState() ParsedState                 // the parsed states
	LoadedSources() (results []LoadedSources) // the loaded sources

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

	// DoBuiltinAction runs a internal action for you.
	//
	// The available internal actions were defined as ActionEnum.
	// Such as ActionShowHelpScreen, or ActionShowVersion.
	DoBuiltinAction(ctx context.Context, action ActionEnum) (handled bool, err error)
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

	// Translate is a helper function, which can interpret the
	// placeholders and translate them to the real value.
	// Translate is used for formatting command/flag's description
	// or examples string.
	//
	// The avaliable placeholders could be: `{{.AppNmae}}`,
	// `{{.AppVersion}}`, `{{.DadCommands}}`, `{{.Commands}}` ...
	Translate(pattern string) (result string)

	DadCommandsText() string
	CommandsText() string
}

type OnInterpretLeadingPlusSign func(w Runner, ctx ParsedState) bool

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

// ActionEnum abstracts the internal OnAction(s).
//
// The available internal actions were defined as ActionEnum.
// Such as ActionShowHelpScreen, or ActionShowVersion.
type ActionEnum int

const (
	ActionShowVersion         ActionEnum = 1 << iota // show version screen
	ActionShowBuiltInfo                              // show build info screen
	ActionShowHelpScreen                             // show help screen
	ActionShowHelpScreenAsMan                        // show help screen in a man-like interactive TUI (by using `less`).
	ActionShowTree                                   // Tree. `~~tree` | show all commands (& flags) as a tree
	ActionShowDebug                                  // Debug. `~~debug` | show debug information for debugging cmdr internal states
	ActionShowDebugEnv                               // with `~~env`
	ActionShowDebugMore                              // with `~~more`
	ActionShowDebugRaw                               // with `~~raw`
	ActionShowDebugValueType                         // with `~~type` (?)
	ActionShowSBOM                                   // show SBOM screen
	// actionShortMode
	// actionDblTildeMode

	ActionDefault // builtin internal action handler
)

func (e ActionEnum) String() string {
	var sb strings.Builder
	if e&ActionShowVersion != 0 {
		_, _ = sb.WriteString("- ShowVersion\n")
	}
	if e&ActionShowBuiltInfo != 0 {
		_, _ = sb.WriteString("- ShowBuiltInfo\n")
	}
	if e&ActionShowHelpScreen != 0 {
		_, _ = sb.WriteString("- ShowHelpScreen\n")
	}
	if e&ActionShowHelpScreenAsMan != 0 {
		_, _ = sb.WriteString("- ShowHelpScreenAsMan\n")
	}
	if e&ActionShowTree != 0 {
		_, _ = sb.WriteString("- ShowTree\n")
	}
	if e&ActionShowDebug != 0 {
		_, _ = sb.WriteString("- ShowDebug\n")
	}
	if e&ActionShowSBOM != 0 {
		_, _ = sb.WriteString("- ShowSBOM\n")
	}
	if e&ActionDefault != 0 {
		_, _ = sb.WriteString("- Default\n")
	}
	return sb.String()
}
