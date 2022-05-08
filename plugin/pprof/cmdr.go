/*
 * Copyright © 2021 Hedzr Yeh.
 */

package pprof

import (
	"strings"

	"github.com/hedzr/cmdr"
)

// AttachToCmdr attaches the profiling options to root command.
//
// For example:
//
//    cmdr.Exec(buildRootCmd())
//
//    func buildRootCmd() (rootCmd *cmdr.RootCommand) {
//        root := cmdr.Root(appName, cmdr.Version).
//            Copyright(copyright, "hedzr").
//            Description(desc, longDesc).
//            Examples(examples)
//        rootCmd = root.RootCommand()
//        pprof.AttachToCmdr(root.RootCmdOpt())
//        return
//    }
//
// The variadic param `types` allows which profiling types are enabled
// by default, such as `cpu`, `mem`. While end-user enables profiling
// with `app -ep`, those profiles will be created and written.
// The available type names are:
// "cpu", "mem", "mutex", "block", "thread-create", "trace", and "go-routine".
func AttachToCmdr(root *cmdr.RootCmdOpt, types ...string) {
	sTypes = append(sTypes, types...)
	addTo(root)
}

// GetCmdrProfilingOptions returns an opt for cmdr.Exec(root, opts...). And,
// it adds profiling options to cmdr system.
//
// For example:
//
//    cmdr.Exec(buildRootCmd(),
//        profile.GetCmdrProfilingOptions(),
//    )
//
// The variadic param `types` allows which profiling types are enabled
// by default, such as `cpu`, `mem`. While end-user enables profiling
// with `app -ep`, those profiles will be created and written.
// The available type names are:
// "cpu", "mem", "mutex", "block", "thread-create", "trace", and "go-routine".
func GetCmdrProfilingOptions(types ...string) cmdr.ExecOption {
	sTypes = append(sTypes, types...)
	return optAddCPUProfileOptions
}

// WithCmdrProfilingOptions returns an opt for cmdr.Exec(root, opts...). And,
// it adds profiling options to cmdr system.
//
// For example:
//
//    cmdr.Exec(buildRootCmd(),
//        profile.GetCmdrProfilingOptions(),
//    )
//
// The variadic param `types` allows which profiling types are enabled
// by default, such as `cpu`, `mem`. While end-user enables profiling
// with `app -ep`, those profiles will be created and written.
// The available type names are:
// "cpu", "mem", "mutex", "block", "thread-create", "trace", and "go-routine".
func WithCmdrProfilingOptions(types ...string) cmdr.ExecOption {
	return GetCmdrProfilingOptions(types...)
}

// WithCmdrProfilingOptionsHidden hides the commands and flags from help screen
func WithCmdrProfilingOptionsHidden(types ...string) cmdr.ExecOption {
	hidden = true
	return WithCmdrProfilingOptions(types...)
}

