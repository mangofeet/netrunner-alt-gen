package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"strings"
	"sync"

	"github.com/StephaneBunel/bresenham"
	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/mangofeet/nrdb-go"
	"github.com/ojrac/opensimplex-go"
)

// number of walkers to draw
const numWalkers = 1000

// the default factor for corrdinates in the noise map
const noiseStepFactor = 0.005

// max total steps to try to escape the canvas, walker will be discarded
// if it does not make it out by this count
const maxSteps = 100_000_000

// the point at which the noise factor will change if the walkers are
// not escaping the canvas
const noiseStepChangeThreshold = 10_000

func drawArt(img draw.Image, card *nrdb.Printing) error {

	seed := card.Attributes.Title + card.Attributes.Text + card.Attributes.CardTypeID + card.Attributes.FactionID

	startX := prng.Next(seed, canvasWidth/2) + canvasWidth/4
	startY := prng.Next(seed, canvasHeight/4) + (canvasHeight / 6)

	if card.Attributes.CardTypeID == "ice" {
		startY = prng.Next(seed, canvasHeight/4) + ((canvasHeight / 3) * 2)
	}
	// startX := img.Bounds().Max.X / 2
	// startY := img.Bounds().Max.Y / 2

	baseColor := getFactionBaseColor(card.Attributes.FactionID)

	noise := opensimplex.New(prng.Next(seed, math.MaxInt64))

	var walkers []*Walker
	for i := 0; i < numWalkers; i++ {

		colorFactor := prng.Next(seed, 128) - 64

		var direction string
		var grid = false
		dirSeed := prng.Next(seed, 41)

		if dirSeed <= 10 {
			direction = "up"
		} else if dirSeed <= 20 {
			direction = "down"
		} else if dirSeed <= 30 {
			direction = "left"
		} else if dirSeed <= 40 {
			direction = "right"
		} else {
			grid = true
		}

		wlk := Walker{
			Seed:      seed,
			Sequence:  i,
			Direction: direction,
			X:         float64(startX),
			Y:         float64(startY),
			Vx:        (float64(prng.Next(seed, 200)) / 100) - 1,
			Vy:        (float64(prng.Next(seed, 200)) / 100) - 1,
			Color: color.RGBA{
				R: uint8(math.Max(0, math.Min(float64(int64(baseColor.R)+colorFactor), 255))),
				G: uint8(math.Max(0, math.Min(float64(int64(baseColor.G)+colorFactor), 255))),
				B: uint8(math.Max(0, math.Min(float64(int64(baseColor.B)+colorFactor), 255))),
				A: 0xff,
			},
			Noise: noise,
			Grid:  grid,
			// Grid:  prng.Next(seed, 10) <= 3,
			// Grid: true,
		}
		walkers = append(walkers, &wlk)
	}

	var walkerPaths []image.Image

	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, wlk := range walkers {
		wg.Add(1)
		go func(wlk *Walker) {
			newPath := wlk.Walk(img.Bounds())
			if wlk.stepCount < maxSteps {
				lock.Lock()
				walkerPaths = append(walkerPaths, newPath)
				lock.Unlock()
				log.Printf("done %s", wlk)
			} else {
				log.Printf("ignoring incomplete: %s", wlk)
			}
			wg.Done()
		}(wlk)
	}
	wg.Wait()

	log.Println("overlaying paths")
	for _, path := range walkerPaths {
		draw.Draw(img, img.Bounds(), path, image.Point{}, draw.Over)
	}
	log.Println("done overlaying paths")

	return nil
}

type point struct {
	x, y float64
}

type Walker struct {
	Seed         string
	Sequence     int
	Direction    string
	X, Y, Vx, Vy float64
	Color        color.Color
	Noise        opensimplex.Noise
	stepCount    int
	prev         *point
	Grid         bool
}

