package art

import "image/color"

func GetFactionBaseColor(factionID string) color.RGBA {

	switch factionID {
	case "shaper", "weyland_consortium":
		return color.RGBA{
			R: 0x7f,
			G: 0x9f,
			B: 0x7f,
			A: 0xff,
		}

	case "anarch":
		return color.RGBA{
			R: 0xdf,
			G: 0xaf,
			B: 0x8f,
			A: 0xff,
		}

	case "criminal":
		return color.RGBA{
			R: 0x8c,
			G: 0xd0,
			B: 0xd3,
			A: 0xff,
		}

	case "nbn":
		return color.RGBA{
			R: 0xf0,
			G: 0xdf,
			B: 0xaf,
			A: 0xff,
		}

	case "jinteki":
		return color.RGBA{
			R: 0xcc,
			G: 0x93,
			B: 0x96,
			A: 0xff,
		}

	case "haas_bioroid":
		return color.RGBA{
			R: 0xc0,
			G: 0xbe,
			B: 0xd1,
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
