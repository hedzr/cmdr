/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/log"
	"gopkg.in/hedzr/errors.v2"
	"io"
	"io/ioutil"
	log2 "log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// ParseComplex converts a string to complex number.
//
// Examples:
//
//    c1 := cmdr.ParseComplex("3-4i")
//    c2 := cmdr.ParseComplex("3.13+4.79i")
func ParseComplex(s string) (v complex128) {
	return a2complexShort(s)
}

// ParseComplexX converts a string to complex number.
// If the string is not valid complex format, return err not nil.
//
// Examples:
//
//    c1 := cmdr.ParseComplex("3-4i")
//    c2 := cmdr.ParseComplex("3.13+4.79i")
func ParseComplexX(s string) (v complex128, err error) {
	return a2complex(s)
}

func a2complexShort(s string) (v complex128) {
	v, _ = a2complex(s)
	return
}

func a2complex(s string) (v complex128, err error) {
	s = strings.TrimSpace(strings.TrimRightFunc(strings.TrimLeftFunc(s, func(r rune) bool {
		return r == '('
	}), func(r rune) bool {
		return r == ')'
	}))

	if i := strings.IndexAny(s, "+-"); i >= 0 {
		rr, ii := s[0:i], s[i:]
		if j := strings.Index(ii, "i"); j >= 0 {
			var ff, fi float64
			ff, err = strconv.ParseFloat(strings.TrimSpace(rr), 64)
			if err != nil {
				return
			}
			fi, err = strconv.ParseFloat(strings.TrimSpace(ii[0:j]), 64)
			if err != nil {
				return
			}

			v = complex(ff, fi)
			return
		}
		err = errors.New("for a complex number, the imaginary part should end with 'i', such as '3+4i'")
		return

		// err = errors.New("not valid complex number.")
	}

	var ff float64
	ff, err = strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return
	}
	v = complex(ff, 0)
	return
}

// FindSubCommand find sub-command with `longName` from `cmd`
// if cmd == nil: finding from root command
func FindSubCommand(longName string, cmd *Command) (res *Command) {
	if cmd == nil {
		cmd = &internalGetWorker().rootCommand.Command
	}
	res = cmd.FindSubCommand(longName)
	return
}

// FindFlag find flag with `longName` from `cmd`
// if cmd == nil: finding from root command
func FindFlag(longName string, cmd *Command) (res *Flag) {
	if cmd == nil {
		cmd = &internalGetWorker().rootCommand.Command
	}
	res = cmd.FindFlag(longName)
	return
}

// FindSubCommandRecursive find sub-command with `longName` from `cmd` recursively
// if cmd == nil: finding from root command
func FindSubCommandRecursive(longName string, cmd *Command) (res *Command) {
	if cmd == nil {
		cmd = &internalGetWorker().rootCommand.Command
	}
	res = cmd.FindSubCommandRecursive(longName)
	return
}

// FindFlagRecursive find flag with `longName` from `cmd` recursively
// if cmd == nil: finding from root command
func FindFlagRecursive(longName string, cmd *Command) (res *Flag) {
	if cmd == nil {
		cmd = &internalGetWorker().rootCommand.Command
	}
	res = cmd.FindFlagRecursive(longName)
	return
}

func manBr(s string) string {
	var lines []string
	for _, l := range strings.Split(s, "\n") {
		lines = append(lines, l+"\n.br")
	}
	return strings.Join(lines, "\n")
}

func manWs(fmtStr string, args ...interface{}) string {
	str := fmt.Sprintf(fmtStr, args...)
	str = replaceAll(strings.TrimSpace(str), " ", `\ `)
	return str
}

func manExamples(s string, data interface{}) string {
	var (
		sources  = strings.Split(s, "\n")
		lines    []string
		lastLine string
	)
	for _, l := range sources {
		if strings.HasPrefix(l, "$ {{.AppName}}") {
			lines = append(lines, `.TP \w'{{.AppName}}\ 'u
.BI {{.AppName}} \ `+manWs(l[14:]))
		} else {
			if len(lastLine) == 0 {
				lastLine = strings.TrimSpace(l)
				// ignore multiple empty lines, compat them as one line.
				if len(lastLine) != 0 {
					lines = append(lines, lastLine+"\n.br")
				}
			} else {
				lastLine = strings.TrimSpace(l)
				lines = append(lines, lastLine+"\n.br")
			}
		}
	}
	return tplApply(strings.Join(lines, "\n"), data)
}

