/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"github.com/hedzr/cmdr/tool"
	"regexp"
	"strings"
)

// HasParent detects whether owner is available or not
func (s *BaseOpt) HasParent() bool {
	return s.owner != nil
}

// GetTitleName returns name/full/short string
func (s *BaseOpt) GetTitleName() string {
	if len(s.Name) != 0 {
		return s.Name
	}
	if len(s.Full) > 0 {
		return s.Full
	}
	if len(s.Short) > 0 {
		return s.Short
	}
	// for _, ss := range s.Aliases {
	// 	return ss
	// }
	return ""
}

// GetTitleNamesArray returns short,full,aliases names
func (s *BaseOpt) GetTitleNamesArray() []string {
	a := s.GetTitleNamesArrayMainly()
	a = uniAddStrs(a, s.Aliases...)
	return a
}

// GetTitleNamesArrayMainly returns short,full names
func (s *BaseOpt) GetTitleNamesArrayMainly() []string {
	var a []string
	if len(s.Short) != 0 {
		a = uniAddStr(a, s.Short)
	}
	if len(s.Full) > 0 {
		a = uniAddStr(a, s.Full)
	}
	return a
}

// GetShortTitleNamesArray returns short name as an array
func (s *BaseOpt) GetShortTitleNamesArray() []string {
	var a []string
	if len(s.Short) != 0 {
		a = uniAddStr(a, s.Short)
	}
	return a
}

// GetLongTitleNamesArray returns long name and aliases as an array
func (s *BaseOpt) GetLongTitleNamesArray() []string {
	var a []string
	if len(s.Full) > 0 {
		a = uniAddStr(a, s.Full)
	}
	a = uniAddStrs(a, s.Aliases...)
	return a
}

// GetTitleNames return the joint string of short,full,aliases names
func (s *BaseOpt) GetTitleNames() string {
	return s.GetTitleNamesBy(", ")
}

// GetTitleNamesBy returns the joint string of short,full,aliases names
func (s *BaseOpt) GetTitleNamesBy(delimChar string) string {
	var a = s.GetTitleNamesArray()
	str := strings.Join(a, delimChar)
	return str
}

// GetTitleZshNames temp
func (s *BaseOpt) GetTitleZshNames() string {
	var a = s.GetTitleNamesArrayMainly()
	str := strings.Join(a, ",")
	return str
}

// GetTitleZshNamesBy temp
func (s *BaseOpt) GetTitleZshNamesBy(delimChar string) (str string) {
	if len(s.Short) != 0 {
		str += s.Short + delimChar
	}
	if len(s.Full) > 0 {
		str += s.Full
	}
	return
}

// GetDescZsh temp
func (s *BaseOpt) GetDescZsh() (desc string) {
	desc = s.Description
	if len(desc) == 0 {
		desc = tool.EraseAnyWSs(s.GetTitleZshNames())
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

var (
	reSQnp = regexp.MustCompile(`'`)
	reBQnp = regexp.MustCompile("`")
	reSQ   = regexp.MustCompile(`'(.*?)'`)
	reBQ   = regexp.MustCompile("`(.*?)`")
	reULs  = regexp.MustCompile("_+")
)
