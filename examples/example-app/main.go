package main

import (
	"fmt"

	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/pprof"
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log"
	"github.com/hedzr/log/buildtags"
	"github.com/hedzr/log/detects"
	"github.com/hedzr/logex/build"
	"gopkg.in/hedzr/errors.v3"
)

func main() {
	Run()
}

func Run() {
	log.Fatal(cmdr.Exec(buildRootCmd(), options...)) // since hedzr/log 1.6.1, log.Fatal/Panic can ignore nil safely
	// root := buildRootCmd()
	// if err := cmdr.Exec(root, options...); err != nil {
	// 	log.Fatalf("error occurs in app running: %+v\n", err)
	// }
}

func buildRootCmd() (rootCmd *cmdr.RootCommand) {
	root := cmdr.Root(appName, version).
		// AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
		//	// cmdr.Set("enable-ueh", true)
		//	return
		// }).
		// AddGlobalPreAction(func(cmd *cmdr.Command, args []string) (err error) {
		//	//fmt.Printf("# global pre-action 2, exe-path: %v\n", cmdr.GetExecutablePath())
		//	return
		// }).
		// AddGlobalPostAction(func(cmd *cmdr.Command, args []string) {
		//	//fmt.Println("# global post-action 1")
		// }).
		// AddGlobalPostAction(func(cmd *cmdr.Command, args []string) {
		//	//fmt.Println("# global post-action 2")
		// }).
		Copyright(copyright, "hedzr").
		Description(desc, longDesc).
		Examples(examples)
	rootCmd = root.RootCommand()

	// for your biz-logic, constructing an AttachToCmdr(root *cmdr.RootCmdOpt) is recommended.
	// see our full sample and template repo: https://github.com/hedzr/cmdr-go-starter
	// core.AttachToCmdr(root.RootCmdOpt())

	// These lines are removable

	cmdr.NewBool(false).
		Titles("enable-ueh", "ueh").
		Description("Enables the unhandled exception handler?").
		AttachTo(root)
	// cmdrPanic(root)
	cmdrSoundex(root)
	// pprof.AttachToCmdr(root.RootCmdOpt())
	return
}

func cmdrSoundex(root cmdr.OptCmd) {
	cmdr.NewSubCmd().Titles("soundex", "snd", "sndx", "sound").
		Description("soundex test").
		Group("Test").
		TailPlaceholder("[text1, text2, ...]").
		Action(func(cmd *cmdr.Command, args []string) (err error) {
			for ix, s := range args {
				fmt.Printf("%5d. %s => %s\n", ix, s, tool.Soundex(s))
			}
			return
		}).
		AttachTo(root)
}

func onUnhandledErrorHandler(err interface{}) {
	if cmdr.GetBoolR("enable-ueh") {
		dumpStacks()
	}

	panic(err) // re-throw it
}

func dumpStacks() {
	fmt.Printf("\n\n=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===\n\n", errors.DumpStacksAsString(true))
}

func init() {
	options = append(options,
		cmdr.WithUnhandledErrorHandler(onUnhandledErrorHandler),

		// cmdr.WithLogxShort(defaultDebugEnabled,defaultLoggerBackend,defaultLoggerLevel),
		cmdr.WithLogx(build.New(build.NewLoggerConfigWith(
			defaultDebugEnabled, defaultLoggerBackend, defaultLoggerLevel,
			log.WithTimestamp(true, "")))),

		cmdr.WithHelpTailLine(`
# Type '-h'/'-?' or '--help' to get command help screen.
# Star me if it's helpful: https://github.com/hedzr/cmdr/examples/example-app
`),
	)

	if isDebugBuild() {
		options = append(options, pprof.GetCmdrProfilingOptions())
	}

	// enable '--trace' command line option to toggle a internal trace mode (can be retrieved by cmdr.GetTraceMode())
	// import "github.com/hedzr/cmdr-addons/pkg/plugins/trace"
	// trace.WithTraceEnable(defaultTraceEnabled)
	// Or:
	optAddTraceOption := cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
		cmdr.NewBool(false).
			Titles("trace", "tr").
			Description("enable trace mode for tcp/mqtt send/recv data dump", "").
			// Action(func(cmd *cmdr.Command, args []string) (err error) { println("trace mode on"); cmdr.SetTraceMode(true); return; }).
			Group(cmdr.SysMgmtGroup).
			AttachToRoot(root)
	}, nil)
	options = append(options, optAddTraceOption)
	// options = append(options, optAddServerExtOpt«ion)

	// allow and search '.<appname>.yml' at first
	locations := []string{".$APPNAME.yml"}
	locations = append(locations, cmdr.GetPredefinedLocations()...)
	options = append(options, cmdr.WithPredefinedLocations(locations...))

	// options = append(options, internal.NewAppOption())
}

var options []cmdr.ExecOption

func isDebugBuild() bool         { return detects.InDebugging() }
func isDockerBuild() bool        { return buildtags.IsDockerBuild() }
func isRunningInDockerEnv() bool { return detects.InDocker() }

//goland:noinspection GoNameStartsWithPackageName
const (
	appName   = "example-app"
	version   = "0.2.5"
	copyright = "example-app - A devops tool - cmdr series"
	desc      = "example-app is an effective devops tool. It make an demo application for 'cmdr'"
	longDesc  = `example-app is an effective devops tool. It make an demo application for 'cmdr'.
`
	examples = `
$ {{.AppName}} gen shell [--bash|--zsh|--fish|--auto]
  generate bash/shell completion scripts
$ {{.AppName}} gen man
  generate linux man page 1
$ {{.AppName}} --help
  show help screen.
$ {{.AppName}} --help --man
  show help screen in manpage viewer (for linux/darwin).
`
	// overview = ``
	// zero = 0

	// defaultTraceEnabled  = true
	defaultDebugEnabled  = false
	defaultLoggerLevel   = "debug"
	defaultLoggerBackend = "logrus"
)
