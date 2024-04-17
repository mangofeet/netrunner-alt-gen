package walker

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
	Seed         string
	Sequence     int
	Direction    string
	X, Y, Vx, Vy float64
	Color        color.Color
	Noise        opensimplex.Noise
	Grid         bool
	StrokeWidth  float64
	prev         *point
}

func (wlk Walker) String() string {
	return fmt.Sprintf("walker %d at (%f, %f), direction=%s, grid=%t", wlk.Sequence, wlk.X, wlk.Y, wlk.Direction, wlk.Grid)
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
	if wlk.Grid {
		switch prng.SequenceNext(wlk.Sequence, wlk.Seed, 4) {
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
