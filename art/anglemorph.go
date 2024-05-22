package art

import (
	"image/color"
	"log"

	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/tdewolff/canvas"
)

type AngleMorphGradient string

const (
	AngleMorphGradientVertical   AngleMorphGradient = "vertical"
	AngleMorphGradientHorizontal AngleMorphGradient = "horizontal"
	AngleMorphGradientBoth       AngleMorphGradient = "both"
	AngleMorphGradientNone       AngleMorphGradient = "none"
)

type AngleMorph struct {
	RNG                               prng.Generator
	Width, Height                     float64
	X, Y                              float64
	MaxShiftFactorX, MaxShiftFactorY  *float64
	ColumnCount, RowCount             int
	InterpolationSteps                *int
	Color                             color.RGBA
	Gradient                          AngleMorphGradient
	ColorShiftMax                     *float64
	StrokeWidthMain, StrokeWidthMinor *float64
	TopRow, BottomRow                 []Point
	LeftColumn, RightColumn           []Point
}

func (drawer *AngleMorph) Draw(ctx *canvas.Context) error {
	columnStep := drawer.Width / float64(drawer.ColumnCount)
	rowStep := drawer.Height / float64(drawer.RowCount)

	cols := make(columns, drawer.ColumnCount+1)

	maxShiftFactorX := 0.5
	if drawer.MaxShiftFactorX != nil {
		maxShiftFactorX = *drawer.MaxShiftFactorX
	}
	maxShiftFactorY := 1.0
	if drawer.MaxShiftFactorY != nil {
		maxShiftFactorY = *drawer.MaxShiftFactorY
	}

	maxShiftX := columnStep * maxShiftFactorX
	maxShiftY := rowStep * maxShiftFactorY

	var thisBottomRow, thisTopRow []Point

	for c := range drawer.ColumnCount + 1 {
		col := make(column, drawer.RowCount+1)

		if c == 0 && len(drawer.LeftColumn) == drawer.RowCount+1 {
			log.Println("using provided left col")
			cols[c] = column(drawer.LeftColumn)
			thisBottomRow = append(thisBottomRow, drawer.LeftColumn[0])
			thisTopRow = append(thisTopRow, drawer.LeftColumn[drawer.RowCount])
			continue
		}
		if c == drawer.ColumnCount && len(drawer.RightColumn) == drawer.RowCount+1 {
			log.Println("using provided right col")
			cols[c] = column(drawer.RightColumn)
			thisBottomRow = append(thisBottomRow, drawer.RightColumn[0])
			thisTopRow = append(thisTopRow, drawer.RightColumn[drawer.RowCount])
			continue
		}

		baseX := columnStep * float64(c)
		for r := range drawer.RowCount + 1 { // go past the top

			if r == 0 && len(drawer.BottomRow) == drawer.ColumnCount+1 {
				col[r] = drawer.BottomRow[c]
				continue
			}
			if r == drawer.RowCount && len(drawer.TopRow) == drawer.ColumnCount+1 {
				col[r] = drawer.TopRow[c]
				continue
			}

			baseY := rowStep * float64(r)

			dX := (float64(drawer.RNG.Next(int64(maxShiftX*200))) - maxShiftX*100) / 100
			dY := float64(drawer.RNG.Next(int64(maxShiftY*100))) / 100

			// go below on first row so it doesn't have a single flat
			// edge on the bottom
			if r == 0 {
				dY *= -1
			}

			x := baseX + dX
			y := baseY + dY

			// add built in X/Y offset for drawing
			col[r] = Point{x + drawer.X, y + drawer.Y}

			if r == 0 {
				thisBottomRow = append(thisBottomRow, col[r])
			}
			if r == drawer.RowCount {
				thisTopRow = append(thisTopRow, col[r])
			}

		}
		cols[c] = col

		if c == 0 {
			drawer.LeftColumn = col
		}
		if c == drawer.ColumnCount {
			drawer.RightColumn = col
		}

	}

	drawer.BottomRow = thisBottomRow
	drawer.TopRow = thisTopRow

	strokeWidthMain := columnStep * 0.02
	if drawer.StrokeWidthMain != nil {
		strokeWidthMain = *drawer.StrokeWidthMain
	}
	strokeWidthMinor := columnStep * 0.02
	if drawer.StrokeWidthMinor != nil {
		strokeWidthMinor = *drawer.StrokeWidthMinor
	}

	if drawer.ColorShiftMax == nil {
		drawer.ColorShiftMax = makePointer(90.0)
	}

	ctx.Push()
	ctx.SetStrokeWidth(strokeWidthMain)
	for i, col := range cols {

		thisColor := drawer.Color
		if drawer.Gradient == AngleMorphGradientBoth || drawer.Gradient == AngleMorphGradientHorizontal {
			colorShift := *drawer.ColorShiftMax * ((float64(i)) / float64(len(cols)/2))
			if i > len(cols)/2 {
				colorShift = *drawer.ColorShiftMax * ((float64(len(cols)) - float64(i)) / float64(len(cols)/2))
			}

			var err error
			thisColor, _, err = Analogous(drawer.Color, colorShift)
			if err != nil {
				panic(err)
			}
		}

		ctx.SetStrokeColor(thisColor)
		drawer.drawCol(ctx, col, thisColor)
		// drawer.drawDots(ctx, col)

		if len(cols) > i+1 {
			ctx.Push()
			ctx.SetStrokeWidth(strokeWidthMinor)
			for _, iCol := range drawer.interpolate(col, cols[i+1]) {
				drawer.drawCol(ctx, iCol, thisColor)
			}
			ctx.Pop()
		}
	}
	ctx.Pop()

	log.Println("done")

	return nil

}

func (drawer AngleMorph) drawCol(ctx *canvas.Context, col column, baseColor color.RGBA) {
	prev := Point{col[0].x, col[0].y}

	for row, p := range col {

		thisColor := baseColor

		if drawer.Gradient == AngleMorphGradientBoth || drawer.Gradient == AngleMorphGradientVertical {
			colorShift := *drawer.ColorShiftMax * ((float64(row)) / float64(len(col)/2))
			if row > len(col)/2 {
				colorShift = *drawer.ColorShiftMax * ((float64(len(col)) - float64(row)) / float64(len(col)/2))
			}
			var err error
			thisColor, _, err = Analogous(baseColor, colorShift)
			if err != nil {
				panic(err)
			}
			ctx.Push()
			ctx.SetStrokeColor(thisColor)
			ctx.MoveTo(prev.x, prev.y)
			ctx.LineTo(p.x, p.y)
			ctx.Stroke()
			ctx.Pop()

			prev = p
		} else { // if not making a gradient this way, we can use single paths
			ctx.Push()
			ctx.SetStrokeColor(thisColor)
			if row == 0 {
				ctx.MoveTo(p.x, p.y)
			}
			ctx.LineTo(p.x, p.y)
			ctx.Pop()
		}

	}

	ctx.Stroke()
}

func (drawer AngleMorph) drawDots(ctx *canvas.Context, col column) {
	ctx.Push()
	ctx.SetStrokeColor(canvas.White)
	ctx.SetFillColor(drawer.Color)
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

			newCol[row] = Point{x, y}
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
		col[i] = Point{x, ys[i]}
	}
	return col
}

type columns []column

type column []Point
