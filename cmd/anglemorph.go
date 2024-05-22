package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art/anglemorph"
	"github.com/spf13/cobra"
)

var anglemorphCmd = &cobra.Command{
	Use:   "anglemorph [card name or printing ID]",
	Args:  cobra.MinimumNArgs(1),
	Short: `Generate a card using the "anglemorph" algorithm`,
	Run: func(cmd *cobra.Command, args []string) {

		cardName := strings.Join(args, " ")

		if err := generateCardAnglemorph(cardName); err != nil {
			log.Println("error:", err)
			os.Exit(1)
		}

	},
}

func generateCardAnglemorph(cardName string) error {
	printing, err := getCardData(cardName)
	if err != nil {
		return err
	}
	log.Printf("generating %s", printing.Attributes.StrippedTitle)

	ns := anglemorph.AngleMorph{
		ColumnCount: 60,
		RowCount:    90,
		Color:       parseColor(baseColor),
		ColorBG:     parseColor(colorBG),
	}

	return generateCard(ns, printing, "anglemorph", "mangofeet")
}
