package prng

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

type Generator interface {
	Next(max int64) int64
	Sample(sample, max int64) int64
	Sequence() int64
}

func NewGenerator(seed string, sequence *int64) Generator {
	return &generator{
		seed:     seed,
		sequence: sequence,
		sample:   0,
	}
}

type generator struct {
	seed     string
	sequence *int64
	sample   int64
}

func (gen *generator) Next(max int64) int64 {
	gen.sample++
	return gen.Sample(gen.sample, max)
}

func (gen *generator) Sample(sample, max int64) int64 {
	return prng(gen.seed, gen.getSalt(sample), max)
}

func (gen generator) Sequence() int64 {
	if gen.sequence == nil {
		return 0
	}
	return *gen.sequence
}

func (gen *generator) getSalt(sample int64) string {
	if gen.sequence == nil {
		return fmt.Sprint(sample)
	}
	return fmt.Sprintf("%d:%d", *gen.sequence, sample)
}

func prng(seed, salt string, sampleMax int64) int64 {

	h := sha256.New()
	fmt.Fprintf(h, "%s,%s", seed, salt)
	hash := h.Sum(nil)

	sum := (&big.Int{}).SetBytes(hash)

	mod := sum.Mod(sum, big.NewInt(sampleMax)).Int64()
	return mod + 1

}
