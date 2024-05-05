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
	RNG                 prng.Generator
	Direction           string
	DirectionVariance   int64
	X, Y, Vx, Vy        float64
	Color               color.RGBA
	Noise               opensimplex.Noise
	Obstacles           []*canvas.Path
	Grid                bool
	StrokeWidth         float64
	stepCount           int
	DirectionChangeStep float64
	prev                *point
	nextGridDirection   string
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

	ctx.MoveTo(wlk.X, wlk.Y)
	ctx.LineTo(wlk.prev.x, wlk.prev.y)
	ctx.Stroke()

	wlk.prev = &point{wlk.X, wlk.Y}
}

func (wlk Walker) drawLine(ctx *canvas.Context, x1, y1, x2, y2 float64) {
}

func (wlk *Walker) Move() {
	wlk.updateVelocity(noiseStepFactor, false)

	wlk.stepCount++

	if wlk.Grid {
		switch wlk.nextGridDirection {
		case "right":
			wlk.X += wlk.Vx
		case "up":
			wlk.Y += wlk.Vy
		case "down":
			wlk.Y -= wlk.Vy
		case "left":
			wlk.X -= wlk.Vx
		}
	} else {
		wlk.X += wlk.Vx
		wlk.Y += wlk.Vy
	}

	wlk.maybeChangeDirection()
}

func (wlk *Walker) updateVelocity(factor float64, hasChangedDirection bool) {
	deltaX := wlk.Noise.Eval2(wlk.X*factor, wlk.Y*factor)
	deltaY := wlk.Noise.Eval2(wlk.Y*factor, wlk.X*factor)

	var newVx, newVy float64

	switch strings.ToLower(wlk.Direction) {
	case "down":
		newVx = wlk.Vx + deltaX
		newVy = wlk.Vy + -1*math.Abs(deltaY)
	case "up":
		newVx = wlk.Vx + deltaX
		newVy = wlk.Vy + math.Abs(deltaY)
	case "left":
		newVx = wlk.Vx + -1*math.Abs(deltaX)
		newVy = wlk.Vy + deltaY
	case "right":
		newVx = wlk.Vx + math.Abs(deltaX)
		newVy = wlk.Vy + deltaY
	default:
		newVx = wlk.Vx + deltaX
		newVy = wlk.Vy + deltaY
	}

	if wlk.Grid {
		switch wlk.RNG.Next(4) {
		case 1:
			wlk.nextGridDirection = "right"
		case 2:
			wlk.nextGridDirection = "up"
		case 3:
			wlk.nextGridDirection = "down"
		case 4:
			wlk.nextGridDirection = "left"

		}
	}

	if wlk.willCollide(newVx, newVy) {
		if !hasChangedDirection {
			wlk.reverse()
			if math.Abs(newVx) > math.Abs(newVy) {
				newVx *= -1
			} else {
				newVy *= -1
			}
		} else {
			wlk.updateVelocity(factor*10, true)
			return
		}

	}

	wlk.Vx = newVx
	wlk.Vy = newVy
}

func (wlk Walker) willCollide(vx, vy float64) bool {

	path := &canvas.Path{}

	x, y := wlk.X, wlk.Y

	var newX, newY float64
	path.MoveTo(x, y)

	if wlk.Grid {
		switch wlk.nextGridDirection {
		case "right":
			newX = x + vx
			newY = y
		case "up":
			newX = x
			newY = y + vy
		case "down":
			newX = x
			newY = y - vy
		case "left":
			newX = x - vx
			newY = y
		}
	} else {
		newX = x + vx
		newY = y + vy
	}

	path.LineTo(newX, newY)

	for _, obs := range wlk.Obstacles {
		if path.Intersects(obs) || obs.Contains(newX, newY) {
			return true
		}
	}

	return false

}

func (wlk *Walker) maybeChangeDirection() {
	if wlk.Direction == "" {
		return
	}

	// if no variance, just return here
	if wlk.DirectionVariance == 0 {
		return
	}

	// see if the diructon change step was set
	if wlk.DirectionChangeStep == 0 {
		wlk.DirectionChangeStep = 30
	}

	// if we are not at a potential direction change, just return
	if wlk.stepCount%int(wlk.DirectionChangeStep) != 0 {
		return
	}

	// up the number for the next direction change
	wlk.DirectionChangeStep *= 3

	// normalize variance
	if wlk.DirectionVariance > 4 {
		wlk.DirectionVariance = 4
	}
	if wlk.DirectionVariance <= 0 {
		wlk.DirectionVariance = 1
	}

	// check for direction change
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
