package match

import (
	artifact2 "l6p.io/kun/api/pkg/core/service/vo/cve/match/artifact"
)

type Artifact struct {
	Name      string                `json:"name"`
	Version   string                `json:"version"`
	Type      string                `json:"type"`
	Locations []*artifact2.Location `json:"locations"`
	Language  string                `json:"language"`
	Licenses  []string              `json:"licenses"`
	Cpes      []string              `json:"cpes"`
	Purl      string                `json:"purl"`
	Metadata  *artifact2.Metadata   `json:"metadata"`
}
