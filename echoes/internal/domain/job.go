package domain

type VideoEncodingJob struct {
	ID  int64
	URl string
}

type JobResult struct {
	ID                 int64
	Label              string
	LocalFileSrc       string
	FileDestinationSrc string
	Type               string
}
