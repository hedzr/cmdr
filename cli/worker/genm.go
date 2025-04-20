package worker

import (
	"context"
	"os"

	"github.com/hedzr/cmdr/v2/cli"
)

// func a(filename string) (err error) {
// 	var f *os.File
// 	f, err = os.Create(filename)
// 	if err != nil {
// 		return
// 	}
// 	defer func() { err = f.Close() }()
// 	return
// }

func genManpage(ctx context.Context, filename string, w *workerS, pc *parseCtx, cmd cli.Cmd, args ...any) (err error) {
	var f *os.File
	f, err = os.Create(filename)
	if err != nil {
		return
	}
	defer func() { err = f.Close() }()

	hp := &helpPrinter{w: w, asManual: true}
	hp.PrintTo(ctx, f, pc, cmd, args...)
	return
}
