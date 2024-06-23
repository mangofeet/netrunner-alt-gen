package basic

import (
	"image/color"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type FrameBasic struct {
	Flavor, FlavorAttribution    string
	TextBoxHeightFactor          *float64
	Designer, Version, Algorithm string

	ColorBG, ColorBorder, ColorText,
	ColorTextStrength, ColorInfluencePips,
	ColorInfluenceLimitBG, ColorMinDeckBG,
	ColorInfluenceBG, ColorStrengthBG, ColorFactionBG *color.RGBA
}

func (fb FrameBasic) getAdditionalText() []additionalText {
	var extra []additionalText
	if fb.Flavor != "" {
		extra = append(extra, additionalText{
			content:  fb.Flavor,
			textType: additionalTextTypeFlavor,
			align:    canvas.Left,
		})
	}
	if fb.FlavorAttribution != "" {
		extra = append(extra, additionalText{
			content:  fb.FlavorAttribution,
			textType: additionalTextTypeAttribution,
			align:    canvas.Right,
		})
	}

	return extra
}

func (fb FrameBasic) getColorBG() color.RGBA {

	if fb.ColorBG != nil {
		return *fb.ColorBG
	}

	return colorDefaultBG
}

func (fb FrameBasic) getColorBorder() color.RGBA {

	if fb.ColorBorder != nil {
		return *fb.ColorBorder
	}

	return colorDefaultText
}

func (fb FrameBasic) getColorTextStrength() color.RGBA {
	if fb.ColorTextStrength != nil {
		return *fb.ColorTextStrength
	}
	return fb.getColorText()
}

func (fb FrameBasic) getColorInfluencePips() color.RGBA {
	if fb.ColorInfluencePips != nil {
		return *fb.ColorInfluencePips
	}
	return fb.getColorText()
}

func (fb FrameBasic) getColorText() color.RGBA {

	if fb.ColorText != nil {
		return *fb.ColorText
	}

	return colorDefaultText
}

func (fb FrameBasic) getColorInfluenceBG(card *nrdb.Printing) color.RGBA {

	if fb.ColorInfluenceBG != nil {
		return *fb.ColorInfluenceBG
	}

	return art.GetFactionBaseColor(card.Attributes.FactionID)
}

func (fb FrameBasic) getColorStrengthBG(card *nrdb.Printing) color.RGBA {

	if fb.ColorStrengthBG != nil {
		return *fb.ColorStrengthBG
	}

	return art.GetFactionBaseColor(card.Attributes.FactionID)

}

func (fb FrameBasic) getColorFactionBG() color.RGBA {

	if fb.ColorFactionBG != nil {
		return *fb.ColorFactionBG
	}

	return colorDefaultOpaqueBG
}

func (fb FrameBasic) getColorInfluenceLimitBG() color.RGBA {

	if fb.ColorInfluenceLimitBG != nil {
		return *fb.ColorInfluenceLimitBG
	}

	return colorDefaultOpaqueBG
}

func (fb FrameBasic) getColorMinDeckBG(card *nrdb.Printing) color.RGBA {

	if fb.ColorMinDeckBG != nil {
		return *fb.ColorMinDeckBG
	}

	return art.GetFactionBaseColor(card.Attributes.FactionID)
}
