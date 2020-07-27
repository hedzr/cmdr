// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"github.com/hedzr/logex"
	"os"
)

// WithLogex enables logex integration
func WithLogex(lvl Level, opts ...logex.Option) ExecOption {
	return func(w *ExecWorker) {
		w.logexInitialFunctor = w.getWithLogexInitializor(lvl, opts...)
	}
}

// WithLogexSkipFrames specify the skip frames to lookup the caller
func WithLogexSkipFrames(skipFrames int) ExecOption {
	return func(w *ExecWorker) {
		w.logexSkipFrames = skipFrames
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
	l := GetIntR("logger-level")
	return Level(l)
}

func (w *ExecWorker) processLevelStr(lvl Level, opts ...logex.Option) (err error) {
	var lvlStr = GetStringRP(w.logexPrefix, "level", lvl.String())
	var l Level

	l, err = ParseLevel(lvlStr)

	if l != OffLevel {
		if InDebugging() || GetDebugMode() {
			if l < DebugLevel {
				l = DebugLevel
			}
		}
		if GetBoolR("trace") || GetBool("trace") || ToBool(os.Getenv("TRACE")) {
			if l < TraceLevel {
				l = TraceLevel
				flog("--> processLevelStr: trace mode switched")
			}
		}
	}

	Set("logger-level", int(l))

	logex.EnableWith(logex.Level(l), opts...)
	// logrus.Tracef("setup logger: lvl=%v", l)
	return
}

func (w *ExecWorker) getWithLogexInitializor(lvl Level, opts ...logex.Option) Handler {
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
		logex.SetupLoggingFormat(format, w.logexSkipFrames)
		//switch format {
		//case "json":
		//	logrus.SetFormatter(&logrus.JSONFormatter{
		//		TimestampFormat:  "2006-01-02 15:04:05.000",
		//		DisableTimestamp: false,
		//		PrettyPrint:      false,
		//	})
		//default:
		//	e := false
		//	if w.logexSkipFrames > 0 {
		//		e = true
		//	}
		//	logrus.SetFormatter(&formatter.TextFormatter{
		//		ForceColors:               true,
		//		DisableColors:             false,
		//		FullTimestamp:             true,
		//		TimestampFormat:           "2006-01-02 15:04:05.000",
		//		Skip:                      w.logexSkipFrames,
		//		EnableSkip:                e,
		//		EnvironmentOverrideColors: true,
		//	})
		//}

		// can_use_log_file, journal_mode := ij(target, foreground)
		// l := GetLoggerLevel()
		// logrus.Tracef("Using logger: format=%v, lvl=%v, target=%v, formatter=%+v", format, l, target, logrus.StandardLogger().Formatter)

		return
	}
}
