/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"errors"
	"fmt"
	"sync"
)

const (
	appNameDefault = "cmdr"

	// UnsortedGroup for commands and flags
	UnsortedGroup = "zzzz.unsorted"
	// SysMgmtGroup for commands and flags
	SysMgmtGroup = "zzz9.Misc"

	// DefaultEditor is 'vim'
	DefaultEditor = "vim"

	// ExternalToolEditor environment variable name, EDITOR is fit for most of shells.
	ExternalToolEditor = "EDITOR"

	// ExternalToolPasswordInput enables secure password input without echo.
	ExternalToolPasswordInput = "PASSWD"
)

type (
	// BaseOpt is base of `Command`, `Flag`
	BaseOpt struct {
		Name string
		// Short rune. short option/command name.
		// single char. example for flag: "a" -> "-a"
		Short string
		// Full full/long option/command name.
		// word string. example for flag: "addr" -> "--addr"
		Full string
		// Aliases are the more synonyms
		Aliases []string
		// Group group name
		Group string

		owner  *Command
		strHit string

		Description     string
		LongDescription string
		Examples        string
		Hidden          bool

		// Deprecated is a version string just like '0.5.9', that means this command/flag was/will be deprecated since `v0.5.9`.
		Deprecated string

		// Action is callback for the last recognized command/sub-command.
		// return: ErrShouldBeStopException will break the following flow and exit right now
		// cmd 是 flag 被识别时已经得到的子命令
		Action func(cmd *Command, args []string) (err error)
	}

	// Command holds the structure of commands and subcommands
	Command struct {
		BaseOpt

		Flags []*Flag

		SubCommands []*Command
		// return: ErrShouldBeStopException will break the following flow and exit right now
		PreAction func(cmd *Command, args []string) (err error)
		// PostAction will be run after Action() invoked.
		PostAction func(cmd *Command, args []string)
		// be shown at tail of command usages line. Such as for TailPlaceHolder="<host-fqdn> <ipv4/6>":
		// austr dns add <host-fqdn> <ipv4/6> [Options] [Parent/Global Options]
		TailPlaceHolder string

		root            *RootCommand
		allCmds         map[string]map[string]*Command // key1: Commnad.Group, key2: Command.Full
		allFlags        map[string]map[string]*Flag    // key1: Command.Flags[#].Group, key2: Command.Flags[#].Fullui
		plainCmds       map[string]*Command
		plainShortFlags map[string]*Flag
		plainLongFlags  map[string]*Flag
		headLikeFlag    *Flag
	}

	// RootCommand holds some application information
	RootCommand struct {
		Command

		AppName    string
		Version    string
		VersionInt uint32

		Copyright string
		Author    string
		Header    string // using `Header` for header and ignore built with `Copyright` and `Author`, and no usage lines too.

		ow   *bufio.Writer
		oerr *bufio.Writer
	}

	// Flag means a flag, a option, or a opt.
	Flag struct {
		BaseOpt

		// ToggleGroup: to-do: Toggle Group
		ToggleGroup string
		// DefaultValuePlaceholder for flag
		DefaultValuePlaceholder string
		// DefaultValue default value for flag
		DefaultValue interface{}
		// ValidArgs to-do
		ValidArgs []string
		// Required to-do
		Required bool

		// ExternalTool to get the value text by invoking external tool
		// It's an environment variable name, such as: "EDITOR" (or cmdr.ExternalToolEditor)
		ExternalTool string

		// HeadLike enables a free-hand option like `head -3`.
		//
		// When a free-hand option presents, it'll be treated as a named option with an integer value.
		//
		// For example, option/flag = `{{Full:"line",Short:"l"},HeadLike:true}`, the command line:
		// `app -3`
		// is equivalent to `app -l 3`, and so on.
		//
		// HeadLike assumed an named option with an integer value, that means, Min and Max can be applied on it too.
		// NOTE: Only one head-like option can be defined in a command/sub-command chain.
		HeadLike bool

		// Min minimal value of a range.
		Min int64
		// Max maximal value of a range.
		Max int64

		onSet func(keyPath string, value interface{})

		// times how many times this flag was triggered.
		times int

		// PostAction treat this flag as a command!
		// PostAction func(cmd *Command, args []string) (err error)

		// by default, a flag is always `optional`.
	}

	// Options is a holder of all options
	Options struct {
		entries   map[string]interface{}
		hierarchy map[string]interface{}
		rw        *sync.RWMutex

		usedConfigFile   string
		usedConfigSubDir string
		configFiles      []string

		onConfigReloadedFunctions map[ConfigReloaded]bool
		rwlCfgReload              *sync.RWMutex
	}

	// OptOne struct {
	// 	Children map[string]*OptOne `yaml:"c,omitempty"`
	// 	Value    interface{}        `yaml:"v,omitempty"`
	// }

	// ConfigReloaded for config reloaded
	ConfigReloaded interface {
		OnConfigReloaded()
	}

	// HookFunc the hook function prototype for SetBeforeXrefBuilding and SetAfterXrefBuilt
	HookFunc func(root *RootCommand, args []string)

	// HookOptsFunc the hook function prototype
	HookOptsFunc func(root *RootCommand, opts *Options)
)

