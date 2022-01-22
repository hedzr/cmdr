/*
 * Copyright © 2019 Hedzr Yeh.
 */

// Package conf are used to store the app-level constants (app name/vaersion) for cmdr and your app.
// The names, such as app name, version, and building tags, would be held.
package conf

var (
	// AppName app name.
	// it'll be rewritten by build-arg.
	AppName string

	// Version app version.
	// it'll be rewritten by build-arg.
	Version = "0.2.1"
	// Buildstamp app built stamp.
	// it'll be rewritten by build-arg.
	Buildstamp = ""
	// Githash app git hash.
	// it'll be rewritten by build-arg.
	Githash = ""
	// GoVersion `go version` string.
	// it'll be rewritten by build-arg.
	GoVersion = ""

	// GitSummary holds the output of git describe --tags --dirty --always
	GitSummary = ""

	// GitShortVersion from `git describe --long`. [NEVER USED]
	GitShortVersion = ""
	// ServerTag app server tag names.[NEVER USED]
	ServerTag = ""
	// ServerID app server id.[NEVER USED]
	ServerID = ""
	// CfgFile never used [NEVER USED]
	CfgFile string
)
