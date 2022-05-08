package randomizer_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/hedzr/cmdr/tool/randomizer"
)

func TestNew(t *testing.T) {
	r := randomizer.New()
	t.Log(r.Next())
	t.Log(r.NextIn(100))
	t.Log(r.NextInRange(20, 30))
}

func BenchmarkRandomizer(b *testing.B) {
	var result int
	r := randomizer.New()
	for n := 0; n < b.N; n++ {
		result = r.NextIn(9139)
	}
	b.Logf("end of: %v", result)
}

func BenchmarkRandomizerHiRes(b *testing.B) {
	var result uint64
	r := randomizer.New().(randomizer.HiresRandomizer) //nolint:errcheck //like it
	for n := 0; n < b.N; n++ {
		result = r.HiresNextIn(9139)
	}
	b.Logf("end of: %v", result)
}

func BenchmarkGlobal(b *testing.B) {
	var result int
	for n := 0; n < b.N; n++ {
		result = rand.Intn(9139) //nolint:gosec //like it
	}
	b.Logf("end of: %v", result)
}

func BenchmarkNative(b *testing.B) {
	var result int
	random := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec //like it
	for n := 0; n < b.N; n++ {
		result = random.Intn(9139)
	}
	b.Logf("end of: %v", result)
}