var (
	//
	// doNotLoadingConfigFiles = false

	// // rootCommand the root of all commands
	// rootCommand *RootCommand
	// // rootOptions *Opt
	// rxxtOptions = NewOptions()

	// usedConfigFile
	// usedConfigFile            string
	// usedConfigSubDir          string
	// configFiles               []string
	// onConfigReloadedFunctions map[ConfigReloaded]bool
	//
	// predefinedLocations = []string{
	// 	"./ci/etc/%s/%s.yml",
	// 	"/etc/%s/%s.yml",
	// 	"/usr/local/etc/%s/%s.yml",
	// 	os.Getenv("HOME") + "/.%s/%s.yml",
	// }

	//
	// defaultStdout = bufio.NewWriterSize(os.Stdout, 16384)
	// defaultStderr = bufio.NewWriterSize(os.Stderr, 16384)

	//
	// currentHelpPainter Painter

	// CurrentDescColor the print color for description line
	CurrentDescColor = FgDarkGray
	// CurrentDefaultValueColor the print color for default value line
	CurrentDefaultValueColor = FgDarkGray
	// CurrentGroupTitleColor the print color for titles
	CurrentGroupTitleColor = DarkColor

	// globalShowVersion   func()
	// globalShowBuildInfo func()

	// beforeXrefBuilding []HookFunc
	// afterXrefBuilt     []HookFunc

	// getEditor sets callback to get editor program
	// getEditor func() (string, error)

	// ErrShouldBeStopException tips `Exec()` cancelled the following actions after `PreAction()`
	ErrShouldBeStopException = errors.New("should be stop right now")

	errWrongEnumValue = errors.New("unexpect enumerable value '%s' for option '%s', under command '%s'")

	defaultStringMetric = JaroWinklerDistance(JWWithThreshold(similarThreshold))
)

const similarThreshold = 0.6666666666666666

// ErrorForCmdr structure
type ErrorForCmdr struct {
	Inner     error
	Ignorable bool
	Msg       string
}

// NewError formats a ErrorForCmdr object
func NewError(ignorable bool, inner error, args ...interface{}) *ErrorForCmdr {
	if len(args) > 0 {
		return &ErrorForCmdr{Inner: nil, Ignorable: ignorable, Msg: fmt.Sprintf(inner.Error(), args...)}
	}
	return &ErrorForCmdr{Inner: inner, Ignorable: ignorable}
}

// NewErrorWithMsg formats a ErrorForCmdr object
func NewErrorWithMsg(msg string, inner error) *ErrorForCmdr {
	return &ErrorForCmdr{Inner: inner, Ignorable: false, Msg: msg}
}

func (s *ErrorForCmdr) Error() string {
	if s.Inner != nil {
		return fmt.Sprintf("Error: %v. Inner: %v", s.Msg, s.Inner.Error())
	}
	return s.Msg
}

// GetStrictMode enables error when opt value missed. such as:
// xxx a b --prefix''   => error: prefix opt has no value specified.
// xxx a b --prefix'/'  => ok.
//
// ENV: use `CMDR_APP_STRICT_MODE=true` to enable strict-mode.
// NOTE: `CMDR_APP_` prefix could be set by user (via: `EnvPrefix` && `RxxtPrefix`).
//
// the flag value of `--strict-mode`.
func GetStrictMode() bool {
	return GetBoolR("strict-mode")
}

// GetDebugMode returns the flag value of `--debug`/`-D`
func GetDebugMode() bool {
	return GetBoolR("debug")
}

// GetVerboseMode returns the flag value of `--verbose`/`-v`
func GetVerboseMode() bool {
	return GetBoolR("verbose")
}

// GetQuietMode returns the flag value of `--quiet`/`-q`
func GetQuietMode() bool {
	return GetBoolR("quiet")
}

// GetNoColorMode return the flag value of `--no-color`
func GetNoColorMode() bool {
	return GetBoolR("no-color")
}

// func init() {
// 	// onConfigReloadedFunctions = make(map[ConfigReloaded]bool)
// 	// SetCurrentHelpPainter(new(helpPainter))
// }
