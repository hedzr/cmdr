/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"strings"
)

func (s *BaseOpt) GetTitleName() string {
	if len(s.Name) != 0 {
		return s.Name
	}
	return s.Full
}

func (s *BaseOpt) GetTitleNamesArray() []string {
	var a []string
	if len(s.Short) != 0 {
		a = append(a, fmt.Sprintf("%s", s.Short))
	}
	if len(s.Full) > 0 {
		a = append(a, s.Full)
	}
	a = append(a, s.Aliases...)
	return a
}

func (s *BaseOpt) GetTitleNames() string {
	return s.GetTitleNamesBy(", ")
}

func (s *BaseOpt) GetTitleNamesBy(delimChar string) string {
	var a []string = s.GetTitleNamesArray()
	str := strings.Join(a, delimChar)
	return str
}

func (s *BaseOpt) GetTitleFlagNames() string {
	return s.GetTitleFlagNamesBy(",")
}

func (s *BaseOpt) GetDescZsh() (desc string) {
	desc = s.Description
	if len(desc) == 0 {
		desc = simpsimp(s.GetTitleZshFlagName())
	}
	// desc = strings.ReplaceAll(desc, " ", "\\ ")
	return
}

func (s *BaseOpt) GetTitleZshFlagName() (str string) {
	if len(s.Full) > 0 {
		str += "--" + s.Full
	} else if len(s.Short) == 1 {
		str += "-" + s.Short
	}
	return
}

func (s *BaseOpt) GetTitleZshFlagNames(delimChar string) (str string) {
	if len(s.Short) == 1 {
		str += "-" + s.Short + delimChar
	}
	if len(s.Full) > 0 {
		str += "--" + s.Full
	}
	return
}

func (s *BaseOpt) GetTitleZshFlagNamesArray() (ary []string) {
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

func (s *BaseOpt) GetTitleFlagNamesBy(delimChar string) string {
	return s.GetTitleFlagNamesByMax(delimChar, -1)
}

func (s *BaseOpt) GetTitleFlagNamesByMax(delimChar string, maxCount int) string {
	var a []string = s.GetTitleNamesArray()
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
				// str += fmt.Sprintf("=\x1b[2m\x1b[%dm%s\x1b[0m", darkColor, s.DefaultValuePlaceholder)
				str += fmt.Sprintf("=%s", s.DefaultValuePlaceholder)
			}
		} else {
			str += delimChar + " --" + sz
		}
	}
	return str
}
