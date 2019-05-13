/*
 * Copyright © 2019 Hedzr Yeh.
 */

package conf

var (
	CfgFile string
	AppName string

	// these 3 variables will be rewrote when app had been building by ci-tool
	Version    = "0.2.1"
	Buildstamp = ""
	Githash    = ""

	ServerTag = ""
	ServerID  = ""
)
