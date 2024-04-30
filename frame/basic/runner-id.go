package basic

import (
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type FrameRunnerID struct{}

func (FrameRunnerID) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	canvasWidth, canvasHeight := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	titleBoxHeight := getTitleBoxHeight(ctx)

	// title box
	titleBoxTop := getTitleBoxTop(ctx)
	titleBoxBottom := titleBoxTop - titleBoxHeight
	titleBoxLeftTop := canvasWidth * 0.15
	titleBoxLeftBottom := canvasWidth * 0.2
	// titleBoxLeft := costIconX
	titleBoxRadius := (canvasHeight / 48)
	titleBoxDangleTopRight := titleBoxBottom - titleBoxHeight*0.4
	titleBoxDangleTopLeft := titleBoxBottom + titleBoxHeight*0.1
	titleBoxDangleBottom := titleBoxBottom - titleBoxHeight*1.3
	titleBoxDangleLeft := canvasWidth * 0.0853
	titleBoxDangleRight := titleBoxDangleLeft + titleBoxHeight*1.07
	// titleBoxArc1StartY := titleBoxTop - titleBoxRadius
	// titleBoxArc1EndX := titleBoxLeft + titleBoxRadius
	// titleBoxArc2StartX := titleBoxLeft + titleBoxRadius
	// titleBoxArc2EndY := titleBoxBottom + titleBoxRadius

	titlePath := &canvas.Path{}
	titlePath.MoveTo(titleBoxLeftTop, titleBoxTop)
	titlePath.LineTo(canvasWidth, titleBoxTop)
	titlePath.LineTo(canvasWidth, titleBoxBottom)
	titlePath.LineTo(titleBoxLeftBottom, titleBoxBottom)

	titlePath.LineTo(titleBoxDangleRight, titleBoxDangleTopRight)
	titlePath.LineTo(titleBoxDangleRight, titleBoxDangleBottom)
	titlePath.LineTo(titleBoxDangleLeft, titleBoxDangleBottom)
	titlePath.LineTo(titleBoxDangleLeft, titleBoxDangleTopLeft)
	titlePath.LineTo(titleBoxLeftTop, titleBoxTop)

	titlePath.Close()

	// subtitle box
	subtitleFactor := 0.6
	subtitleBoxHeight := titleBoxHeight * subtitleFactor
	subtitleBoxTop := titleBoxBottom
	subtitleBoxLeft := titleBoxLeftBottom + canvasWidth*0.1
	subtitleBoxBottom := subtitleBoxTop - subtitleBoxHeight
	subtitleBoxArc2StartX := subtitleBoxLeft + titleBoxRadius
	// subtitleBoxArc2EndY := subtitleBoxBottom + titleBoxRadius

	subtitlePath := &canvas.Path{}
	subtitlePath.MoveTo(subtitleBoxLeft, subtitleBoxTop)
	subtitlePath.LineTo(canvasWidth, subtitleBoxTop)
	subtitlePath.LineTo(canvasWidth, subtitleBoxBottom)
	subtitlePath.LineTo(subtitleBoxArc2StartX, subtitleBoxBottom)
	// subtitlePath.QuadTo(subtitleBoxLeft, subtitleBoxBottom, subtitleBoxLeft, subtitleBoxArc2EndY)
	subtitlePath.LineTo(subtitleBoxLeft, subtitleBoxTop)

	subtitlePath.Close()

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	ctx.DrawPath(0, 0, titlePath)
	ctx.DrawPath(0, 0, subtitlePath)
	ctx.Pop()

	boxText, boxType := drawTextBox(ctx, canvasHeight/96, cornerStraight)

	drawRunnerLimits(ctx, card, boxText)
	drawMU(ctx, card, false)
	drawLink(ctx, card)

	// render card text

	// not sure how these sizes actually correlate to the weird
	// pixel/mm setup I'm using, but these work
	fontSizeTitle := titleBoxHeight * 2
	fontSizeSubtitle := fontSizeTitle * subtitleFactor
	fontSizeCard := titleBoxHeight * 1.2

	factionBubbleWidth := (titleBoxHeight + subtitleBoxHeight)
	factionBubbleX := canvasWidth * 0.85
	factionBubbleY := subtitleBoxBottom + (titleBoxHeight+subtitleBoxHeight)*0.5

	titleTextX := titleBoxLeftBottom + titleBoxHeight*0.3
	titleTextMaxWidth := factionBubbleX - titleTextX - factionBubbleWidth*0.7

	titleText := getTitleText(ctx, card, fontSizeTitle, titleTextMaxWidth, titleBoxHeight, canvas.Left)
	titleTextY := titleBoxTop - titleBoxHeight*0.1
	ctx.DrawText(titleTextX, titleTextY, titleText)

	subtitleTextX := subtitleBoxLeft + subtitleBoxHeight*0.6
	subtitleTextMaxWidth := factionBubbleX - subtitleTextX - factionBubbleWidth*0.7

	subtitleText := getSubtitleText(ctx, card, fontSizeSubtitle, subtitleTextMaxWidth, subtitleBoxHeight, canvas.Left)
	subtitleTextY := (subtitleBoxTop - (subtitleBoxHeight-subtitleText.Bounds().H)*0.5)
	ctx.DrawText(subtitleTextX, subtitleTextY, subtitleText)

	drawCardText(ctx, card, fontSizeCard, boxText.height*0.75, canvasWidth*0.06, boxText)
	drawTypeText(ctx, card, fontSizeCard, boxType)
	drawFactionSybmol(ctx, card, factionBubbleX, factionBubbleY, factionBubbleWidth)

	return nil
}
