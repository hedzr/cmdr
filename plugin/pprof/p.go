/*
 * Copyright © 2021 Hedzr Yeh.
 */

package pprof

import (
	"github.com/hedzr/log"
	"github.com/hedzr/log/trace"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	traceR "runtime/trace"
)

// ProfType represents the profile type
type ProfType int

const (
	// CPUProf represents cpu profiling
	CPUProf ProfType = 1
	// MemProf represents memory profiling
	MemProf ProfType = 2
	// MutexProf represents mutex profiling
	MutexProf ProfType = 4
	// BlockProf represents block profiling
	BlockProf ProfType = 8
	// ThreadCreateProf represents thread creation profiling
	ThreadCreateProf ProfType = 16
	// TraceProf represents trace profiling
	TraceProf ProfType = 32
	// GoRoutineProf represents go-routine profiling
	GoRoutineProf ProfType = 64
)

// DefaultMemProfileRate is the default memory profiling rate.
// See also http://golang.org/pkg/runtime/#pkg-variables
const DefaultMemProfileRate = 4096

// profile represents the profiling options holder
type profile struct {
	path string

	types                []ProfType
	cpuProfName          string
	memProfName          string
	mutexProfName        string
	blockProfName        string
	threadCreateProfName string
	traceProfName        string
	goRoutineProfName    string

	memProfileRate int
	// memProfileType holds the profile type for memory
	// profiles. Allowed values are `heap` and `allocs`.
	memProfileType string
}

// var cpuProfile, memProfile string

// EnableCPUProfile enables cpu profiling.
// And review the pprof result in a web ui:
//
//    go tool pprof -http=:8555 ./cpu.pprof
//
// Now you can open 'http://localhost:8555/ui' in a browser
//
func EnableCPUProfile(cpuProfilePath string) (closer func()) {
	closer = func() {}
	if cpuProfilePath != "" {
		f, err := os.Create(cpuProfilePath)
		if err != nil {
			log.Fatal("could not create cpu profile: %v", err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal("error: %v", err)
		}
		closer = pprof.StopCPUProfile
	}
	runtime.SetBlockProfileRate(20)
	return
}

// EnableMemProfile enables memory profiling.
// And review the pprof result in a web ui:
//
//    go tool pprof -http=:8555 ./mem.pprof
//
// Now you can open 'http://localhost:8555/ui' in a browser
//
func EnableMemProfile(memProfilePath string) (closer func()) {
	closer = func() {}
	if memProfilePath != "" {
		closer = func() {
			f, err := os.Create(memProfilePath)
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			defer f.Close()
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
		}
	}
	return
}

func (s *profile) enableCPUProfile(profilePath string) (closer func()) {
	closer = func() {}
	if profilePath != "" {
		fn := filepath.Join(s.path, profilePath)
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create cou profile %q: %v", fn, err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			return
		}
		log.Debugf("profile: cpu profiling enabled, %s", fn)
		closer = func() {
			pprof.StopCPUProfile()
			_ = f.Close()
			log.Debugf("profile: cpu profiling disabled, %s", fn)
		}
	}
	return
}

func (s *profile) enableMemProfile(profilePath string) (closer func()) {
	closer = func() {}
	if profilePath != "" {
		fn := filepath.Join(s.path, profilePath)
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create memory profile %q: %v", fn, err)
		}
		old := runtime.MemProfileRate
		runtime.MemProfileRate = s.memProfileRate
		if err != nil {
			return
		}
		log.Debugf("profile: memory profiling enabled (rate %d), %s", runtime.MemProfileRate, fn)
		closer = func() {
			_ = pprof.Lookup(s.memProfileType).WriteTo(f, 0)
			_ = f.Close()
			runtime.MemProfileRate = old
			log.Debugf("profile: memory profiling disabled, %s", fn)
		}
	}
	return
}

func (s *profile) enableMutexProfile(profilePath string) (closer func()) {
	closer = func() {}
	if profilePath != "" {
		fn := filepath.Join(s.path, profilePath)
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create mutex profile %q: %v", fn, err)
		}
		runtime.SetMutexProfileFraction(1)
		log.Debugf("profile: mutex profiling enabled, %s", fn)
		closer = func() {
			if mp := pprof.Lookup("mutex"); mp != nil {
				_ = mp.WriteTo(f, 0)
			}
			_ = f.Close()
			runtime.SetMutexProfileFraction(0)
			log.Debugf("profile: mutex profiling disabled, %s", fn)
		}
	}
	return
}

func (s *profile) enableBlockProfile(profilePath string) (closer func()) {
	closer = func() {}
	if profilePath != "" {
		fn := filepath.Join(s.path, profilePath)
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create block profile %q: %v", fn, err)
		}
		runtime.SetBlockProfileRate(1)
		log.Debugf("profile: block profiling enabled, %s", fn)
		closer = func() {
			_ = pprof.Lookup("block").WriteTo(f, 0)
			_ = f.Close()
			runtime.SetBlockProfileRate(0)
			log.Debugf("profile: block profiling disabled, %s", fn)
		}
	}
	return
}

func (s *profile) enableThreadCreateProfile(profilePath string) (closer func()) {
	closer = func() {}
	if profilePath != "" {
		fn := filepath.Join(s.path, profilePath)
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create thread creation profile %q: %v", fn, err)
		}
		runtime.SetBlockProfileRate(1)
		log.Debugf("profile: thread creation profiling enabled, %s", fn)
		closer = func() {
			if mp := pprof.Lookup("threadcreate"); mp != nil {
				_ = mp.WriteTo(f, 0)
			}
			_ = f.Close()
			log.Debugf("profile: thread creation profiling disabled, %s", fn)
		}
	}
	return
}

