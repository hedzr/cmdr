package worker

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/hedzr/cmdr/v2/cli"
)

func (w *genShS) genShellFish(ctx context.Context, writer io.Writer, fullPath string, cmd cli.Cmd, args []string) (err error) {
	gen := gensh{
		ext: "fish",
		tplm: map[whatTpl]string{
			wtHeader: fishCompHead,
			wtBody:   fishCompBody,
			wtTail:   fishCompTail,
		},
		getTargetPath: func(g *gensh) string {
			fullPath = path.Join(g.shConfigDir, "completions", g.appName+"."+g.ext)
			return fullPath
		},
		detectShell: func(g *gensh) { g.shell = isFishShell() },
		endingText:  "It takes affect in a while, try '%v<TAB>' right now.",
	}

	gen.init()
	err = gen.Generate(ctx, writer, fullPath, cmd, args)
	return
}

func isFishShell() bool { return os.Getenv("FISH_VERSION") == os.Getenv("version") }

const (
	fishCompHead = `# Fish Shell Completions for {{.AppName}} {{.Version}}             -*- shell-script -*-
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
