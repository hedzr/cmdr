package worker

import (
	"context"
	"fmt"
	"sync/atomic"

	errorsv3 "gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

func (w *workerS) SetTasksAfterParse(tasks ...taskAfterParse) {
	w.tasksAfterParse = append(w.tasksAfterParse, tasks...)
}

func (w *workerS) parse(ctx context.Context, pc *parseCtx) (err error) {
	ec := errorsv3.New("tasks failed")

	defer func() {
		if len(w.tasksAfterParse) > 0 {
			for _, task := range w.tasksAfterParse {
				if task != nil {
					ec.Attach(task(ctx, w, pc, err))
				}
			}
		}

		ec.Defer(&err)

		if err == nil {
			w.parsingCtx = pc // save pc for later, OnAction might need it.
			// logz.VerboseContext(ctx, "pc.lastcmd", "lastcmd", pc.LastCmd())
			if cc, ok := pc.LastCmd().(*cli.CmdS); ok && cc != nil {
				logz.VerboseContext(ctx, "check binded positional-args receiver", "cc", cc)
				if varptr := cc.PositionalArgsPtr(); varptr != nil {
					logz.VerboseContext(ctx, fmt.Sprintf("checked varptr = %p", varptr), "cc", cc, "positional-args", pc.PositionalArgs())
					*varptr = pc.PositionalArgs()
				}
			}
		}
	}()

	logz.VerboseContext(ctx, "parsing command line args ...", "args", (*pc.argsPtr))

	if err = w.preApplyEnvMatched(ctx, pc); err != nil {
		return
	}

loopArgs:
	for pc.i = 1; pc.i < len(*pc.argsPtr); pc.i++ {
		if (*pc.argsPtr)[pc.i] == "" {
			continue
		}

		if atomic.LoadInt32(&pc.passThruMatched) > 0 || w.errShouldStopParsingLoop(err) {
			pc.positionalArgs = append(pc.positionalArgs, (*pc.argsPtr)[pc.i])
			logz.VerboseContext(ctx, "positional args added", "i", pc.i, "args", pc.positionalArgs)
			continue
		}

		logz.VerboseContext(ctx, "parsing command-line args", "i", pc.i, "arg", (*pc.argsPtr)[pc.i])

		pc.arg, pc.short, pc.dblTilde, pc.leadingPlus, pc.pos = (*pc.argsPtr)[pc.i], false, false, false, 0
		switch c1 := pc.arg[0]; c1 {
		// TODO need more design for form '+flag'.
		// currently, +flag is designed as a bool value flipper
		case '+': // for bool flag it's a flipper;
			if len(pc.arg) > 1 {
				pc.leadingPlus = true
				pc.arg, pc.short, pc.dblTilde = pc.arg[1:], true, false

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
					return w.matchFlag(ctx, pc, true)
				}); !w.errIsSignalOrNil(err) {
					if !pc.LastCmd().IgnoreUnmatched() {
						err = w.onUnknownFlagMatched(ctx, pc)
						break loopArgs
					}
				} else {
					continue
				}
			}
			// single '+': as a positional arg
			pc.positionalArgs = append(pc.positionalArgs, pc.arg)
			logz.VerboseContext(ctx, "positional args added", "i", pc.i, "args", pc.positionalArgs)
			continue

		case '-', '~':
			if len(pc.arg) > 1 { // so, pc.arg >= 2
				if (c1 == '-' && pc.arg[1] == '-') || (c1 == '~' && pc.arg[1] == '~') {
					if len(pc.arg) == 2 && pc.arg[1] == '-' {
						// --: pass-thru found
						err = w.onPassThruCharMatched(ctx, pc)
						continue
					}

					// try match a long flag
					pc.arg, pc.short, pc.dblTilde = pc.arg[2:], false, c1 == '~'
					if err = w.matchFlag(ctx, pc, false); !w.errIsSignalOrNil(err) {
						if !pc.LastCmd().IgnoreUnmatched() {
							err = w.onUnknownFlagMatched(ctx, pc)
							break loopArgs
						}
						pc.positionalArgs = append(pc.positionalArgs, (*pc.argsPtr)[pc.i])
						logz.VerboseContext(ctx, "positional args added", "i", pc.i, "args", pc.positionalArgs)
					}
					continue
				}
				if (c1 == '-' && pc.arg[1] == '~') || (c1 == '~' && pc.arg[1] == '-') {
					if !pc.LastCmd().IgnoreUnmatched() {
						err = w.onUnknownFlagMatched(ctx, pc)
						break loopArgs
					}
					pc.positionalArgs = append(pc.positionalArgs, (*pc.argsPtr)[pc.i])
					logz.VerboseContext(ctx, "positional args added", "i", pc.i, "args", pc.positionalArgs)
					continue
				}

				// try matching a short flag
				pc.arg, pc.short, pc.dblTilde = pc.arg[1:], true, false
				if c1 != '-' {
					if !pc.LastCmd().IgnoreUnmatched() {
						err = w.onUnknownFlagMatched(ctx, pc)
						break loopArgs
					}
				} else if err = w.matchFlag(ctx, pc, true); !w.errIsSignalOrNil(err) {
					if !pc.LastCmd().IgnoreUnmatched() {
						err = w.onUnknownFlagMatched(ctx, pc)
						break loopArgs
					}
				} else {
					continue
				}
				pc.positionalArgs = append(pc.positionalArgs, (*pc.argsPtr)[pc.i])
				logz.VerboseContext(ctx, "positional args added", "i", pc.i, "args", pc.positionalArgs)
				continue
			}
			// single '-' matched, is it a wrong state?
			err = w.onSingleHyphenMatched(ctx, pc)
			continue

		default: // for command
			if pc.NoCandidateChildCommands() || pc.LastCmd().IgnoreUnmatched() {
				pc.positionalArgs = append(pc.positionalArgs, pc.arg)
				logz.VerboseContext(ctx, "positional args added", "i", pc.i, "args", pc.positionalArgs)
				if pc.LastCmd().IgnoreUnmatched() {
					atomic.AddInt32(&pc.passThruMatched, 1)
					logz.VerboseContext(ctx, "entering passThruMode because a possible-cmd was added as positional-arg", "i", pc.i, "cmd", pc.LastCmd(), "arg", pc.arg)
					continue
				}
				continue
			}
			if err = w.matchCommand(ctx, pc); !w.errIsSignalOrNil(err) {
				if err = w.onUnknownCommandMatched(ctx, pc); w.errIsSignalFallback(err) {
					err = nil
					pc.positionalArgs = append(pc.positionalArgs, pc.arg)
					logz.VerboseContext(ctx, "positional args added", "i", pc.i, "args", pc.positionalArgs)
				}
			}
		}
	}
	return
}

