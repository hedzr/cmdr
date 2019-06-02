/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"github.com/hedzr/cmdr"
	"testing"
)

func TestIsDirectory(t *testing.T) {
	if yes, err := cmdr.IsDirectory("./conf.d1"); yes {
		t.Fatal(err)
	}
	if yes, err := cmdr.IsDirectory("./ci"); !yes {
		t.Fatal(err)
	}
	if yes, err := cmdr.IsRegularFile("./doc1.golang"); yes {
		t.Fatal(err)
	}
	if yes, err := cmdr.IsRegularFile("./doc.go"); !yes {
		t.Fatal(err)
	}
}
