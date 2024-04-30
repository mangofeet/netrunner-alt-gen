package basic

import (
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type FrameCorpID struct{}

func (FrameCorpID) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	canvasWidth, canvasHeight := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	titleBoxHeight := getTitleBoxHeight(ctx)
	titleBoxTop := getTitleBoxTop(ctx)
	titleBoxBottom := titleBoxTop - titleBoxHeight
	titleBoxRightOut := canvasWidth * 0.9
	titleBoxRightIn := titleBoxRightOut - titleBoxHeight
	titlePath := &canvas.Path{}
	titlePath.MoveTo(0, titleBoxTop)
	titlePath.LineTo(titleBoxRightOut, titleBoxTop)
	titlePath.LineTo(titleBoxRightIn, titleBoxBottom)
	titlePath.LineTo(0, titleBoxBottom)
	titlePath.Close()

	subtitleFactor := 0.6
	subtitleBoxHeight := titleBoxHeight * subtitleFactor
	subtitleBoxTop := titleBoxBottom
	subtitleBoxBottom := subtitleBoxTop - subtitleBoxHeight
	subtitleBoxRightOut := canvasWidth * 0.7
	subtitleBoxRightIn := subtitleBoxRightOut - subtitleBoxHeight
	subtitlePath := &canvas.Path{}
	subtitlePath.MoveTo(0, subtitleBoxTop)
	subtitlePath.LineTo(subtitleBoxRightOut, subtitleBoxTop)
	subtitlePath.LineTo(subtitleBoxRightIn, subtitleBoxBottom)
	subtitlePath.LineTo(0, subtitleBoxBottom)
	subtitlePath.Close()

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)
	ctx.DrawPath(0, 0, titlePath)
	ctx.DrawPath(0, 0, subtitlePath)
	ctx.Pop()

	boxText, boxType := drawTextBox(ctx, canvasHeight/192, cornerRounded)

	drawCorpLimits(ctx, card, boxText)

	// render card text

	// not sure how these sizes actually correlate to the weird
	// pixel/mm setup I'm using, but these work
	fontSizeTitle := titleBoxHeight * 2
	fontSizeSubtitle := fontSizeTitle * subtitleFactor
	fontSizeCard := titleBoxHeight * 1.2

	titleTextX := canvasWidth * 0.25
	titleTextMaxWidth := titleBoxRightIn - titleTextX

	titleText := getTitleText(ctx, card, fontSizeTitle, titleTextMaxWidth, titleBoxHeight, canvas.Left)
	titleTextY := (titleBoxTop - (titleBoxHeight-titleText.Bounds().H)*0.5)
	ctx.DrawText(titleTextX, titleTextY, titleText)
	// canvas.NewTextLine(getFont(fontSizeTitle, canvas.FontRegular), getTitleText(card), canvas.Left))

	subtitleTextX := titleTextX
	subtitleTextY := subtitleBoxTop - subtitleBoxHeight*0.1
	ctx.DrawText(subtitleTextX, subtitleTextY, getCardText(getSubtitle(card), fontSizeSubtitle, subtitleBoxRightIn, subtitleBoxHeight, canvas.Left))

	drawCardText(ctx, card, fontSizeCard, canvasHeight, 0, boxText)
	drawTypeText(ctx, card, fontSizeCard, boxType)

	drawFactionSybmol(ctx, card, canvasWidth*0.15, subtitleBoxBottom+(titleBoxHeight+subtitleBoxHeight)*0.5, (titleBoxHeight + subtitleBoxHeight))

	return nil

}
