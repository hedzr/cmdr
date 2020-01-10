/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr_test

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"math"
	"math/cmplx"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestHumanReadableSizes(t *testing.T) {

	s := &cmdr.Options{}

	for _, tc := range []struct {
		src      string
		expected uint64
	}{
		{"1234", 1234},
		// just for go1.13+: {"1_234", 1234},
		// {"1,234", 1234},
		{"1.234 kB", 1234},
		{"238.674052 MB", 238674052},
		{"543 B", 543},
		{"8k", 8000},
		{"8GB", 8 * 1000 * 1000 * 1000},
		{"8TB", 8 * 1000 * 1000 * 1000 * 1000},
		{"8pB", 8 * 1000 * 1000 * 1000 * 1000 * 1000},
		{"8EB", 8 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000},
	} {
		tgt := s.FromKilobytes(tc.src)
		if tgt != tc.expected {
			t.Fatalf("StripQuotes(%q): expect %v, but got %v", tc.src, tc.expected, tgt)
		}
	}

	for _, tc := range []struct {
		src      string
		expected uint64
	}{
		{"1234", 1234},
		// just for go1.13+: {"1_234", 1234},
		// {"1,234", 1234},
		{"1.234 kB", 1263},
		{"238.674052 MB", 250267882},
		{"543 B", 543},
		{"8k", 8192},
		{"640K", 655360},
		{"8GB", 8 * 1024 * 1024 * 1024},
		{"8TB", 8 * 1024 * 1024 * 1024 * 1024},
		{"8pB", 8 * 1024 * 1024 * 1024 * 1024 * 1024},
		{"8EB", 8 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024},
	} {
		tgt := s.FromKibibytes(tc.src)
		if tgt != tc.expected {
			t.Fatalf("StripQuotes(%q): expect %v, but got %v", tc.src, tc.expected, tgt)
		}
	}

}

func TestComplex(t *testing.T) {
	// see also: https://golang.org/pkg/fmt/
	fmt.Printf("%v\n", (3 + 4i))
	fmt.Printf("%g\n", (3 + 4i))
	fmt.Printf("%.3g\n", (3.14 + 4.73i))
	fmt.Printf("%.3g\n", (3.14 - 4.73i))
	fmt.Printf("%.2f\n", cmplx.Abs(3+4i))
	fmt.Printf("%.2f\n", cmplx.Exp(1i*math.Pi)+1)

	fmt.Printf("> %g\n", complex64(cmdr.ParseComplex("3 + 4i")))
	fmt.Printf("> %g\n", cmdr.ParseComplex("3+4i"))
	fmt.Printf("> %g\n", complex64(cmdr.ParseComplex("3.14-4.73i")))
	fmt.Printf("> %g\n", cmdr.ParseComplex("3.14+4.73i"))

	_, _ = cmdr.ParseComplexX("3+4i")
	_, _ = cmdr.ParseComplexX("3.14+4.73i")
	_, _ = cmdr.ParseComplexX("3+4qi")
	_, _ = cmdr.ParseComplexX("3z+4qi")
	_, _ = cmdr.ParseComplexX("3x+4t")
	_, _ = cmdr.ParseComplexX("3+4")
	_, _ = cmdr.ParseComplexX("3")
	
	r, theta := cmplx.Polar(2i)
	fmt.Printf("r: %.2f, θ: %.2f*π\n", r, theta/math.Pi)

	// 5.00
	// (0.00+0.00i)
	// r: 2.00, θ: 0.50*π

	// strconv.new
	const a = complex(100, 8)
	const b = complex(8, 100)

	fmt.Println("Complex number a : ", a)
	fmt.Println("Complex number b : ", b)
	fmt.Println("Get the real part of complex number a : ", real(a))
	fmt.Println("Get the imaginary part of complex number b : ", imag(a))

	conjugate := cmplx.Conj(a)

	fmt.Println("Complex number a's conjugate : ", conjugate)

	c := a + b

	fmt.Println("a + b complex number : ", c)
	fmt.Println("Cosine of complex number b : ", cmplx.Cos(b))

	// see https://golang.org/pkg/math/cmplx/
	// for more functions such as sine, log, exponential

	var err error
	var i64 float64
	i64, err = strconv.ParseFloat("0x96c1", 64)
	t.Logf("i: %v | err: %v", i64, err)
}

