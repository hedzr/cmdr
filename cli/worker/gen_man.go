package worker

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/is/dir"
	"gopkg.in/hedzr/errors.v3"
)

type genManS struct{}

func (w *genManS) onAction(ctx context.Context, cmd cli.Cmd, args []string) (err error) { //nolint:revive,unused
	outDir := cmd.Store().MustString("dir")
	all := cmd.Store().MustBool("all")
	fmt.Printf("# generating manpages (output-dir: %s, all: %v) ...\n", outDir, all)

	// fmt.Printf("# app.name = %s\n", UniqueWorker().Name())
	// fmt.Printf("# app.unique = %v\n", UniqueWorker())
	// app := cmd.Root().App()
	// fmt.Printf("# app = %v\n", app)
	// fmt.Printf("# app.version = %s\n", cmdr.AppVersion())

	app := cmd.Root().App()
	runner := app.GetRunner()
	if worker, ok := runner.(*workerS); !ok {
		return errors.New("invalid workerS object")
	} else {
		pc := worker.parsingCtx.(*parseCtx)

		// todo generate manpages for all commands
		dir.EnsureDir(outDir)

		var cx = cmd
		if all {
			cx = cx.Root()
		}

		appName := app.Name()
		cx.Walk(ctx, func(cc cli.Cmd, index, level int) {
			title := appName
			if !cc.IsRoot() {
				dp := cc.GetDottedPath()
				title = appName + "-" + strings.ReplaceAll(dp, ".", "-")
			}
			name := path.Join(outDir, title+".man")
			fmt.Printf("#    writing to %s...\n", name)
			genManpage(ctx, name, worker, pc, cc)
		})
	}

	fmt.Printf("#    DONE.\n")
	return
}

////////////////////////////////////////////////////////

//
//
// /////////////////////////////////////////
//
//
