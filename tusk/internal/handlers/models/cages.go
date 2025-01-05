package models

import (
	"time"

	"github.com/google/uuid"
)

// swagger:model
type CageCreationResponse struct {
	ActivationCode string `json:"activation_code"`
}

// swagger:model
type Cage struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Register    time.Time `json:"register"`
}

// swagger:model
type Cages struct { // TODO: add cursoring
	Data []Cage `json:"data"`
}
