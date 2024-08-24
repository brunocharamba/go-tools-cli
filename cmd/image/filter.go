package image

import (
	"log"

	"github.com/spf13/cobra"
)

var filterCmd = &cobra.Command{
	Use: "filter",
	Short: "Apply filters to images inside a folder",
	Long: `Apply filters to all the images inside a folder. For example:
	got image filter *.jpg -f grayscale`,
	Run: executeFilter,
}

func init() {
	filterCmd.Flags().StringP("filter", "f", "jpg", "Extension of the output. Default: jpg. Options: jpg, png.")
}

func executeFilter(cmd *cobra.Command, args []string) {
	filter, _ := cmd.Flags().GetString("filter")


	log.Println(filter)
}