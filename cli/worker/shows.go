package worker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/conf"
	"github.com/hedzr/cmdr/v2/internal/hs"
	"github.com/hedzr/cmdr/v2/pkg/times"
)

func (w *workerS) showVersion(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd, args ...any) (err error) { //nolint:revive,unused
	// (&helpPrinter{w: w}).Print(ctx, pc, lastCmd)
	if w.OnShowVersion != nil {
		err = w.OnShowVersion(ctx, lastCmd, pc.PositionalArgs())
		return
	}

	if lastCmd.HitTitle() == "V" {
		fp(`v%v`, strings.TrimLeft(conf.Version, "v"))
		return
	}

	ts := conf.Buildstamp
	if ts == "" {
		ts = time.Now().UTC().Format(time.RFC3339)
	}
	dt, err := times.SmartParseTime(ts)
	// dt, err := time.Parse("", ts)
	if err == nil {
		ts = dt.Format(time.RFC3339)
	}

	fp(`v%v
%v
%v
%v
%v
%v
%v
%v
%v`,
		strings.TrimLeft(conf.Version, "v"),
		conf.AppName,
		ts,
		conf.Githash,
		conf.GoVersion,
		conf.GitSummary,
		conf.Serial,
		conf.SerialString,
		conf.BuilderComments,
	)
	return
}

func (w *workerS) showBuiltInfo(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd, args ...any) (err error) { //nolint:revive,unused
	// (&helpPrinter{w: w}).Print(ctx, pc, lastCmd)
	if w.OnShowBuildInfo != nil {
		w.OnShowBuildInfo(ctx, lastCmd, pc.PositionalArgs())
		return
	}

	// initTabStop(defaultTabStop)
	//
	// w.printHeader(w.currentHelpPainter, &w.rootCommand.Command)
	//
	// fp("")

	if conf.GoVersion != "" {
		fp(`           Built by: %v`, conf.GoVersion)
	}

	ts := conf.Buildstamp
	if ts == "" {
		ts = time.Now().UTC().Format(time.RFC3339)
	}
	dt, err := times.SmartParseTime(ts)
	// dt, err := time.Parse("", ts)
	if err == nil {
		ts = dt.Format(time.RFC3339)
	}
	fp(`    Built Timestamp: %v`, ts)

	if lastCmd.HitTitle() == "V" {
		return
	}

	if conf.BuilderComments != "" {
		fp(`   Builder Comments: %v`, conf.BuilderComments)
	}

	if conf.Serial != "" {
		fp(`       Built Serial: %v`, conf.Serial)
	}
	if conf.SerialString != "" {
		fp(`Built Serial String: %v`, conf.SerialString)
	}

	fp(`
         Git Commit: %v
        Git Summary: %v
    Git Description: %v
`, conf.Githash, conf.GitSummary, conf.GitDesc)

	// w.currentHelpPainter.FpPrintHelpTailLine(lastCmd)
	return
}

func fp(fmtStr string, args ...interface{}) {
	s := fmt.Sprintf(fmtStr, args...)
	needln := !strings.HasSuffix(s, "\n")
	_fpz(needln, s)
}

// func fpK(fmtStr string, args ...interface{}) {
// 	s := fmt.Sprintf(fmtStr, args...)
// 	_fpz(false, s)
// }

func _fpz(needln bool, s string) {
	var w io.Writer = os.Stdout
	if wkr := UniqueWorker(); wkr != nil {
		if r, ok := wkr.(interface{ GetHelpScreenWriter() HelpWriter }); ok {
			if hsw := r.GetHelpScreenWriter(); hsw != nil {
				w = hsw
			}
		}
		if w != nil {
			_fpzz(needln, s, w)
		}
	} else {
		_fpzz(needln, s, w)
	}
}

func _fpzz(needln bool, s string, w io.Writer) {
	_, _ = fmt.Fprintf(w, "%s", s)
	if needln {
		_, _ = fmt.Fprintln(w)
	}
}

func (w *workerS) showSBOM(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd, args ...any) (err error) { //nolint:revive,unused
	if w.OnSBOM != nil {
		err = w.OnSBOM(ctx, lastCmd, pc.PositionalArgs())
		return
	}
	err = (&sbomS{}).onAction(ctx, lastCmd, pc.PositionalArgs())
	return
}

func (w *workerS) showHelpScreen(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd, args ...any) (err error) { //nolint:revive,unused
	(&helpPrinter{w: w}).Print(ctx, pc, lastCmd, args...)
	return
}

func (w *workerS) showHelpScreenAsMan(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd, args ...any) (err error) { //nolint:revive,unused
	hp := &helpPrinter{w: w, asManual: true}

	// err = w.pipeToReader(func(stdout io.Writer) {
	// 	ws1 := &ws{stdout}
	// 	hp.PrintTo(ctx, ws1, pc, lastCmd)
	// },
	// 	"man", "1",
	// 	// "less", os.Getenv("LESS"),
	// )

	var f *os.File
	f, err = os.CreateTemp(os.TempDir(), lastCmd.GetDottedPath()+".man.1")
	if err != nil {
		return
	}
	defer func() {
		err = f.Close()
		if keep := lastCmd.Root().Store().MustBool("keep"); keep {
			fmt.Printf("manpage file %q kept.\n", f.Name())
		} else {
			err = os.Remove(f.Name())
		}
	}()

	hp.PrintTo(ctx, f, pc, lastCmd, args...)

	// close xxx.man.1, so that we can launch it into as parameter of `man` command
	_ = f.Close()

	// logz.InfoContext(ctx, "temp manpage written", "path", f.Name())

	program, fname := "man", ""
	fname, err = exec.LookPath(program)
	if err == nil {
		program, err = filepath.Abs(fname)
	}
	if err != nil {
		return
	}

	// logz.InfoContext(ctx, "manpage", "man", program, "path", f.Name())
	cmd := exec.Command(program, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}

// helpSystemAction is the reaction for 'help' command at root level.
func (w *workerS) helpSystemAction(ctx context.Context, cmd cli.Cmd, args []string) (err error) { //nolint:revive,unused
	if len(args) > 0 {
		hp, handled := &helpPrinter{w: w}, cmd
		// trying to recognize the given commands and print help screen of it.
		if handled, err = hs.New(w, cmd, args).FindCmd(ctx, cmd, args); handled == nil {
			return
		}
		hp.Print(ctx, w.parsingCtx, handled)
		return
	}

	// entering an interactive shell mode and listen to the user's commands.
	err = hs.New(w, cmd, args).Run(ctx)
	return
}

func (w *workerS) showTree(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd, args ...any) (err error) { //nolint:revive,unused
	(&helpPrinter{w: w, debugMatches: true, treeMode: true}).Print(ctx, pc, lastCmd, args...)
	return
}

func (w *workerS) showDebugScreen(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd, args ...any) (err error) { //nolint:revive,unused
	(&helpPrinter{w: w, debugScreenMode: true, debugMatches: true}).Print(ctx, pc, lastCmd, args...)
	return
}
