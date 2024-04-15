package prng

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"sync"
)

var globalSampleNum = int64(0)

var sequenceSampleNumbers = map[int]int64{}

var sequenceLock sync.Mutex

func Next(seed string, sampleMax int64) int64 {
	globalSampleNum++
	result := prng(seed, fmt.Sprint(globalSampleNum), sampleMax)
	return result
}

func SequenceNext(sequence int, seed string, sampleMax int64) int64 {
	sequenceLock.Lock()
	defer sequenceLock.Unlock()

	sequenceSampleNumbers[sequence]++
	result := prng(seed, fmt.Sprintf("%d:%d", sequence, sequenceSampleNumbers[sequence]), sampleMax)
	return result
}

func Sample(seed string, sampleNum, sampleMax int64) int64 {
	return prng(seed, fmt.Sprint(sampleNum), sampleMax)
}

func prng(seed, salt string, sampleMax int64) int64 {

	h := sha256.New()
	fmt.Fprintf(h, "%s,%s", seed, salt)
	hash := h.Sum(nil)

	sum := (&big.Int{}).SetBytes(hash)

	mod := sum.Mod(sum, big.NewInt(sampleMax)).Int64()
	return mod + 1

}
