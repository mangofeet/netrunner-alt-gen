package prng

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

var globalSampleNum = int64(0)

func Next(seed string, sampleMax int64) int64 {
	globalSampleNum++
	result := prng(seed, globalSampleNum, sampleMax)
	return result
}

func Sample(seed string, sampleNum, sampleMax int64) int64 {
	return prng(seed, sampleNum, sampleMax)
}

func prng(seed string, sampleNum, sampleMax int64) int64 {

	h := sha256.New()
	fmt.Fprintf(h, "%s,%d", seed, sampleNum)
	hash := h.Sum(nil)

	sum := (&big.Int{}).SetBytes(hash)

	mod := sum.Mod(sum, big.NewInt(sampleMax)).Int64()
	return mod + 1

}
