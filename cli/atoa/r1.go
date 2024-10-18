package atoa

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/hedzr/evendeep/ref"

	"github.com/hedzr/cmdr/v2/pkg/exec"
	"github.com/hedzr/cmdr/v2/pkg/logz"
)

//
// some links which inspired me:
// - https://github.com/golang/go/issues/57975
// - ?

// stepComplexObject _ here
//
// preferKind can be array, slice, map or struct.
func (s *toS) stepComplexObject(
	preferKind reflect.Kind, typObj reflect.Type, runes []rune, fromPos int, meme any,
) (pos int, v any, err error) {
	var rl int
	pos, rl = fromPos, len(runes)

	// try to analysis lexical elements for parsing nested form, such as "[[7,8],[9,0]]".
	for pos < rl {
		ch := runes[pos]

		switch ch {
		case '[':
			pos, v, err = s.stepArrayOrSlice(preferKind, typObj, runes, pos+1, meme)
			return
		case '{':
			pos, v, err = s.stepObjectOrMap(preferKind, typObj, runes, pos+1, meme)
			return
		}

		switch preferKind {
		case reflect.Array, reflect.Slice:
			pos, v, err = s.stepArrayOrSlice(preferKind, typObj, runes, pos, meme)
			return
		case reflect.Map, reflect.Struct:
			pos, v, err = s.stepObjectOrMap(preferKind, typObj, runes, pos, meme)
			return
		}
	}

	// cannot work for nested form, such as "[[7,8],[9,0]]".
	pos, v, err = s.extractByDefaultWay(preferKind, typObj, runes, pos, meme)
	return
}

func (s *toS) extractByDefaultWay( //nolint:revive
	preferKind reflect.Kind, typObj reflect.Type, runes []rune, fromPos int, meme any,
) (position int, v any, err error) { //nolint:unparam
	ssa := strings.Split(string(runes[fromPos:]), ",")

	switch preferKind {
	case reflect.Array:
		elt, l, rv := typObj.Elem(), typObj.Len(), reflect.New(typObj)

		for i := 0; i < min(len(ssa), l); i++ {
			txt := ssa[i]
			tv1, err1 := s.parseImpl(txt, elt, meme)
			if err1 != nil {
				err = err1
				return
			}
			rv.Elem().Index(i).Set(reflect.ValueOf(tv1))
		}

		v = rv.Elem().Interface()

	case reflect.Slice:
		elt, rv := typObj.Elem(), reflect.MakeSlice(typObj, 0, len(ssa))

		for _, txt := range ssa {
			var vv any
			vv, err = s.parseImpl(strings.TrimSpace(txt), elt, meme)
			if err == nil {
				rv = reflect.Append(rv, reflect.ValueOf(vv))
			}
		}
		v = rv.Interface()

	case reflect.Map:
		rv, kt, vt := reflect.MakeMapWithSize(typObj, len(ssa)), typObj.Key(), typObj.Elem()

		re := regexp.MustCompile(`(.*)[=:](.*)`)
		for _, txt := range ssa {
			a := re.FindAllStringSubmatch(txt, -1)
			if len(a) > 0 {
				b := a[0]
				if len(b) > 1 {
					k1, v1 := strings.TrimSpace(b[1]), strings.TrimSpace(b[2])
					kv, err1 := s.parseImpl(k1, kt, meme)
					if err1 != nil {
						err = err1
						return
					}
					vv, err2 := s.parseImpl(v1, vt, meme)
					if err2 != nil {
						err = err2
						return
					}
					rv.SetMapIndex(reflect.ValueOf(kv), reflect.ValueOf(vv))
				}
			}
		}
		v = rv.Interface()
	}
	return
}

