// Copyright © 2019 Hedzr Yeh.

package impl

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"path"
)

// Context of daemon operations
type Context struct {
	// daemon.Context
	PidFileName string
	PidFilePerm int
	LogFileName string
	LogFilePerm int
	WorkDir     string
	Umask       int
	Args        []string
}

// DefaultContext returns a daemon context
func DefaultContext() *Context {
	return &Context{
		PidFileName: "/tmp/daemon.pid",
		PidFilePerm: 0600,
		LogFileName: "/tmp/daemon.log",
		LogFilePerm: 0644,
		WorkDir:     ".",
		Umask:       022,
	}
}

// var defaultContext = &Context{
// 	PidFileName: "/tmp/daemon.pid",
// 	PidFilePerm: 0600,
// 	LogFileName: "/tmp/daemon.log",
// 	LogFilePerm: 0644,
// 	WorkDir:     ".",
// 	Umask:       022,
// }

// var daemonCtx *daemon.Context

// GetContext returns daemon Context object with ref to cmdr options store
func GetContext(cmd *cmdr.Command, args []string) *Context {
	var pidpath, logpath, workdir string

	for _, x := range []string{fmt.Sprintf("/var/log/%s/%%s.log", cmd.GetRoot().AppName), "/tmp/%s.log"} {
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

	xx := fmt.Sprintf("/var/lib/%s", cmd.GetRoot().AppName)
	if cmdr.FileExists(xx) {
		workdir = xx + "/"
	} else {
		workdir = "./"
	}

	return &Context{
		PidFileName: pidpath,
		PidFilePerm: 0644,
		LogFileName: logpath,
		LogFilePerm: 0640,
		WorkDir:     workdir,
		Umask:       027,
		Args:        args,
	}
	// daemonCtx = &daemon.Context{
	// 	PidFileName: pidpath,
	// 	PidFilePerm: 0644,
	// 	LogFileName: logpath,
	// 	LogFilePerm: 0640,
	// 	WorkDir:     workdir,
	// 	Umask:       027,
	// 	Args:        args,
	// }
	// return daemonCtx
}
