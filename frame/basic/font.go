package basic

import (
	"image/color"
	"io"
	"log"

	"github.com/mangofeet/netrunner-alt-gen/assets"
	"github.com/tdewolff/canvas"
)

var fontFamily = canvas.NewFontFamily("cardtext")

func init() {

	if err := loadFont("Ubuntu-R.ttf", "sans-serif", canvas.FontRegular); err != nil {
		panic(err)
	}

	if err := loadFont("Ubuntu-B.ttf", "sans-serif", canvas.FontBold); err != nil {
		panic(err)
	}

	if err := loadFont("Ubuntu-RI.ttf", "sans-serif", canvas.FontItalic); err != nil {
		panic(err)
	}

	if err := loadFont("UbuntuMono-R.ttf", "monospace", canvas.FontBlack); err != nil {
		panic(err)
	}

	// the "extra bold" in the font family is used for unicode
	// symbols, this font seems to be the best at rendering them
	if err := loadFont("DejaVuSans.ttf", "monospace", canvas.FontExtraBold); err != nil {
		panic(err)
	}

}

func loadFont(name, backup string, style canvas.FontStyle) error {

	fontFile, err := assets.FS.Open(name)
	if err != nil {
		return err
	}

	fontFileBytes, err := io.ReadAll(fontFile)
	if err != nil {
		return err
	}

	if err := fontFamily.LoadFont(fontFileBytes, 0, style); err != nil {
		log.Printf(`could not load font "%s", trying "%s"`, name, backup)
		if err := fontFamily.LoadSystemFont(backup, style); err != nil {
			return err
		}
	}
	return nil
}

func (fb FrameBasic) getFont(size float64, style canvas.FontStyle) *canvas.FontFace {
	return fontFamily.Face(size, fb.getColorText(), style)
}

func (fb FrameBasic) getFontWithColor(size float64, style canvas.FontStyle, clr color.Color) *canvas.FontFace {
	return fontFamily.Face(size, clr, style)
}
