package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerThumbnailGet(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid video ID", err)
		return
	}

	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find video", err)
		return
	}

	url := strings.Split(*video.ThumbnailURL, ";")
	mediaType := strings.Split(url[0], ":")
	data := strings.Split(url[1], ":")

	w.Header().Set("Content-Type", mediaType[1])
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data[1])))

	_, err = w.Write([]byte(data[1]))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error writing response", err)
		return
	}
}
