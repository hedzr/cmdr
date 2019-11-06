// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"github.com/hedzr/logex"
	"github.com/hedzr/logex/formatter"
	"github.com/sirupsen/logrus"
	"strings"
)

// WithLogex enables logex integration
func WithLogex(lvl logrus.Level, opts ...logex.LogexOption) ExecOption {
	return func(w *ExecWorker) {
		w.logexInitialFunctor = w.getWithLogexInitializor(lvl, opts...)
	}
}

// WithLogexPrefix specify a prefix string PS.
// 
// In cmdr options store, we will load the logging options under this key path:
//
//    app:
//      logger:
//        level:  DEBUG
//        format: text             # text, json, logfmt
//        target: default          # default, todo: journal
//        
// As showing above, the default prefix is "logger".
// You can replace it with yours, via WithLogexPrefix().
// For example, when you compose WithLogexPrefix("logging"), the following entries would be applied: 
//
//    app:
//      logging:
//        level:  DEBUG
//        format:
//        target:
//
func WithLogexPrefix(prefix string) ExecOption {
	return func(w *ExecWorker) {
		w.logexPrefix = prefix
	}
}

func (w *ExecWorker) getWithLogexInitializor(lvl logrus.Level, opts ...logex.LogexOption) func(cmd *Command, args []string) (err error) {
	return func(cmd *Command, args []string) (err error) {
		logex.EnableWith(lvl, opts...)

		if len(w.logexPrefix) == 0 {
			w.logexPrefix = "logger"
		}

		// var foreground = GetBoolR("server.foreground")
		var lvlStr = GetStringRP(w.logexPrefix, "level", lvl.String())
		var target = GetStringRP(w.logexPrefix, "target")
		var format = GetStringRP(w.logexPrefix, "format")

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
			logrus.SetFormatter(&formatter.TextFormatter{ForceColors: true})
		}

		// 	can_use_log_file, journal_mode := ij(target, foreground)
		l := stringToLevel(lvlStr)
		// 	if cli_common.Debug && l < logrus.DebugLevel {
		// 		l = logrus.DebugLevel
		// 	}
		logrus.SetLevel(l)
		logrus.Tracef("Using logger: format=%v, lvl=%v/%v, target=%v", format, lvlStr, l, target)

		return
	}
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
