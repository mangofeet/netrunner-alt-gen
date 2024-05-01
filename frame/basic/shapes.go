package basic

import (
	"fmt"
	"image/color"
	"image/png"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/assets"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func getStrokeWidth(ctx *canvas.Context) float64 {
	_, canvasHeight := ctx.Size()
	// return canvasHeight * 0.0023
	return canvasHeight * 0.0016
}

func getTitleBoxHeight(ctx *canvas.Context) float64 {
	_, canvasHeight := ctx.Size()
	return (canvasHeight / 16)
}

func getTitleBoxTop(ctx *canvas.Context) float64 {
	_, canvasHeight := ctx.Size()
	return canvasHeight - (canvasHeight / 12)
}

func getCostContainerRadius(ctx *canvas.Context) float64 {
	return getTitleBoxHeight(ctx) * 0.667
}

func getCostContainerStart(ctx *canvas.Context) float64 {
	return getCostContainerRadius(ctx) * 1.3
}

func (fb FrameBasic) getTextBoxHeight(ctx *canvas.Context) float64 {
	_, canvasHeight := ctx.Size()
	factor := 0.3333
	if fb.TextBoxHeightFactor != nil {
		factor = *fb.TextBoxHeightFactor
	}
	return (canvasHeight * factor)
}

func getInfluenceHeight(ctx *canvas.Context) float64 {
	_, canvasHeight := ctx.Size()
	return canvasHeight * 0.28
}

func drawCostCircle(ctx *canvas.Context, bgColor color.Color) {

	strokeWidth := getStrokeWidth(ctx)
	costContainerR := getCostContainerRadius(ctx)
	costContainerStart := getCostContainerStart(ctx)
	titleBoxHeight := getTitleBoxHeight(ctx)
	titleBoxTop := getTitleBoxTop(ctx)

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	path := canvas.Circle(costContainerR)
	ctx.DrawPath(costContainerStart+(costContainerR), titleBoxTop-(titleBoxHeight*0.5), path)

	ctx.Pop()

}

func drawRezCost(ctx *canvas.Context, card *nrdb.Printing, fontSize float64) (*textBoxDimensions, error) {
	canvasWidth, canvasHeight := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	// res cost icon
	rezCostImage, err := loadGameAsset("REZ_COST")
	if err != nil {
		return nil, err
	}
	rezCostImage = rezCostImage.Transform(canvas.Identity.ReflectY()).Scale(0.1, 0.1)

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth * 0.5)

	costIconX := canvasWidth * 0.066
	costIconY := canvasHeight - costIconX

	ctx.DrawPath(costIconX, costIconY, rezCostImage)

	ctx.Pop()

	if card.Attributes.Cost != nil {
		costTextX := costIconX * 1.03
		costTextY := costIconY - rezCostImage.Bounds().H*0.5
		ctx.DrawText(costTextX, costTextY, canvas.NewTextBox(
			getFont(fontSize, canvas.FontBlack), fmt.Sprint(*card.Attributes.Cost),
			rezCostImage.Bounds().W, 0,
			canvas.Center, canvas.Center, 0, 0))
	}

	return &textBoxDimensions{
		top:    costIconY,
		left:   costIconX,
		width:  rezCostImage.Bounds().W,
		height: rezCostImage.Bounds().H,
	}, nil
}

