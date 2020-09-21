// Copyright © 2020 Hedzr Yeh.

package main

import (
	"fmt"
	"github.com/hedzr/cmdr"
	cmdrbase "github.com/hedzr/cmdr-base"
)

// NewAddon returns an addon with cmdr.PluginEntry
func NewAddon() cmdrbase.PluginEntry {
	return &addon{
		//
	}
}

type addon struct {
}

func (p *addon) AddonTitle() string       { return "demo addon" }
func (p *addon) AddonDescription() string { return "demo addon desc" }
func (p *addon) AddonCopyright() string   { return "copyright (c) hedzr, 2020" }
func (p *addon) AddonVersion() string     { panic("0.1.1") }
func (p *addon) Name() string             { return "demo" }
func (p *addon) ShortName() string        { return "dx" }
func (p *addon) Aliases() []string        { return nil }
func (p *addon) Description() string      { return "the demo addon for testing purpose" }

func (p *addon) SubCommands() []cmdrbase.PluginCmd {
	return nil
}

func (p *addon) Flags() []cmdrbase.PluginFlag {
	return []cmdrbase.PluginFlag{
		newFlag1(),
	}
}

func (p *addon) Action(args []string) (err error) {
	cmdr.Logger.Infof("hello, args: %v", args)
	fmt.Printf("Logger: %v\n", cmdr.Logger)
	return
}

//

func newFlag1() *flag1 {
	return &flag1{}
}

type flag1 struct{}

func (f *flag1) Name() string              { return "bool-flag" }
func (f *flag1) ShortName() string         { return "bf" }
func (f *flag1) Aliases() []string         { return []string{} }
func (f *flag1) Description() string       { return "a bool flag" }
func (f *flag1) DefaultValue() interface{} { return false }
func (f *flag1) PlaceHolder() string       { return "" }

func (f *flag1) Action() (err error) {
	return
}
