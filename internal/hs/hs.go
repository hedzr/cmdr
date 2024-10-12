package hs

import (
	"context"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/hedzr/cmdr/v2/cli"
)

//
// hs - help system - not yet, todo
//

func New(w cli.Runner, cmd *cli.Command, args []string) *HelpSystem {
	s := &HelpSystem{worker: w, cmd: cmd, args: args}
	return s
}

type HelpSystem struct {
	worker cli.Runner
	cmd    *cli.Command
	args   []string
}

// for _,arg:=range args{
// 	//
// }

func (s *HelpSystem) Run(ctx context.Context) (err error) {
	if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
		return fmt.Errorf("stdin/stdout should be terminal")
	}

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return err
	}
	defer terminal.Restore(0, oldState)
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term := terminal.NewTerminal(screen, "Fxx: ")
	term.SetPrompt(string(term.Escape.Red) + "Fxx: " + string(term.Escape.Reset))

	rePrefix := string(term.Escape.Cyan) + "Human says:" + string(term.Escape.Reset)

	for {
		line, err := term.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if line == "" {
			continue
		}
		if line == "quit" || line == "exit" || line == "q" {
			fmt.Println("Bye!")
			return nil
		}
		fmt.Fprintln(term, rePrefix, line)
	}
	return
}
