/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"sync"

	"github.com/hedzr/logex"

	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log"
)

const (
	appNameDefault = "cmdr" //nolint:deadcode,unused,varcheck

	// UnsortedGroup for commands and flags
	UnsortedGroup = "zzzg.unsorted"
	// AddonsGroup for commands and flags
	AddonsGroup = "zzzh.Addons"
	// ExtGroup for commands and flags
	ExtGroup = "zzzi.Extensions"
	// AliasesGroup for commands and flags
	AliasesGroup = "zzzj.Aliases"
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
		// Name is reserved for internal purpose.
		Name string `yaml:"name,omitempty" json:"name,omitempty"`
		// Short 'rune' string. short option/command name.
		// single char. example for flag: "a" -> "-a"
		Short string `yaml:"short-name,omitempty" json:"short-name,omitempty"`
		// Full is a full/long form flag/command title name.
		// word string. example for flag: "addr" -> "--addr"
		Full string `yaml:"title,omitempty" json:"title,omitempty"`
		// Aliases are the more synonyms
		Aliases []string `yaml:"aliases,flow,omitempty" json:"aliases,omitempty"`
		// Group specify a group name,
		// A special prefix could sort it, has a form like `[0-9a-zA-Z]+\.`.
		// The prefix will be removed from help screen.
		//
		// Some examples are:
		//    "A001.Host Params"
		//    "A002.User Params"
		//
		// If ToggleGroup specified, Group field can be omitted because we will copy
		// from there.
		Group string `yaml:"group,omitempty" json:"group,omitempty"`

		owner *Command
		// strHit keeps the matched title string from user input in command line
		strHit string

		Description     string `yaml:"desc,omitempty" json:"desc,omitempty"`
		LongDescription string `yaml:"long-desc,omitempty" json:"long-desc,omitempty"`
		Examples        string `yaml:"examples,omitempty" json:"examples,omitempty"`
		Hidden          bool   `yaml:"hidden,omitempty" json:"hidden,omitempty"`
		VendorHidden    bool   `yaml:"vendor-hidden,omitempty" json:"vendor-hidden,omitempty"`

		// Deprecated is a version string just like '0.5.9' or 'v0.5.9', that
		// means this command/flag was/will be deprecated since `v0.5.9`.
		Deprecated string `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`

		// Action is callback for the last recognized command/sub-command.
		// return: ErrShouldBeStopException will break the following flow and exit right now
		// cmd 是 flag 被识别时已经得到的子命令
		Action Handler `yaml:"-" json:"-"`

		onMatched Handler //nolint:structcheck //todo is onMatch used?
	}

	// Handler handles the event on a subcommand matched
	Handler func(cmd *Command, args []string) (err error)
	// Invoker is a Handler but without error returns
	Invoker func(cmd *Command, args []string)

	// Command holds the structure of commands and sub-commands
	Command struct {
		BaseOpt `yaml:",inline"`

		Flags []*Flag `yaml:"flags,omitempty" json:"flags,omitempty"`

		SubCommands []*Command `yaml:"subcmds,omitempty" json:"subcmds,omitempty"`

		// return: ErrShouldBeStopException will break the following flow and exit right now
		PreAction Handler `yaml:"-" json:"-"`
		// PostAction will be run after Action() invoked.
		PostAction Invoker `yaml:"-" json:"-"`
		// TailPlaceHolder contains the description text of positional
		// arguments of a command.
		// It will be shown at tail of command usages line. Suppose there's
		// a TailPlaceHolder  with"<host-fqdn> <ipv4/6>", they will be
		// painted in help screen just like:
		//
		//     austr dns add <host-fqdn> <ipv4/6> [Options] [Parent/Global Options]
		//
		// Deprecated since v1.10.36+, use TailPlaceHolders is recommended
		TailPlaceHolder string `yaml:"tail-placeholder,omitempty" json:"tail-placeholder,omitempty"`
		// TailPlaceHolders gives two places to place the placeholders.
		// It looks like the following form:
		//
		//     austr dns add <placeholder1st> [Options] [Parent/Global Options] <placeholders more...>
		//
		// As shown, you may specify at:
		//
		// - before '[Options] [Parent/Global Options]'
		// - after '[Options] [Parent/Global Options]'
		//
		// In TailPlaceHolders slice, [0] is `placeholder1st``, and others
		// are `placeholders more``.
		//
		// Others:
		//   TailArgsText string [no plan]
		//   TailArgsDesc string [no plan]
		TailPlaceHolders []string ``

		root            *RootCommand
		allCmds         map[string]map[string]*Command // key1: Commnad.Group, key2: Command.Full
		allFlags        map[string]map[string]*Flag    // key1: Command.Flags[#].Group, key2: Command.Flags[#].Full
		plainCmds       map[string]*Command
		plainShortFlags map[string]*Flag
		plainLongFlags  map[string]*Flag
		headLikeFlag    *Flag

		presetCmdLines []string
		// Invoke is a space-separated string which takes Command (name) and extra
		// remain args to be invoked.
		// It invokes a command from the command tree in this app.
		// Invoke field is available for
		Invoke string `yaml:"invoke,omitempty" json:"invoke,omitempty"`
		// InvokeProc is just for cmdr aliases commands
		// invoke the external commands (via: executable)
		InvokeProc string `yaml:"invoke-proc,omitempty" json:"invoke-proc,omitempty"`
		// InvokeShell is just for cmdr aliases commands
		// invoke the external commands (via: shell)
		InvokeShell string `yaml:"invoke-sh,omitempty" json:"invoke-sh,omitempty"`
		// Shell is just for cmdr aliases commands
		// invoke a command under this shell, or /usr/bin/env bash|zsh|...
		Shell string `yaml:"shell,omitempty" json:"shell,omitempty"`

		// times how many times this flag was triggered.
		// To access it with `Command.GetTriggeredTimes()`, or `cmdr.GetHitCountByDottedPath()`.
		times int
	}

	// RootCommand holds some application information
	RootCommand struct {
		Command `yaml:",inline"`

		AppName    string `yaml:"appname,omitempty" json:"appname,omitempty"`
		Version    string `yaml:"version,omitempty" json:"version,omitempty"`
		VersionInt uint32 `yaml:"version-int,omitempty" json:"version-int,omitempty"`

		Copyright string `yaml:"copyright,omitempty" json:"copyright,omitempty"`
		Author    string `yaml:"author,omitempty" json:"author,omitempty"`
		Header    string `yaml:"header,omitempty" json:"header,omitempty"` // using `Header` for header and ignore built with `Copyright` and `Author`, and no usage lines too.

		// RunAsSubCommand give a subcommand and will be invoked
		// while app enter without any subcommands.
		//
		// For example, RunAsSubCommand is set to "build", and
		// entering "app" will equal to entering "app build ...".
		//
		// NOTE that when it's valid, RootCommand.Command.Action handler will be ignored.
		RunAsSubCommand string `yaml:"run-as,omitempty" json:"run-as,omitempty"`

		// PreActions lists all global pre-actions. They will be launched before
		// any command hit would has been invoking.
		PreActions  []Handler `yaml:"-" json:"-"`
		PostActions []Invoker `yaml:"-" json:"-"`

		ow   *bufio.Writer
		oerr *bufio.Writer
	}

	// Flag means a flag, a option, or a opt.
	Flag struct {
		BaseOpt `yaml:",inline"`

		// ToggleGroup for Toggle Group
		ToggleGroup string `yaml:"toggle-group,omitempty" json:"toggle-group,omitempty"`
		// DefaultValuePlaceholder for flag
		DefaultValuePlaceholder string `yaml:"default-placeholder,omitempty" json:"default-placeholder,omitempty"`
		// DefaultValue default value for flag
		DefaultValue interface{} `yaml:"default,omitempty" json:"default,omitempty"`
		// DefaultValueType is a string to indicate the data-type of DefaultValue.
		// It's only available in loading flag definition from a config-file (yaml/json/...).
		// Never used in writing your codes manually.
		DefaultValueType string `yaml:"type,omitempty" json:"type,omitempty"`
		// ValidArgs for enum flag
		ValidArgs []string `yaml:"valid-args,omitempty" json:"valid-args,omitempty"`
		// Required to-do
		Required bool `yaml:"required,omitempty" json:"required,omitempty"`

		// ExternalTool to get the value text by invoking external tool.
		// It's an environment variable name, such as: "EDITOR" (or cmdr.ExternalToolEditor)
		ExternalTool string `yaml:"external-tool,omitempty" json:"external-tool,omitempty"`

		// EnvVars give a list to bind to environment variables manually
		// it'll take effects since v1.6.9
		EnvVars []string `yaml:"envvars,flow,omitempty" json:"envvars,omitempty"`

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
		HeadLike bool `yaml:"head-like,omitempty" json:"head-like,omitempty"`

		// Min minimal value of a range.
		Min int64 `yaml:"min,omitempty" json:"min,omitempty"`
		// Max maximal value of a range.
		Max int64 `yaml:"max,omitempty" json:"max,omitempty"`

		onSet func(keyPath string, value interface{})

		// times how many times this flag was triggered.
		// To access it with `Flag.GetTriggeredTimes()`, `cmdr.GetFlagHitCount()`,
		// `cmdr.GetFlagHitCountRecursively()` or `cmdr.GetHitCountByDottedPath()`.
		times int

		// actionStr: for zsh completion, see action of an optspec in _argument
		actionStr string
		// mutualExclusives is used for zsh completion.
		//
		// For the ToggleGroup group, mutualExclusives is implicit.
		mutualExclusives []string
		// prerequisites flags for this one.
		//
		// In zsh completion, any of prerequisites flags must be present
		// so that user can complete this one.
		//
		// The prerequisites were not present and cmdr would report error
		// and stop parsing flow.
		prerequisites []string
		// justOnce is used for zsh completion.
		justOnce bool
		// circuitBreak is used for zsh completion.
		//
		// A flag can break cmdr parsing flow with return
		// ErrShouldBeStopException in its Action handler.
		// But you' better told zsh system with set circuitBreak
		// to true. At this case, cmdr will generate a suitable
		// completion script.
		circuitBreak bool
		// dblTildeOnly can be used for zsh completion.
		//
		// A DblTildeOnly Flag accepts '~~opt' only, so '--opt' is
		// invalid form and couldn't be used for other Flag
		// anymore.
		dblTildeOnly bool // such as '~~tree'

		// PostAction treat this flag as a command!
		// PostAction Handler

		// by default, a flag is always `optional`.
	}

	// Options is a holder of all options
	Options struct {
		entries   map[string]interface{}
		hierarchy map[string]interface{}
		rw        *sync.RWMutex

		usedConfigFile            string
		usedConfigSubDir          string
		usedAlterConfigFile       string
		usedSecondaryConfigFile   string
		usedSecondaryConfigSubDir string
		configFiles               []string
		filesWatching             []string
		batchMerging              bool
		appendMode                bool // true: append mode, false: replaceMode

		onConfigReloadedFunctions map[ConfigReloaded]bool
		rwlCfgReload              *sync.RWMutex
		rwCB                      sync.RWMutex
		onMergingSet              OnOptionSetCB
		onSet                     OnOptionSetCB
	}

	// OptOne struct {
	// 	Children map[string]*OptOne `yaml:"c,omitempty"`
	// 	Value    interface{}        `yaml:"v,omitempty"`
	// }

	// ConfigReloaded for config reloaded
	ConfigReloaded interface {
		OnConfigReloaded()
	}

	// OnOptionSetCB is a callback function while an option is being set (or merged)
	OnOptionSetCB func(keyPath string, value, oldVal interface{})
	// OnSwitchCharHitCB is a callback function ...
	OnSwitchCharHitCB func(parsed *Command, switchChar string, args []string) (err error)
	// OnPassThruCharHitCB is a callback function ...
	OnPassThruCharHitCB func(parsed *Command, switchChar string, args []string) (err error)

	// HookFunc the hook function prototype for SetBeforeXrefBuilding and SetAfterXrefBuilt
	HookFunc func(root *RootCommand, args []string)

	// HookOptsFunc the hook function prototype
	HookOptsFunc func(root *RootCommand, opts *Options)

	// HookHelpScreenFunc the hook function prototype
	HookHelpScreenFunc func(w *ExecWorker, p Painter, cmd *Command, justFlags bool)
)

