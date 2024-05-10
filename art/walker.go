package art

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/ojrac/opensimplex-go"
	"github.com/tdewolff/canvas"
)

const noiseStepFactor = 0.005

type point struct {
	x, y float64
}

type Walker struct {
	RNG                         prng.Generator
	Direction                   string
	DirectionVariance           int64
	X, Y, Vx, Vy                float64
	Color                       color.Color
	Noise                       opensimplex.Noise
	Grid                        bool
	StrokeWidth                 float64
	stepCount                   int
	DirectionChangeStep         float64
	DirectionChangeStepModifier float64
	prev                        *point
}

func (wlk Walker) String() string {
	return fmt.Sprintf("walker at (%f, %f), direction=%s, grid=%t, steps=%d, directionVariance=%d", wlk.X, wlk.Y, wlk.Direction, wlk.Grid, wlk.stepCount, wlk.DirectionVariance)
}

func (wlk *Walker) Draw(ctx *canvas.Context) {

	ctx.Push()
	defer ctx.Pop()
	ctx.SetStrokeColor(wlk.Color)
	ctx.SetStrokeWidth(wlk.StrokeWidth)

	if wlk.prev == nil {
		wlk.prev = &point{wlk.X, wlk.Y}
	}

	wlk.drawLine(ctx, wlk.X, wlk.Y, wlk.prev.x, wlk.prev.y)

	wlk.prev = &point{wlk.X, wlk.Y}
}

func (wlk Walker) drawPoint(ctx *canvas.Context, x, y float64) {
	ctx.MoveTo(x, y)
	ctx.LineTo(x, y)
	ctx.Stroke()
}

func (wlk Walker) drawLine(ctx *canvas.Context, x1, y1, x2, y2 float64) {
	ctx.MoveTo(x1, y1)
	ctx.LineTo(x2, y2)
	ctx.Stroke()
}

func (wlk *Walker) Velocity() {
	deltaX := wlk.Noise.Eval2(wlk.X*noiseStepFactor, wlk.Y*noiseStepFactor)
	deltaY := wlk.Noise.Eval2(wlk.Y*noiseStepFactor, wlk.X*noiseStepFactor)

	switch strings.ToLower(wlk.Direction) {
	case "down":
		wlk.Vx += deltaX
		wlk.Vy += -1 * math.Abs(deltaY)
	case "up":
		wlk.Vx += deltaX
		wlk.Vy += math.Abs(deltaY)
	case "left":
		wlk.Vx += -1 * math.Abs(deltaX)
		wlk.Vy += deltaY
	case "right":
		wlk.Vx += math.Abs(deltaX)
		wlk.Vy += deltaY
	default:
		wlk.Vx += deltaX
		wlk.Vy += deltaY
	}

}

func (wlk *Walker) Move() {
	wlk.stepCount++

	if wlk.Grid {
		switch wlk.RNG.Next(4) {
		case 1:
			wlk.X += wlk.Vx
		case 2:
			wlk.Y += wlk.Vy
		case 3:
			wlk.Y -= wlk.Vy
		case 4:
			wlk.X -= wlk.Vx
		}
	} else {
		wlk.X += wlk.Vx
		wlk.Y += wlk.Vy
	}

	wlk.maybeChangeDirection()
}

func (wlk *Walker) maybeChangeDirection() {
	if wlk.Direction == "" {
		return
	}

	// if no variance, just return here
	if wlk.DirectionVariance == 0 {
		return
	}

	if wlk.DirectionChangeStep == 0 {
		wlk.DirectionChangeStep = 30
	}
	if wlk.DirectionChangeStepModifier == 0 {
		wlk.DirectionChangeStepModifier = 3
	}

	if wlk.stepCount%int(wlk.DirectionChangeStep) != 0 {
		return
	}

	wlk.DirectionChangeStep *= wlk.DirectionChangeStepModifier
	if wlk.DirectionChangeStep < 1 {
		wlk.DirectionChangeStep = 1
	}

	if wlk.DirectionVariance > 4 {
		wlk.DirectionVariance = 4
	}
	if wlk.DirectionVariance <= 0 {
		wlk.DirectionVariance = 1
	}
	switch wlk.RNG.Sample(int64(wlk.stepCount), wlk.DirectionVariance) {
	case 1:
		wlk.shiftRight()
	case 2:
		wlk.shiftLeft()
	case 3:
		wlk.reverse()
	}

}

func (wlk *Walker) shiftLeft() {
	switch wlk.Direction {
	case "up":
		wlk.Direction = "left"
	case "right":
		wlk.Direction = "up"
	case "down":
		wlk.Direction = "right"
	case "left":
		wlk.Direction = "down"
	}
}

func (wlk *Walker) shiftRight() {
	switch wlk.Direction {
	case "up":
		wlk.Direction = "right"
	case "right":
		wlk.Direction = "down"
	case "down":
		wlk.Direction = "left"
	case "left":
		wlk.Direction = "up"
	}
}

func (wlk *Walker) reverse() {
	switch wlk.Direction {
	case "up":
		wlk.Direction = "down"
	case "right":
		wlk.Direction = "left"
	case "down":
		wlk.Direction = "up"
	case "left":
		wlk.Direction = "right"
	}

}

func (wlk Walker) InBounds(ctx *canvas.Context) bool {
	if wlk.X < 0 {
		return false
	}
	if wlk.Y < 0 {
		return false
	}

	width, height := ctx.Size()
	return wlk.X < width && wlk.Y < height
}
