/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

// func daemon(nochdir, noclose bool) int {
// 	var ret, ret2 uintptr
// 	var err syscall.Errno
// 
// 	darwin := runtime.GOOS == "darwin"
// 
// 	// already a daemon
// 	if syscall.Getppid() == 1 {
// 		log.Println("Already a daemon.")
// 		return 0
// 	}
// 
// 	// fork off the parent process
// 	ret, ret2, err = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
// 	if err != 0 {
// 		log.Println("CAN'T fork: err = ", err)
// 		return ErrnoForkAndDaemonFailed
// 	}
// 	// failure
// 	if ret2 < 0 {
// 		log.Println("CAN'T fork: ret2 = ", ret2)
// 		os.Exit(ErrnoForkAndDaemonFailed)
// 	}
// 
// 	// handle exception for darwin
// 	if darwin && ret2 == 1 {
// 		log.Println("Darwin: ret2 = 1, so ret := 0")
// 		ret = 0
// 	}
// 
// 	// if we got a good PID, then we call exit the parent process.
// 	if ret > 0 {
// 		log.Println("Forked: parent process exiting")
// 		os.Exit(0)
// 	}
// 
// 	log.Printf("Forked: child running, pid=%v, ppid=%v", os.Getpid(), os.Getppid())
// 
// 	/* Change the file mode mask */
// 	_ = syscall.Umask(0)
// 
// 	// create a new SID for the child process
// 	s_ret, s_errno := syscall.Setsid()
// 	if s_errno != nil {
// 		log.Printf("Error: syscall.Setsid errno: %d", s_errno)
// 		return ErrnoForkAndDaemonFailed
// 	}
// 	if s_ret < 0 {
// 		log.Printf("Error: syscall.Setsid s_ret: %d", s_ret)
// 		return ErrnoForkAndDaemonFailed
// 	}
// 	if !nochdir {
// 		s_errno = os.Chdir("/")
// 	}
// 
// 	if !noclose {
// 		f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
// 		if e == nil {
// 			fd := f.Fd()
// 			s_errno = syscall.Dup2(int(fd), int(os.Stdin.Fd()))
// 			if s_errno == nil {
// 				s_errno = syscall.Dup2(int(fd), int(os.Stdout.Fd()))
// 			}
// 			if s_errno == nil {
// 				s_errno = syscall.Dup2(int(fd), int(os.Stderr.Fd()))
// 			}
// 		}
// 	}
// 
// 	savePID(os.Getpid())
// 	return 0
// }
// 
// // wrong!
// func daemonNew() int {
// 	// already a daemon
// 	if isDemonized() {
// 		log.Println("Already a daemon.")
// 		return 0
// 	}
// 
// 	filePath, _ := filepath.Abs(os.Args[0])
// 	cmd := exec.Command(filePath, os.Args[1:]...)
// 	cmd.Stdin = os.Stdin // 给新进程设置文件描述符，可以重定向到文件中
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	err := cmd.Start() // 开始执行新进程，不等待新进程退出
// 	if err != nil {
// 		log.Printf("Unable to create child process: %+v", err)
// 		os.Exit(ErrnoForkAndDaemonFailed)
// 	}
// 	return 0
// }
