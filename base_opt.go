/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"strings"
)

// HasParent detects whether owner is available or not
func (s *BaseOpt) HasParent() bool {
	return s.owner != nil
}

// GetTitleName temp
func (s *BaseOpt) GetTitleName() string {
	if len(s.Name) != 0 {
		return s.Name
	}
	return s.Full
}

// GetTitleNamesArray temp
func (s *BaseOpt) GetTitleNamesArray() []string {
	var a []string
	if len(s.Short) != 0 {
		a = append(a, s.Short)
	}
	if len(s.Full) > 0 {
		a = append(a, s.Full)
	}
	a = append(a, s.Aliases...)
	return a
}

// GetShortTitleNamesArray temp
func (s *BaseOpt) GetShortTitleNamesArray() []string {
	var a []string
	if len(s.Short) != 0 {
		a = append(a, s.Short)
	}
	return a
}

// GetLongTitleNamesArray temp
func (s *BaseOpt) GetLongTitleNamesArray() []string {
	var a []string
	if len(s.Full) > 0 {
		a = append(a, s.Full)
	}
	a = append(a, s.Aliases...)
	return a
}

// GetTitleNames temp
func (s *BaseOpt) GetTitleNames() string {
	return s.GetTitleNamesBy(", ")
}

// GetTitleNamesBy temp
func (s *BaseOpt) GetTitleNamesBy(delimChar string) string {
	var a = s.GetTitleNamesArray()
	str := strings.Join(a, delimChar)
	return str
}
