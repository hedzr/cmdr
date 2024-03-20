package cli

import (
	"fmt"
)

func (c *RootCommand) SelfAssert() {
	c.Command.SelfAssert(c)
}

func (c *Command) SelfAssert(root *RootCommand) { //nolint:revive
	c.selfWalk(root, c, nil, func(cc, oo *Command, ff *Flag) {
		if ff == nil {
			if cc.root != root {
				if !(cc.root == nil && cc == root.Command) {
					panic(fmt.Sprintf("unexpected Command.root: %v | %v", cc.root, root))
				}
			}
			if cc.owner != oo {
				panic("unexpected Command.owner")
			}
			if cc.longCommands == nil || cc.shortCommands == nil {
				panic("internal command maps not ok")
			}
			if cc.longFlags == nil || cc.shortFlags == nil {
				panic("internal flag maps not ok")
			}
			return
		}
		if ff.root != root {
			panic("unexpected Flag.root")
		}
		if ff.owner != cc {
			panic("unexpected Flag.owner" + c.Name() + "," + c.GetQuotedGroupName())
		}
	})
}

func (c *Command) selfWalk(root *RootCommand, cmd, owner *Command, cb func(cc, oo *Command, ff *Flag)) { //nolint:unparam
	cb(cmd, owner, nil)

	for _, cx := range cmd.commands {
		if cx != nil {
			c.selfWalk(root, cx, cmd, cb)
		}
	}

	for _, fx := range cmd.flags {
		if fx != nil {
			cb(cmd, owner, fx)
		}
	}
}
