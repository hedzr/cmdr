package times

import (
	"testing"
	"time"
)

func TestAddKnownTimeFormats(t *testing.T) {
	src := "1979-01-29 11:52:00.678910129"
	var tm time.Time
	var err error

	for i, c := range []struct {
		src       string
		expecting any
	}{
		{"1979-01-29 11:52:00.678910129", 0},
		{"1979-1-29 11:52:0.67891", 0},
	} {
		tm, err = time.Parse("2006-1-2 15:4:5.999999999", c.src)
		if err != nil {
			t.Fatalf("%5d. time.Parse(%q) failed, err: %v.", i, c.src, err)
		}
		t.Logf("%5d. time: %v", i, tm)
	}

	tm, err = SmartParseTime(src)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("time: %v", tm)
}

func TestShortDur(t *testing.T) {
	h, m, s := 5*time.Hour, 4*time.Minute, 3*time.Second
	ds := []time.Duration{
		h + m + s, h + m, h + s, m + s, h, m, s, 0,
	}

	for _, d := range ds {
		t.Logf("%-16v %-16v %-16v\n", d, shortDur(d), MustParseDuration(shortDur(d)))
	}
}
