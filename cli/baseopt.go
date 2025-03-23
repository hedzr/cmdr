package cli

import (
	"bufio"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/is/term/color"
	"github.com/hedzr/store"

	"github.com/hedzr/cmdr/v2/conf"
)

func (c *BaseOpt) Clone() any {
	return &BaseOpt{
		owner:        c.owner,
		root:         c.root,
		name:         c.name,
		Long:         c.Long,
		Short:        c.Short,
		Aliases:      slices.Clone(c.Aliases),
		description:  c.description,
		longDesc:     c.longDesc,
		examples:     c.examples,
		group:        c.group,
		extraShorts:  slices.Clone(c.extraShorts),
		deprecated:   c.deprecated,
		hidden:       c.hidden,
		vendorHidden: c.vendorHidden,
		hitTitle:     c.hitTitle,
		hitTimes:     c.hitTimes,
	}
}

func (c *BaseOpt) errIsSignalFallback(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrShouldFallback)
}

func (c *BaseOpt) SetName(name string) {
	c.name = name
	if c.Long == "" {
		c.Long = name
	}
}

func (c *BaseOpt) SetShorts(shorts ...string) {
	c.extraShorts = append(c.extraShorts, shorts...)
}

func (c *BaseOpt) SetDescription(description string, longDescription ...string) {
	c.description = description
	var (
		lines []string
		lead  = true
	)
	for _, str := range longDescription {
		if s := strings.Trim(str, "\r\n\t "); s != "" {
			lines = append(lines, s)
		} else if !lead {
			lines = append(lines, s)
		} else {
			lead = false
			lines = append(lines, s)
		}
	}
	end := len(lines) - 1
	for i := end; i >= 0; i-- {
		if s := strings.Trim(lines[i], "\r\n\t "); s == "" {
			end--
		}
	}
	lines = lines[:end+1]
	c.longDesc = strings.Join(lines, "\n")
	if description == "" && len(c.longDesc) > 0 {
		c.description = firstNonBlankLine(c.longDesc)
	}
}

func firstNonBlankLine(desc string) string {
	scanner := bufio.NewScanner(strings.NewReader(desc))
	for scanner.Scan() {
		line := scanner.Text()
		if s := strings.Trim(line, "\r\n\t "); s != "" {
			return line
		}
	}
	return ""
}

func (c *BaseOpt) SetExamples(examples ...string) {
	c.examples = strings.Join(examples, "\n")
}

func (c *BaseOpt) SetGroup(group string) {
	c.group = group
}

func (c *BaseOpt) SetDeprecated(deprecated string) {
	c.deprecated = deprecated
}

func (c *BaseOpt) SetHidden(hidden bool, vendorHidden ...bool) {
	c.hidden = hidden
	for _, b := range vendorHidden {
		c.vendorHidden = b
	}
}

// func (c *BaseOpt) Owner() *CmdS                    { return c.owner }            // the owner of this CmdS

func (c *BaseOpt) OwnerOrParent() BacktraceableMin { return c.owner }
func (c *BaseOpt) OwnerIsNil() bool                { return c.owner == nil }
func (c *BaseOpt) OwnerIsNotNil() bool             { return c.owner != nil }
func (c *BaseOpt) OwnerCmd() Cmd                   { return c.owner }
func (c *BaseOpt) OwnerIsRoot() bool               { return c.owner != nil && c.owner.OwnerIsNil() }
func (c *BaseOpt) IsRoot() bool                    { return c.owner == nil }
func (c *BaseOpt) Root() *RootCommand              { return c.root }          // returns Root CmdS (*RootCommand),
func (c *BaseOpt) App() App                        { return c.root.app }      // App returns the current App
func (c *BaseOpt) Set() store.Store                { return c.App().Store() } // Set returns the application Store [store.Store]
func (c *BaseOpt) SetOwner(o *CmdS)                { c.owner = o }
func (c *BaseOpt) SetOwnerCmd(o Cmd)               { c.owner = o }
func (c *BaseOpt) SetRoot(root *RootCommand)       { c.root = root }

// Store returns the commands subset of the application Store.
func (c *BaseOpt) Store() store.Store {
	cs := c.Set().WithPrefix(CommandsStoreKey, c.GetDottedPath())
	return cs
}

