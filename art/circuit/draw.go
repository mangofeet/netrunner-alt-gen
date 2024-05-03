package circuit

import (
	"image/color"
	"log"
	"math"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/mangofeet/nrdb-go"
	"github.com/ojrac/opensimplex-go"
	"github.com/tdewolff/canvas"
)

func Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	seed := card.Attributes.Title + card.Attributes.Text + card.Attributes.CardTypeID + card.Attributes.FactionID

	canvasWidth, canvasHeight := ctx.Size()

	rngGlobal := prng.NewGenerator(seed, nil)

	startX := float64(rngGlobal.Next(int64(canvasWidth/2)) + int64(canvasWidth/4))
	startY := float64(rngGlobal.Next(int64(canvasHeight/6)) + (int64(canvasHeight/8) * 5))

	if card.Attributes.CardTypeID == "ice" {
		startY = float64(rngGlobal.Next(int64(canvasHeight/4)) + (int64(canvasHeight / 6)))
	}

	baseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	cardBGColor := art.Darken(baseColor, 0.623)

	// fill background
	ctx.Push()
	ctx.SetFillColor(cardBGColor)
	ctx.MoveTo(0, 0)
	ctx.LineTo(0, canvasHeight)
	ctx.LineTo(canvasWidth, canvasHeight)
	ctx.LineTo(canvasWidth, 0)
	ctx.Close()
	ctx.Fill()
	ctx.Pop()

	noise := opensimplex.New(rngGlobal.Next(math.MaxInt64))

	strokeWidth := canvasHeight * 0.0032

	sequence := int64(1)

	path := CircuitPath{
		RNG:       prng.NewGenerator(seed, &sequence),
		Color:     baseColor,
		PathWidth: strokeWidth,
		Noise:     noise,
		X:         startX,
		Y:         startY,
	}

	path.Draw(ctx)

	return nil
}

type pathGroup struct {
	space                    float64
	direction, prevDirection string
	paths                    []*canvas.Path
}

func (pg *pathGroup) draw(ctx *canvas.Context) {
	for _, p := range pg.paths {
		ctx.DrawPath(0, 0, p)
	}
}

func (pg *pathGroup) move(distance float64) {
	dx, dy := 0.0, 0.0

	switch pg.direction {
	case "up":
		dy = distance
	case "down":
		dy = -1 * distance
	case "left":
		dx = -1 * distance
	case "right":
		dx = distance
	}

	pathCount := len(pg.paths)

	for i, p := range pg.paths {
		var thisDx, thisDy, extraDx, extraDy float64

		if pg.direction != pg.prevDirection {
			switch pg.direction {
			case "up":
				switch pg.prevDirection {
				case "left":
					extraDx = dx - pg.space*float64(pathCount-i+1)
					thisDy = dy - extraDx
				case "right":
					extraDx = dx + pg.space*float64(i+1)
					thisDy = dy + extraDx
				}

			case "down":
				switch pg.prevDirection {
				case "left":
					extraDx = dx - pg.space*float64(i+1)
					thisDy = dy + extraDx
				case "right":
					extraDx = dx + pg.space*float64(pathCount-i+1)
					thisDy = dy - extraDx
				}

			case "left":
				switch pg.prevDirection {
				case "up":
					extraDy = dy + pg.space*float64(i+1)
					thisDx = dx - extraDy
				case "down":
					extraDy = dy - pg.space*float64(pathCount-i+1)
					thisDx = dx + extraDy
				}

			case "right":
				switch pg.prevDirection {
				case "up":
					extraDy = dy + pg.space*float64(pathCount-i+1)
					thisDx = dx + extraDy
				case "down":
					extraDy = dy - pg.space*float64(i+1)
					thisDx = dx - extraDy
				}

			}
		}

		pos := p.Pos()
		p.LineTo(pos.X+extraDx, pos.Y+extraDy)
		pos = p.Pos()
		p.LineTo(pos.X+thisDx, pos.Y+thisDy)
	}

}

func (pg pathGroup) split(n int) []*pathGroup {
	if n == 0 || n == len(pg.paths) {
		return []*pathGroup{&pg}
	}

	return []*pathGroup{
		{direction: pg.direction, space: pg.space, paths: pg.paths[:n]},
		{direction: pg.direction, space: pg.space, paths: pg.paths[n:]},
	}
}

type CircuitPath struct {
	RNG       prng.Generator
	Color     color.RGBA
	PathWidth float64
	Noise     opensimplex.Noise
	X, Y      float64

	startNodeSize    float64
	startNode        canvas.Rect
	pathsPerNodeSide int
	pathGroups       []*pathGroup
}

