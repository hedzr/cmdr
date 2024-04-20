package cmdr_test

import (
	"testing"

	cmdr "github.com/hedzr/cmdr/v2"
	"github.com/hedzr/cmdr/v2/cli"

	"gopkg.in/hedzr/errors.v3"
)

func TestExecNoRoot(t *testing.T) {
	if err := cmdr.Exec(nil); !errors.Iss(err, cli.ErrEmptyRootCommand) {
		t.Errorf("Error: %v", err)
	}
}

// func TestExecSimple(t *testing.T) {
// 	if err := cmdr.Exec(testdata.BuildCommands(true)); !errors.Iss(err, cli.ErrEmptyRootCommand) {
// 		t.Errorf("Error: %v", err)
// 	}
// }