func (s *toS) stepArrayOrSlice( //nolint:revive
	preferKind reflect.Kind, typArray reflect.Type, runes []rune, fromPos int, meme any,
) (pos int, v any, err error) { //nolint:unparam
	ctx := context.Background()
	r, rl, elt, rvLen := runes[fromPos:], len(runes), typArray, 0
	pos = fromPos
	logz.VerboseContext(ctx, "[stepArrayOrSlice]", "pos", pos, "r", string(r), "el", ref.Typfmt(elt))

	var rv reflect.Value
	switch preferKind {
	case reflect.Array:
		rv, rvLen, elt = reflect.New(typArray), typArray.Len(), typArray.Elem()
	case reflect.Slice:
		rv, elt = reflect.MakeSlice(typArray, 0, 0), typArray.Elem()
	default:
		err = errors.New("only array and slice type are allowed")
		return
	}

	ix := 0
	arraySet := func(rv *reflect.Value, ix *int, el any, err error) {
		if err == nil {
			if preferKind == reflect.Slice {
				*rv = reflect.Append(*rv, reflect.ValueOf(el))
			} else if *ix < rvLen {
				rv.Elem().Index(*ix).Set(reflect.ValueOf(el))
			}
			*ix++
		}
	}

forOneElem:
	for pos < rl && err == nil {
		r = runes[pos:]
		logz.VerboseContext(ctx, "[stepArrayOrSlice] for each of 'r' from 'pos'", "pos", pos, "r", string(r))
		ch, p := preferLookAheadOrEOF(r, '[', ']', '{', '(', ',')
		pos += p

		var el any
		switch ch {
		case '[': // =91
			logz.VerboseContext(ctx, "[stepArrayOrSlice] entering stepArrayOfSlice", "r", string(runes[pos+1:]), "el", ref.Typfmt(elt))
			pos, el, err = s.stepArrayOrSlice(reflect.Slice, elt, runes, pos+1, meme)
			arraySet(&rv, &ix, el, err)
			if ch, pos = skipWSAndNextChar(runes, pos, ']'); ch == ']' {
				break forOneElem
			}
			_, pos = skipWSAndNextChar(runes, pos, ',')
		case '{': // =123
			logz.VerboseContext(ctx, "[stepArrayOrSlice] entering stepObjectOrMap", "r", string(runes[pos+1:]), "el", ref.Typfmt(elt))
			pos, el, err = s.stepObjectOrMap(reflect.Map, elt, runes, pos+1, meme)
			arraySet(&rv, &ix, el, err)
			if ch, pos = skipWSAndNextChar(runes, pos, ']'); ch == ']' {
				break forOneElem
			}
			_, pos = skipWSAndNextChar(runes, pos, ',')
		case ',', ']', 0:
			txt := r[:p]
			// skip ',' or ']'
			if ch != 0 {
				pos++
			}
			logz.VerboseContext(ctx, "[stepArrayOrSlice] append one elem", "el-txt", string(txt), "r", string(runes[pos:]))
			el, err = s.parseImpl(strings.TrimSpace(string(txt)), elt, meme)
			arraySet(&rv, &ix, el, err)
			if ch == ']' {
				_, pos = skipWSAndNextChar(runes, pos, ',')
				logz.VerboseContext(ctx, "[stepArrayOrSlice] end of one elem", "r", string(runes[pos:]))
				break forOneElem
				// } else {
				// pos++ // skipWSAndNextChar ','
			}
		}
	}

	_, pos = skipWSAndNextChar(runes, pos, ']')

	if k := rv.Kind(); preferKind == reflect.Array && k == reflect.Ptr {
		rv = rv.Elem()
	}
	v = rv.Interface()

	// var k []rune
	// if p := lookAhead(r, ':', '='); p >= 0 {
	// 	k, p = r[:p], p+1
	// 	p += skipWSAndNextChar(r[p:], ' ', '\t')
	// 	switch ch := r[p]; ch {
	// 	case '[':
	// 		p, v, err = s.stepArrayOrSlice(reflect.Slice, typArray, runes, fromPos+p, meme)
	// 	case '{':
	// 	}
	// }
	return
}

