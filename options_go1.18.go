// Copyright © 2022 Hedzr Yeh.

//go:build go1.18
// +build go1.18

package cmdr

// NewAny creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewAny(defaultValue ...any) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &stringOpt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv interface{}
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// // NewAny creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
// func NewAny[T any](defaultValue ...T) (opt OptFlag) {
// 	workingFlag := &Flag{}
// 	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
// 	opt = &stringOpt{optFlagImpl: optFlagImpl{working: workingFlag}}
// 	var dv any
// 	for _, v := range defaultValue {
// 		dv = v
// 	}
// 	opt.DefaultValue(dv, "")
// 	return
// }
