package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"

	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/is"
	"github.com/hedzr/is/exec"
	"github.com/hedzr/is/term"
	"github.com/hedzr/is/term/color"

	origterm "golang.org/x/term"
)

var legacy bool

func init() {
	flag.BoolVar(&legacy, "legacy", false, "legacy mode")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	var err error
	if legacy {
		err = promptMode(ctx)
	} else {
		err = improvedPromptMode(ctx)
	}
	if err != nil {
		logz.Error("app error", "err", err)
	}
}

func improvedPromptMode(ctx context.Context) (err error) {
	var dfn, dfn2 func()
	if dfn, err = term.MakeRawWrapped(); err != nil {
		return
	}
	defer dfn()

	var welcomeString = color.New().StripLeftTabsColorful(`
	Type 'help' to print Help Screen, 'help cmd...' for a specified cmd.
	Type 'quit' to end this session and back to Shell.
	`).Build()

	if dfn2, err = term.MakeNewTerm(ctx, welcomeString, promptString, replyPrefix, helpSystemLooper); err != nil {
		return
	}
	defer dfn2()

	return
}

func helpSystemLooper(ctx context.Context, tty term.SmallTerm, replyPrefix string, exitChan <-chan struct{}, closer func()) (err error) {
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
		err = interpretCommand(ctx, line, tty)
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

func interpretCommand(ctx context.Context, line string, term io.Writer) (err error) {
	a := exec.SplitCommandString(line, '\'')
	// _, _ = fmt.Fprintln(term, a)
	// logz.InfoContext(ctx, "cmd line", "a", a)
	if len(a) == 0 {
		return
	}
	switch a[0] {
	case "help":
		// err = s.helpCmd(ctx, a[1:], term)
		_, _ = fmt.Fprintln(term, "command ran ok:", a)
	default:
		err = runSession(ctx, a, term)
		// println("command ran ok:", a)
		_, _ = fmt.Fprintln(term, "command ran ok:", a)
	}
	return
}

func runSession(ctx context.Context, a []string, term io.Writer) (err error) {
	_, _, _ = ctx, a, term
	// panic("unimplemented")
	return
}

func promptMode(ctx context.Context) (err error) {
	// if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
	// 	return fmt.Errorf("stdin/stdout should be terminal")
	// }

	var oldState *origterm.State
	oldState, err = term.MakeRaw(0)
	if err != nil {
		if !errors.Is(err, syscall.ENOTTY) {
			return err
		}
	}
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

		if e1 := term.Restore(0, oldState); e1 != nil {
			if err == nil {
				err = e1
			} else {
				err = fmt.Errorf("%v | %v", e1, err)
			}
		}
	}()
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term := origterm.NewTerminal(screen, promptString)
	term.SetPrompt(string(term.Escape.Red) + promptString + string(term.Escape.Reset))

	rePrefix := string(term.Escape.Cyan) + replyPrefix + string(term.Escape.Reset)

	exitChan := make(chan struct{}, 3)
	defer func() { close(exitChan) }()

	catcher := is.Signals().Catch()
	catcher.
		WithVerboseFn(func(msg string, args ...any) {
			// logz.WithSkip(2).Println(fmt.Sprintf("[verbose] %s\n", fmt.Sprintf(msg, args...)))
			// // server.Verbose(fmt.Sprintf("[verbose] %s", fmt.Sprintf(msg, args...)))
		}).
		WithOnSignalCaught(func(ctx context.Context, sig os.Signal, wg *sync.WaitGroup) {
			println()
			logz.Debug("signal caught", "sig", sig)
			// if err := server.Shutdown(); err != nil {
			// 	logger.Error("server shutdown error", "err", err)
			// }
			// cancel()
			exitChan <- struct{}{}
		}).
		WaitFor(ctx, func(ctx context.Context, closer func()) {
			// server.WithOnShutdown(func(err error, ss net.Server) { wgShutdown.Done() })
			// err := server.ListenAndServe(ctx, nil)
			// if err != nil {
			// 	server.Fatal("server serve failed", "err", err)
			// }

			defer func() {
				_, _ = fmt.Fprintln(term, byeString)
				// stopChan <- syscall.SIGINT
				// wgShutdown.Done()
				closer()
				// _, _ = fmt.Println("end")
			}()

			var line string
			for {
				line, err = term.ReadLine()
				if err == io.EOF {
					return
				}
				if err != nil {
					return
				}
				if line == "" {
					continue
				}
				if line == quitCmd || line == exitCmd || line == "q" {
					break
				}
				_, _ = fmt.Fprintln(term, rePrefix, line)
				err = interpretCommand(ctx, line, term)
				if err != nil {
					return
				}
				select {
				case <-exitChan:
					return
				default:
				}
			}
		})
	return
}

// func interpretCommandOld(line string, term *origterm.Terminal) (err error) {
// 	a := exec.SplitCommandString(line)
// 	_, _ = fmt.Fprintln(term, a)
// 	return
// }

const (
	promptString = "(cmdr): "
	replyPrefix  = "Human says: "
	quitCmd      = "quit"
	exitCmd      = "exit"
	byeString    = "Bye!"
)
