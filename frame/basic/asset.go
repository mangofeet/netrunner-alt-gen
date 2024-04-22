package basic

import (
	"image/color"
	"log"
	"math"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type FrameAsset struct{}

func (FrameAsset) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	canvasWidth, canvasHeight := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	log.Printf("strokeWidth: %f", strokeWidth)

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
	ctx.SetStrokeWidth(strokeWidth)

	titleBoxHeight := getTitleBoxHeight(ctx)
	fontSizeCost := titleBoxHeight * 2.3
	boxResIcon, err := drawRezCost(ctx, card, fontSizeCost)
	if err != nil {
		return err
	}

	titleBoxTop := getTitleBoxTop(ctx)
	titleBoxBottom := titleBoxTop - titleBoxHeight
	titleBoxLeftOut := boxResIcon.left + (boxResIcon.width * 1.2)
	titleBoxLeftIn := titleBoxLeftOut + titleBoxHeight*0.8

	titlePath := &canvas.Path{}
	titlePath.MoveTo(titleBoxLeftOut, titleBoxTop)
	titlePath.LineTo(canvasWidth, titleBoxTop)
	titlePath.LineTo(canvasWidth, titleBoxBottom)
	titlePath.LineTo(titleBoxLeftIn, titleBoxBottom)
	titlePath.Close()

	ctx.DrawPath(0, 0, titlePath)
	ctx.Pop()

	var boxText, boxType textBoxDimensions
	if card.Attributes.TrashCost != nil {
		boxText, boxType = drawTextBoxTrashable(ctx, canvasHeight/192, cornerRounded)
	} else {
		boxText, boxType = drawTextBox(ctx, canvasHeight/192, cornerRounded)
	}

	drawInfluence(ctx, card, boxText.left, factionColor)

	if _, err := drawTrashCost(ctx, card); err != nil {
		return err
	}
	// render card text

	// not sure how these sizes actually correlate to the weird
	// pixel/mm setup I'm using, but these work
	fontSizeTitle := titleBoxHeight * 2
	fontSizeCard := titleBoxHeight * 1.2

	titleTextX := titleBoxLeftIn
	if card.Attributes.IsUnique { // unique diamon fits better in the angled end here
		titleTextX = titleBoxLeftIn - titleBoxHeight*0.2
	}

	titleTextMaxWidth := (canvasWidth * 0.9) - titleBoxLeftIn

	titleText := getTitleText(ctx, card, fontSizeTitle, titleTextMaxWidth, titleBoxHeight)
	titleTextY := (titleBoxTop - (titleBoxHeight-titleText.Bounds().H)*0.5)
	ctx.DrawText(titleTextX, titleTextY, titleText)
	// canvas.NewTextLine(getFont(fontSizeTitle, canvas.FontRegular), getTitleText(card), canvas.Left))

	drawCardText(ctx, card, fontSizeCard, canvasHeight, 0, boxText)
	drawTypeText(ctx, card, fontSizeCard, boxType)

	return nil

}
