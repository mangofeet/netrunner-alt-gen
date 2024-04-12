package main

import (
	"crypto/sha256"
	"fmt"
	"image"
	"log"
	"math/big"

	"github.com/mangofeet/nrdb-go"
)

func drawArt(img image.Image, card *nrdb.Printing) error {

	seed := card.Attributes.Title + card.Attributes.Text

	for i := 0; i < 10; i++ {
		log.Println(prng(seed, int64(i), 100))
	}

	return nil

}

func prng(seed string, sampleNum, sampleMax int64) int64 {

	h := sha256.New()
	fmt.Fprintf(h, "%s,%d", seed, sampleNum)
	hash := h.Sum(nil)

	sumFloat, _, err := (&big.Float{}).Parse(fmt.Sprintf("%0X", hash), 16)
	if err != nil {
		panic(err)
	}

	sum, _ := sumFloat.Int(nil)

	return sum.Mod(sum, big.NewInt(sampleMax)).Int64() + 1

}
