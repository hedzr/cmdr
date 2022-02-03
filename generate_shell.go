/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"fmt"
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log/dir"
	"github.com/hedzr/log/exec"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

func genShell(cmd *Command, args []string) (err error) {
	// logrus.Infof("OK gen shell. %v", *cmd)
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
		if f, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644); err == nil {
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
	if g, ok := w.gsGenMaps()[what]; ok {
		err = g(writer, filePath, cmd, args)
	} else {
		err = w.genShellBash(writer, filePath, cmd, args)
	}

	return
}

var onceShGen sync.Once

type shGenerator func(writer io.Writer, fullPath string, cmd *Command, args []string) (err error)

var shGenMaps map[string]shGenerator

func (w *ExecWorker) gsGenMaps() (m map[string]shGenerator) {
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
		if strings.HasSuffix(shell, "/zsh") {
			what = "zsh"
		} else if strings.HasSuffix(shell, "/bash") {
			what = "bash"
		} else {
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
	shell := os.Getenv("SHELL")
	if !strings.Contains(shell, "/bin/zsh") {
		var zsh string
		if _, zsh, err = exec.RunWithOutput("which", "zsh"); err != nil {
			// err = errors.New("Couldn't find zsh installation, please install zsh and try again")
			err = genZshTo(cmd, args, "-", os.Stdout)
			return
		} else {
			shell = zsh
		}
	}

	if fullPath == "" && len(args) > 0 {
		for _, a := range args {
			if a == "-" {
				err = genZshTo(cmd, args, "-", os.Stdout)
				return
			}
		}

		fullPath = args[0]
	}

	// find fpath and write to the target

	_, fpath, _ := exec.RunWithOutput(shell, "-c", `echo $fpath`)
	//Logger.Infof("fpath = %v", fpath)
	//Logger.Infof("ENV:\n%v", os.Environ())
	//
	// /usr/local/share/zsh/site-functions
	// $HOME/.oh-my-zsh/completions
	// $HOME/.oh-my-zsh/functions
	//
	locs := tool.ReverseStringSlice(strings.Split(strings.TrimRight(fpath, "\n"), " "))
	err = _makeFileIn(writer, fullPath, locs, cmd.root.AppName, genShellZshHO(cmd, args))
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

// // not complete
// func genShellB(cmd *Command, args []string) (err error) {
// 	// var sb strings.Builder
// 	// var sbca []strings.Builder
//
// 	// cx := &cmd.GetRoot().Command
// 	// lvl := findLvl(cx, 0)
// 	// sbca = make([]strings.Builder, lvl+1)
//
// 	return
// }
//
// // not complete
// func genShellA(cmd *Command, args []string) (err error) {
// 	var sb strings.Builder
// 	var sbca []strings.Builder
//
// 	cx := &cmd.GetRoot().Command
// 	lvl := findLvl(cx, 0)
// 	sbca = make([]strings.Builder, lvl+1)
//
// 	sb.WriteString(fmt.Sprintf(`#compdef _%v %v
//
// # zsh completion wrapper for %v
// # version: %v
// # deep: %v
// #
// # Copyright (c) 2019-2025 Hedzr Yeh <hedzrz@gmail.com>
// #
//
// __ac() {
// 	local state
// 	typeset -A words
// 	_arguments \
// `,
// 		cmd.GetRoot().AppName, cmd.GetRoot().AppName, cmd.GetRoot().AppName, cmd.GetRoot().Version, lvl))
//
// 	for i := 1; i < lvl; i++ {
// 		sb.WriteString(fmt.Sprintf("\t\t'%d: :->level%d' \\\n", i, i))
// 	}
// 	sb.WriteString(fmt.Sprintf("\t\t'%d: :_files'\n\n\tcase $state in\n", lvl))
//
// 	cx = &cmd.GetRoot().Command
// 	body1, body2 := genShellLoopCommands(cx, 1, sbca)
// 	// sb.WriteString(body1)
// 	// sb.WriteString(body2)
// 	logrus.Debugf("%v,%v", len(body1), len(body2))
// 	for i := 1; i <= lvl; i++ {
// 		sb.WriteString(fmt.Sprintf("\t\tlevel%d)\n\t\t\tcase $words[%d] in\n", i, i))
// 		sb.WriteString(sbca[i].String())
// 		sb.WriteString(fmt.Sprintf("\t\t\t\t*) _arguments '%d: :_files' ;;\n\t\t\tesac\n\t\t;;\n\n", i))
// 	}
//
// 	sb.WriteString(fmt.Sprintf(`
// 	esac
// }
//
//
// __ac "$@"
//
//
// # Local Variables:
// # mode: Shell-Script
// # sh-indentation: 4
// # indent-tabs-mode: nil
// # sh-basic-offset: 4
// # End:
// # vim: ft=zsh sw=4 ts=4 et
//
// `))
//
// 	err = ioutil.WriteFile("_"+cmd.GetRoot().AppName, []byte(sb.String()), 0644)
// 	if err == nil {
// 		logrus.Infof("_%v written.", cmd.GetRoot().AppName)
// 	}
// 	return
// }
//
// func genShellLoopCommands(cmd *Command, level int, sbca []strings.Builder) (scrFlg, scrCmd string) {
// 	var sbCmds, sbFlags strings.Builder
//
// 	sbca[level].WriteString(fmt.Sprintf("\t\t\t\t%v) _arguments '%d: :(%v)' ;;\n",
// 		cmd.GetName(), level, cmd.GetSubCommandNamesBy(" ")))
//
// 	for _, cc := range cmd.SubCommands {
// 		// sbCmds.WriteString(fmt.Sprintf(`%v:::`, cc.Name))
//
// 		// sbFlags.WriteString(fmt.Sprintf("\t\t\t\n"))
//
// 		// '(- *)'{--version,-V}'[display version info]' \
// 		// '(- *)'{--help,-h}'[display help]' \
// 		// '(--background -b)'{--background,-b}'[run in background]' \
// 		// 		if len(cc.Flags) > 0 {
// 		// 			for _, flg := range cc.Flags {
// 		// 				sbFlags.WriteString(fmt.Sprintf(`		'(%v)'{%v}'[%v]' \
// 		// `, eraseMultiWSs(flg.GetTitleFlagNamesBy(" ")), eraseMultiWSs(flg.GetTitleFlagNames()), flg.Description))
// 		// 			}
// 		// 		}
//
// 		if len(cc.SubCommands) > 0 {
// 			a, b := genShellLoopCommands(cc, level+1, sbca)
// 			// sbChild.WriteString(a)
// 			// sbca[level+1].WriteString(fmt.Sprintf("\t\tlevel%d)\n\t\t\tcase $words[%d] in\n", level+1, level+1))
// 			sbca[level+1].WriteString(a)
// 			// sbFlags.WriteString(fmt.Sprintf("\t\t\t\t*) _arguments '%d: :_files' ;;\n\t\t\tesac\n\t\t;;\n", level+1))
// 			logrus.Debugf("level %v \nflgs:\n%v\ncmds:\n%v", level, a, b)
// 		}
// 	}
//
// 	// sbFlags.WriteString(fmt.Sprintf("\t\tlevel%d)\n\t\t\tcase $words[%d] in\n", level+1, level+1))
// 	// sbFlags.WriteString(sbChild.String())
// 	// sbFlags.WriteString(fmt.Sprintf("\t\t\t\t*) _arguments '%d: :_files' ;;\n\t\t\tesac\n\t\t;;\n", level+1))
//
// 	if level == 0 {
// 		// 		scrFlg = fmt.Sprintf(`	_arguments -s -S \
// 		// %v && return 0
// 		//
// 		// `, sbFlags.String())
// 		// 		scrCmd = fmt.Sprintf(`	_alternative \
// 		// %v
// 		//
// 		// `, sbCmds.String())
// 	} else {
// 		scrFlg = sbFlags.String()
// 		scrCmd = sbCmds.String()
// 	}
// 	return
// }

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
	if err = ioutil.WriteFile(fn, painter.Results(), 0644); err == nil {
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
		if err = ioutil.WriteFile(fn, painter.Results(), 0644); err == nil {
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
		if GetBoolP(prefix, "markdown") {
			painter = newMarkdownPainter()
		} else if GetBoolP(prefix, "pdf") {
			painter = newMarkdownPainter()
			// } else if GetBoolP(prefix, "tex") {
			// 	painter = newMarkdownPainter()
		} else {
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
		if err = ioutil.WriteFile(fn, painter.Results(), 0644); err == nil {
			log.Printf("'%v' generated...", fn)
		}
		return
	})

	return
}
