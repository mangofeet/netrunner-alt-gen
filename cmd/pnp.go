package cmd

import (
	"encoding/csv"
	"github.com/mangofeet/nrdb-go"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
)

var pnpCmd = &cobra.Command{
	Use:   "pnp [CSV file]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Generate a print & play file containing cards from a CSV",
	Run: func(cmd *cobra.Command, args []string) {
		if err := generatePnPFile(args[0]); err != nil {
			log.Println("error:", err)
			os.Exit(1)
		}
	},
}

func generatePnPFile(csvPath string) error {
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	log.Printf("Generating print & play file from %s", csvPath)

	baseColor = "ffffff"
	frameColorBackground = "ffffff"
	frameColorFactionBG = "ffffff"
	frameColorStrengthBG = "ffffff"
	frameColorInfluenceLimitBG = "ffffff"
	frameColorMinDeckBG = "ffffff"
	frameColorBorder = "000000"
	frameColorText = "000000"

	cardID := startRow
	for _, record := range records[startRow - 1:] {
		// Instantiate Printing
		card := &nrdb.Printing{}
		card.Attributes = &nrdb.PrintingAttributes{}
		card.Attributes.CardAbilities = &nrdb.CardAbilities{}

		titleStripper := strings.NewReplacer(" ", "", ".", "", ",", "", "-", "", "!", "", "◆", "")

		// Split summary into sections for ease of use
		summary_sections := strings.Split(record[5], "====")
		lower_sections := summary_sections[1]
		summary_sections = []string{summary_sections[0]}
		summary_sections = append(summary_sections, strings.Split(lower_sections, "----")...)
		for i, section := range summary_sections {
			summary_sections[i] = strings.Trim(section, "\n")
		}

		// Set ID, faction, name, and type
		card.ID = strconv.Itoa(cardID)
		card.Attributes.PositionInSet = cardID
		if record[1] == "Weyland" {
			card.Attributes.FactionID = "weyland_consortium"
		} else {
			card.Attributes.FactionID = strings.ReplaceAll(strings.ToLower(record[1]), "-", "_")
		}
		// card.Attributes.Title = record[3]
		card.Attributes.Title = record[3]
		card.Attributes.StrippedTitle = titleStripper.Replace(record[3])
		card.Attributes.IsUnique = strings.Contains(summary_sections[0], "◆")
		if record[4] == "Runner-ID" {
			card.Attributes.CardTypeID = "runner_identity"
			card.Attributes.CardAbilities.MUProvided = nil
		} else if record[4] == "Corp-ID" {
			card.Attributes.CardTypeID = "corp_identity"
		} else {
			card.Attributes.CardTypeID = strings.ToLower(record[4])
		}

		// Set subtypes
		type_line := strings.Split(summary_sections[1], " ")
		subtypes := []string{}
		for _, token := range type_line {
			if !strings.HasSuffix(token, ":") && token != "-" {
				subtypes = append(subtypes, token)
			}
		}
		card.Attributes.DisplaySubtypes = new(string)
		*card.Attributes.DisplaySubtypes = strings.Join(subtypes, "-")

		// Set numeric values
		cost_line := strings.Split(summary_sections[2], ", ")
		for _, cost := range cost_line {
			tokens := strings.Split(cost, ": ")
			val, _ := strconv.Atoi(tokens[1])
			switch tokens[0] {
			case "Advancements":
				card.Attributes.AdvancementRequirement = &val
			case "Cost":
				card.Attributes.Cost = &val
			case "Deck":
				card.Attributes.MinimumDeckSize = &val
			case "Influence":
				card.Attributes.InfluenceCost = &val
				card.Attributes.InfluenceLimit = &val
			case "Memory":
				card.Attributes.MemoryCost = &val
			case "Points":
				card.Attributes.AgendaPoints = &val
			case "Strength":
				card.Attributes.Strength = &val
			case "Trash":
				card.Attributes.MinimumDeckSize = &val
			}
		}

		// Set text
		card.Attributes.Text = strings.Trim(summary_sections[3], "\n")

		err := generateCard(emptyDrawer{}, card, "", "")
		if err != nil {
			return err
		}

		cardID += 1
	}

	return nil
}
