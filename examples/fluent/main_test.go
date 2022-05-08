// Copyright © 2020 Hedzr Yeh.

package main

import (
	"testing"

	"github.com/hedzr/cmdr"
)

func Test1(t *testing.T) {
	cmdr.Set("app.testing", true)
	main()
}
