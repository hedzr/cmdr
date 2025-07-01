package builder

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"unicode"

	"github.com/hedzr/cmdr/v2/cli"
	"github.com/hedzr/evendeep/ref"
	"github.com/hedzr/is"
	"github.com/hedzr/is/stringtool"
	logz "github.com/hedzr/logg/slog"
)

func WithStructBuilderTitleFormatter(titleFormatType TitleFormatType, customFormatter ...TitleFormatFunc) cli.StructBuilderOpt {
	return func(s any) {
		if ss, ok := s.(*sbS); ok {
			var yours = stringtool.ToKebabCase
			for _, c := range customFormatter {
				yours = c
			}
			var types = map[TitleFormatType]TitleFormatFunc{
				KebabCase:             stringtool.ToKebabCase,
				SnakeCase:             stringtool.ToSnakeCase,
				CamelCase:             stringtool.ToCamelCase,
				SmallCamelCase:        stringtool.ToSmallCamelCase,
				CustomTitleFormatType: yours,
			}
			if v, ok := types[titleFormatType]; ok {
				ss.titleFormatter = v
			}
		}
	}
}

type TitleFormatFunc func(title string) (formatted string)
type TitleFormatType int

const (
	KebabCase             TitleFormatType = iota // 'kebab-case-title'
	SnakeCase                                    // 'snake_case_title'
	CamelCase                                    // 'NormalCamelCaseTitle'
	SmallCamelCase                               // 'smallCamelCaseTitle'
	CustomTitleFormatType                        // yours formatter
)

func newStructBuilder(b buildable, structValue any, opts ...cli.StructBuilderOpt) cli.OptBuilder {
	return newStructBuilderFrom(nil, nil, b, structValue, opts...)
}

func newStructBuilderShort(b buildable, structValue any, opts ...cli.StructBuilderOpt) *sbS {
	return newStructBuilderFrom(nil, nil, b, structValue, opts...)
}

func newStructBuilderFrom(pc *constructCtx, from *cli.CmdS, b buildable, structValue any, opts ...cli.StructBuilderOpt) *sbS {
	s := &sbS{
		0, 0,
		// isAssumedAsRootCmd(assumedAsRootcmd),
		pc,
		b,
		from, // cmd here // new(cli.CmdS),
		from, // parent here
		structValue,
		stringtool.ToKebabCase,
	}

	for _, opt := range opts {
		opt(s)
	}

	// s.Long, s.Short, s.Aliases = theTitles(longTitle, titles...)

	if app, ok := b.(*appS); ok {
		// links app, owner, parent and this building cmd
		app.root.SetApp(app)
		s.cmd = app.root.Cmd.(*cli.CmdS)
		s.cmd.SetRoot(app.root)
		// s.CmdS.SetRoot(app.root)
		// s.CmdS.SetOwner(s.parent)
	} else if sb, ok := b.(*sbS); ok {
		s.cmd = sb.Building()
		// 	s.CmdS.SetOwner(s.parent)
		// 	s.CmdS.SetRoot(s.parent.Root())
	} else if sb, ok := b.(*ccb); ok {
		cc := from
		if from == nil { // replace s.cmd with the parent builder's building CmdS
			cc = sb.CmdS
			s.cmd = cc
			// s.cmd.SetOwner(cc)
			// s.cmd.SetRoot(cc.Root())
		} else {
			cc = from // sb.CmdS
			if cc == sb.CmdS {
				s.cmd = cc
				// s.cmd.SetOwner(cc)
			} else {
				s.cmd = cc
			}
			// s.cmd = new(cli.CmdS)
			// cc = sb.CmdS
			// s.cmd.SetOwner(cc)
			// s.cmd.SetRoot(cc.Root())
		}
		// s.asRoot = false
	}
	return s
}

type sbS struct {
	inCmd int32
	inFlg int32
	// asRoot         bool
	pc             *constructCtx // parent constructCtx if exists
	buildable      buildable
	cmd            *cli.CmdS
	from           *cli.CmdS
	structValue    any
	titleFormatter TitleFormatFunc
}

