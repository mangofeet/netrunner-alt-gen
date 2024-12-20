package basic

import (
	"fmt"
	"math"
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

func (fb FrameBasic) getTitleText(ctx *canvas.Context, card *nrdb.Printing, fontSize, maxWidth, height float64, align canvas.TextAlign) *canvas.Text {
	return fb.getHorizontalFittedText(ctx, getTitle(card), fontSize, maxWidth, height, align)
}

func (fb FrameBasic) getSubtitleText(ctx *canvas.Context, card *nrdb.Printing, fontSize, maxWidth, height float64, align canvas.TextAlign) *canvas.Text {
	return fb.getHorizontalFittedText(ctx, getSubtitle(card), fontSize, maxWidth, height, align)
}

func (fb FrameBasic) getFittedText(ctx *canvas.Context, title string, fontSize, maxWidth, maxHeight float64, align canvas.TextAlign) *canvas.Text {
	return fb.getFittedTextWithFont(ctx, title, fontSize, maxWidth, maxHeight, align, fb.getFont(fontSize, canvas.FontRegular))
}

func (fb FrameBasic) getFittedTextWithFont(ctx *canvas.Context, title string, fontSize, maxWidth, maxHeight float64, align canvas.TextAlign, font *canvas.FontFace) *canvas.Text {
	if maxWidth == 0 {
		return fb.getVerticalFittedText(ctx, title, fontSize, maxWidth, maxHeight, align)
	}
	if maxHeight == 0 {
		return fb.getHorizontalFittedText(ctx, title, fontSize, maxWidth, maxHeight, align)
	}

	text := fb.getCardTextWithFont(title, fontSize, maxWidth*2, maxHeight*2, align, font)

	strokeWidth := getStrokeWidth(ctx)

	for (text.Bounds().W > maxWidth || text.Bounds().H > maxHeight) && fontSize > 0 {
		fontSize -= strokeWidth
		text = fb.getCardTextWithFont(title, fontSize, maxWidth*2, maxHeight*2, align, font)
	}

	return fb.getCardTextWithFont(title, fontSize, maxWidth, maxHeight, align, font)

}

func (fb FrameBasic) getHorizontalFittedText(ctx *canvas.Context, title string, fontSize, maxWidth, height float64, align canvas.TextAlign) *canvas.Text {
	return fb.getHorizontalFittedTextWithFont(ctx, title, fontSize, maxWidth, height, align, fb.getFont(fontSize, canvas.FontRegular))
}

func (fb FrameBasic) getHorizontalFittedTextWithFont(ctx *canvas.Context, title string, fontSize, maxWidth, height float64, align canvas.TextAlign, font *canvas.FontFace) *canvas.Text {

	text := fb.getCardTextWithFont(title, fontSize, maxWidth*2, height, align, font)

	strokeWidth := getStrokeWidth(ctx)

	for text.Bounds().W > maxWidth && fontSize > 0 {
		fontSize -= strokeWidth
		text = fb.getCardTextWithFont(title, fontSize, maxWidth*2, height, align, font)
	}

	// get it a final time to get the width correct
	text = fb.getCardTextWithFont(title, fontSize, maxWidth, height, align, font)

	return text
}

func (fb FrameBasic) getVerticalFittedText(ctx *canvas.Context, title string, fontSize, width, maxHeight float64, align canvas.TextAlign) *canvas.Text {
	return fb.getVerticalFittedTextWithFont(ctx, title, fontSize, width, maxHeight, align, fb.getFont(fontSize, canvas.FontRegular))
}

func (fb FrameBasic) getVerticalFittedTextWithFont(ctx *canvas.Context, title string, fontSize, width, maxHeight float64, align canvas.TextAlign, font *canvas.FontFace) *canvas.Text {

	text := fb.getCardTextWithFont(title, fontSize, width, maxHeight*2, align, font)

	strokeWidth := getStrokeWidth(ctx)

	for text.Bounds().H > maxHeight {
		fontSize -= strokeWidth
		text = fb.getCardTextWithFont(title, fontSize, width, maxHeight*2, align, font)
	}

	// get it a final time to get the width correct
	text = fb.getCardTextWithFont(title, fontSize, width, maxHeight, align, font)

	return text
}

func (fb FrameBasic) getCardText(text string, fontSize, cardTextBoxW, cardTextBoxH float64, align canvas.TextAlign) *canvas.Text {
	regFace := fb.getFont(fontSize, canvas.FontRegular)
	return fb.getCardTextWithFont(text, fontSize, cardTextBoxW, cardTextBoxH, align, regFace)
}

func (fb FrameBasic) getCardTextWithFont(text string, fontSize, cardTextBoxW, cardTextBoxH float64, align canvas.TextAlign, font *canvas.FontFace) *canvas.Text {

	regFace := fontFamily.Face(fontSize, font.Fill.Color, font.Style)
	boldFace := fb.getFont(fontSize, canvas.FontBold)
	italicFace := fb.getFont(fontSize, canvas.FontItalic)
	arrowFace := fb.getFont(fontSize, canvas.FontExtraBold)
	uniqueFace := fb.getFont(fontSize, canvas.FontExtraBold)

	rt := canvas.NewRichText(font)

	var parts []string
	strongParts := strings.Split(text, "<strong>")
	for _, p := range strongParts {
		emParts := strings.Split(p, "<em>")
		parts = append(parts, emParts...)
	}

	for _, part := range parts {

		if strings.Contains(part, "→") {
			subParts := strings.Split(part, "→")
			fb.writeTextPart(rt, subParts[0], regFace, boldFace, italicFace)
			rt.WriteFace(arrowFace, "→")
			part = subParts[1]
		}

		if strings.Contains(part, "♦") {
			subParts := strings.Split(part, "♦")
			fb.writeTextPart(rt, subParts[0], regFace, boldFace, italicFace)
			rt.WriteFace(uniqueFace, "♦")
			part = subParts[1]
		}

		fb.writeTextPart(rt, part, regFace, boldFace, italicFace)
	}

	return rt.ToText(
		cardTextBoxW, cardTextBoxH,
		align, canvas.Top,
		0, 0)

}

func (fb FrameBasic) writeTextPart(rt *canvas.RichText, text string, regFace, boldFace, italicFace *canvas.FontFace) {

	text = strings.ReplaceAll(text, "\n", "\n") // used to be replace text "\n" "\n\n"
	text = strings.ReplaceAll(text, "<BR>", "\n")
	text = strings.Replace(text, "</li>", "", -1)
	text = strings.Replace(text, "</ul>", "", -1)
	text = strings.Replace(text, "<ul>", "", -1)
	text = strings.Replace(text, "<li>", "\n •", -1)

	if strings.Contains(text, "</strong>") {
		subParts := strings.Split(text, "</strong>")
		fb.writeChunk(rt, subParts[0], boldFace)
		text = subParts[1]
	}

	if strings.Contains(text, "</em>") {
		subParts := strings.Split(text, "</em>")
		fb.writeChunk(rt, subParts[0], italicFace)
		text = subParts[1]
	}

	fb.writeChunk(rt, text, regFace)

}

var replacementCheck = regexp.MustCompile(`\[[a-z-]+\]`)

func (fb FrameBasic) replaceSymbol(rt *canvas.RichText, symbol, svgName, text string, face *canvas.FontFace, scaleFactor, translateFactor float64) string {
	if strings.Contains(text, symbol) {
		subParts := strings.Split(text, symbol)
		for _, chunk := range subParts[:len(subParts)-1] {
			fb.writeChunk(rt, chunk, face)
			path := mustLoadGameAsset(svgName).Scale(face.Size*scaleFactor, face.Size*scaleFactor).Transform(canvas.Identity.ReflectY().Translate(0, face.Size*-1*translateFactor))
			rt.WritePath(path, fb.getColorText(), canvas.FontMiddle)
		}
		text = subParts[len(subParts)-1]
		if len(text) == 0 || (text[0] != ' ' && text[0] != ',' && text[0] != '.') {
			text = " " + text
		}
	}

	return text

}

func (fb FrameBasic) writeChunk(rt *canvas.RichText, text string, face *canvas.FontFace) {

	text = fb.replaceSymbol(rt, "[mu]", "MU", text, face, 0.0002, 0.8)
	text = fb.replaceSymbol(rt, "[credit]", "CREDIT", text, face, 0.000025, 0.8)
	text = fb.replaceSymbol(rt, "[recurring-credit]", "RECURRING_CREDIT", text, face, 0.00014, 0.8)
	text = fb.replaceSymbol(rt, "[click]", "CLICK", text, face, 0.0002, 1)
	text = fb.replaceSymbol(rt, "[subroutine]", "SUBROUTINE", text, face, 0.0002, 1)
	text = fb.replaceSymbol(rt, "[trash]", "TRASH_ABILITY", text, face, 0.0002, 1)
	text = fb.replaceSymbol(rt, "[interrupt]", "INTERRUPT", text, face, 0.0002, 1)

	if replacementCheck.MatchString(text) {
		fb.writeChunk(rt, text, face)
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

type additionalTextType string

const (
	additionalTextTypeFlavor      additionalTextType = "flavor"
	additionalTextTypeAttribution additionalTextType = "attribution"
)

type additionalText struct {
	content  string
	textType additionalTextType
	align    canvas.TextAlign
}

func (fb FrameBasic) drawCardText(ctx *canvas.Context, card *nrdb.Printing, fontSize, indentCutoff, indent float64, box textBoxDimensions, extra ...additionalText) {

	if box.align == 0 {
		box.align = canvas.Left
	}

	strokeWidth := getStrokeWidth(ctx)
	_, canvasHeight := ctx.Size()

	originalFontSize := fontSize
	extraFontSize := originalFontSize * 0.7
	maxExtraFontSize := extraFontSize
	textBottom := canvasHeight * 0.0592

	paddingLR, paddingTB := getCardTextPadding(ctx)
	x := box.left + paddingLR
	if box.top == 0 {
		box.top = box.height
	} else { // this is for ice
		textBottom = box.top - box.height + paddingTB
	}
	y := box.top - paddingTB
	w := box.right - box.left - (paddingLR * 2.5)
	h := box.height

	cText := fb.getCardText(card.Attributes.Text, fontSize, w, h, box.align)
	var fTexts []*canvas.Text
	for _, txt := range extra {
		fTexts = append(fTexts, fb.getCardText(txt.content, extraFontSize, w*0.85, h, txt.align))
	}

	var leftoverText string

	lastLineH := cText.Bounds().H

	if len(fTexts) > 0 {
		lastLineH += fontSize * 0.3
	}
	for _, txt := range fTexts {
		extraH := txt.Bounds().H
		lastLineH += extraH
	}

	for y-lastLineH < textBottom {
		fontSize -= strokeWidth
		extraFontSize = math.Min(maxExtraFontSize, fontSize)

		// remake font boxes with new font size
		cText = fb.getCardText(card.Attributes.Text, fontSize, w, h, box.align)
		for i, txt := range extra {
			fTexts[i] = fb.getCardText(txt.content, extraFontSize, w*0.85, h, txt.align)
		}

		// get new last line height
		lastLineH = cText.Bounds().H
		if len(fTexts) > 0 {
			lastLineH += fontSize * 0.3
		}
		for _, txt := range fTexts {
			extraH := txt.Bounds().H
			lastLineH += extraH
		}

	}

	i := 0
	lastLineH = cText.Bounds().H

	for y-lastLineH < indentCutoff {

		i++

		lines := strings.Split(card.Attributes.Text, "\n")

		leftoverText = strings.Join(lines[len(lines)-i:], "\n")
		newText := strings.Join(lines[:len(lines)-i], "\n")

		cText = fb.getCardText(newText, fontSize, w, h, box.align)

		lastLineH = cText.Bounds().H

	}

	ctx.DrawText(x, y, cText)

	newCardTextX := x

	if leftoverText != "" {
		newCardTextX = x + indent
		if !cText.Empty() {
			y = y - (lastLineH + fontSize*0.4)
		}

		cText := fb.getCardText(leftoverText, fontSize, w-(newCardTextX-x)-w*0.03, h, box.align)
		ctx.DrawText(newCardTextX, y, cText)
		lastLineH = cText.Bounds().H
	}

	newCardTextX += w * 0.08
	y = y - (lastLineH + fontSize*0.4)
	widestLine := 0.0
	extraFontSize = math.Min(maxExtraFontSize, fontSize)
	for _, ln := range extra {

		textWidth := w - (newCardTextX - x) - w*0.03
		if ln.textType == additionalTextTypeAttribution {
			textWidth = math.Min(widestLine*1.2, textWidth)
		}

		txt := fb.getCardText(ln.content, extraFontSize, textWidth, h, ln.align)
		ctx.DrawText(newCardTextX, y, txt)

		widestLine = math.Max(txt.Bounds().W, widestLine)
		lastLineH = txt.Bounds().H
		y = y - (lastLineH)
	}

}

func getCardTextPadding(ctx *canvas.Context) (lr, tb float64) {
	canvasWidth, _ := ctx.Size()

	lr = canvasWidth * 0.03
	tb = canvasWidth * 0.02

	return lr, tb

}

func (fb FrameBasic) drawTypeText(ctx *canvas.Context, card *nrdb.Printing, fontSize float64, box textBoxDimensions) {

	if box.align == 0 {
		box.align = canvas.Left
	}

	paddingLR, _ := getCardTextPadding(ctx)

	w := box.right - box.left - (paddingLR * 2)
	h := box.height
	box.top = box.bottom + box.height

	typeText := fb.getTypeText(ctx, card, fontSize, w, h, box.align)

	x := box.left + paddingLR
	y := box.top - (box.height-(typeText.Bounds().H))*0.5

	ctx.DrawText(x, y, typeText)

}

func (fb FrameBasic) getTypeText(ctx *canvas.Context, card *nrdb.Printing, fontSize, w, h float64, align canvas.TextAlign) *canvas.Text {
	var tText *canvas.Text
	typeName := getTypeName(card.Attributes.CardTypeID)

	if card.Attributes.DisplaySubtypes != nil && card.Attributes.TrashCost != nil {
		tText = fb.getFittedText(ctx, fmt.Sprintf("<strong>%s</strong> - %s - Trash: %d", typeName, *card.Attributes.DisplaySubtypes, *card.Attributes.TrashCost), fontSize, w, h, align)
	} else if card.Attributes.DisplaySubtypes != nil && card.Attributes.TrashCost == nil {
		tText = fb.getFittedText(ctx, fmt.Sprintf("<strong>%s</strong> - %s", typeName, *card.Attributes.DisplaySubtypes), fontSize, w, h, align)
	} else if card.Attributes.TrashCost != nil {
		tText = fb.getFittedText(ctx, fmt.Sprintf("<strong>%s</strong> - Trash: %d", typeName, *card.Attributes.TrashCost), fontSize, w, h, align)
	} else {
		tText = fb.getFittedText(ctx, fmt.Sprintf("<strong>%s</strong>", typeName), fontSize, w, h, align)
	}

	return tText
}
