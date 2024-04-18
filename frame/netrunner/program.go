package netrunner

import (
	"image/color"
	"math"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func DrawFrameProgram(ctx *canvas.Context, card *nrdb.Printing) error {

	canvasWidth, canvasHeight := ctx.Size()

	factionBaseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	factionColor := color.RGBA{
		R: uint8(math.Max(0, math.Min(float64(int64(factionBaseColor.R)-48), 255))),
		G: uint8(math.Max(0, math.Min(float64(int64(factionBaseColor.G)-48), 255))),
		B: uint8(math.Max(0, math.Min(float64(int64(factionBaseColor.B)-48), 255))),
		A: 0xff,
	}

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(10)

	titleBoxHeight := (canvasHeight / 16)

	titleBoxTop := canvasHeight - (canvasHeight / 12)
	titleBoxBottom := titleBoxTop - titleBoxHeight
	titleBoxArcStart := canvasWidth - (canvasWidth / 3)
	titleBoxRight := canvasWidth - (canvasWidth / 16)
	titleBoxArcCP1 := titleBoxRight - (canvasWidth / 48)

	costContainerR := titleBoxHeight * 0.667
	costContainerStart := costContainerR

	titlePath := &canvas.Path{}
	titlePath.MoveTo(0, titleBoxTop)

	// background for cost, top
	titlePath.LineTo(costContainerStart, titleBoxTop)
	titlePath.QuadTo(costContainerStart+costContainerR, titleBoxTop+(costContainerR*0.8), costContainerStart+(costContainerR*2), titleBoxTop)

	// arc down on right side
	titlePath.LineTo(titleBoxArcStart, titleBoxTop)
	titlePath.QuadTo(titleBoxArcCP1, titleBoxTop, titleBoxRight, titleBoxBottom)

	// background for cost, bottom
	titlePath.LineTo(costContainerStart+(costContainerR*2), titleBoxBottom)
	titlePath.QuadTo(costContainerStart+costContainerR, titleBoxBottom-(costContainerR*0.8), costContainerStart, titleBoxBottom)

	// finish title box
	titlePath.LineTo(0, titleBoxBottom)
	titlePath.Close()

	ctx.DrawPath(0, 0, titlePath)
	ctx.Pop()

	// outline for cost circle

	ctx.Push()
	ctx.SetFillColor(transparent)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(10)
	costOutline := canvas.Circle(costContainerR)
	ctx.DrawPath(costContainerStart+(costContainerR), titleBoxTop-(titleBoxHeight*0.5), costOutline)
	ctx.Pop()

	// bottom text box

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(10)

	textBoxHeight := canvasHeight / 3
	textBoxLeft := canvasWidth / 8
	textBoxRight := canvasWidth - (canvasWidth / 12)
	textBoxArcRadius := (canvasHeight / 32)
	textBoxArc1StartY := textBoxHeight - textBoxArcRadius
	textBoxArc1EndX := textBoxLeft + textBoxArcRadius

	textBoxArc2StartX := textBoxRight - textBoxArcRadius
	textBoxArc2EndY := textBoxHeight - textBoxArcRadius

	ctx.MoveTo(textBoxLeft, 0)
	ctx.LineTo(textBoxLeft, textBoxArc1StartY)
	ctx.QuadTo(textBoxLeft, textBoxHeight, textBoxArc1EndX, textBoxHeight)

	ctx.LineTo(textBoxArc2StartX, textBoxHeight)
	ctx.QuadTo(textBoxRight, textBoxHeight, textBoxRight, textBoxArc2EndY)

	ctx.LineTo(textBoxRight, 0)

	ctx.FillStroke()
	ctx.Pop()

	ctx.Push()
	ctx.SetFillColor(factionColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(10)

	influenceHeight := textBoxHeight * 0.55
	influenceWidth := canvasHeight / 48

	influenceCost := 0
	if card.Attributes.InfluenceCost != nil {
		influenceCost = *card.Attributes.InfluenceCost
	}
	ctx.DrawPath(textBoxRight-(influenceWidth/2), 0, influence(influenceHeight, influenceWidth, influenceCost))

	ctx.Pop()

	// program strength
	ctx.Push()

	ctx.SetFillColor(factionColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(10)

	ctx.MoveTo(0, canvasHeight/8)
	ctx.CubeTo(canvasWidth/6, canvasHeight/4, canvasWidth/4, canvasHeight/12, canvasWidth/8, 0)
	ctx.LineTo(0, 0)
	ctx.Close()
	ctx.FillStroke()

	ctx.Pop()

	return nil
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
