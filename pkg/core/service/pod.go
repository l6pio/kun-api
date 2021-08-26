package service

import (
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db"
	"l6p.io/kun/api/pkg/core/service/vo"
)

func GetPodCount(conf *core.Config) (*vo.PodCount, error) {
	total, err := db.GetTotalPods(conf)
	if err != nil {
		return nil, err
	}

	countByPhase, err := db.GetPodCountByPhase(conf)
	if err != nil {
		return nil, err
	}

	countByStatus, err := db.GetPodCountByStatus(conf)
	if err != nil {
		return nil, err
	}

	countByNamespace, err := db.GetPodCountByNamespace(conf)
	if err != nil {
		return nil, err
	}

	return &vo.PodCount{
		Total:            total,
		CountByStatus:    countByStatus,
		CountByPhase:     countByPhase,
		CountByNamespace: countByNamespace,
	}, nil
}
