// Copyright © 2020 Hedzr Yeh.

package tool

import (
	"sort"
	"strings"
	"testing"
)

//  - app.env
//  - app.micro-service.all
//  - app.micro-service.id
//  - app.micro-service.list
//  - app.micro-service.money
//  - app.micro-service.name
//  - app.micro-service.retry
//  - app.micro-service.tags
//  - app.micro-service.tags.add
//  - app.mx
//  - app.mx.1
//  - app.mx.1.2
//  - app.mx.2
//  - app.mx-test.fluent
//  - app.server
//  - app.todd
//  - app.todd.xx
var tests = []string{
	"app.server",
	"app.mx-test.fluent",
	"app.env",
	"app.mx",
	"app.todd",
	"app.micro-service.tags.add",
	"app.micro-service.all",
	"app.micro-service.id",
	"app.todd.xx",
	"app.mx.1.2",
	"app.mx.2",
	"app.micro-service.list",
	"app.micro-service.money",
	"app.micro-service.name",
	"app.micro-service.retry",
	"app.micro-service.tags",
	"app.mx.1",
}

func TestSort00(t *testing.T) {
	for _, tt := range []struct {
		s      []string
		expect bool
	}{
		{[]string{"app.micro-service.name", "app.mx-test.fluent"}, true},
		{[]string{"mx", "mx.1"}, true},
		{[]string{"mx.1", "mx-test"}, true},
		{[]string{"mx", "mx-test"}, true},
		{[]string{"app.mx.1.2", "app.mx"}, false},
		{[]string{"app.mx.1.2", "app.mx.2"}, true},
	} {
		ks := byDottedSlice(tt.s)
		ret := ks.Less(0, 1)
		if ret != tt.expect {
			t.Fatalf(" for %v, wanted %v but got %v", ks, tt.expect, ret)
		}
	}
}

func TestSortByDottedSlice(t *testing.T) {
	ks := byDottedSlice(tests)

	sort.Sort(ks)

	for _, s := range ks {
		t.Logf("  - %v", s)
	}
}

func TestSortByDottedSlice2(t *testing.T) {
	ks := tests

	SortAsDottedSlice(ks)

	for _, s := range ks {
		t.Logf("  - %v", s)
	}

	t.Logf("  %v", strings.Compare("mx", "mx-test"))
}

func TestSortByDottedSlice3(t *testing.T) {
	ks := tests
	SortAsDottedSlice(ks)
	for _, s := range ks {
		t.Logf("  - %v", s)
	}

	t.Logf("  %v", strings.Compare("mx", "mx.1"))
	t.Logf("  %v", strings.Compare("mx", "mx-test"))
	t.Logf("  %v", strings.Compare("mx.1", "mx-test"))
}

func TestSortByDottedSliceRev(t *testing.T) {
	ks := tests

	SortAsDottedSliceReverse(ks)

	for _, s := range ks {
		t.Logf("  - %v", s)
	}
}
