package art

import (
	"image/color"
	"math"

	"github.com/crazy3lf/colorconv"
)

func Lighten(baseColor color.RGBA, factor float64) color.RGBA {
	r, g, b, a := baseColor.R, baseColor.G, baseColor.B, baseColor.A

	factorInt := 255 - int64(255.0*factor)

	return color.RGBA{
		R: uint8(math.Max(0, math.Min(float64(int64(r)+factorInt), 255))),
		G: uint8(math.Max(0, math.Min(float64(int64(g)+factorInt), 255))),
		B: uint8(math.Max(0, math.Min(float64(int64(b)+factorInt), 255))),
		A: uint8(a),
	}
}

func Darken(baseColor color.Color, factor float64) color.RGBA {
	r, g, b, a := baseColor.RGBA()

	factorInt := 255 - int64(255.0*factor)

	return color.RGBA{
		R: uint8(math.Max(0, math.Min(float64(int64(r/257)-factorInt), 255))),
		G: uint8(math.Max(0, math.Min(float64(int64(g/257)-factorInt), 255))),
		B: uint8(math.Max(0, math.Min(float64(int64(b/257)-factorInt), 255))),
		A: uint8(a),
	}
}

func Complementary(baseColor color.RGBA) color.RGBA {
	r, g, b, a := baseColor.R, baseColor.G, baseColor.B, baseColor.A

	return color.RGBA{
		R: 255 - r,
		G: 255 - g,
		B: 255 - b,
		A: a,
	}

}

func Analogous(baseColor color.RGBA, degShift float64) (color.RGBA, color.RGBA, error) {
	h, s, l := colorconv.ColorToHSL(baseColor)

	h1 := math.Mod(h+degShift, 360)
	h2 := math.Mod(h-degShift, 360)

	if h1 < 0 {
		h1 = 360 + h1
	}
	if h2 < 0 {
		h2 = 360 + h2
	}

	r1, g1, b1, err := colorconv.HSLToRGB(h1, s, l)
	if err != nil {
		return color.RGBA{}, color.RGBA{}, err
	}
	r2, g2, b2, err := colorconv.HSLToRGB(h2, s, l)
	if err != nil {
		return color.RGBA{}, color.RGBA{}, err
	}

	c1 := color.RGBA{r1, g1, b1, baseColor.A}
	c2 := color.RGBA{r2, g2, b2, baseColor.A}

	return c1, c2, nil

}

func Desaturate(baseColor color.Color, amount float64) (color.RGBA, error) {
	h, s, l := colorconv.ColorToHSL(baseColor)

	r, g, b, err := colorconv.HSLToRGB(h, math.Min(s*amount, 1), l)
	if err != nil {
		return color.RGBA{}, err
	}
	_, _, _, baseA := baseColor.RGBA()
	c1 := color.RGBA{r, g, b, uint8(baseA / 257)}

	return c1, nil
}

func AdjustLevel(baseColor color.Color, amount float64) (color.RGBA, error) {
	h, s, l := colorconv.ColorToHSL(baseColor)

	r, g, b, err := colorconv.HSLToRGB(h, s, math.Min(l*amount, 1))
	if err != nil {
		return color.RGBA{}, err
	}
	_, _, _, baseA := baseColor.RGBA()
	c1 := color.RGBA{r, g, b, uint8(baseA / 257)}

	return c1, nil
}

func GetFactionBaseColor(factionID string) color.RGBA {

	switch factionID {
	case "shaper":
		return color.RGBA{
			R: 0x15,
			G: 0x73,
			B: 0x0a,
			A: 0xff,
		}

	case "anarch":
		return color.RGBA{
			R: 0x8f,
			G: 0x3f,
			B: 0x06,
			A: 0xff,
		}

	case "criminal":
		return color.RGBA{
			R: 0x19,
			G: 0x4c,
			B: 0x9b,
			A: 0xff,
		}

	case "nbn":
		return color.RGBA{
			R: 0x8f,
			G: 0x6a,
			B: 0x06,
			A: 0xff,
		}

	case "jinteki":
		return color.RGBA{
			R: 0x8f,
			G: 0x07,
			B: 0x15,
			A: 0xff,
		}

	case "haas_bioroid":
		return color.RGBA{
			R: 0x52,
			G: 0x23,
			B: 0x69,
			A: 0xff,
		}

	case "weyland_consortium":

		return color.RGBA{
			R: 0x3f,
			G: 0x4f,
			B: 0x3f,
			A: 0xff,
		}

	}

	return color.RGBA{
		R: 0x3f,
		G: 0x3f,
		B: 0x3f,
		A: 0xff,
	}

}
