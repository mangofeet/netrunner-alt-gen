package basic

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func getTitle(card *nrdb.Printing) string {

	if card.Attributes.CardTypeID == "runner_identity" || card.Attributes.CardTypeID == "corp_identity" {
		return strings.Split(card.Attributes.Title, ":")[0]
	}

	if !card.Attributes.IsUnique {
		return card.Attributes.Title
	}

	return fmt.Sprintf("♦ %s", card.Attributes.Title)
}

func getSubtitle(card *nrdb.Printing) string {

	if card.Attributes.CardTypeID != "runner_identity" && card.Attributes.CardTypeID != "corp_identity" {
		return ""
	}

	subtitle := strings.Join(strings.Split(card.Attributes.Title, ":")[1:], ":")

	if card.Attributes.CardTypeID == "runner_identity" && card.Attributes.Pronouns != nil {
		subtitle += fmt.Sprintf(" (%s)", *card.Attributes.Pronouns)
	}

	return subtitle
}

func getTitleText(ctx *canvas.Context, card *nrdb.Printing, fontSize, maxWidth, height float64, align canvas.TextAlign) *canvas.Text {
	return getFittedText(ctx, getTitle(card), fontSize, maxWidth, height, align)
}

func getSubtitleText(ctx *canvas.Context, card *nrdb.Printing, fontSize, maxWidth, height float64, align canvas.TextAlign) *canvas.Text {
	return getFittedText(ctx, getSubtitle(card), fontSize, maxWidth, height, align)
}

func getFittedText(ctx *canvas.Context, title string, fontSize, maxWidth, height float64, align canvas.TextAlign) *canvas.Text {

	text := getCardText(title, fontSize, maxWidth*2, height, align)

	strokeWidth := getStrokeWidth(ctx)

	for text.Bounds().W > maxWidth {
		fontSize -= strokeWidth
		text = getCardText(title, fontSize, maxWidth*2, height, align)
	}

	// get it a final time to get the width correct
	text = getCardText(title, fontSize, maxWidth, height, align)

	return text
}

func getCardText(text string, fontSize, cardTextBoxW, cardTextBoxH float64, align canvas.TextAlign) *canvas.Text {

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
		align, canvas.Top,
		0, 0)

}

func writeTextPart(rt *canvas.RichText, text string, regFace, boldFace, italicFace *canvas.FontFace) {

	text = strings.ReplaceAll(text, "\n", "\n\n")
	text = strings.ReplaceAll(text, "<BR>", "\n")

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
		for _, chunk := range subParts[:len(subParts)-1] {
			writeChunk(rt, chunk, face)
			path := mustLoadGameAsset(svgName).Scale(face.Size*scaleFactor, face.Size*scaleFactor).Transform(canvas.Identity.ReflectY().Translate(0, face.Size*-1*translateFactor))
			rt.WritePath(path, textColor, canvas.FontMiddle)
		}
		text = subParts[len(subParts)-1]
		if len(text) == 0 || (text[0] != ' ' && text[0] != ',' && text[0] != '.') {
			text = " " + text
		}
	}

	return text

}

func writeChunk(rt *canvas.RichText, text string, face *canvas.FontFace) {

	text = replaceSymbol(rt, "[mu]", "MU", text, face, 0.0002, 0.8)
	text = replaceSymbol(rt, "[credit]", "CREDIT", text, face, 0.000025, 0.8)
	text = replaceSymbol(rt, "[recurring-credit]", "RECURRING_CREDIT", text, face, 0.00014, 0.8)
	text = replaceSymbol(rt, "[click]", "CLICK", text, face, 0.0002, 1)
	text = replaceSymbol(rt, "[subroutine]", "SUBROUTINE", text, face, 0.0002, 1)
	text = replaceSymbol(rt, "[trash]", "TRASH_ABILITY", text, face, 0.0002, 1)
	text = replaceSymbol(rt, "[interrupt]", "INTERRUPT", text, face, 0.0002, 1)

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
	case "runner_identity", "corp_identity":
		return "Identity"
	case "ice":
		return "Ice"
	case "asset":
		return "Asset"
	case "upgrade":
		return "Upgrade"
	case "operation":
		return "Operation"
	case "agenda":
		return "Agenda"
	}

	return typeID
}

