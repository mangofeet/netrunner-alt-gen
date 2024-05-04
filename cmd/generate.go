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

func generateCard(drawer art.Drawer, card *nrdb.Printing, algorithm, designer string) error {

	cnv := canvas.New(canvasWidth, canvasHeight)

	ctx := canvas.NewContext(cnv)

	if err := drawer.Draw(ctx, card); err != nil {
		return err
	}

	var backCnv *canvas.Canvas
	if makeBack {
		backCnv = canvas.New(canvasWidth, canvasHeight)
		cnv.RenderTo(backCnv)
	}

	if err := output(cnv, ctx, card, algorithm, designer); err != nil {
		return err
	}

	if makeBack {
		frame = frame + "-back"
		output(backCnv, canvas.NewContext(backCnv), card, algorithm, designer)
	}

	return nil

}

func output(cnv *canvas.Canvas, ctx *canvas.Context, card *nrdb.Printing, algorithm, designer string) error {
	framer, err := getFramer(card, algorithm, designer)
	if err != nil {
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

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/%s.png", outputDir, getFileName(card))
	log.Printf("rendering output to %s", filename)
	if err := renderers.Write(filename, cnv, canvas.DPMM(1)); err != nil {
		return err
	}
	log.Println("done")

	return nil
}
