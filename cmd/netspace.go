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

	var nGridP *float64
	if gridPercent >= 0 {
		nGridP = &gridPercent
	}

	ns := netspace.Netspace{
		MinWalkers:   netspaceWalkersMin,
		MaxWalkers:   netspaceWalkersMax,
		GridPercent:  nGridP,
		Color:        parseColor(baseColor),
		ColorBG:      parseColor(netspaceColorBG),
		WalkerColor1: parseColor(walkerColor1),
		WalkerColor2: parseColor(walkerColor2),
		WalkerColor3: parseColor(walkerColor3),
		WalkerColor4: parseColor(walkerColor4),
		GridColor1:   parseColor(gridColor1),
		GridColor2:   parseColor(gridColor2),
		GridColor3:   parseColor(gridColor3),
		GridColor4:   parseColor(gridColor4),
	}

	return generateCard(ns, printing)
}
