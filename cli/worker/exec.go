package worker

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync/atomic"

	errorsv3 "gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/logz"
	"github.com/hedzr/is/exec"
)

func (w *workerS) SetTasksAfterRun(tasks ...taskAfterRun) {
	w.tasksAfterRun = append(w.tasksAfterRun, tasks...)
}

func (w *workerS) exec(ctx context.Context, pc *parseCtx) (err error) {
	if w.DontExecuteAction {
		return
	}

	lastCmd := pc.LastCmd()
	logz.VerboseContext(ctx, "[cmdr] exec...", "last-matched-cmd", lastCmd)

	forceDefaultAction := pc.forceDefaultAction

	defer func() {
		c := context.Background()
		if err1 := w.finalActions(c, pc, lastCmd); err1 != nil {
			ec := errorsv3.New()
			ec.Attach(err, err1)
			ec.Defer(&err)
		}

		if sb, ok := ctx.Value(cli.CtxKeyHelpScreenWriter).(*strings.Builder); ok {
			for _, task := range w.tasksAfterRun {
				if task != nil {
					err = task(c, w, pc, err, sb)
				}
			}
		}
	}()
	return w.execCmd(ctx, pc, lastCmd, forceDefaultAction)
}

func (w *workerS) execCmd(ctx context.Context, pc *parseCtx, cmd cli.Cmd, forceDefaultAction bool) (err error) {
	var deferActions func(errInvoked error)
	if deferActions, err = w.beforeExec(ctx, pc, cmd); err != nil {
		return
	}
	defer func() {
		deferActions(err) // err must be delayed caught here
		if err1 := w.afterExec(ctx, pc, cmd); err1 == nil {
			// c := context.Background()
			// if err1 = w.finalActions(c, pc, cmd); err1 != nil {
			// 	ec := errorsv3.New()
			// 	ec.Attach(err, err1)
			// 	ec.Defer(&err)
			// }
			ec := errorsv3.New()
			ec.Attach(err, err1)
			ec.Defer(&err)
		}
	}()

	handled, err1 := w.handleActions(ctx, pc)
	if handled || !w.errIsSignalOrNil(err1) {
		err = err1
		return
	}

	if is := cmd.InvokeShell(); is != "" {
		if !cmd.CanInvoke() {
			if sh := cmd.Shell(); sh != "" {
				err = exec.New().WithCommand(sh, "-c", is).RunAndCheckError()
				return
			}
			err = exec.New().WithCommandString(is).RunAndCheckError()
			return
		}
	}
	if ip := cmd.InvokeProc(); ip != "" {
		if !cmd.CanInvoke() {
			if strings.Contains(ip, "sh -c ") {
				err = exec.New().WithCommandString(ip).RunAndCheckError()
			} else {
				err = exec.Call(ip, nil)
			}
			return
		}
	}

	if !forceDefaultAction && cmd.CanInvoke() {
		logz.VerboseContext(ctx, "invoke action of cmd, with args", "cmd", cmd, "args", pc.positionalArgs)
		err = cmd.Invoke(ctx, pc.positionalArgs)
		logz.VerboseContext(ctx, "invoke action ends.", "err", err)
		if !w.errIsSignalFallback(err) {
			return
		}
		pc.forceDefaultAction, err = true, nil
	}

	if redirectTo := cmd.RedirectTo(); redirectTo != "" {
		if cc, ff := cli.DottedPathToCommandOrFlag1(redirectTo, cmd.Root()); cc != nil && ff == nil {
			if cmd1, ok := cc.(*cli.CmdS); ok {
				logz.VerboseContext(ctx, "invoke action of cmd (redirect), with args", "cmd", cmd1, "args", pc.positionalArgs)
				err = w.execCmd(ctx, pc, cmd1, forceDefaultAction)
				logz.VerboseContext(ctx, "invoke action ends (redirect).", "err", err)
				return
			}
		}
	}

	if pc.forceDefaultAction {
		err = w.onDefaultAction(ctx, pc, cmd)
		return
	}

	logz.VerboseContext(ctx, "[cmdr] no onAction associated to cmd", "cmd", cmd)
	err = w.onPrintHelpScreen(ctx, pc, cmd)
	return
}

func (w *workerS) DoBuiltinAction(ctx context.Context, action cli.ActionEnum, args ...any) (handled bool, err error) {
	for k, handler := range w.actions {
		if k&action != 0 {
			logz.VerboseContext(ctx, "[cmdr] Invoking action", "hit-action", k, "actions", w.Actions())
			pc, lastCmd := w.parsingCtx, w.parsingCtx.LastCmd()
			if len(args) > 0 {
				if c, ok := args[0].(cli.Cmd); ok {
					lastCmd = c
					args = args[1:]
				}
			}
			err, handled = handler(ctx, pc.(*parseCtx), lastCmd, args...), true
			break
		}
	}
	return
}

func (w *workerS) handleActions(ctx context.Context, pc *parseCtx) (handled bool, err error) {
	lastCmd := pc.LastCmd()
	for k, action := range w.actions {
		if k&w.actionsMatched != cli.ActionNone {
			logz.VerboseContext(ctx, "[cmdr] Invoking worker.actionsMatched", "hit-action", k, "actions", w.Actions())
			err, handled = action(ctx, pc, lastCmd), true
			break
		}
	}
	return
}

