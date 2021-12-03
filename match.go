// Copyright © 2020 Hedzr Yeh.

package cmdr

import (
	"os"
	"strings"
	"text/template"
)

// Match try parsing the input command-line, the result is the last hit *Command.
func Match(inputCommandlineWithoutArg0 string, opts ...ExecOption) (last *Command, err error) {
	saved := internalGetWorker()
	savedUnknownOptionHandler := unknownOptionHandler
	defer func() {
		uniqueWorkerLock.Lock()
		uniqueWorker = saved
		unknownOptionHandler = savedUnknownOptionHandler
		uniqueWorkerLock.Unlock()
	}()

	rootCmd := internalGetWorker().rootCommand

	w := internalResetWorkerNoLock()

	for _, opt := range opts {
		opt(w)
	}

	w.noDefaultHelpScreen = true
	w.noUnknownCmdTip = true
	w.noCommandAction = true
	unknownOptionHandler = emptyUnknownOptionHandler

	line := os.Args[0] + " " + inputCommandlineWithoutArg0
	last, err = w.InternalExecFor(rootCmd, strings.Split(line, " "))
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
	str := s.GetTitleZshNamesExtBy(",", true, shortTitleOnly, longTitleOnly)
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
	cmd    *Command
	args   *internalShellTemplateArgs
	output *os.File
}

type internalShellTemplateArgs struct {
	*RootCommand
	CmdrVersion string
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
#    Generated with '{{.AppName}} gen zsh' for cmdr version {{.CmdrVersion}}
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
    local line state ret

    _arguments -s -C : \
               "1: :->cmds" \
               "*::arg:->args" \
{{.Flags}}               && ret=0
    case "$state" in
        cmds)
            local commands; commands=(
{{.DescribeCommands}}            )
            __{{.AppName}}_debug "_describe '{{.FuncName}}': ${commands[@]}"
            _describe -t commands '{{.FuncName}} commands' commands "$@"
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
	_fluent
fi

# Local Variables:
# mode: Shell-Script
# sh-indentation: 2
# indent-tabs-mode: nil
# sh-basic-offset: 2
# End:
# vim: ft=zsh sw=2 ts=2 et
`
)