func (s *profile) enableTraceProfile(profilePath string) (closer func()) {
	closer = func() {}
	if profilePath != "" {
		fn := filepath.Join(s.path, profilePath)
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create trace profile %q: %v", fn, err)
		}
		if err := trace.Start(); err != nil {
			log.Fatalf("profile: could not start trace: %v", err)
		}
		if err := traceR.Start(f); err != nil {
			log.Fatalf("profile: could not start trace: %v", err)
		}
		log.Debugf("profile: trace profiling enabled, %s", fn)
		closer = func() {
			trace.Stop()
			traceR.Stop()
			log.Debugf("profile: trace profiling disabled, %s", fn)
		}
	}
	return
}

func (s *profile) enableGoRoutineProfile(profilePath string) (closer func()) {
	closer = func() {}
	if profilePath != "" {
		fn := filepath.Join(s.path, profilePath)
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("profile: could not create goroutine profile %q: %v", fn, err)
		}
		log.Debugf("profile: goroutine profiling enabled, %s", fn)
		closer = func() {
			if mp := pprof.Lookup("goroutine"); mp != nil {
				_ = mp.WriteTo(f, 0)
			}
			_ = f.Close()
			log.Debugf("profile: goroutine profiling disabled, %s", fn)
		}
	}
	return
}

// ProfOpt provides a wrapped functional options type for Start(...)
type ProfOpt func(*profile)

// Start constructs any profiling sessions from `types`.
// `types` is a OR-combined ProfType value, it looks like CPUProf | MemProf.
//
func Start(types ProfType, opts ...ProfOpt) (stop func()) {
	p := &profile{
		path:                 ".",
		cpuProfName:          "cpu.prof",
		memProfName:          "mem.prof",
		mutexProfName:        "mutex.prof",
		blockProfName:        "block.prof",
		threadCreateProfName: "thread-create.prof",
		traceProfName:        "trace.out",
		goRoutineProfName:    "go-routine.prof",
		memProfileRate:       DefaultMemProfileRate,
		memProfileType:       "heap",
	}

	for _, opt := range opts {
		opt(p)
	}

	var closers []func()
	for _, t := range p.types {
		types |= t
	}
	if types&CPUProf == CPUProf {
		c := p.enableCPUProfile(p.cpuProfName)
		closers = append(closers, c)
	}
	if types&MemProf == MemProf {
		c := p.enableMemProfile(p.memProfName)
		closers = append(closers, c)
	}
	if types&MutexProf == MutexProf {
		c := p.enableMutexProfile(p.mutexProfName)
		closers = append(closers, c)
	}
	if types&BlockProf == BlockProf {
		c := p.enableBlockProfile(p.blockProfName)
		closers = append(closers, c)
	}
	if types&ThreadCreateProf == ThreadCreateProf {
		c := p.enableThreadCreateProfile(p.threadCreateProfName)
		closers = append(closers, c)
	}
	if types&TraceProf == TraceProf {
		c := p.enableTraceProfile(p.traceProfName)
		closers = append(closers, c)
	}
	if types&GoRoutineProf == GoRoutineProf {
		c := p.enableGoRoutineProfile(p.goRoutineProfName)
		closers = append(closers, c)
	}

	stop = func() {
		for _, c := range closers {
			c()
		}
	}
	return
}

// WithTypes allows a list of ProfType specified. For Example:
//
//    Start(0, WithTypes(CPUProf, MemProf, BlockProf))
//
func WithTypes(types ...ProfType) ProfOpt {
	return func(profile *profile) {
		profile.types = append(profile.types, types...)
	}
}

// WithOutputDirectory specify a directory for writing profiles. Default is '.'.
func WithOutputDirectory(path string) ProfOpt {
	return func(profile *profile) {
		profile.path = path
	}
}

// WithCPUProfName ..
func WithCPUProfName(filename string) ProfOpt {
	return func(profile *profile) {
		profile.cpuProfName = filename
	}
}

// WithMemProfName ..
func WithMemProfName(filename string) ProfOpt {
	return func(profile *profile) {
		profile.memProfName = filename
	}
}

// WithMutexProfName ..
func WithMutexProfName(filename string) ProfOpt {
	return func(profile *profile) {
		profile.mutexProfName = filename
	}
}

// WithBlockProfName ..
func WithBlockProfName(filename string) ProfOpt {
	return func(profile *profile) {
		profile.blockProfName = filename
	}
}

// WithThreadCreateProfName ..
func WithThreadCreateProfName(filename string) ProfOpt {
	return func(profile *profile) {
		profile.threadCreateProfName = filename
	}
}

// WithTraceProfName ..
func WithTraceProfName(filename string) ProfOpt {
	return func(profile *profile) {
		profile.traceProfName = filename
	}
}

// WithGoRoutineProfName ..
func WithGoRoutineProfName(filename string) ProfOpt {
	return func(profile *profile) {
		profile.goRoutineProfName = filename
	}
}

// WithMemProfileType specify the typename of memory profiling, "heap" or "allocs"
func WithMemProfileType(typ string) ProfOpt {
	return func(profile *profile) {
		profile.memProfileType = typ
	}
}

// WithMemProfileRate enables memory profiling with a special rate
func WithMemProfileRate(rate int) ProfOpt {
	return func(profile *profile) {
		profile.memProfileRate = rate
	}
}
