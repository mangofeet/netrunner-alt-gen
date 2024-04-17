package netrunner

import (
	"image/color"

	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func DrawFrameProgram(ctx *canvas.Context, card *nrdb.Printing, bgColor, textColor color.Color) error {

	canvasWidth, canvasHeight := ctx.Size()

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.MoveTo(0, canvasHeight)
	ctx.LineTo(canvasWidth-(canvasWidth/4), canvasHeight)
	ctx.LineTo(canvasWidth-(canvasWidth/4), canvasHeight-(canvasHeight/12))
	ctx.LineTo(0, canvasHeight-(canvasHeight/12))
	ctx.Close()
	ctx.Fill()

	ctx.MoveTo(canvasWidth/6, 0)
	ctx.LineTo(canvasWidth, 0)
	ctx.LineTo(canvasWidth, canvasHeight/3)
	ctx.LineTo(canvasWidth/6, canvasHeight/3)
	ctx.Close()
	ctx.Fill()

	ctx.Pop()

	return nil
}
