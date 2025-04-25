package common

import (
	"context"
	"testing"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/examples/common"
	"github.com/hedzr/store"
)

func TestMxTest(t *testing.T) {
	cmd := minimalCmd("mx", cli.WithStore(store.New()))
	cmd.Set().Set("mx-test.stdin", true)
	ctx := context.Background()
	_ = common.MxTest(ctx, cmd, []string{"abcdefg"})
}

func TestAttachKvCommand(t *testing.T) {
	_, app := minimalApp(cli.WithStore(store.New()))
	b := app.Cmd("x")

	common.AttachServerCommand(b)
	common.AttachKvCommand(b)
	common.AttachMsCommand(b)

	common.AttachMoreCommandsForTest(b, false)
	common.AttachMoreCommandsForTest(b, true)
}

// func newapp(opts ...cli.Opt) cli.App {
// 	_ = os.Setenv("CMDR_VERSION", Version)
// 	logz.Verbose("setup env-var at earlier time", "CMDR_VERSION", Version)
// 	cfg := cli.NewConfig(opts...)
// 	w := worker.New(cfg)
// 	return builder.New(w)
// }

func minimalApp(opts ...cli.Opt) (root *cli.RootCommand, app cli.App) {
	app = cmdr.New(opts...)
	root = (&cli.RootCommand{
		Cmd: &cli.CmdS{},
	}).SetApp(app)
	return
}

func minimalCmd(longTitle string, opts ...cli.Opt) (cc *cli.CmdS) {
	root, app := minimalApp(opts...)
	cc = root.NewCmd(longTitle)
	_ = app
	return
}
