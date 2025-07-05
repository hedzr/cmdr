package devmode

import (
	"fmt"
	"log"
	"os"

	"github.com/hedzr/cmdr/v2/pkg/logz"
	"github.com/hedzr/is"
	"github.com/hedzr/is/states"
	"github.com/hedzr/is/term"
	logzorig "github.com/hedzr/logg/slog"
)

func InDevelopmentMode() bool { return states.Env().InDevMode() }

func init() {
	log.SetFlags(log.LstdFlags | log.Lmsgprefix | log.LUTC | log.Lshortfile | log.Lmicroseconds)
	log.SetPrefix("")

	// d := dir.GetCurrentDir()
	// cliName := filepath.Dir(d)
	// prjName := filepath.Dir(cliName)
	// if filepath.Base(cliName) == "cli" {
	// 	d = prjName
	// }

	// isCmdrV2 := false
	// devModeFile := filepath.Join(d, ".dev-mode")
	// if devModeFilePresent = dir.FileExists(devModeFile); devModeFilePresent {
	// 	devMode = true
	// } else {
	// 	devModeFile := filepath.Join(d, ".devmode")
	// 	if devModeFilePresent = dir.FileExists(devModeFile); devModeFilePresent {
	// 		devMode = true
	// 	} else if dir.FileExists("go.mod") {
	// 		data, err := os.ReadFile("go.mod")
	// 		if err != nil {
	// 			return
	// 		}
	// 		content := string(data)

	// 		// dev := true
	// 		if strings.Contains(content, "github.com/hedzr/cmdr/v2/pkg/") {
	// 			devMode = false
	// 		}

	// 		// I am tiny-app in cmdr/v2, I will be launched in dev-mode always
	// 		if strings.Contains(content, "module github.com/hedzr/cmdr/v2") {
	// 			isCmdrV2, devMode = true, true
	// 		}
	// 	}
	// }

	boolenvvar := func(name string) (v, exists bool) {
		var val string
		if val, exists = os.LookupEnv(name); exists {
			v = is.StringToBool(val)
		}
		return
	}
	info := func(s string, args ...any) {
		if states.Env().IsVerboseMode() {
			logz.Skip(1).Print(s, args...)
		}
	}

	info(fmt.Sprintf("[dev-mode] the dev-mode is %v", is.InDevMode()))
	isCmdrV2, devModeFilePresent, devMode := false, false, false
	if states.DetectDevModeFileEnabled {
		isCmdrV2, devModeFilePresent, devMode = states.DetectDevModeFile()
	}

	// logzorig.SetLevel(logzorig.InfoLevel)
	// dbglog.Println("[dev-mode] initialize to InfoLevel", "dev-mode", devMode, "cwd", dir.GetCurrentDir())

	var debugMode, traceMode, verboseMode, debuggerAttached bool
	if devMode {
		if v, e := boolenvvar("CMDR_FORCE_DEBUG"); e && v {
			debugMode = true
			if devModeFilePresent {
				info("[dev-mode] .devmode file detected, CMDR_FORCE_DEBUG=1, entering Debug Mode...")
			} else {
				info("[dev-mode] dev-mode is enabled, CMDR_FORCE_DEBUG=1, entering Debug Mode...")
			}
		}
	}

	if is.DebugBuild() {
		debugMode = true
	}
	if is.VerboseBuild() {
		verboseMode = true
	}

	if v, e := boolenvvar("VERBOSE"); e {
		verboseMode = v
	}
	if v, e := boolenvvar("QUIET"); e {
		verboseMode = !v
	}
	if v, e := boolenvvar("DEBUG"); e {
		debugMode = v
	}
	if v, e := boolenvvar("TRACE"); e {
		traceMode = v
	}
	if isCmdrV2 {
		if v, e := boolenvvar("CMDR_NO_FORCE_DEBUG"); e {
			debugMode = !v
		}
		if v, e := boolenvvar("CMDR_FORCE_DEBUG"); e {
			debugMode = v
		}
	} else {
		if v, e := boolenvvar("NO_FORCE_DEBUG"); e {
			debugMode = !v
		}
		if v, e := boolenvvar("FORCE_DEBUG"); e {
			debugMode = v
		}
	}

	is.SetVerboseMode(verboseMode)
	is.SetDebugMode(debugMode)
	is.SetTraceMode(traceMode)

	if !debugMode && !devMode && !traceMode {
		logz.SetLevel(logzorig.WarnLevel)
		if verboseMode {
			info("[dev-mode] .set-level to warn for normal work")
		}

		if debugMode || debuggerAttached {
			logzorig.RemoveFlags(logzorig.Lprivacypathregexp, logzorig.Lprivacypath)
			if logz.GetLevel() < logzorig.DebugLevel {
				logz.SetLevel(logzorig.DebugLevel)
				info("[dev-mode] .set-level to debug", "debugMode", debugMode, "debuggerAttached", debuggerAttached)
			}
		} else {
			logzorig.AddFlags(logzorig.Lprivacypathregexp, logzorig.Lprivacypath)
		}

		if is.Windows() {
			if logz.GetLevel() < logzorig.InfoLevel {
				logz.SetLevel(logzorig.InfoLevel)
				info("[dev-mode] .set-level to info for windows")
			}
		}
	}
	// if traceMode {
	// 	if logz.GetLevel() < logzorig.TraceLevel {
	// 		logz.SetLevel(logzorig.TraceLevel)
	// 		info("[dev-mode] .set-level to trace")
	// 	}
	// }
	if verboseMode {
		if logz.GetLevel() < logzorig.InfoLevel {
			logz.SetLevel(logzorig.InfoLevel)
			info("[dev-mode] .set-level to info because VERBOSE=1 present")
		}
	}

	debuggerAttached = is.DebuggerAttached()
	debugMode = is.DebugMode()
	build, dBuild, vBuild := "", is.DebugBuild(), is.VerboseBuild()
	if dBuild {
		build += "debug('delve')"
	}
	if vBuild {
		build += "verbose('verbose')"
	}
	if build == "" {
		build = "none/normal"
	}
	logz.Debug(fmt.Sprintf(`[dev-mode] devenv states are:
			        logging-level: %v
			             dev-mode: %v
			dev-mode-file-present: %v
			           debug-mode: %v
			    debugger-attached: %v
			           trace-mode: %v
			  build (tags) method: %v
			`,
		logzorig.GetLevel(), devMode, devModeFilePresent,
		debugMode, debuggerAttached,
		traceMode, build,
	))
	// logz.Trace("trace mode enabled", "trace-mode", traceMode, "debug-level", logzorig.GetLevel())
	// logz.Print("debug level is", "debug-level", logzorig.GetLevel())

	n, r, p, t := term.StatStdout()
	if t {
		if term.IsColorful(os.Stdout) {
			is.SetNoColorMode(term.DisableColors())
		}
	} else if p || n || r {
		is.SetNoColorMode(true)
		info(fmt.Sprintf(`[dev-mode] .for %q, switch to no-color mode`, term.StatStdoutString()))
	}
}
