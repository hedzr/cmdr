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
func (s *Flag) GetTriggeredTimes() int {
	return s.times
}

// GetDescZsh temp
func (s *Flag) GetDescZsh() (desc string) {
	desc = s.Description
	if len(desc) == 0 {
		desc = tool.EraseAnyWSs(s.GetTitleZshFlagName())
	}
	// desc = replaceAll(desc, " ", "\\ ")
	desc = reSQ.ReplaceAllString(desc, `*$1*`)
	desc = reBQ.ReplaceAllString(desc, `**$1**`)
	desc = strings.ReplaceAll(desc, ":", "\\:")
	desc = strings.ReplaceAll(desc, "[", "\\[")
	desc = strings.ReplaceAll(desc, "]", "\\]")
	return
}

// GetTitleFlagNames temp
func (s *Flag) GetTitleFlagNames() string {
	return s.GetTitleFlagNamesBy(",")
}

// GetTitleZshFlagName temp
func (s *Flag) GetTitleZshFlagName() (str string) {
	if len(s.Full) > 0 {
		str += "--" + s.Full
	} else if len(s.Short) > 0 {
		str += "-" + s.Short
	}
	return
}

// GetTitleZshFlagShortName temp
func (s *Flag) GetTitleZshFlagShortName() (str string) {
	if len(s.Short) > 0 {
		str += "-" + s.Short
	} else if len(s.Full) > 0 {
		str += "--" + s.Full
	}
	return
}

// GetTitleZshNamesBy temp
func (s *Flag) GetTitleZshNamesBy(delimChar string, allowPrefix bool) (str string) {
	return s.GetTitleZshNamesExtBy(delimChar, allowPrefix, true, true)
}

// GetTitleZshNamesExtBy temp
func (s *Flag) GetTitleZshNamesExtBy(delimChar string, allowPrefix, shortTitleOnly, longTitleOnly bool) (str string) {
	quote := false
	prefix, suffix := "", ""
	if _, ok := s.DefaultValue.(bool); !ok {
		suffix = "="
		//} else if _, ok := s.DefaultValue.(bool); ok {
		//	suffix = "-"
	}
	if allowPrefix && !s.justOnce {
		quote, prefix = true, "*"
	}
	if !longTitleOnly && len(s.Short) > 0 {
		if quote {
			str += "'" + prefix + "-" + s.Short + suffix + "'"
		} else {
			str += prefix + "-" + s.Short + suffix
		}
		if shortTitleOnly {
			return
		}
	}
	if len(s.Full) > 0 {
		if str != "" {
			str += delimChar
		}
		if quote {
			str += "'" + prefix + "--" + s.Full + suffix + "'"
		} else {
			str += "'" + prefix + "--" + s.Full + suffix + "'"
		}
	}
	return
}

// GetTitleZshFlagNamesArray temp
func (s *Flag) GetTitleZshFlagNamesArray() (ary []string) {
	if len(s.Short) == 1 || len(s.Short) == 2 {
		if len(s.DefaultValuePlaceholder) > 0 {
			ary = append(ary, "-"+s.Short+"=") // +s.DefaultValuePlaceholder)
		} else {
			ary = append(ary, "-"+s.Short)
		}
	}
	if len(s.Full) > 0 {
		if len(s.DefaultValuePlaceholder) > 0 {
			ary = append(ary, "--"+s.Full+"=") // +s.DefaultValuePlaceholder)
		} else {
			ary = append(ary, "--"+s.Full)
		}
	}
	return
}

// GetTitleFlagNamesBy temp
func (s *Flag) GetTitleFlagNamesBy(delimChar string) string {
	return s.GetTitleFlagNamesByMax(delimChar, len(s.Short))
}

// GetTitleFlagNamesByMax temp
func (s *Flag) GetTitleFlagNamesByMax(delimChar string, maxShort int) string {
	var sb strings.Builder

	if len(s.Short) == 0 {
		// if no flag.Short,
		sb.WriteString(strings.Repeat(" ", maxShort))
	} else {
		sb.WriteRune('-')
		sb.WriteString(s.Short)
		sb.WriteString(delimChar)
		if len(s.Short) < maxShort {
			sb.WriteString(strings.Repeat(" ", maxShort-len(s.Short)))
		}
	}

	if len(s.Short) == 0 {
		sb.WriteRune(' ')
		sb.WriteRune(' ')
	}
	sb.WriteRune(' ')
	sb.WriteString("--")
	sb.WriteString(s.Full)
	if len(s.DefaultValuePlaceholder) > 0 {
		// str += fmt.Sprintf("=\x1b[2m\x1b[%dm%s\x1b[0m", DarkColor, s.DefaultValuePlaceholder)
		sb.WriteString(fmt.Sprintf("=%s", s.DefaultValuePlaceholder))
	}

	for _, sz := range s.Aliases {
		sb.WriteString(delimChar)
		sb.WriteString("--")
		sb.WriteString(sz)
	}
	return sb.String()
}
