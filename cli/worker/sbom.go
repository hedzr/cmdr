package worker

import (
	"debug/buildinfo"
	"fmt"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/dir"

	"gopkg.in/hedzr/errors.v3"
)

type sbomS struct{}

func (w *sbomS) onAction(cmd *cli.Command, args []string) (err error) { //nolint:revive,unused
	ec := errors.New("[processing executables]")
	if len(args) == 0 {
		args = append(args, dir.GetExecutablePath()) //nolint:revive
	}
	for _, file := range args {
		ec.Attach(w.sbomOne(file))
	}
	return
}

func (w *sbomS) sbomOne(file string) (err error) {
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
