package techcircle

import (
	"image/color"
	"math"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type TechCircle struct {
	Color, ColorBG *color.RGBA
}

func (drawer TechCircle) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	seed := card.Attributes.Title + card.Attributes.Text + card.Attributes.CardTypeID + card.Attributes.FactionID + card.Attributes.Flavor

	canvasWidth, canvasHeight := ctx.Size()

	rngGlobal := prng.NewGenerator(seed, nil)

	centerX := float64(rngGlobal.Next(int64(canvasWidth/2))) + (canvasWidth / 4)
	centerY := float64(rngGlobal.Next(int64(canvasHeight/6))) + ((canvasHeight / 8) * 5)
	// radius := float64(rngGlobal.Next(int64(canvasHeight/8))) + (canvasHeight / 12)

	if card.Attributes.CardTypeID == "ice" {
		centerY = float64(rngGlobal.Next(int64(canvasHeight/4))) + (canvasHeight / 6)
	}

	baseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	if drawer.Color != nil {
		baseColor = *drawer.Color
	}

	cardBGColor := art.Darken(baseColor, 0.623)
	if drawer.ColorBG != nil {
		cardBGColor = *drawer.ColorBG
	}

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

	radius := math.Max(canvasHeight-centerY, centerY) * 1.5
	angle := 0.0

	radiusStart := canvasHeight * 0.03

	circ := TechCircleDrawer{
		RNG:         rngGlobal,
		X:           centerX,
		Y:           centerY,
		Radius:      radius,
		RadiusStart: radiusStart,
		Color:       baseColor,
		StrokeMin:   canvasHeight * 0.06,
		StrokeMax:   canvasHeight * 0.1,
		GetColor:    getColor,
		Angle:       angle,
	}

	if err := circ.Draw(ctx); err != nil {
		return err
	}

	circOverlay := TechCircleDrawer{
		RNG:         rngGlobal,
		X:           centerX,
		Y:           centerY,
		Radius:      radius,
		RadiusStart: radiusStart,
		Color:       baseColor,
		StrokeMin:   canvasHeight * 0.01,
		StrokeMax:   canvasHeight * 0.03,
		GetColor: func(rng prng.Generator, base color.RGBA) (color.RGBA, error) {
			return color.RGBA{
				R: 0xff,
				G: 0xff,
				B: 0xff,
				A: 0x44,
			}, nil
		},
		Angle:         angle,
		SegmentArcMin: 8,
		SegmentArcMax: 15,
	}

	if err := circOverlay.Draw(ctx); err != nil {
		return err
	}

	circBlanker := TechCircleDrawer{
		RNG:         rngGlobal,
		X:           centerX,
		Y:           centerY,
		Radius:      radius,
		RadiusStart: radiusStart,
		Color:       baseColor,
		StrokeMin:   canvasHeight * 0.01,
		StrokeMax:   canvasHeight * 0.03,
		GetColor: func(rng prng.Generator, base color.RGBA) (color.RGBA, error) {
			return cardBGColor, nil
		},
		Angle:         angle,
		SegmentArcMin: 2,
		SegmentArcMax: 5,
	}

	if err := circBlanker.Draw(ctx); err != nil {
		return err
	}

	// ctx.Push()
	// ctx.SetFillColor(baseColor)
	// ctx.DrawPath(centerX, centerY, canvas.Circle(radius*0.1))
	// ctx.Pop()

	return nil
}

type ColorGetter func(rng prng.Generator, base color.RGBA) (color.RGBA, error)

type TechCircleDrawer struct {
	RNG                 prng.Generator
	Color               color.RGBA
	Angle               float64
	X, Y                float64
	Radius, RadiusStart float64
	StrokeMin           float64
	StrokeMax           float64
	GetColor            ColorGetter

	SegmentArcMin, SegmentArcMax float64
	BreakArcMin, BreakArcMax     float64
}

type circleSegment struct {
	start, end  float64
	strokeWidth float64
	strokeColor color.RGBA
	isBlank     bool
}

type circleRing struct {
	segments []circleSegment
	radius   float64
	rotation float64
}

