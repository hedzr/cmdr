package worker

import (
	"context"
	"sync/atomic"

	errorsv3 "gopkg.in/hedzr/errors.v3"

	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/cmdr/v2/cli"
)

func (w *workerS) SetTasksAfterParse(tasks ...taskAfterParse) {
	w.tasksAfterParse = append(w.tasksAfterParse, tasks...)
}

func (w *workerS) parse(ctx context.Context, pc *parseCtx) (err error) { //nolint:revive
	ec := errorsv3.New("tasks failed")
	defer ec.Defer(&err)

	defer func() {
		if len(w.tasksAfterParse) > 0 {
			for _, task := range w.tasksAfterParse {
				if task != nil {
					ec.Attach(task(w, pc, err))
				}
			}
		}
	}()

	logz.VerboseContext(ctx, "[cmdr] parsing command line args ...", "args", w.args)

loopArgs:
	for pc.i = 1; pc.i < len(w.args); pc.i++ {
		if w.args[pc.i] == "" {
			continue
		}

		if atomic.LoadInt32(&pc.passThruMatched) > 0 || errorsv3.Is(err, cli.ErrShouldStop) || w.errIsUnmatchedArg(err) {
			pc.positionalArgs = append(pc.positionalArgs, w.args[pc.i])
			logz.VerboseContext(ctx, "[cmdr] positional args added", "i", pc.i, "args", pc.positionalArgs)
			continue
		}

		logz.VerboseContext(ctx, "[cmdr] parsing command-line args", "i", pc.i, "arg", w.args[pc.i])

		pc.arg, pc.short, pc.pos = w.args[pc.i], false, 0
		switch c1 := pc.arg[0]; c1 {
		// TODO need more design for form '+flag'.
		// currently, +flag is designed as a bool value flipper
		case '+': // for bool flag it's a flipper;
			if len(pc.arg) > 1 {
				if w.interpretLeadingPlusSign(pc) {
					continue
				}

				// try matching a short flag
				plusFound := func(cb func(w *workerS, ctx *parseCtx) error) error {
					pc.prefixPlusSign.Swap(true)
					defer func() { pc.prefixPlusSign.Swap(false) }()
					return cb(w, pc)
				}
				if err = plusFound(func(w *workerS, pc *parseCtx) error {
					pc.arg, pc.short, pc.dblTilde = pc.arg[1:], true, false
					return w.matchFlag(ctx, pc, true)
				}); !w.errIsSignalOrNil(err) {
					err = w.onUnknownFlagMatched(ctx, pc)
					break loopArgs
				}
				continue
			}
			// single '+': as a positional arg
			pc.positionalArgs = append(pc.positionalArgs, pc.arg)
			logz.VerboseContext(ctx, "[cmdr] positional args added", "i", pc.i, "args", pc.positionalArgs)
			continue

		case '-', '~':
			if len(pc.arg) > 1 { // so, pc.arg >= 2
				if (c1 == '-' && pc.arg[1] == '-') || (c1 == '~' && pc.arg[1] == '~') {
					if len(pc.arg) == 3 && pc.arg[2] == '-' { //nolint:revive
						// --: pass-thru found
						err = w.onPassThruCharMatched(ctx, pc)
						continue
					}

					// try match a long flag
					pc.arg, pc.short, pc.dblTilde = pc.arg[2:], false, c1 == '~'
					if err = w.matchFlag(ctx, pc, false); !w.errIsSignalOrNil(err) {
						err = w.onUnknownFlagMatched(ctx, pc)
						break loopArgs
					}
					continue
				}
				if (c1 == '-' && pc.arg[1] == '~') || (c1 == '~' && pc.arg[1] == '-') {
					err = w.onUnknownFlagMatched(ctx, pc)
					break loopArgs
				}

				// try matching a short flag
				pc.arg, pc.short, pc.dblTilde = pc.arg[1:], true, false
				if c1 != '-' {
					err = w.onUnknownFlagMatched(ctx, pc)
					break loopArgs
				}
				if err = w.matchFlag(ctx, pc, true); !w.errIsSignalOrNil(err) {
					err = w.onUnknownFlagMatched(ctx, pc)
					break loopArgs
				}
				continue
			}
			// single '-' matched, is it a wrong state?
			err = w.onSingleHyphenMatched(ctx, pc)
			continue

		default:
			if pc.NoCandidateChildCommands() {
				pc.positionalArgs = append(pc.positionalArgs, pc.arg)
				logz.VerboseContext(ctx, "[cmdr] positional args added", "i", pc.i, "args", pc.positionalArgs)
				continue
			}
			if err = w.matchCommand(ctx, pc); !w.errIsSignalOrNil(err) {
				err = w.onUnknownCommandMatched(ctx, pc)
			}
		}
	}
	return
}

func (w *workerS) interpretLeadingPlusSign(pc *parseCtx) bool {
	if w.onInterpretLeadingPlusSign != nil {
		return w.onInterpretLeadingPlusSign(w, pc)
	}
	return false
}

func isCmdIsNil(cc cli.Cmd) (nilptr bool) {
	if x, ok := cc.(*cli.CmdS); ok {
		nilptr = x == nil
	} else {
		nilptr = cc == nil
	}
	return
}

func isCmdIsNotNil(cc cli.Cmd) (nilptr bool) {
	if x, ok := cc.(*cli.CmdS); ok {
		nilptr = x != nil
	} else {
		nilptr = cc != nil
	}
	return
}

func (w *workerS) matchCommand(ctx context.Context, pc *parseCtx) (err error) {
	err = cli.ErrUnmatchedCommand
	cmd := pc.LastCmd()
	if short, cc := cmd.Match(ctx, pc.arg); isCmdIsNotNil(cc) {
		ms, handled := pc.addCmd(cc, short), false
		handled, err = cc.TryOnMatched(0, ms)
		if err == nil {
			pc.lastCommand, err = len(pc.matchedCommands)-1, nil
		}
		logz.VerboseContext(ctx, "[cmdr] command matched", "short", short, "cmd", pc.LastCmd(), "handled", handled)
	}
	return
}

func (w *workerS) matchFlag(ctx context.Context, pc *parseCtx, short bool) (err error) {
	err = cli.ErrUnmatchedFlag
	cmd, vp := pc.LastCmd(), cli.NewFVP(w.args[pc.i+1:], pc.arg, short, pc.prefixPlusSign.Load(), pc.dblTilde)
	// defer func() { pc.i, vp.AteArgs = pc.i+vp.AteArgs, 0 }()

compactFlags:
	ff, err1 := cmd.MatchFlag(ctx, vp)
	if vp.Matched != "" && ff != nil && w.errIsSignalOrNil(err1) {
		ms, handled := pc.addFlag(ff), false
		handled, err1 = ff.TryOnMatched(0, ms)
		logz.VerboseContext(ctx, "[cmdr] flag matched", "short", vp.Short, "flg", ff, "val-pkg-val", ff.DefaultValue(), "handled", handled)

		pc.i += vp.AteArgs
		vp.AteArgs = 0
		err = err1

		if vp.Remains != "" && vp.PartialMatched {
			pc.arg = vp.Remains
			vp.Reset()
			pc.prefixPlusSign.Swap(false)

			if !errorsv3.Is(err, cli.ErrShouldStop) {
				goto compactFlags // try matching next compact flag. eg: '-avz' => '-a' parsed, remains '-vz'
			}
		}
	}
	return
}
