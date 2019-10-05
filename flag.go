/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"strings"
)

// GetTriggeredTimes returns the matched times
func (s *Flag) GetTriggeredTimes() int {
	return s.times
}

// GetTitleFlagNames temp
func (s *Flag) GetTitleFlagNames() string {
	return s.GetTitleFlagNamesBy(",")
}

// GetDescZsh temp
func (s *Flag) GetDescZsh() (desc string) {
	desc = s.Description
	if len(desc) == 0 {
		desc = eraseAnyWSs(s.GetTitleZshFlagName())
	}
	// desc = strings.ReplaceAll(desc, " ", "\\ ")
	return
}

// GetTitleZshFlagName temp
func (s *Flag) GetTitleZshFlagName() (str string) {
	if len(s.Full) > 0 {
		str += "--" + s.Full
	} else if len(s.Short) == 1 {
		str += "-" + s.Short
	}
	return
}

// GetTitleZshFlagNames temp
func (s *Flag) GetTitleZshFlagNames(delimChar string) (str string) {
	if len(s.Short) == 1 {
		str += "-" + s.Short + delimChar
	}
	if len(s.Full) > 0 {
		str += "--" + s.Full
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
	return s.GetTitleFlagNamesByMax(delimChar, -1)
}

// GetTitleFlagNamesByMax temp
func (s *Flag) GetTitleFlagNamesByMax(delimChar string, maxCount int) string {
	var a = s.GetTitleNamesArray()
	var str string

	if len(s.Short) == 0 {
		// if no flag.Short,
		a = append([]string{""}, a...)
	}

	for ix, sz := range a {
		if ix == 0 {
			if len(sz) == 0 {
				// if no flag.Short,
				str += "  "
			} else {
				str += "-" + sz
			}
		} else if ix == 1 {
			if len(strings.TrimSpace(str)) == 0 {
				// if no flag.Short,
				str += " "
			} else {
				str += delimChar
			}
			if len(str) < 4 {
				// align between -nv and -v
				str += " "
			}
			str += " --" + sz
			if len(s.DefaultValuePlaceholder) > 0 {
				// str += fmt.Sprintf("=\x1b[2m\x1b[%dm%s\x1b[0m", DarkColor, s.DefaultValuePlaceholder)
				str += fmt.Sprintf("=%s", s.DefaultValuePlaceholder)
			}
		} else {
			str += delimChar + " --" + sz
		}
	}
	return str
}
