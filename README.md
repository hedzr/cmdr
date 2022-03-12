# cmdr

<!-- [![Build Status](https://travis-ci.org/hedzr/cmdr.svg?branch=master)](https://travis-ci.org/hedzr/cmdr) -->

![Go](https://github.com/hedzr/cmdr/workflows/Go/badge.svg)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/cmdr.svg?label=release)](https://github.com/hedzr/cmdr/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/cmdr) [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr?ref=badge_shield)
[![go.dev](https://img.shields.io/badge/go.dev-reference-green)](https://pkg.go.dev/github.com/hedzr/cmdr)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/cmdr)](https://goreportcard.com/report/github.com/hedzr/cmdr)
[![codecov](https://codecov.io/gh/hedzr/cmdr/branch/master/graph/badge.svg)](https://codecov.io/gh/hedzr/cmdr)
[![Coverage Status](https://coveralls.io/repos/github/hedzr/cmdr/badge.svg?branch=master)](https://coveralls.io/github/hedzr/cmdr?branch=master)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#command-line)

<!-- https://gowalker.org/github.com/hedzr/cmdr -->

`cmdr` is a POSIX-compliant, command-line UI (CLI) library in Golang.
It is a getopt-like parser of command-line options,
be compatible with the [getopt_long](http://www.gnu.org/s/libc/manual/html_node/Argument-Syntax.html#Argument-Syntax)
command line UI, which is an extension of the syntax recommended by POSIX.

We made many enhancements beyond the standard library `flag`.

There is a fully-functional `Options Store` (configurations) for your
hierarchical configuration dataset too.

The .netCore version [Cmdr.Core](https://github.com/hedzr/Cmdr.Core) is available now. And, a cxx version [`cmdr-cxx`](https://github.com/hedzr/cmdr-cxx) was pre-released just now (Happy Spring Festival 2021).

![ee99d078e2f7](https://user-images.githubusercontent.com/12786150/72876202-f49ee500-3d30-11ea-9de0-434bf8decf90.gif)

<!-- built by https://ezgif.com/ -->

> See the image frames at [#1](https://github.com/hedzr/cmdr/issues/1#issuecomment-567779978).

See our extras:

- [**cmdr-docs**](https://github.com/hedzr/cmdr-docs): documentations (Working)
- [**cmdr-addons**](https://github.com/hedzr/cmdr-addons): a new daemon plugin `dex` for linux/macOS/windows.
- [**cmdr-examples**](https://github.com/hedzr/cmdr-examples): collects the samples for cmdr
- [**cmdr-go-starter**](https://github.com/hedzr/cmdr-go-starter): public template repo to new your cli app

and Bonus of [#cmdr](https://github.com/topics/cmdr) Series:

- dotnetCore: [Cmdr.Core](https://github.com/hedzr/Cmdr.Core)
- C++17 or higher: [cmdr-cxx](https://github.com/hedzr/cmdr-cxx)

## News

- docs (WIP):
  - english documentation NOT completed yet
  - documentation at: https://hedzr.github.io/cmdr-docs/

- v1.10.32 (WIP)

- v1.10.31
  - routine maintenance release

- v1.10.30
  - routine maintenance release
  - add: examples/example-app
  - imp: use relative path in log output
  - reenable error template

- v1.10.29
  - routine maintenance release

- v1.10.27
  - upgrade to [errors.v3](https://github.com/hedzr/errors)

- v1.10.24
  - fix: version command, help screen not work

- v1.10.23
  - fix: unknown switch char not an error now
  - imp: refined gen zsh code, and gen shell codes
  - fea: support fish-shell completion generating now
  - fea: added root.`RunAsSubCommand`, treat 'app' as a synonym of 'app subcmd1 subcmd2'
  - imp/fix/fea: clarify the slice append or replace mode - SetNx & `SetNxOverwrite` for Option Store
  - fea: added `VendorHidden` field for when you wanna a never shown flag or command
  - fea: conf package - add `Serial`, `SerialString` for CI tool
  - imp: erase man1 folder after `--man`
  - fix/imp: prints description with color escaped, when multiline
  - fix: restore Match() but with new name MatchAndTest()
  - fix: high-order fn hold the older copy, so pass it by holding a pointer to original variable
  - imp: review most of the tests
  - NOTE: cleanup the deprecated codes [`cmd.NewFlagV`,`cmd.NewFlag`, `cmd.NewSubCommand`, ...]
  - fea: more completion supports

- v1.10.19
  - temporary build for earlier testing
  - confirmed: backward compatible with go1.12

- v1.10.13
  - fix/imp: adapt windir to *nix name to fit for generic config file loading
  - fea/imp: improved Aliases algor, support more tmpl var substitute
  - fix: fallback the unknown type as string type
  - fea: add flag to control whether write the changes back to alternative config file or not, `WithAlterConfigAutoWriteBack`
  - imp: name/desc fields of builtin commands and flags
  - CHANGE: use [bgo build-tool](https://github.com/hedzr/bgo) now, Makefile thrown

- More details at [CHANGELOG](./CHANGELOG)

## Features

[Features.md](old/Features.md)

> Old README.md: [README.old.md](old/README.old.md)

## For Developer

[For Developer](old/Developer.md)

### Fast Guide

See [example-app](./tree/master/examples/example-app/)

<details>
	<summary> Expand to source codes </summary>

```go
package main

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/examples/internal"
	"github.com/hedzr/cmdr/plugin/pprof"
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log"
	"github.com/hedzr/log/isdelve"
	"github.com/hedzr/logex/build"
	"gopkg.in/hedzr/errors.v3"
)

func main() {
	Entry()
}

func Entry() {
	root := buildRootCmd()
	if err := cmdr.Exec(root, options...); err != nil {
		log.Fatalf("error occurs in app running: %+v\n", err)
	}
}

func buildRootCmd() (rootCmd *cmdr.RootCommand) {
	root := cmdr.Root(appName, version).
		//AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
		//	// cmdr.Set("enable-ueh", true)
		//	return
		//}).
		//AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
		//	//fmt.Printf("# global pre-action 2, exe-path: %v\n", cmdr.GetExecutablePath())
		//	return
		//}).
		//AddGlobalPostAction(func(cmd *cmdr.Command, args []string) {
		//	//fmt.Println("# global post-action 1")
		//}).
		//AddGlobalPostAction(func(cmd *cmdr.Command, args []string) {
		//	//fmt.Println("# global post-action 2")
		//}).
		Copyright(copyright, "hedzr").
		Description(desc, longDesc).
		Examples(examples)
	rootCmd = root.RootCommand()

	// for your biz-logic, constructing an AttachToCmdr(root *cmdr.RootCmdOpt) is recommended.
	// see our full sample and template repo: https://github.com/hedzr/cmdr-go-starter
	//core.AttachToCmdr(root.RootCmdOpt())

	// These lines are removable

	cmdr.NewBool(false).
		Titles("enable-ueh", "ueh").
		Description("Enables the unhandled exception handler?").
		AttachTo(root)
	//cmdrPanic(root)
	cmdrSoundex(root)
	//pprof.AttachToCmdr(root.RootCmdOpt())
	return
}

func cmdrSoundex(root cmdr.OptCmd) {

	cmdr.NewSubCmd().Titles("soundex", "snd", "sndx", "sound").
		Description("soundex test").
		Group("Test").
		TailPlaceholder("[text1, text2, ...]").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			for ix, s := range args {
				fmt.Printf("%5d. %s => %s\n", ix, s, tool.Soundex(s))
			}
			return
		}).
		AttachTo(root)

}

func onUnhandledErrorHandler(err interface{}) {
	if cmdr.GetBoolR("enable-ueh") {
		dumpStacks()
	}

	panic(err) // re-throw it
}

func dumpStacks() {
	fmt.Printf("\n\n=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===\n\n", errors.DumpStacksAsString(true))
}

func init() {
	options = append(options, cmdr.WithUnhandledErrorHandler(onUnhandledErrorHandler))

	options = append(options,
		cmdr.WithLogx(build.New(build.NewLoggerConfigWith(
			defaultDebugEnabled, defaultLoggerBackend, defaultLoggerLevel,
			log.WithTimestamp(true, "")))))

	options = append(options, cmdr.WithHelpTailLine(`
# Type '-h'/'-?' or '--help' to get command help screen.
# Star me if it's helpful: https://github.com/hedzr/cmdr/examples/example-app
`))

	if isDebugBuild() {
		options = append(options, pprof.GetCmdrProfilingOptions())
	}

	// enable '--trace' command line option to toggle a internal trace mode (can be retrieved by cmdr.GetTraceMode())
	// import "github.com/hedzr/cmdr-addons/pkg/plugins/trace"
	// trace.WithTraceEnable(defaultTraceEnabled)
	// Or:
	optAddTraceOption := cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
		cmdr.NewBool(false).
			Titles("trace", "tr").
			Description("enable trace mode for tcp/mqtt send/recv data dump", "").
			//Action(func(cmd *cmdr.Command, args []string) (err error) { println("trace mode on"); cmdr.SetTraceMode(true); return; }).
			Group(cmdr.SysMgmtGroup).
			AttachToRoot(root)
	}, nil)
	options = append(options, optAddTraceOption)
	//options = append(options, optAddServerExtOpt«ion)

	// allow and search '.<appname>.yml' at first
	locations := []string{".$APPNAME.yml"}
	locations = append(locations, cmdr.GetPredefinedLocations()...)
	options = append(options, cmdr.WithPredefinedLocations(locations...))

	options = append(options, internal.NewAppOption())
}

func isDebugBuild() bool { return isdelve.Enabled }

var options []cmdr.ExecOption

//goland:noinspection GoNameStartsWithPackageName
const (
	appName   = "example-app"
	version   = "0.2.5"
	copyright = "example-app - A devops tool - cmdr series"
	desc      = "example-app is an effective devops tool. It make an demo application for 'cmdr'"
	longDesc  = `example-app is an effective devops tool. It make an demo application for 'cmdr'.
`
	examples = `
$ {{.AppName}} gen shell [--bash|--zsh|--fish|--auto]
  generate bash/shell completion scripts
$ {{.AppName}} gen man
  generate linux man page 1
$ {{.AppName}} --help
  show help screen.
$ {{.AppName}} --help --man
  show help screen in manpage viewer (for linux/darwin).
`
	overview = ``

	zero = 0

	defaultTraceEnabled  = true
	defaultDebugEnabled  = false
	defaultLoggerLevel   = "debug"
	defaultLoggerBackend = "logrus"
)
```
</details>

### About the Docker build

Here is a docker build for cmdr/examples/fluent so that you can run it without go building or downloading the release files:

```bash
# from Docker Hub:
$ docker run -it --rm hedzr/cmdr-fluent
$ docker run -it --rm hedzr/cmdr-fluent --help

# from Github Packages (please following the guide of GitHub Packages Site):
$ docker run -it --rm docker.pkg.github.com/hedzr/cmdr/cmdr-fluent
$ docker run -it --rm docker.pkg.github.com/hedzr/cmdr/cmdr-fluent --help
```

### Examples

1. [**short**](./examples/short/README.md)  
   simple codes with structured data style.

2. [demo](./examples/demo/README.md)  
   normal demo with external config files.

3. [wget-demo](./examples/wget-demo/README.md)  
   partial-covered for GNU `wget`.

4. [fluent](./examples/fluent)  
   demostrates how to define your command-ui with the fluent api style.

5. [ffdemo](./examples/ffdemo)

   a demo to show you how to migrate from go `flag` smoothly.

6. [cmdr-http2](https://github.com/hedzr/cmdr-http2)  
   http2 server with daemon supports, graceful shutdown

7. [awesome-tool](https://github.com/hedzr/awesome-tool)  
   `awesome-tool` is a cli app that fetch the repo stars and generate a markdown summary, accordingly with most of awesome-xxx list in github (such as awesome-go).

**See Also the examples index: [Examples.md](old/Examples.md)** _(zh-cn TODO: [Examples.zh-cn.md](old/Examples.zh-cn.md))_

## Uses

-   https://github.com/hedzr/consul-tags
-   https://github.com/hedzr/ini-op
-   https://github.com/hedzr/awesome-tool
-   austr
-   Issue me to adding yours

## Contrib

_Feel free to issue me bug reports and fixes. Many thanks to all contributors._

## Thanks to JODL

Thanks to [JetBrains](https://www.jetbrains.com/?from=cmdr) for donating product licenses to help develop **cmdr** [![jetbrains](https://gist.githubusercontent.com/hedzr/447849cb44138885e75fe46f1e35b4a0/raw/bedfe6923510405ade4c034c5c5085487532dee4/jetbrains-variant-4.svg)](https://www.jetbrains.com/?from=hedzr/cmdr)

[![goland](https://gist.githubusercontent.com/hedzr/447849cb44138885e75fe46f1e35b4a0/raw/ca8ac2694906f5650d585263dbabfda52072f707/logo-goland.svg)](https://www.jetbrains.com/?from=hedzr/cmdr)

## License

MIT

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr?ref=badge_large)