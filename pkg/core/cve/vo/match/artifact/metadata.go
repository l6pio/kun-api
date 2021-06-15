package artifact

type Metadata struct {
	VirtualPath   string `json:"VirtualPath"`
	PomArtifactID string `json:"PomArtifactID"`
	PomGroupID    string `json:"PomGroupID"`
}
