// Copyright © 2022 Hedzr Yeh.

package cli

import (
	"context"
)

const (
	appNameDefault = "cmdr" //nolint:deadcode,unused,varcheck //keep it

	// UnsortedGroup for commands and flags
	UnsortedGroup = "!!!!.Unsorted"
	// AddonsGroup for commands and flags
	AddonsGroup = "zzzh.Addons"
	// ExtGroup for commands and flags
	ExtGroup = "zzzi.Extensions"
	// AliasesGroup for commands and flags
	AliasesGroup = "zzzj.Aliases"
	// SysMgmtGroup for commands and flags
	SysMgmtGroup = "zzz9.Misc"

	// DefaultEditor is 'vim'
	DefaultEditor = "vim"

	// ExternalToolEditor environment variable name, EDITOR is fit for most of shells.
	ExternalToolEditor = "EDITOR"

	// ExternalToolPasswordInput enables secure password input without echo.
	ExternalToolPasswordInput = "PASSWD"
)

// Task is used for WithTasksBeforeRun, WithTasksBeforeParse, and Config.TasksAfterRun ...
//
// extras is a non-typed args, 1st is worker.parseCtx,
// 2nd is parseCtx.PositionalArgs(), ...
//
// There is no TasksAfterParsed, but you can replace it
// with WithTasksBeforeRun/Config.TasksBeforeRun.
type Task func(cmd *Command, runner Runner, extras ...any) (err error)

type Loader interface {
	Load(app App) (err error)
}

type RootCommand struct {
	AppName string
	Version string
	// AppDescription string
	// AppLongDesc    string
	Copyright  string
	Author     string
	HeaderLine string
	FooterLine string

	*Command

	preActions  []OnPreInvokeHandler
	postActions []OnPostInvokeHandler

	app App
}

type BaseOptI interface {
	Owner() *Command
	Root() *RootCommand
	Name() string
}

type BaseOpt struct {
	owner *Command
	root  *RootCommand

	// name is reserved for internal purpose.
	name string
	// Long is a full/long form flag/command title name.
	// word string. example for flag: "addr" -> "--addr"
	Long string
	// Short 'rune' string. short option/command name.
	// single char. example for flag: "a" -> "-a"
	Short string
	// Aliases are the more synonyms
	Aliases     []string
	description string
	longDesc    string
	examples    string
	// group specify a group name,
	// A special prefix could sort it, has a form like `[0-9a-zA-Z]+\.`.
	// The prefix will be removed from help screen.
	//
	// Some examples are:
	//    "A001.Host Params"
	//    "A002.User Params"
	//
	// If ToggleGroup specified, Group field can be omitted because we will copy
	// from there.
	group string

	extraShorts []string // more short titles

	// deprecated is a version string just like '0.5.9' or 'v0.5.9', that
	// means this command/flag was/will be deprecated since `v0.5.9`.
	deprecated   string
	hidden       bool
	vendorHidden bool

	// hitTitle keeps the matched title string from user input in command line
	hitTitle string
	// hitTimes how many times this flag was triggered.
	// To access it with `Flag.GetTriggeredTimes()`, `cmdr.GetFlagHitCount()`,
	// `cmdr.GetFlagHitCountRecursively()` or `cmdr.GetHitCountByDottedPath()`.
	hitTimes int
}

type Command struct {
	BaseOpt

	// tailPlaceHolders gives two places to place the placeholders.
	// It looks like the following form:
	//
	//     austr dns add <placeholder1st> [Options] [Parent/Global Options] <placeholders more...>
	//
	// As shown, you may specify at:
	//
	// - before '[Options] [Parent/Global Options]'
	// - after '[Options] [Parent/Global Options]'
	//
	// In TailPlaceHolders slice, [0] is `placeholder1st``, and others
	// are `placeholders more``.
	//
	// Others:
	//   TailArgsText string [no plan]
	//   TailArgsDesc string [no plan]
	tailPlaceHolders []string

	commands []*Command
	flags    []*Flag

	// preActions will be launched before running OnInvoke.
	// The return value obj.ErrShouldStop will cause the remained
	// following processing flow broken right away.
	preActions []OnPreInvokeHandler
	// onInvoke is the main action or entry point when the command
	// was hit from parsing command-line arguments.
	onInvoke OnInvokeHandler
	// postActions will be launched after running OnInvoke.
	postActions []OnPostInvokeHandler

	onMatched []OnCommandMatchedHandler

	onEvalSubcommands *struct {
		cb OnEvaluateSubCommands
		// invoked bool
	}
	onEvalSubcommandsOnce *struct {
		cb       OnEvaluateSubCommands
		invoked  bool
		commands []*Command
	}
	onEvalFlags *struct {
		cb OnEvaluateFlags
		// invoked bool
	}
	onEvalFlagsOnce *struct {
		cb      OnEvaluateFlags
		invoked bool
		flags   []*Flag
	}

	// redirectTo gives the dotted-path to point to a subcommand.
	//
	// Thd target subcommand will be invoked while this command is being invoked.
	//
	// For example, if RootCommand.redirectTo is set to "build", and
	// entering "app" will equal to entering "app build ...".
	//
	// NOTE:
	//
	//     when redirectTo is valid, Command.OnInvoke handler will be ignored.
	redirectTo string

	//

	presetCmdLines []string

	// invokeProc is just for cmdr aliases commands
	// invoke the external commands (via: executable)
	invokeProc string // `yaml:"invoke-proc,omitempty" `
	// invokeShell is just for cmdr aliases commands
	// invoke the external commands (via: shell)
	invokeShell string // `yaml:"invoke-sh,omitempty" `
	// shell is just for cmdr aliases commands
	// invoke a command under this shell, or /usr/bin/env bash|zsh|...
	shell string // `yaml:"shell,omitempty" `

	// internal indices ------------------

	longCommands  map[string]*Command
	shortCommands map[string]*Command

	longFlags  map[string]*Flag
	shortFlags map[string]*Flag
	// allLongFlags  map[string]*Flag
	// allShortFlags map[string]*Flag

	allCommands map[string]*CmdSlice
	allFlags    map[string]*FlgSlice

	toggles      map[string]*ToggleGroupMatch // key: toggle-group
	headLikeFlag *Flag
}

