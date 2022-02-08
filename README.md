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

- v1.10.23 (WIP)
  - fix: unknown switch char not an error now
  - imp: refined gen zsh code, and gen shell codes
  - fea: support fish-shell completion generating now
  - fea: added root.RunAsSubCommand, treat 'app' as a synonym of 'app subcmd1 subcmd2'
  - imp/fix/fea: clarify the slice append or replace mode - SetNx & SetNxOverwrite for Option Store

- v1.10.13
  - fix/imp: adapt windir to *nix name to fit for generic config file loading
  - fea/imp: improved Aliases algor, support more tmpl var substitute
  - fix: fallback the unknown type as string type
  - fea: add flag to control whether write the changes back to alternative config file or not, `WithAlterConfigAutoWriteBack`
  - imp: name/desc fields of builtin commands and flags
  - CHANGE: use bgo build-tool now, Makefile thrown

- v1.10.11
  - fix: setNx bug at last commit
  - fix: send 1.10.10 failure

- v1.10.9
  - fix: setNx with slices merging
  - fix: aliases might be added to multiple groups
  - fea: secondary config file locations

- v1.10.8
  - fix/fea/imp: make cmdr aliases subsystem better

- v1.10.7
  - fix: generate shell may be lost buffered contents on writing to file

- v1.10.6
  - fix: internal commands and flags has wrong group declarations since last refactored.
  - fea: `-o file` for `generate shell` command.

- v1.10.5
  - fix: logex might crash on a nil skip field

- v1.10.3
  - last release failed because some deps cannot committed due to weak network

- v1.10.1
  - move to go1.17 to get a split declaration
  - fix: added the forgotten long-desc field
  - fix: transfer proper log-level to hedzr/log if in debug/trace mode
  - fix/imp: log.ForDir, ForFile
  - fix: log.LeftPad
  - fea: added InvokeCommand to run a sub-command from somewhere

- v1.10.0
  - fix: toggle-group key not sync while set via envvar
  - imp: speed up by extracting a re compiling code
  - imp: upgrade deps with more enh-helpers from [hedzr/log](https://github.com/hedzr/log)
  - imp: yaml indent size
  - imp: StripOrderPrefix
  - imp/fix: sync debug/trace mode back to hedzr/log
  - fix: options after tail args (positional args) might be ignored
  - fix: ResetOptions not clean up internal hierarchy-list
  - fea: added `Checkpoints` on _Option Store_  
    you may save and restore multiple checkpoints for cmdr _Option Store_, so that some temporary changes can be made.
  - fix/imp: `--man` crashes if manpages not installed - the responding manpage will be generated temporarily and instantly now
  - add `GitSummary` field into conf package
  - imp: speed up by reduce get worker
    - centralize rxxtOptions to store()
    - flatten backtrace(Flg|Cmd)Names, added dottedPathToCommand
  - **NOTE**: _the phrase wrapped by backtick(````) in `Description` field will be extracted as `DefaultValuePlaceholder` field automatically, so **beware** this feature._
  - fea: `-vv` (dup `-v` more than once) will print the hidden commands & flags in help screen NOW.  
    To take a sight of running `fluent generate --help --verbose -verbose`.
  - ...

- More details at [CHANGELOG](./CHANGELOG)

## Features

[Features.md](old/Features.md)

> Old README.md: [README.old.md](old/README.old.md)

## For Developer

[For Developer](old/Developer.md)

### Import cmdr

With go-modules enabled:

```go
import "github.com/hedzr/cmdr"
```

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