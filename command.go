/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"strings"
)

// PrintHelp prints help screen
func (c *Command) PrintHelp(justFlags bool) {
	printHelp(c, justFlags)
}

// PrintVersion prints versions information
func (c *Command) PrintVersion() {
	uniqueWorker.showVersion()
}

// GetRoot returns the `RootCommand`
func (c *Command) GetRoot() *RootCommand {
	return c.root
}

// GetOwner returns the parent command object
func (c *Command) GetOwner() *Command {
	return c.owner
}

// IsRoot returns true if this command is a RootCommand
func (c *Command) IsRoot() bool {
	return c == &c.root.Command
}

// GetHitStr returns the matched command string
func (c *Command) GetHitStr() string {
	return c.strHit
}

// // HasParent detects whether owner is available or not
// func (c *BaseOpt) HasParent() bool {
// 	return c.owner != nil
// }

// GetName returns the name of a `Command`.
func (c *Command) GetName() string {
	if len(c.Full) > 0 {
		return c.Full
	}
	if len(c.Short) > 0 {
		return c.Short
	}
	return c.Name
}

// GetQuotedGroupName returns the group name quoted string.
func (c *Command) GetQuotedGroupName() string {
	if len(strings.TrimSpace(c.Group)) == 0 {
		return ""
	}
	i := strings.Index(c.Group, ".")
	if i >= 0 {
		return fmt.Sprintf("[%v]", c.Group[i+1:])
	}
	return fmt.Sprintf("[%v]", c.Group)
}

// GetExpandableNamesArray returns the names array of command, includes short name and long name.
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

// GetExpandableNames returns the names comma splitted string.
func (c *Command) GetExpandableNames() string {
	a := c.GetExpandableNamesArray()
	if len(a) == 1 {
		return a[0]
	} else if len(a) > 1 {
		return fmt.Sprintf("{%v}", strings.Join(a, ","))
	}
	return c.Name
}

// GetParentName returns the owner command name
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

// GetSubCommandNamesBy returns the joint string of subcommands
func (c *Command) GetSubCommandNamesBy(delimChar string) string {
	var a []string
	for _, sc := range c.SubCommands {
		if !sc.Hidden {
			a = append(a, sc.GetTitleNamesBy(delimChar))
		}
	}
	return strings.Join(a, delimChar)
}
