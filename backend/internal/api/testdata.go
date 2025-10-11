package api

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
)

// createTestImage creates a minimal valid image for testing in the specified format
func createTestImage(mimeType string) []byte {
	// Create a 1x1 pixel image
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255}) // Red pixel
	
	var buf bytes.Buffer
	
	switch mimeType {
	case "image/jpeg":
		jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	case "image/png":
		png.Encode(&buf, img)
	default:
		// Default to PNG
		png.Encode(&buf, img)
	}
	
	return buf.Bytes()
}

