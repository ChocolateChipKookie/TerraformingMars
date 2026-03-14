package imageutil

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

const MaxWidth = 1920
const MaxHeight = 1920
const WebPQuality = 80

// ProcessImage resizes and compresses an image to WebP format
func ProcessImage(imageData []byte, originalMimeType string) ([]byte, string, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %v", err)
	}

	resizedImg := resizeImage(img)
	return encodeAsWebP(resizedImg)
}

func resizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width <= MaxWidth && height <= MaxHeight {
		return img
	}

	newWidth, newHeight := calculateNewSize(width, height, MaxWidth, MaxHeight)
	return resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)
}

func calculateNewSize(width, height, maxWidth, maxHeight int) (int, int) {
	if width <= maxWidth && height <= maxHeight {
		return width, height
	}

	scaleX := float64(maxWidth) / float64(width)
	scaleY := float64(maxHeight) / float64(height)

	// Use the smaller scale to maintain aspect ratio
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	newWidth := int(float64(width) * scale)
	newHeight := int(float64(height) * scale)

	return newWidth, newHeight
}

func encodeAsWebP(img image.Image) ([]byte, string, error) {
	webpData, err := webp.EncodeRGB(img, WebPQuality)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode WebP: %v", err)
	}

	return webpData, "image/webp", nil
}

