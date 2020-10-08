/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log"
)

func main() {
	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	if err := cmdr.Exec(rootCmd,
		// To disable internal commands and flags, uncomment the following codes
		cmdr.WithBuiltinCommands(false, false, false, false, true),
		// daemon.WithDaemon(svr.NewDaemon(), nil, nil, nil),
		cmdr.WithLogx(log.NewStdLoggerWith(log.DebugLevel)),
		// cmdr.WithHelpTabStop(40),
		// cmdr.WithNoColor(true),
	); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

const (
	wgetVersion = "1.20"

	cStartup          = "10.Startup"
	cLogging          = "20.Logging and input file"
	cDownload         = "30.Download"
	cDirectories      = "40.Directories"
	cHTTPOptions      = "50.HTTP Options"
	cHTTPSOptions     = "51.HTTPS (SSL/TLS) options"
	cHstsOptions      = "52.HSTS options"
	cFtpOptions       = "53.FTP options"
	cFtpsOptions      = "54.FTPS options"
	cWarcOptions      = "55.WARC options"
	cRecusiveDownload = "60.Recursive download"
	cRecusiveAccept   = "61.Recursive accept/reject"
)

var (
	rootCmd = &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name: "wget",
			},
			Flags: append(
				startupFlags,
				append(loggerFlags,
					downloadFlags...)...,
			),
			SubCommands: []*cmdr.Command{},
		},

		AppName:    "wget-demo",
		Version:    wgetVersion,
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
				Group: cStartup,
			},
			DefaultValue: wgetVersion,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "h",
				Full:        "help",
				Description: "print this help",
				Action: func(cmd *cmdr.Command, args []string) (err error) {
					// cmd.PrintHelp(false)
					return
				},
				Group: cStartup,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "b",
				Full:        "background",
				Aliases:     []string{"bg"},
				Description: "go to background after startup",
				Group:       cStartup,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Short:       "e",
				Full:        "execute",
				Description: "execute a `.wgetrc'-style command",
				Group:       cStartup,
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "COMMAND",
		},
	}

	loggerFlags = []*cmdr.Flag{
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "001.output-file",
				Short:       "o",
				Full:        "output-file",
				Description: "log messages to FILE",
				Group:       cLogging,
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "FILE",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "011.append-output",
				Short:       "a",
				Full:        "append-output",
				Description: "append messages to FILE",
				Group:       cLogging,
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "FILE",
		},
		{
			// modified, for ~~debug
			BaseOpt: cmdr.BaseOpt{
				Name:        "021.debug",
				Full:        "debug",
				Description: "debug mode",
				Group:       cLogging,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "030.quiet",
				Short:       "q",
				Full:        "quiet",
				Description: "quiet (no output)",
				Group:       cLogging,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "031.verbose",
				Short:       "v",
				Full:        "verbose",
				Description: "be verbose (this is the default)",
				Group:       cLogging,
			},
			DefaultValue: true,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "041.no-verbose",
				Short:       "nv",
				Full:        "no-verbose",
				Description: "turn off verboseness, without being quiet",
				Group:       cLogging,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "051.report-speed",
				Full:        "report-speed",
				Description: "output bandwidth as TYPE.  TYPE can be bits",
				Group:       cLogging,
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "TYPE",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "061.input-file",
				Short:       "i",
				Full:        "input-file",
				Description: "download URLs found in local or external FILE",
				Group:       cLogging,
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "FILE",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "062.force-html",
				Short:       "F",
				Full:        "force-html",
				Description: "treat input file as HTML",
				Group:       cLogging,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "063.base",
				Short:       "B",
				Full:        "base",
				Description: "resolves HTML input-file links (-i -F)  relative to URL",
				Group:       cLogging,
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "URL",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "064.config",
				Full:        "config",
				Description: "specify config file to use",
				Group:       cLogging,
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "FILE",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name: "065.no-config",
				// Short:       "nc", // modified
				Full:        "no-config",
				Description: "do not read any config file",
				Group:       cLogging,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "066.rejected-log",
				Full:        "rejected-log",
				Description: "log reasons for URL rejection to FILE",
				Group:       cLogging,
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "FILE",
		},
	}

	downloadFlags = []*cmdr.Flag{
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "001.trace",
				Short:       "t",
				Full:        "trace",
				Description: "set number of retries to NUMBER (0 unlimits)",
				Group:       cDownload,
			},
			DefaultValue:            0,
			DefaultValuePlaceholder: "NUMBER",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "010.retry-connrefused",
				Full:        "retry-connrefused",
				Description: "retry even if connection is refused",
				Group:       cDownload,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "011.retry-on-http-error",
				Full:        "retry-on-http-error",
				Description: "comma-separated list of HTTP errors to retry",
				Group:       cDownload,
			},
			DefaultValue:            []string{},
			DefaultValuePlaceholder: "ERRORS",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "013.output-document",
				Short:       "O",
				Full:        "output-document",
				Description: "write documents to FILE",
				Group:       cDownload,
			},
			DefaultValue:            "",
			DefaultValuePlaceholder: "FILE",
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "015.no-clobber",
				Short:       "nc",
				Full:        "no-clobber",
				Description: "skip downloads that would download to existing files (overwriting them)",
				Group:       cDownload,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "017.no-netrc",
				Full:        "no-netrc",
				Description: "don't try to obtain credentials from .netrc",
				Group:       cDownload,
			},
			DefaultValue: false,
		},
		{
			BaseOpt: cmdr.BaseOpt{
				Name:        "019.continue",
				Short:       "c",
				Full:        "continue",
				Description: "resume getting a partially-downloaded file",
				Group:       cDownload,
			},
			DefaultValue: false,
		},
	}
)
