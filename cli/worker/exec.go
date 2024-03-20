package worker

import (
	"context"
	"sync/atomic"

	errorsv3 "gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2/cli"
	logz "github.com/hedzr/logg/slog"
)

func (w *workerS) exec(ctx *parseCtx) (err error) {
	lastCmd := ctx.LastCmd()

	w.parsingCtx = ctx // save ctx for later, OnAction might need it.

	if w.Store().Has("app.force-default-action") {
		ctx.forceDefaultAction = w.Store().MustBool("app.force-default-action", false)
	}

	var deferActions func(errInvoked error)
	if deferActions, err = w.beforeExec(ctx, lastCmd); err != nil {
		return
	}
	defer func() {
		deferActions(err) // err must be delayed caught here
		if err = w.afterExec(ctx, lastCmd); err == nil {
			c := context.Background()
			err = w.finalActions(c, ctx, lastCmd)
		}
	}()

	if !ctx.forceDefaultAction && lastCmd.HasOnAction() {
		logz.Verbose("invoke action of cmd, with args", "cmd", lastCmd, "args", ctx.positionalArgs)
		err = lastCmd.Invoke(ctx.positionalArgs)
		if !w.errIsSignalFallback(err) {
			return
		}
		ctx.forceDefaultAction, err = true, nil
	}

	handled, err1 := w.handleActions(ctx)
	// for k, action := range w.actions {
	// 	if k&w.actionsMatched != 0 {
	// 		logz.Verbose("Invoking worker.actionsMatched", "hit-action", k, "actions", w.Actions())
	// 		err, handled = action(ctx, lastCmd), true
	// 		break
	// 	}
	// }
	if handled || !w.errIsSignalOrNil(err1) {
		return
	}

	// if ctx.helpScreen {
	// 	err = w.onPrintHelpScreen(ctx, lastCmd)
	// 	return
	// }

	if ctx.forceDefaultAction {
		err = w.onDefaultAction(ctx, lastCmd)
		return
	}

	logz.Verbose("no onAction associate to cmd", "cmd", lastCmd)
	err = w.onPrintHelpScreen(ctx, lastCmd)
	return
}

func (w *workerS) handleActions(ctx *parseCtx) (handled bool, err error) {
	lastCmd := ctx.LastCmd()
	for k, action := range w.actions {
		if k&w.actionsMatched != 0 {
			logz.Verbose("Invoking worker.actionsMatched", "hit-action", k, "actions", w.Actions())
			err, handled = action(ctx, lastCmd), true
			break
		}
	}
	return
}

type onAction func(ctx *parseCtx, lastCmd *cli.Command) (err error)

func (w *workerS) beforeExec(ctx *parseCtx, lastCmd *cli.Command) (deferActions func(errInvoked error), err error) {
	err = w.checkRequiredFlags(ctx, lastCmd)
	deferActions = func(error) {}
	if err != nil {
		return
	}

	if lastCmd != w.root.Command {
		deferActions, err = w.root.RunPreActions(lastCmd, ctx.positionalArgs)
	}
	return
}

func (w *workerS) checkRequiredFlags(ctx *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive
	lastCmd.WalkBackwards(func(cc *cli.Command, ff *cli.Flag, index, count, level int) {
		if ff != nil {
			if ff.Required() && ff.GetTriggeredTimes() < 0 {
				err = ErrRequiredFlag.FormatWith(ff, lastCmd)
				_, _, _, _, _ = ctx, cc, index, count, level
				return
			}
		}
	})
	return
}

func (w *workerS) afterExec(ctx *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive
	_, _ = ctx, lastCmd
	return
}

func (w *workerS) finalActions(ctx context.Context, pc *parseCtx, lastCmd *cli.Command) (err error) { //nolint:revive,unparam
	_, _ = pc, lastCmd
	err = w.writeBackToLoaders(ctx)
	return
}

func (w *workerS) onPrintHelpScreen(ctx *parseCtx, lastCmd *cli.Command) (err error) { //nolint:unparam
	(&helpPrinter{w: w}).Print(ctx, lastCmd)
	return
}

func (w *workerS) onDefaultAction(ctx *parseCtx, lastCmd *cli.Command) (err error) { //nolint:unparam
	(&helpPrinter{w: w, debugMatches: true}).Print(ctx, lastCmd)
	return
}

func (w *workerS) onPassThruCharMatched(ctx *parseCtx) (err error) { //nolint:unparam
	if atomic.CompareAndSwapInt32(&ctx.passThruMatched, 0, 1) {
		atomic.StoreInt32(&ctx.passThruMatched, int32(ctx.i))
	}
	return
}

func (w *workerS) onSingleHyphenMatched(ctx *parseCtx) (err error) { //nolint:unparam
	if atomic.CompareAndSwapInt32(&ctx.singleHyphenMatched, 0, 1) {
		atomic.StoreInt32(&ctx.singleHyphenMatched, int32(ctx.i))
	}
	return
}

func (w *workerS) onUnknownCommandMatched(ctx *parseCtx) (err error) {
	logz.Warn("UNKNOWN <mark>Command</mark> FOUND", "arg", ctx.arg)
	err = ErrUnmatchedCommand.FormatWith(ctx.arg, ctx.LastCmd())
	return
}

func (w *workerS) onUnknownFlagMatched(ctx *parseCtx) (err error) {
	logz.Warn("UNKNOWN <mark>Flag</mark> FOUND", "arg", ctx.arg)
	err = ErrUnmatchedFlag.FormatWith(ctx.arg, ctx.LastCmd())
	return
}

//

//

//

var (
	// ErrUnmatchedCommand means Unmatched command found. It's just a state, not a real error, see [cli.Config.UnmatchedAsError]
	ErrUnmatchedCommand = errorsv3.New("UNKNOWN Command FOUND: %q | cmd=%v")
	// ErrUnmatchedFlag means Unmatched flag found. It's just a state, not a real error, see [cli.Config.UnmatchedAsError]
	ErrUnmatchedFlag = errorsv3.New("UNKNOWN Flag FOUND: %q | cmd=%v")
	// ErrRequiredFlag means required flag must be set explicitly
	ErrRequiredFlag = errorsv3.New("Flag %q is REQUIRED | cmd=%v")
)
