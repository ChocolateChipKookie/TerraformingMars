package imageutil

import (
	"bytes"
	"fmt"
	"image"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

const MaxWidth = 1920
const MaxHeight = 1920
const WebPQuality = 80

// ProcessImage resizes and compresses an image to WebP format
func ProcessImage(imageData []byte, originalMimeType string) ([]byte, string, error) {
	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %v", err)
	}

	// Resize if needed
	resizedImg := resizeImage(img)

	// Encode as WebP
	return encodeAsWebP(resizedImg)
}

// resizeImage resizes image if it exceeds max dimensions
func resizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Check if resize is needed
	if width <= MaxWidth && height <= MaxHeight {
		return img
	}

	// Calculate new dimensions maintaining aspect ratio
	newWidth, newHeight := calculateNewSize(width, height, MaxWidth, MaxHeight)

	// Use Lanczos resampling for high quality resize
	return resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)
}

// calculateNewSize calculates new dimensions maintaining aspect ratio
func calculateNewSize(width, height, maxWidth, maxHeight int) (int, int) {
	if width <= maxWidth && height <= maxHeight {
		return width, height
	}

	// Calculate scaling factors
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

// encodeAsWebP encodes image as WebP
func encodeAsWebP(img image.Image) ([]byte, string, error) {
	// Encode as WebP with specified quality
	webpData, err := webp.EncodeRGB(img, WebPQuality)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode WebP: %v", err)
	}
	
	return webpData, "image/webp", nil
}