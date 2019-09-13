/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"github.com/hedzr/cmdr"
)

func main() {
	// fmt.Println("Hello, playground")

	// // To disable internal commands and flags, uncomment the following codes
	// cmdr.EnableVersionCommands = false
	// cmdr.EnableVerboseCommands = false
	// cmdr.EnableHelpCommands = false
	// cmdr.EnableGenerateCommands = false
	// cmdr.EnableCmdrCommands = false

	rootCmd := buildRootCmd()
	if err := cmdr.Exec(rootCmd,
		// To disable internal commands and flags, uncomment the following codes
		cmdr.WithBuiltinCommands(false, false, false, false, false),
		// daemon.WithDaemon(svr.NewDaemon(), nil, nil, nil),
	); err != nil {
		panic(err)
	}
}

const (
	versionName = "0.0.1"
	appName     = "cmdrtest"
)

func buildRootCmd() (rootCmd *cmdr.RootCommand) {

	root := cmdr.Root(appName, versionName)

	// root.NewSubCommand().
	// 	Titles("h", "help").
	// 	Description("show help screen", "").
	// 	Action(func(cmd *cmdr.Command, args []string) (err error) {
	// 		fmt.Println("this is help text")
	// 		os.Exit(0)
	// 		return
	// 	})

	// root.NewFlag(cmdr.OptFlagTypeBool).
	// 	Titles("h", "help").
	// 	Description("show help screen", "").
	// 	DefaultValue(false, "").
	// 	OnSet(func(keyPath string, value interface{}) {
	// 		fmt.Println("this is help text")
	// 		os.Exit(0)
	// 		return
	// 	})

	rootCmd = root.RootCommand()

	return
}
