package exec

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	"gopkg.in/hedzr/errors.v3"
)

// New return a calling object to allow you to make the fluent call.
//
// Just like:
//
//	exec.New().WithCommand("bash", "-c", "echo hello world!").Run()
//	err = exec.New().WithCommand("bash", "-c", "echo hello world!").RunAndCheckError()
//
// Processing the invoke result:
//
//	exec.New().
//	    WithCommand("bash", "-c", "echo hello world!").
//	    WithStdoutCaught().
//	    WithOnOK(func(retCode int, stdoutText string) { }).
//	    WithStderrCaught().
//	    WithOnError(func(err error, retCode int, stdoutText, stderrText string) { }).
//	    Run()
//
// Use context:
//
//	exec.New().
//	    WithCommandString("bash -c 'echo hello world!'", '\'').
//	    WithContext(context.TODO()).
//	    Run()
//	// or double quote pieces
//	exec.New().
//	    WithCommandString("bash -c \"echo hello world!\"").
//	    Run()
//
// Auto Left-Padding if WithOnError / WithOnOK /
// WithStderrCaught / WithStdoutCaught specified
// (It's no effects when you caught stdout/stderr with
// the handlers above. In this case, do it with
// LeftPad manually).
//
//	args := []string{"-al", "/usr/local/bin"}
//	err := exec.New().
//	    WithPadding(8).
//	    WithCommandArgs("ls", args...).
//	    RunAndCheckError()
func New(opts ...Opt) *calling {
	c := &calling{}

	c.env = append(c.env, os.Environ()...)

	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *calling) Run()                    { _ = c.run() }
func (c *calling) RunAndCheckError() error { return c.run() }

func (c *calling) WithCommandArgs(cmd string, args ...string) *calling {
	c.Cmd = exec.Command(cmd, args...)
	return c
}

// WithCommandString allows split command-line by quote
// characters (default is double-quote).
func (c *calling) WithCommandString(cmd string, quoteChars ...rune) *calling {
	a := SplitCommandString(cmd, quoteChars...)
	c.Cmd = exec.Command(a[0], a[1:]...)
	return c
}

func (c *calling) WithCommand(cmd ...interface{}) *calling {
	var args []string
	for _, a := range cmd[1:] {
		if as, ok := a.([]string); ok {
			args = append(args, as...)
		} else if as, ok := a.(string); ok {
			args = append(args, as)
		} else {
			args = append(args, toStringSimple(a))
		}
	}
	c.Cmd = exec.Command(toStringSimple(cmd[0]), args...)
	return c
}

func (c *calling) WithEnv(key, value string) *calling {
	if key != "" {
		chk := key + "="
		for i, kv := range c.env {
			if strings.HasPrefix(kv, chk) {
				c.env[i] = chk + value
				return c
			}
		}

		c.env = append(c.env, chk+value)
	}
	return c
}

func (c *calling) WithWorkDir(dir string) *calling {
	c.Cmd.Dir = dir
	return c
}

func (c *calling) WithExtraFiles(files ...*os.File) *calling {
	c.Cmd.ExtraFiles = files
	return c
}

func (c *calling) WithContext(ctx context.Context) *calling {
	c.Cmd = exec.CommandContext(ctx, c.Cmd.Path, c.Cmd.Args[1:]...)
	return c
}

func (c *calling) WithStdoutCaught(writer ...io.Writer) *calling {
	for _, w := range writer {
		c.stdoutWriter = w
	}
	c.prepareStdoutPipe()
	return c
}

func (c *calling) WithOnOK(onOK func(retCode int, stdoutText string)) *calling {
	c.onOK = onOK
	return c
}

func (c *calling) WithStderrCaught(writer ...io.Writer) *calling {
	for _, w := range writer {
		c.stderrWriter = w
	}
	c.prepareStderrPipe()
	return c
}

func (c *calling) WithOnError(onError func(err error, retCode int, stdoutText, stderrText string)) *calling {
	c.onError = onError
	return c
}

// WithPadding apply left paddings to stdout/err if no
// WithOnError / WithOnOK / WithStderrCaught / WithStdoutCaught
// specified.
func (c *calling) WithPadding(leftPadding int) *calling {
	c.leftPadding = leftPadding
	return c
}

