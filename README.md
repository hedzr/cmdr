# cmdr

[![Build Status](https://travis-ci.org/hedzr/cmdr.svg?branch=master)](https://travis-ci.org/hedzr/cmdr)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/cmdr)](https://goreportcard.com/report/github.com/hedzr/cmdr)
![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/cmdr.svg?label=release)

getopt/getopt_long like command-line UI golang library.

A getopt-like parser of command-line options, compatible with the [getopt_long](http://www.gnu.org/s/libc/manual/html_node/Argument-Syntax.html#Argument-Syntax) syntax, which is an extension of the syntax recommended by POSIX.

`cmdr` is a UNIX command-line UI library written by golang.


## Features

- Unix [*getopt*(3)](http://man7.org/linux/man-pages/man3/getopt.3.html) representation but without its programmatic interface.

- Automatic help screen generation

- Support for unlimited multiple sub-commands.

- Support for command short and long name, and aliases names.

- Support for both short and long options (`-o` and `--opt`). Support for multiple aliases

- Automatically allows those formats (applied to long flags too):
  - `-I file`, `-Ifile`, and `-I=files`
  - `-I 'file'`, `-I'file'`, and `-I='files'`
  - `-I "file"`, `-I"file"`, and `-I="files"`

- Support for `-D+`, `-D-` to enable/disable a bool flag.

- Support for **PassThrough** by `--`.

- Support for options being specified multiple times, with different values

- Support for optional arguments.

- Groupable commands and options/flags.

  Sortable group name with `[0-9A-Za-z]+\..+` format, eg:

  `1001.c++`, `1100.golang`, `1200.java`, …;

  `abcd.c++`, `b999.golang`, `zzzz.java`, …;

- Sortable commands and options/flags. Or sorted by alphabetic order.

- Predefined commands and flags:

  - Help: `-h`, `-?`, `--help`, ...
  - Version & Build Info: `--version`/`-V`, `--build-info`/`-#`
  - Verbose & Debug: `—verbose`/`-v`, `—debug`/`-D`, `—quiet`/`-q`
  - Generate Commands:
    - `generate shell`: `—bash`/`—zsh`(*todo*)
    - `generate manual`: *todo*
    - `generate doc`: *todo*

- Generators

  - *Todo: ~~manual generator~~, and ~~markdown~~/docx/pdf generators.*

  - Man Page generator: `bin/demo generate man`

  - Markdown generator: `bin/demo generate [doc|mdk|markdown]`

  - Bash and Zsh (*not yet, todo*) completion.

     ```bash
     bin/wget-demo generate shell --bash
     ```

- Predefined yaml config file locations:
  - `/etc/<appname>/<appname>.yml` and `conf.d` sub-directory.

  - `/usr/local/etc/<appname>/<appname>.yml` and `conf.d` sub-directory.

  - `$HOME/<appname>/<appname>,yml` and `conf.d` sub-directory.

  - Watch `conf.d` directory:
    - `AddOnConfigLoadedListener(c)`
    - `RemoveOnConfigLoadedListener(c)`
    - `SetOnConfigLoadedListener(c, enabled)`

  - As a feature, do NOT watch the changes on `<appname>.yml`.

  - To customize the searching locations yourself:

    - `SetPredefinedLocations(locations)`

      ```go
      SetPredefinedLocations([]string{"./config", "~/.config/cmdr/", "$GOPATH/running-configs/cmdr"})
      ```

- Overrides by environment variables.

  *todo: prior:*

- `cmdr.Get(key)`, `cmdr.GetBool(key)`, `cmdr.GetInt(key)`, `cmdr.GetString(key)`, `cmdr.GetStringSlice(key)` for Option value extractions.

  - bool
  - int, int64, uint, uint64
  - string
  - string slice
  - int slice
  - time duration
  - *todo: time, ~~duration~~, ~~int slice~~, ...*

- `cmdr.Set(key, value)`, `cmdr.SerNx(key, value)`

- Customizable `Painter` interface to loop each commands and flags.

- Uses `WalkAllCommands(walker)` to loop commands.

- Daemon

  > Uses daemon feature with `go-daemon`

  ```golang
  import "github.com/hedzr/cmdr/plugin/daemon"
  func main() {
  	daemon.Enable(NewDaemon())
  	if err := cmdr.Exec(rootCmd); err != nil {
  		log.Fatal("Error:", err)
  	}
  }
  func NewDaemon() daemon.Daemon {
  	return &DaemonImpl{}
  }
  ```

  See full codes in [demo](./examples/demo/) app.

  ```bash
  bin/demo server [start|stop|status|restart|install|uninstall]
  ```

  `install`/`uninstall` sub-commands could install `demo` app as a systemd service.

- `ExecWith(rootCmd *RootCommand, beforeXrefBuilding_, afterXrefBuilt_ HookXrefFunc) (err error)`

  `AddOnBeforeXrefBuilding(cb)`

  `AddOnAfterXrefBuilt(cb)`

- More...



## Examples

1. [**short**](./examples/short/README.md)  
   simple codes.
2. [demo](./examples/demo/README.md)  
   normal demo with external config files.
3. [wget-demo](./examples/wget-demo/README.md)  
   partial-impl wget demo.
   ![image](https://user-images.githubusercontent.com/12786150/58327052-29386500-7e61-11e9-8cd6-372aa1f14bfa.png)



## Documentation

- [*TODO: wiki*](https://github.com/hedzr/cmdr/wiki)



## Uses

- https://github.com/hedzr/consul-tags
- https://github.com/hedzr/ini-op
- Issue me



## Contrib

*Feel free to issue me bug reports and fixes. Many thanks to all contributors.*



## License

MIT.





