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

		cardName := strings.Join(args, "")

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
	}

	return generateCard(ns, printing)
}
