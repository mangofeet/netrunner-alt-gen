package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art/reflection"
	"github.com/spf13/cobra"
)

var reflectionCmd = &cobra.Command{
	Use:   "reflection [card name or printing ID]",
	Args:  cobra.MinimumNArgs(1),
	Short: `Generate a card using the "reflection" algorithm`,
	Run: func(cmd *cobra.Command, args []string) {

		cardName := strings.Join(args, " ")

		if err := generateCardReflection(cardName); err != nil {
			log.Println("error:", err)
			os.Exit(1)
		}

	},
}

func generateCardReflection(cardName string) error {
	printing, err := getCardData(cardName)
	if err != nil {
		return err
	}
	log.Printf("generating %s", printing.Attributes.StrippedTitle)

	ns := reflection.Reflection{
		Color:   parseColor(baseColor),
		ColorBG: parseColor(colorBG),
	}

	return generateCard(ns, printing, "reflection", "mangofeet")
}
