package vo

type PodCount struct {
	Total         int            `json:"total"`
	CountByPhase  map[string]int `json:"countByPhase"`
	CountByStatus map[string]int `json:"countByStatus"`
}
