package image

import (
	"github.com/spf13/cobra"
)

var ResizeCmd = &cobra.Command{
	Use: "resize",
	Short: "Resize images inside a folder",
	Long: `Resize all the images inside a folder. For example:
	got resize *.jpg -h 1080 -w 1920`,
	Run: executeResize,
}

func init() {
	// generate

}

func executeResize(cmd *cobra.Command, args []string) {
	// var height, _ = cmd.Flags().IntP("height", "y", 1080, "height")
	// var width, _ = cmd.Flags().IntP("width", "x", 1920, "width")

	// if height {

	// }

	// for _, file := range args {
	// 	fmt.Println("Processing file:", file)
	// }
}

