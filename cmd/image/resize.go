package image

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/image/draw"
)

var VALID_EXTS = []string{".jpg", ".png"}

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
	resizeCmd.Flags().StringP("extension", "e", "jpg", "Extension of the output. Default: jpg. Options: jpg, png.")
}

func executeResize(cmd *cobra.Command, args []string) {
	width, _ := cmd.Flags().GetInt("width")
	height, _ := cmd.Flags().GetInt("height")
	percentage, _ := cmd.Flags().GetFloat32("percentage")
	extension, _ := cmd.Flags().GetString("extension")

	if len(args) == 0 {
		panic("No images found.")
	}

	extension = handleExtension(extension)

	path := "output"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	count := len(args)
	for i, arg := range args {
		outputName := parseImage(width, height, percentage, extension, arg)
		log.Printf("[%v/%v] Image '%v' successfully! \n", i + 1, count, outputName)
	}
}

func handleExtension(extension string) string{
	// add '.' if it was not added
	if extension[0] != '.' {
		extension = "." + extension
	}

	// check output extension
	if !slices.Contains(VALID_EXTS, extension) {
		log.Printf("Extension '%v' isn't valid %v. Defaulting to the same extension \n", extension, VALID_EXTS)
	}

	return extension
}

func parseImage(width int, height int, percentage float32, extension string, arg string) string {
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
		panic("Pass or a percentage or values for width and/or height")
	}

	// perform the scaling using BiLinear interpolation
	draw.BiLinear.Scale(&destination, destination.Bounds(), img, img.Bounds(), draw.Over, nil)

	// generate file
	outputFile, outputName := generateFile(arg, extension)

	// encode the resized image to the file
	switch extension {
		case ".jpg":
			encodeJpgImage(outputFile, &destination)
		case ".png":
			encodePngImage(outputFile, &destination)
	}

	return outputName
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

func generateFile(originalFileName string, ext string) (*os.File, string) {
	originalExt := filepath.Ext(originalFileName)
	outputName := fmt.Sprintf("%v-resized%v", strings.TrimSuffix(originalFileName, originalExt), ext)

	// Create a new file to save the resized image
	out, err := os.Create(fmt.Sprintf("output/%v", outputName))
	if err != nil {
		log.Fatalf("Failed to create output image file: %v", err)
	}

	return out, outputName
}

func encodeJpgImage(outputFile *os.File, outputImage image.Image) {
	err := jpeg.Encode(outputFile, outputImage, nil)
	defer outputFile.Close()

	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}
}

func encodePngImage(outputFile *os.File, outputImage image.Image) {
	err := png.Encode(outputFile, outputImage)

	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}
}