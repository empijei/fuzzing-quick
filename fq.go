// Package fq provides helpers to wire `go test -fuzz` together with "testing/quick".
package fq

import (
	"math/rand"
	"reflect"
	"testing/quick"
)

// Value uses the fuzzer input to generate an arbitrary value of type T with
// testing/quick.
func Value[T any](fuzzInput []byte) (t T, ok bool) {
	// this is needed to make sure that rand.Intn lookup will eventually terminate
	// for arbitrarily small numbers.
	// This also helps skewing inputs towards zero values.
	fuzzInput = append(fuzzInput, 0x00, 0x00, 0x00, 0x00)
	rng := rndSrc{buf: fuzzInput}
	rnd := rand.New(&rng)
	rv, ok := quick.Value(reflect.TypeFor[T](), rnd)
	if !ok {
		return t, false
	}
	return rv.Interface().(T), true
}

type rndSrc struct {
	pos uint64
	buf []byte
}

func (r *rndSrc) Int63() int64 { return int64(r.Uint64() % (1 << 63)) }

func (r *rndSrc) Uint64() (ret uint64) {
	if len(r.buf) == 0 {
		return 0
	}
	for i := range 8 {
		ret |= uint64(r.buf[r.pos]) << ((i) * 8)
		r.pos = (r.pos + 1) % uint64(len(r.buf))
	}
	return ret
}

func (r *rndSrc) Seed(seed int64) { /* we're not actually a RNG */ }
