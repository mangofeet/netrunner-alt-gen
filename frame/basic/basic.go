package basic

import (
	"image/color"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type FrameBasic struct {
	Flavor, FlavorAttribution string
	Color                     *color.RGBA
}

func (fb FrameBasic) getAdditionalText() []additionalText {
	var extra []additionalText
	if fb.Flavor != "" {
		extra = append(extra, additionalText{
			content: fb.Flavor,
			align:   canvas.Left,
		})
	}
	if fb.FlavorAttribution != "" {
		extra = append(extra, additionalText{
			content: fb.FlavorAttribution,
			align:   canvas.Right,
		})
	}

	return extra
}

func (fb FrameBasic) getColor(card *nrdb.Printing) color.RGBA {

	if fb.Color != nil {
		return *fb.Color
	}

	baseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	return art.Darken(baseColor, 0.811)

}