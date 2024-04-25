package util

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/chai2010/webp"
	"golang.org/x/image/tiff"
)

func GetImageFormat(img image.Image) string {
	switch img.(type) {
	case *image.RGBA:
		return "png"
	case *image.NRGBA:
		return "png"
	case *image.Gray:
		return "jpeg"
	case *image.YCbCr:
		return "jpeg"
	default:
		return "unknown"
	}
}

func EncodeImage(w io.Writer, img image.Image, format string, quality int) error {
	switch strings.ToLower(format) {
	case "png":
        // Map quality to PNG compression levels
        var encoder png.Encoder
        switch {
        case quality < 25:
            encoder.CompressionLevel = png.BestSpeed
        case quality < 50:
            encoder.CompressionLevel = png.NoCompression
        case quality < 75:
            encoder.CompressionLevel = png.DefaultCompression
        default:
            encoder.CompressionLevel = png.BestCompression
        }
        return encoder.Encode(w, img)
	case "gif":
		return gif.Encode(w, img, nil)
	case "tiff":
		return tiff.Encode(w, img, nil)
	case "webp":
		return webp.Encode(w, img, &webp.Options{Lossless: false, Quality: float32(quality)})
	default:
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	}
}
