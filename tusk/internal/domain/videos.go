package domain

import "github.com/google/uuid"

type Video struct {
	ID            uuid.UUID
	Video         MediaFile
	OwnerId       uuid.UUID
	Description   *string
	Name          string
	AnalysedVideo *MediaFile
}

type CreateVideoDto struct {
	VideoUrl    string
	Description *string
	Name        string
}

type AnalystResult struct {
	ID      uuid.UUID
	VideoID uuid.UUID
	MediaID uuid.UUID
	Url     string
}

type VideosCursored struct {
	Data   []Video
	Cursor Cursor
}

type EncodingResult struct {
	ID      uuid.UUID
	VideoID uuid.UUID
	MediaID uuid.UUID
	Url     string
}
