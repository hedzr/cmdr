/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

type (
	// // Opt never used?
	// Opt interface {
	// 	Titles(short, long string, aliases ...string) (opt Opt)
	// 	Short(short string) (opt Opt)
	// 	Long(long string) (opt Opt)
	// 	Aliases(ss ...string) (opt Opt)
	// 	Description(oneLine, long string) (opt Opt)
	// 	Examples(examples string) (opt Opt)
	// 	Group(group string) (opt Opt)
	// 	Hidden(hidden bool) (opt Opt)
	// 	Deprecated(deprecation string) (opt Opt)
	// 	Action(action Handler) (opt Opt)
	// }

	// OptFlag to support fluent api of cmdr.
	// see also: cmdr.Root().NewSubCommand()/.NewFlag()
	//
	// For an option, its default value must be declared with exact
	// type as is
	OptFlag interface {
		// Titles broken API since v1.6.39.
		//
		// If necessary, an ordered prefix can be attached to the long
		// title. The title with prefix will be set to Name field and
		// striped to Long field.
		//
		// An ordered prefix is a dotted string with multiple alphabets
		// and digits. Such as:
		//     "zzzz.", "0001.", "700.", "A1." ...
		Titles(long, short string, aliases ...string) (opt OptFlag)
		// Short gives short sub-command title in string representation.
		//
		// A short command title is often one char, sometimes it can be
		// two chars, but even more is allowed.
		Short(short string) (opt OptFlag)
		// Long gives long flag title, we call it 'Full Title' too.
		//
		// A long title should be one or more words in english, separated
		// by short hypen ('-'), for example 'auto-increment'.
		//
		Long(long string) (opt OptFlag)
		// Name is an internal identity, and an order prefix is optional
		//
		// An ordered prefix is a dotted string with multiple alphabets
		// and digits. Such as:
		//     "zzzz.", "0001.", "700.", "A1." ...
		Name(name string) (opt OptFlag)
		// Aliases give more choices for a command.
		Aliases(ss ...string) (opt OptFlag)
		Description(oneLineDesc string, longDesc ...string) (opt OptFlag)
		Examples(examples string) (opt OptFlag)
		// Group provides flag group name.
		//
		// All flags with same group name will be agreegated in
		// one group.
		//
		// An order prefix is a dotted string with multiple alphabet
		// and digit. Such as:
		//     "zzzz.", "0001.", "700.", "A1." ...
		Group(group string) (opt OptFlag)
		// Hidden flag will not be shown in help screen, expect user
		// entered 'app --help -vvv'.
		//
		// NOTE -v is a builtin flag, -vvv means be triggered 3rd times.
		//
		// And the triggered times of a flag can be retrieved by
		// flag.GetTriggeredTimes(), or via cmdr Option Store API:
		// cmdr.GetFlagHitCount("verbose") / cmdr.GetVerboseHitCount().
		Hidden(hidden bool) (opt OptFlag)
		// VendorHidden flag will never be shown in help screen, even
		// if user was entering 'app --help -vvv'.
		//
		// NOTE -v is a builtin flag, -vvv means be triggered 3rd times.
		//
		// And the triggered times of a flag can be retrieved by
		// flag.GetTriggeredTimes(), or via cmdr Option Store API:
		// cmdr.GetFlagHitCount("verbose") / cmdr.GetVerboseHitCount().
		VendorHidden(hidden bool) (opt OptFlag)
		// Deprecated gives a tip to users to encourage them to move
		// up to some new API/commands/operations.
		//
		// A tip is starting by 'Since v....' typically, for instance:
		//
		// Since v1.10.36, PlaceHolder is deprecated, use PlaceHolders
		// is recommended for more abilities.
		Deprecated(deprecation string) (opt OptFlag)
		// Action will be triggered once being parsed ok
		Action(action Handler) (opt OptFlag)

		// ToggleGroup provides the RADIO BUTTON model for a set of
		// flags.
		//
		// The flags with same toggle group name will be agreegated
		// into one group.
		//
		// One of them are set will clear all of others flags in same
		// toggle group.
		// In programmtically, you may retrieve the final choice by its
		// name, see also toggle-group example (at our cmdr-examples repo).
		//
		// NOTE When the ToggleGroup specified, Group can be ignored.
		//
		//    func tg(root cmdr.OptCmd) {
		//      // toggle-group
		//
		//      c := cmdr.NewSubCmd().Titles("toggle-group", "tg").
		//        Description("soundex test").
		//        Group("Test").
		//        TailPlaceholder("[text1, text2, ...]").
		//        Action(func(cmd *cmdr.Command, args []string) (err error) {
		//          selectedMuxType := cmdr.GetStringR("toggle-group.mux-type")
		//          fmt.Printf("Flag 'echo' = %v\n", cmdr.GetBoolR("toggle-group.echo"))
		//          fmt.Printf("Flag 'gin' = %v\n", cmdr.GetBoolR("toggle-group.gin"))
		//          fmt.Printf("Flag 'gorilla' = %v\n", cmdr.GetBoolR("toggle-group.gorilla"))
		//          fmt.Printf("Flag 'iris' = %v\n", cmdr.GetBoolR("toggle-group.iris"))
		//          fmt.Printf("Flag 'std' = %v\n", cmdr.GetBoolR("toggle-group.std"))
		//          fmt.Printf("Toggle Group 'mux-type' = %v\n", selectedMuxType)
		//          return
		//        }).
		//        AttachTo(root)
		//
		//      cmdr.NewBool(false).Titles("echo", "echo").Description("using 'echo' mux").ToggleGroup("mux-type").Group("Mux").AttachTo(c)
		//      cmdr.NewBool(false).Titles("gin", "gin").Description("using 'gin' mux").ToggleGroup("mux-type").Group("Mux").AttachTo(c)
		//      cmdr.NewBool(false).Titles("gorilla", "gorilla").Description("using 'gorilla' mux").ToggleGroup("mux-type").Group("Mux").AttachTo(c)
		//      cmdr.NewBool(true).Titles("iris", "iris").Description("using 'iris' mux").ToggleGroup("mux-type").Group("Mux").AttachTo(c)
		//      cmdr.NewBool(false).Titles("std", "std").Description("using standardlib http mux mux").ToggleGroup("mux-type").Group("Mux").AttachTo(c)
		//    }
		//
		ToggleGroup(group string) (opt OptFlag)
		// DefaultValue needs an exact typed 'val'.
		//
		// IMPORTANT: cmdr interprets value type of an option based
		// on the underlying default value set.
		DefaultValue(val interface{}, placeholder string) (opt OptFlag)
		Placeholder(placeholder string) (opt OptFlag)
		CompletionActionStr(s string) (opt OptFlag)
		// CompletionMutualExclusiveFlags is a slice of flag full/long titles
		CompletionMutualExclusiveFlags(flags ...string) (opt OptFlag)
		// CompletionPrerequisitesFlags is a slice of flag full/long titles
		CompletionPrerequisitesFlags(flags ...string) (opt OptFlag)
		CompletionJustOnce(once bool) (opt OptFlag)
		CompletionCircuitBreak(once bool) (opt OptFlag)
		// DoubleTildeOnly requests only the form is okay when the
		// double tilde chars ('~~') leads the flag.
		//
		// So, short form (-f) and long form (--flag) will be interpreted as
		// unknown flag found.
		DoubleTildeOnly(once bool) (opt OptFlag)
		// ExternalTool provides a OS-environment variable name,
		// which identify the position of an external tool.
		// When cmdr parsed command-line ok, the external tool
		// will be invoked at first.
		//
		// For example:
		//
		// 'git commit -m' will wake up EDITOR to edit a commit message.
		//
		ExternalTool(envKeyName string) (opt OptFlag)
		// ValidArgs gives a constrained list for input value of
		// a flag.
		//
		// ValidArgs provides ENUM type on command-line.
		ValidArgs(list ...string) (opt OptFlag)
		// HeadLike enables `head -n` mode.
		//
		// 'min', 'max' will be ignored at this version, its might be
		// impl in the future.
		//
		// There's only one head-like flag in one command and its parent
		// and children commands.
		HeadLike(enable bool, min, max int64) (opt OptFlag)

		// EnvKeys is a list of env-var names of binding on this flag.
		EnvKeys(keys ...string) (opt OptFlag)
		// Required flag is MUST BE supplied by user entering from
		// command-line.
		//
		// NOTE Required() sets the required flag to true while
		// it was been invoking with empty params.
		Required(required ...bool) (opt OptFlag)

		// OwnerCommand returns the parent command in OptCmd form.
		//
		// It's avaiable for building time.
		OwnerCommand() (opt OptCmd)
		SetOwner(opt OptCmd)

		RootCommand() *RootCommand

		ToFlag() *Flag

		// AttachTo attach as a flag of `opt` OptCmd object
		AttachTo(parent OptCmd) (opt OptFlag)
		// AttachToCommand attach as a flag of *Command object
		AttachToCommand(cmd *Command) (opt OptFlag)
		// AttachToRoot attach as a flag of *RootCommand object
		AttachToRoot(root *RootCommand) (opt OptFlag)

		OnSet
	}

	// OptCmd to support fluent api of cmdr.
	// see also: cmdr.Root().NewSubCommand()/.NewFlag()
	OptCmd interface {
		// Titles broken API since v1.6.39
		//
		// If necessary, an order prefix can be attached to the long title.
		// The title with prefix will be set to Name field and striped to Long field.
		//
		// An order prefix is a dotted string with multiple alphabet and digit. Such as:
		// "zzzz.", "0001.", "700.", "A1." ...
		//
		Titles(long, short string, aliases ...string) (opt OptCmd)
		// Short gives short sub-command title in string representation.
		//
		// A short command title is often one char, sometimes it can be
		// two chars, but even more is allowed.
		Short(short string) (opt OptCmd)
		// Long gives long sub-command title, we call it 'Full Title' too.
		//
		// A long title should be one or more words in english, separated
		// by short hypen ('-'), for example 'create-new-account'.
		//
		// Of course, in a multiple-level command system, a better
		// hierarchical command structure might be:
		//
		//     ROOT
		//       account
		//         create
		//         edit
		//         suspend
		//         destroy
		//
		Long(long string) (opt OptCmd)
		// Name is an internal identity, and an order prefix is optional
		//
		// An ordered prefix is a dotted string with multiple alphabets
		// and digits. Such as:
		//     "zzzz.", "0001.", "700.", "A1." ...
		Name(name string) (opt OptCmd)
		// Aliases give more choices for a command.
		Aliases(ss ...string) (opt OptCmd)
		Description(oneLine string, long ...string) (opt OptCmd)
		Examples(examples string) (opt OptCmd)
		// Group provides command group name.
		//
		// All commands with same group name will be agreegated in
		// one group.
		//
		// An order prefix is a dotted string with multiple alphabet
		// and digit. Such as:
		//     "zzzz.", "0001.", "700.", "A1." ...
		Group(group string) (opt OptCmd)
		// Hidden command will not be shown in help screen, expect user
		// entered 'app --help -vvv'.
		//
		// NOTE -v is a builtin flag, -vvv means be triggered 3rd
		// times.
		//
		// And the triggered times of a flag can be retrieved by
		// flag.GetTriggeredTimes(), or via cmdr Option Store API:
		// cmdr.GetFlagHitCount("verbose") / cmdr.GetVerboseHitCount().
		Hidden(hidden bool) (opt OptCmd)
		// VendorHidden command will never be shown in help screen, even
		// if user was entering 'app --help -vvv'.
		//
		// NOTE -v is a builtin flag, -vvv means be triggered 3rd times.
		//
		// And the triggered times of a flag can be retrieved by
		// flag.GetTriggeredTimes(), or via cmdr Option Store API:
		// cmdr.GetFlagHitCount("verbose") / cmdr.GetVerboseHitCount().
		VendorHidden(hidden bool) (opt OptCmd)
		// Deprecated gives a tip to users to encourage them to move
		// up to some new API/commands/operations.
		//
		// A tip is starting by 'Since v....' typically, for instance:
		//
		// Since v1.10.36, PlaceHolder is deprecated, use PlaceHolders
		// is recommended for more abilities.
		Deprecated(deprecation string) (opt OptCmd)
		// Action will be triggered after all command-line arguments parsed.
		//
		// Action might be the most important entry for a command.
		//
		// For a nest command system, parent command shouldn't get a
		// valid Action handler because we need to get a chance to
		// step into its children for parsing command-line.
		Action(action Handler) (opt OptCmd)

		// FlagAdd(flg *Flag) (opt OptCmd)
		// SubCommand(cmd *Command) (opt OptCmd)

		// PreAction will be invoked before running Action
		// NOTE that RootCommand.PreAction will be invoked too.
		PreAction(pre Handler) (opt OptCmd)
		// PostAction will be invoked after run Action
		// NOTE that RootCommand.PostAction will be invoked too.
		PostAction(post Invoker) (opt OptCmd)

		// TailPlaceholder gives two places to place the placeholders.
		// It looks like the following form:
		//
		//     austr dns add <placeholder-1st> [Options] [Parent/Global Options] <placeholders-more...>
		//
		// As shown, you may specify at:
		//
		// - before '[Options] [Parent/Global Options]'
		// - after '[Options] [Parent/Global Options]'
		//
		// In TailPlaceHolders slice, [0] is `placeholder-1st``, and others
		// are `placeholders-more`.
		//
		// Others:
		//   TailArgsText string [no plan]
		//   TailArgsDesc string [no plan]
		TailPlaceholder(placeholders ...string) (opt OptCmd)

		// Sets _
		//
		// Reserved API.
		Sets(func(cmd OptCmd)) (opt OptCmd)

		// NewFlag create a new flag object and return it for further operations.
		// Deprecated since v1.6.9, replace it with FlagV(defaultValue)
		//
		// Deprecated since v1.6.50, we recommend the new form:
		//    cmdr.NewBool(false).Titles(...)...AttachTo(ownerCmd)
		// NewFlag(typ OptFlagType) (opt OptFlag)
		// NewFlagV create a new flag object and return it for further operations.
		// the titles in arguments MUST be: longTitle, [shortTitle, [aliasTitles...]]
		//
		// Deprecated since v1.6.50, we recommend the new form:
		//    cmdr.NewBool(false).Titles(...)...AttachTo(ownerCmd)
		// NewFlagV(defaultValue interface{}, titles ...string) (opt OptFlag)
		// NewSubCommand make a new sub-command optcmd object with optional titles.
		// the titles in arguments MUST be: longTitle, [shortTitle, [aliasTitles...]]
		//
		// Deprecated since v1.6.50
		// NewSubCommand(titles ...string) (opt OptCmd)

		// OwnerCommand returns the parent command in OptCmd form.
		//
		// It's avaiable for building time.
		OwnerCommand() (opt OptCmd)
		//
		SetOwner(opt OptCmd)

		RootCommand() *RootCommand
		RootCmdOpt() (root *RootCmdOpt)

		ToCommand() *Command

		AddOptFlag(flag OptFlag)
		AddFlag(flag *Flag)
		// AddOptCmd adds 'opt' OptCmd as a sub-command
		AddOptCmd(opt OptCmd)
		// AddCommand adds a *Command as a sub-command
		AddCommand(cmd *Command)

		// AttachTo attaches itself as a sub-command of 'opt' OptCmd object
		AttachTo(parentOpt OptCmd) (opt OptCmd)
		// AttachToCommand attaches itself as a sub-command of *Command object
		AttachToCommand(cmd *Command) (opt OptCmd)
		// AttachToRoot attaches itself as a sub-command of *RootCommand object
		AttachToRoot(root *RootCommand) (opt OptCmd)
	}

	// OnSet interface
	OnSet interface {
		// OnSet will be callback'd after this flag parsed
		OnSet(f func(keyPath string, value interface{})) (opt OptFlag)
	}

	// OptFlagType to support fluent api of cmdr.
	// see also: OptCmd.NewFlag(OptFlagType)
	//
	// Usage
	//
	//   root := cmdr.Root()
	//   co := root.NewSubCommand()
	//   co.NewFlag(cmdr.OptFlagTypeUint)
	//
	// See also those short-hand constructors: Bool(), Int(), ....
	OptFlagType int
)

