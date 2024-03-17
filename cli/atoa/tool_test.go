package atoa

import "testing"

func TestEat(t *testing.T) {
	str := ":Since:"
	res, ate := eat(str, ':')
	t.Logf("eat() result: %q", res)
	if ate {
		res, ate = eatTail(res, ':')
		t.Logf("eatTail() result: %q", res)
		if !ate {
			t.Fail()
		}
	}
}
