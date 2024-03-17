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
		}
		if showHitStates {
			s.wrDebugScreen = &discardP{}
		} else {
			s.wrDebugScreen = os.Stdout
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

func withTasksBeforeRun(tasks ...cli.Task) cli.Opt { //nolint:unused
	return func(s *cli.Config) {
		s.TasksBeforeRun = tasks
	}
}
