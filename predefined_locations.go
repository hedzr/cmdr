// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/cmdr/tool"
	"os"
	"strings"
)

func (w *ExecWorker) parsePredefinedLocation() (err error) {
	// pre-detects for `--config xxx`, `--config=xxx`, `--configxxx`
	if ix, str, yes := partialContains(os.Args, "--config"); yes {
		var location string
		if i := strings.Index(str, "="); i > 0 {
			location = str[i+1:]
		} else if len(str) > 8 {
			location = str[8:]
		} else if ix+1 < len(os.Args) {
			location = os.Args[ix+1]
		}

		location = tool.StripQuotes(location)
		flog("--> preprocess / buildXref / parsePredefinedLocation: %q", location)

		if len(location) > 0 && FileExists(location) {
			if yes, err = IsDirectory(location); yes {
				if FileExists(location + "/conf.d") {
					setPredefinedLocations(location + "/%s.yml")
				} else {
					setPredefinedLocations(location + "/%s/%s.yml")
				}
			} else if yes, err = IsRegularFile(location); yes {
				setPredefinedLocations(location)
			}
		}
	}
	return
}

func (w *ExecWorker) loadFromPredefinedLocation(rootCmd *RootCommand) (err error) {
	// and now, loading the external configuration files
	for _, s := range w.getExpandedPredefinedLocations() {
		fn := s
		switch strings.Count(fn, "%s") {
		case 2:
			fn = fmt.Sprintf(s, rootCmd.AppName, rootCmd.AppName)
		case 1:
			fn = fmt.Sprintf(s, rootCmd.AppName)
		}

		b := FileExists(fn)
		if !b {
			fn = replaceAll(fn, ".yml", ".yaml")
			b = FileExists(fn)
		}
		if b {
			err = w.rxxtOptions.LoadConfigFile(fn)
			if err == nil {
				conf.CfgFile = fn
				flog("--> preprocess / buildXref / loadFromPredefinedLocation: %q loaded", fn)
				//flog("--> loadFromPredefinedLocation(): %q loaded", fn)
			}
			break
		}
	}
	return
}

// getExpandedPredefinedLocations for internal using
func (w *ExecWorker) getExpandedPredefinedLocations() (locations []string) {
	for _, d := range internalGetWorker().predefinedLocations {
		locations = uniAddStr(locations, normalizeDir(d))
	}
	return
}

// GetPredefinedLocations return the searching locations for loading config files.
func GetPredefinedLocations() []string {
	return internalGetWorker().predefinedLocations
}

// // SetPredefinedLocations to customize the searching locations for loading config files.
// //
// // It MUST be invoked before `cmdr.Exec`. Such as:
// // ```go
// //     SetPredefinedLocations([]string{"./config", "~/.config/cmdr/", "$GOPATH/running-configs/cmdr"})
// // ```
// func SetPredefinedLocations(locations []string) {
// 	uniqueWorker.predefinedLocations = locations
// }

func setPredefinedLocations(locations ...string) {
	internalGetWorker().predefinedLocations = locations
}
