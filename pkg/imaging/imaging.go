package imaging

import (
	"image"
	"image/color"
	"math"
)

// Helper function to convert RGB to HSL
func RgbToHSL(r, g, b uint8) (float64, float64, float64) {
    rf := float64(r) / 255.0
    gf := float64(g) / 255.0
    bf := float64(b) / 255.0
    max := math.Max(rf, math.Max(gf, bf))
    min := math.Min(rf, math.Min(gf, bf))
    l := (max + min) / 2.0

    var h, s float64
    if max == min {
        h = 0
        s = 0
    } else {
        d := max - min
        if l > 0.5 {
            s = d / (2.0 - max - min)
        } else {
            s = d / (max + min)
        }
        switch max {
        case rf:
            h = (gf - bf) / d
            if g < b {
                h += 6
            }
        case gf:
            h = (bf - rf)/d + 2
        case bf:
            h = (rf - gf)/d + 4
        }
        h /= 6
    }
    return h, s, l
}
// Helper function to convert HSL to RGB
func HslToRGB(h, s, l float64) (uint8, uint8, uint8) {
    var r, g, b float64

    if s == 0 {
        r, g, b = l, l, l // achromatic
    } else {
        var q float64
        if l < 0.5 {
            q = l * (1 + s)
        } else {
            q = l + s - l*s
        }
        p := 2*l - q
        r = HueToRGB(p, q, h+1.0/3.0)
        g = HueToRGB(p, q, h)
        b = HueToRGB(p, q, h-1.0/3)
    }

    return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}

func HueToRGB(p, q, t float64) float64 {
    if t < 0 {
        t += 1
    }
    if t > 1 {
        t -= 1
    }
    if t < 1.0/6 {
        return p + (q-p)*6*t
    }
    if t < 1.0/2 {
        return q
    }
    if t < 2.0/3 {
        return p + (q-p)*(2.0/3-t)*6
    }
    return p
}

// Adjust the hue of the entire image
func AdjustHue(img image.Image, deltaHue float64) image.Image {
    bounds := img.Bounds()
    dst := image.NewNRGBA(bounds)
    deltaHue = deltaHue / 360.0 // Normalize deltaHue to [0, 1)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            // Convert RGB from uint32 to uint8
            red, green, blue := uint8(r>>8), uint8(g>>8), uint8(b>>8)
            h, s, l := RgbToHSL(red, green, blue)
            h += deltaHue
            if h > 1 {
                h -= 1
            } else if h < 0 {
                h += 1
            }
            // Compute the new RGB values after hue adjustment
            nr, ng, nb := HslToRGB(h, s, l)
            // Set the new color with correct types
            dst.Set(x, y, color.NRGBA{R: nr, G: ng, B: nb, A: uint8(a >> 8)})
        }
    }
    return dst
}
