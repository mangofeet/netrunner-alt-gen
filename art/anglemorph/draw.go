package anglemorph

import (
	"image/color"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type AngleMorph struct {
	ColumnCount, RowCount int

	InterpolationSteps *int

	Color, ColorBG *color.RGBA
}

func (drawer AngleMorph) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

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

	// width := canvasWidth * 1.2
	// height := canvasHeight * 1.2
	// x := canvasWidth * -0.1
	// y := canvasHeight * -0.1
	width := canvasWidth * 1.2
	height := canvasHeight * 0.6
	x := canvasWidth * -0.1
	y := canvasHeight * -0.1

	// columnCount := rngGlobal.Next(60) + 60
	// rowCount := rngGlobal.Next(120)
	columnCount := 40
	rowCount := int(float64(columnCount) * (height / width))
	rowCount *= 2

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

	first := &art.AngleMorph{
		RNG:                rngGlobal,
		Width:              width,
		Height:             height,
		X:                  x,
		Y:                  y,
		ColumnCount:        int(columnCount),
		RowCount:           int(rowCount),
		InterpolationSteps: drawer.InterpolationSteps,
		Color:              art.Complementary(baseColor),
		Gradient:           art.AngleMorphGradientHorizontal,
		ColorShiftMax:      makePointer(90.0),
		StrokeWidthMain:    makePointer(width * (0.02 / float64(columnCount))),
		StrokeWidthMinor:   makePointer(width * (0.02 / float64(columnCount))),
	}

	first.Draw(ctx)

	y += height

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
