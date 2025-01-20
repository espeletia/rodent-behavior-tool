package models

// swagger:model
type Cursor struct {
	After  *string `json:"after"`
	Before *string `json:"before"`
}