func (drawer TechCircleDrawer) Draw(ctx *canvas.Context) error {
	radius := drawer.RadiusStart

	var rings []circleRing

	segArcMin := drawer.SegmentArcMin
	if segArcMin == 0 {
		segArcMin = 5
	}
	segArcMax := drawer.SegmentArcMax
	if segArcMax == 0 {
		segArcMax = 25
	}

	breakArcMin := drawer.BreakArcMin
	if breakArcMin == 0 {
		breakArcMin = 5
	}
	breakArcMax := drawer.BreakArcMax
	if breakArcMax == 0 {
		breakArcMax = 25
	}

	for radius < drawer.Radius {

		ring := circleRing{
			radius: radius,
		}

		strokeWidthRand := drawer.Radius*(float64(drawer.RNG.Next(5))/100.0) + 0.005
		strokeWidth := math.Min(drawer.StrokeMax, math.Min(math.Max(drawer.StrokeMin, strokeWidthRand), radius*0.5))
		thisColor, err := drawer.GetColor(drawer.RNG, drawer.Color)
		if err != nil {
			return err
		}

		rot := float64(drawer.RNG.Next(90) - 45)
		if drawer.Angle != 0 {
			rot = drawer.Angle
		}

		arcPos := rot

		maxArc := 360 + rot

		for arcPos < maxArc {

			var isBreak bool
			var arcStart = arcPos

			for {

				_, isBreak = getColorOrBreak(drawer.RNG, thisColor)

				var arc float64
				if isBreak {
					arc = math.Min(maxArc-arcPos, float64(drawer.RNG.Next(int64(breakArcMax)-int64(breakArcMin))+int64(breakArcMin)))
				} else {
					arc = math.Min(maxArc-arcPos, float64(drawer.RNG.Next(int64(segArcMax)-int64(segArcMin))+int64(segArcMin)))
				}
				// log.Println("r:", radius, "arc:", arc, "arcPos:", arcPos, "stroke:", strokeWidth, "break", isBreak)

				if isBreak {
					if arcStart != arcPos {
						segment := circleSegment{
							start:       arcStart,
							end:         math.Min(maxArc, arcPos),
							strokeWidth: strokeWidth,
							// strokeWidth: 10,
							strokeColor: thisColor,
						}
						ring.segments = append(ring.segments, segment)
						// log.Printf("segment: %#v", segment)
					}
					spacer := circleSegment{
						start:       arcPos,
						end:         math.Min(maxArc, arcPos+arc),
						strokeWidth: strokeWidth * 0.5,
						isBlank:     true,
						strokeColor: color.RGBA{
							R: 0xff,
							G: 0xff,
							B: 0xff,
							A: 0x33,
						},
					}
					ring.segments = append(ring.segments, spacer)
					arcPos += arc
					break
				}

				arcPos += arc
				// bufferPath := &canvas.Path{}

			}

		}

		// radius += math.Min(math.Max(float64(drawer.RNG.Next(int64(radius))), strokeWidth*2), strokeWidth)
		radius += strokeWidth * (float64(drawer.RNG.Next(60))/100.0 + 0.7)

		ring.rotation = rot
		rings = append(rings, ring)
	}

	for _, ring := range reverse(rings) {

		x, y := drawer.X+ring.radius, drawer.Y

		if drawer.Angle == 0 {
			rotPath := &canvas.Path{}
			rotPath.Arc(ring.radius, ring.radius, 0.1, 0, ring.rotation)
			x += rotPath.Pos().X
			y += rotPath.Pos().Y
		}

		for _, seg := range ring.segments {
			// log.Printf("%#v", seg)
			path := &canvas.Path{}
			path.Arc(ring.radius, ring.radius, 0.1, seg.start, seg.end)

			if !seg.isBlank {
				ctx.Push()
				ctx.SetStrokeColor(seg.strokeColor)
				ctx.SetStrokeWidth(seg.strokeWidth)
				ctx.SetFillColor(canvas.Transparent)
				ctx.DrawPath(x, y, path)
				ctx.Pop()
			}

			x += path.Pos().X
			y += path.Pos().Y
		}

	}

	return nil

}

func reverse[T any](slc []T) []T {
	reversed := make([]T, len(slc))
	for i := range len(slc) {
		ri := len(slc) - 1 - i
		reversed[ri] = slc[i]
	}
	return reversed
}

func getColorOrBreak(rng prng.Generator, base color.RGBA) (color.RGBA, bool) {
	switch rng.Next(3) {
	case 1:
		return canvas.Transparent, true
	}

	return base, false
}

func getColor(rng prng.Generator, base color.RGBA) (color.RGBA, error) {

	var err error
	newColor := base

	switch rng.Next(5) {
	case 1:
		// newColor = color.RGBA{
		// 	R: 0xff,
		// 	G: 0xff,
		// 	B: 0xff,
		// 	A: 0x77,
		// }
	case 2:
		newColor, _, err = art.Analogous(base, float64(rng.Next(80)-40))
		if err != nil {
			return base, err
		}
	case 3:
		newColor, _, err = art.Analogous(base, float64(rng.Next(100)-50))
		if err != nil {
			return base, err
		}
	case 4:
		newColor, _, err = art.Analogous(base, float64(rng.Next(120)-60))
		if err != nil {
			return base, err
		}
	case 5:
		newColor, _, err = art.Analogous(base, float64(rng.Next(140)-70))
		if err != nil {
			return base, err
		}
	case 6:
		// color = canvas.Transparent
	}

	return newColor, nil
}
