package worker

import (
	"sync/atomic"

	errorsv3 "gopkg.in/hedzr/errors.v3"

	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/cmdr/v2/cli"
)

func (w *workerS) parse(ctx *parseCtx) (err error) { //nolint:revive
	defer func() {
		if len(w.tasksAfterParse) > 0 {
			ec := errorsv3.New("tasks failed")
			defer ec.Defer(&err)
			for _, task := range w.tasksAfterParse {
				if task != nil {
					ec.Attach(task(w, ctx, err))
				}
			}
		}
	}()

loopArgs:
	for ctx.i = 1; ctx.i < len(w.args); ctx.i++ {
		if w.args[ctx.i] == "" {
			continue
		}

		if atomic.LoadInt32(&ctx.passThruMatched) > 0 || errorsv3.Is(err, cli.ErrShouldStop) || w.errIsUnmatchedArg(err) {
			ctx.positionalArgs = append(ctx.positionalArgs, w.args[ctx.i])
			logz.Verbose("positional args added", "i", ctx.i, "args", ctx.positionalArgs)
			continue
		}

		logz.Verbose("parsing command-line args", "i", ctx.i, "arg", w.args[ctx.i])

		ctx.arg, ctx.short, ctx.pos = w.args[ctx.i], false, 0
		switch c1 := ctx.arg[0]; c1 {
		case '+': // for short title only
			if len(ctx.arg) > 1 {
				// try matching a short flag
				plusFound := func(cb func(w *workerS, ctx *parseCtx) error) error {
					ctx.prefixPlusSign.Swap(true)
					defer func() { ctx.prefixPlusSign.Swap(false) }()
					return cb(w, ctx)
				}
				if err = plusFound(func(w *workerS, ctx *parseCtx) error {
					ctx.arg, ctx.short, ctx.dblTilde = ctx.arg[1:], true, false
					return w.matchFlag(ctx, true)
				}); !w.errIsSignalOrNil(err) {
					err = w.onUnknownFlagMatched(ctx)
					break loopArgs
				}
				continue
			}
			// single '+': as a positional arg
			ctx.positionalArgs = append(ctx.positionalArgs, ctx.arg)
			logz.Verbose("positional args added", "i", ctx.i, "args", ctx.positionalArgs)
			continue

		case '-', '~':
			if len(ctx.arg) > 1 { // so, ctx.arg >= 2
				if (c1 == '-' && ctx.arg[1] == '-') || (c1 == '~' && ctx.arg[1] == '~') {
					if len(ctx.arg) == 3 && ctx.arg[2] == '-' { //nolint:revive
						// --: pass-thru found
						err = w.onPassThruCharMatched(ctx)
						continue
					}

					// try match a long flag
					ctx.arg, ctx.short, ctx.dblTilde = ctx.arg[2:], false, c1 == '~'
					if err = w.matchFlag(ctx, false); !w.errIsSignalOrNil(err) {
						err = w.onUnknownFlagMatched(ctx)
						break loopArgs
					}
					continue
				}
				if (c1 == '-' && ctx.arg[1] == '~') || (c1 == '~' && ctx.arg[1] == '-') {
					err = w.onUnknownFlagMatched(ctx)
					break loopArgs
				}

				// try matching a short flag
				ctx.arg, ctx.short, ctx.dblTilde = ctx.arg[1:], true, false
				if c1 != '-' {
					err = w.onUnknownFlagMatched(ctx)
					break loopArgs
				}
				if err = w.matchFlag(ctx, true); !w.errIsSignalOrNil(err) {
					err = w.onUnknownFlagMatched(ctx)
					break loopArgs
				}
				continue
			}
			// single '-' matched, is it a wrong state?
			err = w.onSingleHyphenMatched(ctx)
			continue

		default:
			if ctx.NoCandidateChildCommands() {
				ctx.positionalArgs = append(ctx.positionalArgs, ctx.arg)
				logz.Verbose("positional args added", "i", ctx.i, "args", ctx.positionalArgs)
				continue
			}
			if err = w.matchCommand(ctx); !w.errIsSignalOrNil(err) {
				err = w.onUnknownCommandMatched(ctx)
			}
		}
	}
	return
}

func (w *workerS) matchCommand(ctx *parseCtx) (err error) {
	err = ErrUnmatchedCommand
	cmd := ctx.LastCmd()
	if short, cc := cmd.Match(ctx.arg); cc != nil {
		ms, handled := ctx.addCmd(cc, short), false
		handled, err = cc.TryOnMatched(0, ms)
		if err == nil {
			ctx.lastCommand, err = len(ctx.matchedCommands)-1, nil
		}
		logz.Verbose("command matched", "short", short, "cmd", ctx.LastCmd(), "handled", handled)
	}
	return
}

func (w *workerS) matchFlag(ctx *parseCtx, short bool) (err error) {
	err = ErrUnmatchedFlag
	cmd, vp := ctx.LastCmd(), cli.NewFVP(w.args[ctx.i+1:], ctx.arg, short, ctx.prefixPlusSign.Load(), ctx.dblTilde)
	// defer func() { ctx.i, vp.AteArgs = ctx.i+vp.AteArgs, 0 }()

compactFlags:
	ff, err1 := cmd.MatchFlag(vp)
	if vp.Matched != "" && ff != nil && w.errIsSignalOrNil(err1) {
		ms, handled := ctx.addFlag(ff), false
		handled, err1 = ff.TryOnMatched(0, ms)
		logz.Verbose("flag matched", "short", vp.Short, "flg", ff, "val-pkg-val", ff.DefaultValue(), "handled", handled)

		ctx.i += vp.AteArgs
		vp.AteArgs = 0
		err = err1

		if vp.Remains != "" && vp.PartialMatched {
			ctx.arg = vp.Remains
			vp.Reset()
			ctx.prefixPlusSign.Swap(false)

			if !errorsv3.Is(err, cli.ErrShouldStop) {
				goto compactFlags // try matching next compact flag. eg: '-avz' => '-a' parsed, remains '-vz'
			}
		}
	}
	return
}
