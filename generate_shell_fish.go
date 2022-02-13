package cmdr

import (
	"fmt"
	"github.com/hedzr/log/dir"
	"gopkg.in/hedzr/errors.v2"
	"io"
	"os"
	"path"
	"strings"
)

func (w *ExecWorker) genShellFish(writer io.Writer, fullPath string, cmd *Command, args []string) (err error) {
	//fmt.Println(`# todo fish`)
	var gen genfish
	gen.init()
	err = gen.Generate(writer, fullPath, cmd, args)
	return
}

type genfish struct {
	shell    bool
	homeDir  string
	fishDir  string
	fullPath string
	appName  string
}

func (g *genfish) init() {
	g.detectFishShell()
	g.detectFishFolders()
}

//func (g *genfish) isFishShell() bool { return g.shell }
func (g *genfish) detectFishShell() { g.shell = isFishShell() }

func isFishShell() bool { return os.Getenv("FISH_VERSION") == os.Getenv("version") }

func (g *genfish) detectFishFolders() {
	g.homeDir = os.Getenv("HOME") // note that it's available in cmdr system specially for windows since we ever duplicated USERPROFILE as HOME.
	fishDir := path.Join(g.homeDir, ".config", "fish")
	if dir.FileExists(fishDir) {
		g.fishDir = fishDir
	}
}

