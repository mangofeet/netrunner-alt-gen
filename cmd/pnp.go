package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/spf13/cobra"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"github.com/tdewolff/canvas/renderers/pdf"
	"github.com/tdewolff/canvas/renderers/rasterizer"
)

var pnpCmd = &cobra.Command{
	Use:   "pnp [CSV file] [Prefix]",
	Args:  cobra.MinimumNArgs(2),
	Short: "Generate a print & play file containing cards from a CSV. Card Titles can be prepended with a prefix for version tracking",
	Run: func(cmd *cobra.Command, args []string) {
		if err := generatePnPFile(args[0], args[1]); err != nil {
			log.Println("error:", err)
			os.Exit(1)
		}
	},
}

func generatePnPFile(csvPath string, prefix string) error {
	const (
		PAGE_WIDTH_MM  = 216
		PAGE_HEIGHT_MM = 279
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
	pdfFilePath_three := fmt.Sprintf("%s/pnp_3x.pdf", outputDir)
	log.Printf("Generating print & play file at %s", pdfFilePath_three)
	pdfFile_three, err := os.Create(pdfFilePath_three)
	if err != nil {
		return err
	}
	defer pdfFile_three.Close()
	// Instantiate variables
	p_three := pdf.New(pdfFile_three, PAGE_WIDTH_MM, PAGE_HEIGHT_MM, nil)
	pdfCanvas_three := canvas.New(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
	pdfContext_three := canvas.NewContext(pdfCanvas_three)

	pdfFilePath_one := fmt.Sprintf("%s/pnp_1x.pdf", outputDir)
	log.Printf("Generating print & play file at %s", pdfFilePath_one)
	pdfFile_one, err := os.Create(pdfFilePath_one)
	if err != nil {
		return err
	}
	defer pdfFile_one.Close()
	// Instantiate variables
	p_one := pdf.New(pdfFile_one, PAGE_WIDTH_MM, PAGE_HEIGHT_MM, nil)
	pdfCanvas_one := canvas.New(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
	pdfContext_one := canvas.NewContext(pdfCanvas_one)

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
	imageCount_one := 0
	imageCount_three := 0
	for i, record := range records[startRow-1:] {
		if i == 0 {
			imageCount_three = i
			imageCount_one = i
		}
		card := buildCard(record, cardID)

		// prepend 'Dev 8.2' etc
		card.Attributes.Title = fmt.Sprintf("%s %s", prefix, card.Attributes.Title)
		// Generate card image

		imgPath := fmt.Sprintf("piggybank_images/%d.png", cardID)
		_, imgFileErr := os.Stat(imgPath)
		var drawer art.Drawer
		drawer = emptyDrawer{}
		if imgFileErr == nil {
			drawer = imageDrawer{
				filename: imgPath,
			}
		}
		cnv, err := generateCardCanvas(drawer, card, "", "")
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
		// print 3x non IDs
		if card.Attributes.CardTypeID == "runner_identity" || card.Attributes.CardTypeID == "corp_identity" {
			imageIndex := imageCount_three % 9
			pageX := pageMarginX + (float64(imageCount_three%3) * imgWidth)
			pageY := PAGE_HEIGHT_MM - (float64((imageIndex/3)+1) * imgHeight) - pageMarginY
			pdfContext_three.DrawImage(pageX, pageY, cardImg, canvas.DPMM(imgDPMM))
			imageCount_three++

			if imageCount_three%9 == 8 {
				pdfCanvas_three.RenderTo(p_three)
				p_three.NewPage(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
				pdfCanvas_three = canvas.New(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
				pdfContext_three = canvas.NewContext(pdfCanvas_three)
			}

		} else {
			for j := 0; j < 3; j++ {
				imageIndex := imageCount_three % 9
				pageX := pageMarginX + (float64(imageCount_three%3) * imgWidth)
				pageY := PAGE_HEIGHT_MM - (float64((imageIndex/3)+1) * imgHeight) - pageMarginY
				pdfContext_three.DrawImage(pageX, pageY, cardImg, canvas.DPMM(imgDPMM))
				imageCount_three++
				// need to handle after each image or else we cna miss a page accidentally
				if imageCount_three%9 == 8 {
					pdfCanvas_three.RenderTo(p_three)
					p_three.NewPage(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
					pdfCanvas_three = canvas.New(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
					pdfContext_three = canvas.NewContext(pdfCanvas_three)
				}
			}
		}

		imageIndex := imageCount_one % 9
		pageX := pageMarginX + (float64(imageCount_one%3) * imgWidth)
		pageY := PAGE_HEIGHT_MM - (float64((imageIndex/3)+1) * imgHeight) - pageMarginY
		pdfContext_one.DrawImage(pageX, pageY, cardImg, canvas.DPMM(imgDPMM))
		imageCount_one++

		if imageCount_one%9 == 8 {
			pdfCanvas_one.RenderTo(p_one)
			p_one.NewPage(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
			pdfCanvas_one = canvas.New(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
			pdfContext_one = canvas.NewContext(pdfCanvas_one)
		}

		cardName := card.Attributes.Title
		cardImgFilePath := fmt.Sprintf("%s/%d_%s.png", outputDir, cardID, cardName)
		imgFile, err := os.Create(cardImgFilePath)
		if err != nil {
			return err
		}

		cardCanvas := canvas.New(60, 88)
		cardContext := canvas.NewContext(cardCanvas)
		cardContext.DrawImage(0, 0, cardImg, canvas.DPMM(imgDPMM))
		renderers.PNG(canvas.DPMM(imgDPMM))(imgFile, cardCanvas)
		imgFile.Close()
		cardID += 1

		// Quit if we reached the last record
		//if i == len(records)-startRow-1 {
		//	break
		//}

		// Render the page and create a new one

	}

	// Render the last page
	pdfCanvas_three.RenderTo(p_three)

	if err := p_three.Close(); err != nil {
		return err
	}

	pdfCanvas_one.RenderTo(p_one)

	if err := p_one.Close(); err != nil {
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
	summary_sections := strings.Split(record[4], "====")
	lower_sections := summary_sections[1]
	summary_sections = []string{summary_sections[0]}
	summary_sections = append(summary_sections, strings.Split(lower_sections, "----")...)
	for i, section := range summary_sections {
		summary_sections[i] = strings.Trim(section, "\n")
	}

	// Set ID, faction, name, and type
	card.ID = strconv.Itoa(cardID)
	card.Attributes.PositionInSet = cardID
	if record[3] == "Weyland" {
		card.Attributes.FactionID = "weyland_consortium"
	} else {
		card.Attributes.FactionID = strings.ReplaceAll(strings.ToLower(record[3]), "-", "_")
	}
	card.Attributes.Title = strings.ReplaceAll(record[2], "\n", "")
	card.Attributes.StrippedTitle = titleStripper.Replace(record[2])
	card.Attributes.IsUnique = strings.Contains(summary_sections[0], "◆")
	if record[1] == "Runner-ID" {
		card.Attributes.CardTypeID = "runner_identity"
		card.Attributes.CardAbilities.MUProvided = nil
	} else if record[1] == "Corp-ID" {
		card.Attributes.CardTypeID = "corp_identity"
	} else {
		card.Attributes.CardTypeID = strings.ToLower(record[1])
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
			card.Attributes.TrashCost = &val
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
