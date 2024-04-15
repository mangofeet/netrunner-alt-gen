package netrunner

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/mangofeet/nrdb-go"
)

func DrawFrameProgram(img draw.Image, card *nrdb.Printing, bgColor, textColor color.Color) error {

	bg := image.NewUniform(bgColor)

	bounds := img.Bounds()

	titleBounds := image.Rect(0, 0, bounds.Max.X-(bounds.Max.X/3), bounds.Max.Y/12)
	textBounds := image.Rect(bounds.Max.Y/6, (bounds.Max.Y/3)*2, bounds.Max.X, bounds.Max.Y)

	draw.Draw(img, titleBounds, bg, image.Point{}, draw.Over)
	draw.Draw(img, textBounds, bg, image.Point{}, draw.Over)

	return nil

}
