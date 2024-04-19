package basic

import "github.com/tdewolff/canvas"

var fontFamily = canvas.NewFontFamily("cardtext")

func init() {
	if err := fontFamily.LoadSystemFont("Ubuntu", canvas.FontRegular); err != nil {
		panic(err)
	}

	if err := fontFamily.LoadSystemFont("Ubuntu", canvas.FontBold); err != nil {
		panic(err)
	}

	if err := fontFamily.LoadSystemFont("Ubuntu", canvas.FontItalic); err != nil {
		panic(err)
	}

	if err := fontFamily.LoadSystemFont("UbuntuMono Nerd Font", canvas.FontBlack); err != nil {
		panic(err)
	}

	if err := fontFamily.LoadSystemFont("DejaVu Sans", canvas.FontExtraBold); err != nil {
		panic(err)
	}
}

func getFont(size float64, style canvas.FontStyle) *canvas.FontFace {
	return fontFamily.Face(size, textColor, style)
}
