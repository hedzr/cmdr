/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"github.com/hedzr/cmdr"
	"os"
	"strings"
	"testing"
)

func TestHumanReadableSizes(t *testing.T) {

	s := &cmdr.Options{}

	for _, tc := range []struct {
		src      string
		expected uint64
	}{
		{"1234", 1234},
		// just for go1.13+: {"1_234", 1234},
		// {"1,234", 1234},
		{"1.234 kB", 1234},
		{"238.674052 MB", 238674052},
		{"543 B", 543},
		{"8k", 8000},
		{"8GB", 8 * 1000 * 1000 * 1000},
		{"8TB", 8 * 1000 * 1000 * 1000 * 1000},
		{"8pB", 8 * 1000 * 1000 * 1000 * 1000 * 1000},
		{"8EB", 8 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000},
	} {
		tgt := s.FromKilobytes(tc.src)
		if tgt != tc.expected {
			t.Fatalf("StripQuotes(%q): expect %v, but got %v", tc.src, tc.expected, tgt)
		}
	}

	for _, tc := range []struct {
		src      string
		expected uint64
	}{
		{"1234", 1234},
		// just for go1.13+: {"1_234", 1234},
		// {"1,234", 1234},
		{"1.234 kB", 1263},
		{"238.674052 MB", 250267882},
		{"543 B", 543},
		{"8k", 8192},
		{"640K", 655360},
		{"8GB", 8 * 1024 * 1024 * 1024},
		{"8TB", 8 * 1024 * 1024 * 1024 * 1024},
		{"8pB", 8 * 1024 * 1024 * 1024 * 1024 * 1024},
		{"8EB", 8 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024},
	} {
		tgt := s.FromKibibytes(tc.src)
		if tgt != tc.expected {
			t.Fatalf("StripQuotes(%q): expect %v, but got %v", tc.src, tc.expected, tgt)
		}
	}

}

func TestFinds(t *testing.T) {
	t.Log("finds")
	cmdr.InternalResetWorker()
	cmdr.ResetOptions()

	cmdr.Set("no-watch-conf-dir", true)

	// copyRootCmd = rootCmdForTesting
	var rootCmdX = &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name: "consul-tags",
			},
		},
	}
	t.Log("rootCmdForTesting", rootCmdX)

	var commands = []string{
		"consul-tags --help -q",
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc, " ")
		cmdr.SetInternalOutputStreams(nil, nil)

		if err := cmdr.Exec(rootCmdX); err != nil {
			t.Fatal(err)
		}
	}

	if cmdr.InTesting() {
		cmdr.FindSubCommand("generate", nil)
		cmdr.FindFlag("generate", nil)
		cmdr.FindSubCommandRecursive("generate", nil)
		cmdr.FindFlagRecursive("generate", nil)
	} else {
		t.Log("noted")
	}
	resetOsArgs()
}
