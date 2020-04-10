/*
 * Copyright © 2019 Hedzr Yeh.
 */

package conf

var (
	// CfgFile never used
	CfgFile string
	// AppName app name
	AppName string

	// these 3 variables will be rewrote when app had been building by ci-tool

	// Version app version
	Version = "0.2.1"
	// Buildstamp app built stamp
	Buildstamp = ""
	// Githash app git hash
	Githash = ""
	// GoVersion `go version` string
	GoVersion = ""

	// ServerTag app server tag names
	ServerTag = ""
	// ServerID app server id
	ServerID = ""
)