var (
	//
	// doNotLoadingConfigFiles = false

	// // rootCommand the root of all commands
	// rootCommand *RootCommand
	// // rootOptions *Opt
	// rxxtOptions = newOptions()

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

	// CurrentHiddenColor the print color for left part of a hidden opt
	CurrentHiddenColor = FgDarkGray

	// CurrentDeprecatedColor the print color for deprecated opt line
	CurrentDeprecatedColor = FgDarkGray

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

	defaultStringMetric = tool.JaroWinklerDistance(tool.JWWithThreshold(similarThreshold))
)

const similarThreshold = 0.6666666666666666

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

// GetTraceMode returns the flag value of `--trace`/`-tr`
//
// NOTE
//     log.GetTraceMode()/SetTraceMode() have higher universality
//
// the flag value of `--trace` or `-tr` is always stored
// in cmdr Option Store, so you can retrieved it by
// GetBoolR("trace") and set it by Set("trace", true).
// You could also set it with SetTraceMode(b bool).
//
// The `--trace` is not enabled in default, so you have to
// add it manually:
//
//     import "github.com/hedzr/cmdr-addons/pkg/plugins/trace"
//     cmdr.Exec(buildRootCmd(),
//         trace.WithTraceEnable(true),
//     )
func GetTraceMode() bool {
	return GetBoolR("trace") || log.GetTraceMode()
}

