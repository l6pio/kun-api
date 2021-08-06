package cve

type Report struct {
	Matches []*Match `json:"matches"`
	Source  *Source  `json:"source"`
}
