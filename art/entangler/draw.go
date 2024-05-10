package entangler

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/mangofeet/nrdb-go"
	"github.com/ojrac/opensimplex-go"
	"github.com/tdewolff/canvas"
)

type Entangler struct {
	MinWalkers, MaxWalkers                                 int
	GridPercent                                            *float64
	Color, ColorBG                                         *color.RGBA
	WalkerColor1, WalkerColor2, WalkerColor3, WalkerColor4 *color.RGBA
	GridColor1, GridColor2, GridColor3, GridColor4         *color.RGBA
	RingColor1, RingColor2, RingColor3, RingColor4         *color.RGBA
}

func (drawer Entangler) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	seed := card.Attributes.Title + card.Attributes.Text + card.Attributes.CardTypeID + card.Attributes.FactionID + card.Attributes.Flavor

	canvasWidth, canvasHeight := ctx.Size()

	rngGlobal := prng.NewGenerator(seed, nil)

	numWalkers := int(math.Max(float64(drawer.MinWalkers), float64(rngGlobal.Next(int64(drawer.MaxWalkers)))))

	startX := rngGlobal.Next(int64(canvasWidth/2)) + int64(canvasWidth/4)
	startY := rngGlobal.Next(int64(canvasHeight/6)) + (int64(canvasHeight/8) * 5)

	if card.Attributes.CardTypeID == "ice" {
		startY = rngGlobal.Next(int64(canvasHeight/4)) + (int64(canvasHeight / 6))
	}

	baseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	if drawer.Color != nil {
		baseColor = *drawer.Color
	}
	walkerColor1, walkerColor2, err := art.Analogous(baseColor, 10+float64(rngGlobal.Next(20)))
	if err != nil {
		return fmt.Errorf("getting analogous colors: %w", err)
	}
	if drawer.WalkerColor1 != nil {
		walkerColor1 = *drawer.WalkerColor1
	}
	if drawer.WalkerColor2 != nil {
		walkerColor2 = *drawer.WalkerColor2
	}

	walkerColor3, walkerColor4, err := art.Analogous(baseColor, 30+float64(rngGlobal.Next(20)))
	if err != nil {
		return fmt.Errorf("getting third analog: %w", err)
	}
	if drawer.WalkerColor3 != nil {
		walkerColor3 = *drawer.WalkerColor3
	}
	if drawer.WalkerColor4 != nil {
		walkerColor4 = *drawer.WalkerColor4
	}

	gridColor1, gridColor2, gridColor3, gridColor4 := baseColor, baseColor, baseColor, baseColor
	if drawer.GridColor1 != nil {
		gridColor1 = *drawer.GridColor1
	}
	if drawer.GridColor2 != nil {
		gridColor2 = *drawer.GridColor2
	}
	if drawer.GridColor3 != nil {
		gridColor3 = *drawer.GridColor3
	}
	if drawer.GridColor4 != nil {
		gridColor4 = *drawer.GridColor4
	}

	cardBGColor := art.Darken(baseColor, 0.623)
	if drawer.ColorBG != nil {
		cardBGColor = *drawer.ColorBG
	}

	ringColor := color.RGBA{
		R: baseColor.R,
		G: baseColor.G,
		B: baseColor.B,
		A: 0xdd,
	}
	overlayRingColor := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x22}

	var ringColor1 *color.RGBA
	if drawer.RingColor1 != nil {
		ringColor1 = drawer.RingColor1
	}
	var ringColor2 *color.RGBA
	if drawer.RingColor2 != nil {
		ringColor2 = drawer.RingColor2
	}
	var ringColor3 *color.RGBA
	if drawer.RingColor3 != nil {
		ringColor3 = drawer.RingColor3
	}
	var ringColor4 *color.RGBA
	if drawer.RingColor4 != nil {
		ringColor4 = drawer.RingColor4
	}

	ringRadius := canvasWidth * 1.4
	ringRadiusStart := canvasWidth * 0.01
	ringStrokeMin := canvasWidth * 0.1
	ringStrokeMax := canvasWidth * 0.21

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

	var walkers []*art.Walker

	nGrid := 0.0
	if float64(numWalkers)*0.01 >= 1 {
		nGrid = math.Max(float64(numWalkers)*0.01, float64(rngGlobal.Next(int64(float64(numWalkers)*0.02))))
	}
	if drawer.GridPercent != nil {
		nGrid = float64(numWalkers) * *drawer.GridPercent
	}

	dirChangeStep := 45.0
	// dirChangeStep := float64(rngGlobal.Next(15) + 40)

	// do manual seeds for these with high numbers so they didn't
	// affect the walkers
	var altColorDirection1 string
	// switch prng.Sample(seed, 9999999, 4) {
	switch rngGlobal.Next(4) {
	case 1:
		altColorDirection1 = "up"
	case 2:
		altColorDirection1 = "down"
	case 3:
		altColorDirection1 = "left"
	case 4:
		altColorDirection1 = "right"
	}

	altColorDirection2 := altColorDirection1
	for altColorDirection2 == altColorDirection1 {
		// switch prng.Sample(seed, i+10000000, 4) {
		switch rngGlobal.Next(4) {
		case 1:
			altColorDirection2 = "up"
		case 2:
			altColorDirection2 = "down"
		case 3:
			altColorDirection2 = "left"
		case 4:
			altColorDirection2 = "right"
		}
	}

	altColorDirection3 := altColorDirection1
	for altColorDirection3 == altColorDirection1 || altColorDirection3 == altColorDirection2 {
		// switch prng.Sample(seed, i+10000000, 4) {
		switch rngGlobal.Next(4) {
		case 1:
			altColorDirection3 = "up"
		case 2:
			altColorDirection3 = "down"
		case 3:
			altColorDirection3 = "left"
		case 4:
			altColorDirection3 = "right"
		}
	}

	altColorDirection4 := altColorDirection1
	for altColorDirection4 == altColorDirection1 || altColorDirection4 == altColorDirection2 || altColorDirection4 == altColorDirection3 {
		// switch prng.Sample(seed, i+10000000, 4) {
		switch rngGlobal.Next(4) {
		case 1:
			altColorDirection4 = "up"
		case 2:
			altColorDirection4 = "down"
		case 3:
			altColorDirection4 = "left"
		case 4:
			altColorDirection4 = "right"
		}
	}

	var sequence int64
	for i := 0; i < numWalkers; i++ {

		colorFactor := rngGlobal.Next(128) - 64

		var direction string
		var grid = false
		var strokeWidth = 0.3

		thisColor := baseColor

		thisStartX := startX
		thisStartY := startY

		if float64(i) < nGrid {
			colorFactor = -2 * int64(math.Abs(float64(colorFactor)))
			grid = true
			strokeWidth = 1.5
			switch rngGlobal.Next(4) {
			case 1:
				thisColor = gridColor1
			case 2:
				thisColor = gridColor2
			case 3:
				thisColor = gridColor3
			case 4:
				thisColor = gridColor4
			}
			thisColor, err = art.Desaturate(thisColor, float64(colorFactor)/-128.0)
			if err != nil {
				return err
			}
		} else {

			dirSeed := rngGlobal.Next(4)

			if dirSeed <= 1 {
				direction = "up"
				// thisStartY -= int64(ringRadiusStart)
			} else if dirSeed <= 2 {
				direction = "down"
				// thisStartY += int64(ringRadiusStart)
			} else if dirSeed <= 3 {
				direction = "left"
				// thisStartX += int64(ringRadiusStart)
			} else if dirSeed <= 4 {
				direction = "right"
				// thisStartX -= int64(ringRadiusStart)
			}

			switch direction {
			case altColorDirection1:
				thisColor = walkerColor1
			case altColorDirection2:
				thisColor = walkerColor2
			case altColorDirection3:
				thisColor = walkerColor3
			case altColorDirection4:
				thisColor = walkerColor4
			}

			switch rngGlobal.Next(8) {
			case 1:
				thisColor = walkerColor1
			case 2:
				thisColor = walkerColor2
			case 3:
				thisColor = walkerColor3
			case 4:
				thisColor = walkerColor4
			}

		}

		sequence = int64(i)

		wlk := art.Walker{
			RNG:                         prng.NewGenerator(seed, &sequence),
			Direction:                   direction,
			DirectionVariance:           2,
			DirectionChangeStep:         dirChangeStep,
			DirectionChangeStepModifier: 1.5,
			X:                           float64(thisStartX),
			Y:                           float64(thisStartY),
			Vx:                          (float64(rngGlobal.Next(100)) / 100) - 0.5,
			Vy:                          (float64(rngGlobal.Next(100)) / 100) - 0.5,
			Color: color.RGBA{
				R: uint8(math.Max(0, math.Min(float64(int64(thisColor.R)+colorFactor), 255))),
				G: uint8(math.Max(0, math.Min(float64(int64(thisColor.G)+colorFactor), 255))),
				B: uint8(math.Max(0, math.Min(float64(int64(thisColor.B)+colorFactor), 255))),
				A: 0xff,
			},
			Noise:       noise,
			Grid:        grid,
			StrokeWidth: strokeWidth,
		}
		walkers = append(walkers, &wlk)
	}

	// rings under the walkers
	sequence++
	(art.TechRing{
		RNG:         prng.NewGenerator(seed, &sequence),
		X:           float64(startX),
		Y:           float64(startY),
		Radius:      ringRadius,
		RadiusStart: ringRadiusStart,
		StrokeMin:   ringStrokeMin,
		StrokeMax:   ringStrokeMax,
		Color:       ringColor,
		AltColor1:   ringColor1,
		AltColor2:   ringColor2,
		AltColor3:   ringColor3,
		AltColor4:   ringColor4,
		// AltColor1:    &canvas.Red,
		// AltColor2:    &canvas.Red,
		// AltColor3:    &canvas.Red,
		// AltColor4:    &canvas.Red,
		OverlayColor: &canvas.Transparent,
	}).Draw(ctx, card)

	for i, wlk := range walkers {
		wlk.Draw(ctx)
		if i == (len(walkers)/4)*3 {
			sequence++
			(art.TechRing{
				RNG:         prng.NewGenerator(seed, &sequence),
				X:           float64(startX),
				Y:           float64(startY),
				Radius:      ringRadius,
				RadiusStart: ringRadiusStart * 10,
				StrokeMin:   ringStrokeMin,
				StrokeMax:   ringStrokeMax,
				Color:       ringColor,
				AltColor1:   ringColor1,
				AltColor2:   ringColor2,
				AltColor3:   ringColor3,
				AltColor4:   ringColor4,
				// AltColor1:    &canvas.Green,
				// AltColor2:    &canvas.Green,
				// AltColor3:    &canvas.Green,
				// AltColor4:    &canvas.Green,
				OverlayColor: &canvas.Transparent,
			}).Draw(ctx, card)

		}
		for wlk.InBounds(ctx) {
			wlk.Velocity()
			wlk.Move()
			wlk.Draw(ctx)
		}
	}

	// rings over the walkers
	sequence++
	(art.TechRing{
		RNG:         prng.NewGenerator(seed, &sequence),
		X:           float64(startX),
		Y:           float64(startY),
		Radius:      ringRadius,
		RadiusStart: ringRadiusStart * 20,
		StrokeMin:   ringStrokeMin,
		StrokeMax:   ringStrokeMax,
		Color:       ringColor,
		AltColor1:   ringColor1,
		AltColor2:   ringColor2,
		AltColor3:   ringColor3,
		AltColor4:   ringColor4,
		// AltColor1:    &canvas.Blue,
		// AltColor2:    &canvas.Blue,
		// AltColor3:    &canvas.Blue,
		// AltColor4:    &canvas.Blue,
		OverlayColor: &canvas.Transparent,
	}).Draw(ctx, card)

	sequence++
	(art.TechRing{
		RNG:          prng.NewGenerator(seed, &sequence),
		X:            float64(startX),
		Y:            float64(startY),
		Radius:       canvasWidth * 0.5,
		RadiusStart:  canvasWidth * 0.1,
		StrokeMin:    canvasWidth * 0.1,
		StrokeMax:    canvasWidth * 0.2,
		Color:        overlayRingColor,
		AltColor1:    &overlayRingColor,
		AltColor2:    &overlayRingColor,
		AltColor3:    &overlayRingColor,
		AltColor4:    &overlayRingColor,
		OverlayColor: &overlayRingColor,
	}).Draw(ctx, card)

	log.Printf("finished %d walkers", len(walkers))

	return nil
}
