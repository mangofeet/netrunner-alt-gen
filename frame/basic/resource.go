package basic

import (
	"fmt"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func (fb FrameBasic) Resource() art.Drawer {
	return art.DrawerFunc(func(ctx *canvas.Context, card *nrdb.Printing) error {

		canvasWidth, canvasHeight := ctx.Size()

		strokeWidth := getStrokeWidth(ctx)

		ctx.Push()
		ctx.SetFillColor(fb.getColorBG())
		ctx.SetStrokeColor(fb.getColorBorder())
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
		fb.drawCostCircle(ctx, fb.getColorBG())

		boxText, boxType := fb.drawTextBox(ctx, canvasHeight/48, cornerIn)

		fb.drawInfluenceAndOrFactionSymbol(ctx, card, boxText.right)

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
		ctx.DrawText(titleTextX, titleTextY, fb.getCardText(getTitle(card), fontSizeTitle, canvasWidth, titleBoxHeight, canvas.Left))
		// ctx.DrawText(titleTextX, titleTextY, canvas.NewTextLine(getFont(fontSizeTitle, canvas.FontRegular), getTitleText(card), canvas.Left))

		if card.Attributes.Cost != nil {
			costTextX := costContainerStart
			costTextY := titleBoxBottom + titleBoxHeight/2
			ctx.DrawText(costTextX, costTextY, canvas.NewTextBox(
				fb.getFont(fontSizeCost, canvas.FontBlack), fmt.Sprint(*card.Attributes.Cost),
				costContainerR*2, 0,
				canvas.Center, canvas.Center, 0, 0))
		}

		fb.drawCardText(ctx, card, fontSizeCard, 0, 0, boxText, fb.getAdditionalText()...)
		fb.drawTypeText(ctx, card, fontSizeCard, boxType)

		return nil
	})
}
