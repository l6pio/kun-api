package vo

type PodsOverview struct {
	Total         int            `json:"total"`
	TotalRunning  int            `json:"totalRunning"`
	CountByPhase  map[string]int `json:"countByPhase"`
	CountByStatus map[string]int `json:"countByStatus"`
}
