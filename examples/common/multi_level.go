package common

import (
	"fmt"

	"github.com/hedzr/cmdr/v2/cli"
)

func AddMultiLevelTestCommand(parent cli.CommandBuilder) {
	cb := parent.Cmd("mls", "mls").
		Description("multi-level subcommands test").
		Group("Test")

	// Sets(func(cmd obj.CommandBuilder) {
	//	cmdrAddFlags(cmd)
	// })

	// cmd := root.NewSubCommand("mls", "mls").
	//	Description("multi-level subcommands test").
	//	Group("Test")
	AddToggleGroupFlags(cb)
	AddValidArgsFlag(cb)

	cmdrMultiLevel(cb, 1)

	cb.Build()
}

func cmdrMultiLevel(parent cli.CommandBuilder, depth int) {
	if depth > 3 {
		return
	}

	for i := 1; i < 4; i++ {
		t := fmt.Sprintf("subcmd-%v", i)
		// var ttls []string
		// for o, l := parent, obj.CommandBuilder(nil); o != nil && o != l; {
		// 	ttls = append(ttls, o.ToCommand().GetTitleName())
		// 	l, o = o, o.OwnerCommand()
		// }
		// ttl := strings.Join(tool.ReverseStringSlice(ttls), ".")
		ttl := ""

		cb := parent.Cmd(t, fmt.Sprintf("sc%v", i)).
			// cc := parent.NewSubCommand(t, fmt.Sprintf("sc%v", i)).
			Description(fmt.Sprintf("subcommands %v.sc%v test", ttl, i)).
			Group("Test")
		AddToggleGroupFlags(cb)
		AddValidArgsFlag(cb)
		cmdrMultiLevel(cb, depth+1)
		cb.Build()
	}
}
