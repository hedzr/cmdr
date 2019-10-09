/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// Stop stops the daemon process if running
func Stop(appName string, ctx *Context) {
	present, process := FindDaemonProcess(ctx)
	if present {
		if err := sigSendTERM(process); err != nil {
			return
		}
	} else {
		fmt.Printf("%v is stopped.\n", appName)
	}
}

// Reload reloads the daemon process if running
func Reload(appName string, ctx *Context) {
	present, process := FindDaemonProcess(ctx)
	if present {
		if err := sigSendHUP(process); err != nil {
			return
		}
	} else {
		fmt.Printf("%v is stopped.\n", appName)
	}
}

// FindDaemonProcess locates the daemon process if running
func FindDaemonProcess(ctx *Context) (present bool, process *os.Process) {
	if IsPidFileExists(ctx) {
		s, _ := ioutil.ReadFile(ctx.PidFileName)
		pid, err := strconv.ParseInt(string(s), 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		process, err = os.FindProcess(int(pid))
		if err == nil {
			present = true
		}
	}
	return
}

// IsPidFileExists checks if the pid file exists or not
func IsPidFileExists(ctx *Context) bool {
	// check if daemon already running.
	if _, err := os.Stat(ctx.PidFileName); err == nil {
		return true

	}
	return false
}

func removePID(ctx *Context) {
	// remove PID file
	_ = os.Remove(ctx.PidFileName)
}

func savePID(pid int, ctx *Context) {

	file, err := os.Create(ctx.PidFileName)
	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(ErrnoForkAndDaemonFailed)
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))
	log.Printf("pid %v written into %v", pid, ctx.PidFileName)

	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(ErrnoForkAndDaemonFailed)
	}

	err = file.Sync() // flush to disk

	if err != nil {
		log.Printf("Unable to flush pid file : %v\n", err)
		os.Exit(ErrnoForkAndDaemonFailed)
	}
}

// var PIDFile = "/tmp/daemonize.pid"
