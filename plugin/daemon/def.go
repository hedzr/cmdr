/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import (
	"github.com/hedzr/cmdr"
	"os"
)

type Daemon interface {
	OnRun(cmd *cmdr.Command, args []string, stopCh, doneCh chan struct{}) (err error)
	OnStop(cmd *cmdr.Command, args []string) (err error)
	OnReload()
	OnStatus(cxt *Context, cmd *cmdr.Command, p *os.Process) (err error)
	OnInstall(cxt *Context, cmd *cmdr.Command, args []string) (err error)
	OnUninstall(cxt *Context, cmd *cmdr.Command, args []string) (err error)
}

var daemonImpl Daemon
