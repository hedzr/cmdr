package colortable

import (
	"fmt"
	"github.com/hedzr/cmdr"
)

func WithColorTableCommand() cmdr.ExecOption {
	return cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
		addTo(cmdr.RootFrom(root))
	}, nil)
}

func addTo(root *cmdr.RootCmdOpt) {

	c := cmdr.NewSubCmd()

	c.Titles("color-table", "ct").
		Description("print shell escape sequence table", "").
		Group(cmdr.SysMgmtGroup).
		Hidden(hidden).
		Action(printColorTable).
		AttachTo(root)

	const tgName = "Color Bits"
	cmdr.NewBool(true).
		Titles("4", "4").
		Description("3/4 bit color table", "").
		ToggleGroup(tgName).
		AttachTo(c)

	cmdr.NewBool().
		Titles("8", "8").
		Description("enable profiling", "").
		ToggleGroup(tgName).
		Hidden(true).
		AttachTo(c)

	cmdr.NewBool().
		Titles("24", "24").
		Description("enable profiling", "").
		ToggleGroup(tgName).
		Hidden(true).
		AttachTo(c)

	//root.AddGlobalPreAction(onCommandInvoking)
	//root.AddGlobalPostAction(afterCommandInvoked)

}

func printColorTable(cmd *cmdr.Command, args []string) (err error) {
	//println("yes")

	var table = []struct {
		name   string
		fg, bg int
	}{
		{"Black", 30, 40},
		{"Red", 31, 41},
		{"Green", 32, 42},
		{"Yellow", 33, 43},
		{"Blue", 34, 44},
		{"Magenta", 35, 45},
		{"Cyan", 36, 46},
		{"White", 37, 47},
		{"Bright Black", 90, 100},
		{"Bright Red", 91, 101},
		{"Bright Green", 92, 102},
		{"Bright Yellow", 93, 103},
		{"Bright Blue", 94, 104},
		{"Bright Magenta", 95, 105},
		{"Bright Cyan", 96, 106},
		{"Bright White", 97, 107},
	}

	switch bits := cmdr.GetStringR(cmd.GetDottedNamePath() + ".Color Bits"); bits {
	case "4":
		printColorTable4(table)
	case "24":
		printColorTable24(table)
	default:
		fmt.Printf("%q not supported\n", bits)
	}
	return
}

func printColorTable24(table []struct {
	name   string
	fg, bg int
}) {
	str := "#"
	for r := 0; r < 256; r++ {
		for g := 0; g < 256; g++ {
			for b := 0; b < 256; b++ {
				fmt.Printf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, str)
			}
			fmt.Println()
		}
		fmt.Println()
	}
	return
}

func printColorTable4(table []struct {
	name   string
	fg, bg int
}) {
	fmt.Printf("%16s%5s%5s  %-16s %-s\n", "Fg Name", "Fg", "Bg", "Sample", "More Bg...")

	for _, it := range table {
		fmt.Printf("%16s%5d%5d  ", it.name, it.fg, it.bg)

		str := fmt.Sprintf("(%d) Hello World", it.fg)
		fmt.Printf("\x1b[%dm%-16s\x1b[0m ", it.fg, str)

		for j, bgit := range table {
			if j < 8 {
				str := fmt.Sprintf("(%d;%d)", bgit.bg, it.fg)
				fmt.Printf("\x1b[%d;%dm%-7s\x1b[0m ", bgit.bg, it.fg, str)
			} else {
				str := fmt.Sprintf("(%d;%d)", bgit.bg, it.fg)
				fmt.Printf("\x1b[%d;%dm%-8s\x1b[0m ", bgit.bg, it.fg, str)
			}
		}
		fmt.Println()
	}

	println("\n")
	fmt.Printf("%16s%5s  %-20s %-s\n", "Fg Name", "Fg", "Bold,Underline Text", "Dim,Italic Text")
	for _, it := range table {
		fmt.Printf("%16s%5d  ", it.name, it.fg)

		str := fmt.Sprintf("(%d,%d,%d) Hello World", it.fg, 1, 4)
		fmt.Printf("\x1b[%d;%d;%dm%-16s\x1b[0m ", it.fg, 1, 4, str)

		str = fmt.Sprintf("(%d,%d,%d) Hello World", it.fg, 2, 3)
		fmt.Printf("\x1b[%d;%d;%dm%-16s\x1b[0m ", it.fg, 2, 3, str)

		fmt.Println()
	}
	return
}

var hidden bool
