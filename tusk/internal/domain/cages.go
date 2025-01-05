package domain

import (
	"time"

	"github.com/google/uuid"
)

type Cage struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Register    time.Time
}