func (w *workerS) beforeExec(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd) (deferActions func(errInvoked error), err error) {
	deferActions = func(error) {}
	err = w.checkRequiredFlags(ctx, pc, lastCmd)
	if err != nil {
		return
	}
	err = w.checkValidArgs(ctx, pc, lastCmd)
	if err != nil {
		return
	}

	if lastCmd != w.root.Cmd {
		if cx, ok := w.root.Cmd.(*cli.CmdS); ok {
			deferActions, err = cx.RunPreActions(ctx, lastCmd, pc.positionalArgs)
		}
	}
	return
}

func (w *workerS) checkValidArgs(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd) (err error) { //nolint:revive
	for ff, ms := range pc.matchedFlags {
		val := fmt.Sprintf("%v", ms.Value)
		if ff != nil && len(ff.ValidArgs()) > 0 && slices.Contains(ff.ValidArgs(), val) == false {
			err = cli.ErrValidArgs.FormatWith(ff, ff.ValidArgs(), lastCmd)
		}
	}
	_ = ctx
	return
}

func (w *workerS) checkRequiredFlags(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd) (err error) { //nolint:revive
	wbc := &cli.WalkBackwardsCtx{
		Group: true,
		Sort:  false,
	}
	_ = pc
	lastCmd.WalkBackwardsCtx(ctx, func(ctx context.Context, pc *cli.WalkBackwardsCtx, cc cli.Cmd, ff *cli.Flag, index, groupIndex, count, level int) {
		if ff != nil {
			if ff.Required() && ff.GetTriggeredTimes() <= 0 {
				err = cli.ErrRequiredFlag.FormatWith(ff, lastCmd)
				_, _, _, _, _, _ = pc, cc, index, groupIndex, count, level
				return
			}
		}
	}, wbc)
	return
}

func (w *workerS) afterExec(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd) (err error) { //nolint:revive
	_, _, _ = ctx, pc, lastCmd
	return
}

func (w *workerS) finalActions(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd) (err error) { //nolint:revive,unparam
	_, _ = pc, lastCmd
	err = w.writeBackToLoaders(ctx)
	return
}

func (w *workerS) onPrintHelpScreen(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd, args ...any) (err error) { //nolint:unparam
	if cc := pc.LastCmd(); cc.RedirectTo() != "" {
		(&helpPrinter{w: w}).Print(ctx, pc, cc, args...)
		return
	}

	(&helpPrinter{w: w}).Print(ctx, pc, lastCmd, args...)
	return
}

func (w *workerS) onDefaultAction(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd, args ...any) (err error) { //nolint:unparam
	(&helpPrinter{w: w, debugMatches: true}).Print(ctx, pc, lastCmd, args...)
	return
}

func (w *workerS) runHelpSystem(ctx context.Context, pc *parseCtx, lastCmd cli.Cmd, args ...any) (err error) { //nolint:unparam
	err = w.helpSystemAction(ctx, lastCmd, pc.PositionalArgs())
	return
}

func (w *workerS) onPassThruCharMatched(ctx context.Context, pc *parseCtx) (err error) { //nolint:unparam
	if atomic.CompareAndSwapInt32(&pc.passThruMatched, 0, 1) {
		atomic.StoreInt32(&pc.passThruMatched, int32(pc.i))
		if w.OnPassThruCharMatched != nil {
			err = w.OnPassThruCharMatched(ctx, pc)
		}
	}
	return
}

func (w *workerS) onSingleHyphenMatched(ctx context.Context, pc *parseCtx) (err error) { //nolint:unparam
	if atomic.CompareAndSwapInt32(&pc.singleHyphenMatched, 0, 1) {
		atomic.StoreInt32(&pc.singleHyphenMatched, int32(pc.i))
		if w.OnSingleHyphenMatched != nil {
			err = w.OnSingleHyphenMatched(ctx, pc)
		}
	}
	return
}

func (w *workerS) onUnknownCommandMatched(ctx context.Context, pc *parseCtx) (err error) {
	// if ignoreTestArgs && strings.HasPrefix(pc.arg,"test."){
	// 	return
	// }
	err = cli.ErrUnmatchedCommand.FormatWith(pc.arg, pc.LastCmd())
	if w.OnUnknownCommandHandler != nil {
		err = w.OnUnknownCommandHandler(ctx, pc.arg, pc.LastCmd(), err)
	}
	if err == nil || w.errIsSignalFallback(err) {
		return
	}

	logz.WarnContext(ctx, "[cmdr] UNKNOWN <mark>CmdS</mark> FOUND", "arg", pc.arg, "all-args", *pc.argsPtr)
	return
}

func (w *workerS) onUnknownFlagMatched(ctx context.Context, pc *parseCtx) (err error) {
	if ignoreTestArgs && strings.HasPrefix(pc.arg, "test.") {
		return
	}
	err = cli.ErrUnmatchedFlag.FormatWith(pc.arg, pc.LastCmd())
	if w.OnUnknownFlagHandler != nil {
		err = w.OnUnknownFlagHandler(ctx, pc.arg, pc.LastCmd(), err)
	}
	if err == nil || w.errIsSignalFallback(err) {
		return
	}

	logz.WarnContext(ctx, "[cmdr] UNKNOWN <mark>Flag</mark> FOUND", "arg", pc.arg, "all-args", *pc.argsPtr)
	return
}

const ignoreTestArgs = true