func tplApply(tmpl string, data interface{}) string {
	var w = new(bytes.Buffer)
	var tpl = template.Must(template.New("x").Parse(tmpl))
	if err := tpl.Execute(w, data); err != nil {
		log2.Printf("tpl execute error: %v", err)
		return ""
	}
	return w.String()
}

//
// external
//

// Launch executes a command setting both standard input, output and error.
func Launch(cmd string, args ...string) (err error) {
	c := exec.Command(cmd, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	err = c.Run()

	if err != nil {
		if _, isExitError := err.(*exec.ExitError); isExitError {
			err = nil
		}
	}
	return
}

// // LaunchSudo executes a command under "sudo".
// func LaunchSudo(cmd string, args ...string) error {
// 	return Launch("sudo", append([]string{cmd}, args...)...)
// }

//
// editor
//

// func getEditor() (string, error) {
// 	if GetEditor != nil {
// 		return GetEditor()
// 	}
// 	return exec.LookPath(DefaultEditor)
// }

// InDebugging return the status if cmdr was built with debug mode / or the app running under a debugger attached.
//
// To enable the debugger attached mode for cmdr, run `go build` with `-tags=delve` options. eg:
//
//     go run -tags=delve ./cli
//     go build -tags=delve -o my-app ./cli
//
// For Goland, you can enable this under 'Run/Debug Configurations', by adding the following into 'Go tool arguments:'
//
//     -tags=delve
//
// InDebugging() is a synonym to IsDebuggerAttached().
//
// NOTE that `isdelve` algor is from https://stackoverflow.com/questions/47879070/how-can-i-see-if-the-goland-debugger-is-running-in-the-program
//
//noinspection GoBoolExpressions
func InDebugging() bool {
	return log.InDebugging() // isdelve.Enabled
}

// IsDebuggerAttached return the status if cmdr was built with debug mode / or the app running under a debugger attached.
//
// To enable the debugger attached mode for cmdr, run `go build` with `-tags=delve` options. eg:
//
//     go run -tags=delve ./cli
//     go build -tags=delve -o my-app ./cli
//
// For Goland, you can enable this under 'Run/Debug Configurations', by adding the following into 'Go tool arguments:'
//
//     -tags=delve
//
// IsDebuggerAttached() is a synonym to InDebugging().
//
// NOTE that `isdelve` algor is from https://stackoverflow.com/questions/47879070/how-can-i-see-if-the-goland-debugger-is-running-in-the-program
//
//noinspection GoBoolExpressions
func IsDebuggerAttached() bool {
	return log.InDebugging() // isdelve.Enabled
}

// InTesting detects whether is running under go test mode
func InTesting() bool {
	if !strings.HasSuffix(SavedOsArgs[0], ".test") &&
		!strings.Contains(SavedOsArgs[0], "/T/___Test") {

		// [0] = /var/folders/td/2475l44j4n3dcjhqbmf3p5l40000gq/T/go-build328292371/b001/exe/main
		// !strings.Contains(SavedOsArgs[0], "/T/go-build")

		for _, s := range SavedOsArgs {
			if s == "-test.v" || s == "-test.run" {
				return true
			}
		}
		return false

	}
	return true
}

func randomFilename() (fn string) {
	buf := make([]byte, 16)
	fn = os.Getenv("HOME") + ".CMDR_EDIT_FILE"
	if _, err := rand.Read(buf); err == nil {
		fn = fmt.Sprintf("%v/.CMDR_%x", os.Getenv("HOME"), buf)
	}
	return
}

// const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomStringPure(length int) (result string) {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err == nil {
		result = string(buf)
	}
	return
	// source:=rand.NewSource(time.Now().UnixNano())
	// b := make([]byte, length)
	// for i := range b {
	// 	b[i] = charset[source.Int63()%int64(len(charset))]
	// }
	// return string(b)
}

// LaunchEditor launches the specified editor
func LaunchEditor(editor string) (content []byte, err error) {
	return launchEditorWith(editor, randomFilename())
}

// LaunchEditorWith launches the specified editor with a filename
func LaunchEditorWith(editor string, filename string) (content []byte, err error) {
	return launchEditorWith(editor, filename)
}

func launchEditorWith(editor, filename string) (content []byte, err error) {
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		if _, isExitError := err.(*exec.ExitError); !isExitError {
			return
		}
	}

	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, nil
	}
	return
}

// Soundex returns the english word's soundex value, such as: 'tags' => 't322'
func Soundex(s string) (snd4 string) {
	return soundex(s)
}

