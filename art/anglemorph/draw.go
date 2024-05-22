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

	columnStep := canvasWidth / float64(drawer.ColumnCount)
	rowStep := canvasHeight / float64(drawer.RowCount)

	cols := make(columns, drawer.ColumnCount+1)

	maxShiftX := columnStep * 0.5
	maxShiftY := rowStep * 1.0

	baseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	if drawer.Color != nil {
		baseColor = *drawer.Color
	}

	cardBGColor := art.Darken(baseColor, 0.623)
	if drawer.ColorBG != nil {
		cardBGColor = *drawer.ColorBG
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

	for c := range drawer.ColumnCount + 1 {
		col := make(column, drawer.RowCount+1)
		baseX := columnStep * float64(c)
		for r := range drawer.RowCount + 1 { // go past the top
			baseY := rowStep * float64(r)

			dX := (float64(rngGlobal.Next(int64(maxShiftX*200))) - maxShiftX*100) / 100
			dY := float64(rngGlobal.Next(int64(maxShiftY*100))) / 100

			if r == 0 {
				dY = 0
			}

			x := baseX + dX
			y := baseY + dY

			col[r] = point{x, y}

		}
		cols[c] = col
	}

	ctx.Push()
	ctx.SetStrokeWidth(columnStep * 0.02)
	for i, col := range cols {
		// colorShift := 90 * ((float64(i)) / float64(len(cols)/2))
		// if i > len(cols)/2 {
		// 	colorShift = 90 * ((float64(len(cols)) - float64(i)) / float64(len(cols)/2))
		// }

		// thisColor, _, err := art.Analogous(baseColor, colorShift)
		// if err != nil {
		// 	panic(err)
		// }

		thisColor := baseColor

		drawCol(ctx, col, thisColor)

		if len(cols) > i+1 {
			for _, iCol := range drawer.interpolate(col, cols[i+1]) {
				drawCol(ctx, iCol, thisColor)
			}
		}
	}
	ctx.Pop()

	return nil
}

func drawCol(ctx *canvas.Context, col column, baseColor color.RGBA) {
	prev := point{col[0].x, col[0].y}
	for row, p := range col {
		colorShift := 90 * ((float64(row)) / float64(len(col)/2))
		if row > len(col)/2 {
			colorShift = 90 * ((float64(len(col)) - float64(row)) / float64(len(col)/2))
		}

		ctx.Push()
		thisColor, _, err := art.Analogous(baseColor, colorShift)
		if err != nil {
			panic(err)
		}
		ctx.SetStrokeColor(thisColor)
		ctx.MoveTo(prev.x, prev.y)
		ctx.LineTo(p.x, p.y)
		prev = p
		ctx.Stroke()
		ctx.Pop()
	}
}

func drawDots(ctx *canvas.Context, col column) {
	ctx.Push()
	ctx.SetStrokeColor(canvas.Black)
	for _, p := range col {
		dot := canvas.Circle(ctx.StrokeWidth * 3)
		ctx.DrawPath(p.x, p.y, dot)
	}
	ctx.Pop()
}

func (drawer AngleMorph) interpolate(col1, col2 column) columns {
	if len(col1) != len(col2) {
		panic("length mismatch in interpolate")
	}

	steps := 10
	if drawer.InterpolationSteps != nil {
		steps = *drawer.InterpolationSteps
	}

	col1x, col1y := unzip(col1)
	col2x, col2y := unzip(col2)

	var res columns

	for step := range steps {
		interpolationAmt := float64(step) / float64(steps)

		newCol := make(column, len(col1))

		for row := range len(col1) {
			x := lerp(col1x[row], col2x[row], interpolationAmt)
			y := lerp(col1y[row], col2y[row], interpolationAmt)

			newCol[row] = point{x, y}
		}

		res = append(res, newCol)
	}

	return res
}

func lerp(v1, v2, amt float64) float64 {
	if amt == 0 {
		return v1
	}

	diff := v2 - v1
	return v1 + (diff * amt)

}

func unzip(col column) ([]float64, []float64) {
	x := make([]float64, len(col))
	y := make([]float64, len(col))
	for i, p := range col {
		x[i] = p.x
		y[i] = p.y
	}

	return x, y
}

func zip(xs, ys []float64) column {
	if len(xs) != len(ys) {
		panic("length mismatch in zip")
	}

	col := make(column, len(xs))

	for i, x := range xs {
		col[i] = point{x, ys[i]}
	}
	return col
}

type columns []column

type column []point

type point struct {
	x, y float64
}