// func (c *BaseOpt) AppName() string {
// 	if conf.AppName != "" {
// 		return conf.AppName
// 	}
// 	if c.root.name != "" {
// 		return c.root.name
// 	}
// 	return c.root.AppName
// }

func (c *BaseOpt) AppVersion() string {
	if conf.Version != "" {
		return conf.Version
	}
	return c.root.Version
}

func (c *BaseOpt) Title() string {
	if c.name != "" {
		return c.name
	}
	if c.Long != "" {
		return c.Long
	}
	if c.Short != "" {
		return c.Short
	}
	for _, s := range c.Aliases {
		if s != "" {
			return s
		}
	}
	return "> ? <"
}

// Shorts collect and return all short titles
// as one array without duplicated items.
//
// include both the internal Short and extraShorts field.
func (c *BaseOpt) Shorts() (shorts []string) {
	if c.Short != "" {
		shorts = append(shorts, c.Short)
	}
	shorts = uniAddStrS(shorts, c.extraShorts...)
	// shorts = append(shorts, c.extraShorts...)
	return
}

// GetName returns the name of a `CmdS`.
func (c *BaseOpt) GetName() string {
	if len(c.name) > 0 {
		return c.name
	}
	if len(c.Long) > 0 {
		return c.Long
	}
	panic("The `Long` or `Name` must be non-empty for a command or flag")
}

// Name returns the identity string of this command/flag, long title or name only
func (c *BaseOpt) Name() string {
	if c.name != "" {
		return c.name
	}
	if c.Long != "" {
		return c.Long
	}
	return ""
}

func (c *BaseOpt) String() string {
	var sb strings.Builder
	_, _ = sb.WriteString("BaseOpt{'")
	_, _ = sb.WriteString(c.GetTitleName())
	_, _ = sb.WriteString("'}")
	return sb.String()
}

// HasParent detects whether owner is available or not
func (c *BaseOpt) HasParent() bool { return c.owner != nil }

// GetHitStr returns the matched command string
func (c *BaseOpt) GetHitStr() string { return c.hitTitle }

// GetTriggeredTimes returns the matched times
func (c *BaseOpt) GetTriggeredTimes() int { return c.hitTimes }

// GetDottedPath return the dotted key path of this command
// in the options store.
// For example, the returned string just like: 'server.start'.
// NOTE that there is no OptionPrefixes in this key path. For
// more information about Option Prefix, refer
// to [WithOptionsPrefix]
func (c *BaseOpt) GetDottedPath() string {
	return strings.ReplaceAll(strings.ReplaceAll(backtraceCmdNamesG(c, ".", false), "<root>.", ""), "<root>", "")
}

func (c *BaseOpt) GetDottedPathFull() string {
	return strings.ReplaceAll(strings.ReplaceAll(backtraceCmdNamesG(c, ",", true), "<root>.", ""), "<root>", "")
}

func (c *BaseOpt) GetCommandTitles() string {
	return strings.ReplaceAll(strings.ReplaceAll(backtraceCmdNamesG(c, " ", false), "<root>.", ""), "<root>", "")
}

func (c *BaseOpt) GetAutoEnvVarName(prefix string, upperCase ...bool) string {
	t := backtraceCmdNamesG(c, "_", false)
	// last := c.Name()
	u := false
	for _, b := range upperCase {
		u = b
	}
	if u {
		t = strings.ToUpper(t)
		// last = strings.ToUpper(last)
	}
	if prefix != "" {
		return prefix + "_" + t // + "_" + last
	}
	return t // + "_" + last
}

// GetTitleName returns name/full/short string
func (c *BaseOpt) GetTitleName() string {
	if c.name != "" {
		if c.owner == nil {
			return "<root>"
		}
		return c.name
	}
	if c.Long != "" {
		return c.Long
	}
	if c.Short != "" {
		return c.Short
	}
	// for _, ss := range s.Aliases {
	// 	return ss
	// }
	return ""
}

func (c *BaseOpt) Desc() string {
	if c.description != "" {
		return c.description
	}
	return c.longDesc
}

func (c *BaseOpt) DescLong() string {
	if c.longDesc != "" {
		return c.longDesc
	}
	return c.description
}

func (c *BaseOpt) SetDesc(desc string) {
	c.description = desc
}

