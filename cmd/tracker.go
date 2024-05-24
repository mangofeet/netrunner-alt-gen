package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art/tracker"
	"github.com/mangofeet/nrdb-go"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var trackerCmd = &cobra.Command{
	Use:   "tracker [card name or printing ID]",
	Args:  cobra.MinimumNArgs(1),
	Short: `Generate a card using the "tracker" algorithm`,
	Run: func(cmd *cobra.Command, args []string) {

		cardName := strings.Join(args, " ")

		if err := generateCardTracker(cardName); err != nil {
			log.Println("error:", err)
			os.Exit(1)
		}

	},
}

func generateCardTracker(cardName string) error {

	// construct a "card" to use for the printing
	printing := &nrdb.Printing{
		Document: nrdb.Document[nrdb.PrintingAttributes, nrdb.PrintingRelationships]{
			ID: "0",
			Attributes: &nrdb.PrintingAttributes{
				CardAttributes: nrdb.CardAttributes{
					Title:         cases.Title(language.English).String(cardName),
					FactionID:     cardName, // allows auto coloring things
					StrippedTitle: fmt.Sprintf("%s tracker", cardName),
				},
				CardSetID:     "components",
				PositionInSet: 0,
			},
		},
	}

	ns := tracker.Tracker{
		Color:      parseColor(baseColor),
		ColorBG:    parseColor(colorBG),
		RingColor1: parseColor(altColor1),
		RingColor2: parseColor(altColor2),
		RingColor3: parseColor(altColor3),
		RingColor4: parseColor(altColor4),
	}

	// set the frame to be "tracker" specifically
	if frame != "none" {
		frame = frame + "-tracker"
	}

	return generateCard(ns, printing, "tracker", "mangofeet")
}