func (w *workerS) interpretLeadingPlusSign(pc *parseCtx) bool {
	if w.OnInterpretLeadingPlusSign != nil {
		return w.OnInterpretLeadingPlusSign(w, pc)
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
		logz.VerboseContext(ctx, "command matched", "short", short, "cmd", pc.LastCmd(), "handled", handled)

		if pcl := cc.PresetCmdLines(); pcl != nil {
			a := *pc.argsPtr
			a = append(append(a[:pc.i], pcl...), a[pc.i:]...)
			*pc.argsPtr = a
		}
		if nm := cc.PassThruNow(); nm {
			// pass-thru now
			atomic.AddInt32(&pc.passThruMatched, 1)
			logz.VerboseContext(ctx, "entering passThruMode since cmd requested", "i", pc.i, "cmd", cc)
		}
	}
	return
}

func (w *workerS) matchFlag(ctx context.Context, pc *parseCtx, short bool) (err error) {
	err = cli.ErrUnmatchedFlag
	cmd, vp := pc.LastCmd(), cli.NewFVP((*pc.argsPtr)[pc.i+1:], pc.arg, short, pc.prefixPlusSign.Load(), pc.dblTilde)
	// defer func() { pc.i, vp.AteArgs = pc.i+vp.AteArgs, 0 }()

compactFlags:
	ff, err1 := cmd.MatchFlag(ctx, vp)
	if vp.Matched != "" && ff != nil && w.errIsSignalOrNil(err1) {
		ms, handled := pc.addFlag(ff), false
		handled, err1 = ff.TryOnMatched(0, ms)
		logz.VerboseContext(ctx, "flag matched", "short", vp.Short, "flg", ff, "val-pkg-val", ff.DefaultValue(), "handled", handled)

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

func (w *workerS) preApplyEnvMatched(ctx context.Context, pc *parseCtx) (err error) {
	if w.envvarMatched != nil {
		for ff, evm := range w.envvarMatched {
			if ff != nil {
				var handled bool
				ms := pc.addFlag(ff)
				handled, err = ff.TryOnMatched(0, ms)
				logz.DebugContext(ctx, "flag matched by envvar", "flg", ff, "envvar", evm.EnvVar, "value", evm.EnvValue)
				_ = handled
			}
		}
	}
	return
}
