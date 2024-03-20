package worker

import (
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/hedzr/cmdr/v2/cli"
)

type parseCtx struct {
	arg string
	i   int
	pos int

	root               *cli.RootCommand
	forceDefaultAction bool

	// parsed:

	short                 bool                             // parsing short flags?
	dblTilde              bool                             // parsing '~~' flag?
	lastCommand           int                              // index of ...
	matchedCommands       []*cli.Command                   // matched
	matchedCommandsStates map[*cli.Command]*cli.MatchState // matched ...
	matchedFlags          map[*cli.Flag]*cli.MatchState    // matched ...
	positionalArgs        []string                         //
	passThruMatched       int32                            // >0: index of '--'
	singleHyphenMatched   int32                            // >0: index of '-'
	prefixPlusSign        atomic.Bool                      // '+' leading
	// helpScreen            bool
}

func (s *parseCtx) addCmd(cc *cli.Command, short bool) (ms *cli.MatchState) {
	s.matchedCommands = append(s.matchedCommands, cc)
	if s.matchedCommandsStates == nil {
		s.matchedCommandsStates = make(map[*cli.Command]*cli.MatchState)
	}
	if st, ok := s.matchedCommandsStates[cc]; ok {
		st.HitStr, st.HitTimes = cc.GetHitStr(), cc.GetTriggeredTimes()
		ms = st
	} else {
		st = &cli.MatchState{
			Short:    short,
			HitStr:   cc.GetHitStr(),
			HitTimes: cc.GetTriggeredTimes(),
		}
		s.matchedCommandsStates[cc] = st
	}
	return
}

func (s *parseCtx) addFlag(ff *cli.Flag) (ms *cli.MatchState) {
	if s.matchedFlags == nil {
		s.matchedFlags = make(map[*cli.Flag]*cli.MatchState)
	}
	if st, ok := s.matchedFlags[ff]; ok {
		st.HitStr, st.HitTimes, st.Value, st.Short, st.DblTilde = ff.GetHitStr(), ff.GetTriggeredTimes(), ff.DefaultValue(), s.short, s.dblTilde
		ms = st
	} else {
		ms = &cli.MatchState{
			Short:    s.short,
			DblTilde: s.dblTilde,
			HitStr:   ff.GetHitStr(),
			HitTimes: ff.GetTriggeredTimes(),
			Value:    ff.DefaultValue(),
		}
		s.matchedFlags[ff] = ms
	}
	return
}

func (s *parseCtx) argsAre(list ...string) {
	if !reflect.DeepEqual(s.positionalArgs, list) {
		panic(fmt.Sprintf("expect positional args are %v but got %v (for cmd %v)", list, s.positionalArgs, s.LastCmd()))
	}
}

func (s *parseCtx) hitTest(longTitle string, times int) {
	cc := s.LastCmd()
	if f := cc.FindFlagBackwards(longTitle); f == nil {
		panic(fmt.Sprintf("can't found flag: %q", longTitle))
	} else if f.GetTriggeredTimes() != times {
		panic(fmt.Sprintf("expect hit times is %d but got %d (for flag %v)", times, f.GetTriggeredTimes(), f))
	}
}

func (s *parseCtx) valTest(longTitle string, val any) {
	cc := s.LastCmd()
	if f := cc.FindFlagBackwards(longTitle); f == nil {
		panic(fmt.Sprintf("can't found flag: %q", longTitle))
	} else if !reflect.DeepEqual(f.DefaultValue(), val) {
		panic(fmt.Sprintf("expect flag's value is '%v' but got '%v' (for flag %v)", val, f.DefaultValue(), f))
	}
}

func (s *parseCtx) flag(longTitle string) (f *cli.Flag) { //nolint:unused
	cc := s.LastCmd()
	f = cc.FindFlagBackwards(longTitle)
	return
}

func (s *parseCtx) cmd(longTitle string) (c *cli.Command) { //nolint:unused
	// ?? no uses yet ??
	c = s.root.FindSubCommand(longTitle, false)
	return
}

func (s *parseCtx) hasCmd(longTitle string, validator func(cc *cli.Command, state *cli.MatchState) bool) (found bool) { //nolint:revive,unused
	for _, c := range s.matchedCommands {
		if c.Long == longTitle {
			found = validator(c, s.matchedCommandsStates[c])
			break
		}
	}
	return
}

func (s *parseCtx) hasFlag(longTitle string, validator func(ff *cli.Flag, state *cli.MatchState) bool) (found bool) {
	for k, v := range s.matchedFlags {
		if k.Long == longTitle {
			found = validator(k, v)
			break
		}
	}
	return
}

func (s *parseCtx) NoCandidateChildCommands() bool {
	cmd := s.LastCmd()
	return len(cmd.SubCommands()) == 0
}

func (s *parseCtx) LastCmd() *cli.Command {
	cmd := s.root.Command
	if s.lastCommand >= 0 && len(s.matchedCommands) > 0 {
		cmd = s.matchedCommands[s.lastCommand]
	}
	return cmd
}

func (s *parseCtx) parsedCommandsStrings() (ret []string) { //nolint:revive,unused
	for _, cc := range s.matchedCommands {
		ret = append(ret, cc.String())
	}
	return
}
