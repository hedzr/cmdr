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

func BenchmarkGetStringR(b *testing.B) {
	copyRootCmd = rootCmdForTesting
	// cmdr.internalResetWorkerNoLock()
	cmdr.ResetOptions()
	cmdr.Set("no-watch-conf-dir", true)
	cmdr.Set("server.deps.kafka.devel.peers", []string{"192.168.0.11", "192.168.0.12", "192.168.0.13"})
	cmdr.Set("server.deps.kafka.devel.id", "default-kafka")
	os.Args = []string{"consul-tags", "--version"}

	if err := cmdr.Exec(rootCmdForTesting, cmdr.WithCustomShowVersion(func() {}), cmdr.WithNoDefaultHelpScreen(true)); err != nil {
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
