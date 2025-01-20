package models

// swagger:model
type CageCreationResponse struct {
	ActivationCode string `json:"activation_code"`
	SecretToken    string `json:"secret_token"`
}

// swagger:model
type CageMessageRequest struct {
	Revision  int64   `json:"revision"`
	Water     int64   `json:"water"`
	Food      int64   `json:"food"`
	Light     int64   `json:"light"`
	Temp      int64   `json:"temp"`
	Humidity  int64   `json:"humidity"`
	VideoUrl  *string `json:"video_url,omitempty"`
	Timestamp int64   `json:"timestamp"`
}

// swagger:model
type CageMessage struct {
	CageID    string  `json:"cage_id"`
	Revision  int64   `json:"revision"`
	Water     int64   `json:"water"`
	Food      int64   `json:"food"`
	Light     int64   `json:"light"`
	Temp      int64   `json:"temp"`
	Humidity  int64   `json:"humidity"`
	VideoUrl  *string `json:"video_url,omitempty"`
	Timestamp int64   `json:"timestamp"`
}

// swagger:model
type CageMessagesCursored struct {
	Data   []CageMessage `json:"data"`
	Cursor Cursor        `json:"cursor"`
}

// swagger:model
type Cages struct { // TODO: add cursoring
	Data []Cage `json:"data"`
}

// swagger:model
type Cage struct {
	ID          string  `json:"id"`
	UserID      *string `json:"user_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Register    int64   `json:"register"`
}
