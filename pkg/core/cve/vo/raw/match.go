package raw

import (
	match2 "l6p.io/kun/api/pkg/core/cve/vo/raw/match"
)

type Match struct {
	Artifact      *match2.Artifact      `json:"artifact"`
	Vulnerability *match2.Vulnerability `json:"vulnerability"`
}
