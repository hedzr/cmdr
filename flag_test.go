package cmdr

import (
	"testing"
)

func TestFlag_EqualTo(t *testing.T) {
	var c1 *Flag
	cmd := &Flag{
		BaseOpt: BaseOpt{
			Name:  "a",
			Short: "a",
		},
	}

	c1.EqualTo(cmd)

	cmd.EqualTo(c1)
	cmd.GetDescZsh()
	cmd.GetTitleZshFlagName()
}

func TestFlag_Delete(t *testing.T) {
	var c1 *Flag
	cmd := &Command{
		BaseOpt: BaseOpt{
			Name:  "a",
			Short: "a",
		},
	}
	c1 = &Flag{BaseOpt: BaseOpt{Full: "b", owner: cmd}}
	cmd.Flags = uniAddFlg(cmd.Flags, c1)
	c1.GetDottedNamePath()
	c1.Delete()
}

func TestCommand_Delete(t *testing.T) {
	cmd := &Command{
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
		Flags:            nil,
		SubCommands:      nil,
		PreAction:        nil,
		PostAction:       nil,
		TailPlaceHolder:  "",
		TailPlaceHolders: []string{},
		root:             nil,
		allCmds:          nil,
		allFlags:         nil,
		plainCmds:        nil,
		plainShortFlags:  nil,
		plainLongFlags:   nil,
		headLikeFlag:     nil,
	}

	cmd.SubCommands = append(cmd.SubCommands, child)
	child.Delete()
}