func (c *BaseOpt) Examples() string   { return c.examples }
func (c *BaseOpt) Deprecated() string { return c.deprecated }
func (c *BaseOpt) Hidden() bool       { return c.hidden }
func (c *BaseOpt) VendorHidden() bool { return c.vendorHidden }

// SafeGroup return UnsortedGroup if group member not set yet.
func (c *BaseOpt) SafeGroup() string {
	if c.group == "" {
		return UnsortedGroup
	}
	return c.group
}

func (c *BaseOpt) RemoveOrderedPrefix(title string) string {
	return reSortingPrefix.ReplaceAllString(title, "")
}

func RemoveOrderedPrefix(title string) string {
	return reSortingPrefix.ReplaceAllString(title, "")
}

// GroupTitle returns the group title without leading
// ordered pieces.
//
// The ordered prefix of returned title (such as "00ab." in
// "00ab.Group A") was removed, the final title would be
// "Group A".
func (c *BaseOpt) GroupTitle() string {
	return c.RemoveOrderedPrefix(c.SafeGroup())
}

// GroupHelpTitle returns the group title or empty string if
// it's UnsortedGroup.
//
// The title will be printed in help screen. Its ordered prefix
// (such as "00ab." in "00ab.Group A") was removed.
func (c *BaseOpt) GroupHelpTitle() string {
	tmp := c.SafeGroup()
	if tmp == UnsortedGroup {
		return ""
	}
	return c.RemoveOrderedPrefix(tmp)
}

// GetTitleNamesArray returns short,full,aliases names
//
// A title with prefix '_'/'__' will be hidden from the
// result array.
func (c *BaseOpt) GetTitleNamesArray() []string {
	a := c.GetTitleNamesArrayMainly()
	for _, x := range c.extraShorts {
		if x != "" {
			a = uniAddStr(a, x)
		}
	}
	for _, x := range c.Aliases {
		if x == "" || strings.HasPrefix(x, "_") {
			continue
		}
		a = uniAddStr(a, x)
	}
	// a = uniAddStrS(a, c.Aliases...)
	return a
}

// GetTitleNamesArrayMainly returns short,full names
func (c *BaseOpt) GetTitleNamesArrayMainly() []string {
	var a []string
	if c.Short != "" {
		a = uniAddStr(a, c.Short)
	}
	if c.name != "" {
		a = uniAddStr(a, c.name)
	} else if c.Long != "" {
		a = uniAddStr(a, c.Long)
	}
	return a
}

// GetShortTitleNamesArray returns short name as an array
func (c *BaseOpt) GetShortTitleNamesArray() []string {
	var a []string
	for _, ss := range c.Shorts() {
		a = uniAddStr(a, ss)
	}
	return a
}

// GetLongTitleNamesArray returns long name and aliases as an array
func (c *BaseOpt) GetLongTitleNamesArray() []string {
	var a []string
	if n := c.Name(); n != "" {
		a = uniAddStr(a, n)
	}
	if n := c.Long; n != "" {
		a = uniAddStr(a, n)
	}
	a = uniAddStrS(a, c.Aliases...)
	return a
}

// LongTitle returns long name
func (c *BaseOpt) LongTitle() string {
	if n := c.Long; n != "" {
		return n
	}
	if n := c.Name(); n != "" {
		return n
	}
	for _, n := range c.Aliases {
		return n
	}
	return "(noname)"
}

// ShortTitle _
func (c *BaseOpt) ShortTitle() string {
	if s := c.Short; s != "" {
		return s
	}
	for _, s := range c.extraShorts {
		if s != "" {
			return s
		}
	}
	return ""
}

