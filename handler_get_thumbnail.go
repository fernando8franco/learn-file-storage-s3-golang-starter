package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

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

	filename := path.Base(*video.ThumbnailURL)
	filepath := filepath.Join(cfg.assetsRoot, filename)

	fileData, err := os.ReadFile(filepath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't open image", err)
		return
	}

	mediaType := http.DetectContentType(fileData)

	w.Header().Set("Content-Type", mediaType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileData)))

	_, err = w.Write(fileData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error writing response", err)
		return
	}
}
