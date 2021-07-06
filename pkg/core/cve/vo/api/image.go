package api

type Image struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	ArtCount int64  `json:"artCount"`
	VulCount int64  `json:"vulCount"`
}
