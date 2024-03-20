package dir

import (
	"testing"
)

func TestIsMatch(t *testing.T) {
	output := IsWildMatch("aa", "aa")
	expectTrue(t, output)

	output = IsWildMatch("aaaa", "*")
	expectTrue(t, output)

	output = IsWildMatch("ab", "a?")
	expectTrue(t, output)

	output = IsWildMatch("adceb", "*a*b")
	expectTrue(t, output)

	output = IsWildMatch("aa", "a")
	expectFalse(t, output)

	output = IsWildMatch("mississippi", "m??*ss*?i*pi")
	expectFalse(t, output)

	output = IsWildMatch("acdcb", "a*c?b")
	expectFalse(t, output)

	output = IsWildMatch(".config/1.x/1", "*/1.x/*")
	expectTrue(t, output)
}

func expectFalse(t *testing.T, text any) {
	if ToBool(text) {
		t.Fatalf("expecting false, but got %q", text)
	}
}

func expectTrue(t *testing.T, text any) {
	if !ToBool(text) {
		t.Fatalf("expecting true, but got %q", text)
	}
}
