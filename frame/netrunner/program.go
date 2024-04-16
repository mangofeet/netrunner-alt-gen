package netrunner

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/mangofeet/nrdb-go"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func DrawFrameProgram(img draw.Image, card *nrdb.Printing, bgColor, textColor color.Color) error {

	bg := image.NewUniform(bgColor)
	fg := image.NewUniform(textColor)

	bounds := img.Bounds()

	titleBoxHeight := bounds.Max.Y / 12
	titleBoxWidth := bounds.Max.X - (bounds.Max.X / 4)

	textBoxHeight := (bounds.Max.Y / 3)
	textBoxWidth := (bounds.Max.X / 6) * 5

	titleBounds := image.Rect(0, 0, titleBoxWidth, titleBoxHeight)
	textBounds := image.Rect(bounds.Max.X-textBoxWidth, bounds.Max.Y-textBoxHeight, bounds.Max.X, bounds.Max.Y)

	draw.Draw(img, titleBounds, bg, image.Point{}, draw.Over)
	draw.Draw(img, textBounds, bg, image.Point{}, draw.Over)

	fontBytesRegular, err := os.ReadFile("/usr/share/fonts/nerd-fonts-complete/TTF/Ubuntu Light Nerd Font Complete.ttf")
	if err != nil {
		return fmt.Errorf("reading font file: %w", err)
	}
	fontBytesBold, err := os.ReadFile("/usr/share/fonts/nerd-fonts-complete/TTF/Ubuntu Bold Nerd Font Complete.ttf")
	if err != nil {
		return fmt.Errorf("reading font file: %w", err)
	}
	fontBytesSymbol, err := os.ReadFile("/usr/share/fonts/TTF/InputMonoNerdFont-Regular.ttf")
	if err != nil {
		return fmt.Errorf("reading font file: %w", err)
	}

	fontFaceRegular, err := truetype.Parse(fontBytesRegular)
	if err != nil {
		return fmt.Errorf("parsing font file: %w", err)
	}

	fontFaceBold, err := truetype.Parse(fontBytesBold)
	if err != nil {
		return fmt.Errorf("parsing font file: %w", err)
	}

	fontFaceSymbol, err := truetype.Parse(fontBytesSymbol)
	if err != nil {
		return fmt.Errorf("parsing font file: %w", err)
	}

	fontSizeTitle := float64(titleBoxHeight) * 0.75
	fontDPI := float64(72)

	drawerTitle := &font.Drawer{
		Dst: img,
		Src: fg,
		Face: truetype.NewFace(fontFaceRegular, &truetype.Options{
			Size:    fontSizeTitle,
			DPI:     fontDPI,
			Hinting: font.HintingFull,
		}),
	}

	titleTextX := fixed.I(titleBoxWidth) - drawerTitle.MeasureString(card.Attributes.Title) - fixed.I(titleBoxWidth/24)
	titleTextY := fixed.I(titleBoxHeight - (titleBoxHeight / 4))

	drawerTitle.Dot = fixed.Point26_6{
		X: titleTextX,
		Y: titleTextY,
	}

	drawerTitle.DrawString(card.Attributes.Title)

	drawCardText(textBoxHeight, textBoxWidth, card.Attributes.Text, img, textColor, &drawFonts{
		regular: fontFaceRegular,
		bold:    fontFaceBold,
		symbol:  fontFaceSymbol,
	})

	return nil

}
