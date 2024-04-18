package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
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

func getCardData(cardName string) (*nrdb.Printing, error) {
	nrClient := nrdb.NewClient()

	if printingID, err := strconv.Atoi(cardName); err == nil {
		printing, err := nrClient.Printing(cardName)
		if err != nil {
			return nil, fmt.Errorf("no result for printing ID %d", printingID)
		}
		return printing, nil
	}

	cards, err := nrClient.Cards(&nrdb.CardFilter{
		Search: &cardName,
	})
	if err != nil {
		return nil, fmt.Errorf("getting card data: %w", err)
	}

	if len(cards) == 0 {
		return nil, fmt.Errorf("no results")
	}

	if len(cards) != 1 {
		for _, card := range cards {
			log.Printf("%s - %s", card.Title(), card.LatestPrintingID())
		}
		return nil, fmt.Errorf("mulitple results")
	}

	card := cards[0]

	printing, err := nrClient.Printing(card.LatestPrintingID())
	if err != nil {
		return nil, fmt.Errorf("getting latest printing data: %w", err)
	}

	return printing, nil

}

func generateCard(cardName string, drawBleedLines bool) error {

	printing, err := getCardData(cardName)
	if err != nil {
		return err
	}

	log.Printf("%s - %s", printing.Attributes.Title, printing.Attributes.Text)

	cnv := canvas.New(canvasWidth, canvasHeight)

	ctx := canvas.NewContext(cnv)

	if err := netspace.Draw(ctx, printing); err != nil {
		return err
	}

	switch printing.Attributes.CardTypeID {
	case "program":
		if err := netrunner.DrawFrameProgram(ctx, printing); err != nil {
			return err
		}
	case "resource":
		if err := netrunner.DrawFrameResource(ctx, printing); err != nil {
			return err
		}
	}

	if drawBleedLines {

		marginX := (canvasWidth - cardWidth) / 2
		marginY := (canvasHeight - cardHeight) / 2

		dash := canvasHeight * 0.0021

		ctx.Push()
		ctx.SetStrokeColor(color.White)
		ctx.SetStrokeWidth(dash / 2)
		ctx.MoveTo(marginX, marginY)
		ctx.LineTo(marginX, marginY+cardHeight)
		ctx.LineTo(marginX+cardWidth, marginY+cardHeight)
		ctx.LineTo(marginX+cardWidth, marginY)
		ctx.Close()
		ctx.SetDashes(0, dash, dash*2)
		ctx.Stroke()
		ctx.Pop()

	}

	// if err := renderers.Write("out.png", cnv, canvas.DPI(1200)); err != nil {
	if err := renderers.Write("out.png", cnv, canvas.DPMM(1)); err != nil {
		return err
	}

	return nil
}
