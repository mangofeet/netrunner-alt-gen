package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/netrunner-alt-gen/art/netspace"
	"github.com/mangofeet/netrunner-alt-gen/frame/basic"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

// old values

// const canvasWidth = 3288.0
// const canvasHeight = 4488.0
// const cardWidth = 3064.0
// const cardHeight = 4212.0

// based on NSG pdf sizes

// const canvasWidth = 3199.0
// const canvasHeight = 4432.0
// const cardWidth = 2975.0
// const cardHeight = 4156.0

// based on MPC template

const canvasWidth = 3264.0
const canvasHeight = 4450.0
const cardWidth = 2976.0
const cardHeight = 4152.0
const safeWidth = 2736.0
const safeHeight = 3924.0

// real card MM, doesn't work currently, need higehr res for
// generation

// const canvasWidth = 69.35
// const canvasHeight = 94.35
// const cardWidth = 63.0
// const cardHeight = 88.0

func main() {

	cardName := strings.Join(os.Args[1:], " ")

	log.Printf("generating %s", cardName)

	drawBleedLines := true
	// drawBleedLines := false

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

	var framer art.Drawer
	log.Println("card type:", printing.Attributes.CardTypeID)
	switch printing.Attributes.CardTypeID {
	case "program":
		framer = basic.FrameProgram{}
	case "resource":
		framer = basic.FrameResource{}
	case "hardware":
		framer = basic.FrameHardware{}
	case "event":
		framer = basic.FrameEvent{}
	case "ice":
		framer = basic.FrameIce{}
	case "asset":
		framer = basic.FrameAsset{}
	case "upgrade":
		framer = basic.FrameUpgrade{}
	case "operation":
		framer = basic.FrameOperation{}
	case "agenda":
		framer = basic.FrameAgenda{}
	case "runner_identity":
		framer = basic.FrameRunnerID{}
	case "corp_identity":
		framer = basic.FrameCorpID{}
	}

	if framer != nil {
		if err := framer.Draw(ctx, printing); err != nil {
			return err
		}
	}

	if drawBleedLines {

		marginX := (canvasWidth - cardWidth) / 2
		marginY := (canvasHeight - cardHeight) / 2
		safeMarginX := (canvasWidth - safeWidth) / 2
		safeMarginY := (canvasHeight - safeHeight) / 2

		drawMargin(ctx, marginX, marginY, cardWidth, cardHeight, color.White)
		drawMargin(ctx, safeMarginX, safeMarginY, safeWidth, safeHeight, canvas.Red)

	}

	if err := os.MkdirAll("output", os.ModePerm); err != nil {
		return err
	}
	if err := renderers.Write(fmt.Sprintf("output/%s.png", getFileName(printing)), cnv, canvas.DPMM(1)); err != nil {
		return err
	}

	return nil
}

var fileNameRegexp = regexp.MustCompile(`[^A-Za-z0-9]`)

func getFileName(card *nrdb.Printing) string {

	pos := fmt.Sprint(card.Attributes.PositionInSet)

	if card.Attributes.PositionInSet < 10 {
		pos = fmt.Sprintf("00%d", card.Attributes.PositionInSet)
	} else if card.Attributes.PositionInSet < 100 {
		pos = fmt.Sprintf("0%d", card.Attributes.PositionInSet)
	}

	title := fileNameRegexp.ReplaceAllString(card.Attributes.StrippedTitle, "-")
	title = strings.ToLower(title)

	set := fileNameRegexp.ReplaceAllString(card.Attributes.CardSetID, "-")

	return fmt.Sprintf("%s-%s-%s", set, pos, title)

}

func drawMargin(ctx *canvas.Context, x, y, w, h float64, c color.Color) {
	_, canvasHeight := ctx.Size()

	dash := canvasHeight * 0.0021

	ctx.Push()
	ctx.SetStrokeColor(c)
	ctx.SetStrokeWidth(dash / 2)
	ctx.MoveTo(x, y)
	ctx.LineTo(x, y+h)
	ctx.LineTo(x+w, y+h)
	ctx.LineTo(x+w, y)
	ctx.Close()
	ctx.SetDashes(0, dash, dash*2)
	ctx.Stroke()
	ctx.Pop()

}
