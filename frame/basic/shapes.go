package basic

import (
	"fmt"
	"image/color"

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

func getTextBoxHeight(ctx *canvas.Context) float64 {
	_, canvasHeight := ctx.Size()
	return (canvasHeight / 3)
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

func drawAgendaPoints(ctx *canvas.Context, card *nrdb.Printing, fontSize float64) (*textBoxDimensions, error) {
	canvasWidth, _ := ctx.Size()

	strokeWidth := getStrokeWidth(ctx)

	// res cost icon
	icon, err := loadGameAsset("AGENDA")
	if err != nil {
		return nil, err
	}
	icon = icon.Transform(canvas.Identity.ReflectY()).Scale(0.07, 0.07)

	iconX := canvasWidth * 0.085
	iconY := getTextBoxHeight(ctx) + icon.Bounds().H*1.8

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth * 0.5)
	circle := canvas.Circle(icon.Bounds().H * 0.6)
	ctx.DrawPath(iconX+icon.Bounds().W*0.53, iconY-icon.Bounds().H*0.46, circle)
	ctx.Pop()

	ctx.Push()
	ctx.SetFillColor(transparent)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth * 0.5)

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

func drawInfluence(ctx *canvas.Context, card *nrdb.Printing, x float64, bgColor color.RGBA) {

	if card.Attributes.InfluenceCost == nil {
		return
	}
	influenceCost := *card.Attributes.InfluenceCost
	strokeWidth := getStrokeWidth(ctx)

	_, canvasHeight := ctx.Size()

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	influenceHeight := getTextBoxHeight(ctx) * 0.65
	influenceWidth := canvasHeight / 42

	// center around the give point
	boxX := x - (influenceWidth / 2)

	ctx.DrawPath(boxX, 0, influenceBox(influenceHeight, influenceWidth))

	ctx.Pop()

	curveRadius := influenceWidth / 2

	pipR := curveRadius * 0.6

	for i := 0.0; i < 5; i += 1 {

		pipY := influenceHeight - (pipR * ((i + 1) * 4)) + (pipR / 2)

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

func drawTextBox(ctx *canvas.Context, cornerSize float64, cornerType corner) (textBoxDimensions, textBoxDimensions) {
	canvasWidth, _ := ctx.Size()
	textBoxLeft := canvasWidth / 8
	textBoxRight := canvasWidth - (canvasWidth / 8)

	return drawTextBoxToSize(ctx, textBoxLeft, textBoxRight, cornerSize, cornerType)
}

func drawTextBoxTrashable(ctx *canvas.Context, cornerSize float64, cornerType corner) (textBoxDimensions, textBoxDimensions) {
	canvasWidth, _ := ctx.Size()
	textBoxLeft := canvasWidth / 8
	textBoxRight := canvasWidth - (canvasWidth / 6)

	return drawTextBoxToSize(ctx, textBoxLeft, textBoxRight, cornerSize, cornerType)
}

func drawTextBoxToSize(ctx *canvas.Context, textBoxLeft, textBoxRight, cornerSize float64, cornerType corner) (textBoxDimensions, textBoxDimensions) {

	strokeWidth := getStrokeWidth(ctx)

	textBoxHeight := getTextBoxHeight(ctx)

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

	typeBoxHeight := textBoxHeight * 0.17
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
