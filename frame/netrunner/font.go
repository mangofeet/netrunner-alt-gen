package netrunner

import "github.com/tdewolff/canvas"

var fontFamily = canvas.NewFontFamily("cardtext")

func init() {
	if err := fontFamily.LoadFontFile("/usr/share/fonts/TTF/UbuntuNerdFont-Regular.ttf", canvas.FontRegular); err != nil {
		panic(err)
	}

	if err := fontFamily.LoadFontFile("/usr/share/fonts/TTF/UbuntuNerdFont-Bold.ttf", canvas.FontBold); err != nil {
		panic(err)
	}

	if err := fontFamily.LoadFontFile("/usr/share/fonts/TTF/UbuntuMonoNerdFont-Regular.ttf", canvas.FontBlack); err != nil {
		panic(err)
	}
}

func getFont(size float64, style canvas.FontStyle) *canvas.FontFace {
	return fontFamily.Face(size, textColor, style, canvas.FontNormal)
}
