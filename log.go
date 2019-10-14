// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"github.com/hedzr/logex"
	"github.com/sirupsen/logrus"
	"strings"
)

// WithLogex enables logex integration
func WithLogex() ExecOption {
	return func(w *ExecWorker) {
		w.withLogex = true
	}
}

func (w *ExecWorker) initWithLogex(cmd *Command, args []string) (err error) {
	logex.Enable()

	// var foreground = GetBoolR("server.foreground")
	var lvl = GetStringR("logger.level")
	var target = GetStringR("logger.target")
	var format = GetStringR("logger.format")

	if len(target) == 0 {
		target = "default"
	}
	if len(format) == 0 {
		format = "text"
	}
	if target == "journal" {
		format = "text"
	}
	switch format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		// logrus.SetFormatter(&formatter.TextFormatter{ForceColors: true})
	}

	// 	can_use_log_file, journal_mode := ij(target, foreground)
	l := stringToLevel(lvl)
	// 	if cli_common.Debug && l < logrus.DebugLevel {
	// 		l = logrus.DebugLevel
	// 	}
	logrus.SetLevel(l)
	logrus.Debugf("Using logger: format=%v, lvl=%v/%v, target=%v", format, lvl, l, target)

	return
}

// func earlierInitLogger() {
// 	l := "OFF"
// 	if !vxconf.IsProd() {
// 		l = "DEBUG"
// 	}
// 	l = vxconf.GetStringR("server.logger.level", l)
// 	logrus.SetLevel(stringToLevel(l))
// 	if l == "OFF" {
// 		logrus.SetOutput(ioutil.Discard)
// 	}
// }

func stringToLevel(s string) logrus.Level {
	s = strings.ToUpper(s)
	switch s {
	case "TRACE":
		return logrus.TraceLevel
	case "DEBUG", "devel", "dev":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "WARN", "WARNING":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "FATAL":
		return logrus.FatalLevel
	case "PANIC":
		return logrus.PanicLevel
	default:
		return logrus.WarnLevel
	}
}
