package cli

import (
	"testing"
)

func TestBaseOpt_GetDottedPath(t *testing.T) {
	root := rootCmdForTesting()
	for i, tc := range []struct {
		dp string
	}{
		// {"microservices.tags"},
		{""},
		{"server"},
		{"server.head"},
		{"server.tail"},
		{"server.enum"},
		{"server.retry"},
		{"server.start"},
		{"server.start.foreground"},
		{"server.stop"},
		{"server.restart"},
		{"kvstore"},
		{"kvstore.backup"},
		{"kvstore.backup.output"},
		{"microservices.tags"},
		{"microservices.tags.list"},
		{"microservices.tags.add.list"},
		{"microservices.tags.modify"},
		{"microservices.tags.modify.add"},
		{"microservices.tags.modify.rm"},
		{"microservices.tags.modify.ued"},
		{"microservices.tags.toggle.address"},
	} {
		c, f := root.DottedPathToCommandOrFlag(tc.dp)
		if f != nil {
			if s := f.GetDottedPath(); s != tc.dp {
				t.Fatalf("%5d. GetDottedPath expected %q, got %q", i, tc.dp, s)
			}
		} else if c != nil {
			if s := c.GetDottedPath(); s != tc.dp {
				t.Fatalf("%5d. GetDottedPath expected %q, got %q", i, tc.dp, s)
			}
		} else {
			t.Fatalf("%5d. expected %q, but cmd/flg not found", i, tc.dp)
		}
	}
}
