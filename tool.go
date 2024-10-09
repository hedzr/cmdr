// Copyright © 2022 Hedzr Yeh.

package cmdr

import (
	"github.com/hedzr/cmdr/v2/internal/tool"
)

// DataDir returns standard datadir associated with this app.
// In general, it would be /usr/local/share/<appName>,
// $HOME/.local/share/<appName>, and so on.
func DataDir(base ...string) string {
	return tool.DataDir(AppName(), base...)
}

// ConfigDir returns standard configdir associated with this app.
// In general, it would be /etc/<appName>,
// $HOME/.config/<appName>, and so on.
func ConfigDir(base ...string) string {
	return tool.ConfigDir(AppName(), base...)
}

// CacheDir returns standard cachedir associated with this app.
// In general, it would be /var/cache/<appName>,
// $HOME/.cache/<appName>, and so on.
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
