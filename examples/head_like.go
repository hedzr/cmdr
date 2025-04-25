package examples

import (
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
)

func AddHeadLikeFlag(parent cli.CommandBuilder)           { AddHeadLikeFlagImpl(parent) }
func AddHeadLikeFlagWithoutCmd(parent cli.CommandBuilder) { AddHeadLikeFlagImpl(parent, true) }

func AddHeadLikeFlagImpl(parent cli.CommandBuilder, dontHandlingParentCmd ...bool) { //nolint:revive
	common.AddHeadLikeFlagImpl(parent, dontHandlingParentCmd...)
}
