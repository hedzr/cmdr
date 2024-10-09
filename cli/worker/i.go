package worker

import (
	"os"

	"github.com/hedzr/cmdr/v2/cli"
)

type wOpt func(s *workerS)

func WithConfig(c *cli.Config) wOpt {
	return func(s *workerS) {
		s.Config = c
	}
}

func WithHelpScreenSets(showHelpScreen, showHitStates bool) wOpt { //nolint:revive
	return func(s *workerS) {
		if showHelpScreen {
			s.wrHelpScreen = &discardP{}
		} else {
			s.wrHelpScreen = os.Stdout
			if s.HelpScreenWriter != nil {
				s.wrHelpScreen = s.HelpScreenWriter
			}
		}
		if showHitStates {
			s.wrDebugScreen = &discardP{}
		} else {
			s.wrDebugScreen = os.Stdout
			if s.DebugScreenWriter != nil {
				s.wrDebugScreen = s.DebugScreenWriter
			}
		}
	}
}

//

//

//

func withTasksBeforeParse(tasks ...cli.Task) cli.Opt {
	return func(s *cli.Config) {
		s.TasksBeforeParse = tasks
	}
}

// func withTasksAfterParse(tasks ...taskAfterParse) wOpt {
// 	return func(s *workerS) {
// 		s.tasksAfterParse = tasks
// 	}
// }

func withEnv(env map[string]string) cli.Opt {
	return func(s *cli.Config) {
		if s.Env == nil {
			s.Env = make(map[string]string)
		}
		for k, v := range env {
			s.Env[k] = v
		}
	}
}

func withTasksBeforeRun(tasks ...cli.Task) cli.Opt { //nolint:unused
	return func(s *cli.Config) {
		s.TasksBeforeRun = tasks
	}
}

func withHelpScreenWriter(w HelpWriter) cli.Opt {
	return func(s *cli.Config) {
		s.HelpScreenWriter = w
	}
}
