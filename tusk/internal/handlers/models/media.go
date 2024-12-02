package models

// swagger:model
type UploadResponse struct {
	UploadUrl string `json:"upload_url"`
}

// swagger:model
type MediaFile struct {
	ID         string  `json:"id"`
	MimeType   string  `json:"mime_type"`
	Variant    string  `json:"variant"`
	EntityType string  `json:"entity_type"`
	MasterID   *string `json:"master_id,omitempty"`
	Url        string  `json:"url"`
	Duration   *int64  `json:"duration,omitempty"`
	Size       int64   `json:"size"`
	Width      int64   `json:"width"`
	Height     int64   `json:"height"`
	Master     bool    `json:"master"`
}
