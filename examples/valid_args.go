package examples

import (
	"github.com/hedzr/cmdr/v2/cli"
)

func AddValidArgsFlag(c cli.CommandBuilder) { //nolint:revive
	c.Flg("fruit", "fr").
		Default("").
		Description("the message.", "").
		Group("Valid Args").
		PlaceHolder("FRUIT").
		ValidArgs("apple", "banana", "orange").
		Build()
}
