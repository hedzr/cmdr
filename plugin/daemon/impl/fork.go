/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"
)

func forkDaemon() (err error) {
	if isDemonized() {
		// log.Println("Already a daemon.")
		detachFromTty(false, true)
		return
	}

	if IsPidFileExists() {
		log.Printf("Already running or %v file exist.", PIDFile)
		s, _ := ioutil.ReadFile(PIDFile)
		pid, err := strconv.ParseInt(string(s), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		ok, err := pidExists(int(pid))
		log.Printf("    pid: %v | %v", pid, ok)
		if !ok && err != nil {
			log.Println("    pidfile removed because it's finished or not useable.")
			removePID()
		}
		os.Exit(ErrnoForkAndDaemonFailed)
	}

	var pid int
	err = os.Setenv(envvarInDaemonized, "1")
	procAttr := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: fds(os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()),
	}
	pid, _, err = syscall.StartProcess(os.Args[0], os.Args, procAttr)
	if err != nil {
		log.Printf("Fork failed: %+v", err)
		os.Exit(ErrnoForkAndDaemonFailed)
	}

	savePID(pid)
	log.Printf("parent process stop itself now.")
	log.Printf("child process is running and detached at: %v", pid)
	os.Exit(0)
	return
}
