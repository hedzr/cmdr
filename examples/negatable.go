package examples

import (
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
)

func AddNegatableFlag(parent cli.CommandBuilder)           { AddNegatableFlagImpl(parent) }
func AddNegatableFlagWithoutCmd(parent cli.CommandBuilder) { AddNegatableFlagImpl(parent, true) }

func AddNegatableFlagImpl(parent cli.CommandBuilder, dontHandlingParentCmd ...bool) { //nolint:revive
	common.AddNegatableFlagImpl(parent, dontHandlingParentCmd...)
}
