package cmd

import (
	"fmt"
	"image/color"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

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
			log.Printf("%s - %s", card.LatestPrintingID(), card.StrippedTitle())
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

var fileNameRegexp = regexp.MustCompile(`[^A-Za-z0-9]+`)

func getFileName(card *nrdb.Printing) string {

	pos := fmt.Sprint(card.Attributes.PositionInSet)

	if card.Attributes.PositionInSet < 10 {
		pos = fmt.Sprintf("00%d", card.Attributes.PositionInSet)
	} else if card.Attributes.PositionInSet < 100 {
		pos = fmt.Sprintf("0%d", card.Attributes.PositionInSet)
	}

	cardIDInt, err := strconv.Atoi(card.ID)
	if err != nil {
		panic(err)
	}
	cardID := fmt.Sprint(card.ID)

	if cardIDInt < 10 {
		cardID = fmt.Sprintf("00%s", card.ID)
	} else if cardIDInt < 100 {
		cardID = fmt.Sprintf("0%s", card.ID)
	}

	title := fileNameRegexp.ReplaceAllString(strings.ReplaceAll(card.Attributes.StrippedTitle, ":", ""), "-")
	title = strings.ToLower(title)

	set := fileNameRegexp.ReplaceAllString(card.Attributes.CardSetID, "-")

	return fmt.Sprintf("%s-%s-%s-%s", cardID, set, pos, title)

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
