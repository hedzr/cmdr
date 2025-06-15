package hs

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2/cli"

	"github.com/hedzr/is/exec"
	"github.com/hedzr/is/term"
	"github.com/hedzr/is/term/color"

	logz "github.com/hedzr/logg/slog"
)

//
// hs - help system - not yet, todo
//

func New(w cli.Runner, cmd cli.Cmd, args []string) *HelpSystem {
	s := &HelpSystem{worker: w, cmd: cmd, args: args}
	return s
}

type HelpSystem struct {
	worker cli.Runner
	cmd    cli.Cmd
	args   []string
}

// for _,arg:=range args{
// 	//
// }

func (s *HelpSystem) Run(ctx context.Context) (err error) {
	var dfn func()
	if dfn, err = term.MakeRawWrapped(); err != nil {
		return
	}
	defer dfn()

	var welcomeString = color.New().StripLeftTabsColorful(`
	Type 'help' to print Help Screen, 'help cmd...' for a specified cmd.
	Type 'quit' to end this session and back to Shell.
	`).Build()

	if dfn, err = term.MakeNewTerm(ctx, welcomeString, promptString, replyPrefix, s.helpSystemLooper); err != nil {
		return
	}
	defer dfn()

	return
}

func (s *HelpSystem) helpSystemLooper(ctx context.Context, tty term.SmallTerm, replyPrefix string, exitChan <-chan struct{}, closer func()) (err error) {
	defer func() {
		_, _ = fmt.Fprintln(tty, byeString)
		closer()
		// _, _ = fmt.Println("end")
	}()

	var line string
	for {
		line, err = tty.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return
		}
		if line == "" {
			continue
		}
		if line == quitCmd || line == exitCmd || line == "q" {
			break
		}
		_, _ = fmt.Fprintln(tty, replyPrefix, line)
		err = s.interpretCommand(ctx, line, tty)
		if err != nil {
			return
		}
		if ch := ctx.Done(); ch != nil {
			select {
			case <-ch:
				err = ctx.Err()
				return
			case <-exitChan:
				return
			default:
			}
		} else {
			select {
			case <-exitChan:
				return
			default:
			}
		}
	}
	return
}

func (s *HelpSystem) interpretCommand(ctx context.Context, line string, term io.Writer) (err error) {
	a := exec.SplitCommandString(line, '\'')
	// _, _ = fmt.Fprintln(term, a)
	// logz.InfoContext(ctx, "cmd line", "a", a)
	if len(a) == 0 {
		return
	}
	switch a[0] {
	case "help":
		err = s.helpCmd(ctx, a[1:], term)
	default:
		err = s.runSession(ctx, a, term)
	}
	return
}

func (s *HelpSystem) helpCmd(ctx context.Context, args []string, wr io.Writer) (err error) {
	var handled cli.Cmd
	rootCmd := s.cmd.Root().Cmd
	if handled, err = s.FindCmd(ctx, rootCmd, args); handled == nil {
		err = nil
		ttl := strings.Join(args, ".")
		cc, ff := cli.DottedPathToCommandOrFlag1(ttl, rootCmd)
		if ff != nil {
			_, _ = fmt.Fprintf(wr, `
Flag %v FOUND. It belongs to %v.

`, ff, cc)
			return
		} else if cc == nil {
			_, _ = fmt.Fprintf(wr, "%q command not found (from %v).\n\n", ttl, rootCmd)
			return
		}
		handled = cc.(cli.Cmd)
	}
	// logz.InfoContext(ctx, "cmd found", "cmd", handled)
	var sb strings.Builder
	_, err = s.worker.DoBuiltinAction(ctx, cli.ActionShowHelpScreen, handled, &sb)
	// str := strings.Replace(sb.String(), "\n", "\r\n", -1)
	// _, _ = wr.Write([]byte(str))
	str := sb.String()
	for _, line := range strings.Split(str, "\n") {
		// _, _ = fmt.Fprintf(wr, "%v\r\n", line)
		_, _ = wr.Write([]byte(line))
		_, _ = wr.Write([]byte{'\r', '\n'})
	}

	// var f *os.File
	// f, err = os.Create("2.log")
	// if err != nil {
	// 	return
	// }
	// defer f.Close()
	// _, _ = f.WriteString(str)
	return
}

