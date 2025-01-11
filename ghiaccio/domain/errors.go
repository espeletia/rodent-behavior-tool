package domain

type Error struct {
	Message string
	Code    int
}

func (e Error) Error() string { return e.Message }
