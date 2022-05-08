package cmdr

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log"
	"github.com/hedzr/log/dir"
	"github.com/hedzr/log/exec"
)

type genzsh struct {
	shell     string
	locations []string
	fullPath  string
	appName   string
}

func (g *genzsh) Generate(writer io.Writer, fullPath string, cmd *Command, args []string) (err error) {
	if err = g.detectShell(cmd, args); err != nil {
		// if current shell is not zsh, generate to stdout, and return right away
		return
	}

	if fullPath == "" && len(args) > 0 {
		for _, a := range args {
			if a == "-" {
				err = g.genZshTo(cmd, args, "-", os.Stdout)
				return
			}
		}

		fullPath = args[0]
	}
	g.fullPath, g.appName = fullPath, cmd.root.AppName

	// find fpath and write to the target

	_, fpath, _ := exec.RunWithOutput(g.shell, "-c", `echo $fpath`)
	// Logger.Infof("fpath = %v", fpath)
	// Logger.Infof("ENV:\n%v", os.Environ())
	//
	// /usr/local/share/zsh/site-functions
	// $HOME/.oh-my-zsh/completions
	// $HOME/.oh-my-zsh/functions

	g.locations = tool.ReverseStringSlice(strings.Split(strings.TrimRight(fpath, "\n"), " "))
	err = g.generateFileIntoWriter(writer, g.genShellZshHO(cmd, args))
	return
}

func (g *genzsh) detectShell(cmd *Command, args []string) (err error) {
	g.shell = os.Getenv("SHELL")
	if !strings.Contains(g.shell, "/bin/zsh") {
		var zsh string
		if _, zsh, err = exec.RunWithOutput("which", "zsh"); err != nil {
			// err = errors.New("Couldn't find zsh installation, please install zsh and try again")
			if err = g.genZshTo(cmd, args, "-", os.Stdout); err == nil {
				err = ErrShouldBeStopException
			}
			return
		}

		g.shell = zsh
	}
	return
}

func (g *genzsh) generateFileIntoWriter(writer io.Writer, fn func(path string, writer io.Writer) (err error)) (err error) {
	if writer != nil {
		err = fn(g.fullPath, writer)
		return
	}

	linuxRoot := os.Getuid() == 0
	file := "-"
	var wr io.Writer = os.Stdout
	var f *os.File
	for _, s := range g.locations {
		// Logger.Infof("--- checking %s", s)
		if dir.FileExists(s) { //nolint:gocritic //like it
			file = path.Join(s, "_"+g.appName)
			// Logger.Debugf("    try creating %s", file)
			if f, err = os.Create(file); err != nil {
				if !linuxRoot {
					err = nil
					continue // for non-root user, we break file-writing loop and dump scripts to console too.
				}
				return
			}
			Logger.Debugf("    generating %s", file)
			wr = f
			// err = fn(file, f)
			// if !linuxRoot {
			//	break // for non-root user, we break file-writing loop and dump scripts to console too.
			// }
			break
		}
	}

	if f != nil {
		defer func(f *os.File) {
			err = f.Close()
		}(f)
	}

	err = fn(file, wr)
	return
}

func (g *genzsh) genShellZshHO(cmd *Command, args []string) func(path string, writer io.Writer) (err error) {
	return func(path string, writer io.Writer) (err error) {
		err = g.genZshTo(cmd, args, path, writer)
		return
	}
}

func (g *genzsh) genZshTo(cmd *Command, args []string, filepath string, writer io.Writer) (err error) {
	ctx := &genshCtx{
		cmd: cmd,
		theArgs: &internalShellTemplateArgs{
			cmd.root,
			GetString("cmdr.Version"),
			cmd,
			args,
		},
		output: writer,
	}

	err = genshTplExpand(ctx, "zsh.completion.head", zshCompHead, ctx.theArgs)
	if err == nil {
		err = walkFromCommand(&cmd.root.Command, 0, 0,
			func(cx *Command, index, level int) (err error) {
				// generate for command cx
				safeName := g.getZshSubFuncName(cx)
				if level == 0 {
					safeName = "_" + g.safeZshFnName(cx.root.AppName)
				}
				err = g.genZshFnFromTpl(ctx, safeName, cx)
				return
			})
		if err != nil {
			log.Warn(err)
		}

		err = genshTplExpand(ctx, "zsh.completion.tail", zshCompTail, ctx.theArgs)
		fmt.Printf(`
# %q generated.
# Re-login to enable the new zsh completion script.
`, filepath)
	}
	return
}