func (wlk Walker) String() string {
	return fmt.Sprintf("walker %d at (%f, %f), step %d", wlk.Sequence, wlk.X, wlk.Y, wlk.stepCount)
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

	// bresenham.DrawLine(img, int(x1+1), int(y1+1), int(x2+1), int(y2+1), wlk.Color)
	// bresenham.DrawLine(img, int(x1), int(y1+1), int(x2), int(y2+1), wlk.Color)
	// bresenham.DrawLine(img, int(x1+1), int(y1), int(x2+1), int(y2), wlk.Color)

	// bresenham.DrawLine(img, int(x1-1), int(y1-1), int(x2-1), int(y2-1), wlk.Color)
	// bresenham.DrawLine(img, int(x1), int(y1-1), int(x2), int(y2-1), wlk.Color)
	// bresenham.DrawLine(img, int(x1-1), int(y1), int(x2-1), int(y2), wlk.Color)
}

func (wlk Walker) drawPoint(img draw.Image, x, y float64) {

	img.Set(int(x), int(y), wlk.Color)

	// const size = 5

	// for xp, yp := x-size, y-size; xp < x+size && yp < y+size; xp, yp = xp+1, yp+1 {
	// 	img.Set(int(xp), int(yp), wlk.Color)
	// }

	// for xp, yp := x-size, y+size; xp < x+size && yp > y-size; xp, yp = xp+1, yp-1 {
	// 	img.Set(int(xp), int(yp), wlk.Color)

	// }

	// for xp := x - size; xp < x+size; xp = xp + 1 {
	// 	img.Set(int(xp), int(y), wlk.Color)
	// }

	// for yp := y - size; yp < y+size; yp = yp + 1 {
	// 	img.Set(int(x), int(yp), wlk.Color)
	// }

}

func (wlk *Walker) Walk(bounds image.Rectangle) image.Image {

	img := image.NewRGBA(bounds)

	for wlk.inBounds(bounds) && wlk.stepCount < maxSteps {

		wlk.stepCount++

		thisNoiseStepFactor := float64(int(wlk.stepCount/1_000)+1) * noiseStepFactor

		deltaX := wlk.Noise.Eval2(wlk.X*thisNoiseStepFactor, wlk.Y*thisNoiseStepFactor)
		deltaY := wlk.Noise.Eval2(wlk.Y*thisNoiseStepFactor, wlk.X*thisNoiseStepFactor)

		switch strings.ToLower(wlk.Direction) {
		case "up":
			wlk.Vx += deltaX
			wlk.Vy += -1 * math.Abs(deltaY)
		case "down":
			wlk.Vx += deltaX
			wlk.Vy += math.Abs(deltaY)
		case "left":
			wlk.Vx += -1 * math.Abs(deltaX)
			wlk.Vy += deltaY
		case "right":
			wlk.Vx += math.Abs(deltaX)
			wlk.Vy += deltaY
		default:
			wlk.Vx += deltaX
			wlk.Vy += deltaY
		}

		if wlk.Grid {

			switch prng.SequenceNext(wlk.Sequence, wlk.Seed, 4) {
			case 1:
				wlk.X += wlk.Vx
			case 2:
				wlk.Y += wlk.Vy
			case 3:
				wlk.Y -= wlk.Vy
			case 4:
				wlk.X -= wlk.Vx
			}

		} else {
			wlk.X += wlk.Vx
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

func getFactionBaseColor(factionID string) color.RGBA {

	switch factionID {
	case "shaper", "weyland":
		return color.RGBA{
			R: 0x7f,
			G: 0x9f,
			B: 0x7f,
			A: 0xff,
		}

	case "anarch":
		return color.RGBA{
			R: 0xdf,
			G: 0xaf,
			B: 0x8f,
			A: 0xff,
		}

	case "criminal":
		return color.RGBA{
			R: 0x8c,
			G: 0xd0,
			B: 0xd3,
			A: 0xff,
		}

	case "nbn":
		return color.RGBA{
			R: 0xf0,
			G: 0xdf,
			B: 0xaf,
			A: 0xff,
		}

	case "jinteki":
		return color.RGBA{
			R: 0xcc,
			G: 0x93,
			B: 0x96,
			A: 0xff,
		}

	case "haas_bioroid":
		return color.RGBA{
			R: 0xc0,
			G: 0xbe,
			B: 0xd1,
			A: 0xff,
		}

	}

	return color.RGBA{
		R: 0xee,
		G: 0xee,
		B: 0xee,
		A: 0xff,
	}

}
