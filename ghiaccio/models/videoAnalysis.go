package models

// swagger:model
type CreateVideoDto struct {
	VideoUrl    string
	Description *string
	Name        string
}

// swagger:model
type VideoAnalysis struct {
	ID            string     `json:"id"`
	Video         MediaFile  `json:"video"`
	OwnerId       *string    `json:"owner_id"`
	CageId        *string    `json:"cage_id"`
	Description   *string    `json:"description,omitempty"`
	Name          string     `json:"name"`
	AnalysedVideo *MediaFile `json:"analysed_video,omitempty"`
}

// swagger:model
type CursoredVideoAnalysis struct {
	Data   []VideoAnalysis `json:"data"`
	Cursor Cursor          `json:"cursor"`
}
