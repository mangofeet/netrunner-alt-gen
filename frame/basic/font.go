package basic

import (
	"log"

	"github.com/tdewolff/canvas"
)

var fontFamily = canvas.NewFontFamily("cardtext")

func init() {

	if err := loadFont("Ubuntu", "sans-serif", canvas.FontRegular); err != nil {
		panic(err)
	}

	if err := loadFont("Ubuntu", "sans-serif", canvas.FontBold); err != nil {
		panic(err)
	}

	if err := loadFont("Ubuntu", "sans-serif", canvas.FontItalic); err != nil {
		panic(err)
	}

	if err := loadFont("Ubuntu Mono", "monospace", canvas.FontBlack); err != nil {
		panic(err)
	}

	// the "extra bold" in the font family is used for unicode
	// symbols, this font seems to be the best at rendering them
	if err := loadFont("DejaVu Sans", "monospace", canvas.FontExtraBold); err != nil {
		panic(err)
	}

}

func loadFont(name, backup string, style canvas.FontStyle) error {
	if err := fontFamily.LoadSystemFont(name, style); err != nil {
		log.Printf(`could not find font "%s", trying "%s"`, name, backup)
		if err := fontFamily.LoadSystemFont(backup, style); err != nil {
			return err
		}
	}
	return nil
}

func getFont(size float64, style canvas.FontStyle) *canvas.FontFace {
	return fontFamily.Face(size, textColor, style)
}
