package randomizer

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand"
	"time"
)

// New return a tool for randomizer
func New() Randomizer { return &randomizer{} }

// Randomizer enables normal resolution randomizer
type Randomizer interface {
	Next() int
	NextIn(max int) int
	NextInRange(min, max int) int
	AsHires() HiresRandomizer
	AsStrings() StringsRandomizer
}

// HiresRandomizer enables high resolution randomizer
type HiresRandomizer interface {
	HiresNext() uint64
	HiresNextIn(max uint64) uint64
	HiresNextInRange(min, max uint64) uint64
}

// StringsRandomizer interface
type StringsRandomizer interface {
	// NextStringSimple returns a random string with specified length 'n', just in A..Z
	NextStringSimple(n int) string
	// NextString returns a random string with specified length 'n'
	NextString(n int) string

	NextStringByCharset(n int, charset []rune) string

	NextStringWithVariantLength() string
	NextStringWithVariantLengthRange(min, max int) string

	NextStringWithVariantLengthByCharset(min, max int, charset []rune) string
}

type randomizer struct {
	lastErr error
}

// var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
const (
	// Alphabets gets the a to z and A to Z
	Alphabets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Digits gets 0 to 9
	Digits = "0123456789"
	// AlphabetNumerics gets Alphabets and Digits
	AlphabetNumerics = Alphabets + Digits
	// Symbols gets the ascii symbols
	Symbols = "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
	// ASCII gets the ascii characters
	ASCII = AlphabetNumerics + Symbols
)

var hundred = big.NewInt(100)
var seededRand = mrand.New(mrand.NewSource(time.Now().UTC().UnixNano()))

func (r *randomizer) Next() int                    { return seededRand.Int() }
func (r *randomizer) NextIn(max int) int           { return seededRand.Intn(max) }
func (r *randomizer) inRange(min, max int) int     { return seededRand.Intn(max-min) + min }
func (r *randomizer) NextInRange(min, max int) int { return r.inRange(min, max) }
func (r *randomizer) AsHires() HiresRandomizer     { return r }
func (r *randomizer) AsStrings() StringsRandomizer { return r }

func (r *randomizer) HiresNext() uint64 {
	return r.hiresNextIn(hundred)
}

func (r *randomizer) HiresNextIn(max uint64) uint64 {
	return r.hiresNextIn(big.NewInt(int64(max)))
}

func (r *randomizer) hiresNextIn(max *big.Int) uint64 {
	var bi *big.Int
	bi, r.lastErr = rand.Int(rand.Reader, max)
	if r.lastErr == nil {
		return bi.Uint64()
	}
	return 0
}

func (r *randomizer) hiresInRange(min, max uint64) uint64 {
	var bi *big.Int
	bi, r.lastErr = rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if r.lastErr == nil {
		return bi.Uint64() + min
	}
	return 0
}

func (r *randomizer) HiresNextInRange(min, max uint64) uint64 { return r.hiresInRange(min, max) }

//
//
//

// NextStringSimple returns a random string with specified length 'n', just in A..Z
func (r *randomizer) NextStringSimple(n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte(r.inRange(65, 90)) // 'a' .. 'z'
	}
	return string(bytes)
}

// NextString returns a random string with specified length 'n'
func (r *randomizer) NextString(n int) string {
	return r.randStringBaseImpl(n, []rune(AlphabetNumerics))
}

func (r *randomizer) randStringBaseImpl(n int, charset []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (r *randomizer) NextStringByCharset(n int, charset []rune) string {
	return r.randStringBaseImpl(n, charset)
}

// NextStringWithVariantLength returns a random string with random length (1..127)
func (r *randomizer) NextStringWithVariantLength() string {
	n := seededRand.Intn(128)
	return r.NextString(n)
}

func (r *randomizer) NextStringWithVariantLengthRange(min, max int) string {
	return r.NextStringWithVariantLengthByCharset(min, max, []rune(AlphabetNumerics))
}

func (r *randomizer) NextStringWithVariantLengthByCharset(min, max int, charset []rune) string {
	if min <= 0 {
		min = 1
	}
	if max <= min+1 {
		max = min + 4096
	}
	length := r.NextInRange(min, max)
	return r.randStringBaseImpl(length, charset)
}
