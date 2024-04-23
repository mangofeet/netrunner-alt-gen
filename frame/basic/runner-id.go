package basic

import (
	"log"

	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type FrameRunnerID struct{}

func (FrameRunnerID) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	canvasWidth, canvasHeight := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	log.Printf("strokeWidth: %f", strokeWidth)

	titleBoxHeight := getTitleBoxHeight(ctx)

	// title box
	titleBoxTop := getTitleBoxTop(ctx)
	titleBoxBottom := titleBoxTop - titleBoxHeight
	titleBoxLeft := canvasWidth * 0.2
	// titleBoxLeft := costIconX
	titleBoxRadius := (canvasHeight / 48)
	titleBoxArc1StartY := titleBoxTop - titleBoxRadius
	titleBoxArc1EndX := titleBoxLeft + titleBoxRadius
	titleBoxArc2StartX := titleBoxLeft + titleBoxRadius
	titleBoxArc2EndY := titleBoxBottom + titleBoxRadius

	titlePath := &canvas.Path{}
	titlePath.MoveTo(titleBoxLeft, titleBoxArc1StartY)
	titlePath.QuadTo(titleBoxLeft, titleBoxTop, titleBoxArc1EndX, titleBoxTop)

	titlePath.LineTo(canvasWidth, titleBoxTop)
	titlePath.LineTo(canvasWidth, titleBoxBottom)
	titlePath.LineTo(titleBoxArc2StartX, titleBoxBottom)
	titlePath.QuadTo(titleBoxLeft, titleBoxBottom, titleBoxLeft, titleBoxArc2EndY)

	titlePath.Close()

	// subtitle box
	subtitleFactor := 0.6
	subtitleBoxHeight := titleBoxHeight * subtitleFactor
	subtitleBoxTop := titleBoxBottom
	subtitleBoxLeft := titleBoxLeft + canvasWidth*0.1
	subtitleBoxBottom := subtitleBoxTop - subtitleBoxHeight
	subtitleBoxArc2StartX := subtitleBoxLeft + titleBoxRadius
	subtitleBoxArc2EndY := subtitleBoxBottom + titleBoxRadius

	subtitlePath := &canvas.Path{}
	subtitlePath.MoveTo(subtitleBoxLeft, subtitleBoxTop)
	subtitlePath.LineTo(canvasWidth, subtitleBoxTop)
	subtitlePath.LineTo(canvasWidth, subtitleBoxBottom)
	subtitlePath.LineTo(subtitleBoxArc2StartX, subtitleBoxBottom)
	subtitlePath.QuadTo(subtitleBoxLeft, subtitleBoxBottom, subtitleBoxLeft, subtitleBoxArc2EndY)

	subtitlePath.Close()

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	ctx.DrawPath(0, 0, titlePath)
	ctx.DrawPath(0, 0, subtitlePath)
	ctx.Pop()

	boxText, boxType := drawTextBox(ctx, canvasHeight/48, cornerRounded)

	drawMU(ctx, card, false)
	drawLink(ctx, card)

	// render card text

	// not sure how these sizes actually correlate to the weird
	// pixel/mm setup I'm using, but these work
	fontSizeTitle := titleBoxHeight * 2
	fontSizeSubtitle := fontSizeTitle * subtitleFactor
	fontSizeCard := titleBoxHeight * 1.2

	titleTextX := titleBoxLeft + titleBoxHeight*0.2
	titleTextY := titleBoxTop - titleBoxHeight*0.1
	ctx.DrawText(titleTextX, titleTextY, getCardText(getTitle(card), fontSizeTitle, canvasWidth-titleBoxLeft, titleBoxHeight, canvas.Left))

	subtitleTextX := subtitleBoxLeft + subtitleBoxHeight*0.2
	subtitleTextY := subtitleBoxTop - subtitleBoxHeight*0.1
	ctx.DrawText(subtitleTextX, subtitleTextY, getCardText(getSubtitle(card), fontSizeSubtitle, canvasWidth-subtitleBoxLeft, subtitleBoxHeight, canvas.Left))

	drawCardText(ctx, card, fontSizeCard, boxText.height*0.45, canvasWidth*0.06, boxText)
	drawTypeText(ctx, card, fontSizeCard, boxType)

	return nil
}
