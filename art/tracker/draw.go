package tracker

import (
	"image/color"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/netrunner-alt-gen/internal/prng"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type Tracker struct {
	Color, ColorBG                                 *color.RGBA
	RingColor1, RingColor2, RingColor3, RingColor4 *color.RGBA
}

func (drawer Tracker) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	seed := card.Attributes.Title

	canvasWidth, canvasHeight := ctx.Size()

	startX := canvasWidth / 2
	startY := canvasHeight / 2

	baseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	if drawer.Color != nil {
		baseColor = *drawer.Color
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

	ringRadius := canvasHeight * 1.2
	ringRadiusStart := canvasWidth * 0.35
	ringStrokeMin := canvasWidth * 0.04
	ringStrokeMax := canvasWidth * 0.08

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

	ringSequence := int64(2)
	// rings under the walkers
	ringSequence++
	(art.TechRing{
		RNG:          prng.NewGenerator(seed, &ringSequence),
		X:            startX,
		Y:            startY,
		Radius:       ringRadius,
		RadiusStart:  ringRadiusStart,
		StrokeMin:    ringStrokeMin,
		StrokeMax:    ringStrokeMax,
		Color:        ringColor,
		AltColor1:    ringColor1,
		AltColor2:    ringColor2,
		AltColor3:    ringColor3,
		AltColor4:    ringColor4,
		OverlayColor: &canvas.Transparent,
	}).Draw(ctx)

	ringSequence++
	(art.TechRing{
		RNG:          prng.NewGenerator(seed, &ringSequence),
		X:            startX,
		Y:            startY,
		Radius:       ringRadius,
		RadiusStart:  ringRadiusStart * 2,
		StrokeMin:    ringStrokeMin,
		StrokeMax:    ringStrokeMax,
		Color:        ringColor,
		AltColor1:    ringColor1,
		AltColor2:    ringColor2,
		AltColor3:    ringColor3,
		AltColor4:    ringColor4,
		OverlayColor: &canvas.Transparent,
	}).Draw(ctx)

	// rings over the walkers
	ringSequence++
	(art.TechRing{
		RNG:          prng.NewGenerator(seed, &ringSequence),
		X:            startX,
		Y:            startY,
		Radius:       ringRadius,
		RadiusStart:  ringRadiusStart * 5,
		StrokeMin:    ringStrokeMin,
		StrokeMax:    ringStrokeMax,
		Color:        ringColor,
		AltColor1:    &canvas.Transparent,
		AltColor2:    ringColor2,
		AltColor3:    ringColor3,
		AltColor4:    ringColor4,
		OverlayColor: &canvas.Transparent,
	}).Draw(ctx)

	ringSequence++
	(art.TechRing{
		RNG:          prng.NewGenerator(seed, &ringSequence),
		X:            startX,
		Y:            startY,
		Radius:       ringRadius,
		RadiusStart:  ringRadiusStart,
		StrokeMin:    ringStrokeMin,
		StrokeMax:    ringStrokeMax,
		Color:        overlayRingColor,
		AltColor1:    &overlayRingColor,
		AltColor2:    &overlayRingColor,
		AltColor3:    &overlayRingColor,
		AltColor4:    &overlayRingColor,
		OverlayColor: &overlayRingColor,
	}).Draw(ctx)

	return nil
}
