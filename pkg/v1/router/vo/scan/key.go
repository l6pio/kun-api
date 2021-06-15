package scan

type Key struct {
	ImageRepo string `json:"imageRepo" validate:"required"`
	ImageTag  string `json:"imageTag" validate:"required"`
}
