// Copyright © 2022 Hedzr Yeh.

package cli

import (
	"context"

	"github.com/hedzr/is/term/color"
	"github.com/hedzr/store"
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

	// ExternalToolEditor environment variable name, EDITOR is fit for most of the shells.
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
type Task func(ctx context.Context, cmd Cmd, runner Runner, extras ...any) (err error)

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

	Cmd // root command here

	preActions  []OnPreInvokeHandler
	postActions []OnPostInvokeHandler
	linked      int32 // ensureTree called?
	app         App
}

type BaseOpt struct {
	owner Cmd
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

type Cmd interface {
	Backtraceable

	String() string

	// App() App

	// Set returns the application Store [store.Store]
	Set() store.Store
	// Store returns the commands subset of the application Store.
	Store() store.Store

	OwnerIsValid() bool
	OwnerIsNil() bool
	OwnerIsNotNil() bool
	OwnerCmd() Cmd
	SetOwnerCmd(o Cmd)
	Root() *RootCommand
	SetRoot(*RootCommand)

	Name() string
	SetName(name string)
	ShortName() string
	ShortNames() []string
	AliasNames() []string
	Desc() string
	DescLong() string
	Examples() string
	TailPlaceHolder() string
	GetCommandTitles() string

	GroupTitle() string                          // group title, removed the ordered prefix
	GroupHelpTitle() string                      // group title, remove the ordered prefix, or UnsortedGroup
	SafeGroup() string                           // group title, or UnsortedGroup
	AllGroupKeys(chooseFlag, sort bool) []string // subcommand group-key-titles
	Hidden() bool
	VendorHidden() bool
	Deprecated() string
	DeprecatedHelpString(trans func(ss string, clr color.Color) string, clr, clrDefault color.Color) (hs, plain string)

	CountOfCommands() int
	CommandsInGroup(groupTitle string) (list []Cmd)
	FlagsInGroup(groupTitle string) (list []*Flag)
	SubCommands() []*CmdS
	Flags() []*Flag

	HeadLikeFlag() *Flag
	SetHeadLikeFlag(*Flag)

	SetHitTitle(title string)
	HitTitle() string
	HitTimes() int

	SetRedirectTo(dottedPath string)

	CanInvoke() bool
	Invoke(ctx context.Context, args []string) (err error)

	OnEvalSubcommands() OnEvaluateSubCommands
	OnEvalSubcommandsOnce() OnEvaluateSubCommands
	OnEvalSubcommandsOnceInvoked() bool
	OnEvalSubcommandsOnceCache() []Cmd
	OnEvalSubcommandsOnceSetCache(list []Cmd)
	OnEvalFlags() OnEvaluateFlags
	OnEvalFlagsOnce() OnEvaluateFlags
	OnEvalFlagsOnceInvoked() bool
	OnEvalFlagsOnceCache() []*Flag
	OnEvalFlagsOnceSetCache(list []*Flag)

	Match(ctx context.Context, title string) (short bool, cc Cmd)
	TryOnMatched(position int, hitState *MatchState) (handled bool, err error)
	MatchFlag(ctx context.Context, vp *FlagValuePkg) (ff *Flag, err error)

	FindSubCommand(ctx context.Context, longName string, wide bool) (res Cmd)
	FindFlagBackwards(ctx context.Context, longName string) (res *Flag)

	WalkGrouped(ctx context.Context, cb WalkGroupedCB)
	WalkBackwardsCtx(ctx context.Context, cb WalkBackwardsCB, pc *WalkBackwardsCtx)
	WalkEverything(ctx context.Context, cb WalkEverythingCB)

	// Walk(ctx context.Context, cb WalkCB)
}

type CmdPriv interface {
	partialMatchFlag(ctx context.Context, title string, short, dblTildeMode bool, cclist map[string]*Flag) (matched, remains string, ff *Flag, err error)
	findSubCommandIn(ctx context.Context, cc Cmd, children []Cmd, longName string, wide bool) (res Cmd)
	findFlagIn(ctx context.Context, cc Cmd, children []Cmd, longName string, wide bool) (res *Flag)
	findFlagBackwardsIn(ctx context.Context, cc Cmd, children []Cmd, longName string) (res *Flag)
}

type CmdS struct {
	BaseOpt

	// tailPlaceHolders gives two places to place the placeholders.
	// It looks like the following form:
	//
	//     dns-util dns add <placeholder1st> [Options] [Parent/Global Options] <placeholders more...>
	//
	// As shown, you may specify at:
	//
	// - before '[Options] [Parent/Global Options]'
	// - after '[Options] [Parent/Global Options]'
	//
	// In TailPlaceHolders slice, [0] is `placeholder-1st``, and others
	// are `placeholders more``.
	//
	// Others:
	//   TailArgsText string [no plan]
	//   TailArgsDesc string [no plan]
	tailPlaceHolders []string

	commands []*CmdS
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
		commands []Cmd
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
	//     when redirectTo is valid, CmdS.OnInvoke handler will be ignored.
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

	longCommands  map[string]*CmdS
	shortCommands map[string]*CmdS

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
	// But you'd better told zsh system with set circuitBreak
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
	A []*CmdS
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

type OnInvokeHandler func(ctx context.Context, cmd Cmd, args []string) (err error)

type OnPostInvokeHandler func(ctx context.Context, cmd Cmd, args []string, errInvoked error) (err error)

type OnPreInvokeHandler func(ctx context.Context, cmd Cmd, args []string) (err error)

type OnEvaluateSubCommands func(ctx context.Context, c Cmd) (it EvalIterator, err error)

type OnEvaluateFlags func(ctx context.Context, c Cmd) (it EvalFlagIterator, err error)

type EvalIterator func(ctx context.Context) (bo Cmd, hasNext bool, err error)

type EvalFlagIterator func(ctx context.Context) (bo *Flag, hasNext bool, err error)

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

type OnCommandMatchedHandler func(c Cmd, position int, hitState *MatchState) (err error)

type OnMatchedHandler func(f *Flag, position int, hitState *MatchState) (err error)

// OnChangingHandler handles when a flag has been setting by parsing command-line
// args, loading from external sources and other cases.
//
// You can cancel the parsing before received a formal OnChanged event,
// for its validation.
type OnChangingHandler func(f *Flag, oldVal, newVal any) (err error)

type OnChangedHandler func(f *Flag, oldVal, newVal any)

type OnSetHandler func(f *Flag, oldVal, newVal any)
