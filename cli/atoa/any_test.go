package atoa

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hedzr/cmdr/v2/pkg/times"
)

func TestParse(t *testing.T) { //nolint:revive
	// t.Logf("time parsed: %v", times.MustSmartParseTime("11:52"))

	for i, c := range []struct {
		src    string
		meme   any
		expect any
	}{
		{"1979-1-29 11:52:00.6789", &time.Time{}, times.MustSmartParseTimePtr("1979-1-29 11:52:0.6789")},
		{"1979-1-29 11:52:00.6789", time.Time{}, times.MustSmartParseTime("1979-1-29 11:52:0.6789")},
		{"a=1,b=2,", &aS{}, &aS{map[string]int{"a": 1, "b": 2}}},

		// {"apple=1, banana=2, orange=3", map[string]int{}, map[string]int{"apple": 1, "banana": 2, "orange": 3}},
		// {"[  { a : [ 8 , 9 ] , b : [ 9 , 2 ] , c : [ 7 , 5 ] } ,   { e : [ 1 , 3 ] , f : [ 6 , -1 ] }   ]", []map[string][]int{},
		// 	[]map[string][]int{
		// 		{"a": {8, 9}, "b": {9, 2}, "c": {7, 5}},
		// 		{"e": {1, 3}, "f": {6, -1}},
		// 	},
		// },
		// {"[  {a:[8,9],b:[9,2],c:[7,5]},   {e:[1,3],f:[6,-1]}   ]", []map[string][]int{},
		// 	[]map[string][]int{
		// 		{"a": {8, 9}, "b": {9, 2}, "c": {7, 5}},
		// 		{"e": {1, 3}, "f": {6, -1}},
		// 	},
		// },
		// {"{a:[8,9],b:[9,2],c:[7,5]}", map[string][]int{}, map[string][]int{"a": {8, 9}, "b": {9, 2}, "c": {7, 5}}},
		// {"[8,9,7]", [2]int{}, [2]int{8, 9}},

		{"8,9,7", []int{}, []int{8, 9, 7}},
		{"[8,9,7]", []int{}, []int{8, 9, 7}},
		{"apple=1, banana=2, orange=3", map[string]int{}, map[string]int{"apple": 1, "banana": 2, "orange": 3}},
		{"apple:1, banana: 2, orange: 3", map[string]int{}, map[string]int{"apple": 1, "banana": 2, "orange": 3}},
		{"{apple:1, banana: 2, orange: 3}", map[string]int{}, map[string]int{"apple": 1, "banana": 2, "orange": 3}},

		{"on", false, true},
		{"string", "", "string"},
		{`"string"`, "", "string"},
		{"3.14159", 1.1, 3.14159},
		{"0xffff", 1, 0xffff},

		{"[8,9,7]", [2]int{}, [2]int{8, 9}},

		{"[[8,9],[9,2],[7,5]]", [][]int{}, [][]int{{8, 9}, {9, 2}, {7, 5}}},

		{"{a:[8,9],b:[9,2],c:[7,5]}", map[string][]int{}, map[string][]int{"a": {8, 9}, "b": {9, 2}, "c": {7, 5}}},

		{`[{a:[8,9],b:[9,2],c:[7,5]},{e:[1,3],f:[6,-1]}  ]`,
			[]map[string][]int{},
			[]map[string][]int{
				{"a": {8, 9}, "b": {9, 2}, "c": {7, 5}},
				{"e": {1, 3}, "f": {6, -1}},
			},
		},

		{`[  {a:[8,9],b:[9,2],c:[7,5]},   {e:[1,3],f:[6,-1]}   ]`,
			[]map[string][]int{},
			[]map[string][]int{
				{"a": {8, 9}, "b": {9, 2}, "c": {7, 5}},
				{"e": {1, 3}, "f": {6, -1}},
			},
		},

		{`[  { a : [ 8 , 9 ] , b : [ 9 , 2 ] , c : [ 7 , 5 ] } ,   { e : [ 1 , 3 ] , f : [ 6 , -1 ] }   ]`,
			[]map[string][]int{},
			[]map[string][]int{
				{"a": {8, 9}, "b": {9, 2}, "c": {7, 5}},
				{"e": {1, 3}, "f": {6, -1}},
			},
		},

		{`[  { "a" : [ 8 , 9 ] , b : [ 9 , 2 ] , c : [ 7 , 5 ] } ,   { e : [ 1 , 3 ] , f : [ 6 , -1 ] }   ]`, []map[string][]int{},
			[]map[string][]int{
				{"a": {8, 9}, "b": {9, 2}, "c": {7, 5}},
				{"e": {1, 3}, "f": {6, -1}},
			},
		},

		{"1979-1-29 11:52:00.6789", time.Now(), times.MustSmartParseTime("1979-1-29 11:52:0.6789")},
		{"11:52", time.Now(), times.MustSmartParseTime("11:52")},
		{"3s300ms", time.Duration(1), 3*time.Second + 300*time.Millisecond},

		// {`first: {apple: 1, banana: 2, orange: 3}`, map[string]map[string]int{}, map[string]map[string]int{"first": {"apple": 1, "banana": 2, "orange": 3}}},
	} {
		t.Log()
		t.Log()
		t.Log()
		t.Logf("--------------------- %d. SRC = %q", i, c.src)

		val, err := Parse(c.src, c.meme)
		if err != nil {
			t.Fatalf("%5d. test failed, err: %v", i, err)
		}
		if !reflect.DeepEqual(val, c.expect) {
			t.Fatalf("FAILED PARSE: got %v, but expecting %v", val, c.expect)
		}
	}
}

type aS struct {
	v map[string]int
}

func (a *aS) MarshalText() (text []byte, err error) {
	if a.v != nil {
		var sb bytes.Buffer
		for k, v := range a.v {
			_, _ = sb.WriteString(fmt.Sprintf("%s=%d,", k, v))
		}
		text = sb.Bytes()
	}
	return
}

func (a *aS) UnmarshalText(text []byte) error {
	for _, ln := range strings.Split(string(text), ",") {
		pos := strings.Index(ln, "=")
		if pos < 0 {
			continue
		}
		if a.v == nil {
			a.v = make(map[string]int)
		}
		a.v[ln[:pos]] = MustParse(ln[pos+1:], 1).(int) //nolint:revive
	}
	return nil
}