func (c *calling) WithVerboseCommandLine(verbose bool) *calling {
	c.verbose = verbose
	return c
}

// WithQuietOnError do NOT print error internally
func (c *calling) WithQuietOnError(quiet bool) *calling {
	c.quiet = quiet
	return c
}

type calling struct {
	*exec.Cmd

	err          error
	wg           sync.WaitGroup
	stdoutPiper  io.ReadCloser
	stderrPiper  io.ReadCloser
	stdoutWriter io.Writer
	stderrWriter io.Writer
	env          []string
	leftPadding  int
	verbose      bool // print command-line in std-output dumping

	retCode int
	output  bytes.Buffer
	slurp   bytes.Buffer

	quiet bool

	onOK    func(retCode int, stdoutText string)
	onError func(err error, retCode int, stdoutText, stderrText string)
}

func (c *calling) run() (err error) {

	err = c.runNow()

	var ok, er bool
	if err == nil {
		if c.onOK != nil {
			c.onOK(c.retCode, c.output.String())
			ok = true
		}
	}

	if c.onError != nil && err != nil {
		c.onError(err, c.retCode, c.output.String(), c.slurp.String())
		// er = true
		return
	}

	if !c.quiet {
		if c.output.Len() > 0 && !ok {
			if c.leftPadding > 0 {
				fmt.Print(strings.Repeat(" ", c.leftPadding))
			}
			if c.verbose {
				fmt.Printf("OUTPUT // %v:\n", c.Cmd.Args)
			} else {
				fmt.Print("OUTPUT:\n")
			}
			_, _ = fmt.Fprintf(os.Stdout, "%v\n", leftPad(c.output.String(), c.leftPadding))
		}
		if c.slurp.Len() > 0 && !er && c.retCode != 0 {
			if c.leftPadding > 0 {
				_, _ = fmt.Fprintf(os.Stderr, "%v", strings.Repeat(" ", c.leftPadding))
			}
			if c.verbose {
				fmt.Printf("SLURP // %v:\n", c.Cmd.Args)
			} else {
				fmt.Print("SLURP:\n")
			}
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", leftPad(c.slurp.String(), c.leftPadding))
		}
		if err != nil {
			err = errors.New("system call failed (command-line: %q): %v", c.Args, err)
		}
	}
	return
}

func (c *calling) runNow() error {
	if c.Cmd == nil {
		return errors.New("WithCommand() hasn't called yet.")
	}

	c.Cmd.Env = append(c.Cmd.Env, c.env...)

	// log.Debugf("ENV:\n%v", c.Cmd.Env)

	if (c.onOK != nil || c.leftPadding > 0) && c.stdoutPiper == nil {
		c.prepareStdoutPipe()
	}

	if c.stdoutPiper != nil {
		defer c.stdoutPiper.Close()
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			if c.stdoutWriter != nil {
				if c.stdoutWriter == os.Stdout && c.leftPadding > 0 {
					_, _ = io.Copy(&c.output, c.stdoutPiper)
					b := []byte(leftPad(c.output.String(), c.leftPadding))
					_, _ = c.stdoutWriter.Write(b)
					c.output.Reset()
				} else {
					_, _ = io.Copy(c.stdoutWriter, c.stderrPiper)
				}
			} else {
				_, _ = io.Copy(&c.output, c.stdoutPiper)
			}
		}()
	} else {
		c.Cmd.Stdout = os.Stdout
	}

	if (c.onError != nil || c.leftPadding > 0) && c.stderrPiper == nil {
		c.prepareStderrPipe()
	}

	if c.stderrPiper != nil {
		defer c.stderrPiper.Close()
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			if c.stderrWriter != nil {
				if c.stderrWriter == os.Stderr && c.leftPadding > 0 {
					_, _ = io.Copy(&c.slurp, c.stderrPiper)
					b := []byte(leftPad(c.slurp.String(), c.leftPadding))
					_, _ = c.stderrWriter.Write(b)
					c.slurp.Reset()
				} else {
					_, _ = io.Copy(c.stderrWriter, c.stderrPiper)
				}
			} else {
				_, _ = io.Copy(&c.slurp, c.stderrPiper)
			}
		}()
	} else {
		c.Cmd.Stderr = os.Stderr
	}

	if c.err = c.Cmd.Start(); c.err != nil {
		// Problem while copying stdin, stdout, or stderr
		c.err = fmt.Errorf("failed: %v, cmd: %q", c.err, c.Path)
		return c.err
	}

	// Zero exit status
	// Darwin: launchctl can fail with a zero exit status,
	// so check for emtpy stderr
	if c.Path == "launchctl" {
		slurpText, _ := ioutil.ReadAll(c.stderrPiper)
		if len(slurpText) > 0 && !bytes.HasSuffix(slurpText, []byte("Operation now in progress\n")) {
			c.err = fmt.Errorf("failed with stderr: %s, cmd: %q", slurpText, c.Path)
			return c.err
		}
	}

	// slurp, _ := ioutil.ReadAll(stderr)

	c.wg.Wait()

	if c.err = c.Cmd.Wait(); c.err != nil {
		exitStatus, ok := IsExitError(c.err)
		if ok {
			// Command didn't exit with a zero exit status.
			c.retCode = exitStatus
			c.err = errors.New("%q failed with stderr:\n%v\n  ", c.Path, c.slurp.String()).WithErrors(c.err)
			return c.err
		}

		// An error occurred and there is no exit status.
		// return 0, output, fmt.Errorf("%q failed: %v |\n  stderr: %s", command, err.Error(), slurp)
		c.err = errors.New("%q failed with stderr:\n%v\n  ", c.Path, c.slurp.String()).WithErrors(c.err)
		return c.err
	}

	return nil
}

