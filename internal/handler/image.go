package handler

import (
	"fmt"
	"image"
	"net/http"

	"github.com/siinghd/golang-imageoptimizer/internal/httputil"
	"github.com/siinghd/golang-imageoptimizer/internal/service"
	"github.com/siinghd/golang-imageoptimizer/internal/util"
)

func ProcessImageHandler(w http.ResponseWriter, r *http.Request) {
	params := util.ParseQueryParams(r)

	if params.URL == "" && params.Text == "" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, util.Htmlpage)
		return
	}

	var img image.Image
	var err error

	if params.URL != "" {
		img, err = service.FetchImage(params.URL, params.Default)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if params.Text != "" {
			textOverlayParams := params
			textOverlayParams.Background = "transparent"
			textImg, err := service.RenderTextToImage(params.Text, textOverlayParams)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			img = service.OverlayImage(img, textImg, params)
		}

		if params.Output == "json" {
			metadata := service.GetImageMetadata(img)
			httputil.JSONResponse(w, metadata)
			return
		}
	} else if params.Text != "" {
		img, err = service.RenderTextToImage(params.Text, params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Missing URL or Text parameter", http.StatusBadRequest)
		return
	}

	processedImg, err := service.ProcessImageWithParams(img, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.EncodeImageResponse(w, processedImg, params)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
}