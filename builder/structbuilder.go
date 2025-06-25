package builder

import (
	"context"
	"errors"
	"fmt"
	"reflect"
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
	return newStructBuilderFrom(nil, b, structValue, opts...)
}

func newStructBuilderShort(b buildable, structValue any, opts ...cli.StructBuilderOpt) *sbS {
	return newStructBuilderFrom(nil, b, structValue, opts...)
}

func newStructBuilderFrom(from *cli.CmdS, b buildable, structValue any, opts ...cli.StructBuilderOpt) *sbS {
	s := &sbS{
		0, 0,
		isAssumedAsRootCmd(assumedAsRootcmd),
		b,
		from,
		// new(cli.CmdS),
		structValue,
		stringtool.ToKebabCase,
	}

	// s.Long, s.Short, s.Aliases = theTitles(longTitle, titles...)

	if app, ok := b.(*appS); ok {
		// links app, owner, parent and this building cmd
		app.root.SetApp(app)
		s.parent = app.root.Cmd.(*cli.CmdS)
		s.parent.SetRoot(app.root)
		// s.CmdS.SetRoot(app.root)
		// s.CmdS.SetOwner(s.parent)
	} else if sb, ok := b.(*sbS); ok {
		s.parent = sb.Building()
		// 	s.CmdS.SetOwner(s.parent)
		// 	s.CmdS.SetRoot(s.parent.Root())
	} else if sb, ok := b.(*ccb); ok {
		s.parent = sb.CmdS
		// 	s.CmdS.SetOwner(s.parent)
		// 	s.CmdS.SetRoot(s.parent.Root())
	}

	// if s.asRoot {
	// 	// transfer app info from old rootcmd to this building cmd
	// 	c, p := s.CmdS, s.parent
	// 	c.Long = p.Long
	// 	c.Short = p.Short
	// 	c.SetName(p.Name())
	// 	c.SetDescription(p.Desc(), p.DescLong())
	// 	c.SetExamples(p.Examples())
	// 	c.SetGroup(p.SafeGroup())

	// 	if app, ok := b.(*appS); ok {
	// 		old := app.root.Cmd.(*cli.CmdS)
	// 		app.root.Cmd = s.CmdS
	// 		s.CmdS.SetCommands(old.SubCommands()...)
	// 		s.CmdS.SetFlags(old.Flags()...)
	// 		s.CmdS.SetRoot(app.root)
	// 		s.CmdS.SetOwner(nil) // remove owner ref pointer since the biulding cmd will be used as Root
	// 	}
	// }
	return s
}

type sbS struct {
	inCmd  int32
	inFlg  int32
	asRoot bool
	buildable
	parent *cli.CmdS
	// *cli.CmdS
	structValue    any
	titleFormatter TitleFormatFunc
}

func isAssumedAsRootCmd(title string) bool {
	return strings.HasPrefix(title, "(") && strings.HasSuffix(title, ")")
}

func (s *sbS) Buildable() cli.OptBuilder { return s.buildable }
func (s *sbS) Parent() *cli.CmdS         { return s.parent }
func (s *sbS) Building() *cli.CmdS       { return s.parent }

func (s *sbS) Build() {
	if err := s.construct(); err != nil {
		logz.Error("cannot construct cmdr command system from a struct value", "err", err)
		return
	}

	if s.asRoot {
		logz.Verbose(assumedAsRootcmd)
	} else {
		logz.Verbose("normal")
	}
	if a, ok := s.buildable.(adder); ok {
		a.addCommand(nil)
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
	s.parent.AddSubCommand(child)
	if child != nil {
		logz.Trace(fmt.Sprintf("added %v -> %v", child.String(), s.parent))
	}
}

// addFlag adds a in-building Flg into current CmdS as its flag.
// used by adder when ccb.Build.
func (s *sbS) addFlag(child *cli.Flag) {
	atomic.AddInt32(&s.inFlg, -1)
	s.parent.AddFlag(child)
	logz.Trace(fmt.Sprintf("added %v -> %v", child, s.parent))
}

func (s *sbS) construct() (err error) {
	var sv = s.structValue
	rt := reflect.TypeOf(sv)
	if rt.Kind() != reflect.Struct {
		rt = ref.Rdecodetypesimple(rt)
		if rt.Kind() != reflect.Struct {
			return errNotStruct
		}

		rv := ref.Rdecodesimple(reflect.ValueOf(sv))
		childCtx := constructCtx{sv, rt, rv}
		err = s.constructFrom(childCtx)
	} else {
		rv := reflect.ValueOf(sv)
		childCtx := constructCtx{sv, rt, rv}
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
		if tag.Get("cmdr") == "-" {
			continue
		}

		title := nonEmpty(tag.Get("title"), tag.Get("name"))
		shorts := strings.Split(nonEmpty(tag.Get("shorts"), tag.Get("short")), ",")
		alias := strings.Split(nonEmpty(tag.Get("alias"), tag.Get("aliases")), ",")
		desc := nonEmpty(tag.Get("desc"), tag.Get("help"))
		group := tag.Get("group")
		required := tag.Get("required") // just for flag

		// _, _, _, _, _, _, _ = frv, title, shorts, alias, desc, group, required
		title, shortTitle, shortTitles, titles := s.asmTitles(title, fieldName, shorts, alias...)

		if frv.Kind() == reflect.Struct {
			// embedded struct -> command
			logz.Trace("[constructFrom] embedded STRUCT -> command", "Field", fieldName, "TgtCmd", title, "parent", ref.Valfmt(&ctx.rv))
			// s.parent.Long = title
			// s.parent.Short = shortTitle
			// s.parent.SetShorts(shortTitles...)
			// s.parent.Aliases = alias
			// s.parent.SetDesc(desc)

			titles := append([]string{shortTitle}, titles...)
			if inCmd := atomic.LoadInt32(&s.inCmd); inCmd != 0 {
				panic("cannot call Cmd() without Build() last StructBuilder")
			}
			atomic.AddInt32(&s.inCmd, 1)
			var cb = newCommandBuilderShort(s, title, titles...).
				ExtraShorts(shortTitles...).
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
					err = ret[0].Interface().(error)
					return
				})
			}
			if mtd := frv.MethodByName("With"); mtd.IsValid() {
				ret := mtd.Call([]reflect.Value{reflect.ValueOf(cb)})
				_ = ret
			}

			// entering for the embedded struct
			childStructValue := frv.Interface()
			childBuilder := newStructBuilderShort(cb, childStructValue)
			childBuilder.titleFormatter = s.titleFormatter
			childBuilder.Build()

			cb.Build()
		} else {
			// normal field -> flag
			logz.Trace("[constructFrom] normal field -> flag", "Field", fieldName, "TgtFlg", title)
			if inFlg := atomic.LoadInt32(&s.inFlg); inFlg != 0 {
				panic("cannot call Flg() without Build() last StructBuilder")
			}
			atomic.AddInt32(&s.inFlg, 1)
			var fb = newFlagBuilderShort(s,
				title, append([]string{shortTitle}, titles...)...)
			fb.ExtraShorts(shortTitles...).
				Group(group).
				Description(desc).
				DefaultValue(frv.Interface()).
				Required(is.StringToBool(required))
			if shortTitle == "" {
				fb.Short = title // set short-title with long-title if user omitted it
			}
			if mtd := frv.MethodByName(title + "With"); mtd.IsValid() {
				ret := mtd.Call([]reflect.Value{reflect.ValueOf(fb)})
				_ = ret
			}
			fb.Build()
		}
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
}

var (
	errNotStruct = errors.New("structValue is not a struct-based value")
)