func (fb FrameBasic) drawAgendaPoints(ctx *canvas.Context, card *nrdb.Printing, fontSize float64) (*textBoxDimensions, error) {
	canvasWidth, _ := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	// res cost icon
	icon, err := loadGameAsset("AGENDA")
	if err != nil {
		return nil, err
	}
	icon = icon.Transform(canvas.Identity.ReflectY()).Scale(0.07, 0.07)

	iconX := canvasWidth * 0.085
	iconY := fb.getTextBoxHeight(ctx) + icon.Bounds().H*1.8
	iconColor := color.RGBA{
		R: textColor.R,
		G: textColor.G,
		B: textColor.B,
		A: 0x44,
	}

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)
	circle := canvas.Circle(icon.Bounds().H * 0.6)
	ctx.DrawPath(iconX+icon.Bounds().W*0.53, iconY-icon.Bounds().H*0.46, circle)
	ctx.Pop()

	ctx.Push()
	ctx.SetFillColor(iconColor)
	ctx.DrawPath(iconX, iconY, icon)
	ctx.Pop()

	if card.Attributes.AgendaPoints != nil {
		textX := iconX * 1.03
		textY := iconY - icon.Bounds().H*0.4
		ctx.DrawText(textX, textY, canvas.NewTextBox(
			getFont(fontSize, canvas.FontBlack), fmt.Sprint(*card.Attributes.AgendaPoints),
			icon.Bounds().W, 0,
			canvas.Center, canvas.Center, 0, 0))
	}

	return &textBoxDimensions{
		top:    iconY,
		left:   iconX,
		width:  icon.Bounds().W,
		height: icon.Bounds().H,
	}, nil
}

func loadTrashCostPath() *canvas.Path {

	trashScale := 0.005
	paths := make([]*canvas.Path, 5)
	for i := range 5 {
		path := mustLoadGameAsset(fmt.Sprintf("TRASH_COST_%d", i))
		path = path.Transform(canvas.Identity.ReflectY()).Scale(trashScale, trashScale)
		paths[i] = path
	}

	return paths[0].Join(paths[1]).Join(paths[2]).Join(paths[3]).Join(paths[4])

}

func drawTrashCost(ctx *canvas.Context, card *nrdb.Printing) (*textBoxDimensions, error) {

	if card.Attributes.TrashCost == nil {
		return nil, nil
	}

	canvasWidth, canvasHeight := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	// res cost icon
	image := loadTrashCostPath()

	fontSize := image.Bounds().H * 2
	iconX := canvasWidth * 0.815
	iconY := canvasHeight * 0.145

	ctx.Push()
	ctx.SetFill(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)
	ctx.DrawPath(iconX, iconY, image)
	ctx.Pop()

	if card.Attributes.TrashCost != nil {
		textX := iconX + image.Bounds().W*0.45
		textY := iconY - image.Bounds().H
		ctx.DrawText(textX, textY, canvas.NewTextBox(
			getFont(fontSize, canvas.FontBlack), fmt.Sprint(*card.Attributes.TrashCost),
			image.Bounds().W, 0,
			canvas.Center, canvas.Center, 0, 0))
	}

	return &textBoxDimensions{
		top:    iconY,
		left:   iconX,
		width:  image.Bounds().W,
		height: image.Bounds().H,
	}, nil
}

func drawMU(ctx *canvas.Context, card *nrdb.Printing, drawBox bool) {
	canvasWidth, _ := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	// mu icon
	muImage, err := loadGameAsset("MU")
	if err != nil {
		panic(err)
	}
	muImage = muImage.Transform(canvas.Identity.ReflectY()).Scale(0.05, 0.05)

	muBoxX := canvasWidth * 0.0853
	muBoxY := (getTitleBoxTop(ctx) - getTitleBoxHeight(ctx)) - (muImage.Bounds().H * 0.8)
	muBoxW := muImage.Bounds().W + muImage.Bounds().W*0.35
	muBoxH := muImage.Bounds().H + muImage.Bounds().H*0.45

	if drawBox {
		ctx.Push()
		ctx.SetFillColor(bgColor)
		ctx.SetStrokeColor(textColor)
		ctx.SetStrokeWidth(strokeWidth)

		boxPath := &canvas.Path{}
		boxPath.MoveTo(0, 0)
		boxPath.LineTo(muBoxW, 0)
		boxPath.LineTo(muBoxW, -1*muBoxH)
		boxPath.LineTo(0, -1*muBoxH)
		boxPath.Close()
		ctx.DrawPath(muBoxX, muBoxY, boxPath)

		ctx.Pop()
	}

	ctx.Push()
	ctx.SetFillColor(textColor)

	muIconX := muBoxX
	muIconY := muBoxY + (muImage.Bounds().H * 0.2)

	ctx.DrawPath(muIconX, muIconY, muImage)

	ctx.Pop()

	var muText string
	switch card.Attributes.CardTypeID {
	case "program":
		if card.Attributes.MemoryCost != nil {
			muText = fmt.Sprint(*card.Attributes.MemoryCost)
		}
	case "runner_identity":
		// TODO: see if this is actually coming back yet from the preview API
		muText = "4"
		if card.Attributes.CardAbilities.MUProvided != nil {
			muText = fmt.Sprint(*card.Attributes.CardAbilities.MUProvided)
		}

	}

	muTextX := muBoxX - muBoxW*0.08
	muTextY := muBoxY
	fontSize := getTitleBoxHeight(ctx) * 1.2
	ctx.DrawText(muTextX, muTextY, canvas.NewTextBox(
		getFont(fontSize, canvas.FontBlack), muText,
		muBoxW, muBoxH,
		canvas.Center, canvas.Center, 0, 0))

}

