// Copyright © 2025 Hedzr Yeh.

// Package cmdr provides an argument parser for golang cli app.
//
// cmdr/v2 parses end-user's commandline input and forwards to
// nested subcommands.
//
// cmdr also provides a high-performance hierachical configuration
// data manager by integrating with `hedzr/store`(https://github.com/hedzr/store). This part
// loads app settings from pluggable external sources.
//
// The basic feature of cmdr/v2 are
//
//   - Basic command-line arguments parser like POSIX getopt and go stdlib flag.
//
//   - Short flag, single character or a string here to support golang CLI style
//
//     Compact flags if possible. Also the sticky value will be parsed. For example: `-c1b23zv` = `-c 1 -b 23 -z -v`
//
//     Hit info: `-v -v -v` = `-v` (hitCount == 3, hitTitle == 'v')
//
//     Optimized for slice: `-a 1,2,3 -a 4 -a 5,6` => []int{1,2,3,4,5,6}
//
//     Value can be sticked or not. Valid forms: `-c1`, `-c 1`, `-c=1` and quoted: `-c"1"`, `-c'1'`, `-c="1"`, `-c='1'`, etc.
//
//     ...
//
//   - Long flags and aliases
//
//     = Eventual subcommands: an `OnAction` handler can be attached.
//
//   - Eventual subcommands and flags: PreActions, PostAction, OnMatching, OnMatched, ...,
//
//   - Auto bind to environment variables, For instance: command line `HELP=1 app` = `app --help`.
//
//   - Builtin commands and flags:
//
//     `--help`, `-h`
//
//     `--version`, `-V`
//
//     `--verbose`. `-v`
//
//     ...
//
//   - Help Screen: auto generate and print
//
//   - Smart suggestions when wrong cmd or flag parsed. Jaro-winkler distance is used.
//
//   - Loosely parse subcmds and flags:
//
//     Subcommands and flags can be input in any order
//
//     Lookup a flag along with subcommands tree for resolving the duplicated flags
//
//   - Can integrate with [hedzr/store](https://github.com/hedzr/store)
//
//     High-performance in-memory KV store for hierarchical data.
//
//     Extract data to user-spec type with auto-converting
//
//     Loadable external sources: environ, config files, consul, etcd, etc..
//
//     extensible codecs and providers for loading from data sources
//
//   - Generating shell autocompletion scripts
//
//     supported shells are: zsh, bash, fish, powershell, ...
//
//     auto install the scripts for zsh shell.
//
//   - Generating command manpages for software deployment time.
//
// For more documentation please go and check out https://docs.hedzr.com/cmdr.v2/.
//
// A tiny sample app can be:
//
//	package main
//
//	import (
//		"context"
//		"os"
//
//		"github.com/hedzr/cmdr/v2"
//		"github.com/hedzr/cmdr/v2/examples/cmd"
//		"github.com/hedzr/cmdr/v2/pkg/logz"
//	)
//
//	const (
//		appName = "concise"
//		desc    = `concise version of tiny app.`
//		version = cmdr.Version
//		author  = `The Example Authors`
//	)
//
//	func main() {
//		app := cmdr.Create(appName, version, author, desc).
//			WithAdders(cmd.Commands...).
//			Build()
//
//		ctx, cancel := context.WithCancel(context.Background())
//		defer cancel()
//
//		if err := app.Run(ctx); err != nil {
//			logz.ErrorContext(ctx, "Application Error:", "err", err) // stacktrace if in debug mode/build
//			os.Exit(app.SuggestRetCode())
//		} else if rc := app.SuggestRetCode(); rc != 0 {
//			os.Exit(rc)
//		}
//	}
//
// Getting started from [New] or [Create] function.
//
// For more documentation please go and check out https://docs.hedzr.com/cmdr.v2/.
package cmdr

const Version = "v2.1.35" // Version fir hedzr/cmdr/v2
