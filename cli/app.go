package cli

type App interface {
	// NewCommandBuilder(longTitle string, titles ...string) CommandBuilder // starts a closure to build a new sub-command and its children
	// NewFlagBuilder(longTitle string, titles ...string) FlagBuilder       // starts a closure to build a flag

	With(cb func(app App))

	WithOpts(opts ...Opt) App

	// Cmd is a shortcut to NewCommandBuilder and starts a stream building for a new sub-command
	Cmd(longTitle string, titles ...string) CommandBuilder
	// Flg is a shortcut to NewFlagBuilder and starts a stream building for a new flag
	Flg(longTitle string, titles ...string) FlagBuilder

	// // AddCmd starts a closure to build a new sub-command and its children.
	// // After the closure invoked, Build() will be called implicitly.
	// AddCmd(func(b CommandBuilder)) App
	// // AddFlg starts a closure to build a flag
	// // After the closure invoked, Build() will be called implicitly.
	// AddFlg(cb func(b FlagBuilder)) App

	// NewCmdFrom creates a CommandBuilder from 'from' Cmd.
	NewCmdFrom(from *CmdS, cb func(b CommandBuilder)) App
	// NewFlgFrom creates a CommandBuilder from 'from' Cmd.
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

	Runner

	Info(name, version string, desc ...string) App // setup basic information about this app
	Copyright(copy string) App                     // setup copyright declaration about this app
	Author(author string) App                      // setup author or team information
	Description(desc string) App                   //
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
