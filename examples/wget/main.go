package main

import (
	"context"
	"os"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
)

func main() {
	ctx := context.Background() // with cancel can be passed thru in your actions
	app := prepareApp()
	if err := app.Run(ctx); err != nil {
		println("Application Error:", err)
		os.Exit(app.SuggestRetCode())
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

func prepareApp(opts ...cli.Opt) (app cli.App) {
	app = cmdr.New(opts...).
		Info("wget", wgetVersion).
		Author("The Example Authors").
		Header(`GNU Wget 1.20, a non-interactive network retriever.

Usage: wget [OPTION]... [URL]...

Mandatory arguments to long options are mandatory for short options too.`) // .Description(``).Header(``).Footer(``)

	app.With(func(app cli.App) {
		app.Flg("background", "b", "bg").Description(`go to background after startup`).Group(cStartup).Build()
		app.Flg("execute", "e").Description(`execute a '.wgetrc'-style command`).Group(cStartup).Default("").PlaceHolder("COMMAND").Build()

		app.Flg("output-file", "o").Description(`log messages to FILE`).Group(cLogging).PlaceHolder("FILE").Default("").Build()
		app.Flg("append-output", "a").Description(`append messages to FILE`).Group(cLogging).PlaceHolder("FILE").Default("").Build()
		app.Flg("no-verbose", "nv").Description(`turn off verboseness, without being quiet`).Group(cLogging).Default("").Build()
		app.Flg("report-speed", "").Description(`output bandwidth as TYPE.  TYPE can be bits`).Group(cLogging).Default("").PlaceHolder("TYPE").Build()
		app.Flg("input-file", "i").Description(`download URLs found in local or external FILE`).Group(cLogging).Default("").PlaceHolder("FILE").Build()
		app.Flg("force-html", "F").Description(`treat input file as HTML`).Group(cLogging).Build()
		app.Flg("base", "b").Description(`resolves HTML input-file links (-i -F)  relative to URL`).Default("").PlaceHolder("URL").Group(cLogging).Build()
		app.Flg("config").Description(`specify config file to use`).Group(cLogging).Default("").PlaceHolder("FILE").Build()
		app.Flg("no-config").Description(`do not read any config file`).Group(cLogging).Build()
		app.Flg("rejected-log").Description(`log reasons for URL rejection to FILE`).Group(cLogging).Default("").PlaceHolder("FILE").Build()

		app.Flg("retry-connrefused", "").Description(`retry even if connection is refused`).Group(cDownload).Build()
		app.Flg("retry-on-http-error").Description(`comma-separated list of HTTP errors to retry`).Group(cDownload).Default([]string{}).PlaceHolder("ERRORS").Build()
		app.Flg("output-document", "O").Description(`write documents to FILE`).Group(cDownload).Default("").PlaceHolder("FILE").Build()
		app.Flg("no-clobber", "nc").Description(`skip downloads that would download to existing files (overwriting them)`).Group(cDownload).Build()
		app.Flg("no-netrc").Description(`don't try to obtain credentials from .netrc`).Group(cDownload).Build()
		app.Flg("continue", "c").Description(`resume getting a partially-downloaded file`).Group(cDownload).Build()
	})

	app.OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
		println(`wget ran.`)
		_, err = cmd.App().DoBuiltinAction(ctx, cli.ActionShowHelpScreen)
		return
	})

	// app.RootBuilder(func(parent cli.CommandBuilder) {
	// 	// parent.Cmd(...)...
	// 	parent.OnAction(...)
	// })
	return
}
