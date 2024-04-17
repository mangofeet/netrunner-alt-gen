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

	titleBoxTop := canvasHeight - (canvasHeight / 12)
	titleBoxBottom := titleBoxTop - titleBoxHeight
	titleBoxArcStart := canvasWidth - (canvasWidth / 2)
	titleBoxRight := canvasWidth - (canvasWidth / 8)
	titleBoxArcCP1 := titleBoxRight - (canvasWidth / 48)

	costContainerR := titleBoxHeight * 0.667
	costContainerStart := costContainerR

	titlePath := &canvas.Path{}
	titlePath.MoveTo(0, titleBoxTop)

	// background for cost, top
	titlePath.LineTo(costContainerStart, titleBoxTop)
	titlePath.QuadTo(costContainerStart+costContainerR, titleBoxTop+(costContainerR*0.8), costContainerStart+(costContainerR*2), titleBoxTop)

	// arc down on right side
	titlePath.LineTo(titleBoxArcStart, titleBoxTop)
	titlePath.QuadTo(titleBoxArcCP1, titleBoxTop, titleBoxRight, titleBoxBottom)

	// background for cost, bottom
	titlePath.LineTo(costContainerStart+(costContainerR*2), titleBoxBottom)
	titlePath.QuadTo(costContainerStart+costContainerR, titleBoxBottom-(costContainerR*0.8), costContainerStart, titleBoxBottom)

	// finish title box
	titlePath.LineTo(0, titleBoxBottom)
	titlePath.Close()

	ctx.DrawPath(0, 0, titlePath)
	ctx.Fill()
	ctx.Pop()

	// outline for cost circle

	ctx.Push()
	ctx.SetStrokeColor(textColor)
	ctx.SetFillColor(transparent)
	ctx.SetStrokeWidth(10)
	costOutline := canvas.Circle(costContainerR)
	ctx.DrawPath(costContainerStart+(costContainerR), titleBoxTop-(titleBoxHeight*0.5), costOutline)
	ctx.Stroke()
	ctx.Pop()

	// bottom text box

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.MoveTo(canvasWidth/6, 0)
	ctx.LineTo(canvasWidth, 0)
	ctx.LineTo(canvasWidth, canvasHeight/3)
	ctx.LineTo(canvasWidth/6, canvasHeight/3)
	ctx.Close()
	ctx.Fill()

	ctx.Pop()

	return nil
}
