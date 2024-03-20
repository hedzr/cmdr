package worker

import (
	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/internal/hs"
	logz "github.com/hedzr/logg/slog"
)

func (w *workerS) showVersion(ctx *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	// (&helpPrinter{w: w}).Print(ctx, lastCmd)
	return
}

func (w *workerS) showBuiltInfo(ctx *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	// (&helpPrinter{w: w}).Print(ctx, lastCmd)
	return
}

func (w *workerS) showSBOM(ctx *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	// (&helpPrinter{w: w}).Print(ctx, lastCmd)
	return
}

func (w *workerS) showHelpScreen(ctx *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	(&helpPrinter{w: w}).Print(ctx, lastCmd)
	return
}

func (w *workerS) showHelpScreenAsMan(ctx *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	(&helpPrinter{w: w}).Print(ctx, lastCmd)
	return
}

// helpSystemAction is the reaction for 'help' command at root level.
func (w *workerS) helpSystemAction(cmd *cli.Command, args []string) (err error) { //nolint:revive,unused
	if len(args) > 0 {
		// trying to recognize the given commands and print help screen of it.
		cc := cmd.Root().Command
		for _, arg := range args {
			cc = cc.FindSubCommand(arg, true)
			if cc == nil {
				logz.Error("Unknown command found.", "commands", args)
				return errors.New("unknown command %v found", args)
			}
		}
		(&helpPrinter{w: w}).Print(w.parsingCtx, cc)
		return
	}

	// entering an interactive shell mode and listen to the user's commands.
	err = hs.New(w, cmd, args).Run()
	return
}

func (w *workerS) showTree(ctx *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	(&helpPrinter{w: w, debugMatches: true, treeMode: true}).Print(ctx, lastCmd)
	return
}

func (w *workerS) showDebugScreen(ctx *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	(&helpPrinter{w: w, debugScreenMode: true, debugMatches: true}).Print(ctx, lastCmd)
	return
}
