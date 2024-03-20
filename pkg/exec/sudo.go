// Copyright © 2020 Hedzr Yeh.

package exec

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"gopkg.in/hedzr/errors.v3"
)

// Run runs an OS command
func Run(command string, arguments ...string) (err error) {
	_, _, err = RunCommand(command, false, arguments...)
	return
}

// Sudo runs an OS command with sudo prefix
func Sudo(command string, arguments ...string) (retCode int, stdoutText string, err error) {
	var sudocmd string
	sudocmd, err = exec.LookPath("sudo")
	if err != nil {
		return -1, "'sudo' not found", Run(command, arguments...)
	}

	retCode, stdoutText, err = RunCommand(sudocmd, true, append([]string{command}, arguments...)...)
	return
}

// RunWithOutput runs an OS command and collect the result outputting
func RunWithOutput(command string, arguments ...string) (retCode int, stdoutText string, err error) {
	return RunCommand(command, true, arguments...)
}

// RunCommand runs an OS command and return outputs
func RunCommand(command string, readStdout bool, arguments ...string) (retCode int, stdoutText string, err error) {
	var errText string
	retCode, stdoutText, errText, err = RunCommandFull(command, readStdout, arguments...)
	if errText != "" {
		stdoutText += errText
	}
	return
}

// RunCommandFull runs an OS command and return the all outputs
func RunCommandFull(command string, readStdout bool, arguments ...string) (retCode int, stdoutText, stderrText string, err error) {
	cmd := exec.Command(command, arguments...)

	var stdout io.ReadCloser
	var stderr io.ReadCloser
	var output, slurp bytes.Buffer
	var wg sync.WaitGroup

	if readStdout {
		// Connect pipe to read Stdout
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			// Failed to connect pipe
			return 0, "", "", fmt.Errorf("%q failed to connect stdout pipe: %v", command, err)
		}

		defer stdout.Close()
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = io.Copy(&output, stdout)
		}()

	} else {
		cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
	}

	// Connect pipe to read Stderr
	stderr, err = cmd.StderrPipe()
	if err != nil {
		// Failed to connect pipe
		return 0, "", "", fmt.Errorf("%q failed to connect stderr pipe: %v", command, err)
	}

	defer stderr.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = io.Copy(&slurp, stderr)
	}()

	// Do not use cmd.Run()
	if err = cmd.Start(); err != nil {
		// Problem while copying stdin, stdout, or stderr
		return 0, "", "", fmt.Errorf("%q failed: %v", command, err)
	}

	// Zero exit status
	// Darwin: launchctl can fail with a zero exit status,
	// so check for emtpy stderr
	if command == "launchctl" {
		slurpText, _ := ioutil.ReadAll(stderr)
		if len(slurpText) > 0 && !bytes.HasSuffix(slurpText, []byte("Operation now in progress\n")) {
			return 0, "", "", fmt.Errorf("%q failed with stderr: %s", command, slurpText)
		}
	}

	// slurp, _ := ioutil.ReadAll(stderr)

	wg.Wait()

	if err = cmd.Wait(); err != nil {
		exitStatus, ok := IsExitError(err)
		if ok {
			// Command didn't exit with a zero exit status.
			return exitStatus, output.String(), slurp.String(), errors.New("%q failed with stderr:\n%v\n  ", command, slurp.String()).WithErrors(err)
		}

		// An error occurred and there is no exit status.
		// return 0, output, fmt.Errorf("%q failed: %v |\n  stderr: %s", command, err.Error(), slurp)
		return 0, output.String(), slurp.String(), errors.New("%q failed with stderr:\n%v\n  ", command, slurp.String()).WithErrors(err)
	}

	// if readStdout {
	//	var out []byte
	//	out, err = ioutil.ReadAll(stdout)
	//	if err != nil {
	//		return 0, "", fmt.Errorf("%q failed while attempting to read stdout: %v", command, err)
	//	} else if len(out) > 0 {
	//		output = string(out)
	//	}
	// }

	return 0, output.String(), slurp.String(), nil
}

// IsEAccess detects whether err is a EACCESS errno or not
func IsEAccess(err error) bool {
	if e, ok := err.(*os.PathError); ok && e.Err == syscall.EACCES {
		return true
	}
	return false
}