func (g *genzsh) genZshFnFromTpl(ctx *genshCtx, fnName string, cmd *Command) (err error) {
	if len(cmd.SubCommands) > 0 {
		err = g.genZshFnByCommand(ctx, fnName, cmd)
	} else {
		// if fnName == "__fluent_generate_shell" {
		//	print()
		// }
		err = g.genZshFnFlagsByCommand(ctx, fnName, cmd)
	}
	return
}

func (g *genzsh) genZshFnByCommand(ctx *genshCtx, fnName string, cmd *Command) (err error) {
	var flags, descCommands, caseCommands strings.Builder

	g.gzt1(&flags, cmd, true)

	if len(cmd.Flags) > 0 {
		flags.WriteString("               \\\n")
	}

	for _, c := range cmd.SubCommands {
		descCommands.WriteString(fmt.Sprintf("                %v':%v'\n",
			g.zshDescribeNames(c), c.GetDescZsh()))
		caseCommands.WriteString(fmt.Sprintf(`            %v)
                %v
                ;;
`, g.safeZshFnName(c.GetTitleNamesBy("|")), g.getZshSubFuncName(c)))
	}

	err = genshTplExpand(ctx, "zsh.completion.sub-commands", zshCompCommands, stdShellTemplateArgs{
		internalShellTemplateArgs: ctx.theArgs,
		FuncName:                  fnName, FuncNameSeq: seqName(fnName),
		Flags:            flags.String(),
		DescribeCommands: strings.TrimSuffix(descCommands.String(), " \\\n"),
		CaseCommands:     caseCommands.String(),
	})
	return
}

func (g *genzsh) genZshFnFlagsByCommand(ctx *genshCtx, fnName string, cmd *Command) (err error) {
	var descCommands strings.Builder

	if len(cmd.Flags) > 0 {
		descCommands.WriteString("    _arguments -s \\\n")
	}

	g.gzt1(&descCommands, cmd, false)

	desc := strings.TrimSuffix(descCommands.String(), " \\\n")
	decl := ""
	if desc != "" {
		decl = `    typeset -A opt_args
    local context curcontext="$curcontext" state line ret=0
    local I="-h --help --version -V -#"
    local -a commands
`
	}

	err = genshTplExpand(ctx, "zsh.completion.flags", zshCompCommandFlags, stdShellTemplateArgs{
		internalShellTemplateArgs: ctx.theArgs,
		FuncName:                  fnName,
		Declarations:              decl,
		DescribeCommands:          desc,
	})
	return
}

func (g *genzsh) gzt1(descCommands *strings.Builder, cmd *Command, shortTitleOnly bool) {
	g.gzt1ForToggleGroups(descCommands, cmd, shortTitleOnly)

	for ix, c := range cmd.Flags {
		if c.ToggleGroup != "" {
			continue
		}
		g.gzt2(descCommands, cmd, ix, c, "", shortTitleOnly)
	}
}

func (g *genzsh) gzt1ForToggleGroups(descCommands *strings.Builder, cmd *Command, shortTitleOnly bool) {
	tgs := make(map[string][]*Flag)
	for _, f := range cmd.Flags {
		if f.ToggleGroup != "" {
			tgs[f.ToggleGroup] = append(tgs[f.ToggleGroup], f)
		}
	}

	for k, v := range tgs {
		me := g.gzChkMEForToggleGroup(k, v)

		for ix, c := range v {
			// var sb strings.Builder
			// for i, f := range v {
			//	if i != ix {
			//		sb.WriteString(f.GetTitleZshNamesBy(" "))
			//		sb.WriteString(" ")
			//	}
			// }
			g.gzt2(descCommands, cmd, ix, c, me, shortTitleOnly)
		}
	}
}

func (g *genzsh) gzt2(descCommands *strings.Builder, cmd *Command, ix int, f *Flag, mutualExclusives string, shortTitleOnly bool) {
	// if c.Full == "pprof" {
	//	println()
	// }

	switch {
	case len(f.ValidArgs) != 0:
		g.gzAction(descCommands, f, "("+strings.Join(f.ValidArgs, " ")+")", mutualExclusives, shortTitleOnly)
	case f.DefaultValuePlaceholder == "FILE":
		act := "_files"
		if f.actionStr != "" {
			act += " -g " + strconv.Quote(f.actionStr)
		}
		g.gzAction(descCommands, f, act, mutualExclusives, shortTitleOnly)
	case f.DefaultValuePlaceholder == "DIR":
		g.gzAction(descCommands, f, "_files -/", mutualExclusives, shortTitleOnly)
	case f.DefaultValuePlaceholder == "USER":
		g.gzAction(descCommands, f, "_users", mutualExclusives, shortTitleOnly)
	case f.DefaultValuePlaceholder == "GROUP":
		g.gzAction(descCommands, f, "_groups", mutualExclusives, shortTitleOnly)
	case f.DefaultValuePlaceholder == "INTERFACES":
		g.gzAction(descCommands, f, "_net_interfaces", mutualExclusives, shortTitleOnly)
	default:
		mutualExclusives = g.gzChkME(f, mutualExclusives)
		if mutualExclusives != "" {
			descCommands.WriteString(fmt.Sprintf("                \"($I %v)\"%v'[%v]'",
				mutualExclusives, g.zshDescribeFlagNames(f, shortTitleOnly, false), f.GetDescZsh()))
		} else {
			descCommands.WriteString(fmt.Sprintf("                %v'[%v]'",
				g.zshDescribeFlagNames(f, shortTitleOnly, false), f.GetDescZsh()))
		}
	}
	// if ix != len(cmd.Flags)-1 {
	descCommands.WriteString(" \\\n")
	// }
}

