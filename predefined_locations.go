// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/conf"
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

		location = trimQuotes(location)

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

		if FileExists(fn) {
			err = w.rxxtOptions.LoadConfigFile(fn)
			if err != nil {
				return
			}
			conf.CfgFile = fn
			break
		}
	}
	return
}

// getExpandedPredefinedLocations for internal using
func (w *ExecWorker) getExpandedPredefinedLocations() (locations []string) {
	for _, d := range uniqueWorker.predefinedLocations {
		locations = append(locations, normalizeDir(d))
	}
	return
}
