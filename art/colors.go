package art

import (
	"image/color"
	"math"
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

func Darken(baseColor color.RGBA, factor float64) color.RGBA {
	r, g, b, a := baseColor.R, baseColor.G, baseColor.B, baseColor.A

	factorInt := 255 - int64(255.0*factor)

	return color.RGBA{
		R: uint8(math.Max(0, math.Min(float64(int64(r)-factorInt), 255))),
		G: uint8(math.Max(0, math.Min(float64(int64(g)-factorInt), 255))),
		B: uint8(math.Max(0, math.Min(float64(int64(b)-factorInt), 255))),
		A: uint8(a),
	}
}

func GetFactionBaseColor(factionID string) color.RGBA {

	switch factionID {
	case "shaper":
		return color.RGBA{
			R: 0x4c,
			G: 0xb1,
			B: 0x48,
			A: 0xff,
		}

	case "anarch":
		return color.RGBA{
			R: 0xe2,
			G: 0x6b,
			B: 0x35,
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
			R: 0xff,
			G: 0xde,
			B: 0x00,
			A: 0xff,
		}

	case "jinteki":
		return color.RGBA{
			R: 0x94,
			G: 0x2c,
			B: 0x4c,
			A: 0xff,
		}

	case "haas_bioroid":
		return color.RGBA{
			R: 0x5a,
			G: 0x32,
			B: 0x6d,
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
		R: 0xaa,
		G: 0xaa,
		B: 0xaa,
		A: 0xff,
	}

}
