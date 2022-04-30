// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log/dir"
	"os"
	"path"
	"strings"
)

func (w *ExecWorker) parsePredefinedLocation() (err error) {
	// pre-detects for `--config xxx`, `--config=xxx`, `--configxxx`
	if ix, str, yes := partialContains(os.Args, "--config"); yes {
		var location string
		i := strings.Index(str, "=")
		switch {
		case i > 0:
			location = str[i+1:]
		case len(str) > 8:
			location = str[8:]
		case ix+1 < len(os.Args):
			location = os.Args[ix+1]
		}

		location = tool.StripQuotes(location)
		flog("--> preprocess / buildXref / parsePredefinedLocation: %q", location)

		if location != "" && dir.FileExists(location) {
			if yes, err = dir.IsDirectory(location); yes {
				if dir.FileExists(path.Join(location, w.confDFolderName)) {
					setPredefinedLocations(location + "/%s.yml")
				} else {
					setPredefinedLocations(location + "/%s/%s.yml")
				}
			} else if yes, err = dir.IsRegularFile(location); yes {
				setPredefinedLocations(location)
			}
		}
	}
	return
}

func (w *ExecWorker) checkMoreLocations(rootCmd *RootCommand) (err error) {
	if w.watchChildConfigFiles {
		a1, a2 := ".$APPNAME.yml", ".$APPNAME/*.yml"
		a3, a4 := os.ExpandEnv(a1), os.ExpandEnv(a2)
		b := dir.FileExists(a3)
		if b {
			w.predefinedLocations = append(w.predefinedLocations, a3)
		}
		b = dir.FileExists(a4) //nolint:staticcheck //keep it
		if b {                 //nolint:staticcheck //keep it
			//
		}
	}
	return
}

func (w *ExecWorker) loadFromPredefinedLocations(rootCmd *RootCommand) (err error) {
	_ = w.checkMoreLocations(w.rootCommand)

	var mainFile, subDir string
	mainFile, subDir, err = w.loadFromLocations(rootCmd, w.getExpandedPredefinedLocations(), mainConfigFiles)
	if err == nil {
		conf.CfgFile = mainFile
		flog("--> preprocess / buildXref / loadFromPredefinedLocations: %q loaded (CFG_DIR=%v)", mainFile, subDir)
		// flog("--> loadFromPredefinedLocations(): %q loaded", fn)
	}
	return
}

func (w *ExecWorker) loadFromSecondaryLocations(rootCmd *RootCommand) (err error) {
	var mainFile, subDir string
	mainFile, subDir, err = w.loadFromLocations(rootCmd, w.getExpandedSecondaryLocations(), secondaryConfigFiles)
	if err == nil {
		// conf.CfgFile = mainFile
		flog("--> preprocess / buildXref / loadFromSecondaryLocations: %q loaded (CFG_DIR_2NDRY=%v)", mainFile, subDir)
		// flog("--> loadFromPredefinedLocations(): %q loaded", fn)
	}
	return
}

func (w *ExecWorker) loadFromAlterLocations(rootCmd *RootCommand) (err error) {
	var mainFile, subDir string
	mainFile, subDir, err = w.loadFromLocations(rootCmd, w.getExpandedAlterLocations(), alterConfigFile)
	if err == nil {
		flog("--> preprocess / buildXref / loadFromAlterLocations: %q loaded (ALTER_DIR=%v)", mainFile, subDir)
	}
	return
}

func (w *ExecWorker) loadFromLocations(rootCmd *RootCommand, locations []string, cft configFileType) (mainFile, subDir string, err error) {
	// and now, loading the external configuration files
	for _, s := range locations {
		fn := s
		switch strings.Count(fn, "%s") {
		case 2:
			fn = fmt.Sprintf(s, rootCmd.AppName, rootCmd.AppName)
		case 1:
			fn = fmt.Sprintf(s, rootCmd.AppName)
		}

		b := dir.FileExists(fn)
		if !b {
			fn = replaceAll(fn, ".yml", ".yaml")
			b = dir.FileExists(fn)
		}
		if b {
			mainFile, subDir, err = w.rxxtOptions.LoadConfigFile(fn, cft)
			break
		}
	}
	return
}

// getExpandedAlterLocations for internal using
func (w *ExecWorker) getExpandedAlterLocations() (locations []string) {
	for _, d := range internalGetWorker().alterLocations {
		locations = uniAddStr(locations, dir.NormalizeDir(d))
	}
	return
}

// getExpandedSecondaryLocations for internal using
func (w *ExecWorker) getExpandedSecondaryLocations() (locations []string) {
	for _, d := range internalGetWorker().secondaryLocations {
		locations = uniAddStr(locations, dir.NormalizeDir(d))
	}
	return
}

// getExpandedPredefinedLocations for internal using
func (w *ExecWorker) getExpandedPredefinedLocations() (locations []string) {
	for _, d := range internalGetWorker().predefinedLocations {
		locations = uniAddStr(locations, dir.NormalizeDir(d))
	}
	return
}

// GetPredefinedLocations return the primary searching locations for
// loading the main config files.
// cmdr finds these location to create the main config store.
func GetPredefinedLocations() []string {
	return internalGetWorker().predefinedLocations
}

// GetSecondaryLocations return the secondary searching
// locations, and these configs will be merged into main config
// store.
func GetSecondaryLocations() []string {
	return internalGetWorker().secondaryLocations
}

// GetPredefinedAlterLocations return the alternative searching
// locations.
// The alter config file will be merged into main config store
// after secondary config merged.
// The most different things are the alter config file can be
// written back when cmdr
func GetPredefinedAlterLocations() []string {
	return internalGetWorker().alterLocations
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

func setSecondaryLocations(locations ...string) {
	internalGetWorker().secondaryLocations = locations
}

func setAlterLocations(locations ...string) {
	internalGetWorker().alterLocations = locations
}
