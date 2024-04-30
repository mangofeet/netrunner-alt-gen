package cmd

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/netrunner-alt-gen/frame/basic"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

func generateCard(drawer art.Drawer, printing *nrdb.Printing) error {

	cnv := canvas.New(canvasWidth, canvasHeight)

	ctx := canvas.NewContext(cnv)

	if err := drawer.Draw(ctx, printing); err != nil {
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

	filename := fmt.Sprintf("output/%s.png", getFileName(printing))
	log.Printf("rendering output to %s", filename)
	if err := renderers.Write(filename, cnv, canvas.DPMM(1)); err != nil {
		return err
	}
	log.Println("done")

	return nil

}
