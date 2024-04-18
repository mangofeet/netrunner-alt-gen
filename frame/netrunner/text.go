package netrunner

import (
	"strings"

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
		for _, p := range emParts {
			parts = append(parts, p)
		}
	}

	for _, part := range parts {

		if strings.Contains(part, "→") {
			subParts := strings.Split(part, "→")
			writeTextPart(rt, subParts[0], regFace, boldFace, italicFace)
			rt.WriteFace(arrowFace, "→")
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

func writeChunk(rt *canvas.RichText, text string, face *canvas.FontFace) {

	if strings.Contains(text, "[credit]") {
		subParts := strings.Split(text, "[credit]")
		rt.WriteFace(face, subParts[0])

		rt.WritePath(mustLoadGameAsset("CREDIT").Scale(0.0025, 0.0025).Transform(canvas.Identity.ReflectY().Translate(0, face.Size*-0.8)), textColor, canvas.FontMiddle)
		text = subParts[1]
		if len(text) > 0 && text[0] != ' ' {
			text = " " + text
		}
	}

	if strings.Contains(text, "[recurring-credit]") {
		subParts := strings.Split(text, "[recurring-credit]")
		rt.WriteFace(face, subParts[0])

		rt.WritePath(mustLoadGameAsset("RECURRING_CREDIT").Scale(0.017, 0.017).Transform(canvas.Identity.ReflectY().Translate(0, face.Size*-0.8)), textColor, canvas.FontMiddle)
		text = subParts[1]
		if len(text) > 0 && text[0] != ' ' {
			text = " " + text
		}
	}

	if strings.Contains(text, "[click]") {
		subParts := strings.Split(text, "[click]")
		rt.WriteFace(face, subParts[0])

		rt.WritePath(mustLoadGameAsset("CLICK").Scale(0.025, 0.025).Transform(canvas.Identity.ReflectY().Translate(0, face.Size*-1)), textColor, canvas.FontMiddle)
		text = subParts[1]
		if len(text) > 0 && text[0] != ' ' {
			text = " " + text
		}
		// always add a space for this icon
		text = " " + text
	}

	rt.WriteFace(face, text)
}
