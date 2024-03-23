package builder

import (
	"testing"
)

func TestNew(t *testing.T) {
	v := New(nil)
	t.Logf("v = %v", v)
}
