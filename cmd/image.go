package cmd

import (
	"image"
	"log"
	"math"
	"os"
	"strings"

	"github.com/mangofeet/nrdb-go"
	"github.com/spf13/cobra"
	"github.com/tdewolff/canvas"
)

var imageCmd = &cobra.Command{
	Use:   "image [path to image] [card name or printing ID]",
	Args:  cobra.MinimumNArgs(2),
	Short: `Generate a frame for the named card over the specified image file`,
	Run: func(cmd *cobra.Command, args []string) {

		filename := args[0]

		cardName := strings.Join(args[1:], " ")

		if err := generateCardImage(filename, cardName); err != nil {
			log.Println("error:", err)
		}
	},
}

func generateCardImage(filename, cardName string) error {
	printing, err := getCardData(cardName)
	if err != nil {
		return err
	}
	log.Printf("generating %s", printing.Attributes.StrippedTitle)

	drawer := imageDrawer{
		filename: filename,
	}
	return generateCard(drawer, printing, "", designer)
}

type imageDrawer struct {
	filename string
}

func (drawer imageDrawer) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	file, err := os.Open(drawer.filename)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	imgWidth := img.Bounds().Max.X
	imgHeight := img.Bounds().Max.Y

	canvasWidth, canvasHeight := ctx.Size()

	widthScale := canvasWidth / float64(imgWidth)
	heightScale := canvasHeight / float64(imgHeight)

	scale := math.Max(widthScale, heightScale)

	ctx.RenderImage(img, canvas.Identity.Scale(scale, scale))

	return nil
}
