package basic

import (
	"image/color"

	"github.com/mangofeet/nrdb-go"
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

func getCostContainerStart(ctx *canvas.Context) float64 {
	return getCostContainerRadius(ctx) * 1.3
}

func getTextBoxHeight(ctx *canvas.Context) float64 {
	_, canvasHeight := ctx.Size()
	return (canvasHeight / 3)
}

func drawCostCircle(ctx *canvas.Context, bgColor color.Color) {

	strokeWidth := getStrokeWidth(ctx)
	costContainerR := getCostContainerRadius(ctx)
	costContainerStart := getCostContainerStart(ctx)
	titleBoxHeight := getTitleBoxHeight(ctx)
	titleBoxTop := getTitleBoxTop(ctx)

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	path := canvas.Circle(costContainerR)
	ctx.DrawPath(costContainerStart+(costContainerR), titleBoxTop-(titleBoxHeight*0.5), path)

	ctx.Pop()

}

func drawInflence(ctx *canvas.Context, card *nrdb.Printing, x float64, bgColor color.RGBA) {

	if card.Attributes.InfluenceCost == nil {
		return
	}
	influenceCost := *card.Attributes.InfluenceCost
	strokeWidth := getStrokeWidth(ctx)

	_, canvasHeight := ctx.Size()

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	influenceHeight := getTextBoxHeight(ctx) * 0.6
	influenceWidth := canvasHeight / 42

	// center around the give point
	boxX := x - (influenceWidth / 2)

	ctx.DrawPath(boxX, 0, influenceBox(influenceHeight, influenceWidth, influenceCost))

	ctx.Pop()

	curveRadius := influenceWidth / 2

	pipR := curveRadius * 0.6

	for i := 0.0; i < 5; i += 1 {

		pipY := influenceHeight - (pipR * ((i + 1) * 4)) + (pipR / 2)

		ctx.Push()
		ctx.SetStrokeWidth(strokeWidth * 0.75)
		ctx.SetStrokeColor(textColor)
		ctx.SetFill(transparent)

		pip := canvas.Circle(pipR)
		ctx.DrawPath(x, pipY, pip)

		ctx.Pop()

		if i >= 5-float64(influenceCost) {
			ctx.Push()
			ctx.SetFill(textColor)
			pip := canvas.Circle(pipR * 0.5)
			ctx.DrawPath(x, pipY, pip)
			ctx.Pop()
		}
	}

}

func strength(canvasWidth, canvasHeight float64) *canvas.Path {
	path := &canvas.Path{}

	path.MoveTo(0, canvasHeight/8)
	path.CubeTo(canvasWidth/6, canvasHeight/4, canvasWidth/4, canvasHeight/12, canvasWidth/8, 0)
	path.LineTo(0, 0)
	path.Close()

	return path
}

func influenceBox(height, width float64, pips int) *canvas.Path {

	path := &canvas.Path{}

	curveRadius := width / 2
	curveStart := height - curveRadius

	path.MoveTo(0, 0)
	path.LineTo(0, curveStart)
	path.CubeTo(0, height, width, height, width, curveStart)
	path.LineTo(width, 0)
	path.Close()

	return path
}
