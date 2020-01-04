// Copyright © 2020 Hedzr Yeh.

package cmdr

import "fmt"

func defaultOnSwitchCharHit(parsed *Command, switchChar string, args []string) (err error) {
	if parsed != nil {
		fmt.Printf("the last parsed command is %q - %q\n", parsed.GetTitleNames(), parsed.Description)
	}
	fmt.Printf("SwitchChar FOUND: %v\nremains: %v\n\n", switchChar, args)
	return nil // ErrShouldBeStopException
}

func defaultOnPasssThruCharHit(parsed *Command, switchChar string, args []string) (err error) {
	if parsed != nil {
		fmt.Printf("the last parsed command is %q - %q\n", parsed.GetTitleNames(), parsed.Description)
	}
	fmt.Printf("PassThrough flag FOUND: %v\nremains: %v\n\n", switchChar, args)
	return nil // ErrShouldBeStopException
}

func emptyUnknownOptionHandler(isFlag bool, title string, cmd *Command, args []string) (fallbackToDefaultDetector bool) {
	return false
}
