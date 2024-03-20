package exec

import "testing"

func TestSplitCommandString(t *testing.T) {
	in := `bash -c 'echo hello world!'`
	out := SplitCommandString(in, '"', '\'')
	t.Log(out)
}
