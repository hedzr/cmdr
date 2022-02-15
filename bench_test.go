// Copyright © 2019 Hedzr Yeh.

package cmdr_test

import (
	"github.com/hedzr/cmdr"
	"os"
	"strconv"
	"testing"
)

func BenchmarkItoa(b *testing.B) {
	num := 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.Itoa(num)
	}
}

func benchmarkGetStringR(lock bool, b *testing.B) {
	if lock {
		cmdr.InternalResetWorkerForTest()
	} else {
		cmdr.InternalResetWorkerNoLockForTest()
	}
	// resetFlagsAndLog(t)
	resetOsArgs()
	cmdr.ResetOptions()

	//copyRootCmd := rootCmdForTesting
	cmdr.Set("no-watch-conf-dir", true)
	cmdr.Set("server.deps.kafka.devel.peers", []string{"192.168.0.11", "192.168.0.12", "192.168.0.13"})
	cmdr.Set("server.deps.kafka.devel.id", "default-kafka")
	os.Args = []string{"consul-tags", "--version"}

	if err := cmdr.Exec(rootCmdForTesting(), cmdr.WithCustomShowVersion(func() {}), cmdr.WithNoDefaultHelpScreen(true)); err != nil {
		b.Fatal(err)
	}

	if cmdr.GetStringR("server.deps.kafka.devel.id") != "default-kafka" {
		b.Fatal("cmdr core logic failed: expect 'default-kafka'")
	}

	var str string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str = cmdr.GetStringR("server.deps.kafka.devel.id")
	}

	b.Log("got id: ", str)
}

func prepare(lock bool, b *testing.B) {
	if lock {
		cmdr.InternalResetWorkerForTest()
	} else {
		cmdr.InternalResetWorkerNoLockForTest()
	}
	// resetFlagsAndLog(t)
	resetOsArgs()
	cmdr.ResetOptions()

	//copyRootCmd = rootCmdForTesting
	cmdr.Set("no-watch-conf-dir", true)
	cmdr.Set("server.deps.kafka.devel.peers", []string{"192.168.0.11", "192.168.0.12", "192.168.0.13"})
	cmdr.Set("server.deps.kafka.devel.id", "default-kafka")
	os.Args = []string{"consul-tags", "--version"}

	if err := cmdr.Exec(rootCmdForTesting(), cmdr.WithCustomShowVersion(func() {}), cmdr.WithNoDefaultHelpScreen(true)); err != nil {
		b.Fatal(err)
	}

	if cmdr.GetStringR("server.deps.kafka.devel.id") != "default-kafka" {
		b.Fatal("cmdr core logic failed: expect 'default-kafka'")
	}

	b.ResetTimer()
}

func benchmarkGetStringR2(lock bool, pb *testing.PB) {
	cmdr.GetStringR("server.deps.kafka.devel.id")
}

func BenchmarkGetStringR(b *testing.B) {
	prepare(false, b)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkGetStringR2(false, pb)
		}
	})
}

func BenchmarkGetStringRNL(b *testing.B) {
	prepare(true, b)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkGetStringR2(true, pb)
		}
	})
}

func BenchmarkGetStringRLock(b *testing.B) {
	benchmarkGetStringR(false, b)
}

func BenchmarkGetStringRNoLock(b *testing.B) {
	benchmarkGetStringR(true, b)
}

/*

$ go test -bench '^BenchmarkGet.*$' -run '^$'
goos: darwin
goarch: amd64
pkg: github.com/hedzr/cmdr
BenchmarkGetStringR-4         	 9215539	       138 ns/op
BenchmarkGetStringRNL-4       	 7936184	       149 ns/op
BenchmarkGetStringRNoLock-4   	 5068267	       215 ns/op
--- BENCH: BenchmarkGetStringRLock-4
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
BenchmarkGetStringRLock-4     	 5475811	       220 ns/op
--- BENCH: BenchmarkGetStringRNoLock-4
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
PASS
ok  	github.com/hedzr/cmdr	7.607s

######################### So, it is not matter about the `internalGetWorker()` with r-lock on uniqueWorker.

$ go test -bench '^BenchmarkGet.*$' -run '^$' -benchtime 20s
goos: darwin
goarch: amd64
pkg: github.com/hedzr/cmdr
BenchmarkGetStringR-4         	175795060	       134 ns/op
BenchmarkGetStringRNL-4       	162173757	       150 ns/op
BenchmarkGetStringRNoLock-4   	100000000	       219 ns/op
--- BENCH: BenchmarkGetStringRLock-4
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
BenchmarkGetStringRLock-4     	100000000	       220 ns/op
--- BENCH: BenchmarkGetStringRNoLock-4
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
    b_test.go:51: got id:  default-kafka
PASS
ok  	github.com/hedzr/cmdr	123.291s
*/