func (s *HelpSystem) FindCmd(ctx context.Context, cmd cli.Cmd, args []string) (handled cli.Cmd, err error) {
	// trying to recognize the given commands and print help screen of it.
	var cc = cmd.Root().Cmd
	for _, arg := range args {
		cc = cc.FindSubCommand(ctx, arg, true)
		if cc == nil {
			// logz.ErrorContext(ctx, "[cmdr] Unknown command found.", "commands", args)
			handled, err = cc, errors.New("unknown command %v found", args)
			break
		}
	}
	return
}

func (s *HelpSystem) runSession(ctx context.Context, a []string, term io.Writer) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if err == nil {
				if e1, ok := e.(error); ok {
					err = e1
				} else {
					err = fmt.Errorf("%v", e)
				}
			} else {
				err = fmt.Errorf("%v | %v", e, err)
			}
		}
	}()

	return s.runProtectedSession(ctx, a, term)
}

type crlfWriter struct {
	io.Writer
}

func (w *crlfWriter) Write(p []byte) (n int, err error) {
	str := string(p)
	rpl := strings.ReplaceAll(str, "\n", "\r\n")
	return w.Writer.Write([]byte(rpl))
}

func (w *crlfWriter) WriteString(s string) (n int, err error) {
	rpl := strings.ReplaceAll(s, "\n", "\r\n")
	return w.Write([]byte(rpl))
}

func (s *HelpSystem) runProtectedSession(ctx context.Context, a []string, term io.Writer) (err error) {
	savedScreen := struct {
		in  *os.File
		out *os.File
	}{os.Stdin, os.Stdout}
	defer func() {
		os.Stdin, os.Stdout = savedScreen.in, savedScreen.out
	}()

	// _, _ = fmt.Fprintln(term, "Session running...", a)
	logz.SetLevel(logz.DebugLevel)

	err = s.worker.Run(ctx,
		cli.WithArgs(append([]string{os.Args[0]}, a...)...),
		cli.WithHelpScreenWriter(&crlfWriter{term}))
	if err == nil && !s.worker.Error().IsEmpty() {
		err = s.worker.Error()
	}
	if err != nil {
		logz.ErrorContext(ctx, "Error occurred", "err", err)
		err = nil
	}
	return
}

func logOutput() func() {
	logfile := `logfile`
	// open file read/write | create if not exist | clear file at open if exists
	f, _ := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	// save existing stdout | MultiWriter writes to saved stdout and file
	out := os.Stdout
	mw := io.MultiWriter(out, f)

	// get pipe reader and writer | writes to pipe writer come out pipe reader
	r, w, _ := os.Pipe()

	// replace stdout,stderr with pipe writer | all writes to stdout, stderr will go through pipe instead (fmt.print, log)
	os.Stdout = w
	os.Stderr = w

	// writes with log.Print should also write to mw
	logz.Default().AddWriter(mw)

	// create channel to control exit | will block until all copies are finished
	exit := make(chan bool)

	go func() {
		// copy all reads from pipe to multiwriter, which writes to stdout and file
		_, _ = io.Copy(mw, r)
		// when r or w is closed copy will finish and true will be sent to channel
		exit <- true
	}()

	// function to be deferred in main until program exits
	return func() {
		// close writer then block on exit channel | this will let mw finish writing before the program exits
		_ = w.Close()
		<-exit
		// close file after all writes have finished
		_ = f.Close()
	}
}

const (
	promptString = "(cmdr): "
	replyPrefix  = "Human says: "
	quitCmd      = "quit"
	exitCmd      = "exit"
	byeString    = "Bye!"
)
