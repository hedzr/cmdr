// Copyright © 2022 Hedzr Yeh.

//go:build go1.18
// +build go1.18

package cmdr

import (
	"debug/buildinfo"
	"fmt"

	"github.com/hedzr/log"
	"github.com/hedzr/log/exec"
	"gopkg.in/hedzr/errors.v3"
)

func sbomAttach(w *ExecWorker, root *RootCommand) {
	found := false
	for _, sc := range root.SubCommands {
		if sc.Full == "sbom" { // generatorCommands.Full {
			found = true
			break
		}
	}
	if !found {
		// root.SubCommands = append(root.SubCommands, generatorCommands)
		var sbom *Command
		w._cmdAdd(root, "sbom", "Print SBOM Info (Software Bill Of Materials).", func(cx1 *Command) {
			// cx1.Short = "s"
			// cx1.Aliases = []string{}
			cx1.LongDescription = `Print SBOM information of this or specified executable(s).

				The outputs is YAML compliant.
			
				Just another way to run 'go version -m executable-file' but no need to install Go Runtime.`
			cx1.Examples = ``
			cx1.Action = sbomAction
			sbom = cx1
		})
		w._boolFlgAdd1(sbom, "more", "Dump more information.", SysMgmtGroup, func(ff *Flag) {
			ff.Short, ff.EnvVars, ff.VendorHidden = "m", []string{"MORE"}, false
		})
	}
}

func sbomAction(cmd *Command, args []string) (err error) {
	var ec = errors.New("processing executables")
	var caught bool
	for _, file := range args {
		ec.Attach(sbomOne(file))
		caught = true
	}
	if !caught {
		file := exec.GetExecutablePath()
		log.Infof("SBOM on %v", file)
		ec.Attach(sbomOne(file))
	}
	return
	// var ec = errors.New("processing executables")
	// if len(args) == 0 {
	// 	args = append(args, dir.GetExecutablePath())
	// }
	// for _, file := range args {
	// 	ec.Attach(sbomOne(file))
	// }
	// return
}

func sbomOne(file string) (err error) {
	var inf *buildinfo.BuildInfo
	if inf, err = buildinfo.ReadFile(file); err != nil {
		return
	}

	fmt.Printf(`SBOM:
  executable: %q
  go-version: %v
  path: %v
  module-path: %v
  module-version: %v
  module-sum: %v
  module-replace: <ignored>
  settings:
`,
		file, inf.GoVersion, inf.Path,
		inf.Main.Path, inf.Main.Version, inf.Main.Sum,
	)

	for _, d := range inf.Settings {
		fmt.Printf("    - %q: %v\n", d.Key, d.Value)
	}
	fmt.Println("  depends:")
	for _, d := range inf.Deps {
		// str := fmt.Sprintf("%#v", *d)
		fmt.Printf("    - debug-module: { path: %q, version: %q, sum: %q, replace: %#v } \n", d.Path, d.Version, d.Sum, d.Replace)
	}
	return
}
