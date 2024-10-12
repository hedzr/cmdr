package worker

import (
	"context"

	"gopkg.in/hedzr/errors.v3"

	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/internal/hs"
)

func (w *workerS) showVersion(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	// (&helpPrinter{w: w}).Print(ctx, pc, lastCmd)
	return
}

func (w *workerS) showBuiltInfo(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	// (&helpPrinter{w: w}).Print(ctx, pc, lastCmd)
	return
}

func (w *workerS) showSBOM(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	// (&helpPrinter{w: w}).Print(ctx, pc, lastCmd)
	return
}

func (w *workerS) showHelpScreen(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	(&helpPrinter{w: w}).Print(ctx, pc, lastCmd)
	return
}

func (w *workerS) showHelpScreenAsMan(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	(&helpPrinter{w: w}).Print(ctx, pc, lastCmd)
	return
}

// helpSystemAction is the reaction for 'help' command at root level.
func (w *workerS) helpSystemAction(ctx context.Context, cmd *cli.Command, args []string) (err error) { //nolint:revive,unused
	if len(args) > 0 {
		// trying to recognize the given commands and print help screen of it.
		cc := cmd.Root().Command
		for _, arg := range args {
			cc = cc.FindSubCommand(ctx, arg, true)
			if cc == nil {
				logz.Error("[cmdr] Unknown command found.", "commands", args)
				return errors.New("unknown command %v found", args)
			}
		}
		(&helpPrinter{w: w}).Print(ctx, w.parsingCtx, cc)
		return
	}

	// entering an interactive shell mode and listen to the user's commands.
	err = hs.New(w, cmd, args).Run(ctx)
	return
}

func (w *workerS) showTree(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	(&helpPrinter{w: w, debugMatches: true, treeMode: true}).Print(ctx, pc, lastCmd)
	return
}

func (w *workerS) showDebugScreen(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unused
	(&helpPrinter{w: w, debugScreenMode: true, debugMatches: true}).Print(ctx, pc, lastCmd)
	return
}
