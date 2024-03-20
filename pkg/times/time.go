package times

import (
	"regexp"
	"strconv"
	"sync"
	"time"
)

// MustSmartParseTime parses a formatted string and returns the time value it represents.
func MustSmartParseTime(str string) (tm time.Time) {
	tm, _ = smartParseTime(str)
	return
}

// MustSmartParseTimePtr parses a formatted string and returns the time value it represents.
func MustSmartParseTimePtr(str string) (tm *time.Time) {
	var tm1 time.Time
	tm1, _ = smartParseTime(str)
	return &tm1
}

// SmartParseTime parses a formatted string and returns the time value it represents.
//
// The example for [time.Time.Format()] demonstrates the working of the layout string
// in detail and is a good reference.
//
// When parsing (only), the input may contain a fractional second
// field immediately after the seconds field, even if the layout does not
// signify its presence. In that case either a comma or a decimal point
// followed by a maximal series of digits is parsed as a fractional second.
// Fractional seconds are truncated to nanosecond precision.
//
// Elements omitted from the layout are assumed to be zero or, when
// zero is impossible, one, so parsing "3:04pm" returns the time
// corresponding to Jan 1, year 0, 15:04:00 UTC (note that because the year is
// 0, this time is before the zero Time).
// Years must be in the range 0000..9999. The day of the week is checked
// for syntax but it is otherwise ignored.
//
// For layouts specifying the two-digit year 06, a value NN >= 69 will be treated
// as 19NN and a value NN < 69 will be treated as 20NN.
//
// The remainder of this comment describes the handling of time zones.
//
// In the absence of a time zone indicator, Parse returns a time in UTC.
//
// When parsing a time with a zone offset like -0700, if the offset corresponds
// to a time zone used by the current location (Local), then Parse uses that
// location and zone in the returned time. Otherwise it records the time as
// being in a fabricated location with time fixed at the given zone offset.
//
// When parsing a time with a zone abbreviation like MST, if the zone abbreviation
// has a defined offset in the current location, then that offset is used.
// The zone abbreviation "UTC" is recognized as UTC regardless of location.
// If the zone abbreviation is unknown, Parse records the time as being
// in a fabricated location with the given zone abbreviation and a zero offset.
// This choice means that such a time can be parsed and reformatted with the
// same layout losslessly, but the exact instant used in the representation will
// differ by the actual zone offset. To avoid such problems, prefer time layouts
// that use a numeric zone offset, or use ParseInLocation.
func SmartParseTime(str string) (tm time.Time, err error) {
	return smartParseTime(str)
}

func smartParseTime(str string) (tm time.Time, err error) {
	for _, layout := range onceInitTimeFormats() {
		if tm, err = time.Parse(layout, str); err == nil {
			break
		}
	}
	return
}

var knownDateTimeFormats []string
var onceFormats sync.Once

func onceInitTimeFormats() []string {
	onceFormats.Do(func() {
		knownDateTimeFormats = []string{
			"2006-01-02 15:04:05.999999999 -0700",
			"2006-01-02 15:04:05.999999999Z07:00",
			"2006-01-02 15:04:05.999999999",
			"2006-01-02 15:04:05.999",
			"2006-01-02 15:04:05",
			"2006-01-02",
			"2006/01/02",
			"01/02/2006",
			"01-02",

			"2006-1-2 15:4:5.999999999 -0700",
			"2006-1-2 15:4:5.999999999Z07:00",
			"2006-1-2 15:4:5.999999999",
			"2006-1-2 15:4:5.999",
			"2006-1-2 15:4:5",
			"2006-1-2",
			"2006/1/2",
			"1/2/2006",
			"1-2",

			"15:04:05.999999999",
			"15:04.999999999",
			"15:04:05.999",
			"15:04.999",
			"15:04:05",
			"15:04",

			"15:4:5.999999999",
			"15:4.999999999",
			"15:4:5.999",
			"15:4.999",
			"15:4:5",
			"15:4",

			time.RFC3339,
			time.RFC3339Nano,
			time.RFC1123Z,
			time.RFC1123,
			time.RFC850,
			time.RFC822Z,
			time.RFC822,
			time.RubyDate,
			time.UnixDate,
			time.ANSIC,
		}
	})
	return knownDateTimeFormats
}

// AddKnownTimeFormats appends more time layouts to the trying list
// used by SmartParseTime.
func AddKnownTimeFormats(format ...string) {
	a := onceInitTimeFormats()
	a = append(a, format...)
}

// RoundTime strips small ticks from a time.Time value.
// For example:
//
//	assert.True(RoundTime("h", SmartParseTime("5:11:22")), SmartParseTime("5:0:0"))
func RoundTime(roundTo string, value time.Time) string {
	since := time.Since(value)
	if roundTo == "h" {
		since -= since % time.Hour
	}
	if roundTo == "m" {
		since -= since % time.Minute
	}
	if roundTo == "s" {
		since -= since % time.Second
	}
	return since.String()
}

func shortDur(d time.Duration) string {
	s := d.String()

	// if strings.HasSuffix(s, "m0s") {
	// 	s = s[:len(s)-2]
	// }
	// if strings.HasSuffix(s, "h0m") {
	// 	s = s[:len(s)-2]
	// }

	for _, z := range []struct{ reg, rep string }{
		{`(m?[hmnus])0[hmnus]?s?`, "$1"},
	} {
		re := regexp.MustCompile(z.reg)
		s = re.ReplaceAllString(s, z.rep)
	}
	return s
}

func DurationFromFloat(f float64) time.Duration {
	return time.Duration(f * float64(time.Second))
}

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func ParseDuration(s string) (dur time.Duration, err error) {
	dur, err = time.ParseDuration(s)
	return
}

// MustParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func MustParseDuration(s string) (dur time.Duration) {
	dur, _ = time.ParseDuration(s)
	return
}

func SmartParseInt(s string) (ret int64, err error) {
	ret, err = strconv.ParseInt(s, 0, 64)
	return
}

func MustSmartParseInt(s string) (ret int64) {
	ret, _ = strconv.ParseInt(s, 0, 64)
	return
}

func SmartParseUint(s string) (ret uint64, err error) {
	ret, err = strconv.ParseUint(s, 0, 64)
	return
}

func MstSmartParseUint(s string) (ret uint64) {
	ret, _ = strconv.ParseUint(s, 0, 64)
	return
}

func ParseFloat(s string) (ret float64, err error) {
	ret, err = strconv.ParseFloat(s, 64)
	return
}

func MustParseFloat(s string) (ret float64) {
	ret, _ = strconv.ParseFloat(s, 64)
	return
}

func ParseComplex(s string) (ret complex128, err error) {
	ret, err = strconv.ParseComplex(s, 64)
	return
}

func MustParseComplex(s string) (ret complex128) {
	ret, _ = strconv.ParseComplex(s, 64)
	return
}
