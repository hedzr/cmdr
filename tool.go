/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
)

// FindSubCommand find sub-command with `longName` from `cmd`
// if cmd == nil: finding from root command 
func FindSubCommand(longName string, cmd *Command) (res *Command) {
	if cmd == nil {
		cmd = &uniqueWorker.rootCommand.Command
	}
	for _, cx := range cmd.SubCommands {
		if longName == cx.Full {
			res = cx
			return
		}
	}
	return
}

// FindFlag find flag with `longName` from `cmd`
// if cmd == nil: finding from root command 
func FindFlag(longName string, cmd *Command) (res *Flag) {
	if cmd == nil {
		cmd = &uniqueWorker.rootCommand.Command
	}
	for _, cx := range cmd.Flags {
		if longName == cx.Full {
			res = cx
			return
		}
	}
	return
}

// FindSubCommandRecursive find sub-command with `longName` from `cmd` recursively
// if cmd == nil: finding from root command 
func FindSubCommandRecursive(longName string, cmd *Command) (res *Command) {
	if cmd == nil {
		cmd = &uniqueWorker.rootCommand.Command
	}
	for _, cx := range cmd.SubCommands {
		if longName == cx.Full {
			res = cx
			return
		}
	}
	for _, cx := range cmd.SubCommands {
		if len(cx.SubCommands) > 0 {
			if res = FindSubCommandRecursive(longName, cx); res != nil {
				return
			}
		}
	}
	return
}

// FindFlagRecursive find flag with `longName` from `cmd` recursively
// if cmd == nil: finding from root command 
func FindFlagRecursive(longName string, cmd *Command) (res *Flag) {
	if cmd == nil {
		cmd = &uniqueWorker.rootCommand.Command
	}
	for _, cx := range cmd.Flags {
		if longName == cx.Full {
			res = cx
			return
		}
	}
	for _, cx := range cmd.SubCommands {
		// if len(cx.SubCommands) > 0 {
		if res = FindFlagRecursive(longName, cx); res != nil {
			return
		}
		// }
	}
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
	str = strings.ReplaceAll(strings.TrimSpace(str), " ", `\ `)
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
		log.Fatalf("tpl execute error: %v", err)
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

// SavedOsArgs is a copy of os.Args, just for testing
var SavedOsArgs []string

func init() {
	// bug: can't copt slice to slice: _ = StandardCopier.Copy(&SavedOsArgs, &os.Args)
	for _, s := range os.Args {
		SavedOsArgs = append(SavedOsArgs, s)
	}
}

// InTesting detects whether is running under go test mode
func InTesting() bool {
	if strings.HasSuffix(SavedOsArgs[0], ".test") ||
		strings.Contains(SavedOsArgs[0], "/T/___Test") ||
		strings.Contains(SavedOsArgs[0], "/T/go-build") {
		return true
	}
	for _, s := range SavedOsArgs {
		if s == "-test.v" || s == "-test.run" {
			return true
		}
	}
	return false
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

func stripPrefix(s, p string) string {
	if strings.HasPrefix(s, p) {
		return s[len(p):]
	}
	return s
}

// IsDigitHeavy tests if the whole string is digit
func IsDigitHeavy(s string) bool {
	m, err := regexp.MatchString("^\\d", s)
	if err != nil {
		return false
	}
	return m
}
