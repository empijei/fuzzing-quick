package fq

import (
	"math/rand"
	"testing"
)

func TestRandRead(t *testing.T) {
	rng := rndSrc{buf: []byte("this is_ the so_urce")}
	rnd := rand.New(&rng)
	const want = "this is the source"
	var gotb [len(want)]byte
	rnd.Read(gotb[:])
	got := string(gotb[:])
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

type genMe struct {
	Field string
}

func TestCanGenerateI31(t *testing.T) {
	_, ok := Value[genMe]([]byte("\xFF"))
	if !ok {
		t.Fatalf("unexpected !ok")
	}
}

func FuzzGen(f *testing.F) {
	f.Fuzz(func(t *testing.T, buf []byte) {
		_, ok := Value[genMe](buf)
		if !ok {
			t.Fatalf("unexpected !ok for input %q", buf)
		}
	})
}
