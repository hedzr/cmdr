package worker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/is/dir"
)

//
//
// /////////////////////////////////////////
//
//

type gensh struct {
	getTargetPath func(g *gensh) string
	detectShell   func(g *gensh)

	tplm map[whatTpl]string
	ext  string

	homeDir     string
	shConfigDir string
	fullPath    string
	appName     string
	endingText  string

	shell bool
}

func (g *gensh) init() {
	g.detectShell(g)
	g.detectShellConfigFolders()
}

func (g *gensh) detectShellConfigFolders() {
	g.homeDir = os.Getenv("HOME") // note that it's available in cmdr system specially for windows since we ever duplicated USERPROFILE as HOME.
	shDir := path.Join(g.homeDir, ".config", g.ext)
	if dir.FileExists(shDir) {
		g.shConfigDir = shDir
	}
}

func (g *gensh) Generate(ctx context.Context, writer io.Writer, fullPath string, cmd cli.Cmd, args []string) (err error) {
	// log.Printf("fullPath: %v, args: %v", fullPath, args)
	if fullPath == "" && len(args) > 0 {
		for _, a := range args {
			if a == "-" {
				err = g.genTo(ctx, os.Stdout, cmd, args)
				return
			}
		}

		fullPath = args[0] + "." + g.ext
	}
	g.fullPath, g.appName = fullPath, cmd.Root().AppName

	if g.shConfigDir != "" && g.fullPath == "" && writer == nil {
		fullPath = g.getTargetPath(g)
		if d := path.Dir(fullPath); !dir.FileExists(d) {
			err = dir.EnsureDir(d)
			if err != nil {
				return
			}
		}
		g.fullPath = fullPath

		var f *os.File
		if f, err = os.Create(g.fullPath); err != nil {
			return
		}
		defer func(f *os.File) {
			err = f.Close()
		}(f)
		writer = f
	}

	if g.fullPath == "" {
		g.fullPath = "-"
		err = g.genTo(ctx, os.Stdout, cmd, args)
	} else if writer != nil {
		err = g.genTo(ctx, writer, cmd, args)
	}
	return
}

func (g *gensh) genTo(ctx context.Context, writer io.Writer, cmd cli.Cmd, args []string) (err error) {
	c := &genshCtx{
		cmd: cmd,
		theArgs: &internalShellTemplateArgs{
			RootCommand: cmd.Root(),
			CmdrVersion: os.Getenv("CMDR_VERSION"),
			Command:     cmd,
			Args:        args,
		},
		output: writer,
	}
	_ = ctx
	err = genshTplExpand(c, "completion.head", g.tplm[wtHeader], c.theArgs)

	if err == nil {
		err = genshTplExpand(c, "completion.body", g.tplm[wtBody], c.theArgs)
		if err == nil {
			err = genshTplExpand(c, "completion.tail", g.tplm[wtTail], c.theArgs)

			if g.fullPath != "-" {
				fmt.Printf(`

# %q generated.`, g.fullPath)
			}

			fmt.Printf(`

%v`, leftPadStr(fmt.Sprintf(g.endingText, g.appName), "# "))
		}
	}

	return
}

func leftPadStr(s, padStr string) string {
	if padStr == "" {
		return s
	}

	var sb strings.Builder
	scanner := bufio.NewScanner(bufio.NewReader(strings.NewReader(s)))
	for scanner.Scan() {
		sb.WriteString(padStr)
		sb.WriteString(scanner.Text())
		sb.WriteRune('\n')
	}
	return sb.String()
}

type whatTpl int

const (
	wtHeader whatTpl = iota
	wtBody
	wtTail
)