func (g *genfish) Generate(writer io.Writer, fullPath string, cmd *Command, args []string) (err error) {
	if fullPath == "" && len(args) > 0 {
		for _, a := range args {
			if a == "-" {
				err = g.genTo(os.Stdout, cmd, args)
				return
			}
		}

		fullPath = args[0]
	}
	g.fullPath, g.appName = fullPath, cmd.root.AppName

	if g.fishDir != "" && g.fullPath == "" && writer == nil {
		fullPath = path.Join(g.fishDir, "completions", g.appName+".fish")
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

func (g *genfish) genTo(writer io.Writer, cmd *Command, args []string) (err error) {

	var ctx = &genshCtx{
		cmd: cmd,
		theArgs: &internalShellTemplateArgs{
			cmd.root,
			GetString("cmdr.Version"),
			cmd,
			args,
		},
		output: writer,
	}

	err = genshTplExpand(ctx, "fish.completion.head", fishCompHead, ctx.theArgs)
	if err == nil {

		err = genshTplExpand(ctx, "fish.completion.body", fishCompBody, ctx.theArgs)

		err = genshTplExpand(ctx, "fish.completion.tail", fishCompTail, ctx.theArgs)

		if g.fullPath != "-" {
			fmt.Printf(`
# %q generated.
# It takes affect in a while, try '%v<TAB>' right now.
`, g.fullPath, g.appName)
		}
	}
	return
}

// helpSystemAction is __complete command handler.
func (w *ExecWorker) helpSystemAction(cmdComplete *Command, args []string) (err error) {
	var (
		ctx = &queryShcompContext{
			w:           w,
			cmdComplete: cmdComplete,
			args:        args,

			matchedPrecededList:     make(map[string]*Command),
			matchedList:             make(map[string]*Command),
			matchedPrecededFlagList: make(map[string]*Flag),
			matchedFlagList:         make(map[string]*Flag),
			directive:               shellCompDirectiveNoFileComp,
		}
		x     strings.Builder
		keys  []string
		count int
	)

	defer func() {
		x.WriteString(fmt.Sprintf(":%d", ctx.directive))
		fp("%v", x.String())
		_, _ = fmt.Fprintf(os.Stderr, "%v\n%v Items populated.\n", directivesToString(ctx.directive), count)
	}()

	err = ctx.lookupForHelpSystem(cmdComplete, args)
	if err == nil || err == ErrShouldBeStopException {

		var cptLocal = getCPT()

		if ctx.matchedFlag != nil {
			keys = getSortedKeysFromFlgMap(ctx.matchedPrecededFlagList)
			for _, k := range keys {
				c := ctx.matchedPrecededFlagList[k]
				x.WriteString(fmt.Sprintf("%v\t%v\n", k, cptLocal.Translate(c.Description, BgNormal)))
				count++
			}

			keys = getSortedKeysFromFlgMap(ctx.matchedFlagList)
			for _, k := range keys {
				c := ctx.matchedFlagList[k]
				x.WriteString(fmt.Sprintf("%v\t%v\n", k, cptLocal.Translate(c.Description, BgNormal)))
				count++
			}
			return
		}

		keys = getSortedKeysFromCmdMap(ctx.matchedPrecededList)
		for _, k := range keys {
			c := ctx.matchedPrecededList[k]
			x.WriteString(fmt.Sprintf("%v\t%v\n", k, cptLocal.Translate(c.Description, BgNormal)))
			count++
		}

		keys = getSortedKeysFromCmdMap(ctx.matchedList)
		for _, k := range keys {
			c := ctx.matchedList[k]
			x.WriteString(fmt.Sprintf("%v\t%v\n", k, cptLocal.Translate(c.Description, BgNormal)))
			count++
		}

	}
	return
}

type queryShcompContext struct {
	w           *ExecWorker
	cmdComplete *Command
	args        []string

	// returned

	matchedCmd  *Command
	matchedFlag *Flag
	exact       bool

	// dynamic updating

	matchedPrecededList, matchedList         map[string]*Command
	matchedPrecededFlagList, matchedFlagList map[string]*Flag
	directive                                int
}

func (ctx *queryShcompContext) buildMatchedSubCommandList(cmd *Command) {
	ctx.matchedCmd = cmd
	for _, c := range cmd.SubCommands {
		ctx.matchedList[c.GetTitleName()] = c
	}
}

func (ctx *queryShcompContext) lookupForHelpSystem(cmdComplete *Command, args []string) (err error) {
	var list []*Command
	cmd := &cmdComplete.root.Command
	var flg *Flag
	var exact bool
	defer func() {
		ctx.matchedCmd, ctx.matchedFlag, ctx.exact = cmd, flg, exact
	}()

	for i, ic := 0, len(args); i < ic; i++ {
		ttl := args[i]
		l := len(ttl)

		if l == 0 {
			if i == 0 {
				// root command matched
				ctx.buildMatchedSubCommandList(cmd)
				return
			} else if i == ic-1 {
				// for commandline: app __complete generate ''
				ll := len(list)
				if ll >= 1 {
					for i := ll - 1; i >= 0; i-- {
						ctx.deleteFromMatchedList(list[i])
					}
				}
				ctx.rebuildMatchedList(list[ll-1])
			}
			continue
		}

		if l > 0 && strings.Contains(ctx.w.switchCharset, ttl[0:1]) {
			// flags
			flg, err = ctx.lookupFlagsForHelpSystem(ttl, cmd, args, i)
			if err != nil && err != ErrShouldBeStopException {
				break
			}
			continue
		}

		// sub-commands
		cmd, exact, err = ctx.lookupCommandsForHelpSystem(ttl, cmd, args, i, ic)
		if err != nil && err != ErrShouldBeStopException {
			break
		}

		list = append(list, cmd)
	}
	return
}

func (ctx *queryShcompContext) lookupCommandsForHelpSystem(title string, parent *Command, args []string, ix, ic int) (cmdMatched *Command, exact bool, err error) {
	ctx.directive = shellCompDirectiveNoFileComp
	err = errors.NotFound
	ok := false

	for _, c := range parent.SubCommands {
		if c.VendorHidden {
			continue
		}

		if exact, ok = ctx.matchCommandTitle(c, title, ix == ic-1); ok {
			if /*ix == len(args)-1-1 && args[ix+1] == "" &&*/ exact {
				cmdMatched, err = c, nil
			} else {
				err = nil
			}
		}
	}
	return
}

func (ctx *queryShcompContext) rebuildMatchedList(cmd *Command) {
	ctx.deleteFromMatchedList(cmd)
	for _, c := range cmd.SubCommands {
		ctx.matchedList[c.GetTitleName()] = c
	}
}

func (ctx *queryShcompContext) deleteFromMatchedList(cmd *Command) {
	if _, ok := ctx.matchedList[cmd.GetTitleName()]; ok {
		delete(ctx.matchedList, cmd.GetTitleName())
	}
}

func (ctx *queryShcompContext) matchCommandTitle(c *Command, titleChecking string, fuzzy bool) (exact, ok bool) {
	ok = ctx.visitCommandTitles(c, false, func(c *Command, title string) (stopNow bool) {
		if exact, ok = ctx.doMatchCommandTitle(c, title, titleChecking, fuzzy); ok {
			stopNow = true
		}
		return
	})
	return
}

func (ctx *queryShcompContext) visitCommandTitles(c *Command, justFullTitle bool, fn func(c *Command, title string) (stopNow bool)) (ok bool) {
	if c.Full != "" && fn(c, c.Full) {
		return true
	}
	if justFullTitle {
		return
	}
	if c.Short != "" && fn(c, c.Short) {
		return true
	}
	for _, t := range c.Aliases {
		if t != "" && fn(c, t) {
			return true
		}
	}
	return
}

func (ctx *queryShcompContext) doMatchCommandTitle(c *Command, title, titleChecking string, fuzzy bool) (exact, ok bool) {
	if title == titleChecking {
		ctx.matchedList[c.GetTitleName()] = c
		exact, ok = true, true
		return
	}
	if fuzzy && strings.HasPrefix(title, titleChecking) {
		ctx.matchedPrecededList[c.GetTitleName()] = c
		ok = true
	} else if noPartialMatching == false && fuzzy && strings.Contains(title, titleChecking) {
		ctx.matchedPrecededList[c.GetTitleName()] = c
		ok = true
	}
	return
}

func (ctx *queryShcompContext) lookupFlagsForHelpSystem(titleChecking string, parent *Command, args []string, ix int) (flgMatched *Flag, err error) {
	ctx.directive = shellCompDirectiveNoFileComp
	err = errors.NotFound

	sw1 := len(titleChecking) > 0 && strings.ContainsAny(titleChecking[0:1], ctx.w.switchCharset)
	sw2 := len(titleChecking) > 1 && strings.ContainsAny(titleChecking[1:2], ctx.w.switchCharset)

	for _, c := range parent.Flags {
		if c.VendorHidden {
			continue
		}

		if _, ok := ctx.matchFlagTitle(c, titleChecking, sw1, sw2); ok {
			//if /*ix == len(args)-1-1 && args[ix+1] == "" &&*/ exact {
			//	flgMatched, err = c, nil
			//} else {
			//	err = nil
			//}
			flgMatched, err = c, nil
		}
	}
	return
}

func (ctx *queryShcompContext) matchFlagTitle(c *Flag, titleChecking string, sw1, sw2 bool) (exact, ok bool) {
	if len(titleChecking) == 1 && sw1 {
		ctx.matchedPrecededFlagList["--"+c.GetTitleName()] = c
		ok = true
		return
	}

	if len(titleChecking) == 2 && sw1 && sw2 {
		ctx.matchedPrecededFlagList["--"+c.GetTitleName()] = c
		ok = true
		return
	}

	if sw1 && sw2 {
		exact, ok = ctx.matchLongFlagTitle(c, titleChecking)
	}

	if !ok && c.Short != "" && sw1 {
		exact, ok = ctx.matchShortFlagTitle(c, titleChecking)
	}

	return
}

func (ctx *queryShcompContext) matchShortFlagTitle(c *Flag, titleChecking string) (exact, ok bool) {
	t := titleChecking[1:]
	if c.Short == t {
		ctx.matchedFlagList["-"+c.Short] = c
		exact, ok = true, true
	} else if strings.HasPrefix(c.Short, t) {
		ctx.matchedPrecededFlagList["-"+c.Short] = c
		ok = true
	}
	return
}

func (ctx *queryShcompContext) matchLongFlagTitle(c *Flag, titleChecking string) (exact, ok bool) {
	t := titleChecking[2:]
	if c.Full == t {
		ctx.matchedFlagList["--"+c.GetTitleName()] = c
		exact, ok = true, true
	} else if strings.HasPrefix(c.Full, t) {
		ctx.matchedPrecededFlagList["--"+c.GetTitleName()] = c
		ok = true
	} else {
		for _, st := range c.Aliases {
			if st == t {
				ctx.matchedFlagList["--"+c.GetTitleName()] = c
				exact, ok = true, true
				break
			} else if strings.HasPrefix(st, t) {
				ctx.matchedPrecededFlagList["--"+c.GetTitleName()] = c
				ok = true
			}
		}
	}
	return
}

var noPartialMatching = true

const (
	shellCompDirectiveError         = 1
	shellCompDirectiveNoSpace       = 2
	shellCompDirectiveNoFileComp    = 4
	shellCompDirectiveFilterFileExt = 8
	shellCompDirectiveFilterDirs    = 16
)

func directivesToString(d int) string {
	//var sb strings.Builder
	if d&shellCompDirectiveError != 0 {
		return "ShellCompDirectiveError"
	}
	if d&shellCompDirectiveNoSpace != 0 {
		return "ShellCompDirectiveNoSpace"
	}
	if d&shellCompDirectiveNoFileComp != 0 {
		return "ShellCompDirectiveNoFileComp"
	}
	if d&shellCompDirectiveFilterFileExt != 0 {
		return "ShellCompDirectiveFilterFileExt"
	}
	if d&shellCompDirectiveFilterDirs != 0 {
		return "ShellCompDirectiveFilterDirs"
	}
	return ""
}

const (
	fishCompHead = `# Fish Shell Completions for {{.AppName}} v{{.Version}}             -*- shell-script -*-
# Place or symlink to $XDG_CONFIG_HOME/fish/completions/{{.AppName}}.fish ($XDG_CONFIG_HOME is usually set to ~/.config)
#
# Generated with '{{.AppName}} gen sh --fish{{range .Args}} {{.}}{{end}}' for cmdr version {{.CmdrVersion}}
#
# Copyright (c) 2019-2025 Cmdr-(go) Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# 

`

	fishCompBody = `

function __{{.AppName}}_debug
    set -q FISH_COMP_DEBUG_FILE; and set -l file "$FISH_COMP_DEBUG_FILE"; or set -q BASH_COMP_DEBUG_FILE; and set -l file "$BASH_COMP_DEBUG_FILE"
    if test "$file" != ""
      if test ! -n "$file"
        touch "$file"
      end
      echo "$argv" >> $file
    end
end

function __fish_{{.AppName}}_needs_command
    # Figure out if the current invocation already has a command.
    set -l cmd (commandline -opc)
    #set -e cmd[1]
    set -l cnt (count $cmd)

  if test (count $cmd) -eq 1 -a "$cmd[1]" = '{{.AppName}}'
    return 0
  end
  return 1

    # #argparse -s (__fish_{{.AppName}}_global_optspecs) -- $cmd 2>/dev/null
    # #or return 0
    # # These flags function as commands, effectively.
    # #set -q _flag_version; and return 1
    # #set -q _flag_html_path; and return 1
    # #set -q _flag_man_path; and return 1
    # #set -q _flag_info_path; and return 1
    # if set -q argv[1]
    #     # Also print the command, so this can be used to figure out what it is.
    #     echo $argv[1]
    #     return 1
    # end
    # return 0
end

function __fish_{{.AppName}}_using_command
  set -l cmd (commandline -opc)
  __{{.AppName}}_debug ""
  __{{.AppName}}_debug "============================== BEGIN: $cmd"
  if test (count $cmd) -gt 1
    set -l ssa $argv[1..-1]
    set -l cca $cmd[2..-1]
    __{{.AppName}}_debug "  ssa: $ssa"
    __{{.AppName}}_debug "  cca: $cca"
    if test (count $ssa) -eq (count $cca)
      for ix in (seq (count $ssa))
        if test ! $ssa[$ix] = $cca[$ix]
            return 1
        end
      end
      __{{.AppName}}_debug "  yes: '$commandline'"
      commandline -i '-'
      return 0
    end
  end
  return 1
end

function __{{.AppName}}_perform_completion
    __{{.AppName}}_debug "    ===== Starting __{{.AppName}}_perform_completion ====="

    # Extract all args except the last one
    set -l args (commandline -opc)
    # Extract the last arg and escape it in case it is a space
    set -l lastArg (string escape -- (commandline -ct))

    __{{.AppName}}_debug "args: $args"
    __{{.AppName}}_debug "last arg: $lastArg"

    set -l requestComp "$args[1] __complete $args[2..-1] $lastArg"

    __{{.AppName}}_debug "Calling $requestComp"
    set -l results (eval $requestComp 2> /dev/null)

    # Some programs may output extra empty lines after the directive.
    # Let's ignore them or else it will break completion.
    # Ref: https://github.com/spf13/cobra/issues/1279
    for line in $results[-1..1]
        if test (string trim -- $line) = ""
            # Found an empty line, remove it
            set results $results[1..-2]
        else
            # Found non-empty line, we have our proper output
            break
        end
    end

    set -l comps $results[1..-2]
    set -l directiveLine $results[-1]

    # For Fish, when completing a flag with an = (e.g., <program> -n=<TAB>)
    # completions must be prefixed with the flag
    set -l flagPrefix (string match -r -- '-.*=' "$lastArg")

    __{{.AppName}}_debug "Comps: $comps"
    __{{.AppName}}_debug "DirectiveLine: $directiveLine"
    __{{.AppName}}_debug "flagPrefix: $flagPrefix"

    for comp in $comps
        printf "%s%s\n" "$flagPrefix" "$comp"
    end

    printf "%s\n" "$directiveLine"
end


function __{{.AppName}}_prepare_completions
    __{{.AppName}}_debug ""
    __{{.AppName}}_debug "============= starting completion logic =============="

    set -e __{{.AppName}}_comp_results

    set -l results (__{{.AppName}}_perform_completion)
    __{{.AppName}}_debug "Completion results: $results"

    if test -z "$results"
        __{{.AppName}}_debug "No completion, probably due to a failure"
        return 1
    end

    set -l directive (string sub --start 2 $results[-1])
    set --global __{{.AppName}}_comp_results $results[1..-2]

    __{{.AppName}}_debug "Completions are: $__{{.AppName}}_comp_results"
    __{{.AppName}}_debug "Directive is: $directive"

    set -l shellCompDirectiveError 1
    set -l shellCompDirectiveNoSpace 2
    set -l shellCompDirectiveNoFileComp 4
    set -l shellCompDirectiveFilterFileExt 8
    set -l shellCompDirectiveFilterDirs 16

    if test -z "$directive"
        set directive 0
    end

    set -l compErr (math (math --scale 0 $directive / $shellCompDirectiveError) % 2)
    if test $compErr -eq 1
        __{{.AppName}}_debug "Received error directive: aborting."
        return 1
    end

    set -l filefilter (math (math --scale 0 $directive / $shellCompDirectiveFilterFileExt) % 2)
    set -l dirfilter (math (math --scale 0 $directive / $shellCompDirectiveFilterDirs) % 2)
    if test $filefilter -eq 1; or test $dirfilter -eq 1
        __{{.AppName}}_debug "File extension filtering or directory filtering not supported"
        return 1
    end

    set -l nospace (math (math --scale 0 $directive / $shellCompDirectiveNoSpace) % 2)
    set -l nofiles (math (math --scale 0 $directive / $shellCompDirectiveNoFileComp) % 2)

    __{{.AppName}}_debug "nospace: $nospace, nofiles: $nofiles"

    # If we want to prevent a space, or if file completion is NOT disabled,
    # we need to count the number of valid completions.
    # To do so, we will filter on prefix as the completions we have received
    # may not already be filtered so as to allow fish to match on different
    # criteria than the prefix.
    if test $nospace -ne 0; or test $nofiles -eq 0
        set -l prefix (commandline -t | string escape --style=regex)
        __{{.AppName}}_debug "prefix: $prefix"

        set -l completions (string match -r -- "^$prefix.*" $__{{.AppName}}_comp_results)
        set --global __{{.AppName}}_comp_results $completions
        __{{.AppName}}_debug "Filtered completions are: $__{{.AppName}}_comp_results"

        # Important not to quote the variable for count to work
        set -l numComps (count $__{{.AppName}}_comp_results)
        __{{.AppName}}_debug "numComps: $numComps"

        if test $numComps -eq 1; and test $nospace -ne 0
            # We must first split on \t to get rid of the descriptions to be
            # able to check what the actual completion will be.
            # We don't need descriptions anyway since there is only a single
            # real completion which the shell will expand immediately.
            set -l split (string split --max 1 \t $__{{.AppName}}_comp_results[1])

            # Fish won't add a space if the completion ends with any
            # of the following characters: @=/:.,
            set -l lastChar (string sub -s -1 -- $split)
            if not string match -r -q "[@=/:.,]" -- "$lastChar"
                # In other cases, to support the "nospace" directive we trick the shell
                # by outputting an extra, longer completion.
                __{{.AppName}}_debug "Adding second completion to perform nospace directive"
                set --global __{{.AppName}}_comp_results $split[1] $split[1].
                __{{.AppName}}_debug "Completions are now: $__{{.AppName}}_comp_results"
            end
        end

        if test $numComps -eq 0; and test $nofiles -eq 0
            # To be consistent with bash and zsh, we only trigger file
            # completion when there are no other completions
            __{{.AppName}}_debug "Requesting file completion"
            return 1
        end
    end

    return 0
end

`
	fishCompTail = `



if type -q "{{.AppName}}"
    complete --do-complete "{{.AppName}} " > /dev/null 2>&1
end

# Remove any pre-existing completions for the program since we will be handling all of them.
complete -c {{.AppName}} -e

# complete -f -c {{.AppName}}

complete -c {{.AppName}} -n '__{{.AppName}}_prepare_completions' -f -a '$__{{.AppName}}_comp_results'

# Local Variables:
# mode: Shell-Script
# sh-indentation: 2
# indent-tabs-mode: nil
# sh-basic-offset: 2
# End:
# vim: ft=zsh sw=2 ts=2 et`
)
