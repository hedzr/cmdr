package hs

import (
	"context"
	"strings"
	"testing"
)

func TestHelpSystem_helpCmd(t *testing.T) {
	ctx := context.Background()
	worker, root, err := rootCmdForTesting()
	if err != nil {
		t.Fatalf("rootCmdForTesting() failed: %v", err)
	}

	for _, tc := range []struct {
		args []string
	}{
		{args: []string{"help", "server", "start"}},
		{args: []string{"help", "server", "start", "foreground"}},
	} {
		hs := &HelpSystem{
			worker: worker,
			cmd:    root,
			args:   tc.args,
		}
		var sb strings.Builder
		err = hs.helpCmd(ctx, tc.args[1:], &sb)
		if err != nil {
			t.Fatalf("Run() failed: %v", err)
		} else {
			t.Log(sb.String())
		}
	}
}

func TestStr(t *testing.T) {
	str := "1\n22\n"
	str = strings.ReplaceAll(str, "\n", "\r\n")
	t.Log(str)
}
