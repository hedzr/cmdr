# cmdr

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/cmdr) 
[![Build Status](https://travis-ci.org/hedzr/cmdr.svg?branch=master)](https://travis-ci.org/hedzr/cmdr)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/cmdr)](https://goreportcard.com/report/github.com/hedzr/cmdr)
![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/cmdr.svg?label=release)

getopt/getopt_long like command-line UI (CLI) golang library. A replacer for go flags.

`cmdr` is a UNIX/Linux/POSIX command-line UI library written by golang. It is a getopt-like parser of command-line options, compatible with the [getopt_long](http://www.gnu.org/s/libc/manual/html_node/Argument-Syntax.html#Argument-Syntax) command line UI, which is an extension of the syntax recommended by POSIX.



> Before our introduce for CMDR, I decided to post A memory of mine, in each new OSS project. Here it is:
>
> [Youtube - 李宗盛2013最新單曲 山丘 官方完整版音檔](https://www.youtube.com/watch?v=_qNpR1Ew5jA) / [Jonathan Lee - Hill *CHT + ENU*](https://www.youtube.com/watch?v=FNlFe8ftBh0)
>
> Thanks.



![image](https://user-images.githubusercontent.com/12786150/58327052-29386500-7e61-11e9-8cd6-372aa1f14bfa.png)




## Features

- Unix [*getopt*(3)](http://man7.org/linux/man-pages/man3/getopt.3.html) representation but without its programmatic interface.

  - Options with short names (`-h`)
  - Options with long names (`--help`)
  - Options with aliases (`—helpme`, `—usage`, `--info`)
  - Options with and without arguments (bool v.s. other type)
  - Options with optional arguments and default values
  - Multiple option groups each containing a set of options
  - Supports multiple short options -aux
  - Supports namespaces for (nested) option groups

- Automatic help screen generation (*Generate and print well-formatted help message*)

- Strict Mode

  - *false*: Ignoring unknown command line options (default)

  - *true*: Report error on unknown commands and options if strict mode enabled (optional)

    enable strict mode:

    - env var `APP_STRICT_MODE=true`

    - hidden option: `--strict-mode` (if `cmdr.EnableCmdrCommands == true`)

    - entry in config file:

      ```yaml
      app:
        strict-mode: true
      ```

- Support for unlimited multiple sub-commands.

- Supports `-I/usr/include`` -I=/usr/include` `-I /usr/include` option argument specification

  Automatically allows those formats (applied to long option too):

  - `-I file`, `-Ifile`, and `-I=files`
  - `-I 'file'`, `-I'file'`, and `-I='files'`
  - `-I "file"`, `-I"file"`, and `-I="files"`

- Support for `-D+`, `-D-` to enable/disable a bool option.

- Support for **PassThrough** by `--`. (*Passing remaining command line arguments after -- (optional)*)

- Support for options being specified multiple times, with different values

- Groupable commands and options/flags.

  Sortable group name with `[0-9A-Za-z]+\..+` format, eg:

  `1001.c++`, `1100.golang`, `1200.java`, …;

  `abcd.c++`, `b999.golang`, `zzzz.java`, …;

- Sortable commands and options/flags. Or sorted by alphabetic order.

- Predefined commands and flags:

  - Help: `-h`, `-?`, `--help`, ...
  - Version & Build Info: `--version`/`-V`, `--build-info`/`-#`
  - Verbose & Debug: `—verbose`/`-v`, `—debug`/`-D`, `—quiet`/`-q`
  - `--no-env-overrides`, and `--strict-mode`
  - Generate Commands:
    - `generate shell`: `—bash`/`—zsh`(*todo*)/`--auto`
    - `generate manual`:  man 1 ready.
    - `generate doc`: markdown ready.

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

  - `$HOME/<appname>/<appname>.yml` and `conf.d` sub-directory.

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

  - supports configuration file formats:

    - Yaml
    - JSON
    - TOML

- Overrides by environment variables.

  *prior:* `defaultValue -> config-file -> env-var -> command-line opts`

- `cmdr.Get(key)`, `cmdr.GetBool(key)`, `cmdr.GetInt(key)`, `cmdr.GetString(key)`, `cmdr.GetStringSlice(key)` and `cmdr.GetIntSlice(key)` for Option value extractions.

  - bool
  - int, int64, uint, uint64
  - string
  - string slice
  - int slice
  - time duration
  - *todo: float, time, ~~duration~~, ~~int slice~~, ...*
  - *todo: all primitive go types*
  - *todo: maps*

- `cmdr.GetP(prefix, key)`, `cmdr.GetBoolP(prefix, key)`, ….

- `cmdr.Set(key, value)`, `cmdr.SerNx(key, value)`

  `Set()` set value by key without RxxtPrefix, eg: `cmdr.Set("debug", true)` for `--debug`.

  `SetNx()` set value by exact key. so: `cmdr.SetNx("app.debug", true)` for `--debug`.

- Customizable `Painter` interface to loop each commands and flags.

- Uses `WalkAllCommands(walker)` to loop commands.

- Daemon (*Linux Only*)

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

  > Just For Linux

- `ExecWith(rootCmd *RootCommand, beforeXrefBuilding_, afterXrefBuilt_ HookXrefFunc) (err error)`

  `AddOnBeforeXrefBuilding(cb)`

  `AddOnAfterXrefBuilt(cb)`

- Launch external editor by `&Flag{BaseOpt:BaseOpt{},ExternalTool:cmdr.ExternalToolEditor}`:

  just like `git -m`, try this command:

  ```bash
  EDITOR=nano bin/demo -m ~~debug
  ```

  Default is `vim`. And `-m "something"` can skip the launching.

- `ToggleGroup`: make a group of flags as a radio-button group.

- Muiltiple API styles:

  - Data Definitions style (Classical style): see also [root_cmd.go in demo](https://github.com/hedzr/cmdr/blob/master/examples/demo/demo/root_cmd.go)
  - Fluent API style: see also [main.go in fluent](https://github.com/hedzr/cmdr/blob/master/examples/fluent/main.go)

- More...



## Examples

1. [**short**](./examples/short/README.md)  
   simple codes.
2. [demo](./examples/demo/README.md)  
   normal demo with external config files.
3. [wget-demo](./examples/wget-demo/README.md)  
   partial-impl wget demo.
4. [fluent](./examples/fluent)  
   fluent api demo.



## Documentation

- [*TODO: wiki*](https://github.com/hedzr/cmdr/wiki)



### Uses Fluent API

```go
	root := cmdr.Root("aa", "1.0.1").Header("aa - test for cmdr - hedzr")
	rootCmd = root.RootCommand()

	co := root.NewSubCommand().
		Titles("ms", "micro-service").
		Description("", "").
		Group("")

	co.NewFlag(cmdr.OptFlagTypeUint).
		Titles("t", "retry").
		Description("", "").
		Group("").
		DefaultValue(3, "RETRY")

	cTags := co.NewSubCommand().
		Titles("t", "tags").
		Description("", "").
		Group("")

	cTags.NewFlag(cmdr.OptFlagTypeString).
		Titles("a", "addr").
		Description("", "").
		Group("").
		DefaultValue("consul.ops.local", "ADDR")

	cTags.NewSubCommand().
		Titles("ls", "list").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})

	cTags.NewSubCommand().
		Titles("a", "add").
		Description("", "").
		Group("").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			return
		})
```





## Uses

- https://github.com/hedzr/consul-tags
- https://github.com/hedzr/ini-op
- voxr
- austr
- Issue me to adding yours



## Contrib

*Feel free to issue me bug reports and fixes. Many thanks to all contributors.*



## License

MIT.





