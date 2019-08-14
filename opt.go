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
	// 	Action(action func(cmd *Command, args []string) (err error)) (opt Opt)
	// }

	// OptFlag to support fluent api of cmdr.
	// see also: cmdr.Root().NewSubCommand()/.NewFlag()
	OptFlag interface {
		Titles(short, long string, aliases ...string) (opt OptFlag)
		Short(short string) (opt OptFlag)
		Long(long string) (opt OptFlag)
		Aliases(ss ...string) (opt OptFlag)
		Description(oneLine, long string) (opt OptFlag)
		Examples(examples string) (opt OptFlag)
		Group(group string) (opt OptFlag)
		Hidden(hidden bool) (opt OptFlag)
		Deprecated(deprecation string) (opt OptFlag)
		Action(action func(cmd *Command, args []string) (err error)) (opt OptFlag)

		ToggleGroup(group string) (opt OptFlag)
		DefaultValue(val interface{}, placeholder string) (opt OptFlag)
		ExternalTool(envKeyName string) (opt OptFlag)
		ValidArgs(list ...string) (opt OptFlag)
		// HeadLike enables `head -n` mode.
		// min, max will be ignored at this version, its might be impl in the future
		HeadLike(enable bool, min, max int64) (opt OptFlag)

		OnSet(func(keyPath string, value interface{})) (opt OptFlag)
		OwnerCommand() (opt OptCmd)
		SetOwner(opt OptCmd)

		RootCommand() *RootCommand
	}

	// OptCmd to support fluent api of cmdr.
	// see also: cmdr.Root().NewSubCommand()/.NewFlag()
	OptCmd interface {
		Titles(short, long string, aliases ...string) (opt OptCmd)
		Short(short string) (opt OptCmd)
		Long(long string) (opt OptCmd)
		Aliases(ss ...string) (opt OptCmd)
		Description(oneLine, long string) (opt OptCmd)
		Examples(examples string) (opt OptCmd)
		Group(group string) (opt OptCmd)
		Hidden(hidden bool) (opt OptCmd)
		Deprecated(deprecation string) (opt OptCmd)
		Action(action func(cmd *Command, args []string) (err error)) (opt OptCmd)

		// FlagAdd(flg *Flag) (opt OptCmd)
		// SubCommand(cmd *Command) (opt OptCmd)
		PreAction(pre func(cmd *Command, args []string) (err error)) (opt OptCmd)
		PostAction(post func(cmd *Command, args []string)) (opt OptCmd)
		TailPlaceholder(placeholder string) (opt OptCmd)

		NewFlag(typ OptFlagType) (opt OptFlag)
		NewSubCommand() (opt OptCmd)

		OwnerCommand() (opt OptCmd)
		SetOwner(opt OptCmd)

		RootCommand() *RootCommand
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
	// OptFlagTypeString to create a new string flag
	OptFlagTypeString OptFlagType = iota + 5
	// OptFlagTypeStringSlice to create a new string slice flag
	OptFlagTypeStringSlice OptFlagType = iota + 6
	// OptFlagTypeIntSlice to create a new int slice flag
	OptFlagTypeIntSlice OptFlagType = iota + 7
	// OptFlagTypeFloat32 to create a new int float32 flag
	OptFlagTypeFloat32 OptFlagType = iota + 8
	// OptFlagTypeFloat64 to create a new int float64 flag
	OptFlagTypeFloat64 OptFlagType = iota + 9
	// OptFlagTypeDuration to create a new duration flag
	OptFlagTypeDuration OptFlagType = iota + 10
)

type optContext struct {
	current     *Command
	root        *RootCommand
	workingFlag *Flag
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
