/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/hedzr/log/dir"
)

func genShell(cmd *Command, args []string) (err error) {
	w := internalGetWorker()

	var writer io.Writer
	filename := GetStringP(w.getPrefix(), "generate.shell.output")
	filePath := filename
	if filename != "" {
		dirname := GetStringP(w.getPrefix(), "generate.shell.dir")
		if dirname != "" {
			filePath = path.Join(dirname, filename)
		}
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
	}

	what := w.gsWhat(cmd)
	if g, ok := w.lazyGetGenMaps()[what]; ok {
		err = g(writer, filePath, cmd, args)
	} else {
		err = w.genShellBash(writer, filePath, cmd, args)
	}

	return
}

var onceShGen sync.Once

type shGenerator func(writer io.Writer, fullPath string, cmd *Command, args []string) (err error)

var shGenMaps map[string]shGenerator

func (w *ExecWorker) lazyGetGenMaps() (m map[string]shGenerator) {
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

func (w *ExecWorker) gsWhat(cmd *Command) (what string) {
	what = GetStringRP(cmd.GetDottedNamePath(), shTypeGroup, "")
	if what == "" || what == "auto" {
		shell := os.Getenv("SHELL")
		switch {
		case strings.HasSuffix(shell, "/zsh"):
			what = "zsh"
		case strings.HasSuffix(shell, "/bash"):
			what = "bash"
		default:
			what = path.Base(shell)
		}
	}
	return
}

// findDepth returns the depth of a command. rootCommand's deep = 1.
func findDepth(cmd *Command) (deep int) {
	deep = 1
	if cmd.owner != nil {
		deep += findDepth(cmd.owner)
	}
	return
}

// func findLvl(cmd *Command, lvl int) (lvlMax int) {
// 	lvlMax = lvl + 1
// 	for _, cc := range cmd.SubCommands {
// 		l := findLvl(cc, lvl+1)
// 		if l > lvlMax {
// 			lvlMax = l
// 		}
// 	}
// 	return
// }

func (w *ExecWorker) genShellZsh(writer io.Writer, fullPath string, cmd *Command, args []string) (err error) {
	var gen genzsh
	err = gen.Generate(writer, fullPath, cmd, args)
	return
}

//
//
//

func (w *ExecWorker) genShellElvish(writer io.Writer, fullPath string, cmd *Command, args []string) (err error) {
	fmt.Println(`# todo elvish`)
	return
}

func (w *ExecWorker) genShellFig(writer io.Writer, fullPath string, cmd *Command, args []string) (err error) {
	fmt.Println(`# todo fig`)
	return
}

//
//
// /////////////////////////////////////////
//
//

type gensh struct {
	shell         bool
	ext           string
	tplm          map[whatTpl]string
	getTargetPath func(g *gensh) string
	detectShell   func(g *gensh)
	homeDir       string
	shConfigDir   string
	fullPath      string
	appName       string
	endingText    string
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

func (g *gensh) Generate(writer io.Writer, fullPath string, cmd *Command, args []string) (err error) {
	// log.Printf("fullPath: %v, args: %v", fullPath, args)
	if fullPath == "" && len(args) > 0 {
		for _, a := range args {
			if a == "-" {
				err = g.genTo(os.Stdout, cmd, args)
				return
			}
		}

		fullPath = args[0] + "." + g.ext
	}
	g.fullPath, g.appName = fullPath, cmd.root.AppName

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
		err = g.genTo(os.Stdout, cmd, args)
	} else if writer != nil {
		err = g.genTo(writer, cmd, args)
	}
	return
}

func (g *gensh) genTo(writer io.Writer, cmd *Command, args []string) (err error) {
	ctx := &genshCtx{
		cmd: cmd,
		theArgs: &internalShellTemplateArgs{
			RootCommand: cmd.root,
			CmdrVersion: GetString("cmdr.Version"),
			Command:     cmd,
			Args:        args,
		},
		output: writer,
	}

	err = genshTplExpand(ctx, "completion.head", g.tplm[wtHeader], ctx.theArgs)

	if err == nil {
		err = genshTplExpand(ctx, "completion.body", g.tplm[wtBody], ctx.theArgs)
		if err == nil {
			err = genshTplExpand(ctx, "completion.tail", g.tplm[wtTail], ctx.theArgs)

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

//
//
// /////////////////////////////////////////
//
//

func genManualForCommand(cmd *Command) (fn string, err error) {
	painter := newManPainter()
	w := internalGetWorker()
	prefix := strings.Join(append(w.rxxtPrefixes, "generate.manual"), ".")
	dirname := GetStringP(prefix, "dir")
	if err = dir.EnsureDir(dirname); err != nil {
		return
	}

	fn = cmd.root.AppName
	if !cmd.IsRoot() {
		cmds := replaceAll(backtraceCmdNames(cmd, false), ".", "-")
		// if cmds == "generate" {
		// 	cmds += ""
		// }
		if len(cmds) > 0 {
			fn += "-" + cmds
		}
	}
	fn = fmt.Sprintf("%s/%v.1", dirname, fn)

	w.paintFromCommand(painter, cmd, false)
	if err = dir.WriteFile(fn, painter.Results(), 0o600); err == nil {
		log.Printf("%q generated...", fn)
	}
	return
}

func genManual(command *Command, args []string) (err error) {
	w := internalGetWorker()
	painter := newManPainter()
	prefix := strings.Join(append(w.rxxtPrefixes, "generate.manual"), ".")
	// logrus.Debugf("OK gen manual: hit=%v", cmd.strHit)
	// paintFromCommand(newManPainter(), &rootCommand.Command, false)
	err = WalkAllCommands(func(cmd *Command, index, level int) (err error) {
		painter.Reset()

		dirname := GetStringP(prefix, "dir")
		if err = dir.EnsureDir(dirname); err != nil {
			return
		}

		fn := cmd.root.AppName
		if !cmd.IsRoot() {
			cmds := replaceAll(backtraceCmdNames(cmd, false), ".", "-")
			// if cmds == "generate" {
			// 	cmds += ""
			// }
			if len(cmds) > 0 {
				fn += "-" + cmds
			}
		}
		fn = fmt.Sprintf("%s/%v.1", dirname, fn)

		w.paintFromCommand(painter, cmd, false)
		if err = dir.WriteFile(fn, painter.Results(), 0o600); err == nil {
			log.Printf("'%v' generated...", fn)
		}
		return
	})
	return
}

//
//
// /////////////////////////////////////////
//
//

func genDoc(command *Command, args []string) (err error) {
	prefix := strings.Join(append(internalGetWorker().rxxtPrefixes, "generate.doc"), ".")
	// logrus.Infof("OK gen doc: hit=%v", cmd.strHit)
	var painter Painter
	switch command.strHit {
	case "mkd", "m", "markdown":
		painter = newMarkdownPainter()
	case "pdf":
		painter = newMarkdownPainter()
	// case "man", "manual", "manpage", "man-page":
	// 	painter = newManPainter()
	// case "docx":
	// 	painter = newMarkdownPainter()
	// case "tex":
	// 	painter = newMarkdownPainter()
	default: // , "doc", "d"
		if GetBoolP(prefix, "markdown") { //nolint:gocritic //like it
			painter = newMarkdownPainter()
		} else if GetBoolP(prefix, "pdf") { //nolint:gocritic //like it
			painter = newMarkdownPainter()
			// } else if GetBoolP(prefix, "tex") {
			// 	painter = newMarkdownPainter()
		} else { //nolint:gocritic //like it
			painter = newMarkdownPainter()
		}
	}

	// fmt.Printf("  .  . args = [%v]\n", args)
	w := internalGetWorker()
	err = WalkAllCommands(func(cmd *Command, index, level int) (err error) {
		painter.Reset()
		// fmt.Printf("  .  .  cmd = %v\n", cmd.GetTitleNames())

		dirname := GetStringP(prefix, "dir")
		if err = dir.EnsureDir(dirname); err != nil {
			return
		}

		fn := cmd.root.AppName
		if !cmd.IsRoot() {
			cmds := replaceAll(backtraceCmdNames(cmd, false), ".", "-")
			if len(cmds) > 0 {
				fn += "-" + cmds
			}
		}
		fn = fmt.Sprintf("%s/%v.md", dirname, fn)

		w.paintFromCommand(painter, cmd, false)
		if err = dir.WriteFile(fn, painter.Results(), 0o600); err == nil {
			log.Printf("'%v' generated...", fn)
		}
		return
	})

	return
}
