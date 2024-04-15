package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"sync"

	"github.com/StephaneBunel/bresenham"
	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/mangofeet/nrdb-go"
	"github.com/ojrac/opensimplex-go"
)

const numWalkers = 15

func drawArt(img draw.Image, card *nrdb.Printing) error {

	seed := card.Attributes.Title + card.Attributes.Text

	// startX := prng.Next(seed, canvasWidth)
	// startY := prng.Next(seed, canvasHeight)
	startX := img.Bounds().Max.X / 2
	startY := img.Bounds().Max.Y / 2

	var walkers []*Walker
	for i := 0; i < numWalkers; i++ {

		colorFactor := uint8(prng.Next(seed, 16))

		wlk := Walker{
			Seed:     seed,
			Sequence: i,
			X:        float64(startX),
			Y:        float64(startY),
			Vx:       1,
			Vy:       1,
			Color: color.RGBA{
				R: 128 + colorFactor,
				G: 128 + colorFactor,
				B: 128 + colorFactor,
				A: 128,
			},
			Noise: opensimplex.New(prng.Next(seed, math.MaxInt64)),
		}
		walkers = append(walkers, &wlk)
	}

	var walkerPaths []image.Image

	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, wlk := range walkers {
		wg.Add(1)
		go func(wlk *Walker) {
			log.Printf("starting walker %d", wlk.Sequence)
			newPath := wlk.Walk(img.Bounds())
			lock.Lock()
			walkerPaths = append(walkerPaths, newPath)
			lock.Unlock()
			log.Printf("done walker %d", wlk.Sequence)
			wg.Done()
		}(wlk)
	}
	wg.Wait()

	for _, path := range walkerPaths {
		draw.Draw(img, img.Bounds(), path, image.Point{}, draw.Over)
	}

	return nil
}

type point struct {
	x, y float64
}

type Walker struct {
	Seed         string
	Sequence     int
	X, Y, Vx, Vy float64
	Color        color.Color
	Noise        opensimplex.Noise
	stepCount    int
	prev         *point
}

func (wlk Walker) String() string {
	return fmt.Sprintf("walker %d at (%f, %f) - %s", wlk.Sequence, wlk.X, wlk.Y, wlk.Color)
}

func (wlk *Walker) Draw(img draw.Image) {

	if wlk.prev == nil {
		wlk.drawPoint(img, wlk.X, wlk.Y)
		goto SET_PREV
	}

	wlk.drawLine(img, wlk.X, wlk.Y, wlk.prev.x, wlk.prev.y)

SET_PREV:

	wlk.prev = &point{wlk.X, wlk.Y}
}

func (wlk Walker) drawLine(img draw.Image, x1, y1, x2, y2 float64) {
	bresenham.DrawLine(img, int(x1), int(y1), int(x2), int(y2), wlk.Color)
}

func (wlk Walker) drawPoint(img draw.Image, x, y float64) {

	const size = 5

	for xp, yp := x-size, y-size; xp < x+size && yp < y+size; xp, yp = xp+1, yp+1 {
		img.Set(int(xp), int(yp), wlk.Color)
	}

	for xp, yp := x-size, y+size; xp < x+size && yp > y-size; xp, yp = xp+1, yp-1 {
		img.Set(int(xp), int(yp), wlk.Color)

	}

	for xp := x - size; xp < x+size; xp = xp + 1 {
		img.Set(int(xp), int(y), wlk.Color)
	}

	for yp := y - size; yp < y+size; yp = yp + 1 {
		img.Set(int(x), int(yp), wlk.Color)
	}

	// img.Set(int(x), int(y), wlk.Color)
	// img.Set(int(x+1), int(y), wlk.Color)
	// img.Set(int(x+2), int(y), wlk.Color)
	// img.Set(int(x-1), int(y), wlk.Color)
	// img.Set(int(x-2), int(y), wlk.Color)
	// img.Set(int(x), int(y+1), wlk.Color)
	// img.Set(int(x), int(y+2), wlk.Color)
	// img.Set(int(x), int(y-1), wlk.Color)
	// img.Set(int(x), int(y-2), wlk.Color)
	// img.Set(int(x+1), int(y+1), wlk.Color)
	// img.Set(int(x+2), int(y+2), wlk.Color)
	// img.Set(int(x-1), int(y-1), wlk.Color)
	// img.Set(int(x-2), int(y-2), wlk.Color)
	// img.Set(int(x-1), int(y+1), wlk.Color)
	// img.Set(int(x-2), int(y+2), wlk.Color)
	// img.Set(int(x+1), int(y-1), wlk.Color)
	// img.Set(int(x+2), int(y-2), wlk.Color)
}

func (wlk *Walker) Walk(bounds image.Rectangle) image.Image {

	img := image.NewRGBA(bounds)

	for wlk.inBounds(bounds) {

		wlk.stepCount++

		wlk.Vx += wlk.Noise.Eval2(wlk.X, float64(wlk.stepCount))
		wlk.Vy += wlk.Noise.Eval2(float64(wlk.stepCount), wlk.Y)
		// wlk.Vx += wlk.Noise.Eval2(wlk.X, wlk.Y)
		// wlk.Vy += wlk.Noise.Eval2(wlk.X, wlk.Y)

		direction := prng.SequenceNext(wlk.Sequence, wlk.Seed, 4)

		switch direction {
		case 1: // left
			wlk.Y -= wlk.Vy
		case 2: // down
			wlk.X += wlk.Vx
		case 3: // up
			wlk.X -= wlk.Vx
		case 4: // right
			wlk.Y += wlk.Vy
		}

		wlk.Draw(img)
	}

	return img
}

func (wlk Walker) inBounds(bounds image.Rectangle) bool {
	if int(wlk.X) > bounds.Max.X {
		return false
	}
	if int(wlk.X) < bounds.Min.X {
		return false
	}
	if int(wlk.Y) > bounds.Max.Y {
		return false
	}
	if int(wlk.Y) < bounds.Min.Y {
		return false
	}

	return true
}
