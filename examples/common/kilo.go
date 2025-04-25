package common

import (
	"context"

	"github.com/hedzr/cmdr/v2/cli"
)

func AddKilobytesFlag(parent cli.CommandBuilder) {
	parent.Flg("size", "s").
		Default("1k").
		Description("max message size. Valid formats: 2k, 2kb, 2kB, 2KB. Suffixes: k, m, g, t, p, e.", "").
		Group("Kilo Bytes").
		Build()
}

func AddKilobytesCommand(parent cli.CommandBuilder) {
	// kb-print

	cb := parent.Cmd("kb-print", "kb").
		Description("kilobytes test", "test kibibytes' input,\nverbose long descriptions here.").
		Group("Kilo Bytes").
		Examples(`
$ {{.AppName}} kb --size 5kb
  5kb = 5,120 bytes
$ {{.AppName}} kb --size 8T
  8TB = 8,796,093,022,208 bytes
$ {{.AppName}} kb --size 1g
  1GB = 1,073,741,824 bytes
		`).
		OnAction(kbPrint)

	cb.Flg("size", "s").
		Default("1k").
		Description("max message size. Valid formats: 2k, 2kb, 2kB, 2KB. Suffixes: k, m, g, t, p, e.", "").
		// Group("").
		Build()

	cb.Build()
}

func kbPrint(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	// fmt.Printf("Got size: %v (literal: %v)\n\n", cmdr.GetKibibytesR("kb-print.size"), cmdr.GetStringR("kb-print.size"))
	_, _ = cmd, args
	return
}
