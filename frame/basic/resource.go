package basic

import (
	"fmt"
	"image/color"
	"math"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type FrameResource struct{}

func (FrameResource) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	canvasWidth, canvasHeight := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

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

	costContainerR := getCostContainerRadius(ctx)
	costContainerStart := getCostContainerStart(ctx)

	titleBoxHeight := getTitleBoxHeight(ctx)

	titleBoxTop := getTitleBoxTop(ctx)
	titleBoxBottom := titleBoxTop - titleBoxHeight
	titleBoxLeft := costContainerR * 3.25

	titlePath := &canvas.Path{}
	titlePath.MoveTo(titleBoxLeft, titleBoxTop)
	titlePath.QuadTo(titleBoxLeft+(costContainerR*0.5), titleBoxBottom+(titleBoxHeight*0.5), titleBoxLeft, titleBoxBottom)
	titlePath.LineTo(canvasWidth, titleBoxBottom)
	titlePath.LineTo(canvasWidth, titleBoxTop)
	titlePath.Close()

	ctx.DrawPath(0, 0, titlePath)
	ctx.Pop()

	// outline for cost circle
	drawCostCircle(ctx, bgColor)

	boxText, boxType := drawTextBox(ctx, canvasHeight/48, cornerIn)

	drawInfluence(ctx, card, boxText.right, factionColor)

	// render card text

	// not sure how these sizes actually correlate to the weird
	// pixel/mm setup I'm using, but these work
	fontSizeTitle := titleBoxHeight * 2
	fontSizeCost := titleBoxHeight * 3
	fontSizeCard := titleBoxHeight * 1.2

	titleTextX := titleBoxLeft + costContainerR*0.5
	if card.Attributes.IsUnique {
		titleTextX = titleBoxLeft + (costContainerR * 0.4)
	}
	titleTextY := titleBoxTop - titleBoxHeight*0.1
	ctx.DrawText(titleTextX, titleTextY, getCardText(getTitleText(card), fontSizeTitle, canvasWidth, titleBoxHeight, canvas.Left))
	// ctx.DrawText(titleTextX, titleTextY, canvas.NewTextLine(getFont(fontSizeTitle, canvas.FontRegular), getTitleText(card), canvas.Left))

	if card.Attributes.Cost != nil {
		costTextX := costContainerStart
		costTextY := titleBoxBottom + titleBoxHeight/2
		ctx.DrawText(costTextX, costTextY, canvas.NewTextBox(
			getFont(fontSizeCost, canvas.FontBlack), fmt.Sprint(*card.Attributes.Cost),
			costContainerR*2, 0,
			canvas.Center, canvas.Center, 0, 0))
	}

	drawCardText(ctx, card, fontSizeCard, canvasHeight, 0, boxText)
	drawTypeText(ctx, card, fontSizeCard, boxType)

	return nil
}