func drawLink(ctx *canvas.Context, card *nrdb.Printing) {
	canvasWidth, _ := ctx.Size()

	// link icon
	icon, err := loadGameAsset("LINK")
	if err != nil {
		panic(err)
	}
	icon = icon.Transform(canvas.Identity.ReflectY()).Scale(0.015, 0.015)

	boxX := canvasWidth * 0.1
	boxY := getTitleBoxTop(ctx) - getTitleBoxHeight(ctx)*0.6
	boxW := icon.Bounds().W + icon.Bounds().W*2.7
	boxH := icon.Bounds().H + icon.Bounds().H*1.8

	ctx.Push()
	ctx.SetFillColor(textColor)

	iconX := boxX + (icon.Bounds().W * 2.5)
	iconY := boxY + (icon.Bounds().H * 1.1)

	ctx.DrawPath(iconX, iconY, icon)

	ctx.Pop()

	text := "0"
	if card.Attributes.BaseLink != nil {
		text = fmt.Sprint(*card.Attributes.BaseLink)
	}

	textX := boxX + boxW*0.0
	textY := boxY + boxW*0.2
	fontSize := getTitleBoxHeight(ctx) * 2
	ctx.DrawText(textX, textY, canvas.NewTextBox(
		getFont(fontSize, canvas.FontBlack), text,
		boxW, boxH,
		canvas.Center, canvas.Center, 0, 0))
}

func (fb FrameBasic) drawInfluenceAndOrFactionSymbol(ctx *canvas.Context, card *nrdb.Printing, x float64, bgColor color.RGBA) {

	strokeWidth := getStrokeWidth(ctx)

	_, canvasHeight := ctx.Size()

	// influenceHeight := math.Max(fb.getTextBoxHeight(ctx)*0.8, canvasHeight*0.26664)
	influenceHeight := getInfluenceHeight(ctx)
	influenceWidth := canvasHeight / 42
	factionY := influenceHeight*0.2 + influenceWidth*1.2

	if card.Attributes.InfluenceCost != nil {

		influenceCost := *card.Attributes.InfluenceCost

		ctx.Push()
		ctx.SetFillColor(bgColor)
		ctx.SetStrokeColor(textColor)
		ctx.SetStrokeWidth(strokeWidth)

		// center around the give point
		boxX := x - (influenceWidth / 2)

		ctx.DrawPath(boxX, 0, influenceBox(influenceHeight, influenceWidth))

		ctx.Pop()

		curveRadius := influenceWidth / 2

		pipR := curveRadius * 0.6

		var pipY float64
		for i := 0.0; i < 5; i += 1 {

			pipY = influenceHeight - (pipR * ((i + 1) * 4)) + (pipR / 2)

			ctx.Push()
			ctx.SetStrokeWidth(strokeWidth * 0.75)
			ctx.SetStrokeColor(textColor)
			ctx.SetFill(transparent)

			pip := canvas.Circle(pipR)
			ctx.DrawPath(x, pipY, pip)

			ctx.Pop()

			if i >= 5-float64(influenceCost) {
				ctx.Push()
				ctx.SetFill(textColor)
				pip := canvas.Circle(pipR * 0.5)
				ctx.DrawPath(x, pipY, pip)
				ctx.Pop()
			}
		}
	}

	if err := drawFactionSybmol(ctx, card, x, factionY, influenceWidth*2); err != nil {
		panic(err)
	}

}