type textBoxDimensions struct {
	left, right, bottom, top float64
	width, height            float64
	align                    canvas.TextAlign
}

type additionalText struct {
	content string
	align   canvas.TextAlign
}

func drawCardText(ctx *canvas.Context, card *nrdb.Printing, fontSize, indentCutoff, indent float64, box textBoxDimensions, extra ...additionalText) {

	if box.align == 0 {
		box.align = canvas.Left
	}

	strokeWidth := getStrokeWidth(ctx)

	paddingLR, paddingTB := getCardTextPadding(ctx)
	x := box.left + paddingLR
	y := box.height - paddingTB
	if box.top != 0 {
		y = box.top - paddingTB
	}
	w := box.right - box.left - (paddingLR * 2.5)
	h := box.height

	cText := getCardText(card.Attributes.Text, fontSize, w, h, box.align)
	var fTexts []*canvas.Text
	for _, txt := range extra {
		fTexts = append(fTexts, getCardText(txt.content, fontSize, w*0.85, h, txt.align))
	}

	var leftoverText string

	_, lastLineH := cText.Heights()

	for _, txt := range fTexts {
		_, extraH := txt.Heights()
		lastLineH += extraH
	}

	for lastLineH > h*0.5 {
		fontSize -= strokeWidth
		cText = getCardText(card.Attributes.Text, fontSize, w, h, box.align)
		fTexts = []*canvas.Text{}
		for _, txt := range extra {
			fTexts = append(fTexts, getCardText(txt.content, fontSize, w*0.85, h, txt.align))
		}
		_, lastLineH = cText.Heights()
		for _, txt := range fTexts {
			_, extraH := txt.Heights()
			lastLineH += extraH
		}
	}

	i := 0
	_, lastLineH = cText.Heights()
	for lastLineH > indentCutoff {

		i++

		lines := strings.Split(card.Attributes.Text, "\n")

		leftoverText = strings.Join(lines[len(lines)-i:], "\n")
		newText := strings.Join(lines[:len(lines)-i], "\n")

		cText = getCardText(newText, fontSize, w, h, box.align)

		_, lastLineH = cText.Heights()

	}

	ctx.DrawText(x, y, cText)

	newCardTextX := x

	if leftoverText != "" {
		newCardTextX = x + indent
		if !cText.Empty() {
			y = y - (lastLineH + fontSize*0.4)
		}

		cText := getCardText(leftoverText, fontSize, w-(newCardTextX-x)-w*0.03, h, box.align)
		ctx.DrawText(newCardTextX, y, cText)
		_, lastLineH = cText.Heights()
	}

	newCardTextX += w * 0.08
	y = y - (lastLineH + fontSize*0.4)
	for _, txt := range fTexts {
		ctx.DrawText(newCardTextX, y, txt)
		_, lastLineH = txt.Heights()
		y = y - (lastLineH)
	}

}

func getCardTextPadding(ctx *canvas.Context) (lr, tb float64) {
	canvasWidth, _ := ctx.Size()

	lr = canvasWidth * 0.03
	tb = canvasWidth * 0.02

	return lr, tb

}

func drawTypeText(ctx *canvas.Context, card *nrdb.Printing, fontSize float64, box textBoxDimensions) {

	if box.align == 0 {
		box.align = canvas.Left
	}

	paddingLR, _ := getCardTextPadding(ctx)

	w := box.right - box.left - (paddingLR * 2)
	h := box.height
	box.top = box.bottom + box.height

	typeText := getTypeText(card, fontSize, w, h, box.align)

	x := box.left + paddingLR
	y := box.top - (box.height-(typeText.Bounds().H))*0.5

	ctx.DrawText(x, y, typeText)

}

func getTypeText(card *nrdb.Printing, fontSize, w, h float64, align canvas.TextAlign) *canvas.Text {
	var tText *canvas.Text
	typeName := getTypeName(card.Attributes.CardTypeID)

	if card.Attributes.DisplaySubtypes != nil {
		tText = getCardText(fmt.Sprintf("<strong>%s</strong> - %s", typeName, *card.Attributes.DisplaySubtypes), fontSize, w, h, align)
	} else {
		tText = getCardText(fmt.Sprintf("<strong>%s</strong>", typeName), fontSize, w, h, align)
	}

	return tText
}
