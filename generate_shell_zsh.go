package cmdr

import (
	"fmt"
	"github.com/hedzr/log/dir"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

func _makeFileIn(writer io.Writer, fullPath string, locations []string, appName string, fn func(path string, writer io.Writer) (err error)) (err error) {
	if writer != nil {
		err = fn(fullPath, writer)
		return
	}

	linuxRoot := os.Getuid() == 0
	for _, s := range locations {
		//Logger.Infof("--- checking %s", s)
		if dir.FileExists(s) {
			file := path.Join(s, "_"+appName)
			//Logger.Debugf("    try creating %s", file)
			var f *os.File
			if f, err = os.Create(file); err != nil {
				if !linuxRoot {
					err = nil
					continue // for non-root user, we break file-writing loop and dump scripts to console too.
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

func genShellZshHO(cmd *Command, args []string) func(path string, writer io.Writer) (err error) {
	return func(path string, writer io.Writer) (err error) {
		err = genZshTo(cmd, args, path, writer)
		return
	}
}

func genZshTo(cmd *Command, args []string, path string, writer io.Writer) (err error) {
	ctx := &genZshCtx{
		cmd: cmd,
		theArgs: &internalShellTemplateArgs{
			cmd.root,
			GetString("cmdr.Version"),
			cmd,
			args,
		},
		output: writer,
	}

	err = zshTplExpand(ctx, "zsh.completion.head", zshCompHead, ctx.theArgs)
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

		err = zshTplExpand(ctx, "zsh.completion.tail", zshCompTail, ctx.theArgs)
		fmt.Printf(`
# %q generated.
# Re-login to enable the new zsh completion script.
`, path)
	}
	return
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
		ctx.theArgs,
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
		ctx.theArgs,
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

func safeZshFnName(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

func zshDescribeNames(s *Command) string {
	str := s.GetTitleZshNamesBy(",")
	if strings.Contains(str, ",") {
		return "{" + str + "}"
	}
	return str
}

func zshDescribeFlagNames(s *Flag, shortTitleOnly, longTitleOnly bool) string {
	//if shortTitleOnly {
	//	return s.GetTitleZshFlagShortName()
	//}
	//if longTitleOnly {
	//	return s.GetTitleZshFlagName()
	//}
	str := s.GetTitleZshNamesExtBy(",", true, true, shortTitleOnly, longTitleOnly)
	if strings.Contains(str, ",") {
		return "{" + str + "}"
	}
	return str
}

func getZshSubFuncName(cmd *Command) (safeFuncName string) {
	safeFuncName = safeZshFnName("__" + cmd.root.AppName + "_" + strings.ReplaceAll(cmd.GetDottedNamePath(), ".", "_"))
	return
}

func zshTplExpand(ctx *genZshCtx, name, tmplString string, args interface{}) (err error) {
	var tmpl *template.Template
	tmpl, err = template.New(name).Parse(tmplString)
	if err == nil {
		err = tmpl.Execute(ctx.output, args)
	}
	return
}

type genZshCtx struct {
	cmd     *Command
	theArgs *internalShellTemplateArgs
	output  io.Writer
}

type internalShellTemplateArgs struct {
	*RootCommand
	CmdrVersion string
	Command     *Command
	Args        []string
}

func seqName(s string) string {
	s = strings.TrimLeft(s, "_")
	s = reULs.ReplaceAllString(s, " ")
	return "[" + s + "]"
}

func unquote(s string) string {
	return regexp.MustCompile(`'(.*?)'`).ReplaceAllString(s, "$1")
}

var (
	zshCompHead = `#compdef _{{.AppName}} {{.AppName}}
#
# version: {{.Version}}
#
# Copyright (c) 2019-2025 Cmdr-(go) Authors
# All rights reserved.
#
#  Zsh completion script for cmdr-series CLI apps (https://github.com/topics/cmdr).
#
#  Status: See FIXME and TODO tags.
#
#  Source: https://github.com/zsh-users/zsh-completions
#
#  Description:
#
#    Generated with '{{.AppName}} gen sh --zsh{{range .Args}} {{.}}{{end}}' for cmdr version {{.CmdrVersion}}
#
#    To install, move/rename this file as $HOME/.zsh-completions/_{{.AppName}}
#    and edit your .zshrc file to include these lines (uncommented):
# 
#    fpath=($HOME/.zsh-completions $fpath)
# 
#    autoload -U compinit
#    compinit
#
# ------------------------------------------------------------------------------
# -*- mode: zsh; sh-indentation: 2; indent-tabs-mode: nil; sh-basic-offset: 2; -*-
# vim: ft=zsh sw=2 ts=2 et
# ------------------------------------------------------------------------------


# autoload
# typeset -A opt_args
autoload -U is-at-least

# reload_zsh_autocomp
# reset_zsh_autocomp
# unfunction _{{.AppName}} && autoload -U _{{.AppName}}
# find_zsh_autocomp_script _{{.AppName}}

__{{.AppName}}_debug() {
    local altfile=""
    [[ ${ENABLE_ZSH_AUTOCOMP_DEBUG:-0} -ne 0 ]] && altfile=/tmp/1 && touch $altfile
    local file="${BASH_COMP_DEBUG_FILE:-$altfile}"
    if [[ -n ${file} ]]; then
        echo "$@" >> "${file}"
    fi
}

`
	zshCompCommands = `
{{.FuncName}}() {
    typeset -A opt_args
    local -a commands
    local context curcontext="$curcontext" line state ret=0
    local I="-h --help --version -V -#"

    _arguments -s -C : \
               "1: :->cmds" \
               "*::arg:->args" \
{{.Flags}}               && ret=0
    case "$state" in
        cmds)
            commands=(
{{.DescribeCommands}}            )
            __{{.AppName}}_debug "_describe '{{.FuncNameSeq}}': ${commands[@]}"
            _describe -t commands '{{.FuncNameSeq}} commands' commands "$@"
            ;;
        args)
            case $line[1] in
{{.CaseCommands}}
            esac
            ;;
    esac
}
`
	zshCompCommandFlags = `
{{.FuncName}}() {
{{.Declarations}}{{.DescribeCommands}}
}
`
	zshCompTail = `

# don't run the completion function when being source-ed or eval-ed
if [ "$funcstack[1]" = "_{{.AppName}}" ]; then
	_{{.AppName}}
fi

# Local Variables:
# mode: Shell-Script
# sh-indentation: 2
# indent-tabs-mode: nil
# sh-basic-offset: 2
# End:
# vim: ft=zsh sw=2 ts=2 et`
)
