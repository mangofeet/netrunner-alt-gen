package netspace

import (
	"image/color"
	"log"
	"math"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/mangofeet/nrdb-go"
	"github.com/ojrac/opensimplex-go"
	"github.com/tdewolff/canvas"
)

// number of walkers to draw
const numWalkersMin = 1000
const numWalkersMax = 10000

// const numWalkersMin = 100
// const numWalkersMax = 100

func Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	seed := card.Attributes.Title + card.Attributes.Text + card.Attributes.CardTypeID + card.Attributes.FactionID

	canvasWidth, canvasHeight := ctx.Size()

	numWalkers := int(math.Max(float64(numWalkersMin), float64(prng.Next(seed, int64(numWalkersMax)))))

	startX := prng.Next(seed, int64(canvasWidth/2)) + int64(canvasWidth/4)
	startY := prng.Next(seed, int64(canvasHeight/6)) + (int64(canvasHeight/8) * 5)

	if card.Attributes.CardTypeID == "ice" {
		startY = prng.Next(seed, int64(canvasHeight/4)) + (int64(canvasHeight / 6))
	}

	baseColor := art.GetFactionBaseColor(card.Attributes.FactionID)

	cardBGColor := art.Darken(baseColor, 0.623)

	// fill background
	ctx.Push()
	ctx.SetFillColor(cardBGColor)
	ctx.MoveTo(0, 0)
	ctx.LineTo(0, canvasHeight)
	ctx.LineTo(canvasWidth, canvasHeight)
	ctx.LineTo(canvasWidth, 0)
	ctx.Close()
	ctx.Fill()
	ctx.Pop()

	noise := opensimplex.New(prng.Next(seed, math.MaxInt64))

	var walkers []*art.Walker

	nGrid := math.Max(float64(numWalkers)*0.01, float64(prng.Next(seed, int64(float64(numWalkers)*0.02))))

	dirChangeStep := 30.0

	for i := 0; i < numWalkers; i++ {

		colorFactor := prng.Next(seed, 128) - 64

		var direction string
		var grid = false
		var strokeWidth = 0.3

		if float64(i) < nGrid {
			colorFactor = -2 * int64(math.Abs(float64(colorFactor)))
			grid = true
			strokeWidth = 1.5
		} else {

			dirSeed := prng.Next(seed, 4)

			if dirSeed <= 1 {
				direction = "up"
			} else if dirSeed <= 2 {
				direction = "down"
			} else if dirSeed <= 3 {
				direction = "left"
			} else if dirSeed <= 4 {
				direction = "right"
			}
		}

		wlk := art.Walker{
			Seed:                seed,
			Sequence:            i,
			Direction:           direction,
			DirectionVariance:   prng.Next(seed, 4),
			DirectionChangeStep: dirChangeStep,
			X:                   float64(startX),
			Y:                   float64(startY),
			Vx:                  (float64(prng.Next(seed, 100)) / 100) - 0.5,
			Vy:                  (float64(prng.Next(seed, 100)) / 100) - 0.5,
			Color: color.RGBA{
				R: uint8(math.Max(0, math.Min(float64(int64(baseColor.R)+colorFactor), 255))),
				G: uint8(math.Max(0, math.Min(float64(int64(baseColor.G)+colorFactor), 255))),
				B: uint8(math.Max(0, math.Min(float64(int64(baseColor.B)+colorFactor), 255))),
				A: 0xff,
			},
			Noise:       noise,
			Grid:        grid,
			StrokeWidth: strokeWidth,
		}
		walkers = append(walkers, &wlk)
	}

	for _, wlk := range walkers {
		wlk.Draw(ctx)
		for wlk.InBounds(ctx) {
			wlk.Velocity()
			wlk.Move()
			wlk.Draw(ctx)
		}
		log.Printf("finished %s", wlk)
	}

	return nil
}
