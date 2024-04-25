package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/siinghd/golang-imageoptimizer/internal/httputil"
	"github.com/siinghd/golang-imageoptimizer/internal/util"
	pkgimaging "github.com/siinghd/golang-imageoptimizer/pkg/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)


func ProcessImageWithParams(img image.Image, params util.QueryParams) (image.Image, error) {
	// Resize image if width or height is provided
	if params.Width > 0 || params.Height > 0 {
		width := params.Width
		height := params.Height
		if width == 0 {
			width = img.Bounds().Dx()
		}
		if height == 0 {
			height = img.Bounds().Dy()
		}
		img = imaging.Resize(img, width, height, imaging.Lanczos)
	}

	// Apply fit mode
	if params.Fit != "" {
		switch params.Fit {
		case "contain":
			img = imaging.Fit(img, params.Width, params.Height, imaging.Lanczos)
		case "fill":
			img = imaging.Fill(img, params.Width, params.Height, imaging.Center, imaging.Lanczos)
		case "inside":
			img = imaging.Fit(img, params.Width, params.Height, imaging.Lanczos)
		case "outside":
			img = imaging.Resize(img, params.Width, params.Height, imaging.Lanczos)
		}
	}

	// Apply background color if provided
	if params.Background != "" {
		bgColor,_ := util.ParseColor(params.Background)
		img = imaging.OverlayCenter(imaging.New(img.Bounds().Dx(), img.Bounds().Dy(), bgColor), img, 1.0)
	}

	// Apply blur if provided
	if params.Blur > 0 {
		img = imaging.Blur(img, params.Blur)
	}

	// Apply gamma correction if provided
	if params.Gamma > 0 {
		img = imaging.AdjustGamma(img, params.Gamma)
	}

	// Apply image modulation if provided
	if params.Modulate != "" {
		parts := strings.Split(params.Modulate, ",")
		if len(parts) == 3 {
			brightness :=  util.ParseFloat(parts[0], 1.0)
			saturation :=  util.ParseFloat(parts[1], 1.0)
			hue :=  util.ParseFloat(parts[2], 0.0)
			img = imaging.AdjustBrightness(img, brightness)
			img = imaging.AdjustSaturation(img, saturation)
			img = pkgimaging.AdjustHue(img, hue)
		}
	}

	// Apply sharpening if provided
	if params.Sharpen > 0 {
		img = imaging.Sharpen(img, params.Sharpen)
	}
// Apply device pixel ratio (DPR)
	if params.DPR != 1 {
		width := int(float64(img.Bounds().Dx()) * params.DPR)
		height := int(float64(img.Bounds().Dy()) * params.DPR)
		img = imaging.Resize(img, width, height, imaging.Lanczos)
	}

	// Apply contain background color
	if params.ContainBgColor != "" {
		bgColor, err := util.ParseColor(params.ContainBgColor)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse contain background color")
		}
		img = imaging.Resize(img, params.Width, params.Height, imaging.Lanczos)
		background := imaging.New(params.Width, params.Height, bgColor)
		img = imaging.PasteCenter(background, img)
	}

	// Apply without enlargement
	if params.WithoutEnlarge {
		originalWidth, originalHeight := img.Bounds().Dx(), img.Bounds().Dy()
		if params.Width > originalWidth || params.Height > originalHeight {
			img = imaging.Resize(img, originalWidth, originalHeight, imaging.Lanczos)
		}
	}

	// Apply tint
	if params.Tint != "" {
		tintColor, err := util.ParseColor(params.Tint)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse tint color")
		}
		img = imaging.Overlay(img, &image.Uniform{tintColor}, image.Pt(0, 0), 0.5)
	}

	// Apply interlacing
	if params.Interlace {
		// TODO: Implement interlacing (not supported by the imaging library)
	}

	// Apply pages and page selection
	if params.Pages != 0 || params.Page != 0 {
		// TODO: Implement multi-page support (not supported by the imaging library)
	}
	return img, nil
}

