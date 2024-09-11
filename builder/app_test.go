package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/cmdr/v2/cli"
)

func TestAppS_AddCmd(t *testing.T) {
	a := &appS{}
	a.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("ask", "a")
	})
	assert.Equal(t, a.root.Long, "ask")
	assert.Equal(t, a.root.Short, "a")

	a.AddCmd(func(b cli.CommandBuilder) {
		b.Titles("bunny", "b")
	})
	child := a.root.SubCommands()[0]
	assert.Equal(t, child.Long, "bunny")
	assert.Equal(t, child.Short, "b")
}

func TestAppS_Run(t *testing.T) {
	a := &appS{inCmd: 1}
	err := a.Run()
	assert.Error(t, err)

	a = &appS{inFlg: 2}
	err = a.Run()
	assert.Error(t, err)

	a = &appS{}
	err = a.Run()
	assert.Equal(t, err, cli.ErrEmptyRootCommand)
}