func soundex(s string) (snd4 string) {
	// if len(s) == 0 {
	// 	return
	// }

	var src, tgt []rune
	src = []rune(s)

	i := 0
	for ; i < len(src); i++ {
		if !(src[i] == '-' || src[i] == '~' || src[i] == '+') {
			// first char
			tgt = append(tgt, src[i])
			break
		}
	}

	for ; i < len(src); i++ {
		ch := src[i]
		switch ch {
		case 'a', 'e', 'i', 'o', 'u', 'y', 'h', 'w': // do nothing to remove it
		case 'b', 'f', 'p', 'v':
			tgt = append(tgt, '1')
		case 'c', 'g', 'j', 'k', 'q', 's', 'x', 'z':
			tgt = append(tgt, '2')
		case 'd', 't':
			tgt = append(tgt, '3')
		case 'l':
			tgt = append(tgt, '4')
		case 'm', 'n':
			tgt = append(tgt, '5')
		case 'r':
			tgt = append(tgt, '6')
		}
	}

	snd4 = string(tgt)
	return
}

const stringMetricFactor = 100000000000

type (
	// StringDistance is an interface for string metric.
	// A string metric is a metric that measures distance between two strings.
	// In most case, it means that the edit distance about those two strings.
	// This is saying, it is how many times are needed while you were
	// modifying string to another one, note that inserting, deleting,
	// substing one character means once.
	StringDistance interface {
		Calc(s1, s2 string, opts ...DistanceOption) (distance int)
	}

	// DistanceOption is a functional options prototype
	DistanceOption func(StringDistance)
)

// JaroWinklerDistance returns an calculator for two strings distance metric, with Jaro-Winkler algorithm.
func JaroWinklerDistance(opts ...DistanceOption) StringDistance {
	x := &jaroWinklerDistance{threshold: 0.7, factor: stringMetricFactor}
	for _, c := range opts {
		c(x)
	}
	return x
}

// JWWithThreshold sets the threshold for Jaro-Winkler algorithm.
func JWWithThreshold(threshold float64) DistanceOption {
	return func(distance StringDistance) {
		if v, ok := distance.(*jaroWinklerDistance); ok {
			v.threshold = threshold
		}
	}
}

type jaroWinklerDistance struct {
	threshold float64
	factor    float64

	matches        int
	maxLength      int
	transpositions int // transpositions is a double number here
	prefix         int

	distance float64
}

func (s *jaroWinklerDistance) Calc(src1, src2 string, opts ...DistanceOption) (distance int) {
	s1, s2 := []rune(src1), []rune(src2)
	lenMax, lenMin := len(s1), len(s2)

	var sMax, sMin []rune
	if lenMax > lenMin {
		sMax, sMin = s1, s2
	} else {
		sMax, sMin = s2, s1
		lenMax, lenMin = lenMin, lenMax
	}
	s.maxLength = lenMax

	iMatchIndexes, matchFlags := s.match(sMax, sMin, lenMax, lenMin)
	s.findTranspositions(sMax, sMin, lenMax, lenMin, iMatchIndexes, matchFlags)

	// println("  matches, transpositions, prefix: ", s.matches, s.transpositions, s.prefix)

	if s.matches == 0 {
		s.distance = 0
		return 0
	}

	m := float64(s.matches)
	jaroDistance := m/float64(lenMax) + m/float64(lenMin)
	jaroDistance += (m - float64(s.transpositions)/2) / m
	jaroDistance /= 3

	var jw float64
	if jaroDistance < s.threshold {
		jw = jaroDistance
	} else {
		jw = jaroDistance + math.Min(0.1, 1/float64(s.maxLength))*float64(s.prefix)*(1-jaroDistance)
	}

	// println("  jaro, jw: ", jaroDistance, jw)

	s.distance = jw * s.factor
	distance = int(math.Round(s.distance))
	return
}

func (s *jaroWinklerDistance) match(sMax, sMin []rune, lenMax, lenMin int) (iMatchIndexes []int, matchFlags []bool) {
	iRange := max(lenMax/2-1, 0)
	iMatchIndexes = make([]int, lenMin)
	for i := 0; i < lenMin; i++ {
		iMatchIndexes[i] = -1
	}

	s.prefix, s.matches = 0, 0
	for mi := 0; mi < len(sMin); mi++ {
		if sMax[mi] == sMin[mi] {
			s.prefix++
		} else {
			break
		}
	}
	s.matches = s.prefix

	matchFlags = make([]bool, lenMax)

	for mi := s.prefix; mi < lenMin; mi++ {
		c1 := sMin[mi]
		xi, xn := max(mi-iRange, s.prefix), lenMax // min(mi+iRange-1, lenMax)
		for ; xi < xn; xi++ {
			if !matchFlags[xi] && c1 == sMax[xi] {
				iMatchIndexes[mi] = xi
				matchFlags[xi] = true
				s.matches++
				break
			}
		}
	}
	return
}