// GetTitleNames return the joint string of short,full,aliases names
func (c *BaseOpt) GetTitleNames(maxWidth ...int) (title, rest string) {
	maxW := 0
	for _, x := range maxWidth {
		maxW = x
	}

	if maxW <= 0 {
		a := c.GetTitleNamesArray()
		title = strings.Join(a, ", ")
		return
	}

	var sb, sbr strings.Builder
	var delimiter = ""
	var ll = 3
	if c.Short != "" {
		sb.WriteString(c.Short)
		if len(c.Short) == 1 {
			delimiter = ", "
		} else {
			delimiter = ","
		}
		sb.WriteString(delimiter)
		delimiter = ","
	} else {
		sb.WriteString("   ")
	}
	if c.Long != "" {
		sb.WriteString(c.Long)
		delimiter = ","
		ll += len(c.Long) + 1
	}

	si := len(c.Aliases)
	for i, x := range c.Aliases {
		if x == "" || strings.HasPrefix(x, "_") {
			continue
		}

		if i > si {
			sbr.WriteString(delimiter)
			sbr.WriteString(x)
			continue
		}
		if ll+len(x)+1 >= maxW {
			si = i
			sbr.WriteString(x)
			continue
		}

		ll += len(x) + len(delimiter)
		sb.WriteString(delimiter)
		sb.WriteString(x)
	}

	if sbr.Len() > 0 {
		si = -1
	} else {
		si = 0
	}
	delimiter = ","
	for i, x := range c.extraShorts {
		if i > si {
			sbr.WriteString(delimiter)
		}
		sbr.WriteString(x)
	}

	title = sb.String()
	rest = sbr.String()
	return
	// return c.GetTitleNamesBy(", ", maxWidth...)
}

// GetTitleNamesBy returns the joint string of short,full,aliases names
func (c *BaseOpt) GetTitleNamesBy(delimiter string, maxWidth ...int) (title, rest string) {
	a := c.GetTitleNamesArray()
	maxW := 0
	for _, x := range maxWidth {
		maxW = x
	}
	if maxW > 0 {
		var sb, sbr strings.Builder
		var ri int = len(a)
		for i, x := range a {
			if i >= ri || sb.Len()+len(delimiter)+len(x) > maxW {
				if ri > i {
					ri = i
					sbr.WriteString(x)
				} else {
					sbr.WriteString(delimiter)
					sbr.WriteString(x)
				}
				continue
			}
			if i > 0 {
				sb.WriteString(delimiter)
			}
			sb.WriteString(x)
		}
		title, rest = sb.String(), sbr.String()
	} else {
		title = strings.Join(a, delimiter)
	}
	return
}

// GetTitleZshNames temp
func (c *BaseOpt) GetTitleZshNames() string {
	a := c.GetTitleNamesArrayMainly()
	str := strings.Join(a, ",")
	return str
}

// GetTitleZshNamesBy temp
func (c *BaseOpt) GetTitleZshNamesBy(delimiter string) (str string) {
	if c.Short != "" {
		str += c.Short + delimiter
	}
	if c.Long != "" {
		str += c.Long
	}
	return
}

// GetDescZsh temp
func (c *BaseOpt) GetDescZsh() (desc string) {
	desc = c.description
	if desc == "" {
		desc = EraseAnyWSs(c.GetTitleZshNames())
	}
	// desc = replaceAll(desc, " ", "\\ ")
	desc = reSQ.ReplaceAllString(desc, `*$1*`)
	desc = reBQ.ReplaceAllString(desc, `**$1**`)
	desc = reSQnp.ReplaceAllString(desc, "''")
	desc = reBQnp.ReplaceAllString(desc, "\\`")
	desc = strings.ReplaceAll(desc, ":", "\\:")
	desc = strings.ReplaceAll(desc, "[", "\\[")
	desc = strings.ReplaceAll(desc, "]", "\\]")
	return
}

// EraseAnyWSs eats any whitespaces inside the giving string s.
func EraseAnyWSs(s string) string {
	return reSimpSimp.ReplaceAllString(s, "")
}

func (c *BaseOpt) DeprecatedHelpString(trans func(ss string, clr color.Color) string, clr, clrDefault color.Color) (hs, plain string) {
	if c.deprecated != "" {
		re := regexp.MustCompile(`[Ss]ince:? `)
		dep := re.ReplaceAllString(c.deprecated, "")
		plain = fmt.Sprintf("[Since: %s]", dep)
		hs = trans(fmt.Sprintf("[Since: <font color=%v>%s</font>]", color.ToColorString(clr), dep), clrDefault)
	}
	return
}

var (
	reSQnp = regexp.MustCompile(`'`)
	reBQnp = regexp.MustCompile("`")
	reSQ   = regexp.MustCompile(`'(.*?)'`)
	reBQ   = regexp.MustCompile("`(.*?)`")
	// reULs  = regexp.MustCompile("_+")

	reSortingPrefix = regexp.MustCompile(`[0-9A-Za-z!#$%^&]+\.`)

	reSimpSimp = regexp.MustCompile(`[ \t]+`)
	// reSimp = regexp.MustCompile(`[ \t][ \t]+`)
)
