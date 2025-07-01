package worker

import (
	"context"
	"strings"
	"sync/atomic"
	"text/template"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

type parseCtx struct {
	arg string
	i   int
	pos int

	argsPtr *[]string // the parsing command-line args

	root               *cli.RootCommand
	forceDefaultAction bool

	// parsed:

	short                 bool                          // parsing short flags? ('-' or '--')
	dblTilde              bool                          // parsing '~~' flag?
	leadingPlus           bool                          // parsing '+' flag?
	lastCommand           int                           // index of ...
	matchedCommands       []cli.Cmd                     // matched
	matchedCommandsStates map[cli.Cmd]*cli.MatchState   // matched ...
	matchedFlags          map[*cli.Flag]*cli.MatchState // matched ...
	positionalArgs        []string                      //
	passThruMatched       int32                         // >0: index of '--'
	singleHyphenMatched   int32                         // >0: index of '-'
	prefixPlusSign        atomic.Bool                   // '+' leading

	// helpScreen            bool
}

func (s *parseCtx) LeadingIs(r rune) bool {
	switch r {
	case '+':
		return s.leadingPlus
	case '~':
		return s.dblTilde
	case '-':
		return !s.leadingPlus && !s.dblTilde
	case '/':
		return !s.leadingPlus && !s.dblTilde
	default:
		return false
	}
}

func (s *parseCtx) Translate(pattern string) (result string) {
	if tpl, err := template.New("cool").Parse(pattern); err != nil {
		logz.Warn("given pattern cannot be transalted or expanded", "pattern", pattern, "err", err)
		return
	} else {
		var sb strings.Builder
		cmd := s.LastCmd()
		if err = tpl.Execute(&sb, struct {
			AppName     string
			AppVersion  string
			DadCommands string // for curr-cmd is `jump to`: "jump"
			Commands    string // for curr-cmd is `jump to`: "jump to"
			*parseCtx
		}{
			cmd.App().Name(),
			cmd.App().Version(),
			s.DadCommandsText(),
			s.CommandsText(),
			s,
		}); err != nil {
			logz.Warn("given pattern cannot be transalted", "pattern", pattern, "err", err)
			return
		}
		result = sb.String()
	}
	return
}

func (s *parseCtx) DadCommandsText() (result string) {
	if s != nil && len(s.matchedCommands) > 1 {
		var ss []string
		for _, z := range s.matchedCommands[:len(s.matchedCommands)-1] {
			ss = append(ss, z.Name())
		}
		result = strings.Join(ss, " ")
	}
	return
}

func (s *parseCtx) CommandsText() (result string) {
	if s != nil && len(s.matchedCommands) > 0 {
		var ss []string
		for _, z := range s.matchedCommands {
			ss = append(ss, z.Name())
		}
		result = strings.Join(ss, " ")
	}
	return
}

func (s *parseCtx) CommandMatchedState(c cli.Cmd) (ms *cli.MatchState) {
	if s != nil {
		if m, ok := s.matchedCommandsStates[c]; ok {
			return m
		}
	}
	return nil
}

func (s *parseCtx) FlagMatchedState(f *cli.Flag) (ms *cli.MatchState) {
	if s != nil {
		if m, ok := s.matchedFlags[f]; ok {
			return m
		}
	}
	return nil
}

func (s *parseCtx) matchedCommand(longTitle string) (cc cli.Cmd) {
	for _, cc = range s.matchedCommands {
		if cc.Name() == longTitle {
			return cc
		}
	}
	return nil
}

func (s *parseCtx) matchedFlag(ctx context.Context, longTitle string) (ff *cli.Flag) {
	ff = s.flag(ctx, longTitle)
	if _, ok := s.matchedFlags[ff]; ok {
		return ff
	}
	return nil
}

func (s *parseCtx) addCmd(cc cli.Cmd, short bool) (ms *cli.MatchState) {
	if cc == nil {
		logz.Panic("the adding command shouldn't be nil")
		panic("")
	}
	s.matchedCommands = append(s.matchedCommands, cc)
	if s.matchedCommandsStates == nil {
		s.matchedCommandsStates = make(map[cli.Cmd]*cli.MatchState)
	}
	if st, ok := s.matchedCommandsStates[cc]; ok {
		st.HitStr, st.HitTimes = cc.HitTitle(), cc.HitTimes()
		ms = st
	} else {
		st = &cli.MatchState{
			Short:    short,
			HitStr:   cc.HitTitle(),
			HitTimes: cc.HitTimes(),
		}
		s.matchedCommandsStates[cc] = st
	}
	return
}

func (s *parseCtx) addFlag(ff *cli.Flag) (ms *cli.MatchState) {
	if ff == nil {
		logz.Panic("the adding flag shouldn't be nil")
		panic("")
	}
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
			Plus:     s.leadingPlus,
			HitStr:   ff.GetHitStr(),
			HitTimes: ff.GetTriggeredTimes(),
			Value:    ff.DefaultValue(),
		}
		s.matchedFlags[ff] = ms
	}

	// save the matched value into option store
	cmdstore := ff.Store() // get the owner's store
	title := ff.Name()
	val := ff.DefaultValue()
	_, _ = cmdstore.Set(title, val)
	if varptr := ff.VarPtr(); varptr != nil {
		ff.WriteBoundValue(val)
	}
	return
}

func (s *parseCtx) flag(ctx context.Context, longTitle string) (f *cli.Flag) { //nolint:unused
	cc := s.LastCmd()
	f = cc.FindFlagBackwards(ctx, longTitle)
	return
}

func (s *parseCtx) cmd(ctx context.Context, longTitle string) (c cli.Cmd) { //nolint:unused
	// ?? no uses yet ??
	ret := s.root.FindSubCommand(ctx, longTitle, false)
	if rc, ok := ret.(*cli.CmdS); ok {
		c = rc
	}
	return
}

func (s *parseCtx) HasCmd(longTitle string, validator func(cc cli.Cmd, state *cli.MatchState) bool) (found bool) { //nolint:revive,unused
	if s == nil {
		return false
	}
	for _, c := range s.matchedCommands {
		if c.Name() == longTitle {
			found = validator(c, s.matchedCommandsStates[c])
			break
		}
	}
	return
}

func (s *parseCtx) HasFlag(longTitle string, validator func(ff *cli.Flag, state *cli.MatchState) bool) (found bool) {
	if s == nil {
		return false
	}
	for k, v := range s.matchedFlags {
		if k.Long == longTitle {
			found = validator(k, v)
			break
		}
	}
	return
}

func (s *parseCtx) NoCandidateChildCommands() bool {
	if s == nil {
		return false
	}
	cmd := s.LastCmd()
	return len(cmd.SubCommands()) == 0
}

func (s *parseCtx) LastCmd() cli.Cmd {
	var cmd = s.root.Cmd
	if s != nil {
		if s.lastCommand >= 0 && len(s.matchedCommands) > 0 {
			cmd = s.matchedCommands[s.lastCommand]
		}
	}
	return cmd
}

func (s *parseCtx) PositionalArgs() []string {
	if s != nil {
		return s.positionalArgs
	}
	return nil
}

func (s *parseCtx) MatchedCommands() []cli.Cmd {
	if s != nil {
		return s.matchedCommands
	}
	return nil
}

func (s *parseCtx) MatchedFlags() map[*cli.Flag]*cli.MatchState {
	if s != nil {
		return s.matchedFlags
	}
	return nil
}

func (s *parseCtx) parsedCommandsStrings() (ret []string) { //nolint:revive,unused
	for _, cc := range s.matchedCommands {
		ret = append(ret, cc.String())
	}
	return
}
