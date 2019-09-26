/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"os"
	"testing"
)

// Worker returns unexported worker for testing
func Worker() *ExecWorker {
	return uniqueWorker
}

func ResetWorker() {
	uniqueWorker = &ExecWorker{
		envPrefixes:  []string{"CMDR"},
		rxxtPrefixes: []string{"app"},

		predefinedLocations: []string{
			"./ci/etc/%s/%s.yml",
			"/etc/%s/%s.yml",
			"/usr/local/etc/%s/%s.yml",
			"$HOME/.%s/%s.yml",
			"$HOME/.config/%s/%s.yml",
		},

		shouldIgnoreWrongEnumValue: true,

		enableVersionCommands:  true,
		enableHelpCommands:     true,
		enableVerboseCommands:  true,
		enableCmdrCommands:     true,
		enableGenerateCommands: true,

		doNotLoadingConfigFiles: false,

		currentHelpPainter: new(helpPainter),

		defaultStdout: bufio.NewWriterSize(os.Stdout, 16384),
		defaultStderr: bufio.NewWriterSize(os.Stderr, 16384),

		rxxtOptions: NewOptions(),
	}
}

func TestTrapSignals(t *testing.T) {
	TrapSignals(func(s os.Signal) {
		//
	})
}
