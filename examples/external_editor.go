package examples

import (
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
)

func AddExternalEditorFlag(c cli.CommandBuilder) { //nolint:revive
	common.AddExternalEditorFlag(c)
}
