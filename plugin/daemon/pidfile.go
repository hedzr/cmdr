/*
 * Copyright © 2019 Hedzr Yeh.
 */

package daemon

import "github.com/hedzr/cmdr"

type pidFileStruct struct {
}

var pidfile = &pidFileStruct{}

func (pf *pidFileStruct) Create(cmd *cmdr.Command) {
	//
}

func (pf *pidFileStruct) Destroy() {
	//
}

type loggerStruct struct {
}

var logger = &loggerStruct{}

func (l *loggerStruct) Setup(cmd *cmdr.Command) {
	//
}
