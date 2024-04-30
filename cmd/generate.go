package cmd

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

func generateCard(drawer art.Drawer, card *nrdb.Printing) error {

	cnv := canvas.New(canvasWidth, canvasHeight)

	ctx := canvas.NewContext(cnv)

	framer, err := getFramer(card)
	if err != nil {
		return err
	}

	if err := drawer.Draw(ctx, card); err != nil {
		return err
	}

	if err := framer.Draw(ctx, card); err != nil {
		return err
	}

	if drawMarginLines {

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

	filename := fmt.Sprintf("output/%s.png", getFileName(card))
	log.Printf("rendering output to %s", filename)
	if err := renderers.Write(filename, cnv, canvas.DPMM(1)); err != nil {
		return err
	}
	log.Println("done")

	return nil

}
