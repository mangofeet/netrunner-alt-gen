package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art/netspace"
	"github.com/mangofeet/netrunner-alt-gen/frame/netrunner"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

// sizes in pixels, this is ~1200 DPI
const canvasWidth = 3288.0
const canvasHeight = 4488.0

const cardWidth = 3064.0
const cardHeight = 4212.0

// const canvasWidth = 69.35
// const canvasHeight = 94.35

// const cardWidth = 63.0
// const cardHeight = 88.0

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

	cnv := canvas.New(canvasWidth, canvasHeight)

	ctx := canvas.NewContext(cnv)

	if err := netspace.Draw(ctx, printing); err != nil {
		return err
	}

	if err := netrunner.DrawFrameProgram(ctx, printing); err != nil {
		return err
	}

	if drawBleedLines {

		marginX := (canvasWidth - cardWidth) / 2
		marginY := (canvasHeight - cardHeight) / 2

		ctx.Push()
		ctx.SetStrokeColor(color.White)
		ctx.SetStrokeWidth(5)
		ctx.MoveTo(marginX, marginY)
		ctx.LineTo(marginX, marginY+cardHeight)
		ctx.LineTo(marginX+cardWidth, marginY+cardHeight)
		ctx.LineTo(marginX+cardWidth, marginY)
		ctx.Close()
		ctx.SetDashes(0, 10, 20)
		ctx.Stroke()
		ctx.Pop()

	}

	// if err := renderers.Write("out.png", cnv, canvas.DPI(1200)); err != nil {
	if err := renderers.Write("out.png", cnv, canvas.DPMM(1)); err != nil {
		return err
	}

	return nil
}
