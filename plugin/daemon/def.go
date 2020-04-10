/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/daemon/impl"
	"net"
	"os"
)

// Daemon interface should be implemented when you are using `daemon.Enable()`.
type Daemon interface {
	OnRun(cmd *cmdr.Command, args []string, stopCh, doneCh chan struct{}, listener net.Listener) (err error)
	OnStop(cmd *cmdr.Command, args []string) (err error)
	OnReload()
	OnStatus(ctx *impl.Context, cmd *cmdr.Command, p *os.Process) (err error)
	OnInstall(ctx *impl.Context, cmd *cmdr.Command, args []string) (err error)
	OnUninstall(ctx *impl.Context, cmd *cmdr.Command, args []string) (err error)
}

// HotReloadable enables hot-restart/hot-reload feature
type HotReloadable interface {
	OnHotReload(ctx *impl.Context) (err error)
}

var daemonImpl Daemon
