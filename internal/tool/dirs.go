package tool

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/hedzr/cmdr/v2/pkg/dir"
	"github.com/hedzr/is"
	logz "github.com/hedzr/logg/slog"
)

// DataDir is used to store app normal runtime data.
//
// For darwin and linux it's generally at "$HOME/.local/share/<app>",
// or "/usr/local/share/<app>" and "/usr/share/<app>" in some builds.
//
// For windows it is "%APPDATA%/<app>/Data".
//
// In your application, it shall look up config files from ConfigDir,
// save the runtime data (or persistent data) into DataDir, use
// CacheDir to store the cache data which can be file and folder
// or file content indexes, the response replies from remote api,
// and so on.
// TempDir is used to store any temporary content which can be
// erased at any time.
//
// UsrLibDir is the place which an application should be installed
// at, in linux.
//
// VarRunDir is the place which a .pid, running socket file handle,
// and others files that can be shared in all processes of this
// application, sometimes for any apps.
func DataDir(appName string, base ...string) string {
	// appName := App().Name()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(append([]string{homeDir(), ".local", "share", appName}, base...)...)
		// return filepath.Join(homeDir(), "Library", "Application Supports", base)
	case "windows":
		for _, ev := range []string{"APPDATA", "CSIDL_APPDATA", "TEMP", "TMP"} {
			if v := os.Getenv(ev); v != "" {
				pre := filepath.Join(append([]string{v, appName}, base...)...)
				return filepath.Join(pre, "Data")
			}
		}
		// Worst case:
		return filepath.Join(append([]string{homeDir(), ".local", "share", appName}, base...)...)

	case "plan9":
		dir := os.Getenv("home")
		if dir == "" {
			logz.Error("[cmdr] $home is not defined")
			return ""
		}
		// ?
		return filepath.Join(append([]string{dir, "lib", "data", appName}, base...)...)
	}
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(append([]string{xdg, appName}, base...)...)
	}
	return filepath.Join(append([]string{homeDir(), ".local", "share", appName}, base...)...)
}

// ConfigDir returns the default root directory to use for user-specific
// configuration data. Users should create their own application-specific
// subdirectory within this one and use that.
//
// On Unix systems, it returns $XDG_CONFIG_HOME as specified by
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
// non-empty, else $HOME/.config/<appName>.
// On Darwin, it would be $HOME/.config/<appName>. ~~it returns $HOME/Library/Application Support.~~
// On Windows, it returns %AppData%/<appName>.
// On Plan 9, it returns $home/lib/<appName>.
//
// If the location cannot be determined (for example, $HOME is not defined),
// then it will return an error.
func ConfigDir(appName string, base ...string) string {
	// appName := App().Name()
	switch runtime.GOOS {
	case "darwin":
		t := filepath.Join(append([]string{homeDir(), ".config", appName}, base...)...)
		if dir.FileExists(t) {
			return t
		}
		r := filepath.Join(append([]string{homeDir(), "." + appName}, base...)...)
		if dir.FileExists(r) {
			return r
		}
		return t
		// return filepath.Join(homeDir(), "Library", "Application Supports", base)
	case "windows":
		for _, ev := range []string{"APPDATA", "CSIDL_APPDATA", "TEMP", "TMP"} {
			if v := os.Getenv(ev); v != "" {
				pre := filepath.Join(append([]string{v, appName}, base...)...)
				return filepath.Join(pre, "Config")
			}
		}
		// Worst case:
		return filepath.Join(append([]string{homeDir(), ".config", appName}, base...)...)

	case "plan9":
		dir := os.Getenv("home")
		if dir == "" {
			logz.Error("[cmdr] $home is not defined")
			return ""
		}
		return filepath.Join(append([]string{dir, "lib", appName}, base...)...)
	}

	// Unix
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(append([]string{xdg, appName}, base...)...)
	}
	return filepath.Join(append([]string{homeDir(), ".config", appName}, base...)...)
}

