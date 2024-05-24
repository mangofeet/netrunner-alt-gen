package basic

import (
	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func (fb FrameBasic) Tracker() art.Drawer {

	return art.DrawerFunc(func(ctx *canvas.Context, card *nrdb.Printing) error {

		canvasWidth, canvasHeight := ctx.Size()

		fontSize := canvasHeight * 0.25
		textMaxHeight := canvasHeight * 0.15
		// textMaxWidth := canvasWidth * 0.7

		text := fb.getVerticalFittedText(ctx, card.Attributes.Title, fontSize, 0, textMaxHeight, canvas.Center)

		boxMarginTB := text.Bounds().H * 0.2

		boxTop := text.Bounds().H + (boxMarginTB * 2) + (canvasHeight * 0.03348)
		boxBottom := 0.0
		boxRadius := canvasWidth * 0.03
		boxWidth := text.Bounds().W * 1.1
		boxLeft := (canvasWidth * 0.5) - boxWidth*0.5
		boxRight := (canvasWidth * 0.5) + boxWidth*0.5
		boxMiddle := boxLeft + boxWidth*0.5

		oppositeBoxTop := canvasHeight
		oppositeBoxBottom := canvasHeight - text.Bounds().H - (boxMarginTB * 2) - (canvasHeight * 0.03348)

		fb.drawRoundedBox(ctx, boxTop, boxBottom, boxLeft, boxRight, boxRadius)
		fb.drawRoundedBox(ctx, oppositeBoxTop, oppositeBoxBottom, boxLeft, boxRight, boxRadius)

		textX := boxMiddle
		textY := boxTop - boxMarginTB

		oppositeTextX := boxMiddle
		oppositeTextY := oppositeBoxBottom + boxMarginTB

		ctx.DrawText(textX, textY, text)

		ctx.Push()
		ctx.Rotate(180)
		ctx.DrawText(oppositeTextX, oppositeTextY, text)
		ctx.Pop()

		return nil
	})
}