func TestNumberParsing(t *testing.T) {
	var err error
	var i64 int64
	var f64 float64

	i64, err = strconv.ParseInt("0x96c1", 0, 64)
	t.Logf("i: %v | err: %v", i64, err)
	i64, err = strconv.ParseInt("0xffb969d28651e43c", 0, 64)
	t.Logf("i: %v | err: %v", i64, err) // i: 9223372036854775807 | err: strconv.ParseInt: parsing "0xffb969d28651e43c": value out of range
	i64, err = strconv.ParseInt("0x7fb969d28651e43c", 0, 64)
	t.Logf("i: %v | err: %v", i64, err)
	i64, err = strconv.ParseInt("011000011b", 0, 64)
	t.Logf("i: %v | err: %v", i64, err) // i: 0 | err: strconv.ParseInt: parsing "011000011b": invalid syntax
	i64, err = strconv.ParseInt("0b00101101", 0, 64)
	t.Logf("i: %v | err: %v", i64, err)
	i64, err = strconv.ParseInt("0755", 0, 64)
	t.Logf("i: %v | err: %v", i64, err)
	i64, err = strconv.ParseInt("0o755", 0, 64)
	t.Logf("i: %v | err: %v", i64, err)
	i64, err = strconv.ParseInt("123_456", 0, 64)
	t.Logf("i: %v | err: %v", i64, err)

	i64, err = strconv.ParseInt("0x1p-2", 0, 64)
	t.Logf("i: %v | err: %v", i64, err) // i: 0 | err: strconv.ParseInt: parsing "0x1p-2": invalid syntax

	// errors_test.go:39: i: 38593 | err: <nil>
	// errors_test.go:41: i: 9223372036854775807 | err: strconv.ParseInt: parsing "0xffb969d28651e43c": value out of range
	// errors_test.go:43: i: 9203503666425881660 | err: <nil>
	// errors_test.go:45: i: 0 | err: strconv.ParseInt: parsing "011000011b": invalid syntax
	// errors_test.go:47: i: 45 | err: <nil>
	// errors_test.go:49: i: 493 | err: <nil>
	// errors_test.go:51: i: 493 | err: <nil>
	// errors_test.go:53: i: 123456 | err: <nil>
	// errors_test.go:56: i: 0 | err: strconv.ParseInt: parsing "0x1p-2": invalid syntax

	f64, err = strconv.ParseFloat("3.4e39", 64)
	t.Logf("f: %v | err: %v", f64, err)
	f64, err = strconv.ParseFloat("3.1415926535897935", 64)
	t.Logf("f: %v | err: %v", f64, err)

	v := "3.1415926535"
	if s, err := strconv.ParseFloat(v, 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	if s, err := strconv.ParseFloat(v, 64); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	if s, err := strconv.ParseFloat("NaN", 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	// ParseFloat is case insensitive
	if s, err := strconv.ParseFloat("nan", 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	if s, err := strconv.ParseFloat("inf", 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	if s, err := strconv.ParseFloat("+Inf", 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	if s, err := strconv.ParseFloat("-Inf", 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	if s, err := strconv.ParseFloat("-0", 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
	if s, err := strconv.ParseFloat("+0", 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)
	}
}

func TestFinds(t *testing.T) {
	t.Log("finds")
	cmdr.InternalResetWorker()
	cmdr.ResetOptions()

	cmdr.Set("no-watch-conf-dir", true)

	// copyRootCmd = rootCmdForTesting
	var rootCmdX = &cmdr.RootCommand{
		Command: cmdr.Command{
			BaseOpt: cmdr.BaseOpt{
				Name: "consul-tags",
			},
		},
	}
	t.Log("rootCmdForTesting", rootCmdX)

	var commands = []string{
		"consul-tags --help -q",
	}
	for _, cc := range commands {
		os.Args = strings.Split(cc, " ")
		cmdr.SetInternalOutputStreams(nil, nil)

		if err := cmdr.Exec(rootCmdX); err != nil {
			t.Fatal(err)
		}
	}

	if cmdr.InTesting() {
		cmdr.FindSubCommand("generate", nil)
		cmdr.FindFlag("generate", nil)
		cmdr.FindSubCommandRecursive("generate", nil)
		cmdr.FindFlagRecursive("generate", nil)
	} else {
		t.Log("noted")
	}
	resetOsArgs()
}
