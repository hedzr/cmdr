// Copyright © 2020 Hedzr Yeh.

package main

import (
	"github.com/hedzr/cmdr"
	"testing"
)

func Test1(t *testing.T) {
	cmdr.Set("app.testing", true)
	main()
}
