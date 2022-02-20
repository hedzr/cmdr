// Copyright © 2020 Hedzr Yeh.

package cmdr_test

import (
	"github.com/hedzr/cmdr"
	"gopkg.in/hedzr/errors.v3"
	"testing"
	"time"
)

func TestCommand_EqualTo(t *testing.T) {
	var c1 *cmdr.Command
	cmd := &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Name: "a",
		},
	}

	c1.EqualTo(cmd)
	cmd.EqualTo(c1)
}

func TestCommand_GetName(t *testing.T) {
	cmd := &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Name: "a",
		},
	}
	if cmd.GetName() != "a" {
		t.Fatal("want 'a'")
	}
	if cmd.GetQuotedGroupName() != "" {
		t.Fatal("want empty group name")
	}

	cmd = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Name:  "",
			Short: "a",
		},
	}
	a := cmd.GetExpandableNames()
	if a != "a" {
		t.Fatal("want 'a'")
	}

	cmd = &cmdr.Command{
		BaseOpt: cmdr.BaseOpt{
			Name:  "a",
			Short: "",
		},
	}
	a = cmd.GetExpandableNames()
	if a != "a" {
		t.Fatal("want 'a'")
	}

}

func TestRootCmdOpt_RunAsSubCommand(t *testing.T) {
	testFramework(t, func() *cmdr.RootCommand {
		x := rootCmdForTesting()
		x.RunAsSubCommand = "microservices.tags"
		return x
	}, testCases{

		"consul-tags -p8500 --prefix=1 --prefix2 -ui123 --uint 2 -dur3h -flt 9.8 --uint 139 --prefix 3": func(t *testing.T, c *cmdr.Command, e error) (err error) {
			// all ok,
			//err = cmdr.InvokeCommand("microservices.tags")

			if cmdr.GetInt("app.ms.tags.port") != 8500 || cmdr.GetString("app.ms.tags.prefix") != "3" ||
				cmdr.GetUint("app.ms.tags.uint") != uint(139) || cmdr.GetFloat32("app.ms.tags.float") != 9.8 ||
				cmdr.GetDuration("app.ms.tags.duration") != 3*time.Hour ||
				cmdr.GetBool("debug") || cmdr.GetVerboseMode() {
				return errors.New("something wrong 3. |%v|%v|%v|%v|%v|%v",
					cmdr.GetInt("app.ms.tags.port"), cmdr.GetString("app.ms.tags.prefix"),
					cmdr.GetUint("app.ms.tags.uint"), cmdr.GetFloat32("app.ms.tags.float"),
					cmdr.GetDuration("app.ms.tags.duration"),
					cmdr.GetBool("debug"), cmdr.GetVerboseMode())
			}

			return
		},
	},
		cmdr.WithInternalDefaultAction(true),
	)
}
