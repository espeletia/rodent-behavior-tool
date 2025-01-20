package domain

type ApiPath struct {
	Name    string
	Methods []string
}

type CursorInput struct {
	Before *string
	After  *string
	Limit  *int
}

type OffsetLimit struct {
	Offset    int64
	Limit     int32
	Ascending bool
}

type Cursor struct {
	After  *string
	Before *string
}
