package image

import "github.com/spf13/cobra"

var ImageCmd = &cobra.Command{
	Use: "image",
	Short: "Handles image",
	Long: `Handles image. For example:

	got image`,
}

func init() {
	ImageCmd.AddCommand(resizeCmd)
	ImageCmd.AddCommand(filterCmd)
}