// CacheDir returns the default root directory to use for user-specific
// cached data. Users should create their own application-specific subdirectory
// within this one and use that.
//
// On Unix systems, it returns $XDG_CACHE_HOME as specified by
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
// non-empty, else $HOME/.cache.
// On Darwin, it would be $HOME/.cache/<appName>. ~~it returns $HOME/Library/Caches.~~
// On Windows, it returns %LocalAppData%/<appName>.
// On Plan 9, it returns $home/lib/cache/<appName>.
//
// If the location cannot be determined (for example, $HOME is not defined),
// then it will return an error.
func CacheDir(appName string, base ...string) string {
	// appName := App().Name()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(append([]string{homeDir(), ".cache", appName}, base...)...)
		// return filepath.Join(append([]string{homeDir(), "Library", "Caches", appName}, base...)...)
	case "windows":
		for _, ev := range []string{"APPDATA", "CSIDL_APPDATA", "TEMP", "TMP"} {
			if v := os.Getenv(ev); v != "" {
				return filepath.Join(append([]string{v, appName}, base...)...)
			}
		}
		// Worst case:
		return filepath.Join(append([]string{homeDir(), "." + appName}, base...)...)
	case "plan9":
		dir := os.Getenv("home")
		if dir == "" {
			logz.Error("[cmdr] $home is not defined")
			return ""
		}
		return filepath.Join(append([]string{dir, "lib", "cache", appName}, base...)...)
	}
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(append([]string{xdg, appName}, base...)...)
	}
	return filepath.Join(append([]string{homeDir(), ".cache", appName}, base...)...)
}

func HomeDir() string { return homeDir() }

func homeDir() string {
	home, _ := os.UserHomeDir()
	return home
	// if runtime.GOOS == "windows" {
	// 	return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	// }
	// if h := os.Getenv("HOME"); h != "" {
	// 	return h
	// }
	// return "/"
}

// TempDir returns the default directory to use for temporary files.
//
// On Unix systems, it returns $TMPDIR/<appName> if non-empty,
// else /tmp/<appName>.
// On Windows, it uses GetTempPath, returning the first non-empty
// value from %TMP%, %TEMP%, %USERPROFILE%, or the Windows directory.
// On Plan 9, it returns /tmp/<appName>.
//
// The directory is neither guaranteed to exist nor have accessible
// permissions.
func TempDir(appName string, base ...string) string {
	return filepath.Join(append([]string{os.TempDir(), appName}, base...)...)
}

func TempFileName(fileNamePattern, defaultFileName string, appName string, base ...string) (filename string) {
	tmpDir := TempDir(appName, base...)
	err := dir.EnsureDir(tmpDir)
	if err != nil {
		logz.Error("[cmdr] cannot creating tmpdir", "tmpdir", tmpDir, "err", err)
		return defaultFileName
	}

	f, err := os.CreateTemp(tmpDir, fileNamePattern)
	if err != nil {
		logz.Error("[cmdr] cannot create temporary file for flag", "err", err)
		return defaultFileName
	}
	filename = f.Name()
	_ = f.Close()
	return
}

// VarRunDir is todo, not exact right yet.
func VarRunDir(appName string, base ...string) string {
	// appName := App().Name()
	switch runtime.GOOS {
	case "darwin":
		// t := filepath.Join(append([]string{homeDir(), ".config", appName}, base...)...)
		// return filepath.Join(homeDir(), "Library", "Application Supports", base)
		t := filepath.Join(append([]string{"var", "run", appName}, base...)...)
		return t
	case "windows":
		return filepath.Join(append([]string{homeDir(), ".var", "run", appName}, base...)...)

	case "plan9":
		dir := os.Getenv("home")
		if dir == "" {
			logz.Error("[cmdr] $home is not defined")
			return ""
		}
		return filepath.Join(append([]string{dir, ".var", "run", appName}, base...)...)
	}

	// Unix
	return filepath.Join(append([]string{"var", "run", appName}, base...)...)
}

// UsrLibDir is todo, not exact right yet.
func UsrLibDir(appName string, base ...string) string {
	// appName := App().Name()
	switch runtime.GOOS {
	case "darwin":
		if is.Root() {
			return filepath.Join(append([]string{"usr", "lib", appName}, base...)...)
		}
		// t := filepath.Join(append([]string{homeDir(), ".config", appName}, base...)...)
		// return filepath.Join(homeDir(), "Library", "Application Supports", base)
		t := filepath.Join(append([]string{"usr", "local", "lib", appName}, base...)...)
		return t
	case "windows":
		return filepath.Join(append([]string{homeDir(), ".usr", "lib", appName}, base...)...)

	case "plan9":
		dir := os.Getenv("home")
		if dir == "" {
			logz.Error("[cmdr] $home is not defined")
			return ""
		}
		if is.Root() {
			return filepath.Join(append([]string{dir, "usr", "lib", appName}, base...)...)
		}
		return filepath.Join(append([]string{dir, "usr", "local", "lib", appName}, base...)...)
	}

	// Unix
	if is.Root() {
		return filepath.Join(append([]string{"usr", "lib", appName}, base...)...)
	}
	return filepath.Join(append([]string{"usr", "local", "lib", appName}, base...)...)
}