func drawFactionSybmol(ctx *canvas.Context, card *nrdb.Printing, x, y, width float64) error {
	strokeWidth := getStrokeWidth(ctx)

	bubbleRadius := width * 0.6

	circle := canvas.Circle(bubbleRadius)
	ctx.Push()
	ctx.SetFill(bgColorOpaque)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)
	ctx.DrawPath(x, y, circle)
	ctx.Pop()

	factionImageFile, err := assets.FS.Open(fmt.Sprintf("%s.png", card.Attributes.FactionID))
	if err != nil {
		return err
	}
	factionImage, err := png.Decode(factionImageFile)
	if err != nil {
		return err
	}

	factionImageWidth := float64(factionImage.Bounds().Max.X)
	factionImageHeight := float64(factionImage.Bounds().Max.Y)
	factionScaleFactor := width / factionImageWidth
	if strings.Contains(card.Attributes.FactionID, "neutral") {
		factionScaleFactor = (width * 0.85) / factionImageWidth
	}

	factionImageX := x - factionImageWidth*0.5*factionScaleFactor
	factionImageY := y - factionImageHeight*factionScaleFactor*0.5

	ctx.RenderImage(
		factionImage,
		canvas.Identity.
			Translate(factionImageX, factionImageY).
			Scale(factionScaleFactor, factionScaleFactor),
	)

	return nil

}

func strength(canvasWidth, canvasHeight float64) *canvas.Path {
	path := &canvas.Path{}

	path.MoveTo(0, canvasHeight*0.12)
	// path.CubeTo(canvasWidth*0.166, canvasHeight*0.25, canvasWidth*0.25, canvasHeight*0.08333, canvasWidth*0.125, 0)
	path.CubeTo(canvasWidth*0.166, canvasHeight*0.25, canvasWidth*0.35, canvasHeight*0.08, canvasWidth*0.17, 0)
	path.LineTo(0, 0)
	path.Close()

	return path
}

func influenceBox(height, width float64) *canvas.Path {

	path := &canvas.Path{}

	curveRadius := width / 2
	curveStart := height - curveRadius

	path.MoveTo(0, 0)
	path.LineTo(0, curveStart)
	path.CubeTo(0, height, width, height, width, curveStart)
	path.LineTo(width, 0)
	path.Close()

	return path
}

type corner func(ctx *canvas.Context, cx, cy, x, y float64)

var cornerRounded = corner(func(ctx *canvas.Context, cx, cy, x, y float64) {
	ctx.QuadTo(cx, cy, x, y)
})

var cornerIn = corner(func(ctx *canvas.Context, cx, cy, x, y float64) {
	var cxNew, cyNew float64

	factor := 0.65

	if x-cx > 0 {
		cxNew = cx + (x-cx)*factor
		cyNew = cy - (x-cx)*factor
	} else {
		cxNew = cx - (cy-y)*factor
		cyNew = cy - (cy-y)*factor
	}

	ctx.QuadTo(cxNew, cyNew, x, y)
})

var cornerStraight = corner(func(ctx *canvas.Context, _, _, x, y float64) {
	ctx.LineTo(x, y)
})

var cornerNone = corner(func(ctx *canvas.Context, cx, cy, x, y float64) {
	ctx.LineTo(cx, cy)
	ctx.LineTo(x, y)
})

