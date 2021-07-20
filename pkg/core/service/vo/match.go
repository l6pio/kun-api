package vo

import (
	"l6p.io/kun/api/pkg/core/service/vo/match"
)

type Match struct {
	Artifact               *match.Artifact        `json:"artifact"`
	Vulnerability          *match.Vulnerability   `json:"vulnerability"`
	RelatedVulnerabilities []*match.Vulnerability `json:"relatedVulnerabilities"`
}
