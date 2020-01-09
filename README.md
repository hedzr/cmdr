# cmdr

[![Build Status](https://travis-ci.org/hedzr/cmdr.svg?branch=master)](https://travis-ci.org/hedzr/cmdr)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/cmdr.svg?label=release)](https://github.com/hedzr/cmdr/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/cmdr) 
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/cmdr)](https://goreportcard.com/report/github.com/hedzr/cmdr)
[![codecov](https://codecov.io/gh/hedzr/cmdr/branch/master/graph/badge.svg)](https://codecov.io/gh/hedzr/cmdr)
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


**| [News](#news) | [Features](#features) |
[Examples](#examples) | 
[Docs: For Developer](#for-developer) |
[Docs: Uses Fluent API](#uses-fluent-api) | [At Playground](#at-playground) |**


> [Youtube - 李宗盛2013最新單曲 山丘 官方完整版音檔](https://www.youtube.com/watch?v=_qNpR1Ew5jA) / [Jonathan Lee - Hill *CHT + ENU*](https://www.youtube.com/watch?v=FNlFe8ftBh0)

<!-- ![image](https://user-images.githubusercontent.com/12786150/58327052-29386500-7e61-11e9-8cd6-372aa1f14bfa.png) -->
<!-- ![image](https://user-images.githubusercontent.com/12786150/71229810-f46e8c80-2321-11ea-8c0d-d1952f47dad3.png) -->
![gif](https://user-images.githubusercontent.com/12786150/71230660-7f04bb00-2325-11ea-8662-673839a968c8.gif)
<!-- built by https://gifmaker.me/ -->
> To review the image frames, go surfing at <https://github.com/hedzr/cmdr/issues/1#issuecomment-567779978>

## Import

for non-go-modules user:

```go
import "gopkg.in/hedzr/cmdr.v1"
```

with go-modules enabled:

```go
import "github.com/hedzr/cmdr"
```



## News

- v1.6.20
  - fixed build error
  
- v1.6.19
  - added: `GetKibibytesXXX`, `GetKilobytesXXX`. with further info, run `go run ./examples/fluent kb --size --help`.
  - added hot-reload feature for daemon plugin:
    `bin/cmdr-http2 server [hot-restart|hot-relaod]` will restart the app, transfer the living listener to the child process without break it.
    `SIGUSR2` has the same effect.
    > - `Daemon.OnRun()` broken.
    > - daemon plugin can't be used for non-linux platform.
  - added: `Commnad.GetDottedNamePath`
    
- v1.6.18
  - added: `WithOnSwitchCharHit`, `WithOnPassThruCharHit`
  - bug fixed
  - reviewing and refactoring

- For more information to refer to [CHANGELOG](./CHANGELOG)




## Features

cmdr has rich features:

- builds multi-level command and sub-commands
- builds short, long and alias options with kinds of data types
- defines commands and options via fluent api style
- or defines its with enhanced stdlib `flag` style
- full featured `Options Store` for hosted any application configurations
  - watchable external config file and child directory `conf.d`
  - watchable option value merging event: while option value modified in external config file and loaded automatically.
  - watchable option value modifying event: while option value modified (from config file, or programmatically)
  - connectable with external configuration center

### More

- Unix [*getopt*(3)](http://man7.org/linux/man-pages/man3/getopt.3.html) representation but without its programmatic interface.

  - Options with short names (`-h`)
  - Options with long names (`--help`)
  - Options with aliases (`--helpme`, `--usage`, `--info`)
  - Options with and without arguments (bool v.s. other type)
  - Options with optional arguments and default values
  - Multiple option groups each containing a set of options
  - Supports the compat short options `-aux` == `-a -u -x`
  - Supports namespaces for (nested) option groups

- Automatic help screen generation (*Generates and prints well-formatted help message*)

- Supports the Fluent API style
  
  <details><summary>Sample codes</summary>
  
  ```go
  root := cmdr.Root("aa", "1.0.3")
      // Or  // .Copyright("All rights reserved", "sombody@example.com")
      .Header("aa - test for cmdr - hedzr")
  rootCmd = root.RootCommand()
  
  co := root.NewSubCommand().
  	Titles("ms", "micro-service").
  	Description("", "").
  	Group("")
  
  // deprecated since v1.6.9
  // co.NewFlag(cmdr.OptFlagTypeUint).
  //  	Titles("t", "retry").
  // 	Description("", "").
  // 	Group("").
  // 	DefaultValue(3, "RETRY")

  co.NewFlagV(3).
  	Titles("t", "retry").
  	Description("", "").
  	Group("").
  	Palceholder("RETRY")
  
  cTags := co.NewSubCommand().
  	Titles("t", "tags").
  	Description("", "").
  	Group("")
  ```
  
  </details>
  
- Muiltiple API styles:

  - Data Definitions style (Classical style): see also [root_cmd.go in demo](https://github.com/hedzr/cmdr/blob/master/examples/demo/demo/root_cmd.go)
  - Fluent API style: see also [main.go in fluent](https://github.com/hedzr/cmdr/blob/master/examples/fluent/main.go)
  - **go `flag`-like API style**: see also [main.go in ffmain](https://github.com/hedzr/cmdr/blob/master/examples/ffmain/main.go)

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

- Supports for unlimited multi-level sub-commands.

- Supports `-I/usr/include -I=/usr/include` `-I /usr/include` option argument specifications
  
  Automatically allows those formats (applied to long option too):

  - `-I file`, `-Ifile`, and `-I=files`
  - `-I 'file'`, `-I'file'`, and `-I='files'`
  - `-I "file"`, `-I"file"`, and `-I="files"`

- Supports for `-D+`, `-D-` to enable/disable a bool option.

- Supports for **PassThrough** by `--`. (*Passing remaining command line arguments after -- (optional)*)

- Supports for options being specified multiple times, with different values

  > since v1.5.0:
  >
  > - and multiple flags `-vvv` == `-v -v -v`, then `cmdr.FindFlagRecursive("verbose", nil).GetTriggeredTimes()` should be `3`
  >
  > - for bool, string, int, ... flags, last one will be kept and others abandoned:
  >
  >   `-t 1 -t 2 -t3` == `-t 3`
  >
  > - for slice flags, all of its will be merged (NOTE that duplicated entries are as is):
  >
  >   slice flag overlapped
  >
  >   - `--root A --root B,C,D --root Z,A` == `--root A,B,C,D,Z`
  >     cmdr.GetStringSliceR("root") will return `[]string{"A","B","C","D","Z"}`

- Smart suggestions for wrong command and flags

  since v1.1.3, using [Jaro-Winkler distance](https://en.wikipedia.org/wiki/Jaro%E2%80%93Winkler_distance) instead of soundex.

- Groupable commands and options/flags.

  Sortable group name with `[0-9A-Za-z]+\..+` format, eg:

  - `1001.c++`, `1100.golang`, `1200.java`, …;
  - `abcd.c++`, `b999.golang`, `zzzz.java`, …;

- Sortable commands and options/flags. Or sorted by alphabetic order.

- Predefined commands and flags:

  - Help: `-h`, `-?`, `--help`, `--info`, `--usage`, `--helpme`, ...
  - Version & Build Info: `--version`/`--ver`/`-V`, `--build-info`/`-#`
    - Simulating version at runtime with `—version-sim 1.9.1`
    - generally, `conf.AppName` and `conf.Version` are originally.
    - `~~tree`: list all commands and sub-commands.
    - `--config <location>`: specify the location of the root config file.
  - Verbose & Debug: `—verbose`/`-v`, `—debug`/`-D`, `—quiet`/`-q`
  - Generate Commands:
    - `generate shell`: `—bash`/`—zsh`(*todo*)/`--auto`
    - `generate manual`:  man 1 ready.
    - `generate doc`: markdown ready.
  - `cmdr` Specials:
    - `--no-env-overrides`, and `--strict-mode`
    - `--no-color`: print the plain text to console without ANSI colors.

- Generators
 
  - *Todo: ~~manual generator~~, and ~~markdown~~/docx/pdf generators.*

  - Man Page generator: `bin/demo generate man`

  - Markdown generator: `bin/demo generate [doc|mdk|markdown]`

  - Bash and Zsh (*not yet, todo*) completion.

     ```bash
     $ bin/wget-demo generate shell --bash
     ```

- Predefined external config file locations:
  - `/etc/<appname>/<appname>.yml` and `conf.d` sub-directory.
  - `/usr/local/etc/<appname>/<appname>.yml` and `conf.d` sub-directory.
  - `$HOME/.config/<appname>/<appname>.yml` and `conf.d` sub-directory.
  - `$HOME/.<appname>/<appname>.yml` and `conf.d` sub-directory.
  - all predefined locations are:
  
    ```go
    predefinedLocations: []string{
  		"./ci/etc/%s/%s.yml",       // for developer
    	"/etc/%s/%s.yml",           // regular location: /etc/$APPNAME/$APPNAME.yml
  		"/usr/local/etc/%s/%s.yml", // regular macOS HomeBrew location
    	"$HOME/.config/%s/%s.yml",  // per user: $HOME/.config/$APPNAME/$APPNAME.yml
  		"$HOME/.%s/%s.yml",         // ext location per user
    	"$THIS/%s.yml",             // executable's directory
  		"%s.yml",                   // current directory
    },
    ```
    
  - since v1.5.0, uses `cmdr.WithPredefinedLocations("a","b",...),`
  
- Watch `conf.d` directory:
  
  - `cmdr.WithConfigLoadedListener(listener)`
    
    - `AddOnConfigLoadedListener(c)`
    - `RemoveOnConfigLoadedListener(c)`  
    - `SetOnConfigLoadedListener(c, enabled)`
    
  - As a feature, do NOT watch the changes on `<appname>.yml`.
  
    - *since v1.6.9*, `WithWatchMainConfigFileToo(true)` allows the main config file `<appname>.yml`  to be watched.

  - on command-line:
  
    ```bash
    $ bin/demo --configci/etc/demo-yy ~~debug
    $ bin/demo --config=ci/etc/demo-yy/any.yml ~~debug
    $ bin/demo --config ci/etc/demo-yy/any.yml ~~debug
    ```
  
  - supports muiltiple file formats:
  
    - Yaml
    - JSON
    - TOML
  
  - `cmdr.Exec(root, cmdr.WithNoLoadConfigFiles(false))`: disable loading external config files.
  
- Overrides by environment variables.

  *priority level:* `defaultValue -> config-file -> env-var -> command-line opts`

- `Option Store` - Unify option value extraction:

  - `cmdr.Get(key)`, `cmdr.GetBool(key)`, `cmdr.GetInt(key)`, `cmdr.GetString(key)`, `cmdr.GetStringSlice(key, defaultValues...)` and `cmdr.GetIntSlice(key, defaultValues...)`, `cmdr.GetDuration(key)` for Option value extractions.

    - bool
    - int, int64, uint, uint64, float32, float64
      ```bash
      $ app -t 1    #  float: 1.1, 1e10, hex: 0x9d, oct: 0700, bin: 0b00010010
      ```
    - string
    - string slice, int slice (comma-separated)
      ```bash
      $ app -t apple,banana      # => []string{"apple", "banana"}
      $ app -t apple -t banana   # => []string{"apple", "banana"}
      ```
    - time duration (1ns, 1ms, 1s, 1m, 1h, 1d, ...)
      ```bash
      $ app -t 1ns -t 1ms -t 1s -t 1m -t 1h -t 1d
      ```
    - ~~*todo: float, time, duration, int slice, …, all primitive go types*~~
    - map
    - struct: `cmdr.GetSectionFrom(sectionKeyPath, &holderStruct)`
    
  - `cmdr.Set(key, value)`, `cmdr.SerNx(key, value)`

    - `Set()` set value by key without RxxtPrefix, eg: `cmdr.Set("debug", true)` for `--debug`.
    - `SetNx()` set value by exact key. so: `cmdr.SetNx("app.debug", true)` for `--debug`.

  - Fast Guide for `Get`, `GetP` and `GetR`:

    - `cmdr.GetP(prefix, key)`, `cmdr.GetBoolP(prefix, key)`, ….
    - `cmdr.GetR(key)`, `cmdr.GetBoolR(key)`, …, `cmdr.GetMapR(key)`
    - `cmdr.GetRP(prefix, key)`, `cmdr.GetBoolRP(prefix, key)`, ….
  
    As a fact, `cmdr.Get("app.server.port")` == `cmdr.GetP("app.server", "port")` == `cmdr.GetR("server.port")`
	(*if cmdr.RxxtPrefix == ["app"]*); so:

    ```go
    cmdr.Set("server.port", 7100)
    assert cmdr.GetR("server.port") == 7100
    assert cmdr.Get("app.server.port") == 7100
    ```
    
    In most cases, **GetXxxR()** are recommended.
    
    While extracting string value, the evnvar will be expanded automatically but raw version `GetStringNoExpandXXX()` available since v1.6.7. For example:
    
    ```go
    fmt.Println(cmdr.GetStringNoExpandR("kk"))  // = $HOME/Downloads
    fmt.Println(cmdr.GetStringR("kk"))          // = /home/ubuntu/Downloads
    ``` 

- cmdr Options Store

  internal `rxxtOptions`

- Walkable

  - Customizable `Painter` interface to loop *each* command and flag.
  - Walks on all commands with `WalkAllCommands(walker)`.

- Daemon (*Linux Only*)

  > rewrote since v1.6.0

  ```go
  import "github.com/hedzr/cmdr/plugin/daemon"
  func main() {
  	if err := cmdr.Exec(rootCmd,
	    daemon.WithDaemon(NewDaemon(), nil,nil,nil),
		); err != nil {
  		log.Fatal("Error:", err)
  	}
  }
  func NewDaemon() daemon.Daemon {
  	return &DaemonImpl{}
  }
  ```

  See full codes in [demo](./examples/demo/) app, and [**cmdr-http2**](https://github.com/hedzr/cmdr-http2).

  ```bash
  $ bin/demo server [start|stop|status|restart|install|uninstall]
  ```

  `install`/`uninstall` sub-commands could install `demo` app as a systemd service.

  > Just for Linux

- ~~`ExecWith(rootCmd *RootCommand, beforeXrefBuilding_, afterXrefBuilt_ HookXrefFunc) (err error)`~~

  ~~`AddOnBeforeXrefBuilding(cb)`~~

  ~~`AddOnAfterXrefBuilt(cb)`~~

- `cmdr.WithXrefBuildingHooks(beforeXrefBuilding, afterXrefBuilding)`

- Debugging options:

  - `~~debug`: dump all key value pairs in parsed options store

    ```bash
    $ bin/demo -? ~~debug
    $ bin/demo -? ~~debug ~~raw  # without envvar expanding
    $ bin/demo -? ~~debug ~~env  # print envvar k-v pairs too
    $ bin/demo -? ~~debug --more
    ```
    
    `~~debug` depends on `--help` present (or invoking a command which have one ore more children)

  - `InDebugging()`, isdelve (refer to [here](https://stackoverflow.com/questions/47879070/how-can-i-see-if-the-goland-debugger-is-running-in-the-program/47890273#47890273)) - To use it, add `-tags=delve`:
  
    ```bash
    go build -tags=delve cli/main.go
    go run -tags=delve cli/main.go --help
    ```
  
- `~~tree`: dump all sub-commands
  
    ```bash
    $ bin/demo ~~tree
    ```

   `~~tree` is a special option/flag like a command.
  

- More Advanced features

  - Launches external editor by `&Flag{BaseOpt:BaseOpt{},ExternalTool:cmdr.ExternalToolEditor}`:

    just like `git -m`, try this command:

     ```bash
     $ EDITOR=nano bin/demo -m ~~debug
     ```

     Default is `vim`. And `-m "something"` can skip the launching.

  - `ToggleGroup`: make a group of flags as a radio-button group.

  - Safe password input for end-user: `cmdr.ExternalToolPasswordInput`

  - `head`-like option: treat `app do sth -1973` as `app do sth -a=1973`, just like `head -1`.

    ```go
    Flags: []*cmdr.Flag{
        {
            BaseOpt: cmdr.BaseOpt{
                Short:       "h",
                Full:        "head",
                Description: "head -1 like",
            },
            DefaultValue: 0,
            HeadLike:     true,
        },
    },
    ```

  - limitation with enumerable values:

    ```go
    Flags: []*cmdr.Flag{
        {
            BaseOpt: cmdr.BaseOpt{
                Short:       "e",
                Full:        "enum",
                Description: "enum tests",
            },
            DefaultValue: "", // "apple",
            ValidArgs:    []string{"apple", "banana", "orange"},
        },
    },
    ```

    While a non-in-list value found, An error (`*ErrorForCmdr`) will be thrown:

    ```go
    cmdr.ShouldIgnoreWrongEnumValue = true
    if err := cmdr.Exec(rootCmd); err != nil {
        if e, ok := err(*cmdr.ErrorForCmdr); ok {
            // e.Ignorable is a copy of [cmdr.ShouldIgnoreWrongEnumValue]
            if e.Ignorable {
                logrus.Warning("Non-recognaizable value found: ", e)
                os.Exit(0)
            }
        }
        logrus.Fatal(err)
    }
    ```

  - `cmdr.TrapSignals(fn, signals...)`

    It is a helper to simplify your infidonite loop before exit program:

    <details><summary>Sample codes</summary>
      Here is sample fragment:
      ```go
      func enteringLoop() {
     	  waiter := cmdr.TrapSignals(func(s os.Signal) {
     	    logrus.Debugf("receive signal '%v' in onTrapped()", s)
     	  })
     	  go waiter()
      }
      ```
    </details>


- More...




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

## Documentation

- [*TODO: wiki*](https://github.com/hedzr/cmdr/wiki)



### For Developer

To build and test `cmdr`:

```bash
$ make help   # see all available sub-targets
$ make info   # display building environment
$ make build  # build binary files for examples
$ make gocov  # test

# customizing
$ GOPROXY_CUSTOM=https://goproxy.io make info
$ GOPROXY_CUSTOM=https://goproxy.io make build
$ GOPROXY_CUSTOM=https://goproxy.io make gocov
```



### Uses Fluent API

<details>
	<summary> Expand to source codes </summary>

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

</details>



## At Playground

Try its out :

- https://play.golang.org/p/ieExm3V1Pcx
- wget-demo at playground: https://play.golang.org/p/wpEZgQGzKyt
- demo with daemon plugin: https://play.golang.org/p/wJUA59uGu2M


## Uses

- https://github.com/hedzr/consul-tags
- https://github.com/hedzr/ini-op
- https://github.com/hedzr/awesome-tool
- austr
- Issue me to adding yours



## Contrib

*Feel free to issue me bug reports and fixes. Many thanks to all contributors.*



## License

MIT.





