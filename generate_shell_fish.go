package cmdr

import (
	"fmt"
	"io"
)

func (w *ExecWorker) genShellFish(writer io.Writer, fullPath string, cmd *Command, args []string) (err error) {
	fmt.Println(`# todo fish`)
	return
}

const (
	fishHeader = `# Fish Shell Completions
# Place or symlink to $XDG_CONFIG_HOME/fish/completions/bat.fish ($XDG_CONFIG_HOME is usually set to ~/.config)
#
# For: {{.AppName}}
# Version: {{.Version}}
#
# Generated with '{{.AppName}} gen sh --fish{{range .Args}} {{.}}{{end}}' for cmdr version {{.CmdrVersion}}
#
# Copyright (c) 2019-2025 Cmdr-(go) Authors
# All rights reserved.
#

`
)
