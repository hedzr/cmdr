package cli

// StructBuilderOpt represents some options when building command
// system from a struct value via his tag definitions.
type StructBuilderOpt func(s any)

// App object
type App interface {
	// NewCommandBuilder(longTitle string, titles ...string) CommandBuilder // starts a closure to build a new sub-command and its children
	// NewFlagBuilder(longTitle string, titles ...string) FlagBuilder       // starts a closure to build a flag

	// With a callback to bind your logics into app object.
	//
	// For example (from wget/main.go):
	//
	//	app.With(func(app cli.App) {
	//		app.Flg("background", "b", "bg").Description(`go to background after startup`).Group(cStartup).Build()
	//		app.Flg("execute", "e").Description(`execute a '.wgetrc'-style command`).Group(cStartup).Default("").PlaceHolder("COMMAND").Build()
	//		// ...
	//	})
	With(cb func(app App))

	// WithOpts binds more Opt definitions.
	//
	// Sometimes WithOpts can be used to apply couples of opts:
	//
	//	app.WithOpts(cli.WithArgs(os.Args[0], "~~tree"), ...)
	WithOpts(opts ...Opt) App

	// FromStruct constructs the whole app commands system from a given struct value.
	//
	// See also [Creator.BuildFrom].
	//
	// The struct members' tag will be recoganized to build the command system.
	//
	// Sample codes:
	//
	//	app.Cmd("multi", "m", "").
	//		Description("multi-level test and imported form struct").
	//		// Group("Test").
	//		TailPlaceHolders("[text1, text2, ...]").
	//		OnAction(soundex).
	//		FromStruct(root{}). // mount struct fields as commands/flags to callee (Cmd("multi"))
	//		With(func(b cli.CommandBuilder) {
	//			// b.FromStruct(&root{})
	//		})
	//
	//	type root struct {
	//		b   bool // unexported values ignored
	//		Int int  `cmdr:"-"` // ignored
	//		A   `title:"a-cmd" shorts:"a,a1,a2" alias:"a1-cmd,a2-cmd" desc:"A command for demo" required:"true"`
	//		B
	//		C
	//		F1 int
	//		F2 string
	//	}
	//
	//	type A struct {
	//		D
	//		F1 int
	//		F2 string
	//	}
	//	type B struct {
	//	...
	//
	// The fully codes at: https://github.com/hedzr/cmdr#from-a-struct-value
	//
	// See also: https://docs.hedzr.com/docs/cmdr/v2/guide/steps/g13-build-from-struct
	FromStruct(structValue any, opts ...StructBuilderOpt) StructBuilder

	// Cmd is a shortcut to NewCommandBuilder and starts a stream building for a new sub-command.
	Cmd(longTitle string, titles ...string) CommandBuilder
	// Flg is a shortcut to NewFlagBuilder and starts a stream building for a new flag
	Flg(longTitle string, titles ...string) FlagBuilder

	// ToggleableFlags creates a batch of toggleable flags associated with root command.
	//
	// For example:
	//
	//	app.ToggleableFlags("fruit",
	//	  BatchToggleFlag{L: "apple", S: "a"},
	//	  BatchToggleFlag{L: "banana"},
	//	  BatchToggleFlag{L: "orange", S: "o", DV: true},
	//	)
	ToggleableFlags(toggleGroupName string, items ...BatchToggleFlag)

	// // AddCmd starts a closure to build a new sub-command and its children.
	// // After the closure invoked, Build() will be called implicitly.
	// AddCmd(func(b CommandBuilder)) App
	// // AddFlg starts a closure to build a flag
	// // After the closure invoked, Build() will be called implicitly.
	// AddFlg(cb func(b FlagBuilder)) App

	// NewCmdFrom creates a CommandBuilder from 'from' Cmd.
	// For example:
	//
	//	app.NewCmdFrom(p, func(b cli.CommandBuilder) {
	//		b.Titles("version", "ver", "versions").
	//			Description("Show app versions information").
	//			Group(cli.SysMgmtGroup).
	//			Hidden(true, false).
	//			OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	//				w.actionsMatched |= cli.ActionShowVersion
	//				_, err = w.DoBuiltinAction(ctx, w.actionsMatched)
	//				return
	//			})
	//	})
	NewCmdFrom(from *CmdS, cb func(b CommandBuilder)) App
	// NewFlgFrom creates a CommandBuilder from 'from' Cmd.
	// For example:
	//
	//	app.NewFlgFrom(p, false, func(b cli.FlagBuilder) {
	//		b.Titles("version", "V", "ver", "versions").
	//			Description("Show app versions information").
	//			Group(cli.SysMgmtGroup).
	//			Hidden(true, false).
	//			OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) {
	//				// app.Store().Set("app.cmd.show.version", true)
	//				w.actionsMatched |= cli.ActionShowVersion
	//				return
	//			}).
	//			CompCircuitBreak(true).
	//			CompJustOnce(true)
	//	})
	NewFlgFrom(from *CmdS, defaultValue any, cb func(b FlagBuilder)) App

	// RootBuilder starts a session so you can build the subcommands from here.
	//
	// Sample code:
	//
	//    app.RootBuilder(func(parent CommandBuilder){
	//        parent.Cmd("sub1").Build()
	//        parent.Cmd("sub2").Build()
	//        // ...
	//    })
	RootBuilder(cb func(parent CommandBuilder)) App

	// Runner interface
	Runner

	GetRunner() Runner // return the embedded runner object

	Info(name, version string, desc ...string) App // setup basic information about this app
	Copyright(copy string) App                     // setup copyright declaration about this app
	Author(author string) App                      // setup author or team information
	Description(desc string) App                   // setup description and/or long-desc feild
	Header(headerLine string) App                  // setup header line(s) instead of copyright+author fields
	Footer(footerLine string) App                  // setup footer line(s)

	OnAction(handler OnInvokeHandler) App

	// Examples(examples ...string) App               // set examples field of root command

	RootCommand() *RootCommand            // get root command
	SetRootCommand(root *RootCommand) App // setup root command
	WithRootCommand(func(root *RootCommand)) App

	Name() string    // this app name
	Version() string // this app version

	Args() []string // command-line args
}

// BatchToggleFlag for [CommandBuilder].[ToggleableFlags(string, ...BatchToggleFlag)]
type BatchToggleFlag struct {
	L  string // long title
	S  string // short title
	DV bool   // default value
}
