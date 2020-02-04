// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"github.com/hedzr/logex"
	"github.com/hedzr/logex/formatter"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

// WithLogex enables logex integration
func WithLogex(lvl Level, opts ...logex.LogexOption) ExecOption {
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
//        level:  DEBUG            # panic, fatal, error, warn, info, debug, trace, off
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

// GetLoggerLevel returns the current logger level after parsed.
func GetLoggerLevel() Level {
	l:=GetIntR("logger-level")
	return Level(l)
}

func (w *ExecWorker) processLevelStr(lvl Level, opts ...logex.LogexOption) (err error) {
	var lvlStr = GetStringRP(w.logexPrefix, "level", lvl.String())
	var l Level

	l, err = ParseLevel(lvlStr)

	if InDebugging() || GetDebugMode() {
		if l < DebugLevel {
			l = DebugLevel
		}
	}
	if GetBoolR("trace") || GetBool("trace") || toBool(os.Getenv("TRACE")) {
		if l < TraceLevel {
			l = TraceLevel
		}
	}

	Set("logger-level", int(l))

	if l == OffLevel {
		logex.EnableWith(logrus.ErrorLevel, opts...)
		logrus.SetOutput(ioutil.Discard)
	} else {
		logex.EnableWith(logrus.Level(l), opts...)
	}
	logrus.Tracef("setup logger: lvl=%v", l)
	return
}

func (w *ExecWorker) getWithLogexInitializor(lvl Level, opts ...logex.LogexOption) func(cmd *Command, args []string) (err error) {
	return func(cmd *Command, args []string) (err error) {

		if len(w.logexPrefix) == 0 {
			w.logexPrefix = "logger"
		}

		err = w.processLevelStr(lvl, opts...)

		// var foreground = GetBoolR("server.foreground")
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
			logrus.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat:  "2006-01-02 15:04:05.000",
				DisableTimestamp: false,
				PrettyPrint:      false,
			})
		default:
			logrus.SetFormatter(&formatter.TextFormatter{
				ForceColors:     true,
				DisableColors:   false,
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05.000",
			})
		}

		// can_use_log_file, journal_mode := ij(target, foreground)
		l := GetLoggerLevel()
		logrus.Tracef("Using logger: format=%v, lvl=%v, target=%v", format, l, target)

		return
	}
}
