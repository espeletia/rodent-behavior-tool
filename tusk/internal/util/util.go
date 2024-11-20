package util

func ToPointer[T any](v T) *T {
	return &v
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
