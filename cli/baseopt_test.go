package cli

import (
	"testing"

	"github.com/hedzr/evendeep"
)

func TestBaseOpt_GetDottedPath(t *testing.T) {
	root := rootCmdForTesting()
	for i, tc := range []struct {
		dp string
	}{
		// {"microservices.tags"},
		{""},
		{"server"},
		{"server.head"},
		{"server.tail"},
		{"server.enum"},
		{"server.retry"},
		{"server.start"},
		{"server.start.foreground"},
		{"server.stop"},
		{"server.restart"},
		{"kvstore"},
		{"kvstore.backup"},
		{"kvstore.backup.output"},
		{"microservices.tags"},
		{"microservices.tags.list"},
		{"microservices.tags.add.list"},
		{"microservices.tags.modify"},
		{"microservices.tags.modify.add"},
		{"microservices.tags.modify.rm"},
		{"microservices.tags.modify.ued"},
		{"microservices.tags.toggle.address"},
	} {
		c, f := root.DottedPathToCommandOrFlag(tc.dp)
		if f != nil {
			if s := f.GetDottedPath(); s != tc.dp {
				t.Fatalf("%5d. GetDottedPath expected %q, got %q", i, tc.dp, s)
			}
		} else if c != nil {
			if s := c.GetDottedPath(); s != tc.dp {
				t.Fatalf("%5d. GetDottedPath expected %q, got %q", i, tc.dp, s)
			}
		} else {
			t.Fatalf("%5d. expected %q, but cmd/flg not found", i, tc.dp)
		}
	}
}

func TestClone(t *testing.T) {
	root := &CmdS{
		BaseOpt:               BaseOpt{},
		tailPlaceHolders:      nil,
		commands:              nil,
		flags:                 nil,
		preActions:            nil,
		onInvoke:              nil,
		postActions:           nil,
		onMatched:             nil,
		onEvalSubcommands:     nil,
		onEvalSubcommandsOnce: nil,
		onEvalFlags:           nil,
		onEvalFlagsOnce:       nil,
		redirectTo:            "",
		presetCmdLines:        nil,
		invokeProc:            "",
		invokeShell:           "",
		shell:                 "",
		longCommands:          nil,
		shortCommands:         nil,
		longFlags:             nil,
		shortFlags:            nil,
		allCommands:           nil,
		allFlags:              nil,
		toggles:               nil,
		headLikeFlag:          nil,
	}
	rc := &RootCommand{
		AppName:     "",
		Version:     "",
		Copyright:   "",
		Author:      "",
		HeaderLine:  "",
		FooterLine:  "",
		Cmd:         root,
		preActions:  nil,
		postActions: nil,
		linked:      0,
		app:         nil,
	}
	ff := &Flag{
		BaseOpt: BaseOpt{
			owner:        root,
			root:         rc,
			name:         "",
			Long:         "warning",
			Short:        "",
			Aliases:      nil,
			description:  "",
			longDesc:     "",
			examples:     "",
			group:        "",
			extraShorts:  nil,
			deprecated:   "",
			hidden:       false,
			vendorHidden: false,
			hitTitle:     "",
			hitTimes:     0,
		},
		toggleGroup:    "toggleGroup",
		placeHolder:    "placeHolder",
		defaultValue:   []int{2, 3, 4, 5},
		envVars:        []string{"VAR1", "VAR2"},
		externalEditor: "EDITOR",
		validArgs:      []string{"apple", "banana"},
		min:            -1,
		max:            99,
		headLike:       true,
		requited:       true,
		onParseValue: func(f *Flag, position int, hitCaption string, hitValue string, moreArgs []string) (newVal any, remainPartInHitValue string, err error) {
			return
		},
		onMatched:        nil,
		onChanging:       nil,
		onChanged:        nil,
		onSet:            nil,
		actionStr:        "ACTION_STR",
		mutualExclusives: nil,
		prerequisites:    nil,
		justOnce:         true,
		circuitBreak:     true,
		dblTildeOnly:     true,
		negatable:        true,
	}

	f := evendeep.MakeClone(ff)
	if f == nil {
		t.Fatalf("MakeClone returned nil")
	}
	if f1, ok := f.(*Flag); ok && f1.Long == "warning" {
		t.Log("OK")
		return
	} else if f1, ok := f.(Flag); ok && f1.Long == "warning" {
		t.Log("OK")
		return
	}

	t.Fatalf("f.Long!='warning', got %v", f)
}
