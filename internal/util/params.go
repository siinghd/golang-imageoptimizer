package util

import (
	"net/http"
	"strconv"
)

type QueryParams struct {
	URL            string
	Width          int
	Height         int
	DPR            float64
	Fit            string
	ContainBgColor string
	WithoutEnlarge bool
	Background     string
	Blur           float64
	Gamma          float64
	Modulate       string
	Sharpen        float64
	Tint           string
	Encoding       string
	MaxAge         int
	Compression    int
	Default        string
	Filename       string
	Interlace      bool
	Pages          int
	Output         string
	Page           int
	Quality        int
	TextColor      string
	FontSize       int
	FontFamily     string
	Text           string
	TextAlign      string
	RoundedCorners bool
	CornerRadius   int
	TextBaseline   string
}

func ParseQueryParams(r *http.Request) QueryParams {
	query := r.URL.Query()
	params := QueryParams{
		URL:            query.Get("url"),
		Width:          ParseInt(query.Get("w"), 0),
		Height:         ParseInt(query.Get("h"), 0),
		DPR:            ParseFloat(query.Get("dpr"), 1),
		Fit:            query.Get("fit"),
		ContainBgColor: query.Get("cbg"),
		WithoutEnlarge: query.Get("we") == "true",
		Background:     query.Get("bg"),
		Blur:           ParseFloat(query.Get("blur"), 0),
		Gamma:          ParseFloat(query.Get("gam"), 0),
		Modulate:       query.Get("mod"),
		Sharpen:        ParseFloat(query.Get("sharp"), 0),
		Tint:           query.Get("tint"),
		Encoding:       query.Get("encoding"),
		MaxAge:         ParseInt(query.Get("maxage"), 31536000),
		Compression:    ParseInt(query.Get("l"), 6),
		Default:        query.Get("default"),
		Filename:       query.Get("filename"),
		Interlace:      query.Get("il") == "true",
		Pages:          ParseInt(query.Get("n"), 0),
		Output:         query.Get("output"),
		Page:           ParseInt(query.Get("page"), 0),
		Quality:        ParseInt(query.Get("q"), 80),
		TextColor:      query.Get("txtColor"),
		FontSize:       ParseInt(query.Get("fontSize"), 48),
		FontFamily:     query.Get("fontFamily"),
		Text:           query.Get("text"),
		TextAlign:      query.Get("textAlign"),
		RoundedCorners: query.Get("roundedCorners") == "true",
		CornerRadius:   ParseInt(query.Get("cornerRadius"), 20),
		TextBaseline:   query.Get("textBaseline"),
	}
	return params
}
func ParseInt(str string, defaultValue int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return value
}

func ParseFloat(str string, defaultValue float64) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return defaultValue
	}
	return value
}