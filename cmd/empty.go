package cmd

import (
	"log"
	"strings"

	"github.com/mangofeet/nrdb-go"
	"github.com/spf13/cobra"
	"github.com/tdewolff/canvas"
)

var emptyCmd = &cobra.Command{
	Use:   "empty [card name or printing ID]",
	Args:  cobra.MinimumNArgs(1),
	Short: `Generate a frame for the named card`,
	Run: func(cmd *cobra.Command, args []string) {

		cardName := strings.Join(args, " ")

		if err := generateCardEmpty(cardName); err != nil {
			log.Println("error:", err)
		}
	},
}

func generateCardEmpty(cardName string) error {
	printing, err := getCardData(cardName)
	if err != nil {
		return err
	}
	log.Printf("generating %s", printing.Attributes.StrippedTitle)

	drawer := emptyDrawer{}
	return generateCard(drawer, printing)
}

type emptyDrawer struct {
}

func (drawer emptyDrawer) Draw(ctx *canvas.Context, card *nrdb.Printing) error {
	return nil
}
