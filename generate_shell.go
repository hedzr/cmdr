/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log/exec"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

func genShell(cmd *Command, args []string) (err error) {
	// logrus.Infof("OK gen shell. %v", *cmd)
	w := internalGetWorker()
	what := w.gsWhat()

	switch what {
	case "zsh":
		err = genShellZsh(cmd, args)
	case "fish":
		err = genShellFish(cmd, args)
	case "fig":
		err = genShellFig(cmd, args)
	case "elvish":
		err = genShellElvish(cmd, args)
	case "powershell":
		err = genShellPowershell(cmd, args)
	case "bash":
		fallthrough
	default:
		err = genShellBash(cmd, args)
	}

	//if GetBoolP(w.getPrefix(), "generate.shell.zsh") {
	//	// if !GetBoolP(getPrefix(), "quiet") {
	//	// 	logrus.Debugf("zsh-dump")
	//	// }
	//	// printHelpZsh(command, justFlags)
	//
	//	// not yet
	//} else if GetBoolP(w.getPrefix(), "generate.shell.bash") {
	//	err = genShellBash(cmd, args)
	//} else {
	//	// auto
	//	// shell := os.Getenv("SHELL")
	//	// if strings.HasSuffix(shell, "/bash") || GetBoolP(getPrefix(), "generate.shell.force-bash") {
	//	// 	err = genShellBash(cmd, args)
	//	// } else if strings.HasSuffix(shell, "/zsh") {
	//	// 	// not yet
	//	// }
	//	err = genShellAuto(cmd, args)
	//	// } else {
	//	// 	_, _ = fmt.Fprint(os.Stderr, "Unknown shell. ignored.")
	//	// err = genShellB(cmd, args)
	//}
	return
}

