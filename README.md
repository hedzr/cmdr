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

`cmdr` is a POSIX-compliant,  command-line UI (CLI) library in Golang.
It is a getopt-like parser of command-line options, 
be compatible with the [getopt_long](http://www.gnu.org/s/libc/manual/html_node/Argument-Syntax.html#Argument-Syntax) 
command line UI, which is an extension of the syntax recommended by POSIX.

We made many enhancements beyond the standard library `flag`.

There is a fully-functional `Options Store` (configurations) for your
hierarchical configuration dataset too.

The .netCore version [Cmdr.Core](https
://github.com/hedzr/Cmdr.Core) is available now. And, a cxx version [`cmdr-cxx`](https://github.com/hedzr/cmdr-cxx) was pre-released just now (Happy Spring Festival 2021).

![ee99d078e2f7](https://user-images.githubusercontent.com/12786150/72876202-f49ee500-3d30-11ea-9de0-434bf8decf90.gif)
<!-- built by https://ezgif.com/ -->
> To review the image frames, go surfing at <https://github.com/hedzr/cmdr/issues/1#issuecomment-567779978>


## Table of Contents

* [cmdr](#cmdr)
  * [Table of Contents](#table-of-contents)
  * [Import](#import)
  * [News](#news)
  * [Features](#features)
  * [For Developer](#for-developer)
  * [Examples](#examples)
  * [Uses](#uses)
  * [Contrib](#contrib)
  * [Thanks to JODL](#thanks-to-jodl)
  * [License](#license)

<!-- Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc.go) -->

> [Youtube - 李宗盛2013最新單曲 山丘 官方完整版音檔](https://www.youtube.com/watch?v=_qNpR1Ew5jA) / [Jonathan Lee - Hill *CHT + ENU*](https://www.youtube.com/watch?v=FNlFe8ftBh0)




## Import

The better choice is importing with go-modules enabled:

```go
import "github.com/hedzr/cmdr"
```

See our extras:

- [**cmdr-docs**](https://github.com/hedzr/cmdr-docs): documentations (Working)
- [**cmdr-addons**](https://github.com/hedzr/cmdr-addons): a new daemon plugin `dex` for linux/macOS/windows.
- [**cmdr-examples**](https://github.com/hedzr/cmdr-examples): collects the samples for cmdr

and Bonus of #cmdr Series:

- dotnetCore: [Cmdr.Core](https://github.com/hedzr/Cmdr.Core)
- C++17 or higher: [cmdr-cxx](https://github.com/hedzr/cmdr-cxx)


## News


- docs (WIP):
  - english documentation NOT completed yet
  - documentation at: https://hedzr.github.io/cmdr-docs/

- v1.9.0 (WIP)
  - .fossa.yaml so a pre-release scan can be launched locally
  - BREAK: remove plugin/daemon - use cmdr-addons/pkg/plugins/dex instead

- v1.8.7
  - updated `log`, added: AutoStart Peripheral interface

- v1.8.5
  - updated `log`, fixed: forwarding systemd log to file

- v1.8.2
  - compliant with plan9,bsd,...
  - some data racing PRBs in parallel testing

- v1.8.1
  - fixed the CI error by imported from `log`
  - fixed a data racing in config files watching
  - small imp: pprof - added validArgs for cmdr-opt `profiling-type`
  - update deps: log & logex(logrus indirect, ...)

- v1.8.0
  - BREAK: removed support to golang 1.11 and below
  - update deps: log & logex(logrus indirect, ...)

- v1.7.46
  - added `plugin/pprof` package to simplify pprof integration
  - slight improvements

- More details at [CHANGELOG](./CHANGELOG)



## Features

[Features.md](old/Features.md)

> Old README.md: [README.old.md](old/README.old.md)


## About the Docker build

Here is a docker build for cmdr/examples/fluent so that you can run it without go building or downloading the release files:

```bash
# from Docker Hub:
$ docker run -it --rm hedzr/cmdr-fluent
$ docker run -it --rm hedzr/cmdr-fluent --help

# from Github Packages (please following the guide of GitHub Packages Site):
$ docker run -it --rm docker.pkg.github.com/hedzr/cmdr/cmdr-fluent
$ docker run -it --rm docker.pkg.github.com/hedzr/cmdr/cmdr-fluent --help
```


## For Developer

[For Developer](old/Developer.md)



## Examples

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
   


**See Also the examples index: [Examples.md](old/Examples.md)** *(zh-cn TODO: [Examples.zh-cn.md](old/Examples.zh-cn.md))*




## Uses

- https://github.com/hedzr/consul-tags
- https://github.com/hedzr/ini-op
- https://github.com/hedzr/awesome-tool
- austr
- Issue me to adding yours



## Contrib

*Feel free to issue me bug reports and fixes. Many thanks to all contributors.*


## Thanks to JODL

Thanks to [JetBrains](https://www.jetbrains.com/?from=cmdr) for donating product licenses to help develop **cmdr** [![jetbrains](https://gist.githubusercontent.com/hedzr/447849cb44138885e75fe46f1e35b4a0/raw/bedfe6923510405ade4c034c5c5085487532dee4/jetbrains-variant-4.svg)](https://www.jetbrains.com/?from=hedzr/cmdr)

[![goland](https://gist.githubusercontent.com/hedzr/447849cb44138885e75fe46f1e35b4a0/raw/ca8ac2694906f5650d585263dbabfda52072f707/logo-goland.svg)](https://www.jetbrains.com/?from=hedzr/cmdr)


## License

MIT


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr?ref=badge_large)
