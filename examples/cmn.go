package examples

import (
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
)

func AttachModifyFlags(bdr cli.CommandBuilder) {
	common.AttachModifyFlags(bdr)
}

func AttachConsulConnectFlags(bdr cli.CommandBuilder) {
	common.AttachConsulConnectFlags(bdr)
}