func (fb FrameBasic) drawTextBox(ctx *canvas.Context, cornerSize float64, cornerType corner) (textBoxDimensions, textBoxDimensions) {
	canvasWidth, _ := ctx.Size()
	textBoxLeft := canvasWidth / 8
	textBoxRight := canvasWidth - (canvasWidth / 8)

	return fb.drawTextBoxToSize(ctx, textBoxLeft, textBoxRight, cornerSize, cornerType)
}

func (fb FrameBasic) drawTextBoxTrashable(ctx *canvas.Context, cornerSize float64, cornerType corner) (textBoxDimensions, textBoxDimensions) {
	canvasWidth, _ := ctx.Size()
	textBoxLeft := canvasWidth / 8
	textBoxRight := canvasWidth - (canvasWidth / 6)

	return fb.drawTextBoxToSize(ctx, textBoxLeft, textBoxRight, cornerSize, cornerType)
}

func (fb FrameBasic) drawTextBoxToSize(ctx *canvas.Context, textBoxLeft, textBoxRight, cornerSize float64, cornerType corner) (textBoxDimensions, textBoxDimensions) {

	_, canvasHeight := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	textBoxHeight := fb.getTextBoxHeight(ctx)

	textBoxArc2StartX := textBoxRight - cornerSize
	textBoxArc2EndY := textBoxHeight - cornerSize

	// text box
	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	ctx.MoveTo(textBoxLeft, 0)
	ctx.LineTo(textBoxLeft, textBoxHeight)

	ctx.LineTo(textBoxArc2StartX, textBoxHeight)

	cornerType(ctx, textBoxRight, textBoxHeight, textBoxRight, textBoxArc2EndY)

	ctx.LineTo(textBoxRight, 0)

	ctx.FillStroke()
	ctx.Pop()

	// type box
	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	// typeBoxHeight := math.Max(textBoxHeight*0.17, canvasHeight*0.056661)
	typeBoxHeight := canvasHeight * 0.056661
	typeBoxBottom := textBoxHeight + strokeWidth*0.5
	typeBoxLeft := textBoxLeft
	typeBoxRight := textBoxRight * 0.9

	typeBoxArcRadius := cornerSize
	typeBoxArc1StartY := typeBoxBottom + typeBoxHeight - typeBoxArcRadius
	typeBoxArc1EndX := typeBoxLeft + typeBoxArcRadius

	typeBoxArc2StartX := typeBoxRight - typeBoxArcRadius
	typeBoxArc2EndY := typeBoxBottom + typeBoxHeight - typeBoxArcRadius

	ctx.MoveTo(typeBoxLeft, typeBoxBottom)
	ctx.LineTo(typeBoxLeft, typeBoxArc1StartY)
	cornerType(ctx, typeBoxLeft, typeBoxHeight+typeBoxBottom, typeBoxArc1EndX, typeBoxHeight+typeBoxBottom)

	ctx.LineTo(typeBoxArc2StartX, typeBoxHeight+typeBoxBottom)
	cornerType(ctx, typeBoxRight, typeBoxHeight+typeBoxBottom, typeBoxRight, typeBoxArc2EndY)

	ctx.LineTo(typeBoxRight, typeBoxBottom)

	ctx.FillStroke()

	ctx.Pop()

	return textBoxDimensions{
			left:   textBoxLeft,
			right:  textBoxRight,
			height: textBoxHeight,
		}, textBoxDimensions{
			left:   typeBoxLeft,
			right:  typeBoxRight,
			height: typeBoxHeight,
			bottom: typeBoxBottom,
		}

}

