package cmdr

import (
	"fmt"
	"github.com/hedzr/log/dir"
	"io"
	"os"
	"path"
	"text/template"
)

func (w *ExecWorker) genShellBash(writer io.Writer, fullPath string, cmd *Command, args []string) (err error) {
	var tmpl *template.Template
	tmpl, err = template.New("bash.completion").Parse(`

#

# bash completion wrapper for {{.AppName}}
# version: {{.Version}}
#
# Copyright (c) 2019-2025 cmdr Authors
# See also: https://githubc.com/hedzr/cmdr
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

		if writer == nil {
			for _, s := range []string{"/etc/bash_completion.d", "/usr/local/etc/bash_completion.d", "/tmp"} {
				if dir.FileExists(s) {
					file := path.Join(s, cmd.root.AppName)
					var f *os.File
					if f, err = os.Create(file); err != nil {
						if !linuxRoot {
							continue
						}
						return
					}

					writer, fullPath = f, file
					break

					//					err = tmpl.Execute(f, cmd.root)
					//					if err == nil {
					//						fmt.Printf(`''%v generated.
					// Re-login to enable the new bash completion script.
					// `, file)
					//					}
					//					if !linuxRoot {
					//						break // for non-root user, we break file-writing loop and dump scripts to console too.
					//					}
					//					return

				}
			}
		}

		if writer == nil {
			err = tmpl.Execute(os.Stdout, cmd.root)
		} else {
			err = tmpl.Execute(writer, cmd.root)

			if !linuxRoot {
				// for non-root user, we break file-writing loop and dump scripts to console too.
				err = tmpl.Execute(os.Stdout, cmd.root)
			}

			if err == nil && fullPath != "" {
				fmt.Printf(`
# %q generated.
# Re-login to enable the new bash completion script.
`, fullPath)
			}

		}
	}
	return
}
