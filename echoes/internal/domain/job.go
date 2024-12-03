package domain

import "github.com/google/uuid"

type VideoEncodingJob struct {
	ID  uuid.UUID
	URl string
}

type JobResult struct {
	ID                 uuid.UUID
	Label              string
	LocalFileSrc       string
	FileDestinationSrc string
	Type               string
}