func RenderTextToImage(text string, options util.QueryParams) (image.Image, error) {
    width := options.Width
    if width == 0 {
        width = 800 // Default width
    }

    height := options.Height
    if height == 0 {
        height = 800 // Default height
    }

    backgroundColor := options.Background
    if backgroundColor == "" {
        backgroundColor = "white"
    }

    textColor := options.TextColor
    if textColor == "" {
        textColor = "black"
    }

    fontSize := options.FontSize
    if fontSize == 0 {
        fontSize = 48 // Default font size
    }

    fontFamily := options.FontFamily
    if fontFamily == "" {
        fontFamily = "Arial" // Default font family
    }

    textAlign := options.TextAlign
    if textAlign == "" {
        textAlign = "center" // Default text alignment
    }

    roundedCorners := options.RoundedCorners

    cornerRadius := options.CornerRadius
    if cornerRadius == 0 {
        cornerRadius = 20 // Default corner radius
    }

    textBaseline := options.TextBaseline
    if textBaseline == "" {
        textBaseline = "middle"
    }

    dc := gg.NewContext(width, height)
	// TODO: this does not work 
    if roundedCorners {
        dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), float64(cornerRadius))
        dc.Clip()
    }

    // Parse background color
    bgColor, err := util.ParseColor(backgroundColor)
    if err != nil {
        return nil, fmt.Errorf("failed to parse background color: %v", err)
    }
    dc.SetColor(bgColor)
    dc.Clear()

    // Parse text color
    txtColor, err := util.ParseColor(textColor)
    if err != nil {
        return nil, fmt.Errorf("failed to parse text color: %v", err)
    }
    dc.SetColor(txtColor)

    var face font.Face
    if fontFamily != "" {
        // Attempt to load the user-provided font
        loadedFont, err := gg.LoadFontFace(fontFamily, float64(fontSize))
        if err == nil {
            face = loadedFont
        }
    }

    if face == nil {
        // If no user-provided font was specified or if loading the font failed, use the default "Go Regular" font
        font, err := truetype.Parse(goregular.TTF)
        if err != nil {
            return nil, errors.Wrap(err, "failed to load default font")
        }
        face = truetype.NewFace(font, &truetype.Options{Size: float64(fontSize)})
    }

    dc.SetFontFace(face)

    var textY float64
	switch textBaseline {
		case "top":
			textY = float64(fontSize) / 2
		case "hanging":
			
			textY = float64(fontSize) - float64(face.Metrics().Ascent.Floor())
		case "middle":
			textY = float64(height) / 2
		case "alphabetic":
			textY = float64(height)/2 +  float64(face.Metrics().Ascent.Floor())/2
		case "ideographic":
			
			textY = float64(height)/2 + float64(face.Metrics().Ascent.Floor())
		case "bottom":
			textY = float64(height) - float64(fontSize)/2
		default:
			textY = float64(height) / 2
		}
    switch textAlign {
    case "left":
        dc.DrawStringAnchored(text, 0, textY, 0, 0.5)
    case "center":
        dc.DrawStringAnchored(text, float64(width)/2, textY, 0.5, 0.5)
    case "right":
        dc.DrawStringAnchored(text, float64(width), textY, 1, 0.5)
    }

    return dc.Image(), nil
}

func FetchImage(url, defaultURL string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		if defaultURL != "" {
			resp, err = http.Get(defaultURL)
			if err != nil {
				return nil, errors.Wrap(err, "failed to fetch default image")
			}
		} else {
			return nil, errors.Wrap(err, "failed to fetch image")
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: %s", resp.Status)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode image")
	}

	return img, nil
}

func GetImageMetadata(img image.Image) map[string]interface{} {
	bounds := img.Bounds()
	metadata := map[string]interface{}{
		"width":  bounds.Dx(),
		"height": bounds.Dy(),
		"format": util.GetImageFormat(img),
	}
	return metadata
}
func OverlayImage(baseImg, overlayImg image.Image, params util.QueryParams) image.Image {
	bounds := baseImg.Bounds()
	centerX := bounds.Dx() / 2
	centerY := bounds.Dy() / 2

	offset := image.Pt(centerX-overlayImg.Bounds().Dx()/2, centerY-overlayImg.Bounds().Dy()/2)
	return imaging.Overlay(baseImg, overlayImg, offset, 1.0)
}

func EncodeImageResponse(w http.ResponseWriter, img image.Image, params util.QueryParams) {
	format := params.Output
	if format == "" {
		format = "jpeg"
	}

	w.Header().Set("Content-Type", "image/"+format)
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", params.MaxAge))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if params.Filename != "" {
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, params.Filename))
	}

	if params.Encoding == "base64" {
		var buf bytes.Buffer
		if err := util.EncodeImage(&buf, img, format, params.Quality); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		base64Data := base64.StdEncoding.EncodeToString(buf.Bytes())
		httputil.JSONResponse(w, map[string]string{
			"data": fmt.Sprintf("data:image/%s;base64,%s", format, base64Data),
		})
		return
	}

	if err := util.EncodeImage(w, img, format, params.Quality); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}