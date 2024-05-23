package reflection

import (
	"image/color"
	"math"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/mangofeet/nrdb-go"
	"github.com/ojrac/opensimplex-go"
	"github.com/tdewolff/canvas"
)

type Reflection struct {
	ColumnCount, RowCount int

	InterpolationSteps *int

	Color, ColorBG *color.RGBA
}

func (drawer Reflection) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	seed := card.Attributes.Title + card.Attributes.Text + card.Attributes.CardTypeID + card.Attributes.FactionID + card.Attributes.Flavor

	canvasWidth, canvasHeight := ctx.Size()

	rngGlobal := prng.NewGenerator(seed, nil)

	baseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	if drawer.Color != nil {
		baseColor = *drawer.Color
	}

	cardBGColor := art.Darken(baseColor, 0.623)
	if drawer.ColorBG != nil {
		cardBGColor = *drawer.ColorBG
	}

	splitFactor := (float64(rngGlobal.Next(33)) / 100.0) + 0.33

	width := canvasWidth * 1.2
	height := canvasHeight * (splitFactor + 0.1)
	x := canvasWidth * -0.1
	y := canvasHeight * -0.1

	// columnCount := rngGlobal.Next(60) + 60
	// rowCount := rngGlobal.Next(120)
	columnCount := 60
	rowCount := int(float64(columnCount) * (height / width))

	if drawer.InterpolationSteps == nil {
		drawer.InterpolationSteps = makePointer(10)
	}

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

	bottomColor, _, err := art.Analogous(baseColor, 45)
	if err != nil {
		panic(err)
	}

	first := &art.AngleMorph{
		RNG:                rngGlobal,
		Width:              width,
		Height:             height,
		X:                  x,
		Y:                  y,
		ColumnCount:        int(columnCount),
		RowCount:           int(rowCount),
		InterpolationSteps: drawer.InterpolationSteps,
		Color:              bottomColor,
		Gradient:           art.AngleMorphGradientHorizontal,
		ColorShiftMax:      makePointer(90.0),
		StrokeWidthMain:    makePointer(width * (0.02 / float64(columnCount))),
		StrokeWidthMinor:   makePointer(width * (0.02 / float64(columnCount))),
	}

	first.Draw(ctx)

	y += height
	height = canvasHeight - height + (canvasHeight * 0.1)

	second := &art.AngleMorph{
		RNG:                rngGlobal,
		Width:              width,
		Height:             height,
		X:                  x,
		Y:                  y,
		MaxShiftFactorY:    makePointer(0.2),
		ColumnCount:        int(columnCount),
		RowCount:           int(rowCount),
		InterpolationSteps: drawer.InterpolationSteps,
		Color:              baseColor,
		Gradient:           art.AngleMorphGradientHorizontal,
		ColorShiftMax:      makePointer(90.0),
		StrokeWidthMain:    makePointer(width * (0.03 / float64(columnCount))),
		StrokeWidthMinor:   makePointer(width * (0.03 / float64(columnCount))),
		BottomRow:          first.TopRow,
	}

	second.Draw(ctx)

	var walkers []*art.Walker

	noise := opensimplex.New(rngGlobal.Next(math.MaxInt64))

	for i := 0; i < 5000; i++ {
		colorFactor := rngGlobal.Next(128) - 64

		baseVx := float64(rngGlobal.Next(10) - 5)

		wlk := art.Walker{
			RNG:             rngGlobal,
			Direction:       "down",
			X:               canvasWidth / 2,
			Y:               canvasHeight - 1,
			Vx:              baseVx + (float64(rngGlobal.Next(50)) / 100) - 0.2,
			Vy:              (float64(rngGlobal.Next(50)) / 100) * -1,
			NoiseStepFactor: 0.005,
			NoiseDimensions: 3,
			Noise:           noise,
			StrokeWidth:     width * (0.01 / float64(columnCount)),

			Color: color.RGBA{
				R: uint8(math.Max(0, math.Min(float64(int64(bottomColor.R)+colorFactor), 255))),
				G: uint8(math.Max(0, math.Min(float64(int64(bottomColor.G)+colorFactor), 255))),
				B: uint8(math.Max(0, math.Min(float64(int64(bottomColor.B)+colorFactor), 255))),
				A: 0xff,
			},
		}
		walkers = append(walkers, &wlk)
	}

	for _, wlk := range walkers {
		wlk.Draw(ctx)
		var hasTurned bool
		for wlk.InBounds(ctx) {
			wlk.Velocity()
			wlk.Move()
			wlk.Draw(ctx)

			if !hasTurned && wlk.Y < canvasHeight*splitFactor*1.75 {
				wlk.Vy *= 0.9
				wlk.Vx *= 0.9
			}
			if !hasTurned && wlk.Y < canvasHeight*splitFactor*1.5 {
				wlk.Vy *= 0.9
				wlk.Vx *= 0.9
			}
			if !hasTurned && wlk.Y < canvasHeight*splitFactor*1.25 {
				wlk.Vy *= 0.9
				wlk.Vx *= 0.9
			}
			if !hasTurned && wlk.Y < canvasHeight*splitFactor {
				hasTurned = true
				if wlk.Vx > 0 {
					wlk.Direction = "right"
				} else {
					wlk.Direction = "left"
				}
			}
		}
	}

	// third := &art.AngleMorph{
	// 	RNG:                rngGlobal,
	// 	Width:              canvasWidth * 1.2,
	// 	Height:             canvasHeight * 1.2,
	// 	X:                  canvasWidth * -0.1,
	// 	Y:                  canvasHeight * -0.1,
	// 	ColumnCount:        int(columnCount),
	// 	RowCount:           int(rowCount * 2),
	// 	InterpolationSteps: drawer.InterpolationSteps,
	// 	Color:              cardBGColor,
	// 	Gradient:           art.AngleMorphGradientNone,
	// 	ColorShiftMax:      makePointer(90.0),
	// 	StrokeWidthMain:    makePointer(width * (0.02 / float64(columnCount))),
	// 	StrokeWidthMinor:   makePointer(width * (0.02 / float64(columnCount))),
	// }

	// third.Draw(ctx)

	return nil
}

func makePointer[T any](thing T) *T {
	return &thing
}
