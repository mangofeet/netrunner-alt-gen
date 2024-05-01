package basic

import (
	"fmt"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func (fb FrameBasic) Agenda() art.Drawer {

	return art.DrawerFunc(func(ctx *canvas.Context, card *nrdb.Printing) error {

		canvasWidth, canvasHeight := ctx.Size()

		strokeWidth := getStrokeWidth(ctx)

		factionColor := fb.getColor(card)

		ctx.Push()
		ctx.SetFillColor(bgColor)
		ctx.SetStrokeColor(textColor)
		ctx.SetStrokeWidth(strokeWidth)

		costContainerR := getCostContainerRadius(ctx)

		titleBoxHeight := getTitleBoxHeight(ctx)

		titleBoxTop := getTitleBoxTop(ctx)
		titleBoxBottom := titleBoxTop - titleBoxHeight
		titleBoxRight := canvasWidth - costContainerR*3.25
		costContainerStart := titleBoxRight

		titlePath := &canvas.Path{}
		titlePath.MoveTo(0, titleBoxTop)
		titlePath.LineTo(titleBoxRight, titleBoxTop)
		titlePath.QuadTo(titleBoxRight-(costContainerR*0.5), titleBoxBottom+(titleBoxHeight*0.5), titleBoxRight, titleBoxBottom)
		titlePath.LineTo(0, titleBoxBottom)
		titlePath.Close()

		ctx.DrawPath(0, 0, titlePath)
		ctx.Pop()

		// outline for cost circle
		ctx.Push()
		ctx.SetFillColor(bgColor)
		ctx.SetStrokeColor(textColor)
		ctx.SetStrokeWidth(strokeWidth)

		path := canvas.Circle(costContainerR)
		ctx.DrawPath(costContainerStart+(costContainerR), titleBoxTop-(titleBoxHeight*0.5), path)

		ctx.Pop()

		var boxText, boxType textBoxDimensions
		if card.Attributes.TrashCost != nil {
			boxText, boxType = fb.drawTextBoxTrashable(ctx, canvasHeight/192, cornerRounded)
		} else {
			boxText, boxType = fb.drawTextBox(ctx, canvasHeight/192, cornerRounded)
		}

		// render card text

		// not sure how these sizes actually correlate to the weird
		// pixel/mm setup I'm using, but these work
		fontSizeTitle := titleBoxHeight * 2
		fontSizeCost := titleBoxHeight * 3
		fontSizeCard := titleBoxHeight * 1.2

		fb.drawAgendaPoints(ctx, card, fontSizeCost)

		titleTextMaxWidth := titleBoxRight * 0.8
		titleTextX := titleBoxRight - titleTextMaxWidth - titleBoxHeight*0.5
		titleText := getTitleText(ctx, card, fontSizeTitle, titleTextMaxWidth, titleBoxHeight, canvas.Right)
		titleTextY := (titleBoxTop - (titleBoxHeight-titleText.Bounds().H)*0.5)

		ctx.DrawText(titleTextX, titleTextY, titleText)

		if card.Attributes.AdvancementRequirement != nil {
			costTextX := costContainerStart
			costTextY := titleBoxBottom + titleBoxHeight/2
			ctx.DrawText(costTextX, costTextY, canvas.NewTextBox(
				getFont(fontSizeCost, canvas.FontBlack), fmt.Sprint(*card.Attributes.AdvancementRequirement),
				costContainerR*2, 0,
				canvas.Center, canvas.Center, 0, 0))
		}

		drawCardText(ctx, card, fontSizeCard, boxText.height*0.6, canvasWidth*0.02, boxText, fb.getAdditionalText()...)
		drawTypeText(ctx, card, fontSizeCard, boxType)

		fb.drawInfluenceAndOrFactionSymbol(ctx, card, boxText.left, factionColor)

		return nil
	})
}