func (w *ExecWorker) gsWhat() (what string) {
	what = "bash"
	if GetBoolP(w.getPrefix(), "generate.shell.zsh") {
		what = "zsh"
	} else if GetBoolP(w.getPrefix(), "generate.shell.elvish") {
		what = "elvish"
	} else if GetBoolP(w.getPrefix(), "generate.shell.fig") {
		what = "fig"
	} else if GetBoolP(w.getPrefix(), "generate.shell.fish") {
		what = "fish"
	} else if GetBoolP(w.getPrefix(), "generate.shell.powershell") {
		what = "powershell"
	} else if !GetBoolP(w.getPrefix(), "generate.shell.bash") {
		shell := os.Getenv("SHELL")
		if strings.HasSuffix(shell, "/zsh") {
			what = "zsh"
		} else if !strings.HasSuffix(shell, "/bash") {
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

func _makeFileIn(locations []string, appName string, fn func(path string, f *os.File) (err error)) (err error) {
	linuxRoot := os.Getuid() == 0
	for _, s := range locations {
		//Logger.Debugf("--- checking %s", s)
		if FileExists(s) {
			file := path.Join(s, "_"+appName)
			//Logger.Debugf("    try creating %s", file)
			var f *os.File
			if f, err = os.Create(file); err != nil {
				if !linuxRoot {
					break // for non-root user, we break file-writing loop and dump scripts to console too.
				}
				return
			}
			Logger.Debugf("    generating %s", file)
			err = fn(file, f)
			//if !linuxRoot {
			//	break // for non-root user, we break file-writing loop and dump scripts to console too.
			//}
			return
		}
	}
	err = fn("-", os.Stdout)
	return
}

func genShellZsh(cmd *Command, args []string) (err error) {
	shell := os.Getenv("SHELL")
	_, fpath, _ := exec.RunWithOutput(shell, "-c", `echo $fpath`)
	//Logger.Infof("fpath = %v", fpath)
	//Logger.Infof("ENV:\n%v", os.Environ())
	//
	// /usr/local/share/zsh/site-functions
	// $HOME/.oh-my-zsh/completions
	// $HOME/.oh-my-zsh/functions
	//
	locs := tool.ReverseStringSlice(strings.Split(strings.TrimRight(fpath, "\n"), " "))
	err = _makeFileIn(locs, cmd.root.AppName, genShellZshHO(cmd))
	return
}

func genShellZshHO(cmd *Command) func(path string, f *os.File) (err error) {
	return func(path string, f *os.File) (err error) {
		ctx := &genZshCtx{
			cmd: cmd,
			args: &internalShellTemplateArgs{cmd.root,
				GetString("cmdr.Version"),
			},
			output: f,
		}

		err = zshTplExpand(ctx, "zsh.completion.head", zshCompHead, ctx.args)
		if err == nil {
			err = walkFromCommand(&cmd.root.Command, 0, 0, func(cx *Command, index, level int) (err error) {
				//generate for command cx
				safeName := getZshSubFuncName(cx)
				if level == 0 {
					safeName = "_" + safeZshFnName(cx.root.AppName)
				}
				err = genZshFnFromTpl(ctx, safeName, cx)
				return
			})

			err = zshTplExpand(ctx, "zsh.completion.tail", zshCompTail, ctx.args)
			fmt.Printf(`# '%v' generated.
# Re-login to enable the new bash completion script.
`, path)
		}
		return
	}
}

func genZshFnFromTpl(ctx *genZshCtx, fnName string, cmd *Command) (err error) {
	if len(cmd.SubCommands) > 0 {
		err = genZshFnByCommand(ctx, fnName, cmd)
	} else {
		//if fnName == "__fluent_generate_shell" {
		//	print()
		//}
		err = genZshFnFlagsByCommand(ctx, fnName, cmd)
	}
	return
}

func genZshFnByCommand(ctx *genZshCtx, fnName string, cmd *Command) (err error) {
	var flags, descCommands, caseCommands strings.Builder

	gzt1(&flags, cmd, true)

	if len(cmd.Flags) > 0 {
		flags.WriteString("               \\\n")
	}

	for _, c := range cmd.SubCommands {
		descCommands.WriteString(fmt.Sprintf("                %v':%v'\n",
			zshDescribeNames(c), c.GetDescZsh()))
		caseCommands.WriteString(fmt.Sprintf(`            %v)
                %v
                ;;
`, safeZshFnName(c.GetTitleNamesBy("|")), getZshSubFuncName(c)))
	}

	err = zshTplExpand(ctx, "zsh.completion.sub-commands", zshCompCommands, struct {
		*internalShellTemplateArgs
		FuncName         string
		FuncNameSeq      string // "__fluent_ms" => "[fluent ms]"
		Flags            string
		DescribeCommands string
		CaseCommands     string
	}{
		ctx.args,
		fnName, seqName(fnName),
		flags.String(),
		strings.TrimSuffix(descCommands.String(), " \\\n"),
		caseCommands.String(),
	})
	return
}

func genZshFnFlagsByCommand(ctx *genZshCtx, fnName string, cmd *Command) (err error) {
	var descCommands strings.Builder

	if len(cmd.Flags) > 0 {
		descCommands.WriteString("    _arguments -s \\\n")
	}

	gzt1(&descCommands, cmd, false)

	desc := strings.TrimSuffix(descCommands.String(), " \\\n")
	decl := ""
	if desc != "" {
		decl = `    typeset -A opt_args
    local context curcontext="$curcontext" state line ret=0
    local I="-h --help --version -V -#"
    local -a commands
`
	}

	err = zshTplExpand(ctx, "zsh.completion.flags", zshCompCommandFlags, struct {
		*internalShellTemplateArgs
		FuncName         string
		Declarations     string
		DescribeCommands string
		CaseCommands     string
	}{
		ctx.args,
		fnName,
		decl,
		desc,
		"",
	})
	return
}

func gzt1(descCommands *strings.Builder, cmd *Command, shortTitleOnly bool) {
	gzt1ForToggleGroups(descCommands, cmd, shortTitleOnly)

	for ix, c := range cmd.Flags {
		if c.ToggleGroup != "" {
			continue
		}
		gzt2(descCommands, cmd, ix, c, "", shortTitleOnly)
	}
}

func gzt1ForToggleGroups(descCommands *strings.Builder, cmd *Command, shortTitleOnly bool) {
	var tgs = make(map[string][]*Flag)
	for _, f := range cmd.Flags {
		if f.ToggleGroup != "" {
			tgs[f.ToggleGroup] = append(tgs[f.ToggleGroup], f)
		}
	}

	for k, v := range tgs {
		me := gzChkMEForToggleGroup(k, v)

		for ix, c := range v {
			//var sb strings.Builder
			//for i, f := range v {
			//	if i != ix {
			//		sb.WriteString(f.GetTitleZshNamesBy(" "))
			//		sb.WriteString(" ")
			//	}
			//}
			gzt2(descCommands, cmd, ix, c, me, shortTitleOnly)
		}
	}
}

func gzt2(descCommands *strings.Builder, cmd *Command, ix int, f *Flag, mutualExclusives string, shortTitleOnly bool) {

	//if c.Full == "pprof" {
	//	println()
	//}

	if len(f.ValidArgs) != 0 {
		gzAction(descCommands, f, "("+strings.Join(f.ValidArgs, " ")+")", mutualExclusives, shortTitleOnly)
	} else if f.DefaultValuePlaceholder == "FILE" {
		act := "_files"
		if f.actionStr != "" {
			act += " -g " + strconv.Quote(f.actionStr)
		}
		gzAction(descCommands, f, act, mutualExclusives, shortTitleOnly)
	} else if f.DefaultValuePlaceholder == "DIR" {
		gzAction(descCommands, f, "_files -/", mutualExclusives, shortTitleOnly)
	} else if f.DefaultValuePlaceholder == "USER" {
		gzAction(descCommands, f, "_users", mutualExclusives, shortTitleOnly)
	} else if f.DefaultValuePlaceholder == "GROUP" {
		gzAction(descCommands, f, "_groups", mutualExclusives, shortTitleOnly)
	} else if f.DefaultValuePlaceholder == "INTERFACES" {
		gzAction(descCommands, f, "_net_interfaces", mutualExclusives, shortTitleOnly)
	} else {
		mutualExclusives = gzChkME(f, mutualExclusives)
		if mutualExclusives != "" {
			descCommands.WriteString(fmt.Sprintf("                \"($I %v)\"%v'[%v]'",
				mutualExclusives, zshDescribeFlagNames(f, shortTitleOnly, false), f.GetDescZsh()))
		} else {
			descCommands.WriteString(fmt.Sprintf("                %v'[%v]'",
				zshDescribeFlagNames(f, shortTitleOnly, false), f.GetDescZsh()))
		}
	}
	//if ix != len(cmd.Flags)-1 {
	descCommands.WriteString(" \\\n")
	//}
}

//func gzAction(descCommands *strings.Builder, c *Flag, action, mutualExclusives string) {
//	gzAction_(descCommands,c, action, mutualExclusives, false)
//}

func gzAction(descCommands *strings.Builder, f *Flag, action, mutualExclusives string, shortTitleOnly bool) {
	if f.dblTildeOnly {
		return
	}

	names := zshDescribeFlagNames(f, shortTitleOnly, false)
	title := f.Full
	if f.DefaultValuePlaceholder != "" {
		title = f.DefaultValuePlaceholder
	}
	mutualExclusives = gzChkME(f, mutualExclusives)
	descCommands.WriteString(fmt.Sprintf("                \"(%v)\"%v'[%v]':%v:'%v'",
		unquote(mutualExclusives), names, f.GetDescZsh(), title, action))
}

// gzChkME checks mutual exclusive flags and builds the leading section for zsh completion system.
// A mutual exclusive section looks like:
//
//      '(--debug -D --quiet -q)'
//
// and the responding optspec will be:
//
//      '(--debug -D --quiet -q)'{--quiet,-q}'[Quiet Mode]'
//      '(--debug -D --quiet -q)'{--debug,-D}'[Debug Mode]'
//
func gzChkME(f *Flag, mutualExclusives string) string {
	const quoted = false
	if mutualExclusives == "" {
		if len(f.mutualExclusives) > 0 {
			var sb strings.Builder
			for _, t := range f.mutualExclusives {
				if tgt, ok := f.owner.plainLongFlags[t]; ok {
					sb.WriteString(tgt.GetTitleZshNamesExtBy(" ", false, quoted, false, false))
				}
			}
			sb.WriteString(f.GetTitleZshNamesExtBy(" ", false, quoted, false, false))
			mutualExclusives = strings.TrimRight(sb.String(), " ")
		} else if f.circuitBreak {
			mutualExclusives = "- *"
		} else {
			mutualExclusives = f.GetTitleZshNamesExtBy(" ", false, quoted, false, false)
		}
	}
	return mutualExclusives
}

func gzChkMEForToggleGroup(toggleGroupName string, v []*Flag) (mutualExclusives string) {
	const quoted = false
	var sb strings.Builder
	for _, f := range v {
		sb.WriteString(f.GetTitleZshNamesExtBy(" ", false, quoted, false, false))
		sb.WriteString(" ")
	}
	mutualExclusives = strings.TrimRight(sb.String(), " ")
	return
}

func unquote(s string) string {
	return regexp.MustCompile(`'(.*?)'`).ReplaceAllString(s, "$1")
}

//
//
//

func genShellFish(cmd *Command, args []string) (err error) {
	fmt.Println(`# todo fish`)
	return
}

func genShellElvish(cmd *Command, args []string) (err error) {
	fmt.Println(`# todo elvish`)
	return
}

func genShellFig(cmd *Command, args []string) (err error) {
	fmt.Println(`# todo fig`)
	return
}

func genShellPowershell(cmd *Command, args []string) (err error) {
	fmt.Println(`# todo powershell`)
	return
}

func genShellBash(cmd *Command, args []string) (err error) {
	var tmpl *template.Template
	tmpl, err = template.New("bash.completion").Parse(`

#

# bash completion wrapper for {{.AppName}}
# version: {{.Version}}
#
# Copyright (c) 2019-2025 Hedzr Yeh <hedzrz@gmail.com>
#

_cmdr_cmd_help_events () {
  $* --help|grep "^  [^ \[\$\#\!/\\@\"']"|awk -F'   ' '{print $1}'|awk -F',' '{for (i=1;i<=NF;i++) print $i}'
}


_cmdr_cmd_{{.AppName}}() {
  local cmd="{{.AppName}}" cur prev words
  _get_comp_words_by_ref cur prev words
  if [ "$prev" != "" ]; then
    unset 'words[${#words[@]}-1]'
  fi

  COMPREPLY=()
  #pre=${COMP_WORDS[COMP_CWORD-1]}
  #cur=${COMP_WORDS[COMP_CWORD]}

  case "$prev" in
    --help|--version)
      COMPREPLY=()
      return 0
      ;;
    $cmd)
      COMPREPLY=( $(compgen -W "$(_cmdr_cmd_help_events $cmd)" -- ${cur}) )
      return 0
      ;;
    *)
      COMPREPLY=( $(compgen -W "$(_cmdr_cmd_help_events ${words[@]})" -- ${cur}) )
      return 0
      ;;
  esac

  #opts="--help --version -q --quiet -v --verbose --system --dest="
  #opts="--help upgrade version deploy undeploy log ls ps start stop restart"
  opts="--help"
  cmds=$($cmd --help|grep "^  [^ \[\$\#\!/\\@\"']"|awk -F'   ' '{print $1}'|awk -F',' '{for (i=1;i<=NF;i++) print $i}')

  COMPREPLY=( $(compgen -W "${opts} ${cmds}" -- ${cur}) )

} # && complete -F _cmdr_cmd_{{.AppName}} {{.AppName}}

if type complete >/dev/null 2>&1; then
	# bash
	complete -F _cmdr_cmd_{{.AppName}} {{.AppName}}
else if type compdef >/dev/null 2>&1; then
	# zsh
	_cmdr_cmd_{{.AppName}}_zsh() { compadd $(_cmdr_cmd_{{.AppName}}); }
	compdef _cmdr_cmd_{{.AppName}}_zsh {{.AppName}}
fi; fi
`)
	if err == nil {
		linuxRoot := os.Getuid() == 0

		for _, s := range []string{"/etc/bash_completion.d", "/usr/local/etc/bash_completion.d", "/tmp"} {
			if FileExists(s) {
				file := path.Join(s, cmd.root.AppName)
				var f *os.File
				if f, err = os.Create(file); err != nil {
					if !linuxRoot {
						continue
					}
					return
				}

				err = tmpl.Execute(f, cmd.root)
				if err == nil {
					fmt.Printf(`''%v generated.
Re-login to enable the new bash completion script.
`, file)
				}
				if !linuxRoot {
					break // for non-root user, we break file-writing loop and dump scripts to console too.
				}
				return

			}
		}

		err = tmpl.Execute(os.Stdout, cmd.root)
	}
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
	dir := GetStringP(prefix, "dir")
	if err = EnsureDir(dir); err != nil {
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
	fn = fmt.Sprintf("%s/%v.1", dir, fn)

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

		dir := GetStringP(prefix, "dir")
		if err = EnsureDir(dir); err != nil {
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
		fn = fmt.Sprintf("%s/%v.1", dir, fn)

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

		dir := GetStringP(prefix, "dir")
		if err = EnsureDir(dir); err != nil {
			return
		}

		fn := cmd.root.AppName
		if !cmd.IsRoot() {
			cmds := replaceAll(backtraceCmdNames(cmd, false), ".", "-")
			if len(cmds) > 0 {
				fn += "-" + cmds
			}
		}
		fn = fmt.Sprintf("%s/%v.md", dir, fn)

		w.paintFromCommand(painter, cmd, false)
		if err = ioutil.WriteFile(fn, painter.Results(), 0644); err == nil {
			log.Printf("'%v' generated...", fn)
		}
		return
	})

	return
}
