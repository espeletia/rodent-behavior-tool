package util

import (
	"strconv"
	"tusk/internal/domain"
)

func ToPointer[T any](v T) *T {
	return &v
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func BuildCursorWithOffsetCursor[T any](data []T, offset int64, limit int32) domain.Cursor {
	var before *string
	if offset >= int64(limit) {
		beforeStr := strconv.FormatInt(offset-int64(limit), 10)
		before = &beforeStr
	}
	var after *string
	if int64(limit) <= int64(len(data)) {
		afterStr := strconv.FormatInt(offset+int64(limit), 10)
		after = &afterStr
	}

	return domain.Cursor{
		After:  after,
		Before: before,
	}
}
