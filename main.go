package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/frame/netrunner"
	"github.com/mangofeet/nrdb-go"
)

// sizes in pixels, this is ~1200 DPI
const canvasWidth = 3288
const canvasHeight = 4488

const cardWidth = 3064
const cardHeight = 4212

func main() {

	cardName := strings.Join(os.Args[1:], " ")

	log.Printf("generating %s", cardName)

	drawBleedLines := true

	if err := generateCard(cardName, drawBleedLines); err != nil {
		log.Printf("error: %s", err)
	}

}

func generateCard(cardName string, drawBleedLines bool) error {

	nrClient := nrdb.NewClient()

	cards, err := nrClient.Cards(&nrdb.CardFilter{
		Search: &cardName,
	})
	if err != nil {
		return fmt.Errorf("getting card data: %w", err)
	}

	if len(cards) != 1 {
		for _, card := range cards {
			log.Printf("%s", card.Title())
		}
		return fmt.Errorf("mulitple results")
	}

	card := cards[0]

	printing, err := nrClient.Printing(card.LatestPrintingID())
	if err != nil {
		return fmt.Errorf("getting latest printing data: %w", err)
	}

	log.Printf("%s - %s", printing.Attributes.Title, printing.Attributes.Text)

	img := initImage(canvasWidth, canvasHeight)

	if drawBleedLines {
		marginLR := (canvasWidth - cardWidth) / 2
		marginTB := (canvasHeight - cardHeight) / 2
		drawBorder(img, marginLR, marginTB)
	}

	if err := drawArt(img, printing); err != nil {
		return fmt.Errorf("drawing art: %w", err)
	}

	// bgColor := getFactionBaseColor(printing.Attributes.FactionID)
	// bgColor.A = 0x44
	bgColor := color.RGBA{
		R: 0xbb,
		G: 0xbb,
		B: 0xbb,
		A: 0xcc,
	}

	if err := netrunner.DrawFrameProgram(img, printing, bgColor, color.Black); err != nil {
		return fmt.Errorf("drawing frame: %w", err)
	}

	outFile, err := os.Create("out.png")
	if err != nil {
		return fmt.Errorf("opening output file: %w", err)
	}

	if err := png.Encode(outFile, img); err != nil {
		return fmt.Errorf("encoding file: %w", err)
	}

	return nil
}

func initImage(width, height int) *image.CMYK {
	bg := image.Black
	cmyk := image.NewCMYK(image.Rect(0, 0, width, height))
	draw.Draw(cmyk, cmyk.Bounds(), bg, image.Point{}, draw.Src)

	return cmyk
}

func drawBorder(img draw.Image, marginLR, marginTB int) {
	ruler := color.RGBA{0xff, 0xff, 0xff, 0xff}

	bounds := img.Bounds()

	for i := 0; i < bounds.Max.X-(marginLR*2); i++ {
		img.Set(marginLR+i, marginTB, ruler)
		img.Set(marginLR+i, marginTB+1, ruler)
		img.Set(marginLR+i, bounds.Max.Y-marginTB, ruler)
		img.Set(marginLR+i, bounds.Max.Y-marginTB-1, ruler)
	}

	for i := 0; i < bounds.Max.Y-(marginTB*2); i++ {
		img.Set(marginLR, marginTB+i, ruler)
		img.Set(marginLR+1, marginTB+i, ruler)
		img.Set(bounds.Max.X-marginLR, marginTB+i, ruler)
		img.Set(bounds.Max.X-marginLR-1, marginTB+i, ruler)
	}
}
