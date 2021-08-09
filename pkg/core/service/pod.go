package service

import (
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/service/vo"
)

func GetPodsOverview(conf *core.Config) (*vo.PodsOverview, error) {
	count, err := db.GetTotalPods(conf)
	if err != nil {
		return nil, err
	}

	countByStatus, err := db.GetPodCountByStatus(conf)
	if err != nil {
		return nil, err
	}

	countByPhase, err := db.GetPodCountByPhase(conf)
	if err != nil {
		return nil, err
	}

	return &vo.PodsOverview{
		Count:         count,
		CountByStatus: countByStatus,
		CountByPhase:  countByPhase,
	}, nil
}
