package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/mangofeet/nrdb-go"
)

func drawArt(img draw.Image, card *nrdb.Printing) error {

	seed := card.Attributes.Title + card.Attributes.Text

	startX := prng.Next(seed, canvasWidth)
	startY := prng.Next(seed, canvasHeight)

	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		wlk := Walker{
			Seed: seed,
			X:    startX,
			Y:    startY,
			Color: color.RGBA{
				R: uint8(prng.Next(seed, 256)),
				G: uint8(prng.Next(seed, 256)),
				B: uint8(prng.Next(seed, 256)),
				A: 155,
			},
			Velocity: 1,
		}
		go func(wlk *Walker) {
			wlk.Walk(img)
			wg.Done()
		}(&wlk)
	}

	wg.Wait()

	return nil
}

type Walker struct {
	Seed     string
	X, Y     int64
	Color    color.Color
	Velocity int
}

func (wlk Walker) String() string {
	return fmt.Sprintf("walker at (%d, %d) - %s", wlk.X, wlk.Y, wlk.Color)
}

func (wlk Walker) Draw(img draw.Image) {
	img.Set(int(wlk.X), int(wlk.Y), wlk.Color)
	// img.Set(int(wlk.X+1), int(wlk.Y), wlk.Color)
	// img.Set(int(wlk.X-1), int(wlk.Y), wlk.Color)
	// img.Set(int(wlk.X), int(wlk.Y+1), wlk.Color)
	// img.Set(int(wlk.X), int(wlk.Y-1), wlk.Color)
	// img.Set(int(wlk.X+1), int(wlk.Y+1), wlk.Color)
	// img.Set(int(wlk.X-1), int(wlk.Y-1), wlk.Color)
	// img.Set(int(wlk.X-1), int(wlk.Y+1), wlk.Color)
	// img.Set(int(wlk.X+1), int(wlk.Y-1), wlk.Color)
}

func (wlk Walker) Walk(img draw.Image) {
	for wlk.inBounds(img) {
		direction := prng.Next(wlk.Seed, 4)

		// make sure it moves
		if wlk.Velocity < 1 {
			wlk.Velocity = 1
		}

		switch direction {
		case 1: // left
			wlk.Y -= int64(wlk.Velocity)
		case 2: // down
			wlk.X += int64(wlk.Velocity)
		case 3: // up
			wlk.X -= int64(wlk.Velocity)
		case 4: // right
			wlk.Y += int64(wlk.Velocity)
		}

		wlk.Draw(img)
	}
}

func (wlk Walker) inBounds(img image.Image) bool {
	bounds := img.Bounds()

	if wlk.X > int64(bounds.Max.X) {
		return false
	}
	if wlk.X < int64(bounds.Min.X) {
		return false
	}
	if wlk.Y > int64(bounds.Max.Y) {
		return false
	}
	if wlk.Y < int64(bounds.Min.Y) {
		return false
	}

	return true
}
