package art

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
)

type TechRing struct {
	RNG                                        prng.Generator
	Angle                                      float64
	X, Y                                       float64
	Radius, RadiusStart                        float64
	StrokeMin                                  float64
	StrokeMax                                  float64
	Color                                      color.RGBA
	AltColor1, AltColor2, AltColor3, AltColor4 *color.RGBA
	OverlayColor                               *color.RGBA
}

func (drawer TechRing) Draw(ctx *canvas.Context, card *nrdb.Printing) error {
	canvasWidth, canvasHeight := ctx.Size()

	ringCnv := canvas.New(canvasWidth, canvasHeight)
	ringCtx := canvas.NewContext(ringCnv)

	circ := techCircleDrawer{
		RNG:         drawer.RNG,
		X:           drawer.X,
		Y:           drawer.Y,
		Radius:      drawer.Radius,
		RadiusStart: drawer.RadiusStart,
		StrokeMin:   drawer.StrokeMin,
		StrokeMax:   drawer.StrokeMax,
		GetColor:    drawer.getColor(drawer.Color),
		Angle:       drawer.Angle,
	}

	if err := circ.Draw(ringCtx); err != nil {
		return err
	}

	overlayColor := color.RGBA{
		R: 0xff,
		G: 0xff,
		B: 0xff,
		A: 0x44,
	}
	if drawer.OverlayColor != nil {
		overlayColor = *drawer.OverlayColor
	}

	circOverlay := techCircleDrawer{
		RNG:         drawer.RNG,
		X:           drawer.X,
		Y:           drawer.Y,
		Radius:      drawer.Radius,
		RadiusStart: drawer.RadiusStart,
		StrokeMin:   drawer.StrokeMin * 0.16666667,
		StrokeMax:   drawer.StrokeMax * 0.3,
		GetColor: func(rng prng.Generator) (color.Color, error) {
			return overlayColor, nil
		},
		Angle:         drawer.Angle,
		SegmentArcMin: 8,
		SegmentArcMax: 15,
	}

	if err := circOverlay.Draw(ringCtx); err != nil {
		return err
	}

	ringImg := rasterizer.Draw(ringCnv, canvas.DPMM(1), canvas.DefaultColorSpace)

	maskCnv := canvas.New(canvasWidth, canvasHeight)
	maskCtx := canvas.NewContext(maskCnv)

	circBlanker := techCircleDrawer{
		RNG:         drawer.RNG,
		X:           drawer.X,
		Y:           drawer.Y,
		Radius:      drawer.Radius,
		RadiusStart: drawer.RadiusStart,
		StrokeMin:   drawer.StrokeMin * 0.16666667,
		StrokeMax:   drawer.StrokeMax * 0.3,
		GetColor: func(rng prng.Generator) (color.Color, error) {
			return canvas.Black, nil
		},
		Angle:         drawer.Angle,
		SegmentArcMin: 2,
		SegmentArcMax: 5,
	}

	if err := circBlanker.Draw(maskCtx); err != nil {
		return err
	}

	maskImg := rasterizer.Draw(maskCnv, canvas.DPMM(1), canvas.DefaultColorSpace)

	// invert the mask image
	for i, pxl := range maskImg.Pix {
		if pxl == 0 {
			maskImg.Pix[i] = 255
		} else {
			maskImg.Pix[i] = 0
		}
	}

	ringsFinal := image.NewRGBA(image.Rect(0, 0, int(canvasWidth), int(canvasHeight)))
	draw.DrawMask(ringsFinal, ringsFinal.Bounds(), ringImg, image.Point{}, maskImg, image.Point{}, draw.Over)

	ctx.RenderImage(ringsFinal, canvas.Identity)

	return nil
}

type colorGetter func(rng prng.Generator) (color.Color, error)

type techCircleDrawer struct {
	RNG                 prng.Generator
	Angle               float64
	X, Y                float64
	Radius, RadiusStart float64
	StrokeMin           float64
	StrokeMax           float64
	GetColor            colorGetter

	SegmentArcMin, SegmentArcMax float64
	BreakArcMin, BreakArcMax     float64
}

type circleSegment struct {
	start, end  float64
	strokeWidth float64
	strokeColor color.Color
	isBlank     bool
}

type circleRing struct {
	segments []circleSegment
	radius   float64
	rotation float64
}

func (drawer techCircleDrawer) Draw(ctx *canvas.Context) error {
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
		thisColor, err := drawer.GetColor(drawer.RNG)
		if err != nil {
			return err
		}

		// darkFactor := drawer.Radius / (drawer.Radius - radius) * 0.7
		// thisColor, _ = Desaturate(thisColor, darkFactor)
		// darkFactor := drawer.Radius / (drawer.Radius - radius) * 0.5
		// thisColor, _ = AdjustLevel(thisColor, darkFactor)

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

func getColorOrBreak(rng prng.Generator, base color.Color) (color.Color, bool) {
	switch rng.Next(3) {
	case 1:
		return canvas.Transparent, true
	}

	return base, false
}

func (drawer TechRing) getColor(base color.RGBA) colorGetter {

	return func(rng prng.Generator) (color.Color, error) {

		var err error
		newColor := base

		switch rng.Next(4) {
		case 1:
			newColor, _, err = Analogous(base, float64(rng.Next(80)-40))
			if err != nil {
				return base, err
			}
			if drawer.AltColor1 != nil {
				return *drawer.AltColor1, nil
			}
		case 2:
			newColor, _, err = Analogous(base, float64(rng.Next(100)-50))
			if err != nil {
				return base, err
			}
			if drawer.AltColor2 != nil {
				return *drawer.AltColor2, nil
			}
		case 3:
			newColor, _, err = Analogous(base, float64(rng.Next(120)-60))
			if err != nil {
				return base, err
			}
			if drawer.AltColor3 != nil {
				return *drawer.AltColor3, nil
			}
		case 4:
			newColor, _, err = Analogous(base, float64(rng.Next(140)-70))
			if err != nil {
				return base, err
			}
			if drawer.AltColor4 != nil {
				return *drawer.AltColor4, nil
			}
		}

		return newColor, nil
	}
}
