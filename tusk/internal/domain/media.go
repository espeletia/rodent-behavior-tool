package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	MediaVariantOriginal     string = "original"
	MediaVariantAnalysedRaw  string = "analysed_raw"
	MediaVariantAnalysedX264 string = "analysed_x264"
)

type FileMetadata struct {
	ContentType   string
	FileExtension string
	Size          int64
}

type MediaFile struct {
	ID         uuid.UUID
	MimeType   string
	Variant    string
	EntityType string
	MasterID   *uuid.UUID
	Url        string
	Created    time.Time
	Duration   *int64
	Size       int64
	Width      int64
	Height     int64
	Master     bool
}
