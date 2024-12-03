package domain

import "github.com/google/uuid"

type MessageWrapper[T any] struct {
	Message T
	Err     *string
}

type VideoEncodingMessage struct {
	ID  uuid.UUID `json:"id"`
	Url string    `json:"url"`
}

type AnalystJobMessage struct {
	ID  uuid.UUID `json:"job_id"`
	Url string    `json:"url"`
}
