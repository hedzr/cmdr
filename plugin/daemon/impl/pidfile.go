/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// Stop stops the daemon process if running
func Stop(appName string, ctx *Context) {
	present, process := FindDaemonProcess(ctx)
	if present {
		const prefix = "server.stop"
		var err error
		switch {
		case cmdr.GetBoolRP(prefix, "hup"):
			log.Printf("sending SIGHUP to pid %v", process.Pid)
			err = sigSendHUP(process)
		case cmdr.GetBoolRP(prefix, "quit"):
			log.Printf("sending SIGQUIT to pid %v", process.Pid)
			err = sigSendQUIT(process)
		case cmdr.GetBoolRP(prefix, "kill"):
			log.Printf("sending SIGKILL to pid %v", process.Pid)
			err = sigSendKILL(process)
		case cmdr.GetBoolRP(prefix, "usr2"):
			log.Printf("sending SIGUSR2 to pid %v", process.Pid)
			err = sigSendUSR2(process)
		case cmdr.GetBoolRP(prefix, "term"):
			fallthrough
		default:
			log.Printf("sending SIGTERM to pid %v", process.Pid)
			err = sigSendTERM(process)
		}
		if err != nil {
			log.Fatal(err)
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
		log.Printf("sending SIGHUP to pid %v", process.Pid)
		if err := sigSendHUP(process); err != nil {
			log.Fatal(err)
			return
		}
	} else {
		fmt.Printf("%v is stopped.\n", appName)
	}
}

// HotReload reloads the daemon process if running
func HotReload(appName string, ctx *Context) {
	present, process := FindDaemonProcess(ctx)
	if present {
		log.Printf("sending SIGUSR2 to pid %v", process.Pid)
		if err := sigSendUSR2(process); err != nil {
			log.Fatal(err)
			return
		}
		Stop(appName, ctx)
	} else {
		fmt.Printf("%v is stopped.\n", appName)
	}
}

// FindDaemonProcess locates the daemon process if running
func FindDaemonProcess(ctx *Context) (present bool, process *os.Process) {
	if IsPidFileExists(ctx) {
		s, _ := ioutil.ReadFile(ctx.PidFileName)
		pid, err := strconv.ParseInt(string(s), 0, 64)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("cat %v ... pid = %v", ctx.PidFileName, pid)

		process, err = os.FindProcess(int(pid))
		if err == nil {
			present = true
		}
	} else {
		log.Printf("cat %v ... app stopped", ctx.PidFileName)
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
