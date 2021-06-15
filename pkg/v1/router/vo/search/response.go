package search

import "l6p.io/kun/api/pkg/core/cve/vo"

type Response struct {
	Request interface{}  `json:"request"`
	Reports []*vo.Report `json:"reports"`
}
