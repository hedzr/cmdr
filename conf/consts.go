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

	// Version app semantic version. The result of `git describe --tags --abbrev=0`.
	//     Sample: '0.2.1' or 'v0.2.1'
	//     Sample: 'v0.3.23-9-g2239632a776f9d7a'
	// it'll be rewritten by build-arg.
	Version = "0.2.1"
	// Buildstamp app built stamp.
	//
	// The result of `date +'%Y-%m-%dT%H:%M:%S.%s+%Z'`, or `date +'%Y-%m-%dT%H:%M:%S+%Z`.
	// The recommend format is RFC3339, such as:
	//
	//     Sample: '2023-01-22T09:26:07+08:00'
	//
	// it'll be rewritten by build-arg.
	Buildstamp = ""
	// Githash app git hash. The result of `git rev-parse --short HEAD`.
	//     Sample: '2827a31'
	// it'll be rewritten by build-arg.
	Githash = ""
	// GoVersion `go version` string.
	// it'll be rewritten by build-arg.
	GoVersion = ""

	// GitSummary holds the output of `git describe --tags --dirty --always`
	//    Sample: 'v0.3.23-9-g2239632-dirty'
	// it'll be rewritten by build-arg.
	GitSummary = ""

	// GitDesc holds the output of `git log --oneline -1`
	//    Sample: '2239632 (HEAD -> master) improved `sbom` command description line.'
	GitDesc = ""

	// BuilderComments can be rewitten by build-arg
	BuilderComments = ""

	// GitShortVersion from `git describe --long`. [NEVER USED]
	GitShortVersion = ""
	// ServerTag app server tag names.[NEVER USED]
	ServerTag = ""
	// ServerID app server id.[NEVER USED]
	ServerID = ""
	// CfgFile never used [NEVER USED]
	CfgFile string

	// Serial is a serial number (int64) from build-tool.
	// `bgo` can hold and manage a serial number in her runtime environment.
	Serial string
	// SerialString is a random string from build-tool.
	// `bgo can hold and manage a random string in her runtime environment.
	SerialString string
)
