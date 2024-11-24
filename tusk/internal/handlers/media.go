package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"tusk/internal/handlers/models"
	"tusk/internal/ports"
	"tusk/internal/util"

	"github.com/nextap-solutions/openapi3Struct"

	"tusk/internal/usecases"

	"github.com/getkin/kin-openapi/openapi3"

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

var UploadOp = openapi3Struct.Path{
	Path: "/upload",
	Item: openapi3.PathItem{
		Put: &openapi3.Operation{
			Tags:        []string{"Media"},
			OperationID: "uploadMedia",
			Description: "Upload a media file (image or video).",
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "File upload in `multipart/form-data` format. Allowed file types: JPEG, PNG, MP4, MPEG.",
					Required:    true,
					Content: map[string]*openapi3.MediaType{
						"multipart/form-data": {
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: "object",
									Properties: map[string]*openapi3.SchemaRef{
										"file": {
											Value: &openapi3.Schema{
												Type:   "string",
												Format: "binary",
											},
										},
									},
									Required: []string{"file"},
								},
							},
						},
					},
				},
			},
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Description: util.ToPointer("Uploaded file s3 url"),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: openapi3.NewSchemaRef("#/components/schemas/UploadResponse", nil),
							},
						},
					},
				},
				"400": {
					Value: &openapi3.Response{
						Description: util.ToPointer("Invalid file or bad request."),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "object",
										Properties: map[string]*openapi3.SchemaRef{
											"error": {
												Value: &openapi3.Schema{
													Type:    "string",
													Example: "Invalid file type.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"413": {
					Value: &openapi3.Response{
						Description: util.ToPointer("File too large."),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "object",
										Properties: map[string]*openapi3.SchemaRef{
											"error": {
												Value: &openapi3.Schema{
													Type:    "string",
													Example: "File too big.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"500": {
					Value: &openapi3.Response{
						Description: util.ToPointer("Internal server error."),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "object",
										Properties: map[string]*openapi3.SchemaRef{
											"error": {
												Value: &openapi3.Schema{
													Type:    "string",
													Example: "Failed to upload file.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
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
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		// Reset file pointer
		file.Seek(0, 0)

		tempFile, err := ports.NewTempFile(tmpDir, file)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		url, err := mh.mediaUsecase.DefaultFileUpload(ctx, tempFile.Path(), contentType, filename)
		if err != nil {
			http.Error(w, "Failed to upload file", http.StatusInternalServerError)
			zap.L().Error("S3 Upload Error:", zap.Error(err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(models.UploadResponse{
			UploadUrl: url,
		})
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
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
	dir, err := os.MkdirTemp("", fmt.Sprintf("%s_*", str))
	if err != nil {
		return "", err
	}
	return dir, nil
}
