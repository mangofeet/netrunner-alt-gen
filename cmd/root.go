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
	drawMarginLines               bool
	frame, outputDir              string
	baseColor                     string
	skipFlavor                    bool
	flavorText, flavorAttribution string
	textBoxFactor                 float64

	netspaceWalkersMin, netspaceWalkersMax int

	// set by ldflags
	version string = "local"
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&drawMarginLines, "draw-margin-lines", "", false, `Draw bleed and "safe area" lines`)
	rootCmd.PersistentFlags().StringVarP(&frame, "frame", "f", "basic", `Frame to draw, use "none" to skip drawing a frame`)
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "output", `Output directory name`)

	rootCmd.PersistentFlags().StringVarP(&flavorText, "flavor", "", "", `Flavor text to add to the generated card`)
	rootCmd.PersistentFlags().StringVarP(&flavorAttribution, "flavor-attribution", "", "", `Flavor text attribution to add to the generated card, for "quotes"`)
	rootCmd.PersistentFlags().BoolVarP(&skipFlavor, "skip-flavor", "", false, `Don't render default flavor text`)
	rootCmd.PersistentFlags().StringVarP(&baseColor, "base-color", "c", "", `Alternate base color for the card, defaults to pre-defined faction colors`)
	rootCmd.PersistentFlags().Float64VarP(&textBoxFactor, "text-box-height", "", 33.3, `Percentage of total card height taken up by the main text box`)

	netspaceCmd.Flags().IntVarP(&netspaceWalkersMin, "min-walkers", "m", 3000, `Minimum amount of walkers`)
	netspaceCmd.Flags().IntVarP(&netspaceWalkersMax, "max-walkers", "M", 10000, `Maximum amount of walkers`)

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