// SetTraceMode setup the tracing mode status in Option Store
func SetTraceMode(b bool) {
	Set("trace", b)
	logex.SetTraceMode(b)
}

// GetDebugMode returns the flag value of `--debug`/`-D`
//
// NOTE
//     log.GetDebugMode()/SetDebugMode() have higher universality
//
// the flag value of `--debug` or `-D` is always stored
// in cmdr Option Store, so you can retrieved it by
// GetBoolR("debug") and set it by Set("debug", true).
// You could also set it with SetDebugMode(b bool).
func GetDebugMode() bool {
	return GetBoolR("debug") || log.GetDebugMode()
}

// SetDebugMode setup the debug mode status in Option Store
func SetDebugMode(b bool) {
	Set("debug", b)
	logex.SetDebugMode(b)
}

// NewLoggerConfig returns a default LoggerConfig
func NewLoggerConfig() *log.LoggerConfig {
	lc := NewLoggerConfigWith(false, "sugar", "error")
	return lc
}

// NewLoggerConfigWith returns a default LoggerConfig
func NewLoggerConfigWith(enabled bool, backend, level string) *log.LoggerConfig {
	log.SetTraceMode(GetTraceMode())
	log.SetDebugMode(GetDebugMode())
	lc := log.NewLoggerConfigWith(enabled, backend, level)
	_ = GetSectionFrom("logger", &lc)
	return lc
}

