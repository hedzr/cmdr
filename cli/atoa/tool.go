package atoa

import (
	"reflect"
)

func eat(str string, chars ...rune) (ret string, ate bool) {
	r := []rune(str)
	if len(r) > 0 {
		for _, ch := range chars {
			if r[0] == ch {
				ate, r = true, r[1:]
			}
		}
	}
	ret = string(r)
	return
}

func eatTail(str string, chars ...rune) (ret string, ate bool) {
	r := []rune(str)
	for _, ch := range chars {
		if rl := len(r); rl > 0 && r[rl-1] == ch {
			ate, r = true, r[:rl-1]
		}
	}
	ret = string(r)
	return
}

//

//

//

// skipWSAndNextChar looks ahead the 'chars' and advances 'pos' after its.
// If 'chars' not found, 'fromPos' returned as 'pos'.
// If next char is not in 'chars' set, return right away.
//
// For example:
//
//	ch, pos = skipWSAndNextChar(runes, pos, '}')
//
// This line checks next char and advances 'pos' if it's '}', else
// nothing to do and 'pos' kept not changed.
func skipWSAndNextChar(runes []rune, fromPos int, chars ...rune) (ch rune, pos int) { //nolint:unused,revive
	pos, i, rl := fromPos, skipWS(runes, fromPos), len(runes)
	for ; i < rl; i++ {
		r, matched := runes[i], false
		for _, ch1 := range chars {
			if ch1 == r {
				matched, pos, ch = true, i+1, ch1
				break
			}
		}
		if !matched {
			break
		}
	}
	return
}

func skipWS(runes []rune, fromPos int) (pos int) {
	i, rl := fromPos, len(runes)
	for ; i < rl; i++ {
		r := runes[i]
		if !(r == ' ' || r == '\t' || r == '\r' || r == '\n') {
			break
		}
	}
	pos = i
	return
}

func preferLookAndSkip(runes []rune, chars ...rune) (ch rune, pos int) { //nolint:revive,unused
	i, rl := skipWS(runes, 0), len(runes)
	for ; i < rl; i++ {
		r, matched := runes[i], false
		for _, ch1 := range chars {
			if ch1 == r {
				matched, pos, ch = true, i+1, ch1
				break
			}
		}
		if !matched {
			return
		}
	}
	pos = i + 1
	return
}

// preferLookAhead looks ahead for the rune character of 'chars'.
// return 'ch' == 0 if nothing found, or 'pos' to point to
// the matched character.
func preferLookAhead(runes []rune, chars ...rune) (ch rune, pos int) {
	i, rl := skipWS(runes, 0), len(runes)
	for ; i < rl; i++ {
		r := runes[i]
		for _, ch1 := range chars {
			if ch1 == r {
				ch, pos = ch1, i
				return
			}
		}
	}
	// pos = i
	return
}

func preferLookAheadOrEOF(runes []rune, chars ...rune) (ch rune, pos int) {
	ch, pos = preferLookAhead(runes, chars...)
	if ch == 0 {
		pos = len(runes)
	}
	return
}

// lookAhead looks ahead for any rune character of 'chars'.
// return 'pos' to point the matched character or end of 'runes'.
func lookAhead(runes []rune, chars ...rune) (pos int) { //nolint:unused
	var ch rune
	ch, pos = preferLookAhead(runes, chars...)
	if ch == 0 {
		pos = len(runes)
	}
	return
}

//

//

//

func intByKind(v int64, k reflect.Kind) (ret any) {
	switch k {
	case reflect.Int:
		return int(v)
	case reflect.Int8:
		return int8(v)
	case reflect.Int16:
		return int16(v)
	case reflect.Int32:
		return int32(v)
	case reflect.Int64:
		return v
	}
	return
}

func uintByKind(v uint64, k reflect.Kind) (ret any) {
	switch k {
	case reflect.Uint:
		return uint(v)
	case reflect.Uint8:
		return uint8(v)
	case reflect.Uint16:
		return uint16(v)
	case reflect.Uint32:
		return uint32(v)
	case reflect.Uint64:
		return v
	}
	return
}

func fltByKind(v float64, k reflect.Kind) (ret any) {
	switch k {
	case reflect.Float64:
		return v
	case reflect.Float32:
		return float32(v)
	}
	return
}

func comByKind(v complex128, k reflect.Kind) (ret any) {
	switch k {
	case reflect.Complex128:
		return v
	case reflect.Complex64:
		return complex64(v)
	}
	return
}
