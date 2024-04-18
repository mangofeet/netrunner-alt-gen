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
			rt.WriteFace(arrowFace, "♦")
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
	}

	return typeID
}

func getTitleText(card *nrdb.Printing) string {

	if !card.Attributes.IsUnique {
		return card.Attributes.Title
	}

	return fmt.Sprintf("♦ %s", card.Attributes.Title)
}
