package domain

import "github.com/google/uuid"

type MessageWrapper[T any] struct {
	Message T
	Err     *string
}

type VideoEncodingMessage struct {
	ID      uuid.UUID `json:"id"`
	VideoID uuid.UUID `json:"video_id"`
	MediaID uuid.UUID `json:"media_id"`
	Url     string    `json:"url"`
}

type VideoEncodingResultMessage struct {
	ID      uuid.UUID `json:"id"`
	VideoID uuid.UUID `json:"video_id"`
	MediaID uuid.UUID `json:"media_id"`
	Url     string    `json:"url"`
}

type AnalystJobMessage struct {
	ID      uuid.UUID `json:"job_id"`
	VideoID uuid.UUID `json:"video_id"`
	MediaID uuid.UUID `json:"media_id"`
	Url     string    `json:"url"`
}

type AnalystJobResultMessage struct {
	ID      uuid.UUID `json:"job_id"`
	VideoID uuid.UUID `json:"video_id"`
	MediaID uuid.UUID `json:"media_id"`
	Url     string    `json:"url"`
}