func (c *calling) prepareStdoutPipe() {
	c.stdoutPiper, c.err = c.Cmd.StdoutPipe()
	if c.err != nil {
		// Failed to connect pipe
		c.err = fmt.Errorf("failed to connect stdout pipe: %v, cmd: %q", c.err, c.Path)
	}
}

func (c *calling) prepareStderrPipe() {
	c.stderrPiper, c.err = c.Cmd.StderrPipe()
	if c.err != nil {
		// Failed to connect pipe
		c.err = fmt.Errorf("failed to connect stderr pipe: %v, cmd: %q", c.err, c.Path)
	}
}

func (c *calling) OutputText() string { return c.output.String() }
func (c *calling) SlurpText() string  { return c.slurp.String() }
func (c *calling) RetCode() int       { return c.retCode }
func (c *calling) Error() error       { return c.err }

type Opt func(*calling)

func WithCommandArgs(cmd string, args ...string) Opt {
	return func(c *calling) {
		c.WithCommandArgs(cmd, args...)
	}
}

func WithCommandString(cmd string) Opt {
	return func(c *calling) {
		c.WithCommandString(cmd)
	}
}

func WithCommand(cmd ...interface{}) Opt {
	return func(c *calling) {
		c.WithCommand(cmd...)
	}
}

func WithEnv(key, value string) Opt {
	return func(c *calling) {
		c.WithEnv(key, value)
	}
}

func WithWorkDir(dir string) Opt {
	return func(c *calling) {
		c.WithWorkDir(dir)
	}
}

func WithExtraFiles(files ...*os.File) Opt {
	return func(c *calling) {
		c.WithExtraFiles(files...)
	}
}

func WithContext(ctx context.Context) Opt {
	return func(c *calling) {
		c.WithContext(ctx)
	}
}

func WithStdoutCaught(writer ...io.Writer) Opt {
	return func(c *calling) {
		c.WithStdoutCaught(writer...)
	}
}

func WithOnOK(onOK func(retCode int, stdoutText string)) Opt {
	return func(c *calling) {
		c.WithOnOK(onOK)
	}
}

func WithStderrCaught(writer ...io.Writer) Opt {
	return func(c *calling) {
		c.WithStderrCaught(writer...)
	}
}

func WithOnError(onError func(err error, retCode int, stdoutText, stderrText string)) Opt {
	return func(c *calling) {
		c.WithOnError(onError)
	}
}

func WithPadding(leftPadding int) Opt {
	return func(c *calling) {
		c.WithPadding(leftPadding)
	}
}

func WithVerboseCommandLine(verbose bool) Opt {
	return func(c *calling) {
		c.WithVerboseCommandLine(verbose)
	}
}

// WithQuietOnError do NOT print error internally
func WithQuietOnError(quiet bool) Opt {
	return func(c *calling) {
		c.WithQuietOnError(quiet)
	}
}
