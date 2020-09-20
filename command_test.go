// Copyright © 2020 Hedzr Yeh.

package cmdr

import "testing"

func TestCommand_GetName(t *testing.T) {
	cmd := &Command{
		BaseOpt: BaseOpt{
			Name: "a",
		},
	}
	if cmd.GetName() != "a" {
		t.Fatal("want 'a'")
	}
	if cmd.GetQuotedGroupName() != "" {
		t.Fatal("want empty group name")
	}

	cmd = &Command{
		BaseOpt: BaseOpt{
			Name:  "",
			Short: "a",
		},
	}
	a := cmd.GetExpandableNames()
	if a != "a" {
		t.Fatal("want 'a'")
	}

	cmd = &Command{
		BaseOpt: BaseOpt{
			Name:  "a",
			Short: "",
		},
	}
	a = cmd.GetExpandableNames()
	if a != "a" {
		t.Fatal("want 'a'")
	}

	cmd = &Command{
		BaseOpt: BaseOpt{
			Name: "",
			Full: "full",
		},
	}
	child := &Command{
		BaseOpt: BaseOpt{
			Name:    "u",
			Short:   "v",
			Full:    "w",
			Aliases: nil,
			Group:   "",
			owner:   cmd,
		},
	}
	if child.GetParentName() != "full" {
		t.Fatal("want 'full'")
	}

	root := &RootCommand{AppName: "aa"}
	child.root = root
	child.owner = nil
	if child.GetParentName() != "aa" {
		t.Fatal("want 'aa'")
	}

	child = &Command{
		BaseOpt: BaseOpt{
			Name:            "u",
			Short:           "v",
			Full:            "w",
			Aliases:         nil,
			Group:           "",
			owner:           cmd,
			strHit:          "",
			Description:     "",
			LongDescription: "",
			Examples:        "",
			Hidden:          false,
			Deprecated:      "",
			Action:          nil,
		},
		Flags:           nil,
		SubCommands:     nil,
		PreAction:       nil,
		PostAction:      nil,
		TailPlaceHolder: "",
		root:            nil,
		allCmds:         nil,
		allFlags:        nil,
		plainCmds:       nil,
		plainShortFlags: nil,
		plainLongFlags:  nil,
		headLikeFlag:    nil,
	}

}
