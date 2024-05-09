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
	Color, ColorBG                             *color.RGBA
	AltColor1, AltColor2, AltColor3, AltColor4 *color.RGBA
}

func (drawer TechCircle) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	seed := card.Attributes.Title + card.Attributes.Text + card.Attributes.CardTypeID + card.Attributes.FactionID + card.Attributes.Flavor

	canvasWidth, canvasHeight := ctx.Size()

	rngGlobal := prng.NewGenerator(seed, nil)

	centerX := float64(rngGlobal.Next(int64(canvasWidth/2))) + (canvasWidth / 4)
	centerY := float64(rngGlobal.Next(int64(canvasHeight/6))) + ((canvasHeight / 8) * 5)
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

	ringer := art.TechRing{
		RNG:         rngGlobal,
		Angle:       angle,
		X:           centerX,
		Y:           centerY,
		Radius:      radius,
		RadiusStart: radiusStart,
		StrokeMin:   canvasHeight * 0.06,
		StrokeMax:   canvasHeight * 0.1,
		Color:       baseColor,
		AltColor1:   drawer.AltColor1,
		AltColor2:   drawer.AltColor2,
		AltColor3:   drawer.AltColor3,
		AltColor4:   drawer.AltColor4,
	}

	return ringer.Draw(ctx, card)

}
