/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/conf"
)

// Demonize enables the demonized mode for this app.
// It fork a new child process and detach it from linux tty session, and the parent process exit itself.
func Demonize(ctx *Context) (err error) {
	cmdr.Set("APPNAME", conf.AppName)
	err = forkDaemon(ctx)
	return
}
