/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"strings"
)

func (c *Command) PrintHelp(justFlags bool) {
	printHelp(c, justFlags)
}

func (c *Command) PrintVersion() {
	showVersion()
}

func (c *Command) GetRoot() *RootCommand {
	return c.root
}

func (c *Command) HasParent() bool {
	return c.owner != nil
}

func (c *Command) GetName() string {
	if len(c.Full) > 0 {
		return c.Full
	}
	if len(c.Short) > 0 {
		return c.Short
	}
	return c.Name
}

func (c *Command) GetQuotedGroupName() string {
	if len(c.Group) == 0 {
		return ""
	}
	i := strings.Index(c.Group, ".")
	if i >= 0 {
		return fmt.Sprintf("[%v]", c.Group[i+1:])
	}
	return fmt.Sprintf("[%v]", c.Group)
}

func (c *Command) GetExpandableNamesArray() []string {
	var a []string
	if len(c.Full) > 0 {
		a = append(a, c.Full)
	}
	if len(c.Short) > 0 {
		a = append(a, c.Short)
	}
	return a
}
func (c *Command) GetExpandableNames() string {
	a := c.GetExpandableNamesArray()
	if len(a) == 1 {
		return a[0]
	} else if len(a) > 1 {
		return fmt.Sprintf("{%v}", strings.Join(a, ","))
	}
	return c.Name
}

func (c *Command) GetParentName() string {
	if c.owner != nil {
		if len(c.owner.Full) > 0 {
			return c.owner.Full
		}
		if len(c.owner.Short) > 0 {
			return c.owner.Short
		}
		if len(c.owner.Name) > 0 {
			return c.owner.Name
		}
	}
	return c.GetRoot().AppName
}

func (c *Command) GetSubCommandNamesBy(delimChar string) string {
	var a []string
	for _, sc := range c.SubCommands {
		a = append(a, sc.GetTitleNamesBy(delimChar))
	}
	return strings.Join(a, delimChar)
}
