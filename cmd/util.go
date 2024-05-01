package cmd

import (
	"fmt"
	"image/color"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/netrunner-alt-gen/frame/basic"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func getFramer(card *nrdb.Printing) (art.Drawer, error) {

	switch frame {
	case "basic":
		var framer art.Drawer
		log.Println("card type:", card.Attributes.CardTypeID)

		if textBoxFactor > 1 {
			textBoxFactor /= 100.0
		}

		frm := basic.FrameBasic{
			Color:               parseColor(baseColor),
			TextBoxHeightFactor: &textBoxFactor,
		}

		if flavorText != "" {
			frm.Flavor = "<em>" + flavorText + "</em>"
		}
		if flavorAttribution != "" {
			frm.FlavorAttribution = "<em>- " + flavorAttribution + "</em>"
		}

		switch card.Attributes.CardTypeID {
		case "program":
			framer = frm.Program()
		case "resource":
			framer = frm.Resource()
		case "hardware":
			framer = frm.Hardware()
		case "event":
			framer = frm.Event()
		case "ice":
			framer = frm.Ice()
		case "asset":
			framer = frm.Asset()
		case "upgrade":
			framer = frm.Upgrade()
		case "operation":
			framer = frm.Operation()
		case "agenda":
			framer = frm.Agenda()
		case "runner_identity":
			framer = frm.RunnerID()
		case "corp_identity":
			framer = frm.CorpID()
		default:
			return nil, fmt.Errorf(`unknown card type "%s"`, card.Attributes.CardTypeID)
		}

		return framer, nil
	case "none":
		return art.NoopDrawer{}, nil
	}

	return nil, fmt.Errorf(`unknown frame type "%s"`, frame)
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
		return nil, fmt.Errorf("no results for %s", cardName)
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

func parseColor(colorStr string) *color.RGBA {

	if colorStr == "" {
		return nil
	}

	if colorStr[0] == '#' {
		colorStr = colorStr[1:]
	}

	if len(colorStr) != 6 && len(colorStr) != 8 {
		log.Printf(`!! Could not parse color string "#%s", using defaults`, colorStr)
		return nil
	}

	rStr := colorStr[0:2]
	gStr := colorStr[2:4]
	bStr := colorStr[4:6]
	aStr := "ff"
	if len(colorStr) == 8 {
		aStr = colorStr[6:8]
	}

	r, err := strconv.ParseInt(rStr, 16, 16)
	if err != nil {
		log.Printf(`!! Could not parse color string "#%s", using defaults`, colorStr)
		return nil
	}

	g, err := strconv.ParseInt(gStr, 16, 16)
	if err != nil {
		log.Printf(`!! Could not parse color string "#%s", using defaults`, colorStr)
		return nil
	}

	b, err := strconv.ParseInt(bStr, 16, 16)
	if err != nil {
		log.Printf(`!! Could not parse color string "#%s", using defaults`, colorStr)
		return nil
	}

	a, err := strconv.ParseInt(aStr, 16, 16)
	if err != nil {
		log.Printf(`!! Could not parse color string "#%s", using defaults`, colorStr)
		return nil
	}

	return &color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}

}
