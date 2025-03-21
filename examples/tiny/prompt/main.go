package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"

	"github.com/hedzr/is"
	"github.com/hedzr/is/exec"
	logz "github.com/hedzr/logg/slog"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	err := promptMode()
	if err != nil {
		logz.Fatal("app error", "err", err)
	}
}

func promptMode() (err error) {
	// if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
	// 	return fmt.Errorf("stdin/stdout should be terminal")
	// }

	var oldState *terminal.State
	oldState, err = terminal.MakeRaw(0)
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

		if e1 := terminal.Restore(0, oldState); e1 != nil {
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
	term := terminal.NewTerminal(screen, promptString)
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
		WithOnSignalCaught(func(sig os.Signal, wg *sync.WaitGroup) {
			println()
			logz.Debug("signal caught", "sig", sig)
			// if err := server.Shutdown(); err != nil {
			// 	logger.Error("server shutdown error", "err", err)
			// }
			// cancel()
			exitChan <- struct{}{}
		}).
		Wait(func(stopChan chan<- os.Signal, wgShutdown *sync.WaitGroup) {
			// server.WithOnShutdown(func(err error, ss net.Server) { wgShutdown.Done() })
			// err := server.ListenAndServe(ctx, nil)
			// if err != nil {
			// 	server.Fatal("server serve failed", "err", err)
			// }

			defer func() {
				_, _ = fmt.Fprintln(term, byeString)
				stopChan <- syscall.SIGINT
				wgShutdown.Done()
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
				err = interpretCommand(line, term)
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

func interpretCommand(line string, term *terminal.Terminal) (err error) {
	a := exec.SplitCommandString(line)
	_, _ = fmt.Fprintln(term, a)
	return
}

const (
	promptString = "(cmdr): "
	replyPrefix  = "Human says: "
	quitCmd      = "quit"
	exitCmd      = "exit"
	byeString    = "Bye!"
)
