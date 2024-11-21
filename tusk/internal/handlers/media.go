package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"tusk/internal/ports"
	"tusk/internal/usecases"

	"go.uber.org/zap"
	// "go.uber.org/zap"
)

type MediaHandler struct {
	mediaUsecase  *usecases.MediaUsecase
	maxUploadSize int64
}

var allowedTypes = []string{"image/jpeg", "image/png", "video/mp4", "video/mpeg"}

func NewMediaHandler(mediaUsecase *usecases.MediaUsecase) *MediaHandler {
	return &MediaHandler{
		mediaUsecase:  mediaUsecase,
		maxUploadSize: 50 * 1024 * 1024,
	}
}

func (mh *MediaHandler) Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r.Body = http.MaxBytesReader(w, r.Body, mh.maxUploadSize)
		if err := r.ParseMultipartForm(mh.maxUploadSize); err != nil {
			http.Error(w, "File too big", http.StatusBadRequest)
			return
		}

		// Parse the uploaded file
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Invalid file upload", http.StatusInternalServerError)
			zap.L().Error("Invalid file upload")
			return
		}
		defer file.Close()

		// Get the content type
		buffer := make([]byte, 512)
		if _, err := file.Read(buffer); err != nil {
			http.Error(w, "Unable to read file", http.StatusInternalServerError)
			zap.L().Error("Unable to read file")
			return
		}
		contentType := http.DetectContentType(buffer)

		// Validate file type
		if !isValidFileType(contentType) {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			zap.L().Error("Invalid file type", zap.String("contentType", contentType))
			return
		}

		// Sanitize the filename
		zap.L().Info("debug line", zap.String("filename", handler.Filename))
		filename := filepath.Base(handler.Filename)

		tmpDir, err := createTempDir("tmp")

		// Reset file pointer
		file.Seek(0, 0)

		tempFile, err := ports.NewTempFile(tmpDir, file)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		err = mh.mediaUsecase.DefaultFileUpload(ctx, tempFile.Path(), contentType, filename)
		if err != nil {
			http.Error(w, "Failed to upload file", http.StatusInternalServerError)
			zap.L().Error("S3 Upload Error:", zap.Error(err))
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "File uploaded successfully: %s", filename)
	}
}

func isValidFileType(contentType string) bool {
	for _, t := range allowedTypes {
		if t == contentType {
			return true
		}
	}
	return false
}

func createTempDir(str string) (string, error) {
	dir, err := os.MkdirTemp("", fmt.Sprintf("%d_*", str))
	if err != nil {
		return "", err
	}
	return dir, nil
}
