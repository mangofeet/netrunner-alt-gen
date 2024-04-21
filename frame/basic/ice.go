package basic

import (
	"fmt"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type FrameIce struct{}

func (FrameIce) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	canvasWidth, canvasHeight := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	factionBaseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	factionColor := art.Darken(factionBaseColor, 0.811)

	// res cost icon
	rezCostImage, err := loadGameAsset("REZ_COST")
	if err != nil {
		return err
	}
	rezCostImage = rezCostImage.Transform(canvas.Identity.ReflectY()).Scale(0.13, 0.13)

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth * 0.5)

	costIconX := canvasWidth * 0.04
	costIconY := canvasHeight - costIconX

	ctx.DrawPath(costIconX, costIconY, rezCostImage)

	ctx.Pop()

	// title box

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	titleBoxHeight := getTitleBoxHeight(ctx)

	titleBoxTop := getTitleBoxTop(ctx)
	titleBoxBottom := titleBoxTop - titleBoxHeight
	// titleBoxRight := canvasWidth - (canvasWidth / 16)
	titleBoxLeft := costIconX + (rezCostImage.Bounds().W * 1.1)
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

	ctx.DrawPath(0, 0, titlePath)
	ctx.Pop()

	// type box
	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)
	typeBoxWidth := titleBoxHeight * 0.75
	typeBoxTop := costIconY - (rezCostImage.Bounds().H * 1.1)
	typeBoxLeft := costIconX + (rezCostImage.Bounds().W * 0.52) - (typeBoxWidth * 0.5)

	typeBoxRadius := titleBoxRadius * 0.75
	typeBoxArc1StartY := typeBoxTop - typeBoxRadius
	typeBoxArc1EndX := typeBoxLeft + typeBoxRadius
	typeBoxArc2StartX := typeBoxLeft + typeBoxWidth - typeBoxRadius
	typeBoxArc2EndY := typeBoxTop - typeBoxRadius

	typePath := &canvas.Path{}
	typePath.MoveTo(typeBoxLeft, 0)

	typePath.LineTo(typeBoxLeft, typeBoxArc1StartY)
	typePath.QuadTo(typeBoxLeft, typeBoxTop, typeBoxArc1EndX, typeBoxTop)
	typePath.LineTo(typeBoxArc2StartX, typeBoxTop)
	typePath.QuadTo(typeBoxLeft+typeBoxWidth, typeBoxTop, typeBoxLeft+typeBoxWidth, typeBoxArc2EndY)
	typePath.LineTo(typeBoxLeft+typeBoxWidth, 0)

	typePath.Close()
	ctx.DrawPath(0, 0, typePath)

	ctx.Pop()

	// text box
	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	textBoxBottom := canvasHeight * 0.5
	textBoxTop := titleBoxBottom - titleBoxHeight*0.5
	textBoxLeft := typeBoxLeft + typeBoxWidth + titleBoxHeight*0.5

	textBoxArc1StartY := textBoxTop - titleBoxRadius
	textBoxArc1EndX := textBoxLeft + titleBoxRadius
	textBoxArc2StartX := textBoxLeft + titleBoxRadius
	textBoxArc2EndY := textBoxBottom + titleBoxRadius

	textPath := &canvas.Path{}
	textPath.MoveTo(textBoxLeft, textBoxArc1StartY)
	textPath.QuadTo(textBoxLeft, textBoxTop, textBoxArc1EndX, textBoxTop)
	textPath.LineTo(canvasWidth, textBoxTop)
	textPath.LineTo(canvasWidth, textBoxBottom)
	textPath.LineTo(textBoxArc2StartX, textBoxBottom)
	textPath.QuadTo(textBoxLeft, textBoxBottom, textBoxLeft, textBoxArc2EndY)
	textPath.Close()

	ctx.DrawPath(0, 0, textPath)

	ctx.Pop()

	// program strength
	ctx.Push()

	ctx.SetFillColor(factionColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	// ctx.DrawPath(canvasWidth*0.035, 0, strength(canvasWidth, canvasHeight))
	ctx.DrawPath(0, 0, strength(canvasWidth, canvasHeight))

	ctx.Pop()

	influenceX := canvasWidth - (canvasWidth / 8)
	drawInfluence(ctx, card, influenceX, factionColor)

	fontSizeTitle := titleBoxHeight * 2
	fontSizeCost := titleBoxHeight * 3
	fontSizeStr := titleBoxHeight * 4
	fontSizeCard := titleBoxHeight * 1.2

	titleTextX := titleBoxLeft + titleBoxHeight*0.3
	titleTextY := titleBoxTop - titleBoxHeight*0.1
	ctx.DrawText(titleTextX, titleTextY, getCardText(getTitleText(card), fontSizeTitle, canvasWidth-titleBoxLeft, titleBoxHeight, canvas.Left))
	// ctx.DrawText(titleTextX, titleTextY, canvas.NewTextLine(getFont(fontSizeTitle, canvas.FontRegular), getTitleText(card), canvas.Left))

	if card.Attributes.Cost != nil {
		costTextX := costIconX * 1.07
		costTextY := costIconY - rezCostImage.Bounds().H*0.5
		ctx.DrawText(costTextX, costTextY, canvas.NewTextBox(
			getFont(fontSizeCost, canvas.FontBlack), fmt.Sprint(*card.Attributes.Cost),
			rezCostImage.Bounds().W, 0,
			canvas.Center, canvas.Center, 0, 0))
	}

	strengthText := "-"
	if card.Attributes.Strength != nil {
		strengthText = fmt.Sprint(*card.Attributes.Strength)
	}

	strTextX := canvasWidth * 0.13
	strTextY := canvasHeight * 0.02
	ctx.Push()
	ctx.Rotate(90)
	ctx.DrawText(strTextX, strTextY, canvas.NewTextBox(
		getFont(fontSizeStr, canvas.FontBlack), strengthText,
		canvasWidth/5, 0,
		canvas.Center, canvas.Center, 0, 0))
	ctx.Pop()

	drawCardText(ctx, card, fontSizeCard, textBoxTop-textBoxBottom, 0, textBoxDimensions{
		left:   textBoxLeft,
		right:  canvasWidth - canvasWidth*0.05,
		height: textBoxTop - textBoxBottom,
		bottom: textBoxBottom,
		top:    textBoxTop,
	})

	ctx.Push()
	ctx.Rotate(90)

	typeText := getTypeText(card, fontSizeCard, typeBoxTop-typeBoxWidth*0.3, typeBoxWidth, canvas.Right)
	ctx.DrawText(typeBoxLeft+typeBoxWidth*0.16, 0, typeText)

	ctx.Pop()

	return nil
}