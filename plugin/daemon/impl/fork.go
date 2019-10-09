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

func forkDaemon(ctx *Context) (err error) {
	if isDemonized() {
		// log.Println("Already a daemon.")
		detachFromTty(ctx.WorkDir, false, true)
		return
	}

	if IsPidFileExists(ctx) {
		log.Printf("Already running or %v file exist.", ctx.PidFileName)
		s, _ := ioutil.ReadFile(ctx.PidFileName)
		pid, err := strconv.ParseInt(string(s), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		ok, err := pidExists(int(pid))
		log.Printf("    pid: %v | %v", pid, ok)
		if !ok && err != nil {
			log.Println("    pidfile removed because it's finished or not useable.")
			removePID(ctx)
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

	savePID(pid, ctx)
	log.Printf("parent process stop itself now.")
	log.Printf("child process is running and detached at: %v", pid)
	os.Exit(0) // parent process exit itself here.
	return
}
