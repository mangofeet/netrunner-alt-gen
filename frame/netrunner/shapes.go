package netrunner

import (
	"image/color"

	"github.com/tdewolff/canvas"
)

func getStrokeWidth(ctx *canvas.Context) float64 {
	_, canvasHeight := ctx.Size()
	return canvasHeight * 0.0023
}

func getTitleBoxHeight(ctx *canvas.Context) float64 {
	_, canvasHeight := ctx.Size()
	return (canvasHeight / 16)
}

func getTitleBoxTop(ctx *canvas.Context) float64 {
	_, canvasHeight := ctx.Size()
	return canvasHeight - (canvasHeight / 12)
}

func getCostContainerRadius(ctx *canvas.Context) float64 {
	return getTitleBoxHeight(ctx) * 0.667
}

func drawCostCircle(ctx *canvas.Context, bgColor color.Color) {

	strokeWidth := getStrokeWidth(ctx)
	costContainerR := getCostContainerRadius(ctx)
	titleBoxHeight := getTitleBoxHeight(ctx)
	titleBoxTop := getTitleBoxTop(ctx)

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	path := canvas.Circle(costContainerR)
	ctx.DrawPath(costContainerR+(costContainerR), titleBoxTop-(titleBoxHeight*0.5), path)

	ctx.Pop()

}

func strength(canvasWidth, canvasHeight float64) *canvas.Path {
	path := &canvas.Path{}

	path.MoveTo(0, canvasHeight/8)
	path.CubeTo(canvasWidth/6, canvasHeight/4, canvasWidth/4, canvasHeight/12, canvasWidth/8, 0)
	path.LineTo(0, 0)
	path.Close()

	return path
}

func influence(height, width float64, pips int) *canvas.Path {

	path := &canvas.Path{}

	curveRadius := width / 2
	curveStart := height - curveRadius

	path.MoveTo(0, 0)
	path.LineTo(0, curveStart)
	path.CubeTo(0, height, width, height, width, curveStart)
	path.LineTo(width, 0)
	path.Close()

	pipR := curveRadius * 0.6
	pipX := width - ((width - (pipR * 2)) / 2)

	for i := 0.0; i < 5; i += 1 {

		pipY := height - (pipR * ((i + 1) * 4)) + (pipR / 2)

		path.MoveTo(pipX, pipY)

		if i >= 5-float64(pips) {
			path.Arc(pipR, pipR, 0, 0, 360)
		} else {
			path.Arc(pipR, pipR, 0, 360, 0)
		}
	}

	return path

}
