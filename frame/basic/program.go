package basic

import (
	"fmt"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func (fb FrameBasic) Program() art.Drawer {
	return art.DrawerFunc(func(ctx *canvas.Context, card *nrdb.Printing) error {

		canvasWidth, canvasHeight := ctx.Size()

		strokeWidth := getStrokeWidth(ctx)

		ctx.Push()
		ctx.SetFillColor(fb.getColorBG())
		ctx.SetStrokeColor(fb.getColorBorder())
		ctx.SetStrokeWidth(strokeWidth)

		titleBoxHeight := getTitleBoxHeight(ctx)

		titleBoxTop := getTitleBoxTop(ctx)
		titleBoxBottom := titleBoxTop - titleBoxHeight
		titleBoxArcStart := canvasWidth - (canvasWidth / 3)
		titleBoxRight := canvasWidth - (canvasWidth / 16)
		titleBoxArcCP1 := titleBoxRight - (canvasWidth / 48)

		costContainerR := getCostContainerRadius(ctx)
		costContainerStart := getCostContainerStart(ctx)

		titlePath := &canvas.Path{}
		titlePath.MoveTo(0, titleBoxTop)

		// background for cost, top
		titlePath.LineTo(costContainerStart, titleBoxTop)
		titlePath.QuadTo(costContainerStart+costContainerR, titleBoxTop+(costContainerR), costContainerStart+(costContainerR*2), titleBoxTop)

		// arc down on right side
		titlePath.LineTo(titleBoxArcStart, titleBoxTop)
		titlePath.QuadTo(titleBoxArcCP1, titleBoxTop, titleBoxRight, titleBoxBottom)

		// background for cost, bottom
		titlePath.LineTo(costContainerStart+(costContainerR*2), titleBoxBottom)
		titlePath.QuadTo(costContainerStart+costContainerR, titleBoxBottom-(costContainerR), costContainerStart, titleBoxBottom)

		// finish title box
		titlePath.LineTo(0, titleBoxBottom)
		titlePath.Close()

		ctx.DrawPath(0, 0, titlePath)
		ctx.Pop()

		fb.drawCostCircle(ctx, transparent)

		boxText, boxType := fb.drawTextBox(ctx, canvasHeight/48, cornerRounded)

		fb.drawInfluenceAndOrFactionSymbol(ctx, card, boxText.right)

		// program strength
		ctx.Push()

		ctx.SetFillColor(fb.getColorStrenthBG(card))
		ctx.SetStrokeColor(fb.getColorBorder())
		ctx.SetStrokeWidth(strokeWidth)

		strengthPath := strength(canvasWidth, canvasHeight)
		ctx.DrawPath(canvasWidth*-0.04, 0, strengthPath)

		ctx.Pop()

		fb.drawMU(ctx, card, true)
		// render card text

		// not sure how these sizes actually correlate to the weird
		// pixel/mm setup I'm using, but these work
		fontSizeTitle := titleBoxHeight * 2
		fontSizeCost := titleBoxHeight * 3
		fontSizeStr := titleBoxHeight * 4
		fontSizeCard := titleBoxHeight * 1.5

		titleTextX := costContainerStart + (costContainerR * 2) + (costContainerR / 3)
		titleTextY := titleBoxTop - titleBoxHeight*0.1
		ctx.DrawText(titleTextX, titleTextY, fb.getCardText(getTitle(card), fontSizeTitle, titleBoxRight, titleBoxHeight, canvas.Left))
		// ctx.DrawText(titleTextX, titleTextY, canvas.NewTextLine(getFont(fontSizeTitle, canvas.FontRegular), getTitleText(card), canvas.Left))

		if card.Attributes.Cost != nil {
			costTextX := costContainerStart
			costTextY := titleBoxBottom + titleBoxHeight/2
			ctx.DrawText(costTextX, costTextY, canvas.NewTextBox(
				fb.getFont(fontSizeCost, canvas.FontBlack), fmt.Sprint(*card.Attributes.Cost),
				costContainerR*2, 0,
				canvas.Center, canvas.Center, 0, 0))
		}

		strengthText := "-"
		if card.Attributes.Strength != nil {
			strengthText = fmt.Sprint(*card.Attributes.Strength)
		}

		strTextX := canvasWidth * 0.01
		strTextY := canvasHeight / 10
		ctx.DrawText(strTextX, strTextY, canvas.NewTextBox(
			fb.getFontWithColor(fontSizeStr, canvas.FontBlack, fb.getColorTextStrength()), strengthText,
			canvasWidth/5, 0,
			canvas.Center, canvas.Center, 0, 0))

		fb.drawCardText(ctx, card, fontSizeCard, strengthPath.Bounds().H, canvasWidth*0.06, boxText, fb.getAdditionalText()...)
		fb.drawTypeText(ctx, card, fontSizeCard, boxType)

		return nil
	})
}
