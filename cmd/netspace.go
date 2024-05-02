package cmd

import (
	"log"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art/netspace"
	"github.com/spf13/cobra"
)

var netspaceCmd = &cobra.Command{
	Use:   "netspace [card name or printing ID]",
	Args:  cobra.MinimumNArgs(1),
	Short: `Generate a card using the "netspace" algorithm`,
	Run: func(cmd *cobra.Command, args []string) {

		cardName := strings.Join(args, " ")

		if err := generateCardNetspace(cardName); err != nil {
			log.Println("error:", err)
		}
	},
}

func generateCardNetspace(cardName string) error {
	printing, err := getCardData(cardName)
	if err != nil {
		return err
	}
	log.Printf("generating %s", printing.Attributes.StrippedTitle)

	ns := netspace.Netspace{
		MinWalkers: netspaceWalkersMin,
		MaxWalkers: netspaceWalkersMax,
		Color:      parseColor(baseColor),
		AltColor1:  parseColor(altColor1),
		AltColor2:  parseColor(altColor2),
		AltColor3:  parseColor(altColor3),
		GridColor1: parseColor(gridColor1),
		GridColor2: parseColor(gridColor2),
		GridColor3: parseColor(gridColor3),
	}

	return generateCard(ns, printing)
}
