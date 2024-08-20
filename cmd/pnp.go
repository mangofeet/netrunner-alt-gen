package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mangofeet/nrdb-go"
	"github.com/spf13/cobra"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/pdf"
	"github.com/tdewolff/canvas/renderers/rasterizer"
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
	const (
		PAGE_WIDTH_MM  = 210
		PAGE_HEIGHT_MM = 297
		CARD_WIDTH_MM  = 60.0
	)

	// Load CSV file
	csvFile, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	// Open PDF file
	pdfFilePath := fmt.Sprintf("%s/pnp.pdf", outputDir)
	log.Printf("Generating print & play file at %s", pdfFilePath)
	pdfFile, err := os.Create(pdfFilePath)
	if err != nil {
		return err
	}
	defer pdfFile.Close()

	// Instantiate variables
	p := pdf.New(pdfFile, PAGE_WIDTH_MM, PAGE_HEIGHT_MM, nil)
	pdfCanvas := canvas.New(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
	pdfContext := canvas.NewContext(pdfCanvas)
	pageMarginX := (PAGE_WIDTH_MM - (CARD_WIDTH_MM * 3)) / 2
	pageMarginY := -1.0

	// Override colors to black & white
	baseColor = "ffffff"
	frameColorBackground = "ffffff"
	frameColorFactionBG = "ffffff"
	frameColorStrengthBG = "ffffff"
	frameColorInfluenceBG = "ffffff"
	frameColorInfluenceLimitBG = "ffffff"
	frameColorMinDeckBG = "ffffff"
	frameColorBorder = "000000"
	frameColorText = "000000"

	cardID := startRow
	for i, record := range records[startRow-1:] {
		card := buildCard(record, cardID)

		// Generate card image
		cnv, err := generateCardCanvas(emptyDrawer{}, card, "", "")
		if err != nil {
			return err
		}

		// Calculate card dimensions
		cardImg := rasterizer.Draw(cnv, canvas.DPMM(1), canvas.DefaultColorSpace)
		imgDPMM := float64(cardImg.Bounds().Max.X) / CARD_WIDTH_MM
		imgWidth := float64(cardImg.Bounds().Max.X) / imgDPMM
		imgHeight := float64(cardImg.Bounds().Max.Y) / imgDPMM
		if pageMarginY == -1 {
			pageMarginY = (PAGE_HEIGHT_MM - (imgHeight * 3)) / 2
		}

		// Draw image
		imageIndex := i % 9
		pageX := pageMarginX + (float64(i%3) * imgWidth)
		pageY := PAGE_HEIGHT_MM - (float64((imageIndex/3)+1) * imgHeight) - pageMarginY
		pdfContext.DrawImage(pageX, pageY, cardImg, canvas.DPMM(imgDPMM))

		cardID += 1

		// Quit if we reached the last record
		if i == len(records)-startRow-1 {
			break
		}

		// Render the page and create a new one
		if i%9 == 8 {
			pdfCanvas.RenderTo(p)

			p.NewPage(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
			pdfCanvas = canvas.New(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
			pdfContext = canvas.NewContext(pdfCanvas)
		}
	}

	// Render the last page
	pdfCanvas.RenderTo(p)

	if err := p.Close(); err != nil {
		return err
	}

	return nil
}

func buildCard(record []string, cardID int) *nrdb.Printing {
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
	card.Attributes.Title = strings.ReplaceAll(record[3], "\n", "")
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
			card.Attributes.AdvancementRequirement = &tokens[1]
		case "Cost":
			card.Attributes.Cost = &tokens[1]
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
	card.Attributes.Text = strings.ReplaceAll(card.Attributes.Text, "{mu}", "[mu]")
	card.Attributes.Text = strings.ReplaceAll(card.Attributes.Text, "{c}", "[credit]")
	card.Attributes.Text = strings.ReplaceAll(card.Attributes.Text, "{recurring}", "[recurring-credit]")
	card.Attributes.Text = strings.ReplaceAll(card.Attributes.Text, "{click}", "[click]")
	card.Attributes.Text = strings.ReplaceAll(card.Attributes.Text, "{sub}", "[subroutine]")
	card.Attributes.Text = strings.ReplaceAll(card.Attributes.Text, "{trash}", "[trash]")
	card.Attributes.Text = strings.ReplaceAll(card.Attributes.Text, "{interrupt}", "[interrupt]")

	return card
}
