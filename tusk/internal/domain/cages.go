package domain

import (
	"time"

	"github.com/google/uuid"
)

type Cage struct {
	ID          uuid.UUID
	UserID      *uuid.UUID
	Name        string
	Description *string
	Register    time.Time
}

type CageMessasgesCursored struct {
	Data   []CageMessage
	Cursor Cursor
}

type CageMessage struct {
	ID        int64
	CageID    uuid.UUID
	Revision  int64
	Water     int64
	Food      int64
	Light     int64
	Temp      int64
	Humidity  int64
	VideoUrl  *string
	VideoID   *uuid.UUID
	Timestamp time.Time
}

type CageMessageData struct {
	Revision  int64
	Water     int64
	Food      int64
	Light     int64
	Temp      int64
	Humidity  int64
	VideoUrl  *string
	Timestamp time.Time
}
