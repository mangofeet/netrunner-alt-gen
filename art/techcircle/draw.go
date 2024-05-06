package techcircle

import (
	"image/color"
	"log"
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

	radius := math.Max(canvasHeight-centerY, centerY) * 0.4

	circ := TechCircleDrawer{
		RNG:         rngGlobal,
		X:           centerX,
		Y:           centerY,
		Radius:      radius,
		RadiusStart: canvasHeight * 0.01,
		Color:       baseColor,
		StrokeMin:   canvasHeight * 0.03,
		StrokeMax:   canvasHeight * 0.05,
		GetColor:    getColor,
	}

	if err := circ.Draw(ctx); err != nil {
		return err
	}

	circOverlay := TechCircleDrawer{
		RNG:         rngGlobal,
		X:           centerX,
		Y:           centerY,
		Radius:      radius,
		RadiusStart: canvasHeight * 0.01,
		Color:       baseColor,
		StrokeMin:   canvasHeight * 0.005,
		StrokeMax:   canvasHeight * 0.01,
		GetColor: func(rng prng.Generator, base color.RGBA) (color.RGBA, error) {
			return color.RGBA{
				R: 0xff,
				G: 0xff,
				B: 0xff,
				A: 0x44,
			}, nil
		},
	}

	if err := circOverlay.Draw(ctx); err != nil {
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
	X, Y                float64
	Radius, RadiusStart float64
	StrokeMin           float64
	StrokeMax           float64
	GetColor            ColorGetter
}

func (drawer TechCircleDrawer) Draw(ctx *canvas.Context) error {
	radius := drawer.RadiusStart
	for radius < drawer.Radius {

		var lastPath *canvas.Path
		x, y := drawer.X+radius, drawer.Y
		strokeWidthRand := drawer.Radius*(float64(drawer.RNG.Next(5))/100.0) + 0.005
		strokeWidth := math.Min(drawer.StrokeMax, math.Min(math.Max(drawer.StrokeMin, strokeWidthRand), radius*0.5))
		color, err := drawer.GetColor(drawer.RNG, drawer.Color)
		if err != nil {
			return err
		}

		rot := 0.0
		arcPos := rot

		for arcPos < 360+rot {

			color := getColorOrBreak(drawer.RNG, color)

			path := &canvas.Path{}
			// bufferPath := &canvas.Path{}
			arc := math.Min(360-arcPos+rot, float64(drawer.RNG.Next(20)+5))
			log.Println("rmax:", drawer.Radius, "r:", radius, "arc:", arc, "arcPos:", arcPos, "stroke:", strokeWidth)

			theta0 := arcPos
			theta1 := math.Min(360+rot, arcPos+arc)

			// theta0Buffer := theta0 - 1
			// theta1Buffer := math.Min(360, theta1+1)

			// path.Arc(radius, radius, 2, arcPos-1, arcPos+arc+1)
			path.Arc(radius, radius, 0, theta0, theta1)
			// bufferPath.Arc(radius, radius, 2, theta0Buffer, theta1Buffer)

			arcPos += arc

			if lastPath != nil {
				x += lastPath.Pos().X
				y += lastPath.Pos().Y
			}

			// lastPath = bufferPath.Copy()
			lastPath = path.Copy()

			// log.Println(path)
			ctx.Push()
			ctx.SetStrokeColor(color)
			ctx.SetStrokeWidth(strokeWidth)
			ctx.SetFillColor(canvas.Transparent)
			ctx.DrawPath(x, y, path)
			ctx.Pop()
		}

		// radius += math.Min(math.Max(float64(drawer.RNG.Next(int64(radius))), strokeWidth*2), strokeWidth)
		radius += strokeWidth * (float64(drawer.RNG.Next(60))/100.0 + 0.7)

	}

	return nil

}

func getColorOrBreak(rng prng.Generator, base color.RGBA) color.RGBA {
	switch rng.Next(3) {
	case 1:
		return canvas.Transparent
	}

	return base
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
