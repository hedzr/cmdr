# cmdr

[![Build Status](https://travis-ci.org/hedzr/cmdr.svg?branch=master)](https://travis-ci.org/hedzr/cmdr)
![Go](https://github.com/hedzr/cmdr/workflows/Go/badge.svg)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/cmdr.svg?label=release)](https://github.com/hedzr/cmdr/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/cmdr) [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr?ref=badge_shield)

[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/cmdr)](https://goreportcard.com/report/github.com/hedzr/cmdr)
[![codecov](https://codecov.io/gh/hedzr/cmdr/branch/master/graph/badge.svg)](https://codecov.io/gh/hedzr/cmdr)
[![Coverage Status](https://coveralls.io/repos/github/hedzr/cmdr/badge.svg?branch=master)](https://coveralls.io/github/hedzr/cmdr?branch=master)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#command-line)

`cmdr` is a POSIX/GNU style,  command-line UI (CLI) Go library.
It is a getopt-like parser of command-line options, 
be compatible with the [getopt_long](http://www.gnu.org/s/libc/manual/html_node/Argument-Syntax.html#Argument-Syntax) 
command line UI, which is an extension of the syntax recommended
by POSIX.

There are couples of enhancements beyond the standard 
library `flag`.

There is a full `Options Store` (configurations) for your
hierarchy configuration data too.


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




[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Fcmdr?ref=badge_large)

## Import

Need go-modules enabled:

```go
import "github.com/hedzr/cmdr"
```

See our extras:

- [cmdr-addons](https://github.com/hedzr/cmdr-addons): a new daemon plugin `dex` for linux/macOS/windows.
- [cmdr-examples](https://github.com/hedzr/cmdr-examples)



## News

- v1.6.48
  - code reviewed
  - maintained
  - unnecessary deps removed.
  - small fixes

- v1.6.47
  - fixed/improved: reset slice value if an empty slice was been setting
  - improved: add logging output in delve debugging mode
  - fixed: matching the longest short flag for combining flags
  - **BROKEN API**: the param `defaultValue` is optional now: cmdr.NewBool(), cmdr.NewInt(), ...
  - added `cmdr.NewUintSlice()`

- v1.6.45
  - fixed/improved: `ToBool(value, defval...) bool`
  - fixed: flag.OnSet trigger for envvar hit
  - fixed/improved: friendly error msg

- v1.6.43
  - fixed/improved: the matching algorithm and remained args

- v1.6.41
  - `WithPagerEnabled()`: enables OS pager for help screen output

- v1.6.39
  - **BROKEN API**: the params order exchanged, their new prototypes are `OptFlag.Titles(long, short, aliases)` and `OptCmd.Titles(long, short, alases)`.
  - improved help screen
  - bug fixed:
    - the value of remained args could be wrong sometimes
    - stop flag split in parsing
    - some coverage test errors

- v1.6.36
  - `ToggleGroup`:
    - assume the empty Group field with ToggleGroup
    - set "command-path.toggleGroupName" to the hit flag full name as flipping a toggle-group.  
      For example, supposed a toggle-group 'KK' under 'server' command with 3 choices/flags: apple, banana, orange. For the input '--orange', these entries will be set in option store:  
      `server.orange` <== true;  
      `server.KK` <== 'orange';  
  - fixed: `GetStringSliceXxx()` return the value array without expand the envvar.
  - improved: some supports for plan9
  - fixed: can't expand envvar correectly at earlier initializing.

- For more information to refer to [CHANGELOG](./CHANGELOG)



## Features

[Features.md](./Features.md)

> Old README.md: [README.old.md](./README.old.md)


## For Developer

[For Developer](Developer.md)



## Examples

1. [**short**](./examples/short/README.md)  
   simple codes with structured data style.

2. [demo](./examples/demo/README.md)  
   normal demo with external config files.

3. [wget-demo](./examples/wget-demo/README.md)  
   partial-covered for GNU `wget`.

4. [fluent](./examples/fluent)  
   demostrates how to define your command-ui with the fluent api style.

5. [ffmain](./examples/ffmain)

   a demo to show you how to migrate from go `flag` smoothly.

6. [cmdr-http2](https://github.com/hedzr/cmdr-http2)  
   http2 server with daemon supports, graceful shutdown

7. [awesome-tool](https://github.com/hedzr/awesome-tool)  
   `awesome-tool` is a cli app that fetch the repo stars and generate a markdown summary, accordingly with most of awesome-xxx list in github (such as awesome-go).
   


**See Also the examples index: [Examples.md](./Examples.md)** *(zh-cn TODO: [Examples.zh-cn.md](./Examples.zh-cn.md))*




## Uses

- https://github.com/hedzr/consul-tags
- https://github.com/hedzr/ini-op
- https://github.com/hedzr/awesome-tool
- austr
- Issue me to adding yours



## Contrib

*Feel free to issue me bug reports and fixes. Many thanks to all contributors.*


## Thanks to JODL

[JODL (JetBrains OpenSource Development License)](https://www.jetbrains.com/community/opensource/) is good:

[![goland](https://gist.githubusercontent.com/hedzr/447849cb44138885e75fe46f1e35b4a0/raw/ca8ac2694906f5650d585263dbabfda52072f707/logo-goland.svg)](https://www.jetbrains.com/?from=hedzr/cmdr)
[![jetbrains](https://gist.githubusercontent.com/hedzr/447849cb44138885e75fe46f1e35b4a0/raw/bedfe6923510405ade4c034c5c5085487532dee4/jetbrains-variant-4.svg)](https://www.jetbrains.com/?from=hedzr/cmdr)


## License

MIT