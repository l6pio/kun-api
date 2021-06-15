package vo

import (
	"l6p.io/kun/api/pkg/core/cve/vo/match"
)

type Match struct {
	Vulnerability *match.Vulnerability `json:"vulnerability"`
	Artifact      *match.Artifact      `json:"artifact"`
}