func (s *toS) stepObjectOrMap( //nolint:revive
	preferKind reflect.Kind, typObjOrMap reflect.Type, runes []rune, fromPos int, meme any, //nolint:revive
) (pos int, v any, err error) { //nolint:unparam
	ctx := context.Background()
	r, rl, elt := runes[fromPos:], len(runes), typObjOrMap.Elem()
	pos = fromPos
	logz.VerboseContext(ctx, "[stepObjectOrMap]", "pos", pos, "r", string(r), "el", ref.Typfmt(elt))

	var rv reflect.Value
	switch preferKind {
	case reflect.Map:
		rv = reflect.MakeMap(typObjOrMap)
	case reflect.Struct:
		err = errors.New("not implements")
		return
	default:
		err = errors.New("only map and struct type are allowed")
		return
	}

	ix := 0
	mapSet := func(rv *reflect.Value, ix *int, key, el any, err error) {
		if err == nil {
			if preferKind == reflect.Map {
				rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(el))
			} else { //nolint:staticcheck,revive
				// rv.Elem().Index(*ix).Set(reflect.ValueOf(el))
			}
			*ix++ // skip the ending '}'
		}
	}

forOneElem:
	for pos < rl && err == nil {
		r = runes[pos:]
		logz.VerboseContext(ctx, "[stepObjectOrMap] for each of 'r' from 'pos'", "pos", pos, "r", string(r))
		_, keyPos := preferLookAhead(r, '=', ':')
		ch, p := preferLookAheadOrEOF(r, '[', '{', '}', '(', ',')
		if keyPos >= 0 && p >= 0 && keyPos < p {
			key := exec.StripQuotes(strings.TrimSpace(string(r[:keyPos])))
			keyPos++
			pos += keyPos

			var el any
			switch ch {
			case '[':
				pos += p - keyPos + 1
				logz.VerboseContext(ctx, "[stepObjectOrMap] entering stepArrayOfSlice", "key", key, "r", string(runes[pos:]), "el", ref.Typfmt(elt))
				pos, el, err = s.stepArrayOrSlice(reflect.Slice, elt, runes, pos, meme)
				mapSet(&rv, &ix, key, el, err)
				if ch, pos = skipWSAndNextChar(runes, pos, '}'); ch == '}' {
					break forOneElem
				}
				_, pos = skipWSAndNextChar(runes, pos, ',')
			case '{':
				pos += p - keyPos + 1
				logz.VerboseContext(ctx, "[stepObjectOrMap] entering stepObjectOrMap", "key", key, "r", string(runes[pos:]), "el", ref.Typfmt(elt))
				pos, el, err = s.stepObjectOrMap(reflect.Map, elt, runes, pos, meme)
				mapSet(&rv, &ix, key, el, err)
				if ch, pos = skipWSAndNextChar(runes, pos, '}'); ch == '}' {
					break forOneElem
				}
				_, pos = skipWSAndNextChar(runes, pos, ',')
			case ',', '}', 0:
				txt := r[keyPos:p]
				pos += p - keyPos
				if ch != 0 {
					pos++
				}
				logz.VerboseContext(ctx, "[stepObjectOrMap] append one elem", "key", key, "el-txt", string(txt), "r", string(runes[pos:]))
				el, err = s.parseImpl(strings.TrimSpace(string(txt)), elt, meme)
				mapSet(&rv, &ix, key, el, err)
				if ch == '}' {
					_, pos = skipWSAndNextChar(runes, pos, ',')
					logz.VerboseContext(ctx, "[stepObjectOrMap] end of one elem", "r", string(runes[pos:]))
					break forOneElem
					// } else {
					// 	pos++ // skipWSAndNextChar ','
				}
			}
		}
	}

	_, pos = skipWSAndNextChar(runes, pos, '}')

	v = rv.Interface()
	return
}
