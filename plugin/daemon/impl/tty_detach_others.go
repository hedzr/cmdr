// +build !windows,!nacl,!plan9,!solaris

// Copyright © 2019 Hedzr Yeh.

package impl

import (
	"log"
	"os"
	"syscall"
)

func detachFromTty(workDir string, nochdir, noclose bool) {
	/* Change the file mode mask */
	_ = syscall.Umask(0)

	// create a new SID for the child process
	sRet, sErrno := syscall.Setsid()
	if sErrno != nil {
		log.Printf("Error: syscall.Setsid errno: %d", sErrno)
		os.Exit(ErrnoForkAndDaemonFailed)
	}
	if sRet < 0 {
		log.Printf("Error: syscall.Setsid sRet: %d", sRet)
		os.Exit(ErrnoForkAndDaemonFailed)
	}
	if !nochdir {
		sErrno = os.Chdir(workDir)
		if sErrno == nil {
			log.Printf("Error: syscall.Setsid errno: %d", sErrno)
		}
	}

	if !noclose {
		fds := fds(0, 0, 0)
		sErrno = syscall.Dup2(int(fds[0]), int(os.Stdin.Fd()))
		if sErrno == nil {
			sErrno = syscall.Dup2(int(fds[1]), int(os.Stdout.Fd()))
		}
		if sErrno == nil {
			sErrno = syscall.Dup2(int(fds[2]), int(os.Stderr.Fd()))
		}
		if sErrno == nil {
			log.Printf("Error: syscall.Setsid errno: %d", sErrno)
		}
	}
}
