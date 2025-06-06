package common

import (
	"context"
	"fmt"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/text"
)

func AddSoundexCommand(parent cli.CommandBuilder) {
	parent.Cmd("soundex", "snd", "sndx", "sound").
		Description("soundex test").
		Group("Test").
		TailPlaceHolders("[text1, text2, ...]").
		OnAction(soundex).
		Build()
}

func soundex(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
	_, _ = cmd, args
	for ix, s := range args {
		fmt.Printf("%5d. %s => %s\n", ix, s, text.Soundex(s))
	}
	return
}