// GetVerboseMode returns the flag value of `--verbose`/`-v`
func GetVerboseMode() bool {
	return GetBoolR("verbose")
}

// GetQuietMode returns the flag value of `--quiet`/`-q`
func GetQuietMode() bool {
	return GetBoolR("quiet")
}

// GetNoColorMode return the flag value of `--no-color`/`-nc`
func GetNoColorMode() bool {
	return GetBoolR("no-color")
}

// GetVerboseModeHitCount returns how many times `--verbose`/`-v` specified
func GetVerboseModeHitCount() int { return GetFlagHitCount("verbose") }

// GetQuietModeHitCount returns how many times `--quiet`/`-q` specified
func GetQuietModeHitCount() int { return GetFlagHitCount("quiet") }

// GetNoColorModeHitCount returns how many times `--no-color`/`-nc` specified
func GetNoColorModeHitCount() int { return GetFlagHitCount("no-color") }

// GetDebugModeHitCount returns how many times `--debug`/`-D` specified
func GetDebugModeHitCount() int { return GetFlagHitCount("debug") }

// GetTraceModeHitCount returns how many times `--trace`/`-tr` specified
func GetTraceModeHitCount() int { return GetFlagHitCount("trace") }

// GetFlagHitCount return how manu times a top-level Flag was specified from command-line.
func GetFlagHitCount(longName string) int {
	w := internalGetWorker()
	return getFlagHitCount(&w.rootCommand.Command, longName)
}

func getFlagHitCount(c *Command, longName string) int {
	if f := c.root.FindFlag(longName); f != nil {
		return f.times
	}
	return 0
}

// GetFlagHitCountRecursively return how manu times a Flag was specified from command-line.
// longName will be search recursively.
func GetFlagHitCountRecursively(longName string) int {
	w := internalGetWorker()
	return getFlagHitCountRecursively(&w.rootCommand.Command, longName)
}

func getFlagHitCountRecursively(c *Command, longName string) int {
	if f := c.root.FindFlagRecursive(longName); f != nil {
		return f.times
	}
	return 0
}

// GetHitCountByDottedPath return how manu times a Flag or a Command was specified from command-line.
func GetHitCountByDottedPath(dottedPath string) int {
	c, f := dottedPathToCommandOrFlag(dottedPath, nil)
	if c != nil {
		return c.times
	}
	if f != nil {
		return f.times
	}
	return 0
}

// func init() {
// 	// onConfigReloadedFunctions = make(map[ConfigReloaded]bool)
// 	// SetCurrentHelpPainter(new(helpPainter))
// }
