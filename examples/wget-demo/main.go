/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"github.com/hedzr/cmdr"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true,})
	cmdr.EnableVersionCommands = false
	cmdr.EnableVerboseCommands = false
	cmdr.EnableHelpCommands = false
	if err := cmdr.Exec(rootCmd); err != nil {
		logrus.Errorf("Error: %v", err)
	}
}

const (
	VERSION = "1.20"

	STARTUP           = "10.Startup"
	LOGGING           = "20.Logging and input file"
	DOWNLOAD          = "30.Download"
	DIRECTORIES       = "40.Directories"
	HTTP_OPTIONS      = "50.HTTP Options"
	HTTPS_OPTIONS     = "51.HTTPS (SSL/TLS) options"
	HSTS_OPTIONS      = "52.HSTS options"
	FTP_OPTIONS       = "53.FTP options"
	FTPS_OPTIONS      = "54.FTPS options"
	WARC_OPTIONS      = "55.WARC options"
	RECUSIVE_DOWNLOAD = "60.Recursive download"
	RECUSIVE_ACCEPT   = "61.Recursive accept/reject"
)

var (
	rootCmd = &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name: "wget",
				Flags: append(
					startupFlags,
					append(loggerFlags,
						downloadFlags...)...
				),
			},
			SubCommands: []*cmdr.Command{
			},
		},

		AppName:    "wget-demo",
		Version:    VERSION,
		VersionInt: 0x011400,
		Header: `GNU Wget 1.20, a non-interactive network retriever.

Usage: wget [OPTION]... [URL]...

Mandatory arguments to long options are mandatory for short options too.`,
		Author: "Hedzr Yeh <hedzrz@gmail.com>",
	}

	startupFlags = []*cmdr.Flag{
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "V",
				Full:        "version",
				Description: "display the version of Wget and exit",
				Action: func(cmd *cmdr.Command, args []string) (err error) {
					cmd.PrintVersion()
					return
				},
				Group: STARTUP,
			},
			DefaultValue: VERSION,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "h",
				Full:        "help",
				Description: "print this help",
				Action: func(cmd *cmdr.Command, args []string) (err error) {
					cmd.PrintHelp(false)
					return
				},
				Group: STARTUP,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "b",
				Full:        "background",
				Aliases:     []string{"bg",},
				Description: "go to background after startup",
				Group:       STARTUP,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:                   "e",
				Full:                    "execute",
				Description:             "execute a `.wgetrc'-style command",
				Group:                   STARTUP,
				DefaultValuePlaceholder: "COMMAND",
			},
			DefaultValue: "",
		},
	}

	loggerFlags = []*cmdr.Flag{
		{
			BaseOpt: cmdr.BaseOpt{
				Name:                    "001.output-file",
				Short:                   "o",
				Full:                    "output-file",
				Description:             "log messages to FILE",
				Group:                   LOGGING,
				DefaultValuePlaceholder: "FILE",
			},
			DefaultValue: "",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:                    "011.append-output",
				Short:                   "a",
				Full:                    "append-output",
				Description:             "append messages to FILE",
				Group:                   LOGGING,
				DefaultValuePlaceholder: "FILE",
			},
			DefaultValue: "",
		},
		{
			// modified, for ~~debug
			BaseOpt: cmdr.BaseOpt{
				Name:        "021.debug",
				Full:        "debug",
				Description: "debug mode",
				Group:       LOGGING,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "030.quiet",
				Short:       "q",
				Full:        "quiet",
				Description: "quiet (no output)",
				Group:       LOGGING,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "031.verbose",
				Short:       "v",
				Full:        "verbose",
				Description: "be verbose (this is the default)",
				Group:       LOGGING,
			},
			DefaultValue: true,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "041.no-verbose",
				Short:       "nv",
				Full:        "no-verbose",
				Description: "turn off verboseness, without being quiet",
				Group:       LOGGING,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:                    "051.report-speed",
				Full:                    "report-speed",
				Description:             "output bandwidth as TYPE.  TYPE can be bits",
				Group:                   LOGGING,
				DefaultValuePlaceholder: "TYPE",
			},
			DefaultValue: "",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:                    "061.input-file",
				Short:                   "i",
				Full:                    "input-file",
				Description:             "download URLs found in local or external FILE",
				Group:                   LOGGING,
				DefaultValuePlaceholder: "FILE",
			},
			DefaultValue: "",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "062.force-html",
				Short:       "F",
				Full:        "force-html",
				Description: "treat input file as HTML",
				Group:       LOGGING,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:                    "063.base",
				Short:                   "B",
				Full:                    "base",
				Description:             "resolves HTML input-file links (-i -F)  relative to URL",
				Group:                   LOGGING,
				DefaultValuePlaceholder: "URL",
			},
			DefaultValue: "",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:                    "064.config",
				Full:                    "config",
				Description:             "specify config file to use",
				Group:                   LOGGING,
				DefaultValuePlaceholder: "FILE",
			},
			DefaultValue: "",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name: "065.no-config",
				// Short:       "nc", // modified
				Full:        "no-config",
				Description: "do not read any config file",
				Group:       LOGGING,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:                    "066.rejected-log",
				Full:                    "rejected-log",
				Description:             "log reasons for URL rejection to FILE",
				Group:                   LOGGING,
				DefaultValuePlaceholder: "FILE",
			},
			DefaultValue: "",
		},
	}

	downloadFlags = []*cmdr.Flag{
		{
			BaseOpt: cmdr.BaseOpt{
				Name:                    "001.trace",
				Short:                   "t",
				Full:                    "trace",
				Description:             "set number of retries to NUMBER (0 unlimits)",
				Group:                   DOWNLOAD,
				DefaultValuePlaceholder: "NUMBER",
			},
			DefaultValue: 0,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "010.retry-connrefused",
				Full:        "retry-connrefused",
				Description: "retry even if connection is refused",
				Group:       DOWNLOAD,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:                    "011.retry-on-http-error",
				Full:                    "retry-on-http-error",
				Description:             "comma-separated list of HTTP errors to retry",
				Group:                   DOWNLOAD,
				DefaultValuePlaceholder: "ERRORS",
			},
			DefaultValue: []string{},
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:                    "013.output-document",
				Short:                   "O",
				Full:                    "output-document",
				Description:             "write documents to FILE",
				Group:                   DOWNLOAD,
				DefaultValuePlaceholder: "FILE",
			},
			DefaultValue: "",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "015.no-clobber",
				Short:       "nc",
				Full:        "no-clobber",
				Description: "skip downloads that would download to existing files (overwriting them)",
				Group:       DOWNLOAD,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "017.no-netrc",
				Full:        "no-netrc",
				Description: "don't try to obtain credentials from .netrc",
				Group:       DOWNLOAD,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "019.continue",
				Short:       "c",
				Full:        "continue",
				Description: "resume getting a partially-downloaded file",
				Group:       DOWNLOAD,
			},
			DefaultValue: false,
		},
	}
)
