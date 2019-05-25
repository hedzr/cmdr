/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/sevlyar/go-daemon"
	"path"
)

type Context struct {
	daemon.Context
}

var daemonCtx *daemon.Context

func getContext(cmd *cmdr.Command, args []string) *daemon.Context {
	var pidpath, logpath string

	for _, x := range []string{"/var/log/%s/%s.log", "/tmp/%s.log"} {
		xx := fmt.Sprintf(x, cmd.GetRoot().AppName)
		if cmdr.FileExists(path.Dir(xx)) {
			logpath = xx
			break
		}
	}

	for _, x := range []string{"/var/run/%s/%s.pid", "/tmp/%s.pid"} {
		xx := fmt.Sprintf(x, cmd.GetRoot().AppName)
		if cmdr.FileExists(path.Dir(xx)) {
			pidpath = xx
			break
		}
	}

	daemonCtx = &daemon.Context{
		PidFileName: pidpath,
		PidFilePerm: 0644,
		LogFileName: logpath,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        args,
	}
	return daemonCtx
}
