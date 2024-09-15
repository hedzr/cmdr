package examples

import (
	"testing"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/store"
)

func TestServerStartup(t *testing.T) {
	var cmd *cli.Command
	var args []string

	_ = serverStartup(cmd, args)
	_ = serverStop(cmd, args)
	_ = serverShutdown(cmd, args)
	_ = serverRestart(cmd, args)
	_ = serverLiveReload(cmd, args)
	_ = serverInstall(cmd, args)
	_ = serverUninstall(cmd, args)

	_ = serverStatus(cmd, args)
	_ = serverPause(cmd, args)
	_ = serverResume(cmd, args)

	_ = kvBackup(cmd, args)
	_ = kvRestore(cmd, args)

	_ = msList(cmd, args)
	_ = msTagsList(cmd, args)
	_ = msTagsAdd(cmd, args)
	_ = msTagsRemove(cmd, args)
	_ = msTagsModify(cmd, args)
	_ = msTagsToggle(cmd, args)
}

func TestMxTest(t *testing.T) {
	cmd := minimalCmd("mx", cli.WithStore(store.New()))
	cmd.Set().Set("mx-test.stdin", true)

	_ = mxTest(cmd, []string{"abcdefg"})
}

func TestXyPrint(t *testing.T) {
	var cmd *cli.Command
	_ = xyPrint(cmd, []string{"abcdefg"})
	_ = kbPrint(cmd, []string{"abcdefg"})
}

func TestSoundex(t *testing.T) {
	var cmd *cli.Command
	_ = soundex(cmd, []string{"abcdefg"})
	_ = ttySize(cmd, []string{})
}

func TestAttachKvCommand(t *testing.T) {
	_, app := minimalApp(cli.WithStore(store.New()))
	b := app.NewCommandBuilder("x")

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
