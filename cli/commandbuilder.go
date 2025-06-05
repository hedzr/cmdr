package cli

type OptBuilder interface {
	// Build connects the built command into the building command system.
	Build()
}

type CommandBuilder interface {
	OptBuilder

	// Titles should be specified with this form:
	//
	//     longTitle, shortTitle, aliases...
	//
	// The Long-Title is must-required, and the rest are optional.
	//
	// For Flag, Long-Title and Aliases are POSIX long parameters with the
	// leading double hyphen string '--'. And Short-Title has single
	// hyphen '-' as leading.
	//
	// For example, A flag with longTitle "debug" means that an end-user
	// should type "--debug" for it.
	//
	// For the multi-level command and subcommands (Cmd), long, short and
	// aliases will be used as is.
	Titles(longTitle string, titles ...string) CommandBuilder
	// Description specifies the one-line description and a multi-line
	// description (optional).
	//
	// For the longDescription, some tips are:
	//
	// 1. the extra leading and following blank lines will be stripped.
	// 2. the given params will be joint together as a one big multi-line
	// string.
	Description(description string, longDescription ...string) CommandBuilder
	// Examples can be a multi-line string.
	//
	// Tip: the extra leading and following blank lines will be stripped.
	Examples(examples string) CommandBuilder
	// Group specify a group name,
	// A special prefix could sort it, has a form like `[0-9a-zA-Z]+\.`.
	// The prefix will be removed from help screen.
	//
	// Some examples are:
	//    "A001.Host Params"
	//    "A002.User Params"
	//
	// If ToggleGroup specified, Group field can be omitted because we will copy
	// from there.
	Group(group string) CommandBuilder
	// Deprecated is a version string just like '0.5.9' or 'v0.5.9', that
	// means this command/flag was/will be deprecated since `v0.5.9`.
	Deprecated(deprecated string) CommandBuilder
	// Hidden command/flag won't be shown in help-screen and others output.
	//
	// The Hidden command/flag may be printed normally if very verbose mode
	// specified (typically '-vv' detected).
	//
	// The VendorHidden commands/flags will be hidden at any time even if
	// in vert verbose mode.
	Hidden(hidden bool, vendorHidden ...bool) CommandBuilder

	// ExtraShorts sets more short titles
	ExtraShorts(shorts ...string) CommandBuilder

	// TailPlaceHolders gives two places to place the placeholders.
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
	TailPlaceHolders(placeHolders ...string) CommandBuilder

	// RedirectTo gives the dotted-path to point to a subcommand.
	//
	// Thd target subcommand will be invoked while this command is being invoked.
	//
	// For example, if RootCommand.RedirectTo is set to "build", and
	// entering "app" will equal to entering "app build ...".
	//
	// NOTE:
	//
	//   when redirectTo is valid, CmdS.OnInvoke handler will be ignored.
	RedirectTo(dottedPath string) CommandBuilder

	// OnAction is the main action or entry point when the command
	// was hit from parsing command-line arguments.
	//
	// a call to `OnAction(nil)` will set the underlying onAction handlet empty.
	OnAction(handler OnInvokeHandler) CommandBuilder
	// OnPreAction will be launched before running OnInvoke.
	// The return value obj.ErrShouldStop will cause the remained
	// following processing flow broken right away.
	OnPreAction(handlers ...OnPreInvokeHandler) CommandBuilder
	// OnPostAction will be launched after running OnInvoke.
	OnPostAction(handlers ...OnPostInvokeHandler) CommandBuilder

	// OnMatched _.
	OnMatched(handler OnCommandMatchedHandler) CommandBuilder

	OnEvaluateSubCommands(handler OnEvaluateSubCommands) CommandBuilder
	OnEvaluateSubCommandsOnce(handler OnEvaluateSubCommands) CommandBuilder
	OnEvaluateFlags(handler OnEvaluateFlags) CommandBuilder
	OnEvaluateFlagsOnce(handler OnEvaluateFlags) CommandBuilder

	OnEvaluateSubCommandsFromConfig(path ...string) CommandBuilder

	// PresetCmdLines provides a set of args so that end-user can
	// type the command-line bypass its.
	PresetCmdLines(args ...string) CommandBuilder

	// IgnoreUnmatched specifies whether this command should ignore
	// unmatched flags or not.
	// If set to true, unmatched flags will be treated as positional
	// arguments, and will be added to the command's positionalArgs.
	//
	// If set to false, unmatched flags will be treated as errors,
	// and will cause the command to return an error.
	IgnoreUnmatched(ignore ...bool) CommandBuilder

	// InvokeProc specifies an executable path which will be launched
	// on this command hit and being invoked
	InvokeProc(executablePath string) CommandBuilder
	// InvokeShell specifies a shell command-line which will be launched
	// on this command hit and being invoked.
	//
	// NOTE the command-line string will be launched under the specified
	// shell environment, if it's defined by UseShell().
	InvokeShell(commandLine string) CommandBuilder
	// UseShell specifies a shell environment.
	//
	// It should be a valid path to a shell, such as '/bin/bash',
	// '/bin/zsh', and so on.
	UseShell(shellPath string) CommandBuilder

	//

	// NewCommandBuilder returns a command builder to help you to
	// add a subcommand to the current command builder.
	//
	// It's special because it can be called before Build() completed.
	// NewCommandBuilder(longTitle string, titles ...string) CommandBuilder
	// NewFlagBuilder(longTitle string, titles ...string) FlagBuilder

	// Cmd is a shortcut to NewCommandBuilder and starts a stream
	// building for a new sub-command.
	//
	// It can only be called after current command builder built
	// (Build() called).
	//
	// The right usage of Cmd is:
	//
	//    b := newCommandBuilderShort(parent, "help", "h", "helpme", "info")
	//    nb := b.Cmd("command", "c", "cmd", "commands")
	//    nb.ExtraShorts("cc")
	//    nb.Build()
	//    // ...
	//    b.Build()
	//
	// You may prefer to use AddCmd and Closure:
	//
	//    b := newCommandBuilderShort(parent, "help", "h", "helpme", "info")
	//    b.AddCmd(func(b CommandBuilder){
	//        b.Titles("command", "c", "cmd", "commands")
	//        b.ExtraShorts("cc")
	//    })
	//    // ...
	//    b.Build()
	//
	Cmd(longTitle string, titles ...string) CommandBuilder

	// Flg is a shortcut to NewFlagBuilder and starts a stream
	// building for a new flag.
	//
	// It can only be called after current command builder built
	// (Build() called).
	//
	// The right usage of Flg is:
	//
	//    b := newCommandBuilderShort(parent, "help", "h", "helpme", "info")
	//    nb := b.Flg("use-less", "l", "less", "more")
	//    nb.ExtraShorts("m")
	//    nb.Build()
	//    // ...
	//    b.Build()
	//
	// You may prefer to use AddFlg and Closure:
	//
	//    b := newCommandBuilderShort(parent, "help", "h", "helpme", "info")
	//    b.AddFlg(func(b FlagBuilder){
	//        b.Titles("use-less", "l", "less", "more")
	//        b.ExtraShorts("m")
	//    })
	//    // ...
	//    b.Build()
	//
	Flg(longTitle string, titles ...string) FlagBuilder

	// ToggleableFlags creates a batch of toggleable flags.
	//
	// For example:
	//
	//	s.ToggleableFlags("fruit",
	//	  cli.BatchToggleFlag{L: "apple", S: "a"},
	//	  cli.BatchToggleFlag{L: "banana"},
	//	  cli.BatchToggleFlag{L: "orange", S: "o", DV: true},
	//	)
	ToggleableFlags(toggleGroupName string, items ...BatchToggleFlag)

	// With starts a closure to help you make changes on this builder.
	//
	// For example,
	//
	//    app := cmdr.New()
	//        Info("tiny-app", "0.3.1").
	//        Author("hedzr")
	//    app.Cmd("jump").
	//        Description("jump command").
	//        Examples(`jump example`).
	//        Deprecated(`v1.1.0`).
	//        // Group(cli.UnsortedGroup).
	//        Hidden(false).
	//        With(func(b cli.CommandBuilder) {
	//            b.Cmd("to").
	//                Description("to command").
	//                Examples(``).
	//                Deprecated(`v0.1.1`).
	//                // Group(cli.UnsortedGroup).
	//                Hidden(false).
	//                OnAction(func(cmd *cli.CmdS, args []string) (err error) {
	//                    app.Store().Set("app.demo.working", dir.GetCurrentDir())
	//                    println()
	//                    println(dir.GetCurrentDir())
	//                    println()
	//                    println(app.Store().Dump())
	//                    return // handling command action here
	//                }).
	//                With(func(b cli.CommandBuilder) {
	//                    b.Flg("full", "f").
	//                        Default(false).
	//                        Description("full command").
	//                        // Group(cli.UnsortedGroup).
	//                        Build()
	//                })
	//           b.Flg("dry-run", "n").
	//                Default(false).
	//                Description("run all but without committing").
	//                Group(cli.UnsortedGroup).
	//                Build()
	//        })
	With(cb func(b CommandBuilder))

	// WithSubCmd starts a closure to build a new sub-command
	// and its children.
	//
	// After the closure invoked, new command's Build() will be called
	// implicitly.
	//
	// It can only be called after current command builder built
	// (Build() called).
	//
	// WithSubCmd is a extension of Build.
	WithSubCmd(cb func(b CommandBuilder))

	// // AddCmd starts a closure to build a new sub-command and its children.
	// // After the closure invoked, new command's Build() will be called
	// // implicitly.
	// //
	// // It can only be called after current command builder built
	// // (Build() called).
	// //
	// // Deprecated v2.1.0
	// AddCmd(func(b CommandBuilder)) CommandBuilder
	// // AddFlg starts a closure to build a new flag.
	// // After the closure invoked, new flag's Build() will be
	// // called implicitly.
	// //
	// // It can only be called after current command builder built
	// // (Build() called).
	// //
	// // Deprecated v2.1.0
	// AddFlg(cb func(b FlagBuilder)) CommandBuilder
}