func (s *jaroWinklerDistance) findTranspositions(sMax, sMin []rune, lenMax, lenMin int, iMatchIndexes []int, matchFlags []bool) {
	ms1, ms2 := make([]rune, s.matches), make([]rune, s.matches)
	for i, si := 0, 0; i < lenMin; i++ {
		if iMatchIndexes[i] != -1 {
			ms1[si] = sMin[i]
			si++
		}
	}
	for i, si := 0, 0; i < lenMax; i++ {
		if matchFlags[i] {
			ms2[si] = sMax[i]
			si++
		}
	}
	// fmt.Printf("iMatchIndexes, s1, s2: %v, %v, %v\n", iMatchIndexes, string(sMax), string(sMin))
	// println("     ms1, ms2: ", string(ms1), string(ms2))

	s.transpositions = 0
	for mi := 0; mi < len(ms1); mi++ {
		if ms1[mi] != ms2[mi] {
			s.transpositions++
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// StripPrefix strips the prefix 'p' from a string 's'
func StripPrefix(s, p string) string {
	return stripPrefix(s, p)
}

func stripPrefix(s, p string) string {
	if strings.HasPrefix(s, p) {
		return s[len(p):]
	}
	return s
}

// StripOrderPrefix strips the prefix string fragment for sorting order.
// see also: Command.Group, Flag.Group, ...
// An order prefix is a dotted string with multiple alphabet and digit. Such as:
// "zzzz.", "0001.", "700.", "A1." ...
func StripOrderPrefix(s string) string {
	if xre.MatchString(s) {
		s = s[strings.Index(s, ".")+1:]
	}
	return s
}

// HasOrderPrefix tests whether an order prefix is present or not.
// An order prefix is a dotted string with multiple alphabet and digit. Such as:
// "zzzz.", "0001.", "700.", "A1." ...
func HasOrderPrefix(s string) bool {
	return xre.MatchString(s)
}

var (
	xre = regexp.MustCompile(`^[0-9A-Za-z]+\.(.+)$`)
)

// IsDigitHeavy tests if the whole string is digit
func IsDigitHeavy(s string) bool {
	m, _ := regexp.MatchString("^\\d+$", s)
	// if err != nil {
	// 	return false
	// }
	return m
}

func (w *ExecWorker) setupRootCommand(rootCmd *RootCommand) {
	w.rootCommand = rootCmd

	w.rootCommand.ow = w.defaultStdout
	w.rootCommand.oerr = w.defaultStderr

	if len(conf.AppName) == 0 {
		conf.AppName = w.rootCommand.AppName
		conf.Version = w.rootCommand.Version
		_ = os.Setenv("APPNAME", conf.AppName)
	}
	if len(conf.Buildstamp) == 0 {
		conf.Buildstamp = time.Now().Format(time.RFC1123)
	}
}

func (w *ExecWorker) getPrefix() string {
	return strings.Join(w.rxxtPrefixes, ".")
}

func (w *ExecWorker) getRemainArgs(pkg *ptpkg, args []string) []string {
	return pkg.remainArgs
}

// PressEnterToContinue lets program pause and wait for user's ENTER key press in console/terminal
func PressEnterToContinue(in io.Reader, msg ...string) (input string) {
	if len(msg) > 0 && len(msg[0]) > 0 {
		fmt.Print(msg[0])
	} else {
		fmt.Print("Press 'Enter' to continue...")
	}
	b, _ := bufio.NewReader(in).ReadBytes('\n')
	return strings.TrimRight(string(b), "\n")
}

// PressAnyKeyToContinue lets program pause and wait for user's ANY key press in console/terminal
func PressAnyKeyToContinue(in io.Reader, msg ...string) (input string) {
	if len(msg) > 0 && len(msg[0]) > 0 {
		fmt.Print(msg[0])
	} else {
		fmt.Print("Press any key to continue...")
	}
	_, _ = fmt.Fscanf(in, "%s", &input)
	return
}

// SavedOsArgs is a copy of os.Args, just for testing
var SavedOsArgs []string

func init() {
	if SavedOsArgs == nil {
		// bug: can't copt slice to slice: _ = StandardCopier.Copy(&SavedOsArgs, &os.Args)
		for _, s := range os.Args {
			SavedOsArgs = append(SavedOsArgs, s)
		}
	}
}