func isAssumedAsRootCmd(title string) bool {
	return strings.HasPrefix(title, "(") && strings.HasSuffix(title, ")")
}

func (s *sbS) Buildable() cli.OptBuilder { return s.buildable }
func (s *sbS) This() *cli.CmdS           { return s.cmd }
func (s *sbS) Building() *cli.CmdS       { return s.cmd }
func (s *sbS) Parent() *cli.CmdS {
	if cc, ok := s.cmd.OwnerCmd().(*cli.CmdS); ok {
		return cc
	}
	return nil
}

func (s *sbS) Build() {
	if err := s.construct(); err != nil {
		logz.Error("cannot construct cmdr command system from a struct value", "err", err)
		return
	}

	if a, ok := s.buildable.(adder); ok {
		if s.from == nil {
			logz.Verbose(assumedAsRootcmd)
			a.addCommand(nil)
		} else {
			logz.Verbose("normal")
			if ss, ok := s.buildable.(*sbS); ok { // struct builder
				if ss.cmd == s.cmd {
					a.addCommand(nil)
				} else {
					a.addCommand(s.cmd)
				}
			} else { // command builder
				a.addCommand(s.cmd)
			}
		}
	}
	atomic.StoreInt32(&s.inCmd, 0)
	atomic.StoreInt32(&s.inFlg, 0)
}

func (s *sbS) StructValue(structValue any) cli.StructBuilder {
	s.structValue = structValue
	return s
}

const assumedAsRootcmd = "(assumed-as-rootcmd)"

// addCommand adds a in-building Cmd into current CmdS as a child-/sub-command.
// used by adder when ccb.Build.
func (s *sbS) addCommand(child *cli.CmdS) {
	atomic.AddInt32(&s.inCmd, -1) // reset increased inCmd at AddCmd or Cmd
	s.cmd.AddSubCommand(child)
	if child != nil {
		logz.Trace(fmt.Sprintf("                      added cmd %v -> %v", child.String(), s.cmd))
	}
}

// addFlag adds a in-building Flg into current CmdS as its flag.
// used by adder when ccb.Build.
func (s *sbS) addFlag(child *cli.Flag) {
	atomic.AddInt32(&s.inFlg, -1)
	s.cmd.AddFlag(child)
	logz.Trace(fmt.Sprintf("                      added flg %v -> %v", child, s.cmd))
	// logz.Trace(fmt.Sprintf("[constructFrom]     | added %v -> %v", child, s.cmd))
}

func (s *sbS) construct() (err error) {
	rt := reflect.TypeOf(s.structValue)
	rv := reflect.ValueOf(s.structValue)
	logz.Debug(fmt.Sprintf("[constrcut()] \n    structValue = %p / %+v", s.structValue, s.structValue),
		"rv.type", ref.Typfmt(rt))
	if rt.Kind() != reflect.Struct { // is a Ptr?
		rt = ref.Rdecodetypesimple(rt)
		if rt.Kind() != reflect.Struct {
			return errNotStruct
		}

		rv = ref.Rdecodesimple(rv)
		childCtx := constructCtx{s.structValue, rt, rv, s.pc}
		err = s.constructFrom(childCtx)
	} else {
		childCtx := constructCtx{s.structValue, rt, rv, s.pc}
		err = s.constructFrom(childCtx)
	}
	return
}

