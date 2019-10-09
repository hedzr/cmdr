/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"log"
	"os"
	"strconv"
)

func IsPidFileExists() bool {
	// check if daemon already running.
	if _, err := os.Stat(PIDFile); err == nil {
		return true
		
	}
	return false
}

func removePID() {
	// remove PID file
	_ = os.Remove(PIDFile)
}

func savePID(pid int) {

	file, err := os.Create(PIDFile)
	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(ErrnoForkAndDaemonFailed)
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))
	log.Printf("pid %v written into %v", pid, PIDFile)

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

var PIDFile = "/tmp/daemonize.pid"
