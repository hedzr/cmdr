package examples

import (
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
)

func AddSoundexCommand(parent cli.CommandBuilder) {
	common.AddSoundexCommand(parent)
}
