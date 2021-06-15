package search

type ByImageID struct {
	ImageID string `json:"imageID" validate:"required"`
}
