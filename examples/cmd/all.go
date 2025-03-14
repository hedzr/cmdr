package cmd

import (
	"github.com/hedzr/cmdr/v2/cli"
)

var Commands = []cli.CmdAdder{
	jumpCmd{},
	wrongCmd{},
	invokeCmd{},
	presetCmd{},
}
