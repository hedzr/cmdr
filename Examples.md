# cmdr-examples

see also: https://github.com/hedzr/cmdr-examples





#### At Playground

Try its out :

- https://play.golang.org/p/ieExm3V1Pcx
- wget-demo at playground: https://play.golang.org/p/wpEZgQGzKyt
- demo with daemon plugin: https://play.golang.org/p/wJUA59uGu2M





#### Getting Start

- Simple Guide

  [`simple`](https://github.com/hedzr/cmdr-examples/tree/master/examples/simple)

- Make a Copy for `wget` (CLI only)

  `wget-demo` - use data definitions by nested struct

- Using Fluent API style

  [`simple`](https://github.com/hedzr/cmdr-examples/tree/master/examples/simple), [`flags`](https://github.com/hedzr/cmdr-examples/tree/master/examples/flags), ...

- Writing the web-server with daemon supports

  obseleted

- With daemonex supports (New)

  [`service`](https://github.com/hedzr/cmdr-examples/tree/master/examples/service)







#### Concepts

- What is POSIX for getopt and CLI?

- Why the Long/Full Name is so important in `cmdr`?

  To extract an option value from `Option Store`, we need a `keyPath` with the form `<subcommand1>.<subcommand2>...<subcommandn>.<flag>` .

  The `cmdr.GetXXX(keyPath)` functions will return the option value of a `keyPath` with the desired type. For example:

  ```go
  var hostPort string = cmdr.GetStringR("server.start.host")
  fmt.Println(hostPort)
  assert(hostPort == cmdr.GetString("app.server.start.host"))
  ```

  > `GetStringR(keyPath)` will wrap the envPrefixes onto `keyPath` and `GetString(keyPath)` not.

  

- Subcommands

  - how the nested commands works?
  - [`subcommands`](https://github.com/hedzr/cmdr-examples/tree/master/examples/subcommands)

- Subcommands categories

  - sorted group

- Flags: [`flags`](https://github.com/hedzr/cmdr-examples/tree/master/examples/flags), [`actions`](https://github.com/hedzr/cmdr-examples/tree/master/examples/actions)

- Options: flags and app settings

  - Placeholder values

  - Positional Arguments

    

  - Tail Args

    By using `cmdr`, it takes this form in command-line:

    ```bash
  app <sub-commands> <all options or ascestor of it> <tail args>
    ```
  
    The *tail arguments* normally indicates a set of filenames, words, and so on.  For example, for a command-line `c++-cc compile -v a.cc b.cc`, we have the tail args `a.cc b.cc`.

    For `cmdr`, 

    

  - Action: [`actions`](https://github.com/hedzr/cmdr-examples/tree/master/examples/actions)

    We can associate an *Action* with a command. Here is a sample:

    ```go
	root := cmdr.Root(appName, "1.0.1").
    		Copyright(copyright, "hedzr").
		Description(desc, longDesc).
    		Examples(examples)
    	// rootCmd = root.RootCommand()
    	soundex(root)
      // ...
    
    func soundex(root cmdr.OptCmd) {
    	// soundex
    
      // To test for this subcommand, type command in shell:
      // $ go run ./examples/simple soundex quick fox
    	root.NewSubCommand("soundex", "snd", "sndx", "sound").
    		Description("soundex test").
    		Group("Test").
    		TailPlaceholder("[text1, text2, ...]").
    		Action(func(cmd *cmdr.Command, remainArgs []string) (err error) {
    			for ix, s := range remainArgs {
    				fmt.Printf("%5d. %s => %s\n", ix, s, cmdr.Soundex(s))
    			}
    			return
    		})
    }
    ```
    
    An *Action* Handler is a callback function with prototype `func(cmd *cmdr.Command, remainArgs []string) (err error)`.  As to `soundex`, it loops for each `remainArgs` and  calculate its soundex value. For this shell command line:
    
    ```bash
$ go run ./examples/simple soundex quick fox
    ```

    The `remainArgs` are `quick` and `fox`.
    
    

  - Alternate names (aliases)

    - [`flags`](https://github.com/hedzr/cmdr-examples/tree/master/examples/flags)

  - Ordering

    - [`flags`](https://github.com/hedzr/cmdr-examples/tree/master/examples/flags)

  - Values from Environment

    - [`envvars`](https://github.com/hedzr/cmdr-examples/tree/master/examples/envvars)

  - Values from external config files (YAML/JSON/TOML, ...)

    - [`configfile`](https://github.com/hedzr/cmdr-examples/tree/master/examples/configfile)

  - Values from external config center: TODO

  - Precedence

    - The precedence for flag value sources is:
    1. Environment variable if specified
      2. Command line flag vlaue from user
    3. Config file(s) if valid
      4. Default defined on the flag

  - Combining short flags: [`action`](https://github.com/hedzr/cmdr-examples/tree/master/examples/actions), [`flags`](https://github.com/hedzr/cmdr-examples/tree/master/examples/flags)

  - Array and flags: [`flags`](https://github.com/hedzr/cmdr-examples/tree/master/examples/flags)
  
  - Head-like: [head-like](https://github.com/hedzr/cmdr-examples/tree/master/examples/head-like)
  
  - Toggleable - `ToogleGroup` : [toggle-group](https://github.com/hedzr/cmdr-examples/tree/master/examples/toggle-group)

  - Password input:

  - for TUI:

  - inline shell mode: [shell-mode](https://github.com/hedzr/cmdr-examples/tree/master/examples/shell-mode)

  - using progress bar

  - unmatched command or flag found: Smart suggestions

  - GetStrictMode, GetDebugMode, GetVerboseMode, GetQuietMode, GetNoColorMode

  - 

  - ...

- Actions: [`actions`](https://github.com/hedzr/cmdr-examples/tree/master/examples/actions)

  - The Right Point to Handling a (sub)command
  - Pre, Post
  - Exit Code
  - Unhandled Exception

- External config files/api: [`configfile`](https://github.com/hedzr/cmdr-examples/tree/master/examples/configfile)

- Using `Option Store`

  - GetBoolR
  - GetMap
  - Dump

- Using `--debug`

- Using `--tree`

- The common preset flags: `--help`, `--version`, `--build-info`, ...

- Help Screen and Customizing It

- Shell features

  - Auto-completion
  - manpage generator
  - ...

- Walking for all commands, and flags

- Migrating from golang `flag`

- Hooks in `cmdr`

- Fluent API Style

  - Why the struct tags style is not supported?

- Others

  - error handling and raising: [`panics`](https://github.com/hedzr/cmdr-examples/tree/master/examples/panics)

  - Show help screen with pager:

    you can pipe cmdr help screen to OS page (like `less`), with `cmdr.WithPagerEnabled()`:

    ```go
    	if err := cmdr.Exec(buildRootCmd(),
    		cmdr.WithPagerEnabled(),
    	); err != nil {
    		logrus.Fatalf("error: %+v", err)
    	}
    ```

    see also [`service`](https://github.com/hedzr/cmdr-examples/tree/master/examples/service) and take a look for:

    ```bash
    go run ./examples/service s run -?
    ```

    

  - 

- More

  - [`service`](https://github.com/hedzr/cmdr-examples/tree/master/examples/service): new daemon plugin extension [`dex`](https://github.com/hedzr/cmdr-addons/pkg/plugins/dex/)
  - `winsvc`
  - `non-cmdr`

- TUI

  - progress bar (just for VT100-compatible Terminal): [`progressbar`](https://github.com/hedzr/cmdr-examples/tree/master/examples/progressbar)

    make progress bar work with cmdr

  - shell mode: [`shell-mode`](https://github.com/hedzr/cmdr-examples/tree/master/examples/shell-mode)

    ```bash
    go run ./examples/shell-mode
    ```

    

  - Interactive prompt: [`interactive-prompt`](https://github.com/hedzr/cmdr-examples/tree/master/examples/interactive-prompt)



## Option Store



- Kibibytes:

  Extracting `Kibibyte` with `cmdr.GetKibibytesR()`: [`kilo-bytes`](https://github.com/hedzr/cmdr-examples/tree/master/examples/kilo-bytes)





