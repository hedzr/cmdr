package cli

import (
	"fmt"
)

func (c *RootCommand) SelfAssert() {
	if x, ok := c.Cmd.(*CmdS); ok {
		x.SelfAssert(c)
	}
}

func (c *CmdS) SelfAssert(root *RootCommand) { //nolint:revive
	c.selfWalk(root, c, nil, func(cc, oo Cmd, ff *Flag) {
		if ff == nil {
			if cc.Root() != root {
				if !(cc.Root() == nil && cc == root.Cmd) {
					panic(fmt.Sprintf("unexpected CmdS.root: %v | root: %v | cc: %v", cc.Root(), root, cc))
				}
			}
			if cc.OwnerCmd() != oo {
				if !(cc.OwnerIsNil() && (cc.Root() == cc || oo == nil)) {
					panic(fmt.Sprintf("unexpected CmdS.owner: %v | oo: %v | cc: %v", cc.OwnerCmd(), oo, cc))
				}
			}
			if cx, ok := cc.(*CmdS); ok {
				if cx.longCommands == nil || cx.shortCommands == nil {
					panic("internal command maps (longCommands or shortCommands) not ok")
				}
				if cx.longFlags == nil || cx.shortFlags == nil {
					panic("internal flag maps (longFlags or shortFlags) not ok")
				}
			}
			return
		}
		if ff.root != root {
			panic(fmt.Sprintf("unexpected Flag.root: ff = %v", ff))
		}
		if ff.owner != cc {
			panic("unexpected Flag.owner" + c.Name() + "," + c.GetQuotedGroupName())
		}
	})
}

func (c *CmdS) selfWalk(root *RootCommand, cmd, owner Cmd, cb func(cc, oo Cmd, ff *Flag)) { //nolint:unparam
	cb(cmd, owner, nil)

	for _, cx := range cmd.SubCommands() {
		if cx != nil {
			c.selfWalk(root, cx, cmd, cb)
		}
	}

	for _, fx := range cmd.Flags() {
		if fx != nil {
			cb(cmd, owner, fx)
		}
	}
}
