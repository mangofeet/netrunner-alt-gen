package basic

import (
	"image/color"
)

var colorDefaultBG = color.RGBA{
	R: 0x1c,
	G: 0x1c,
	B: 0x1c,
	A: 0x99,
}

var colorDefaultOpaqueBG = color.RGBA{
	R: 0x1c,
	G: 0x1c,
	B: 0x1c,
	A: 0xff,
}

var colorDefaultText = color.RGBA{
	R: 0xdc,
	G: 0xdc,
	B: 0xcc,
	A: 0xff,
}

var transparent = color.RGBA{
	R: 0,
	G: 0,
	B: 0,
	A: 0,
}
