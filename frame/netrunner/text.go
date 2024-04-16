package netrunner

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type drawFonts struct {
	regular, bold, symbol             *truetype.Font
	faceRegular, faceBold, faceSymbol font.Face
}

func drawCardText(textBoxHeight, textBoxWidth int, cardText string, img draw.Image, textColor color.Color, fonts *drawFonts) error {

	fontDPI := float64(72)
	bounds := img.Bounds()

	fg := image.NewUniform(textColor)

	fontSizeText := float64(textBoxHeight) * 0.08
	spacing := float64(fontSizeText / 2)

	if fonts.regular == nil {
		return fmt.Errorf("at least regular font required")
	}

	fonts.faceRegular = truetype.NewFace(fonts.regular, &truetype.Options{
		Size:    fontSizeText,
		DPI:     fontDPI,
		Hinting: font.HintingFull,
	})

	if fonts.bold == nil {
		fonts.bold = fonts.regular
	}

	fonts.faceBold = truetype.NewFace(fonts.bold, &truetype.Options{
		Size:    fontSizeText,
		DPI:     fontDPI,
		Hinting: font.HintingFull,
	})

	if fonts.symbol == nil {
		fonts.symbol = fonts.regular
	}

	fonts.faceSymbol = truetype.NewFace(fonts.symbol, &truetype.Options{
		Size:    fontSizeText,
		DPI:     fontDPI,
		Hinting: font.HintingFull,
	})

	text := strings.Split(cardText, "\n")

	drawerText := &font.Drawer{
		Dst:  img,
		Src:  fg,
		Face: fonts.faceRegular,
	}

	textTextX := fixed.I(bounds.Max.X - textBoxWidth + int(spacing))
	textTextY := fixed.I(bounds.Max.Y - textBoxHeight + int(fontSizeText) + int(spacing))
	maxLineWidth := fixed.I(textBoxWidth - int(spacing*2))

	dy := int(math.Ceil((fontSizeText + spacing) * (fontDPI / 72)))
	log.Print(fontSizeText, spacing, fontDPI, fontSizeText*spacing*fontDPI, dy)

	isBold := false
	isSybmol := false

	for _, ln := range text {

		words := strings.Split(ln, " ")

		for len(words) > 0 {

			firstWord := words[0]
			firstWord = strings.ReplaceAll(firstWord, "<strong>", "")
			firstWord = strings.ReplaceAll(firstWord, "</strong>", "")

			thisLine := []string{firstWord}

			switch true {
			case isBold:
				drawerText.Face = fonts.faceBold
			case isSybmol:
				drawerText.Face = fonts.faceSymbol
			default:
				drawerText.Face = fonts.faceRegular
			}

			for _, word := range words[1:] {
				// if !isBold && strings.Contains(word, "<strong>") {
				// 	isBold = true
				// 	break
				// }

				// if isBold && strings.Contains(word, "</strong>") {
				// 	isBold = false
				// }

				word = strings.ReplaceAll(word, "<strong>", "")
				word = strings.ReplaceAll(word, "</strong>", "")

				lineText := strings.Join(thisLine, " ")

				if drawerText.MeasureString(lineText+" "+word) <= maxLineWidth {
					thisLine = append(thisLine, word)
				}
			}

			words = words[len(thisLine):]

			drawerText.Dot = fixed.Point26_6{
				X: textTextX,
				Y: textTextY,
			}

			log.Printf("printing line at %s", drawerText.Dot)
			drawerText.DrawString(strings.Join(thisLine, " "))

			textTextY += fixed.I(dy)
		}

		textTextY += fixed.I(int(spacing))

	}

	return nil
}
