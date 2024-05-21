package basic

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func (fb FrameBasic) Back() art.Drawer {

	return art.DrawerFunc(func(ctx *canvas.Context, card *nrdb.Printing) error {

		canvasWidth, canvasHeight := ctx.Size()

		attributionBoxTop := canvasHeight * 0.12
		attributionBoxBottom := canvasHeight * 0.0591
		attributionBoxHeight := attributionBoxTop - attributionBoxBottom
		// attributionBoxLeft := canvasWidth * 0.25
		// attributionBoxRight := canvasWidth * 0.75
		// attributionBoxRadius := canvasWidth * 0.01

		// fb.drawRoundedBox(ctx, attributionBoxTop, attributionBoxBottom, attributionBoxLeft, attributionBoxRight, attributionBoxRadius)

		attributionFontSize := attributionBoxHeight * 0.8
		// attributionTextMaxWidth := (attributionBoxRight - attributionBoxLeft) * 0.9
		attributionTextMaxHeight := (attributionBoxTop - attributionBoxBottom) * 0.9

		attributionString := fmt.Sprintf("%s<BR>Generated using \"%s\" by %s<BR>netrunner-alt-gen %s", card.Attributes.Title, fb.Algorithm, fb.Designer, fb.Version)

		if fb.Algorithm == "" && fb.Designer == "" {
			attributionString = fmt.Sprintf("%s<BR>Layout by netrunner-alt-gen %s", card.Attributes.Title, fb.Version)
		}
		if fb.Algorithm == "" && fb.Designer != "" {
			attributionString = fmt.Sprintf("%s<BR>Design by %s<BR>Layout by netrunner-alt-gen %s", card.Attributes.Title, fb.Designer, fb.Version)
		}

		// attributionTextX := attributionBoxLeft + ((attributionBoxRight - attributionBoxLeft) * 0.05)
		// attributionTextY := (attributionBoxTop - (attributionBoxHeight-attributionText.Bounds().H)*0.5)
		// ctx.DrawText(attributionTextX, attributionTextY, attributionText)

		cliFontSize := attributionFontSize * 0.8
		cliTextMaxHeight := canvasHeight * 0.7
		cliString := getCLIText()

		attributionText := fb.getVerticalFittedText(ctx, attributionString, attributionFontSize, 0, attributionTextMaxHeight, canvas.Center)
		cliText := fb.getVerticalFittedText(ctx, cliString, cliFontSize, 0, cliTextMaxHeight, canvas.Left)

		cliBoxWidth := math.Max(cliText.Bounds().W, attributionText.Bounds().W) * 1.1
		cliBoxHeight := cliText.Bounds().H + (attributionText.Bounds().H)
		cliBoxMarginTB := attributionText.Bounds().H * 0.2
		if len(cliString) == 0 {
			cliBoxHeight += cliBoxMarginTB * 2
		} else {
			cliBoxHeight += cliBoxMarginTB * 3
		}

		cliBoxBottom := canvasHeight * 0.0591
		cliBoxTop := cliBoxBottom + cliBoxHeight
		cliBoxRight := (canvasWidth / 2) + cliBoxWidth*0.5
		cliBoxLeft := (canvasWidth / 2) - cliBoxWidth*0.5
		cliBoxRadius := canvasWidth * 0.01
		cliBoxMiddle := cliBoxLeft + cliBoxWidth*0.5

		attributionTextX := cliBoxMiddle
		attributionTextY := cliBoxTop - cliBoxMarginTB

		cliTextX := cliBoxMiddle - (cliText.Bounds().W)*0.5
		cliTextY := cliBoxBottom + cliText.Bounds().H + cliBoxMarginTB

		fb.drawRoundedBox(ctx, cliBoxTop, cliBoxBottom, cliBoxLeft, cliBoxRight, cliBoxRadius)

		ctx.DrawText(attributionTextX, attributionTextY, attributionText)
		ctx.DrawText(cliTextX, cliTextY, cliText)

		return nil
	})
}

func (fb FrameBasic) drawRoundedBox(ctx *canvas.Context, top, bottom, left, right, radius float64) {

	strokeWidth := getStrokeWidth(ctx)

	ctx.Push()
	ctx.SetFillColor(fb.getColorBG())
	ctx.SetStrokeColor(fb.getColorBorder())
	ctx.SetStrokeWidth(strokeWidth)

	path := &canvas.Path{}
	path.MoveTo(left, top-radius)
	path.QuadTo(left, top, left+radius, top)
	path.LineTo(right-radius, top)
	path.QuadTo(right, top, right, top-radius)
	path.LineTo(right, bottom+radius)
	path.QuadTo(right, bottom, right-radius, bottom)
	path.LineTo(left+radius, bottom)
	path.QuadTo(left, bottom, left, bottom+radius)
	path.Close()

	ctx.DrawPath(0, 0, path)
	ctx.Pop()
}

func getCLIText() string {
	var args []string

	var isFlag bool
	var thisOpt string
	for _, arg := range os.Args {

		if len(arg) < 1 {
			thisOpt = ""
			continue
		}

		if arg == "-o" {
			thisOpt = ""
			continue
		}

		if arg == "--make-back" {
			thisOpt = ""
			continue
		}

		if arg[0] == '-' {
			isFlag = true
		}

		if !isFlag {
			thisOpt = ""
			continue
		}

		if arg[0] != '-' {
			isFlag = false
			args = append(args, thisOpt+" "+arg)
		} else {
			thisOpt = arg
		}

	}

	return strings.Join(args, "<BR>")
}
