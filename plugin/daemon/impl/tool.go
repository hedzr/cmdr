/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/conf"
	"io/ioutil"
	"log"
	"os"
	"path"
	"syscall"
)

func fds(fdin, fdout, fderr uintptr) []uintptr {
	// logFile, logErrFile := nullDev, nullDev
	logDir := cmdr.NormalizeDir(cmdr.GetStringR("logger.dir"))
	if len(logDir) == 0 || logDir == "-" {
		logDir = os.TempDir()
	}

	logFile := path.Join(logDir, fmt.Sprintf("%v.log", conf.AppName))
	logErrFile := logFile
	if cmdr.GetBoolR("logger.splitted") {
		logErrFile = path.Join(logDir, fmt.Sprintf("%v.err.log", conf.AppName))
	}
	log.Printf("logfile: %v", logFile)
	_ = ioutil.WriteFile("/tmp/11", []byte(fmt.Sprintf("%v\n%v\n%v", logFile, cmdr.NormalizeDir(cmdr.GetStringR("logger.dir")), os.TempDir())), 0644)

	fDiscard, e := os.OpenFile(nullDev, os.O_RDWR, 0)
	if e != nil {
		log.Printf("%+v", e)
	} else {
		if fdin == 0 {
			fdin = fDiscard.Fd()
		}
		if fdout == 0 {
			fdout = fDiscard.Fd()
		}
		if fderr == 0 {
			fderr = fDiscard.Fd()
		}

		f, fErr := fDiscard, fDiscard
		if logFile != nullDev {
			f, e = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
			if e != nil {
				log.Printf("%+v", e)
			} else {
				fdout = f.Fd()
				// log.Printf("using logfile: %v\n", logFile)
			}
		}
		if logErrFile != nullDev && logErrFile != logFile {
			fErr, e = os.OpenFile(logErrFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
			if e != nil {
				log.Printf("%+v", e)
			} else {
				fderr = fErr.Fd()
			}
		} else if logErrFile == logFile {
			fderr = fdout
		}
	}

	return []uintptr{fdin, fdout, fderr}
}

func isDemonized() bool {
	return syscall.Getppid() == 1 || os.Getenv(envvarInDaemonized) == "1"
}

func pidExists(pid int) (bool, error) {
	// pid, err := strconv.ParseInt(p, 10, 64)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	process, err := os.FindProcess(int(pid))
	if err != nil {
		// fmt.Printf("Failed to find process: %s\n", err)
		return false, nil
	}

	err = nilSigSend(process)
	log.Printf("process.Signal on pid %d returned: %v\n", pid, err)
	return err == nil, err
}

const nullDev = "/dev/null"