const (
	// OptFlagTypeBool to create a new bool flag
	OptFlagTypeBool OptFlagType = iota
	// OptFlagTypeInt to create a new int flag
	OptFlagTypeInt OptFlagType = iota + 1
	// OptFlagTypeUint to create a new uint flag
	OptFlagTypeUint OptFlagType = iota + 2
	// OptFlagTypeInt64 to create a new int64 flag
	OptFlagTypeInt64 OptFlagType = iota + 3
	// OptFlagTypeUint64 to create a new uint64 flag
	OptFlagTypeUint64 OptFlagType = iota + 4
	// OptFlagTypeFloat32 to create a new int float32 flag
	OptFlagTypeFloat32 OptFlagType = iota + 8
	// OptFlagTypeFloat64 to create a new int float64 flag
	OptFlagTypeFloat64 OptFlagType = iota + 9
	// OptFlagTypeComplex64 to create a new int complex64 flag
	OptFlagTypeComplex64 OptFlagType = iota + 10
	// OptFlagTypeComplex128 to create a new int complex128 flag
	OptFlagTypeComplex128 OptFlagType = iota + 11
	// OptFlagTypeString to create a new string flag
	OptFlagTypeString OptFlagType = iota + 12
	// OptFlagTypeStringSlice to create a new string slice flag
	OptFlagTypeStringSlice OptFlagType = iota + 13
	// OptFlagTypeIntSlice to create a new int slice flag
	OptFlagTypeIntSlice OptFlagType = iota + 14
	// OptFlagTypeInt64Slice to create a new int slice flag
	OptFlagTypeInt64Slice OptFlagType = iota + 15
	// OptFlagTypeUint64Slice to create a new int slice flag
	OptFlagTypeUint64Slice OptFlagType = iota + 16
	// OptFlagTypeDuration to create a new duration flag
	OptFlagTypeDuration OptFlagType = iota + 17
	// OptFlagTypeHumanReadableSize to create a new human readable size flag
	OptFlagTypeHumanReadableSize OptFlagType = iota + 18
)

type optContext struct {
	current     *Command
	root        *RootCommand
	workingFlag *Flag
	temp        *Command
}

var optCtx *optContext

// Root for fluent api, to create a new [*RootCmdOpt] object.
func Root(appName, version string) (opt *RootCmdOpt) {
	root := &RootCommand{AppName: appName, Version: version, Command: Command{BaseOpt: BaseOpt{Name: appName}}}
	// rootCommand = root
	opt = RootFrom(root)
	return
}

// RootFrom for fluent api, to create the new [*RootCmdOpt] object from an existed [RootCommand]
func RootFrom(root *RootCommand) (opt *RootCmdOpt) {
	optCtx = &optContext{current: &root.Command, root: root, workingFlag: nil}

	opt = &RootCmdOpt{optCommandImpl: optCommandImpl{working: optCtx.current}}
	opt.parent = opt
	return
}
