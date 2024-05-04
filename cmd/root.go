package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// old values

// const canvasWidth = 3288.0
// const canvasHeight = 4488.0
// const cardWidth = 3064.0
// const cardHeight = 4212.0

// based on NSG pdf sizes

// const canvasWidth = 3199.0
// const canvasHeight = 4432.0
// const cardWidth = 2975.0
// const cardHeight = 4156.0

// based on MPC template

const canvasWidth = 3264.0
const canvasHeight = 4450.0
const cardWidth = 2976.0
const cardHeight = 4152.0
const safeWidth = 2736.0
const safeHeight = 3924.0

// real card MM, doesn't work currently, need higehr res for
// generation

// const canvasWidth = 69.35
// const canvasHeight = 94.35
// const cardWidth = 63.0
// const cardHeight = 88.0

var (
	drawMarginLines, makeBack                                         bool
	outputDir                                                         string
	baseColor, walkerColor1, walkerColor2, walkerColor3, walkerColor4 string
	skipFlavor                                                        bool
	flavorText, flavorAttribution                                     string
	textBoxFactor                                                     float64

	frame, frameColorBackground, frameColorBorder, frameColorText,
	frameColorInfluenceBG, frameColorStrengthBG, frameColorFactionBG,
	frameColorInfluenceLimitBG, frameColorMinDeckBG string

	// netspace
	netspaceWalkersMin, netspaceWalkersMax         int
	netspaceColorBG                                string
	gridColor1, gridColor2, gridColor3, gridColor4 string
	gridPercent                                    float64

	// image
	designer string

	// set by ldflags
	version string = "local"
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&drawMarginLines, "draw-margin-lines", "", false, `Draw bleed and "safe area" lines`)
	rootCmd.PersistentFlags().BoolVarP(&makeBack, "make-back", "", false, `Also create a file for a card back. Uses "${frame}-back" as frame name.`)
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "output", `Output directory name`)

	rootCmd.PersistentFlags().StringVarP(&flavorText, "flavor", "", "", `Flavor text to add to the generated card`)
	rootCmd.PersistentFlags().StringVarP(&flavorAttribution, "flavor-attribution", "", "", `Flavor text attribution to add to the generated card, for "quotes"`)
	rootCmd.PersistentFlags().BoolVarP(&skipFlavor, "skip-flavor", "", false, `Don't render default flavor text`)
	rootCmd.PersistentFlags().StringVarP(&baseColor, "base-color", "c", "", `Alternate base color for the card, defaults to pre-defined faction colors`)
	rootCmd.PersistentFlags().Float64VarP(&textBoxFactor, "text-box-height", "", 33.3, `Percentage of total card height taken up by the main text box`)

	rootCmd.PersistentFlags().StringVarP(&frame, "frame", "f", "basic", `Frame to draw, use "none" to skip drawing a frame`)
	rootCmd.PersistentFlags().StringVarP(&frameColorBackground, "frame-color-background", "", "1c1c1c99", `Background color for frame text boxes`)
	rootCmd.PersistentFlags().StringVarP(&frameColorBorder, "frame-color-border", "", "dcdccc", `Border color for frame text boxes`)
	rootCmd.PersistentFlags().StringVarP(&frameColorText, "frame-color-text", "", "dcdccc", `Text color for frame text boxes`)
	rootCmd.PersistentFlags().StringVarP(&frameColorFactionBG, "frame-color-faction-bg", "", "1c1c1c", `Background color for the faction symbol`)
	rootCmd.PersistentFlags().StringVarP(&frameColorInfluenceBG, "frame-color-influence-bg", "", "",
		`Background color for the influence cost indicator
Defaults to pre-defined faction colors or specified base color
If set to "faction", it will use the faction color regardless of the base color`)
	rootCmd.PersistentFlags().StringVarP(&frameColorStrengthBG, "frame-color-strength-bg", "", "",
		`Background color for the strength bubble on ice and programs
Defaults to pre-defined faction colors or specified base color
If set to "faction", it will use the faction color regardless of the base color`)
	rootCmd.PersistentFlags().StringVarP(&frameColorInfluenceLimitBG, "frame-color-influence-limit-bg", "", "3f3f3f",
		`Background color for the influence limit indicator on identities`)
	rootCmd.PersistentFlags().StringVarP(&frameColorMinDeckBG, "frame-color-min-deck-bg", "", "",
		`Background color for the min deck size on identites
Defaults to pre-defined faction colors or specified base color
If set to "faction", it will use the faction color regardless of the base color`)

	netspaceCmd.Flags().IntVarP(&netspaceWalkersMin, "min-walkers", "m", 3000, `Minimum amount of walkers`)
	netspaceCmd.Flags().IntVarP(&netspaceWalkersMax, "max-walkers", "M", 10000, `Maximum amount of walkers`)
	netspaceCmd.Flags().StringVarP(&netspaceColorBG, "color-bg", "", "", `Background color for the generated art, defaults to --base-color value`)
	netspaceCmd.Flags().StringVarP(&walkerColor1, "walker-color-1", "", "", `Alternate walker color for the card, defaults to pre-defined faction color analogue +10 - +30`)
	netspaceCmd.Flags().StringVarP(&walkerColor2, "walker-color-2", "", "", `Alternate walker color for the card, defaults to pre-defined faction color analogue -10 - -30`)
	netspaceCmd.Flags().StringVarP(&walkerColor3, "walker-color-3", "", "", `Alternate walker color for the card, defaults to pre-defined faction color analogue +30 - +50`)
	netspaceCmd.Flags().StringVarP(&walkerColor3, "walker-color-4", "", "", `Alternate walker color for the card, defaults to pre-defined faction color analogue -30 - -50`)
	netspaceCmd.Flags().StringVarP(&gridColor1, "grid-color-1", "", "",
		`Alternate grid color for the grid pattern on the card, defaults to --alt-color-1, will be randomly desaturated by algorithm`)
	netspaceCmd.Flags().StringVarP(&gridColor2, "grid-color-2", "", "",
		`Alternate grid color for the grid pattern on the card, defaults to --alt-color-2, will be randomly desaturated by algorithm`)
	netspaceCmd.Flags().StringVarP(&gridColor3, "grid-color-3", "", "",
		`Alternate grid color for the grid pattern on the card, defaults to --alt-color-3, will be randomly desaturated by algorithm`)
	netspaceCmd.Flags().StringVarP(&gridColor4, "grid-color-4", "", "",
		`Alternate grid color for the grid pattern on the card, defaults to --alt-color-4, will be randomly desaturated by algorithm`)
	netspaceCmd.PersistentFlags().Float64VarP(&gridPercent, "grid-percent", "", -1, `Percentage of total walkers that will run on a grid`)

	imageCmd.Flags().StringVarP(&designer, "designer", "", "", `Name of the designer for the card back attribution`)

	rootCmd.AddCommand(netspaceCmd)
	rootCmd.AddCommand(emptyCmd)
	rootCmd.AddCommand(imageCmd)
}

var rootCmd = &cobra.Command{
	Use:   "netrunner-alt-gen",
	Short: "netrunner-alt-gen generates alt arts for Netrunner",
	Long: `A generative art tool to create alternate art cards with tournament legal frames for Netrunner.
  Complete documentation is available at https://github.com/mangofeet/netrunner-alt-gen`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Version:", version)
		cmd.Usage()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
