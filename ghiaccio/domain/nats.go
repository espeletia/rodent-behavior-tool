package domain

type MessageWrapper[T any] struct {
	Message T
	Err     *string
}

type VideoEncodingMessage struct {
	ID  int64  `json:"id"`
	Url string `json:"url"`
}
