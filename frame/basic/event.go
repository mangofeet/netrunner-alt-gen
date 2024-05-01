package basic

import (
	"fmt"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func (fb FrameBasic) Event() art.Drawer {

	return art.DrawerFunc(func(ctx *canvas.Context, card *nrdb.Printing) error {

		canvasWidth, canvasHeight := ctx.Size()

		strokeWidth := getStrokeWidth(ctx)

		factionColor := fb.getColor(card)

		ctx.Push()
		ctx.SetFillColor(bgColor)
		ctx.SetStrokeColor(textColor)
		ctx.SetStrokeWidth(strokeWidth)

		titleBoxHeight := getTitleBoxHeight(ctx)

		titleBoxTop := getTitleBoxTop(ctx)
		titleBoxBottom := titleBoxTop - titleBoxHeight
		titleBoxRight := canvasWidth - (canvasWidth / 16)
		titleBoxRadius := (canvasHeight / 48)
		titleBoxArc1StartX := titleBoxRight - titleBoxRadius
		titleBoxArc1EndY := titleBoxTop - titleBoxRadius
		titleBoxArc2StartY := titleBoxBottom + titleBoxRadius
		titleBoxArc2EndX := titleBoxRight - titleBoxRadius

		costContainerR := getCostContainerRadius(ctx)
		costContainerStart := getCostContainerStart(ctx)

		titlePath := &canvas.Path{}
		titlePath.MoveTo(0, titleBoxTop)

		// background for cost, top
		titlePath.LineTo(costContainerStart, titleBoxTop)
		titlePath.QuadTo(costContainerStart+costContainerR, titleBoxTop+(costContainerR), costContainerStart+(costContainerR*2), titleBoxTop)

		// right side
		titlePath.LineTo(titleBoxArc1StartX, titleBoxTop)
		titlePath.QuadTo(titleBoxRight, titleBoxTop, titleBoxRight, titleBoxArc1EndY)
		titlePath.LineTo(titleBoxRight, titleBoxArc2StartY)
		titlePath.QuadTo(titleBoxRight, titleBoxBottom, titleBoxArc2EndX, titleBoxBottom)

		// background for cost, bottom
		titlePath.LineTo(costContainerStart+(costContainerR*2), titleBoxBottom)
		titlePath.QuadTo(costContainerStart+costContainerR, titleBoxBottom-(costContainerR), costContainerStart, titleBoxBottom)

		// finish title box
		titlePath.LineTo(0, titleBoxBottom)
		titlePath.Close()

		ctx.DrawPath(0, 0, titlePath)
		ctx.Pop()

		drawCostCircle(ctx, transparent)

		boxText, boxType := drawTextBox(ctx, canvasHeight/48, cornerRounded)

		drawInfluenceAndOrFactionSymbol(ctx, card, boxText.right, factionColor)

		// render card text

		// not sure how these sizes actually correlate to the weird
		// pixel/mm setup I'm using, but these work
		fontSizeTitle := titleBoxHeight * 2
		fontSizeCost := titleBoxHeight * 3
		fontSizeCard := titleBoxHeight * 1.2

		titleTextX := costContainerStart + (costContainerR * 2) + (costContainerR / 3)
		titleTextY := titleBoxTop - titleBoxHeight*0.1
		ctx.DrawText(titleTextX, titleTextY, getCardText(getTitle(card), fontSizeTitle, titleBoxRight, titleBoxHeight, canvas.Left))
		// ctx.DrawText(titleTextX, titleTextY, canvas.NewTextLine(getFont(fontSizeTitle, canvas.FontRegular), getTitleText(card), canvas.Left))

		if card.Attributes.Cost != nil {
			costTextX := costContainerStart
			costTextY := titleBoxBottom + titleBoxHeight/2
			ctx.DrawText(costTextX, costTextY, canvas.NewTextBox(
				getFont(fontSizeCost, canvas.FontBlack), fmt.Sprint(*card.Attributes.Cost),
				costContainerR*2, 0,
				canvas.Center, canvas.Center, 0, 0))
		}

		drawCardText(ctx, card, fontSizeCard, canvasHeight, 0, boxText, fb.getAdditionalText()...)
		drawTypeText(ctx, card, fontSizeCard, boxType)

		return nil
	})
}
