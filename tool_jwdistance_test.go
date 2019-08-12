/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"math"
	"testing"
)

var distanceTestsJW = []struct {
	first  string
	second string
	wanted float64
}{
	{"mame", "name", 0.7222222222222222},
	{"mv", "mx", 0.6666666666666666},
	{"mv", "mx-test", 0.5476190476190476},
	{"mv", "micro-service", 0.7396449704142012},
	{"update-cc", "update-cv", 0.9851851851851852},
	{"AL", "AL", 1},
	{"MARTHA", "MARHTA", 0.9611111111111111},
	{"JONES", "JOHNSON", 0.8323809523809523},
	{"POTATO", "POTATTO", 0.9761904761904762},
	{"kitten", "sitting", 0.7460317460317460},
	{"MOUSE", "HOUSE", 0.8666666666666667},
}

func TestJaroWinkler(t *testing.T) {
	jw := JaroWinklerDistance()
	for ix, vt := range distanceTestsJW {
		// s1, s2 := "POTATO", "POTATTO"
		d := jw.Calc(vt.first, vt.second)
		t.Logf("%5d. distance of '%v' and '%v': %d",
			ix, vt.first, vt.second, d,
		)
		if d != int(math.Round(vt.wanted*stringMetricFactor)) {
			t.Errorf("wrong distance: for '%v' and '%v', expected distance is %v, but got %v",
				vt.first, vt.second, vt.wanted, d)
		}
	}
}
