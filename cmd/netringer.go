package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art/techcircle"
	"github.com/spf13/cobra"
)

var netringerCmd = &cobra.Command{
	Use:   "netringer [card name or printing ID]",
	Args:  cobra.MinimumNArgs(1),
	Short: `Generate a card using the "netringer" algorithm`,
	Run: func(cmd *cobra.Command, args []string) {

		cardName := strings.Join(args, " ")

		if err := generateCardNetringer(cardName); err != nil {
			log.Println("error:", err)
			os.Exit(1)
		}

	},
}

func generateCardNetringer(cardName string) error {
	printing, err := getCardData(cardName)
	if err != nil {
		return err
	}
	log.Printf("generating %s", printing.Attributes.StrippedTitle)

	ns := techcircle.TechCircle{
		Color:     parseColor(baseColor),
		ColorBG:   parseColor(colorBG),
		AltColor1: parseColor(altColor1),
		AltColor2: parseColor(altColor2),
		AltColor3: parseColor(altColor3),
		AltColor4: parseColor(altColor4),
	}

	return generateCard(ns, printing, "netringer", "mangofeet")
}
