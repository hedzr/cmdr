package cmdr

import "github.com/hedzr/cmdr"

var optAddTraceOption, optAddServerExtOption cmdr.ExecOption

func init() {
	// attaches `--trace` to root command
	optAddTraceOption = cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
		cmdr.NewBool(false).
			Titles("tr", "trace").
			Description("enable trace mode for tcp/mqtt send/recv data dump", "").
			AttachToRoot(root)
	}, nil)

	// the following statements show you how to attach an option to a sub-command
	optAddServerExtOption = cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
		serverCmd := cmdr.FindSubCommandRecursive("server", nil)
		serverStartCmd := cmdr.FindSubCommand("start", serverCmd)
		cmdr.NewInt(5100).
			Titles("vnc", "vnc-server").
			Description("start as a vnc server (just a demo)", "").
			Placeholder("PORT").
			AttachTo(cmdr.NewCmdFrom(serverStartCmd))
	}, nil)
}
