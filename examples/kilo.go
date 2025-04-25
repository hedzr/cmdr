package examples

import (
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
)

func AddKilobytesFlag(parent cli.CommandBuilder) {
	common.AddKilobytesFlag(parent)
}

func AddKilobytesCommand(parent cli.CommandBuilder) {
	common.AddKilobytesCommand(parent)
}
