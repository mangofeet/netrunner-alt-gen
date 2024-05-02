package netspace

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/mangofeet/nrdb-go"
	"github.com/ojrac/opensimplex-go"
	"github.com/tdewolff/canvas"
)

type Netspace struct {
	MinWalkers, MaxWalkers                 int
	Color, AltColor1, AltColor2, AltColor3 *color.RGBA
	GridColor1, GridColor2, GridColor3     *color.RGBA
}

func (drawer Netspace) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	seed := card.Attributes.Title + card.Attributes.Text + card.Attributes.CardTypeID + card.Attributes.FactionID

	canvasWidth, canvasHeight := ctx.Size()

	numWalkers := int(math.Max(float64(drawer.MinWalkers), float64(prng.Next(seed, int64(drawer.MaxWalkers)))))

	startX := prng.Next(seed, int64(canvasWidth/2)) + int64(canvasWidth/4)
	startY := prng.Next(seed, int64(canvasHeight/6)) + (int64(canvasHeight/8) * 5)

	if card.Attributes.CardTypeID == "ice" {
		startY = prng.Next(seed, int64(canvasHeight/4)) + (int64(canvasHeight / 6))
	}

	baseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	if drawer.Color != nil {
		baseColor = *drawer.Color
	}
	altColor1, altColor2, err := art.Analogous(baseColor, 20)
	if err != nil {
		return fmt.Errorf("getting analogous colors: %w", err)
	}
	if drawer.AltColor1 != nil {
		altColor1 = *drawer.AltColor1
	}
	if drawer.AltColor2 != nil {
		altColor2 = *drawer.AltColor2
	}

	altColor3, _, err := art.Analogous(baseColor, 40)
	if err != nil {
		return fmt.Errorf("getting third analog: %w", err)
	}
	altColor3, err = art.Desaturate(altColor3, 0.5)
	if err != nil {
		return fmt.Errorf("getting desaturated compliment: %w", err)
	}
	if drawer.AltColor3 != nil {
		altColor3 = *drawer.AltColor3
	}

	gridColor1, gridColor2, gridColor3 := altColor1, altColor2, altColor3
	if drawer.GridColor1 != nil {
		gridColor1 = *drawer.GridColor1
	}
	if drawer.GridColor2 != nil {
		gridColor2 = *drawer.GridColor2
	}
	if drawer.GridColor3 != nil {
		gridColor3 = *drawer.GridColor3
	}

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

	nGrid := 0.0
	if float64(numWalkers)*0.01 >= 1 {
		nGrid = math.Max(float64(numWalkers)*0.01, float64(prng.Next(seed, int64(float64(numWalkers)*0.02))))
	}

	dirChangeStep := 30.0

	// do manual seeds for these with high numbers so they didn't
	// affect the walkers
	var altColorDirection1 string
	// switch prng.Sample(seed, 9999999, 4) {
	switch prng.Next(seed, 4) {
	case 1:
		altColorDirection1 = "up"
	case 2:
		altColorDirection1 = "down"
	case 3:
		altColorDirection1 = "left"
	case 4:
		altColorDirection1 = "right"
	}

	altColorDirection2 := altColorDirection1
	for altColorDirection2 == altColorDirection1 {
		// switch prng.Sample(seed, i+10000000, 4) {
		switch prng.Next(seed, 4) {
		case 1:
			altColorDirection2 = "up"
		case 2:
			altColorDirection2 = "down"
		case 3:
			altColorDirection2 = "left"
		case 4:
			altColorDirection2 = "right"
		}
	}

	altColorDirection3 := altColorDirection1
	for altColorDirection3 == altColorDirection1 || altColorDirection3 == altColorDirection2 {
		// switch prng.Sample(seed, i+10000000, 4) {
		switch prng.Next(seed, 4) {
		case 1:
			altColorDirection3 = "up"
		case 2:
			altColorDirection3 = "down"
		case 3:
			altColorDirection3 = "left"
		case 4:
			altColorDirection3 = "right"
		}
	}

	for i := 0; i < numWalkers; i++ {

		colorFactor := prng.Next(seed, 128) - 64

		var direction string
		var grid = false
		var strokeWidth = 0.3

		thisColor := baseColor

		if float64(i) < nGrid {
			colorFactor = -2 * int64(math.Abs(float64(colorFactor)))
			grid = true
			strokeWidth = 1.5
			switch prng.Next(seed, 6) {
			case 1:
				thisColor = gridColor1
			case 2:
				thisColor = gridColor2
			case 3:
				thisColor = gridColor3
			}
			thisColor, err = art.Desaturate(thisColor, float64(colorFactor)/-128.0)
			if err != nil {
				return err
			}
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
			switch direction {
			case altColorDirection1:
				thisColor = altColor1
			case altColorDirection2:
				thisColor = altColor2
			case altColorDirection3:
				thisColor = altColor3
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
				R: uint8(math.Max(0, math.Min(float64(int64(thisColor.R)+colorFactor), 255))),
				G: uint8(math.Max(0, math.Min(float64(int64(thisColor.G)+colorFactor), 255))),
				B: uint8(math.Max(0, math.Min(float64(int64(thisColor.B)+colorFactor), 255))),
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
	}
	log.Printf("finished %d walkers", len(walkers))

	return nil
}
