package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art/netwalker"
	"github.com/spf13/cobra"
)

var netwalkerCmd = &cobra.Command{
	Use:   "netwalker [card name or printing ID]",
	Args:  cobra.MinimumNArgs(1),
	Short: `Generate a card using the "netwalker" algorithm`,
	Run: func(cmd *cobra.Command, args []string) {

		cardName := strings.Join(args, " ")

		if err := generateCardNetwalker(cardName); err != nil {
			log.Println("error:", err)
			os.Exit(1)
		}

	},
}

func generateCardNetwalker(cardName string) error {
	printing, err := getCardData(cardName)
	if err != nil {
		return err
	}
	log.Printf("generating %s", printing.Attributes.StrippedTitle)

	var nGridP *float64
	if gridPercent >= 0 {
		nGridP = &gridPercent
	}

	ns := netwalker.NetWalker{
		MinWalkers:   walkersMin,
		MaxWalkers:   walkersMax,
		GridPercent:  nGridP,
		Color:        parseColor(baseColor),
		ColorBG:      parseColor(colorBG),
		WalkerColor1: parseColor(walkerColor1),
		WalkerColor2: parseColor(walkerColor2),
		WalkerColor3: parseColor(walkerColor3),
		WalkerColor4: parseColor(walkerColor4),
		GridColor1:   parseColor(gridColor1),
		GridColor2:   parseColor(gridColor2),
		GridColor3:   parseColor(gridColor3),
		GridColor4:   parseColor(gridColor4),
	}

	return generateCard(ns, printing, "netwalker", "mangofeet")
}
