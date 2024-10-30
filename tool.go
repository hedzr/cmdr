// Copyright © 2022 Hedzr Yeh.

package cmdr

import (
	"github.com/hedzr/cmdr/v2/internal/tool"
)

// DataDir returns standard datadir associated with this app.
// In general, it would be /usr/local/share/<appName>,
// $HOME/.local/share/<appName>, and so on.
//
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
func DataDir(base ...string) string {
	return tool.DataDir(AppName(), base...)
}

// ConfigDir returns standard configdir associated with this app.
// In general, it would be /etc/<appName>,
// $HOME/.config/<appName>, and so on.
//
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
func ConfigDir(base ...string) string {
	return tool.ConfigDir(AppName(), base...)
}

// CacheDir returns standard cachedir associated with this app.
// In general, it would be /var/cache/<appName>,
// $HOME/.cache/<appName>, and so on.
//
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
func CacheDir(base ...string) string {
	return tool.CacheDir(AppName(), base...)
}

// HomeDir returns the current user's home directory.
// In general, it would be /Users/<username>, /home/<username>, etc.
func HomeDir() string { return tool.HomeDir() }

func TempDir(base ...string) string {
	return tool.TempDir(AppName(), base...)
}

func TempFileName(fileNamePattern, defaultFileName string, base ...string) (filename string) {
	return tool.TempFileName(fileNamePattern, defaultFileName, AppName(), base...)
}

// VarLogDir is todo, not exact right yet.
func VarLogDir(base ...string) string { return tool.VarLogDir(AppName(), base...) }

// VarRunDir is the runtime temp dir. "/var/run/<app>/"
func VarRunDir(base ...string) string { return tool.VarRunDir(AppName(), base...) }

// UsrLibDir is the runtime temp dir. "/usr/local/lib/<app>/" or "/usr/lib/<app>" in root mode.
func UsrLibDir(base ...string) string { return tool.UsrLibDir(AppName(), base...) }