type ToggleGroupMatch struct {
	Flags        map[string]*Flag // key: flg.Long
	Matched      *Flag
	MatchedTitle string
}

func (s *ToggleGroupMatch) MatchedFlag() *Flag { return s.Matched }

type Flag struct {
	BaseOpt

	toggleGroup  string
	placeHolder  string
	defaultValue any
	envVars      []string

	externalEditor string   // env-var name of the external editor
	validArgs      []string // enum values
	min, max       int
	headLike       bool
	requited       bool

	onParseValue OnParseValueHandler // allows user-defined value parsing, converting and validating
	onMatched    OnMatchedHandler    // cancellable, after parsed from cmdline, new value got, and before old value got
	onChanging   OnChangingHandler   // cancellable notifier (a validator) before a formal on-changed notification, = OnValidating
	onChanged    OnChangedHandler    // modified generally (programmatically, cmdline parsing, cfg file, ...)
	onSet        OnSetHandler        // modified programmatically

	// actionStr: for zsh completion, see action of an optspec in _argument
	actionStr string
	// mutualExclusives is used for zsh completion.
	//
	// For the ToggleGroup group, mutualExclusives is implicit.
	mutualExclusives []string
	// prerequisites flags for this one.
	//
	// In zsh completion, any of prerequisites flags must be present
	// so that user can complete this one.
	//
	// The prerequisites were not present and cmdr would report error
	// and stop parsing flow.
	prerequisites []string
	// justOnce is used for zsh completion.
	justOnce bool
	// circuitBreak is used for zsh completion.
	//
	// A flag can break cmdr parsing flow with return
	// ErrShouldStop in its Action handler.
	// But you' better told zsh system with set circuitBreak
	// to true. At this case, cmdr will generate a suitable
	// completion script.
	circuitBreak bool
	// dblTildeOnly can be used for zsh completion.
	//
	// A DblTildeOnly Flag accepts '~~opt' only, so '--opt' is
	// invalid form and couldn't be used for other Flag
	// anymore.
	dblTildeOnly bool // such as '~~tree'
}

type CmdSlice struct {
	A []*Command
}

type FlgSlice struct {
	A []*Flag
}

type MatchState struct {
	Short, DblTilde bool
	HitStr          string
	HitTimes        int
	Value           any
}

type OnInvokeHandler func(cmd *Command, args []string) (err error)

type OnPostInvokeHandler func(cmd *Command, args []string, errInvoked error) (err error)

type OnPreInvokeHandler func(cmd *Command, args []string) (err error)

type OnCommandMatchedHandler func(c *Command, position int, hitState *MatchState) (err error)

type OnEvaluateSubCommands func(ctx context.Context, c *Command) (it EvalIterator, err error)

type OnEvaluateFlags func(ctx context.Context, c *Command) (it EvalIterator, err error)

type EvalIterator func() (bo BaseOptI, hasNext bool, err error)

// OnParseValueHandler could be used for parsing value string as you want,
// and/or check the validation of the input value or flag, and so on.
//
// return err == obj.ErrShouldFallback: let cmdr fallback to the default implementation;
// return err == obj.ErrShouldStop: let cmdr stop parsing action.
type OnParseValueHandler func(
	f *Flag,
	position int,
	hitCaption string,
	hitValue string,
	moreArgs []string,
) (
	newVal any,
	remainPartInHitValue string,
	err error,
)

type OnMatchedHandler func(f *Flag, position int, hitState *MatchState) (err error)

// OnChangingHandler handles when a flag is been setting by parsing command-line
// args, loading from external sources and other cases.
//
// You can cancel the parsing before received a formal OnChanged event,
// for its validation.
type OnChangingHandler func(f *Flag, oldVal, newVal any) (err error)

type OnChangedHandler func(f *Flag, oldVal, newVal any)

type OnSetHandler func(f *Flag, oldVal, newVal any)
