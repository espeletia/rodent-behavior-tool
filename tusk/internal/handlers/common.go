package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"tusk/internal/domain"
	"tusk/internal/util"

	"go.uber.org/zap"
)

func getCursorFromRequest(r *http.Request) (*domain.CursorInput, error) {
	cursorInput := &domain.CursorInput{}

	if before := r.URL.Query().Get("before"); before != "" {
		cursorInput.Before = &before
	}

	if after := r.URL.Query().Get("after"); after != "" {
		cursorInput.After = &after
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			cursorInput.Limit = util.ToPointer(limit)
		} else {
			zap.L().Error("failed to parse limit", zap.Error(err))
			return nil, errors.New("limit is not a number")
		}
	}

	return cursorInput, nil
}

func mapCursorInput(ctx context.Context, cursor *domain.CursorInput) (*domain.OffsetLimit, error) {
	if cursor == nil {
		return defaultOffsetLimit(), nil
	}

	offset, ascending, err := parseCursorOffsets(cursor)
	if err != nil {
		return nil, err
	}

	return &domain.OffsetLimit{
		Offset:    offset,
		Limit:     mapCursorLimit(cursor.Limit),
		Ascending: ascending,
	}, nil
}

func parseCursorOffsets(cursor *domain.CursorInput) (offset int64, ascending bool, err error) {
	if cursor.After != nil && *cursor.After != "" {
		offset, err = strconv.ParseInt(*cursor.After, 10, 64)
	}

	if cursor.Before != nil && *cursor.Before != "" {
		if offset, err = strconv.ParseInt(*cursor.Before, 10, 64); err == nil {
			ascending = true
		}
	}

	return offset, ascending, err
}

func mapCursorLimit(cursorLimit *int) int32 {
	if cursorLimit == nil || *cursorLimit <= 0 {
		return 10
	}
	if *cursorLimit > 50 {
		return 50
	}
	return int32(*cursorLimit)
}

func defaultOffsetLimit() *domain.OffsetLimit {
	return &domain.OffsetLimit{
		Offset:    0,
		Limit:     10,
		Ascending: false,
	}
}

