/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"errors"
)

const (
	APP_NAME_DEFAULT = "cmdr"

	UNSORTED_GROUP = "zzzz.unsorted"
	SYSMGMT        = "zzz9.Misc"
)

type (
	getopt struct {
		Name string
	}

	BaseOpt struct {
		Name string
		// single char. example for flag: "a" -> "-a"
		// Short rune.
		Short string
		// word string. example for flag: "addr" -> "--addr"
		Full string
		// more synonyms
		Aliases []string
		// group name
		Group string

		owner *Command

		Flags []*Flag

		Description             string
		LongDescription         string
		Examples                string
		Hidden                  bool
		DefaultValuePlaceholder string

		// cmd 是 flag 被识别时已经得到的子命令
		// return: ShouldBeStopException will break the following flow and exit right now
		Action func(cmd *Command, args []string) (err error)
	}

	Command struct {
		BaseOpt
		SubCommands []*Command
		// return: ShouldBeStopException will break the following flow and exit right now
		PreAction func(cmd *Command, args []string) (err error)
		// PostAction will be run after Action() invoked.
		PostAction func(cmd *Command, args []string)
		// be shown at tail of command usages line. Such as for TailPlaceHolder="<host-fqdn> <ipv4/6>":
		// austr dns add <host-fqdn> <ipv4/6> [Options] [Parent/Global Options]
		TailPlaceHolder string

		root       *RootCommand
		allCmds    map[string]map[string]*Command // key1: Commnad.Group, key2: Command.Full
		allFlags   map[string]map[string]*Flag    // key1: Command.Flags[#].Group, key2: Command.Flags[#].Full
		plainCmds  map[string]*Command
		plainFlags map[string]*Flag
	}

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

	Flag struct {
		BaseOpt

		// default value for flag
		DefaultValue interface{}
		ValidArgs    []string
		Required     bool

		// by default, a flag is always `optional`.
	}

	Options struct {
		entries   map[string]interface{}
		hierarchy map[string]interface{}
	}

	// OptOne struct {
	// 	Children map[string]*OptOne `yaml:"c,omitempty"`
	// 	Value    interface{}        `yaml:"v,omitempty"`
	// }

	ConfigReloaded interface {
		OnConfigReloaded()
	}
)

var (
	EnableVersionCommands  bool = true
	EnableHelpCommands     bool = true
	EnableVerboseCommands  bool = true
	EnableGenerateCommands bool = true

	rootCommand *RootCommand
	// rootOptions *OptOne
	RxxtOptions *Options = NewOptions()

	// RxxtPrefix create a top-level namespace, which contains all normalized `Flag`s.
	RxxtPrefix = []string{"app",}

	EnvPrefix = []string{"CMDR",}

	usedConfigFile            string
	usedConfigSubDir          string
	configFiles               []string
	onConfigReloadedFunctions map[ConfigReloaded]bool
	//

	globalShowVersion   func()
	globalShowBuildInfo func()

	//
	ShouldBeStopException = errors.New("should be stop right now")
)
