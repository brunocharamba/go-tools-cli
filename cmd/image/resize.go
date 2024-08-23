package image

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/image/draw"
)

var ImageCmd = &cobra.Command{
	Use: "image",
	Short: "Handles image",
	Long: `Handles image. For example:

	got image`,
}

var resizeCmd = &cobra.Command{
	Use: "resize",
	Short: "Resize images inside a folder",
	Long: `Resize all the images inside a folder. For example:
	got resize *.jpg -h 1080 -w 1920`,
	Run: executeResize,
}

func init() {
	resizeCmd.Flags().IntP("width", "x", 0, "Width of the output image")
	resizeCmd.Flags().IntP("height", "y", 0, "Height of the output image")
	resizeCmd.Flags().Float32P("percentage", "p", 0, "Percentage of the output image")
	ImageCmd.AddCommand(resizeCmd)
}

func executeResize(cmd *cobra.Command, args []string) {
	width, _ := cmd.Flags().GetInt("width")
	height, _ := cmd.Flags().GetInt("height")
	percentage, _ := cmd.Flags().GetFloat32("percentage")

	if len(args) == 0 {
		panic("No images found.")
	}

	path := "output"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	for _, arg := range args {
		parseImage(width, height, percentage, arg)
	}
}

func parseImage(width int, height int, percentage float32, arg string) {
// Open the input image file
	file, err := os.Open(arg)
	if err != nil {
		panic(fmt.Sprintf("Failed to open image: %v", err))
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode image: %v", err))
	}

	var destination image.RGBA


	if percentage > 0 {
		// parse as percentage
		newWidth := int(float32(img.Bounds().Dx()) * percentage)
		newHeight := int(float32(img.Bounds().Dy()) * percentage)
		destination = *image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	} else if width > 0 || height > 0 {
		// parse as default h/w
		newWidth, newHeight := getNewBounds(width, height, &img)
		destination = *image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	} else {
		// run default
		panic("Pass or a percentage or values for width and/or height")
	}

	// Perform the scaling using BiLinear interpolation
	draw.BiLinear.Scale(&destination, destination.Bounds(), img, img.Bounds(), draw.Over, nil)

	outputName := fmt.Sprintf("%v-resized.jpg", strings.TrimSuffix(arg, filepath.Ext(arg)))

	// Create a new file to save the resized image
	out, err := os.Create(fmt.Sprintf("output/%v", outputName))
	if err != nil {
		log.Fatalf("Failed to create output image file: %v", err)
	}
	defer out.Close()

	// Encode the resized image to the file in JPEG format
	err = jpeg.Encode(out, &destination, nil)
	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}

	log.Printf("Image '%v' successfully! \n", outputName)
}

func getNewBounds(width int, height int, img *image.Image) (int, int) {
	originalWidth := (*img).Bounds().Dx()
	originalHeight := (*img).Bounds().Dy()

	if width > 0 && height > 0 {
		return width, height
	} else if width > 0 {
		ratio := float64(width) / float64(originalWidth)
		return width, int(float64(originalHeight) * ratio)
	}

	ratio := float64(height) / float64(originalHeight)
	return int(float64(originalWidth) * ratio), height
}