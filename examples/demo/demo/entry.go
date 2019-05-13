/*
 * Copyright © 2019 Hedzr Yeh.
 */

package demo

import (
	"github.com/hedzr/cmdr"
	"github.com/sirupsen/logrus"
)

func Entry() {

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true,})

	cmdr.EnableVersionCommands = true
	cmdr.EnableVerboseCommands = true
	cmdr.EnableHelpCommands = true
	cmdr.EnableGenerateCommands = true
	if err := cmdr.Exec(rootCmd); err != nil {
		logrus.Errorf("Error: %v", err)
	}

}
