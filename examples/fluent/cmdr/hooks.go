package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"gopkg.in/hedzr/errors.v3"
	"strings"
)

func onSwitchCharHit(parsed *cmdr.Command, switchChar string, args []string) (err error) {
	if parsed != nil {
		fmt.Printf("the last parsed command is %q - %q\n", parsed.GetTitleNames(), parsed.Description)
	}
	fmt.Printf("SwitchChar FOUND: %v\nremains: %v\n\n", switchChar, args)
	return // cmdr.ErrShouldBeStopException
}

func onPassThruCharHit(parsed *cmdr.Command, switchChar string, args []string) (err error) {
	if parsed != nil {
		fmt.Printf("the last parsed command is %q - %q\n", parsed.GetTitleNames(), parsed.Description)
	}
	fmt.Printf("PassThrough flag FOUND: %v\nremains: %v\n\n", switchChar, args)
	return // ErrShouldBeStopException
}

func onUnhandledErrorHandler(err interface{}) {
	// debug.PrintStack()
	// pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	dumpStacks()
}

func dumpStacks() {
	fmt.Printf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===\n", errors.DumpStacksAsString(true))
}

func onUnknownOptionHandler(isFlag bool, title string, cmd *cmdr.Command, args []string) (fallbackToDefaultDetector bool) {
	return true
}

func onOptionMergeModifying(keyPath string, value, oldVal interface{}) {
	cmdr.Logger.Printf("%%-> -> %q: %v -> %v", keyPath, oldVal, value)
	if strings.HasSuffix(keyPath, ".mqtt.server.stats.enabled") {
		// mqttlib.FindServer().EnableSysStats(!vxconf.ToBool(value))
	}
	if strings.HasSuffix(keyPath, ".mqtt.server.stats.log.enabled") {
		// mqttlib.FindServer().EnableSysStatsLog(!vxconf.ToBool(value))
	}
}
