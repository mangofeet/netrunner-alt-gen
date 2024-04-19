package netrunner

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func getCardText(text string, fontSize, cardTextBoxW, cardTextBoxH float64) *canvas.Text {

	regFace := getFont(fontSize, canvas.FontRegular)
	boldFace := getFont(fontSize, canvas.FontBold)
	italicFace := getFont(fontSize, canvas.FontItalic)
	arrowFace := getFont(fontSize, canvas.FontExtraBold)
	uniqueFace := getFont(fontSize, canvas.FontExtraBold)

	rt := canvas.NewRichText(regFace)

	var parts []string
	strongParts := strings.Split(text, "<strong>")
	for _, p := range strongParts {
		emParts := strings.Split(p, "<em>")
		parts = append(parts, emParts...)
	}

	for _, part := range parts {

		if strings.Contains(part, "→") {
			subParts := strings.Split(part, "→")
			writeTextPart(rt, subParts[0], regFace, boldFace, italicFace)
			rt.WriteFace(arrowFace, "→")
			part = subParts[1]
		}

		if strings.Contains(part, "♦") {
			subParts := strings.Split(part, "♦")
			writeTextPart(rt, subParts[0], regFace, boldFace, italicFace)
			rt.WriteFace(uniqueFace, "♦")
			part = subParts[1]
		}

		writeTextPart(rt, part, regFace, boldFace, italicFace)
	}

	return rt.ToText(
		cardTextBoxW, cardTextBoxH,
		canvas.Left, canvas.Top,
		0, 0)

}

func writeTextPart(rt *canvas.RichText, text string, regFace, boldFace, italicFace *canvas.FontFace) {

	text = strings.ReplaceAll(text, "\n", "\n\n")

	if strings.Contains(text, "</strong>") {
		subParts := strings.Split(text, "</strong>")
		writeChunk(rt, subParts[0], boldFace)
		text = subParts[1]
	}

	if strings.Contains(text, "</em>") {
		subParts := strings.Split(text, "</em>")
		writeChunk(rt, subParts[0], italicFace)
		text = subParts[1]
	}

	writeChunk(rt, text, regFace)

}

var replacementCheck = regexp.MustCompile(`\[[a-z-]+\]`)

func replaceSymbol(rt *canvas.RichText, symbol, svgName, text string, face *canvas.FontFace, scaleFactor, translateFactor float64) string {
	if strings.Contains(text, symbol) {
		subParts := strings.Split(text, symbol)
		rt.WriteFace(face, subParts[0])

		path := mustLoadGameAsset(svgName).Scale(face.Size*scaleFactor, face.Size*scaleFactor).Transform(canvas.Identity.ReflectY().Translate(0, face.Size*-1*translateFactor))

		rt.WritePath(path, textColor, canvas.FontMiddle)
		text = subParts[1]
		if len(text) == 0 || (text[0] != ' ' && text[0] != ',') {
			text = " " + text
		}
	}

	return text

}

func writeChunk(rt *canvas.RichText, text string, face *canvas.FontFace) {

	text = replaceSymbol(rt, "[mu]", "Mu", text, face, 0.0002, 0.8)
	text = replaceSymbol(rt, "[credit]", "CREDIT", text, face, 0.000025, 0.8)
	text = replaceSymbol(rt, "[recurring-credit]", "RECURRING_CREDIT", text, face, 0.00014, 0.8)
	text = replaceSymbol(rt, "[click]", "CLICK", text, face, 0.0002, 1)
	text = replaceSymbol(rt, "[subroutine]", "SUBROUTINE", text, face, 0.0002, 1)
	text = replaceSymbol(rt, "[trash]", "TRASH_ABILITY", text, face, 0.0002, 1)

	if replacementCheck.MatchString(text) {
		writeChunk(rt, text, face)
		return
	}

	rt.WriteFace(face, text)

}

func getTypeName(typeID string) string {
	switch typeID {
	case "program":
		return "Program"
	case "resource":
		return "Resource"
	case "hardware":
		return "Hardware"
	case "event":
		return "Event"
	}

	return typeID
}

func getTitleText(card *nrdb.Printing) string {

	if !card.Attributes.IsUnique {
		return card.Attributes.Title
	}

	return fmt.Sprintf("♦ %s", card.Attributes.Title)
}

type textBoxDimensions struct {
	left, right, height, bottom, top float64
}

func drawCardText(ctx *canvas.Context, card *nrdb.Printing, fontSize, indentCutoff, indent float64, box, typeBox textBoxDimensions) {

	canvasWidth, canvasHeight := ctx.Size()
	strokeWidth := canvasHeight * 0.0023

	cardTextPadding := canvasWidth * 0.02
	cardTextX := box.left + cardTextPadding
	cardTextY := box.height - cardTextPadding
	typeTextX := cardTextX
	typeTextY := typeBox.bottom + typeBox.height - cardTextPadding
	cardTextBoxW := box.right - box.left - (cardTextPadding * 2.5)
	cardTextBoxH := box.height
	typeTextBoxW := typeBox.right - typeBox.left - (cardTextPadding * 2)
	typeTextBoxH := typeBox.height

	var tText *canvas.Text

	typeName := getTypeName(card.Attributes.CardTypeID)

	if card.Attributes.DisplaySubtypes != nil {
		tText = getCardText(fmt.Sprintf("<strong>%s</strong> - %s", typeName, *card.Attributes.DisplaySubtypes), fontSize, typeTextBoxW, typeTextBoxH)
	} else {
		tText = getCardText(fmt.Sprintf("<strong>%s</strong>", typeName), fontSize, typeTextBoxW, typeTextBoxH)
	}

	ctx.DrawText(typeTextX, typeTextY, tText)

	cText := getCardText(card.Attributes.Text, fontSize, cardTextBoxW, cardTextBoxH)

	var leftoverText string

	_, lastLineH := cText.Heights()

	for lastLineH > cardTextBoxH*0.75 {
		fontSize -= strokeWidth
		cText = getCardText(card.Attributes.Text, fontSize, cardTextBoxW, cardTextBoxH)
		_, lastLineH = cText.Heights()
	}

	i := 0
	_, lastLineH = cText.Heights()
	for lastLineH > indentCutoff {

		i++

		lines := strings.Split(card.Attributes.Text, "\n")

		leftoverText = strings.Join(lines[len(lines)-i:], "\n")
		newText := strings.Join(lines[:len(lines)-i], "\n")

		cText = getCardText(newText, fontSize, cardTextBoxW, cardTextBoxH)

		_, lastLineH = cText.Heights()

	}

	ctx.DrawText(cardTextX, cardTextY, cText)

	if leftoverText != "" {
		newCardTextX := cardTextX + indent
		if !cText.Empty() {
			cardTextY = cardTextY - (lastLineH + fontSize*0.4)
		}

		cText := getCardText(leftoverText, fontSize, cardTextBoxW-(newCardTextX-cardTextX)-cardTextBoxW*0.03, cardTextBoxH)
		ctx.DrawText(newCardTextX, cardTextY, cText)
	}

}