func (s *sbS) constructFrom(ctx constructCtx) (err error) {
	for i := range ctx.rv.NumField() {
		frv := ctx.rv.Field(i)  // field value (reflect)
		frt := ctx.typ.Field(i) // field type
		tag := frt.Tag
		fieldName := frt.Name
		if fieldName == "" || unicode.IsLower([]rune(fieldName)[0]) {
			continue
		}

		cmdr := tag.Get("cmdr") // just for flag
		if cmdr == "-" {
			continue
		}
		cmdrSlice := strings.Split(cmdr, ",")
		positional := slices.Contains(cmdrSlice, "positional")

		title := nonEmpty(tag.Get("title"), tag.Get("name"))
		shorts := strings.Split(nonEmpty(tag.Get("shorts"), tag.Get("short")), ",")
		alias := strings.Split(nonEmpty(tag.Get("alias"), tag.Get("aliases")), ",")
		desc := nonEmpty(tag.Get("desc"), tag.Get("help"))
		group := tag.Get("group")

		// _, _, _, _, _, _, _ = frv, title, shorts, alias, desc, group, required
		title, shortTitle, shortTitles, titles := s.asmTitles(title, fieldName, shorts, alias...)

		if positional {
			if _, ok := frv.Interface().([]string); ok {
				if frv.CanAddr() {
					varptr := frv.Addr().Interface().(*[]string)
					logz.Trace(fmt.Sprintf("[constructFrom]       bind positional args ptr %p to Field %q", varptr, fieldName))
					s.cmd.BindPositionalArgsPtr(varptr)
				} else {
					logz.Warn("[constructFrom] CANNOT bind a field to parsed positional args because it cannot be addressed.", "rv", ref.Valfmtv(frv))
				}
			}
			continue
		}

		if frv.Kind() == reflect.Struct {
			// embedded struct -> command
			logz.Trace("[constructFrom] embedded STRUCT -> command", "Field", fieldName, "TgtCmd", title, "parent", ref.Valfmt(&ctx.rv))
			// s.parent.Long = title
			// s.parent.Short = shortTitle
			// s.parent.SetShorts(shortTitles...)
			// s.parent.Aliases = alias
			// s.parent.SetDesc(desc)

			titles := append([]string{shortTitle}, titles...)
			// if inCmd := atomic.LoadInt32(&s.inCmd); inCmd != 0 {
			// 	panic("cannot call Cmd() without Build() last StructBuilder")
			// }
			// atomic.AddInt32(&s.inCmd, 1)
			if atomic.CompareAndSwapInt32(&s.inCmd, 0, 1) == false {
				panic("cannot call Cmd() without Build() last StructBuilder")
			}
			var cb = newCommandBuilderFrom(s.cmd, s, title, titles...)
			logz.Trace("[constructFrom]     | applying command-builder", "ccb.CmdS", cb.CmdS, "ccb.parent", cb.parent)
			cb.ExtraShorts(shortTitles...).
				Group(group).
				Description(desc)
			// logz.Trace(fmt.Sprintf("frv.typ: %v", ref.Typfmt(frv.Type())))
			// logz.Trace(fmt.Sprintf("frt    : %v", ref.Typfmt(frt.Type)))
			// for i := 0; i < frv.Type().NumMethod(); i++ {
			// 	mtd := frv.Type().Method(i)
			// 	logz.Trace(fmt.Sprintf("method #%v: %v", mtd.Index, mtd.Name))
			// }
			if mtd := frv.MethodByName("Action"); mtd.IsValid() {
				cb.OnAction(func(ctx context.Context, cmd cli.Cmd, args []string) (err error) {
					ret := mtd.Call([]reflect.Value{
						reflect.ValueOf(ctx),
						reflect.ValueOf(cmd),
						reflect.ValueOf(args),
					})
					var ok bool
					if err, ok = ret[0].Interface().(error); !ok {
						if ret[0].Interface() != nil {
							logz.Fatal("expecting Action() returning error object", "actual", ref.Valfmtv(ret[0]))
						}
					}
					return
				})
			}
			if mtd := frv.MethodByName("With"); mtd.IsValid() {
				ret := mtd.Call([]reflect.Value{reflect.ValueOf(cb)})
				_ = ret
			}

			// entering for the embedded struct
			childStructValue := frv.Interface()
			if frv.Kind() != reflect.Ptr && frv.CanAddr() {
				childStructValue = frv.Addr().Interface()
			}
			childBuilder := newStructBuilderFrom(s.pc, cb.CmdS, cb, childStructValue)
			childBuilder.titleFormatter = s.titleFormatter
			logz.Trace("[constructFrom]     | applying struct-builder for child-struct-value", "cb.cmd", childBuilder.cmd, "cb.parent-from", cb.parent, "child-struct-value", childStructValue)
			childBuilder.Build()
			logz.Trace("[constructFrom]     | applied struct-builder", "cb.cmd", childBuilder.cmd, "cb.parent-from", cb.parent, "child-struct-value", childStructValue)

			cb.Build()
			logz.Trace("[constructFrom]     | applied command-builder", "cb.cmd", childBuilder.cmd, "cb.parent-from", cb.parent, "child-struct-value", childStructValue)
			continue
		}

		// normal field -> flag
		logz.Trace("[constructFrom]   normal field -> flag", "Field", fieldName, "TgtFlg", title, "owner-cmd", s.cmd, "parent-of-owner-cmd", s.cmd.OwnerCmd())
		if inFlg := atomic.LoadInt32(&s.inFlg); inFlg != 0 {
			panic("cannot call Flg() without Build() last StructBuilder")
		}
		atomic.AddInt32(&s.inFlg, 1)
		var fb = newFlagBuilderFrom(s.cmd, s, frv.Interface(),
			title, append([]string{shortTitle}, titles...)...)
		fb.ExtraShorts(shortTitles...).
			Group(group).
			Description(desc)
			// DefaultValue(frv.Interface()).

		required := is.StringToBool(tag.Get("required"))
		if required {
			fb.Required(required)
		}

		envvars := strings.Split(nonEmpty(tag.Get("env"), tag.Get("envvars")), ",")
		if len(envvars) > 0 {
			fb.EnvVars(envvars...)
		}

		headLikeA := strings.Split(nonEmpty(tag.Get("head-like"), tag.Get("headLike")), ",")
		if len(headLikeA) > 1 {
			headLike := is.StringToBool(headLikeA[0])
			var bounds []int
			for _, t := range headLikeA[1:] {
				i, _ := strconv.Atoi(t)
				bounds = append(bounds, i)
			}
			fb.HeadLike(headLike, bounds...)
		}

		var val any
		if frv.CanAddr() {
			val = frv.Addr().Interface()
		} else {
			val = frv.Interface()
		}
		fb.BindVarPtr(val)

		if shortTitle == "" {
			fb.Short = title // set short-title with long-title if user omitted it
		}
		if mtd := frv.MethodByName(title + "With"); mtd.IsValid() {
			ret := mtd.Call([]reflect.Value{reflect.ValueOf(fb)})
			_ = ret
		}
		fb.Build()
	}
	return
}

func (s *sbS) asmTitles(title, fieldName string, shorts []string, alias ...string) (longTitle, shortTitle string, shortTitles, titles []string) {
	longArray := append(append([]string{title, fieldName, shortTitle}, shortTitles...), titles...)
	longTitle = s.titleFormatter(nonEmpty(longArray...))
	if len(shorts) > 0 {
		shortTitle, shortTitles = s.titleFormatter(shorts[0]), shorts[1:]
	} else {
		shortTitle = fieldName
	}
	titles = alias
	return
}

func nonEmpty(ss ...string) string {
	// if s1 != "" {
	// 	return s1
	// }
	// if s2 != "" {
	// 	return s2
	// }
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

// func extractSectionFromTag(tag, tagName string) (result string) {
// 	for _, s := range strings.Split(tag, " ") {
// 		a := strings.Split(s, ":")
// 		if len(a) > 0 && a[0] == tagName {
// 			if len(a) > 1 {
// 				result = a[1]
// 			}
// 			return
// 		}
// 	}
// 	return
// }

type constructCtx struct {
	value any
	typ   reflect.Type
	rv    reflect.Value
	pc    *constructCtx
}

var (
	errNotStruct = errors.New("structValue is not a struct-based value")
)
