# cmdr

<!-- [![Build Status](https://travis-ci.org/hedzr/cmdr.svg?branch=master)](https://travis-ci.org/hedzr/cmdr) -->

![Go](https://github.com/hedzr/cmdr/workflows/Go/badge.svg)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/cmdr.svg?label=release)](https://github.com/hedzr/cmdr/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/cmdr) [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr?ref=badge_shield)
[![go.dev](https://img.shields.io/badge/go.dev-reference-green)](https://pkg.go.dev/github.com/hedzr/cmdr)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/cmdr)](https://goreportcard.com/report/github.com/hedzr/cmdr)
[![codecov](https://codecov.io/gh/hedzr/cmdr/branch/master/graph/badge.svg)](https://codecov.io/gh/hedzr/cmdr)<!--
[![Coverage Status](https://coveralls.io/repos/github/hedzr/cmdr/badge.svg?branch=master)](https://coveralls.io/github/hedzr/cmdr?branch=master)-->
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#command-line)

`cmdr` is a POSIX-compliant, command-line argument parser library in Golang.
It is a getopt-like parser of command-line options,
be compatible with the [getopt_long](http://www.gnu.org/s/libc/manual/html_node/Argument-Syntax.html#Argument-Syntax)
command line UI, which is an extension of the syntax recommended by POSIX.

It's a replacement of the standard library `flag`.

Also, there is a fully-functional `Options Store` (configurations) integrated for your hierarchical configuration dataset.

The .netCore version [Cmdr.Core](https://github.com/hedzr/Cmdr.Core) is available now. A cxx version [`cmdr-cxx`](https://github.com/hedzr/cmdr-cxx) was released (Happy Spring Festival 2021).

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

> **FRZ**: It means I'm planning and in working for cmdr.v2, but the plan is delayed. In recently months, sadly, some things took my focus. But the v2 will be here soon.  
> In our v2's wishlist, codes will be rewritten completely since I understand coding more in golang, well, perhaps, I hope so.

- docs (WIP):
  - english documentation: NOT yet
  - documentation at: <https://hedzr.github.io/cmdr-docs/>

- v1.11.26 (FRZ)
  - security patch: upgrade golang.org/x/crypto to 0.17.0

- v1.11.25 (FRZ)
  - security patch: google/x/net and more
  - upgrade deps

- v1.11.23 (FRZ)
  - security patch: google/x/net and more
  - improved some writings
  - fixed/improved cmdr.Get() to extract some fields more deeply
  - added Flag.GetDottedNamePathEx()
  - improved defaultActionImpl to flush all outputs at once
  - make default action outputs colorful

- v1.11.21 (FRZ)
  - NOTE: switched toolchain to go1.18 (see the go.mod)
  - now you may return `ErrNotImpl` to fallback to default action at early developing

    <details>
      <summary> Expand to details </summary>

       There are 3 ways to make or fallback to cmdr default impl to get some debug information:
         
         1. set action to nil;
         2. use envvar `FORCE_DEFAULT_ACTION=1`;
         3. return `cmdr.ErrNotImpl` in your action handler:
       
       ```go
       listCmd := cmdr.NewSubCmd().Titles("list", "ls").
         Description("list projects in .poi.yml").
         Action(func(cmd *cmdr.Command, args []string) (err error) {
           err = cmdr.ErrNotImpl // notify caller fallback to defaultAction
           return
         }).
         AttachTo(root)
       ```
    
       > The default action of a command handler will print flags hit info and others useful info.

    </details>

- More details at [CHANGELOG](https://github.com/hedzr/cmdr/blob/master/CHANGELOG)

## Features

[Features.md](old/Features.md)

> Old README.md: [README.old.md](old/README.old.md)

## For Developer

[For Developer](old/Developer.md)

### Fast Guide

See [example-app](https://github.com/hedzr/cmdr/tree/master/examples/example-app/), [examples/](https://github.com/hedzr/cmdr/tree/master/examples/), and [**cmdr-examples**](https://github.com/hedzr/cmdr-examples)

<details>
  <summary> Expand to source codes </summary>

```go
package main

import (
	"fmt"

	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/pprof"
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log"
	"github.com/hedzr/log/buildtags"
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
		// AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
		//	// cmdr.Set("enable-ueh", true)
		//	return
		// }).
		// AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
		//	//fmt.Printf("# global pre-action 2, exe-path: %v\n", cmdr.GetExecutablePath())
		//	return
		// }).
		// AddGlobalPostAction(func(cmd *cmdr.Command, args []string) {
		//	//fmt.Println("# global post-action 1")
		// }).
		// AddGlobalPostAction(func(cmd *cmdr.Command, args []string) {
		//	//fmt.Println("# global post-action 2")
		// }).
		Copyright(copyright, "hedzr").
		Description(desc, longDesc).
		Examples(examples)
	rootCmd = root.RootCommand()

	// for your biz-logic, constructing an AttachToCmdr(root *cmdr.RootCmdOpt) is recommended.
	// see our full sample and template repo: https://github.com/hedzr/cmdr-go-starter
	// core.AttachToCmdr(root.RootCmdOpt())

	// These lines are removable

	cmdr.NewBool(false).
		Titles("enable-ueh", "ueh").
		Description("Enables the unhandled exception handler?").
		AttachTo(root)
	// cmdrPanic(root)
	cmdrSoundex(root)
	// pprof.AttachToCmdr(root.RootCmdOpt())
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
	options = append(options,
		cmdr.WithUnhandledErrorHandler(onUnhandledErrorHandler),

		// cmdr.WithLogxShort(defaultDebugEnabled,defaultLoggerBackend,defaultLoggerLevel),
		cmdr.WithLogx(build.New(build.NewLoggerConfigWith(
			defaultDebugEnabled, defaultLoggerBackend, defaultLoggerLevel,
			log.WithTimestamp(true, "")))),

		cmdr.WithHelpTailLine(`
# Type '-h'/'-?' or '--help' to get command help screen.
# Star me if it's helpful: https://github.com/hedzr/cmdr/examples/example-app
`),
	)

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
			// Action(func(cmd *cmdr.Command, args []string) (err error) { println("trace mode on"); cmdr.SetTraceMode(true); return; }).
			Group(cmdr.SysMgmtGroup).
			AttachToRoot(root)
	}, nil)
	options = append(options, optAddTraceOption)
	// options = append(options, optAddServerExtOpt«ion)

	// allow and search '.<appname>.yml' at first
	locations := []string{".$APPNAME.yml"}
	locations = append(locations, cmdr.GetPredefinedLocations()...)
	options = append(options, cmdr.WithPredefinedLocations(locations...))

	// options = append(options, internal.NewAppOption())
}

var options []cmdr.ExecOption

func isDebugBuild() bool         { return isdelve.Enabled }
func isDockerBuild() bool        { return buildtags.IsDockerBuild() }
func isRunningInDockerEnv() bool { return cmdr.InDockerEnv() }

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
	// overview = ``
	// zero = 0

	// defaultTraceEnabled  = true
	defaultDebugEnabled  = false
	defaultLoggerLevel   = "debug"
	defaultLoggerBackend = "logrus"
)
```

</details>

### Tips for Building Your App

As building your app with cmdr, some build tags are suggested:

```bash
export GIT_REVISION="$(git rev-parse --short HEAD)"
export GIT_SUMMARY="$(git describe --tags --dirty --always)"
export GOVERSION="$(go version)"
export BUILDTIME="$(date -u '+%Y-%m-%d_%H-%M-%S')"
export VERSION="$(grep -E "Version[ \t]+=[ \t]+" doc.go|grep -Eo "[0-9.]+")"
export W_PKG="github.com/hedzr/cmdr/conf"
export LDFLAGS="-s -w \
    -X '$W_PKG.Githash=$GIT_REVISION' \
    -X '$W_PKG.GitSummary=$GIT_SUMMARY' \
    -X '$W_PKG.GoVersion=$GOVERSION' \
    -X '$W_PKG.Buildstamp=$BUILDTIME' \
    -X '$W_PKG.Version=$VERSION'"

go build -v -ldflags "$LDFLAGS" -o ./bin/your-app ./your-app/
```

## Contrib

_Feel free to issue me bug reports and fixes. Many thanks to all contributors._

## Thanks to JODL

Thanks to [JetBrains](https://www.jetbrains.com/?from=cmdr) for donating product licenses to help develop **cmdr**  
	[![jetbrains](https://gist.githubusercontent.com/hedzr/447849cb44138885e75fe46f1e35b4a0/raw/bedfe6923510405ade4c034c5c5085487532dee4/jetbrains-variant-4.svg)](https://www.jetbrains.com/?from=hedzr/cmdr)
[![goland](https://gist.githubusercontent.com/hedzr/447849cb44138885e75fe46f1e35b4a0/raw/ca8ac2694906f5650d585263dbabfda52072f707/logo-goland.svg)](https://www.jetbrains.com/?from=hedzr/cmdr)

## License

MIT

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr?ref=badge_large)
