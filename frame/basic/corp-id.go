package basic

import (
	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func (fb FrameBasic) CorpID() art.Drawer {
	return art.DrawerFunc(func(ctx *canvas.Context, card *nrdb.Printing) error {

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
		ctx.SetFillColor(fb.getColorBG())
		ctx.SetStrokeColor(fb.getColorBorder())
		ctx.SetStrokeWidth(strokeWidth)
		ctx.DrawPath(0, 0, titlePath)
		ctx.DrawPath(0, 0, subtitlePath)
		ctx.Pop()

		boxText, boxType := fb.drawTextBox(ctx, canvasHeight/192, cornerRounded)

		fb.drawCorpLimits(ctx, card, boxText)

		// render card text

		// not sure how these sizes actually correlate to the weird
		// pixel/mm setup I'm using, but these work
		fontSizeTitle := titleBoxHeight * 1.5
		fontSizeSubtitle := fontSizeTitle * subtitleFactor
		fontSizeCard := titleBoxHeight * 1.5

		titleTextX := canvasWidth * 0.25
		titleTextMaxWidth := titleBoxRightIn - titleTextX

		titleText := fb.getTitleText(ctx, card, fontSizeTitle, titleTextMaxWidth, titleBoxHeight, canvas.Left)
		titleTextY := (titleBoxTop - (titleBoxHeight-titleText.Bounds().H)*0.5)
		ctx.DrawText(titleTextX, titleTextY, titleText)
		// canvas.NewTextLine(getFont(fontSizeTitle, canvas.FontRegular), getTitleText(card), canvas.Left))

		subtitleTextX := titleTextX
		subtitleTextMaxWidth := subtitleBoxRightIn - subtitleTextX
		subtitleText := fb.getSubtitleText(ctx, card, fontSizeSubtitle, subtitleTextMaxWidth, subtitleBoxHeight, canvas.Left)

		subtitleTextY := (subtitleBoxTop - (subtitleBoxHeight-subtitleText.Bounds().H)*0.5)
		ctx.DrawText(subtitleTextX, subtitleTextY, subtitleText)

		fb.drawCardText(ctx, card, fontSizeCard, 0, 0, boxText, fb.getAdditionalText()...)
		fb.drawTypeText(ctx, card, fontSizeCard, boxType)

		fb.drawFactionSybmol(ctx, card, canvasWidth*0.15, subtitleBoxBottom+(titleBoxHeight+subtitleBoxHeight)*0.5, (titleBoxHeight + subtitleBoxHeight))

		return nil

	})
}
