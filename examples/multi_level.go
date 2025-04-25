package examples

import (
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
)

func AddMultiLevelTestCommand(parent cli.CommandBuilder) {
	common.AddMultiLevelTestCommand(parent)
}
