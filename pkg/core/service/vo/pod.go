package vo

type PodsOverview struct {
	Count         int            `json:"count"`
	CountByStatus map[string]int `json:"countByStatus"`
}
