package basic

import (
	"image/color"

	"github.com/tdewolff/canvas"
)

type FrameBasic struct {
	Flavor, FlavorAttribution string
	Color                     color.RGBA
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