func addTo(root *cmdr.RootCmdOpt) {
	const grpName = "Profiling"
	const grpOptsName = "Profiling-Options"

	cmdr.NewBool().
		Titles("pprof", "ep", "prof", "enable-profile").
		Description("enable profiling", "").
		Group(grpName).
		Hidden(hidden).
		AttachTo(root)

	cmdr.NewString("cpu.prof").
		Titles("cpu-profile-path", "cpu").
		Description("enable cpu profiling mode", "").
		Group(grpOptsName).
		Hidden(hidden).
		AttachTo(root)
	cmdr.NewString("mem.prof").
		Titles("mem-profile-path", "mem").
		Description("enable mem profiling mode", "").
		Group(grpOptsName).
		Hidden(hidden).
		AttachTo(root)
	cmdr.NewString("mutex.prof").
		Titles("mutex-profile-path", "mutex").
		Description("enable mutex profiling mode", "").
		Group(grpOptsName).
		Hidden(hidden).
		AttachTo(root)
	cmdr.NewString("block.prof").
		Titles("block-profile-path", "block").
		Description("enable block profiling mode", "").
		Group(grpOptsName).
		Hidden(hidden).
		AttachTo(root)
	cmdr.NewString("trace.out").
		Titles("trace-profile-path", "tpp").
		Description("enable trace profiling mode", "").
		Group(grpOptsName).
		Hidden(hidden).
		AttachTo(root)
	cmdr.NewString("thread-create.prof").
		Titles("thread-create-profile-path", "thread-create", "tc").
		Description("enable thread create profiling mode", "").
		Group(grpOptsName).
		Hidden(hidden).
		AttachTo(root)
	cmdr.NewString("go-routine.prof").
		Titles("go-routine-profile-path", "go-routine").
		Description("enable go routine profiling mode", "").
		Group(grpOptsName).
		Hidden(hidden).
		AttachTo(root)

	cmdr.NewString("heap").
		Titles("mem-profile-type", "mpt").
		Description("the profile type for memory profiles, 'heap' or 'allocs'", "").
		Placeholder("TYPE").
		Group("Profiling-Memory").
		Hidden(hidden).
		AttachTo(root)
	cmdr.NewInt(DefaultMemProfileRate).
		Titles("mem-profile-rate", "mpr").
		Description("the rate for the memory profile", "").
		Placeholder("RATE").
		Group("Profiling-Memory").
		Hidden(hidden).
		AttachTo(root)

	if len(sTypes) == 0 {
		sTypes = append(sTypes, "cpu", "mem", "mutex") //	"block", "trace", "thread-create", "go-routine"
	}
	cmdr.NewStringSlice(sTypes...).
		Titles("profiling-types", "pt").
		Description("specify types of profiling, such as: cpu, mem, mutex, block, trace, thread-create, goroutine", "").
		Group(grpName).
		ValidArgs("cpu", "mem", "mutex", "block", "trace", "thread-create", "goroutine").
		Hidden(hidden).
		AttachTo(root)

	cmdr.NewString(".").
		Titles("profile-output-dir", "pod").
		Description("the output directory", "").
		Group(grpName).
		Hidden(hidden).
		AttachTo(root)

	root.AddGlobalPreAction(onCommandInvoking)
	root.AddGlobalPostAction(afterCommandInvoked)
}

func onCommandInvoking(cmd *cmdr.Command, args []string) (err error) {
	if cmdr.GetBoolR("pprof") {
		var types ProfType
		for _, str := range cmdr.GetStringSliceR("profiling-types") {
			switch strings.ToLower(str) {
			case "cpu":
				types |= CPUProf
			case "mem", "memory":
				types |= MemProf
			case "mutex":
				types |= MutexProf
			case "block":
				types |= BlockProf
			case "thread-create", "thread-creation":
				types |= ThreadCreateProf
			case "trace":
				types |= TraceProf
			case "goroutine", "go-routine":
				types |= GoRoutineProf
			}
		}
		closer = Start(types,
			WithOutputDirectory(cmdr.GetStringR("profile-output-dir")),

			WithMemProfileRate(cmdr.GetIntR("mem-profile-rate")),
			WithMemProfileType(cmdr.GetStringR("mem-profile-type")),

			WithCPUProfName(cmdr.GetStringR("cpu-profile-path")),
			WithMemProfName(cmdr.GetStringR("mem-profile-path")),
			WithMutexProfName(cmdr.GetStringR("mutex-profile-path")),
			WithBlockProfName(cmdr.GetStringR("block-profile-path")),
			WithThreadCreateProfName(cmdr.GetStringR("thread-create-profile-path")),
			WithTraceProfName(cmdr.GetStringR("trace-profile-path")),
			WithGoRoutineProfName(cmdr.GetStringR("go-routine-profile-path")),
		)
	}
	return
}

func afterCommandInvoked(cmd *cmdr.Command, args []string) {
	if closer != nil {
		closer()
	}
}

func init() {
	optAddCPUProfileOptions = cmdr.WithXrefBuildingHooks(func(root *cmdr.RootCommand, args []string) {
		addTo(cmdr.RootFrom(root))
	}, nil)
}

var (
	optAddCPUProfileOptions cmdr.ExecOption
	sTypes                  []string
	closer                  func()
	hidden                  bool
)
