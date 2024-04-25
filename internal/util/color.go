package util

import (
	"image/color"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"
)

// parseColor parses a hex color string and returns a color.Color.
func ParseColor(colorStr string) (color.RGBA, error) {
    if colorStr == "transparent" {
        return color.RGBA{0, 0, 0, 0}, nil
    }

    // Check if it's a named color
    if namedColor, ok := colornames.Map[strings.ToLower(colorStr)]; ok {
        return color.RGBA{
            R: uint8(namedColor.R),
            G: uint8(namedColor.G),
            B: uint8(namedColor.B),
            A: uint8(namedColor.A),
        }, nil
    }

    // Parse as a hexadecimal color
    if colorStr[0] == '#' {
        colorStr = colorStr[1:]
    }
    c, err := strconv.ParseUint(colorStr, 16, 64)
    if err != nil {
        return color.RGBA{}, err
    }
    return color.RGBA{
        R: uint8(c >> 24),
        G: uint8(c >> 16),
        B: uint8(c >> 8),
        A: uint8(c),
    }, nil
}