func (wlk *CircuitPath) Draw(ctx *canvas.Context) {

	// set up things
	// wlk.pathsPerNodeSide = 2 * int(prng.SequenceNext(wlk.Sequence, wlk.Seed, 4))
	wlk.pathsPerNodeSide = 4

	canvasWidth, _ := ctx.Size()

	wlk.startNodeSize = canvasWidth*0.03 + float64(wlk.RNG.Next(int64(canvasWidth*0.05)))

	wlk.startNode = canvas.Rect{
		X: wlk.X - wlk.startNodeSize*0.5,
		Y: wlk.Y - wlk.startNodeSize*0.5,
		W: wlk.startNodeSize,
		H: wlk.startNodeSize,
	}

	space := wlk.startNodeSize / (float64(wlk.pathsPerNodeSide))
	offset := (wlk.startNodeSize - (space * (float64(wlk.pathsPerNodeSide) - 1))) / 2
	stub := 40.0

	// four sides to the start node
	wlk.pathGroups = make([]*pathGroup, 1)
	for i := range 1 {
		var direction string
		switch i {
		case 0: // left
			direction = "left"
		case 1: // right
			direction = "right"
		case 2: // top
			direction = "up"
		case 3: // bottom
			direction = "down"
		}
		wlk.pathGroups[i] = &pathGroup{
			direction: direction,
			space:     space,
			paths:     make([]*canvas.Path, wlk.pathsPerNodeSide),
		}
		for j := range wlk.pathsPerNodeSide {

			var (
				x, y, x2, y2 float64
			)

			switch i {
			case 0: // left
				x = wlk.startNode.X
				y = wlk.startNode.Y + space*float64(j) + offset
				x2 = x - stub
				y2 = y
			case 1: // right
				x = wlk.startNode.X + wlk.startNode.W
				y = wlk.startNode.Y + space*float64(j) + offset
				x2 = x + stub
				y2 = y
			case 2: // top
				x = wlk.startNode.X + space*float64(j) + offset
				y = wlk.startNode.Y + wlk.startNode.H
				x2 = x
				y2 = y + stub
			case 3: // bottom
				x = wlk.startNode.X + space*float64(j) + offset
				y = wlk.startNode.Y
				x2 = x
				y2 = y - stub
			}

			path := &canvas.Path{}

			path.MoveTo(x, y)
			path.LineTo(x2, y2)

			wlk.pathGroups[i].paths[j] = path
		}
	}

	for !wlk.allOutOfBounds(ctx) {
		// for range 5 {

		var newGroups []*pathGroup

		for _, group := range wlk.pathGroups {

			movement := canvasWidth*0.05 + float64(wlk.RNG.Next(int64(canvasWidth*0.1)))
			group.move(movement)

			// split := int(wlk.RNG.Next(int64(len(group.paths))))
			// splitGroups := group.split(split)
			splitGroups := group.split(0)

			for _, g := range splitGroups {
				dirChange := int(wlk.RNG.Next(4))

				g.prevDirection = g.direction

				switch dirChange {
				case 1:
					if g.direction != "down" {
						g.direction = "up"
					}
				case 2:
					if g.direction != "up" {
						g.direction = "down"
					}
				case 3:
					if g.direction != "left" {
						g.direction = "right"
					}
				case 4:
					if g.direction != "right" {
						g.direction = "left"
					}
				}

			}

			newGroups = append(newGroups, splitGroups...)
		}

		wlk.pathGroups = newGroups

		log.Println(len(wlk.pathGroups), "path groups")
	}

	// wlk.pathGroups[1].move(100)
	// wlk.pathGroups[1].prevDirection = "right"
	// wlk.pathGroups[1].direction = "up"
	// wlk.pathGroups[1].move(100)

	// wlk.pathGroups[0].move(100)
	// wlk.pathGroups[0].prevDirection = "left"
	// wlk.pathGroups[0].direction = "up"
	// wlk.pathGroups[0].move(100)

	// wlk.pathGroups[2].move(100)
	// wlk.pathGroups[2].prevDirection = "up"
	// wlk.pathGroups[2].direction = "left"
	// wlk.pathGroups[2].move(100)

	// wlk.pathGroups[3].move(100)
	// wlk.pathGroups[3].prevDirection = "down"
	// wlk.pathGroups[3].direction = "left"
	// wlk.pathGroups[3].move(100)

	ctx.Push()
	ctx.SetFillColor(canvas.Transparent)

	wlk.drawStartNode(ctx)

	for _, group := range wlk.pathGroups {
		group.draw(ctx)
	}

	ctx.Pop()

}

func (wlk *CircuitPath) allOutOfBounds(ctx *canvas.Context) bool {

	for _, group := range wlk.pathGroups {
		for _, p := range group.paths {
			point := p.Pos()
			if point.X <= ctx.Width() && point.X > 0 && point.Y <= ctx.Height() && point.Y > 0 {
				return false
			}
		}
	}

	return true
}

func (wlk *CircuitPath) drawStartNode(ctx *canvas.Context) {

	ctx.SetStrokeColor(wlk.Color)
	ctx.SetStrokeWidth(wlk.PathWidth)

	ctx.DrawPath(0, 0, wlk.startNode.ToPath())

}
