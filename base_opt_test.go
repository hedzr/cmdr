/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"github.com/hedzr/cmdr"
	"testing"
)

func TestHasParent(t *testing.T) {
	s := cmdr.BaseOpt{
		Name:  "A",
		Short: "A",
		Full:  "Abcuse",
	}
	if s.HasParent() {
		t.Failed()
	}
	if s.GetTitleNames() != "A, Abcuse" {
		t.Failed()
	}
}

func TestClone1(t *testing.T) {
	var aa = "dsajkld"
	var b int

	// cmdr.Clone(b, aa)

	cmdr.Clone(b, &aa)

	cmdr.Clone(&b, &aa)

	cmdr.Clone(&b, nil)

	var c, d bool
	cmdr.Clone(&c, &d)

	var e, f int
	cmdr.Clone(&e, &f)
	var e1, f1 int8
	f1 = 1
	cmdr.Clone(&e1, &f1)
	if e1 != 1 {
		t.Failed()
	}
	var e2, f2 int16
	cmdr.Clone(&e2, &f2)
	var e3, f3 int32
	cmdr.Clone(&e3, &f3)
	var e4, f4 int64
	f4 = 9
	cmdr.Clone(&e4, &f4)
	if e1 != 9 {
		t.Failed()
	}

	var g, h string
	cmdr.Clone(&g, &h)
}

func TestClone2(t *testing.T) {

}
