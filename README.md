# cmdr

![Go](https://github.com/hedzr/cmdr/workflows/Go/badge.svg)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/hedzr/store)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/cmdr.svg?label=release)](https://github.com/hedzr/cmdr/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/cmdr) [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr?ref=badge_shield)
[![go.dev](https://img.shields.io/badge/go.dev-reference-green)](https://pkg.go.dev/github.com/hedzr/cmdr)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/cmdr)](https://goreportcard.com/report/github.com/hedzr/cmdr)
[![codecov](https://codecov.io/gh/hedzr/cmdr/branch/master/graph/badge.svg)](https://codecov.io/gh/hedzr/cmdr)<!--
[![Coverage Status](https://coveralls.io/repos/github/hedzr/cmdr/badge.svg?branch=master)](https://coveralls.io/github/hedzr/cmdr?branch=master)-->
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#command-line)

> **Unstable v2**:
> Working in progress, the main API might be stable till v2.1.0

`cmdr` is a POSIX-compliant, command-line argument parser library with Golang.

Since v2, our license moved to Apache 2.0.

![cover](https://user-images.githubusercontent.com/12786150/72876202-f49ee500-3d30-11ea-9de0-434bf8decf90.gif)<!-- built by https://ezgif.com/ -->

> ~~See the image frames at [#1](https://github.com/hedzr/cmdr/issues/1#issuecomment-567779978).~~

## Motivation

There are many dirty codes in the cmdr.v1 which cannot be refactored as well. It prompted we reimplment a new one as v2.

The passing winter, we did rewrite the cmdr.v2 to keep it clean and absorbed in parsing and dispatching.
Some abilities were removed and relayouted to new modules.
That's why the `Option Store` has been split as a standalone module [hedzr/store](https://github.com/hedzr/store)[^1].
A faster and colorful slog-like logger has been implemented freshly as [hedzr/logg](https://github.com/hedzr/logg)[^3].
[hedzr/evendeep](https://github.com/hedzr/evendeep)[^2] provides a deep fully-functional object copy tool. It helps to deep copy some internal objects easily. It is also ready for you.
[hedzr/is](https://github.com/hedzr/is)[^4] is an environment detecting framework with many out-of-the-box detectors, such as `is.InTesting` and `is.InDebugging`.

Anyway, the whole supply chain painted:

```mermaid
graph BT
  hzis(hedzr/is)-->hzlogg(hedzr/logg/slog)
  hzis-->hzdiff(hedzr/evendeep)
  hzlogg-->hzdiff
  hzerrors(gopkg.in/hedzr/errors.v3)-->hzdiff
  hzerrors-->hzstore(hedzr/store)
  hzis-->hzstore(hedzr/store)
  hzlogg-->hzstore(hedzr/store)
  hzdiff-->hzstore(hedzr/store)
  hzlogg-->cmdr(hedzr/cmdr/v2)
  hzis-->cmdr
  hzlogg-->cmdr
  hzdiff-->cmdr
  hzstore-->cmdr

```

> 1. The .netCore version [Cmdr.Core](https://github.com/hedzr/Cmdr.Core) is available now.
> 2. A cxx version [`cmdr-cxx`](https://github.com/hedzr/cmdr-cxx) was released (Happy Spring Festival 2021).

## Features

v2 is in earlier state but the baseline is stable:

- Basic command-line arguments parser like POSIX getopt and go stdlib flag.
  - Short flag, single character or a string here to support golang CLI style
    - Compact flags if possible. Also the sticking value will be parsed. For example: `-c1b23zv` = `-c 1 -b 23 -z -v`
    - Hit info: `-v -v -v` = `-v` (hitCount == 3, hitTitle == 'v')
    - Optimized for slice: `-a 1,2,3 -a 4 -a 5,6` => []int{1,2,3,4,5,6}
    - Value can be sticked or not. Valid forms: `-c1`, `-c 1`, `-c=1` and quoted: `-c"1"`, `-c'1'`, `-c="1"`, `-c='1'`, etc.
    - ...

  - Long flags and aliases
  - Eventual subcommands: an `OnAction` handler can be attached.
  - Eventual subcommands and flags: PreActions, PostAction, OnMatching, OnMatched, ...,
  - Auto bind to environment variables, For instance: command line `HELP=1 app` = `app --help`.
  - Builtin commands and flags:
    - `--help`, `-h`
    - `--version`, `-V`
    - `--verbose`. `-v`
    - ...

  - Help Screen: auto generate and print
  - Smart suggestions when wrong cmd or flag parsed. Jaro-winkler distance is used.

- Loosely parse subcmds and flags:
  - Subcommands and flags can be input in any order
  - Lookup a flag along with subcommands tree for resolving the duplicated flags

- Can integrate with [hedzr/store](https://github.com/hedzr/store)[^1]
  - High-performance in-memory KV store for hierarchical data.
  - Extract data to user-spec type with auto-converting
  - Loadable external sources: environ, config files, consul, etcd, etc..
    - extensible codecs and providers for loading from data sources

- Three kinds of config files are searched and loaded via `loaders.NewConfigFileLoader()`:
  - Primary: main config, shipped with installable package.
  - Secondary: 2ndry config. Wrapped by reseller(s).
  - Alternative: user's local config, writeable. The runtime changeset will be written back to this file while app stopping.

- TODO
  - Shell autocompletion
  - ...

[^1]: `hedzr/store` is a high-performance configure management library
[^2]: `hedzr/evendeep` offers a customizable deepcopy tool to you. There are also deepequal, deepdiff tools in it.
[^3]: `hedzr/logg` provides a slog like and colorful logging library
[^4]: `hedzr/is` is a basic environ detectors library

More minor details need to be evaluated and reimplemented if it's still meaningful in v2.

Since v2.0.3, loaders had been splitted as a standalone repo so that we can keep cmdr v2 smaller and independer. See the
relevant subproject [cmdr-loaders](https://github.com/hedzr/cmdr-loaders).

## History

v2 is staying in earlier state:

- Full list: [CHANGELOG](https://github.com/hedzr/cmdr/blob/master/CHANGELOG)

## Guide

A simple cli-app can be:

```go
package main

import (
	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/cmdr/v2/pkg/dir"
	"github.com/hedzr/store"
)

func main() {
	ctx := context.Background()
	app := prepareApp(
		cmdr.WithStore(store.New()), // use a option store, if not specified by store.New(), a dummy store allocated

		// cmdr.WithExternalLoaders(
		// 	local.NewConfigFileLoader(),      // import "github.com/hedzr/cmdr-loaders/local" to get in advanced external loading features
		// 	local.NewEnvVarLoader(),
		// ),

		cmdr.WithTasksBeforeRun(func(ctx context.Context, cmd cli.Cmd, runner cli.Runner, extras ...any) (err error) {
			logz.DebugContext(ctx, "command running...", "cmd", cmd, "runner", runner, "extras", extras)
			return
		}),

		// true for debug in developing time, it'll disable onAction on each Cmd.
		// for productive mode, comment this line.
		cmdr.WithForceDefaultAction(true),
	)
	if err := app.Run(ctx); err != nil {
		logz.ErrorContext(ctx, "Application Error:", "err", err)
		os.Exit(app.SuggestRetCode())
	}
}

func prepareApp() (app cli.App) {
	app = cmdr.New().
		Info("demo-app", "0.3.1").
		Author("hedzr")

	app.Flg("no-default").
		Description("disable force default action").
		OnMatched(func(f *cli.Flag, position int, hitState *cli.MatchState) (err error) {
			app.Store().Set("app.force-default-action", false)
			return
		}).
		Build()

	b := app.Cmd("jump").
		Description("jump command").
		Examples(`jump example`).
		Deprecated(`jump is a demo command`).
		Hidden(false)

	b1 := b.Cmd("to").
		Description("to command").
		Examples(``).
		Deprecated(`v0.1.1`).
		Hidden(false).
		OnAction(func(cmd *cli.Command, args []string) (err error) {
			app.Store().Set("app.demo.working", dir.GetCurrentDir())
			println()
			println(dir.GetCurrentDir())
			println()
			println(app.Store().Dump())
			return // handling command action here
		})
	b1.Flg("full", "f").
		Default(false).
		Description("full command").
		Build()
	b1.Build()

	b.Build()

	app.Flg("dry-run", "n").
		Default(false).
		Description("run all but without committing").
		Build()

	app.Flg("wet-run", "w").
		Default(false).
		Description("run all but with committing").
		Build() // no matter even if you're adding the duplicated one.
	return
}
```

More examples please go to [cmdr-tests/examples](https://github.com/hedzr/cmdr-tests/tree/master/examples).

## Thanks to JODL

Thanks to [JetBrains](https://www.jetbrains.com/?from=cmdr) for donating product licenses to help develop **cmdr**  
[![jetbrains](https://gist.githubusercontent.com/hedzr/447849cb44138885e75fe46f1e35b4a0/raw/bedfe6923510405ade4c034c5c5085487532dee4/jetbrains-variant-4.svg)](https://www.jetbrains.com/?from=hedzr/cmdr)
[![goland](https://gist.githubusercontent.com/hedzr/447849cb44138885e75fe46f1e35b4a0/raw/ca8ac2694906f5650d585263dbabfda52072f707/logo-goland.svg)](https://www.jetbrains.com/?from=hedzr/cmdr)

## License

Since v2, our license moved to Apache 2.0.

The v1 keeps under MIT itself.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr?ref=badge_large)
