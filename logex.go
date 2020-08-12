// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log"
	"github.com/hedzr/logex"
	"os"
	"strings"
)

// WithLogx enables github.com/hedzr/logex/logx integration
//
// Samplle:
//
//    WithLogx(log.NewDummyLogger()),	// import "github.com/hedzr/log"
//    WithLogx(log.NewStdLogger()),	// import "github.com/hedzr/log"
//    WithLogx(logrus.New()),		// import "github.com/hedzr/logex/logx/logrus"
//    WithLogx(sugar.New()),		// import "github.com/hedzr/logex/logx/zap/sugar"
//    WithLogx(zap.New()),			// import "github.com/hedzr/logex/logx/zap"
//
func WithLogx(logger log.Logger) ExecOption {
	return func(w *ExecWorker) {
		Logger = logger
	}
}

// Logger for cmdr
var Logger log.Logger = log.NewDummyLogger()

// WithLogex enables github.com/hedzr/logex integration
func WithLogex(lvl Level, opts ...logex.Option) ExecOption {
	return func(w *ExecWorker) {
		w.logexInitialFunctor = w.getWithLogexInitializer(lvl, opts...)
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

	logex.EnableWith(log.Level(l), opts...)
	// cmdr.Logger.Tracef("setup logger: lvl=%v", l)
	return
}

func (w *ExecWorker) getWithLogexInitializer(lvl Level, opts ...logex.Option) Handler {
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

// InDebugging return the status if cmdr was built with debug mode / or the app running under a debugger attached.
//
// To enable the debugger attached mode for cmdr, run `go build` with `-tags=delve` options. eg:
//
//     go run -tags=delve ./cli
//     go build -tags=delve -o my-app ./cli
//
// For Goland, you can enable this under 'Run/Debug Configurations', by adding the following into 'Go tool arguments:'
//
//     -tags=delve
//
// InDebugging() is a synonym to IsDebuggerAttached().
//
// NOTE that `isdelve` algor is from https://stackoverflow.com/questions/47879070/how-can-i-see-if-the-goland-debugger-is-running-in-the-program
//
//noinspection GoBoolExpressions
func InDebugging() bool {
	return log.InDebugging() // isdelve.Enabled
}

// IsDebuggerAttached return the status if cmdr was built with debug mode / or the app running under a debugger attached.
//
// To enable the debugger attached mode for cmdr, run `go build` with `-tags=delve` options. eg:
//
//     go run -tags=delve ./cli
//     go build -tags=delve -o my-app ./cli
//
// For Goland, you can enable this under 'Run/Debug Configurations', by adding the following into 'Go tool arguments:'
//
//     -tags=delve
//
// IsDebuggerAttached() is a synonym to InDebugging().
//
// NOTE that `isdelve` algor is from https://stackoverflow.com/questions/47879070/how-can-i-see-if-the-goland-debugger-is-running-in-the-program
//
//noinspection GoBoolExpressions
func IsDebuggerAttached() bool {
	return log.InDebugging() // isdelve.Enabled
}

// InTesting detects whether is running under go test mode
func InTesting() bool {
	if !strings.HasSuffix(tool.SavedOsArgs[0], ".test") &&
		!strings.Contains(tool.SavedOsArgs[0], "/T/___Test") {

		// [0] = /var/folders/td/2475l44j4n3dcjhqbmf3p5l40000gq/T/go-build328292371/b001/exe/main
		// !strings.Contains(SavedOsArgs[0], "/T/go-build")

		for _, s := range tool.SavedOsArgs {
			if s == "-test.v" || s == "-test.run" {
				return true
			}
		}
		return false

	}
	return true
}
