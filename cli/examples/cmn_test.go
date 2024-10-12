package examples

import (
	"context"
	"testing"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/store"
)

func TestServerStartup(t *testing.T) {
	var cmd *cli.Command
	var args []string

	ctx := context.Background()
	_ = serverStartup(ctx, cmd, args)
	_ = serverStop(ctx, cmd, args)
	_ = serverShutdown(ctx, cmd, args)
	_ = serverRestart(ctx, cmd, args)
	_ = serverLiveReload(ctx, cmd, args)
	_ = serverInstall(ctx, cmd, args)
	_ = serverUninstall(ctx, cmd, args)

	_ = serverStatus(ctx, cmd, args)
	_ = serverPause(ctx, cmd, args)
	_ = serverResume(ctx, cmd, args)

	_ = kvBackup(ctx, cmd, args)
	_ = kvRestore(ctx, cmd, args)

	_ = msList(ctx, cmd, args)
	_ = msTagsList(ctx, cmd, args)
	_ = msTagsAdd(ctx, cmd, args)
	_ = msTagsRemove(ctx, cmd, args)
	_ = msTagsModify(ctx, cmd, args)
	_ = msTagsToggle(ctx, cmd, args)
}

func TestMxTest(t *testing.T) {
	cmd := minimalCmd("mx", cli.WithStore(store.New()))
	cmd.Set().Set("mx-test.stdin", true)
	ctx := context.Background()
	_ = mxTest(ctx, cmd, []string{"abcdefg"})
}

func TestXyPrint(t *testing.T) {
	var cmd *cli.Command
	ctx := context.Background()
	_ = xyPrint(ctx, cmd, []string{"abcdefg"})
	_ = kbPrint(ctx, cmd, []string{"abcdefg"})
}

func TestSoundex(t *testing.T) {
	var cmd *cli.Command
	ctx := context.Background()
	_ = soundex(ctx, cmd, []string{"abcdefg"})
	_ = ttySize(ctx, cmd, []string{})
}

func TestAttachKvCommand(t *testing.T) {
	_, app := minimalApp(cli.WithStore(store.New()))
	b := app.Cmd("x")

	AttachServerCommand(b)
	AttachKvCommand(b)
	AttachMsCommand(b)

	AttachMoreCommandsForTest(b, false)
	AttachMoreCommandsForTest(b, true)
}

func minimalApp(opts ...cli.Opt) (root *cli.RootCommand, app cli.App) {
	app = cmdr.New(opts...)
	root = (&cli.RootCommand{
		Command: &cli.Command{},
	}).SetApp(app)
	return
}

func minimalCmd(longTitle string, opts ...cli.Opt) (cc *cli.Command) {
	root, app := minimalApp(opts...)
	cc = root.NewCmd(longTitle)
	_ = app
	return
}
