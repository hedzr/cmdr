package cmd

import (
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/cmd"
)

var Commands = append(
	[]cli.CmdAdder{
		sndx{},
	},
	cmd.Commands..., // add commands for testing: jump, wrong
)
