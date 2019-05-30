/*
 * Copyright © 2019 Hedzr Yeh.
 */

package demo

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/examples/demo/svr"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/sirupsen/logrus"
)

// Entry is app main entry
func Entry() {

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	// To disable internal commands and flags, uncomment the following codes
	// cmdr.EnableVersionCommands = false
	// cmdr.EnableVerboseCommands = false
	// cmdr.EnableCmdrCommands = false
	// cmdr.EnableHelpCommands = false
	// cmdr.EnableGenerateCommands = false

	daemon.Enable(svr.NewDaemon(), nil, nil)

	if err := cmdr.Exec(rootCmd); err != nil {
		logrus.Errorf("Error: %v", err)
	}

}
