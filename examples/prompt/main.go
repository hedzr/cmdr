package main

import (
	"context"
	"flag"
	"fmt"
	"io"

	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/is/exec"
	"github.com/hedzr/is/term"
	"github.com/hedzr/is/term/color"
)

func main() {
	flag.Parse()

	ctx := context.Background()

	var err error
	err = improvedPromptMode(ctx)
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

	if dfn2, err = term.MakeNewTerm(ctx, &term.PromptModeConfig{
		Name:              "cmdr.v2.examples.prompt",
		WelcomeText:       welcomeString,
		PromptText:        promptString,
		ReplyText:         replyPrefix,
		MainLooperHandler: helpSystemLooper,
		// PostInitTerminal:  postInitTerminal,
	}); err != nil {
		return
	}
	defer dfn2()

	return
}

// func postInitTerminal(t *origterm.Terminal) {}

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

const (
	promptString = "(cmdr): "
	replyPrefix  = "Human says: "
	quitCmd      = "quit"
	exitCmd      = "exit"
	byeString    = "Bye!"
)
