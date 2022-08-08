// Copyright © 2022 Hedzr Yeh.

package cmdr

import (
	"encoding"
	"encoding/base64"
	"fmt"
	"github.com/hedzr/log"
)

// TextVar _
type TextVar interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

// NewTextVar creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewTextVar(defaultValue ...TextVar) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &stringOpt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv TextVar
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// GetTextVar returns the text-var value of an `Option` key.
func GetTextVar(key string, defaultVal ...TextVar) TextVar {
	return currentOptions().GetTextVar(key, defaultVal...)
}

// GetTextVarP returns the text-var value of an `Option` key.
func GetTextVarP(prefix, key string, defaultVal ...TextVar) TextVar {
	return currentOptions().GetTextVar(fmt.Sprintf("%s.%s", prefix, key), defaultVal...)
}

// GetTextVarR returns the text-var value of an `Option` key.
func GetTextVarR(key string, defaultVal ...TextVar) TextVar {
	w := internalGetWorker()
	return w.rxxtOptions.GetTextVar(w.wrapWithRxxtPrefix(key), defaultVal...)
}

// GetTextVarRP returns the text-var value of an `Option` key.
//
// A 'text-var' is a type which implements these interfaces:
//
//	encoding.TextMarshaler and encoding.TextUnmarshaler
//
// The types typically are: net.IPv4, or time.Time.
func GetTextVarRP(prefix, key string, defaultVal ...TextVar) TextVar {
	w := internalGetWorker()
	return w.rxxtOptions.GetTextVar(w.wrapWithRxxtPrefix(fmt.Sprintf("%s.%s", prefix, key)), defaultVal...)
}

// GetTextVar returns the time duration value of an `Option` key.
func (s *Options) GetTextVar(key string, defaultVal ...TextVar) (ir TextVar) {
	if v, ok := s.hasKey(key); !ok {
		goto FALLBACK
	} else {
		switch vr := v.(type) {
		case TextVar:
			ir = vr
		case []byte: // NOTE: the codes in this branch are bad, never used, but kept for ref in future
			if err := ir.UnmarshalText(vr); err != nil {
				goto FALLBACK
			}
		case string: // NOTE: the codes in this branch are bad, never used, but kept for ref in future
			if data, err := base64.StdEncoding.DecodeString(vr); err != nil {
				goto FALLBACK
			} else if err = ir.UnmarshalText(data); err != nil {
				goto FALLBACK
			}
		default:
			log.Errorf("unrecognized default value in Option Store found for a TextVar slot: %v (%T)", v, v)
		}
	}
	return

FALLBACK:
	for _, vv := range defaultVal {
		ir = vv
	}
	return
}