func (fb FrameBasic) drawRunnerLimits(ctx *canvas.Context, card *nrdb.Printing, box textBoxDimensions) {

	canvasWidth, _ := ctx.Size()

	factionColor := fb.getColor(card)
	influenceColor := color.RGBA{
		R: 0x3f,
		G: 0x3f,
		B: 0x3f,
		A: 0xff,
	}
	strokeWidth := getStrokeWidth(ctx)

	width := canvasWidth * 0.1
	height := canvasWidth * 0.08
	deckBoxLeft := box.left - width*0.5
	influenceBoxLeft := box.right - width*0.5
	bottom := box.height * 0.18
	top := bottom + height
	corner := width * 0.25

	ctx.Push()
	ctx.SetFillColor(factionColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	ctx.MoveTo(deckBoxLeft, bottom)
	ctx.LineTo(deckBoxLeft, top-corner)
	ctx.LineTo(deckBoxLeft+corner, top)
	ctx.LineTo(deckBoxLeft+width, top)
	ctx.LineTo(deckBoxLeft+width, bottom+corner)
	ctx.LineTo(deckBoxLeft+width-corner, bottom)
	ctx.Close()

	ctx.FillStroke()
	ctx.Pop()

	ctx.Push()
	ctx.SetFillColor(influenceColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	ctx.MoveTo(influenceBoxLeft, bottom+corner)
	ctx.LineTo(influenceBoxLeft, top)
	ctx.LineTo(influenceBoxLeft+width-corner, top)
	ctx.LineTo(influenceBoxLeft+width, top-corner)
	ctx.LineTo(influenceBoxLeft+width, bottom)
	ctx.LineTo(influenceBoxLeft+corner, bottom)
	ctx.Close()

	ctx.FillStroke()
	ctx.Pop()

	// text
	fontSize := height * 2

	textDeckX := deckBoxLeft
	textDeckY := top
	ctx.DrawText(textDeckX, textDeckY, canvas.NewTextBox(
		getFont(fontSize, canvas.FontBlack), fmt.Sprint(*card.Attributes.MinimumDeckSize),
		width, height,
		canvas.Center, canvas.Center, 0, 0))

	textInfluenceX := influenceBoxLeft
	textInfluenceY := top
	ctx.DrawText(textInfluenceX, textInfluenceY, canvas.NewTextBox(
		getFont(fontSize, canvas.FontBlack), fmt.Sprint(*card.Attributes.InfluenceLimit),
		width, height,
		canvas.Center, canvas.Center, 0, 0))
}

func (fb FrameBasic) drawCorpLimits(ctx *canvas.Context, card *nrdb.Printing, box textBoxDimensions) {

	canvasWidth, _ := ctx.Size()

	factionColor := fb.getColor(card)
	influenceColor := color.RGBA{
		R: 0x3f,
		G: 0x3f,
		B: 0x3f,
		A: 0xff,
	}
	strokeWidth := getStrokeWidth(ctx)

	radius := canvasWidth * 0.04
	y := box.height
	deckX := box.left - radius*0.5
	influenceX := box.right + radius*0.5

	ctx.Push()
	ctx.SetFillColor(factionColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)
	ctx.DrawPath(deckX, y, canvas.Circle(radius))
	ctx.Pop()

	ctx.Push()
	ctx.SetFillColor(influenceColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)
	ctx.DrawPath(influenceX, y, canvas.Circle(radius))
	ctx.Pop()

	// text
	fontSize := radius * 4

	textDeckX := deckX - radius
	textDeckY := y + radius
	ctx.DrawText(textDeckX, textDeckY, canvas.NewTextBox(
		getFont(fontSize, canvas.FontBlack), fmt.Sprint(*card.Attributes.MinimumDeckSize),
		radius*2, radius*2,
		canvas.Center, canvas.Center, 0, 0))

	textInfluenceX := influenceX - radius
	textInfluenceY := y + radius
	ctx.DrawText(textInfluenceX, textInfluenceY, canvas.NewTextBox(
		getFont(fontSize, canvas.FontBlack), fmt.Sprint(*card.Attributes.InfluenceLimit),
		radius*2, radius*2,
		canvas.Center, canvas.Center, 0, 0))
}
