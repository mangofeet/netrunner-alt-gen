package netrunner

import (
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func DrawFrameProgram(ctx *canvas.Context, card *nrdb.Printing) error {

	canvasWidth, canvasHeight := ctx.Size()

	ctx.Push()
	ctx.SetFillColor(bgColor)

	titleBoxHeight := (canvasHeight / 12)

	titleBoxTop := canvasHeight - (canvasHeight / 32)
	titleBoxBottom := titleBoxTop - titleBoxHeight
	titleBoxArcStart := canvasWidth - (canvasWidth / 2)
	titleBoxRight := canvasWidth - (canvasWidth / 8)
	titleBoxArcCP1 := titleBoxRight - (canvasWidth / 48)

	costContainerR := titleBoxHeight * 0.667
	costContainerStart := costContainerR * 0.667

	titlePath := &canvas.Path{}
	titlePath.MoveTo(0, titleBoxTop)

	titlePath.LineTo(costContainerStart, titleBoxTop)

	titlePath.QuadTo(costContainerStart+costContainerR, titleBoxTop+costContainerR, costContainerStart+(costContainerR*2), titleBoxTop)

	titlePath.LineTo(titleBoxArcStart, titleBoxTop)
	titlePath.QuadTo(titleBoxArcCP1, titleBoxTop, titleBoxRight, titleBoxBottom)

	titlePath.LineTo(costContainerStart+(costContainerR*2), titleBoxBottom)
	titlePath.QuadTo(costContainerStart+costContainerR, titleBoxBottom-costContainerR, costContainerStart, titleBoxBottom)

	titlePath.LineTo(0, titleBoxBottom)
	titlePath.Close()

	ctx.DrawPath(0, 0, titlePath)
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
