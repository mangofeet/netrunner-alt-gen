package cmd

import (
	"encoding/csv"
	"fmt"
	"image"
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

const (
	PAGE_WIDTH_MM  = 216
	PAGE_HEIGHT_MM = 279
	CARD_WIDTH_MM  = 60.0
)

type cardImage struct {
	Width  float64
	Height float64
	DPMM   float64
	Data   *image.RGBA
}

func newCardImage(data *image.RGBA) cardImage {
	imgDPMM := float64(data.Bounds().Max.X) / CARD_WIDTH_MM

	return cardImage{
		Width:  float64(data.Bounds().Max.X) / imgDPMM,
		Height: float64(data.Bounds().Max.Y) / imgDPMM,
		DPMM:   imgDPMM,
		Data:   data,
	}
}

var pnpCmd = &cobra.Command{
	Use:   "pnp [CSV file]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Generate a print & play file containing cards from a CSV.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := generatePnPFile(args[0]); err != nil {
			log.Println("error:", err)
			os.Exit(1)
		}
	},
}

func generatePnPFile(csvPath string) error {
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

	// Open 3x PDF file
	pdfFilePath_three := fmt.Sprintf("%s/pnp_3x.pdf", outputDir)
	log.Printf("Generating print & play file at %s", pdfFilePath_three)
	pdfFile_three, err := os.Create(pdfFilePath_three)
	if err != nil {
		return err
	}
	defer pdfFile_three.Close()

	// Instantiate 3x variables
	p_three := pdf.New(pdfFile_three, PAGE_WIDTH_MM, PAGE_HEIGHT_MM, nil)
	pdfCanvas_three := canvas.New(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
	pdfContext_three := canvas.NewContext(pdfCanvas_three)
	imageCount_three := 0

	// Open 1x PDF file
	pdfFilePath_one := fmt.Sprintf("%s/pnp_1x.pdf", outputDir)
	log.Printf("Generating print & play file at %s", pdfFilePath_one)
	pdfFile_one, err := os.Create(pdfFilePath_one)
	if err != nil {
		return err
	}
	defer pdfFile_one.Close()

	// Instantiate 1x variables
	p_one := pdf.New(pdfFile_one, PAGE_WIDTH_MM, PAGE_HEIGHT_MM, nil)
	pdfCanvas_one := canvas.New(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
	pdfContext_one := canvas.NewContext(pdfCanvas_one)
	imageCount_one := 0

	pageMarginX := (PAGE_WIDTH_MM - (CARD_WIDTH_MM * 3)) / 2
	pageMarginY := -1.0

	// Helper func that draws a card to a PDF file
	var addCardToPage = func(imageCount *int, img cardImage, canv *canvas.Canvas, ctx *canvas.Context, p *pdf.PDF) (*canvas.Canvas, *canvas.Context) {
		imageIndex := *imageCount % 9
		pageX := pageMarginX + (float64(*imageCount%3) * img.Width)
		pageY := PAGE_HEIGHT_MM - (float64((imageIndex/3)+1) * img.Height) - pageMarginY
		ctx.DrawImage(pageX, pageY, img.Data, canvas.DPMM(img.DPMM))
		*imageCount++

		if *imageCount%9 == 8 {
			canv.RenderTo(p)
			p.NewPage(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
			canv = canvas.New(PAGE_WIDTH_MM, PAGE_HEIGHT_MM)
			ctx = canvas.NewContext(canv)
		}

		return canv, ctx
	}

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
		if i == 0 {
			imageCount_three = i
			imageCount_one = i
		}
		card := buildCard(record, cardID)

		// Prepend card name prefix
		if cardNamePrefix != "" {
			card.Attributes.Title = fmt.Sprintf("%s %s", cardNamePrefix, card.Attributes.Title)
		}

		// Create image drawer
		imgPath := fmt.Sprintf("%s/%d.png", imageDir, cardID)
		_, imgFileErr := os.Stat(imgPath)
		var drawer art.Drawer
		drawer = emptyDrawer{}
		if imgFileErr == nil {
			drawer = imageDrawer{
				filename: imgPath,
			}
		}

		// Create canvas
		cnv, err := generateCardCanvas(drawer, card, "", "")
		if err != nil {
			return err
		}

		// Create new card image
		cardImg := newCardImage(rasterizer.Draw(cnv, canvas.DPMM(1), canvas.DefaultColorSpace))

		if pageMarginY == -1 {
			pageMarginY = (PAGE_HEIGHT_MM - (cardImg.Height * 3)) / 2
		}

		// Draw to 3x PDF (only 1x for identities)
		if card.Attributes.CardTypeID == "runner_identity" || card.Attributes.CardTypeID == "corp_identity" {
			pdfCanvas_three, pdfContext_three = addCardToPage(&imageCount_three, cardImg, pdfCanvas_three, pdfContext_three, p_three)
		} else {
			for j := 0; j < 3; j++ {
				pdfCanvas_three, pdfContext_three = addCardToPage(&imageCount_three, cardImg, pdfCanvas_three, pdfContext_three, p_three)
			}
		}

		// Draw to 1x PDF
		pdfCanvas_one, pdfContext_one = addCardToPage(&imageCount_one, cardImg, pdfCanvas_one, pdfContext_one, p_one)

		// Create individual card file
		if genIndividualImages {
			cardName := card.Attributes.Title
			cardImgFilePath := fmt.Sprintf("%s/%d_%s.png", outputDir, cardID, cardName)
			imgFile, err := os.Create(cardImgFilePath)
			if err != nil {
				return err
			}
			cardCanvas := canvas.New(60, 88)
			cardContext := canvas.NewContext(cardCanvas)
			cardContext.DrawImage(0, 0, cardImg.Data, canvas.DPMM(cardImg.DPMM))
			renderers.PNG(canvas.DPMM(cardImg.DPMM))(imgFile, cardCanvas)
			imgFile.Close()
		}

		cardID += 1
	}

	// Render the last 3x PDF page
	pdfCanvas_three.RenderTo(p_three)
	if err := p_three.Close(); err != nil {
		return err
	}

	// Render the last 1x PDF page
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
