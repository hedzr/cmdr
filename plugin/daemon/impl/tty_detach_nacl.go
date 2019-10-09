// Copyright © 2019 Hedzr Yeh.

package impl

func detachFromTty(nochdir, noclose bool) {
	// /* Change the file mode mask */
	// _ = syscall.Umask(0)
	// 
	// // create a new SID for the child process
	// s_ret, s_errno := syscall.Setsid()
	// if s_errno != nil {
	// 	log.Printf("Error: syscall.Setsid errno: %d", s_errno)
	// 	os.Exit(ErrnoForkAndDaemonFailed)
	// }
	// if s_ret < 0 {
	// 	log.Printf("Error: syscall.Setsid s_ret: %d", s_ret)
	// 	os.Exit(ErrnoForkAndDaemonFailed)
	// }
	// if !nochdir {
	// 	s_errno = os.Chdir("/")
	// }
	// 
	// if !noclose {
	// 	fds := fds(0, 0, 0)
	// 	s_errno = syscall.Dup2(int(fds[0]), int(os.Stdin.Fd()))
	// 	if s_errno == nil {
	// 		s_errno = syscall.Dup2(int(fds[1]), int(os.Stdout.Fd()))
	// 	}
	// 	if s_errno == nil {
	// 		s_errno = syscall.Dup2(int(fds[2]), int(os.Stderr.Fd()))
	// 	}
	// }
}
