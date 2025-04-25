package common

import (
	"context"
	"testing"

	"github.com/hedzr/cmdr/v2/cli"
)

func TestServerStartup(t *testing.T) {
	var cmd *cli.CmdS
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

func TestXyPrint(t *testing.T) {
	var cmd *cli.CmdS
	ctx := context.Background()
	_ = xyPrint(ctx, cmd, []string{"abcdefg"})
	_ = kbPrint(ctx, cmd, []string{"abcdefg"})
}

func TestSoundex(t *testing.T) {
	var cmd *cli.CmdS
	ctx := context.Background()
	_ = soundex(ctx, cmd, []string{"abcdefg"})
	_ = ttySize(ctx, cmd, []string{})
}
