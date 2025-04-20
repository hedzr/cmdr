package worker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/hedzr/cmdr/v2/cli"
)

//
// -----------------------------------------------------
//

type genS struct{}

func (w *genS) onAction(ctx context.Context, cmd cli.Cmd, args []string) (err error) { //nolint:revive,unused
	outDir := cmd.Store().MustString("dir")
	fmt.Printf("# generating (output-dir: %s) ...\n", outDir)
	return
}

//
// -----------------------------------------------------
//

type genShS struct{}

func (w *genShS) onAction(ctx context.Context, cmd cli.Cmd, args []string) (err error) { //nolint:revive,unused
	outDir := cmd.Store().MustString("dir")
	outputFilename := cmd.Store().MustString("output")
	auto := cmd.Store().MustBool("auto")
	shellMode := cmd.Store().MustString("Shell")
	whichShell := w.gsWhat(cmd, shellMode)

	fmt.Printf("# generating shell autocompletion script (output-dir: %s, file: %s, all: %v, mode: %v, whichShell: %s) ...\n", outDir, outputFilename, auto, shellMode, whichShell)

	var filePath string
	if outDir != "" {
		filePath = path.Join(outDir, outputFilename)
	}

	if auto {
		if g, ok := w.lazyGetGenMaps()[whichShell]; ok {
			err = g(ctx, nil, filePath, cmd, args)
		} else {
			err = w.genShellBash(ctx, nil, filePath, cmd, args)
		}
		return
	}

	var writer io.Writer
	if outputFilename != "" {
		var f *os.File
		if f, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0o644); err == nil {
			ww := bufio.NewWriter(f)
			defer func() {
				err = ww.Flush()
				err = f.Close()
			}()
			writer = ww
		} else {
			return
		}

		if g, ok := w.lazyGetGenMaps()[whichShell]; ok {
			err = g(ctx, writer, filePath, cmd, args)
		} else {
			err = w.genShellBash(ctx, writer, filePath, cmd, args)
		}
	}
	return
}

func (w *genShS) gsWhat(cmd cli.Cmd, shellMode string) (what string) {
	if shellMode == "" || shellMode == "auto" {
		shell := os.Getenv("SHELL")
		switch {
		case strings.HasSuffix(shell, "/zsh"):
			what = "zsh"
		case strings.HasSuffix(shell, "/bash"):
			what = "bash"
		default:
			what = path.Base(shell)
		}
	} else {
		what = shellMode
	}
	_ = cmd
	return
}

var onceShGen sync.Once

type shGenerator func(ctx context.Context, writer io.Writer, fullPath string, cmd cli.Cmd, args []string) (err error)

var shGenMaps map[string]shGenerator

func (w *genShS) lazyGetGenMaps() (m map[string]shGenerator) {
	onceShGen.Do(func() {
		shGenMaps = map[string]shGenerator{
			"bash":       w.genShellBash,
			"zsh":        w.genShellZsh,
			"fish":       w.genShellFish,
			"powershell": w.genShellPowershell,
			"fig":        w.genShellFig,
			"elvish":     w.genShellElvish,
		}
	})
	return shGenMaps
}

func (w *genShS) genShellZsh(ctx context.Context, writer io.Writer, fullPath string, cmd cli.Cmd, args []string) (err error) {
	var gen genzsh
	err = gen.Generate(ctx, writer, fullPath, cmd, args)
	return
}

func (w *genShS) genShellElvish(ctx context.Context, writer io.Writer, fullPath string, cmd cli.Cmd, args []string) (err error) {
	fmt.Println(`# todo elvish`)
	return
}

func (w *genShS) genShellFig(ctx context.Context, writer io.Writer, fullPath string, cmd cli.Cmd, args []string) (err error) {
	fmt.Println(`# todo fig`)
	return
}
