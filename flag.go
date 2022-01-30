/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/tool"
	"strings"
)

// GetTriggeredTimes returns the matched times
func (f *Flag) GetTriggeredTimes() int {
	return f.times
}

// GetDescZsh temp
func (f *Flag) GetDescZsh() (desc string) {
	desc = f.Description
	if len(desc) == 0 {
		desc = tool.EraseAnyWSs(f.GetTitleZshFlagName())
	}
	// desc = replaceAll(desc, " ", "\\ ")
	desc = reSQ.ReplaceAllString(desc, `*$1*`)
	desc = reBQ.ReplaceAllString(desc, `**$1**`)
	desc = reSQnp.ReplaceAllString(desc, "''")
	desc = reBQnp.ReplaceAllString(desc, "\\`")
	desc = strings.ReplaceAll(desc, ":", "\\:")
	desc = strings.ReplaceAll(desc, "[", "\\[")
	desc = strings.ReplaceAll(desc, "]", "\\]")
	return
}

// GetTitleFlagNames temp
func (f *Flag) GetTitleFlagNames() string {
	return f.GetTitleFlagNamesBy(",")
}

// GetTitleZshFlagName temp
func (f *Flag) GetTitleZshFlagName() (str string) {
	if len(f.Full) > 0 {
		str += "--" + f.Full
	} else if len(f.Short) > 0 {
		str += "-" + f.Short
	}
	return
}

// GetTitleZshFlagShortName temp
func (f *Flag) GetTitleZshFlagShortName() (str string) {
	if len(f.Short) > 0 {
		str += "-" + f.Short
	} else if len(f.Full) > 0 {
		str += "--" + f.Full
	}
	return
}

// GetTitleZshNamesBy temp
func (f *Flag) GetTitleZshNamesBy(delimChar string, allowPrefix, quoted bool) (str string) {
	return f.GetTitleZshNamesExtBy(delimChar, allowPrefix, quoted, true, true)
}

// GetTitleZshNamesExtBy temp
func (f *Flag) GetTitleZshNamesExtBy(delimChar string, allowPrefix, quoted, shortTitleOnly, longTitleOnly bool) (str string) {
	// quote := false
	prefix, suffix := "", ""
	if _, ok := f.DefaultValue.(bool); !ok {
		suffix = "="
		//} else if _, ok := s.DefaultValue.(bool); ok {
		//	suffix = "-"
	}
	if allowPrefix && !f.justOnce {
		quoted, prefix = true, "*"
	}
	if !longTitleOnly && len(f.Short) > 0 {
		if quoted {
			str += "'" + prefix + "-" + f.Short + suffix + "'"
		} else {
			str += prefix + "-" + f.Short + suffix
		}
		if shortTitleOnly {
			return
		}
	}
	if len(f.Full) > 0 {
		if str != "" {
			str += delimChar
		}
		if quoted {
			str += "'" + prefix + "--" + f.Full + suffix + "'"
		} else {
			str += prefix + "--" + f.Full + suffix
		}
	}
	return
}

// GetTitleZshFlagNamesArray temp
func (f *Flag) GetTitleZshFlagNamesArray() (ary []string) {
	if len(f.Short) == 1 || len(f.Short) == 2 {
		if len(f.DefaultValuePlaceholder) > 0 {
			ary = append(ary, "-"+f.Short+"=") // +s.DefaultValuePlaceholder)
		} else {
			ary = append(ary, "-"+f.Short)
		}
	}
	if len(f.Full) > 0 {
		if len(f.DefaultValuePlaceholder) > 0 {
			ary = append(ary, "--"+f.Full+"=") // +s.DefaultValuePlaceholder)
		} else {
			ary = append(ary, "--"+f.Full)
		}
	}
	return
}

// GetTitleFlagNamesBy temp
func (f *Flag) GetTitleFlagNamesBy(delimChar string) string {
	return f.GetTitleFlagNamesByMax(delimChar, len(f.Short))
}

// GetTitleFlagNamesByMax temp
func (f *Flag) GetTitleFlagNamesByMax(delimChar string, maxShort int) string {
	var sb strings.Builder

	if len(f.Short) == 0 {
		// if no flag.Short,
		sb.WriteString(strings.Repeat(" ", maxShort))
	} else {
		sb.WriteRune('-')
		sb.WriteString(f.Short)
		sb.WriteString(delimChar)
		if len(f.Short) < maxShort {
			sb.WriteString(strings.Repeat(" ", maxShort-len(f.Short)))
		}
	}

	if len(f.Short) == 0 {
		sb.WriteRune(' ')
		sb.WriteRune(' ')
	}
	sb.WriteRune(' ')
	sb.WriteString("--")
	sb.WriteString(f.Full)
	if len(f.DefaultValuePlaceholder) > 0 {
		// str += fmt.Sprintf("=\x1b[2m\x1b[%dm%s\x1b[0m", DarkColor, s.DefaultValuePlaceholder)
		sb.WriteString(fmt.Sprintf("=%s", f.DefaultValuePlaceholder))
	}

	for _, sz := range f.Aliases {
		sb.WriteString(delimChar)
		sb.WriteString("--")
		sb.WriteString(sz)
	}
	return sb.String()
}

// Delete removes myself from the command owner.
func (f *Flag) Delete() {
	if f == nil || f.owner == nil {
		return
	}

	for i, cc := range f.owner.Flags {
		if f == cc {
			f.owner.Flags = append(f.owner.Flags[0:i], f.owner.Flags[i+1:]...)
			return
		}
	}
}
