package devmode

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hedzr/cmdr/v2/pkg/dir"
	"github.com/hedzr/cmdr/v2/pkg/logz"
	"github.com/hedzr/is"
	logzorig "github.com/hedzr/logg/slog"
)

func InDevelopmentMode() bool { return devMode }

var onceDev sync.Once
var devMode bool

func init() {
	// onceDev is a redundant operation, but we still keep it to
	// fit for defensive programming style.
	onceDev.Do(func() {
		log.SetFlags(log.LstdFlags | log.Lmsgprefix | log.LUTC | log.Lshortfile | log.Lmicroseconds)
		log.SetPrefix("")

		d := dir.GetCurrentDir()
		cliName := filepath.Dir(d)
		prjName := filepath.Dir(cliName)
		if filepath.Base(cliName) == "cli" {
			d = prjName
		}

		devModeFile := filepath.Join(d, ".dev-mode")
		if dir.FileExists(devModeFile) {
			devMode = true
		} else if dir.FileExists("go.mod") {
			data, err := os.ReadFile("go.mod")
			if err != nil {
				return
			}
			content := string(data)

			// dev := true
			if strings.Contains(content, "github.com/hedzr/cmdr/v2/pkg/") {
				devMode = false
			}

			// I am tiny-app in cmdr/v2, I will be launched in dev-mode always
			if strings.Contains(content, "module github.com/hedzr/cmdr/v2") {
				devMode = true
			}
		}

		logzorig.SetLevel(logzorig.InfoLevel)

		// dbglog.Println("[dev-mode] initialize to InfoLevel", "dev-mode", devMode, "cwd", dir.GetCurrentDir())

		if devMode {
			is.SetDebugMode(true)
			logz.SetLevel(logzorig.DebugLevel)
			logz.Debug(".dev-mode file detected, entering Debug Mode...")
		}

		if is.DebugBuild() {
			is.SetDebugMode(true)
			logz.SetLevel(logzorig.DebugLevel)
		}

		if is.VerboseBuild() {
			is.SetVerboseMode(true)
			if logz.GetLevel() < logzorig.InfoLevel {
				logz.SetLevel(logzorig.InfoLevel)
				logz.Debug(".set-level to info")
			}
			if logz.GetLevel() < logzorig.TraceLevel {
				logz.SetLevel(logzorig.TraceLevel)
				logz.Debug(".set-level to trace")
			}
		}
		logz.Debug(".logging-level is", "level", logzorig.GetLevel())

		if is.DebugMode() || is.DebuggerAttached() {
			logzorig.RemoveFlags(logzorig.Lprivacypathregexp, logzorig.Lprivacypath)
			if logz.GetLevel() < logzorig.DebugLevel {
				logz.SetLevel(logzorig.DebugLevel)
				logz.Debug(".set-level to debug")
			}
		} else {
			logzorig.AddFlags(logzorig.Lprivacypathregexp, logzorig.Lprivacypath)
		}
		if is.Windows() {
			if logz.GetLevel() < logzorig.InfoLevel {
				logz.SetLevel(logzorig.InfoLevel)
				logz.Debug(".set-level to info")
			}
		}
	})
}
