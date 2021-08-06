package cve

import (
	match2 "l6p.io/kun/api/pkg/core/service/vo/cve/match"
)

type Match struct {
	Artifact               *match2.Artifact        `json:"artifact"`
	Vulnerability          *match2.Vulnerability   `json:"vulnerability"`
	RelatedVulnerabilities []*match2.Vulnerability `json:"relatedVulnerabilities"`
}
