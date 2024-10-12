package worker

import (
	"context"
	"sync/atomic"

	errorsv3 "gopkg.in/hedzr/errors.v3"

	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/cmdr/v2/cli"
)

func (w *workerS) exec(ctx context.Context, pc *parseCtx) (err error) {
	lastCmd := pc.LastCmd()

	w.parsingCtx = pc // save pc for later, OnAction might need it.

	forceDefaultAction := pc.forceDefaultAction

	var deferActions func(errInvoked error)
	if deferActions, err = w.beforeExec(ctx, pc, lastCmd); err != nil {
		return
	}
	defer func() {
		deferActions(err) // err must be delayed caught here
		if err1 := w.afterExec(ctx, pc, lastCmd); err1 == nil {
			c := context.Background()
			if err1 = w.finalActions(c, pc, lastCmd); err1 != nil {
				ec := errorsv3.New()
				ec.Attach(err, err1)
				ec.Defer(&err)
			}
		}
	}()

	if !forceDefaultAction && lastCmd.HasOnAction() {
		logz.Verbose("[cmdr] invoke action of cmd, with args", "cmd", lastCmd, "args", pc.positionalArgs)
		err = lastCmd.Invoke(ctx, pc.positionalArgs)
		logz.Verbose("[cmdr] invoke action ends.", "err", err)
		if !w.errIsSignalFallback(err) {
			return
		}
		pc.forceDefaultAction, err = true, nil
	}

	handled, err1 := w.handleActions(ctx, pc)
	// for k, action := range w.actions {
	// 	if k&w.actionsMatched != 0 {
	// 		logz.Verbose("[cmdr] Invoking worker.actionsMatched", "hit-action", k, "actions", w.Actions())
	// 		err, handled = action(pc, lastCmd), true
	// 		break
	// 	}
	// }
	if handled || !w.errIsSignalOrNil(err1) {
		err = err1
		return
	}

	// if pc.helpScreen {
	// 	err = w.onPrintHelpScreen(pc, lastCmd)
	// 	return
	// }

	if pc.forceDefaultAction {
		err = w.onDefaultAction(ctx, pc, lastCmd)
		return
	}

	logz.Verbose("[cmdr] no onAction associate to cmd", "cmd", lastCmd)
	err = w.onPrintHelpScreen(ctx, pc, lastCmd)
	return
}

func (w *workerS) handleActions(ctx context.Context, pc *parseCtx) (handled bool, err error) {
	lastCmd := pc.LastCmd()
	for k, action := range w.actions {
		if k&w.actionsMatched != 0 {
			logz.Verbose("[cmdr] Invoking worker.actionsMatched", "hit-action", k, "actions", w.Actions())
			err, handled = action(ctx, pc, lastCmd), true
			break
		}
	}
	return
}

type onAction func(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error)

func (w *workerS) beforeExec(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (deferActions func(errInvoked error), err error) {
	err = w.checkRequiredFlags(ctx, pc, lastCmd)
	deferActions = func(error) {}
	if err != nil {
		return
	}

	if lastCmd != w.root.Command {
		deferActions, err = w.root.RunPreActions(ctx, lastCmd, pc.positionalArgs)
	}
	return
}

func (w *workerS) checkRequiredFlags(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive
	lastCmd.WalkBackwards(ctx, func(ctx context.Context, pc *cli.WalkBackwardsCtx, cc *cli.Command, ff *cli.Flag, index, groupIndex, count, level int) {
		if ff != nil {
			if ff.Required() && ff.GetTriggeredTimes() < 0 {
				err = cli.ErrRequiredFlag.FormatWith(ff, lastCmd)
				_, _, _, _, _, _ = pc, cc, index, groupIndex, count, level
				return
			}
		}
	})
	return
}

func (w *workerS) afterExec(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive
	_, _, _ = ctx, pc, lastCmd
	return
}

func (w *workerS) finalActions(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unparam
	_, _ = pc, lastCmd
	err = w.writeBackToLoaders(ctx)
	return
}

func (w *workerS) onPrintHelpScreen(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:unparam
	(&helpPrinter{w: w}).Print(ctx, pc, lastCmd)
	return
}

func (w *workerS) onDefaultAction(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:unparam
	(&helpPrinter{w: w, debugMatches: true}).Print(ctx, pc, lastCmd)
	return
}

func (w *workerS) onPassThruCharMatched(pc *parseCtx) (err error) { //nolint:unparam
	if atomic.CompareAndSwapInt32(&pc.passThruMatched, 0, 1) {
		atomic.StoreInt32(&pc.passThruMatched, int32(pc.i))
	}
	return
}

func (w *workerS) onSingleHyphenMatched(pc *parseCtx) (err error) { //nolint:unparam
	if atomic.CompareAndSwapInt32(&pc.singleHyphenMatched, 0, 1) {
		atomic.StoreInt32(&pc.singleHyphenMatched, int32(pc.i))
	}
	return
}

func (w *workerS) onUnknownCommandMatched(pc *parseCtx) (err error) {
	logz.Warn("[cmdr] UNKNOWN <mark>Command</mark> FOUND", "arg", pc.arg)
	err = cli.ErrUnmatchedCommand.FormatWith(pc.arg, pc.LastCmd())
	return
}

func (w *workerS) onUnknownFlagMatched(pc *parseCtx) (err error) {
	logz.Warn("[cmdr] UNKNOWN <mark>Flag</mark> FOUND", "arg", pc.arg)
	err = cli.ErrUnmatchedFlag.FormatWith(pc.arg, pc.LastCmd())
	return
}

//

//

//
