# cmdr

getopt/getopt_long like command-line UI golang library.

A getopt-like parser of command-line options, compatible with the [getopt_long](http://www.gnu.org/s/libc/manual/html_node/Argument-Syntax.html#Argument-Syntax) syntax, which is an extension of the syntax recommended by POSIX.

`cmdr` is a UNIX command-line UI library written by golang.


## Features

- Unix [*getopt*(3)](http://man7.org/linux/man-pages/man3/getopt.3.html) representation but without its programmatic interface.
- Automatic help screen generation
- Support for unlimited multiple sub-commands.
- Support for command short and long name, and aliases names.
- Support for both short and long options (`-o` and `—opt`). Support for multiple aliases
- Automatically allows both `-I file` and `-Ifile`, and `-I=files` formats.
- Support for `-D+`, `-D-` to enable/disable a bool flag.
- Support for circuit-break by `--`.
- Support for options being specified multiple times, with different values
- Support for optional arguments.
- Groupable commands and options/flags.
- Sortable commands and options/flags. Or sorted by alphabetic order.
- Bash and Zsh (not yet) completion.
- Predefined yaml config file locations:
  - `/etc/<appname>/<appname>.yml` and `conf.d` sub-directory.
  - `$HOME/<appname>/<appname>,yml` and `conf.d` sub-directory.
  - Watch `conf.d` directory:
    - `AddOnConfigLoadedListener(c)`
    - `RemoveOnConfigLoadedListener(c)`
    - `SetOnConfigLoadedListener(c, enabled)`
  - As a feature, do NOT watch the changes on `<appname>.yml`.
- Overrides by environment variables.
- `cmdr.GetBool(key)`, `cmdr.GetInt(key)`, `cmdr.GetString(key)`, `cmdr.GetStringSlice(key)` for Option value extraction.





## LICENSE

MIT.