// func (g *genzsh) gzAction(descCommands *strings.Builder, c *Flag, action, mutualExclusives string) {
//	g.gzAction_(descCommands,c, action, mutualExclusives, false)
// }

func (g *genzsh) gzAction(descCommands *strings.Builder, f *Flag, action, mutualExclusives string, shortTitleOnly bool) {
	if f.dblTildeOnly {
		return
	}

	names := g.zshDescribeFlagNames(f, shortTitleOnly, false)
	title := f.Full
	if f.DefaultValuePlaceholder != "" {
		title = f.DefaultValuePlaceholder
	}
	mutualExclusives = g.gzChkME(f, mutualExclusives)
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
func (g *genzsh) gzChkME(f *Flag, mutualExclusives string) string {
	const quoted = false
	if mutualExclusives == "" {
		switch {
		case len(f.mutualExclusives) > 0:
			var sb strings.Builder
			for _, t := range f.mutualExclusives {
				if tgt, ok := f.owner.plainLongFlags[t]; ok {
					sb.WriteString(tgt.GetTitleZshNamesExtBy(" ", false, quoted, false, false))
				}
			}
			sb.WriteString(f.GetTitleZshNamesExtBy(" ", false, quoted, false, false))
			mutualExclusives = strings.TrimRight(sb.String(), " ")
		case f.circuitBreak:
			mutualExclusives = "- *"
		default:
			mutualExclusives = f.GetTitleZshNamesExtBy(" ", false, quoted, false, false)
		}
	}
	return mutualExclusives
}

func (g *genzsh) gzChkMEForToggleGroup(toggleGroupName string, v []*Flag) (mutualExclusives string) {
	const quoted = false
	var sb strings.Builder
	for _, f := range v {
		sb.WriteString(f.GetTitleZshNamesExtBy(" ", false, quoted, false, false))
		sb.WriteString(" ")
	}
	mutualExclusives = strings.TrimRight(sb.String(), " ")
	return
}

func (g *genzsh) safeZshFnName(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

func (g *genzsh) zshDescribeNames(s *Command) string {
	str := s.GetTitleZshNamesBy(",")
	if strings.Contains(str, ",") {
		return "{" + str + "}"
	}
	return str
}

func (g *genzsh) zshDescribeFlagNames(s *Flag, shortTitleOnly, longTitleOnly bool) string {
	// if shortTitleOnly {
	//	return s.GetTitleZshFlagShortName()
	// }
	// if longTitleOnly {
	//	return s.GetTitleZshFlagName()
	// }
	str := s.GetTitleZshNamesExtBy(",", true, true, shortTitleOnly, longTitleOnly)
	if strings.Contains(str, ",") {
		return "{" + str + "}"
	}
	return str
}

func (g *genzsh) getZshSubFuncName(cmd *Command) (safeFuncName string) {
	safeFuncName = g.safeZshFnName("__" + cmd.root.AppName + "_" + strings.ReplaceAll(cmd.GetDottedNamePath(), ".", "_"))
	return
}

// func (g *genzsh) zshTplExpand(ctx *genZshCtx, name, tmplString string, args interface{}) (err error) {
//	var tmpl *template.Template
//	tmpl, err = template.New(name).Parse(tmplString)
//	if err == nil {
//		err = tmpl.Execute(ctx.output, args)
//	}
//	return
// }

func genshTplExpand(ctx *genshCtx, tmplName, tmplString string, data interface{}) (err error) {
	var tmpl *template.Template
	tmpl, err = template.New(tmplName).Parse(tmplString)
	if err == nil {
		err = tmpl.Execute(ctx.output, data)
	}
	return
}

type genshCtx struct {
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

type stdShellTemplateArgs struct {
	*internalShellTemplateArgs
	FuncName         string
	FuncNameSeq      string // "__fluent_ms" => "[fluent ms]"
	Flags            string
	Declarations     string
	DescribeCommands string
	CaseCommands     string
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
