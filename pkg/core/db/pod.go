package db

import (
	"gopkg.in/mgo.v2/bson"
	v1 "k8s.io/api/core/v1"
	"l6p.io/kun/api/pkg/core"
	"l6p.io/kun/api/pkg/core/db/vo"
)

func SavePod(conf *core.Config, pod *vo.Pod) error {
	session, col, err := GetCol(conf, "pod")
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = col.Upsert(bson.M{"name": pod.Name}, pod)
	return err
}

func GetTotalPods(conf *core.Config) (int, error) {
	session, col, err := GetCol(conf, "pod")
	if err != nil {
		return 0, err
	}
	defer session.Close()

	var ret map[string]int
	err = col.Pipe([]bson.M{
		{"$count": "count"},
	}).One(&ret)
	if err != nil {
		return 0, err
	}
	return ret["count"], nil
}

func GetPodCountByPhase(conf *core.Config) (map[string]int, error) {
	session, col, err := GetCol(conf, "pod")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var res []map[string]interface{}
	err = col.Pipe([]bson.M{
		{"$group": bson.M{
			"_id":   "$phase",
			"count": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{
			"_id":   false,
			"phase": "$_id",
			"count": "$count",
		}},
	}).All(&res)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]int)
	for _, d := range res {
		phase := d["phase"].(string)
		if phase == "" {
			phase = "Unknown"
		}
		ret[phase] = d["count"].(int)
	}
	return ret, nil
}

func GetPodCountByStatus(conf *core.Config) (map[string]int, error) {
	session, col, err := GetCol(conf, "pod")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var res []map[string]interface{}
	err = col.Pipe([]bson.M{
		{"$group": bson.M{
			"_id":   "$status",
			"count": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{
			"_id":    false,
			"status": "$_id",
			"count":  "$count",
		}},
	}).All(&res)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]int)
	for _, d := range res {
		status := d["status"].(string)
		if status == "" {
			status = "Unknown"
		}
		ret[status] = d["count"].(int)
	}
	return ret, nil
}

func GetPodCountByNamespace(conf *core.Config) (map[string]int, error) {
	session, col, err := GetCol(conf, "pod")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var res []map[string]interface{}
	err = col.Pipe([]bson.M{
		{"$group": bson.M{
			"_id":   "$namespace",
			"count": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{
			"_id":       false,
			"namespace": "$_id",
			"count":     "$count",
		}},
	}).All(&res)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]int)
	for _, d := range res {
		namespace := d["namespace"].(string)
		ret[namespace] = d["count"].(int)
	}
	return ret, nil
}

func FindRunningPodTimeline(conf *core.Config) ([]interface{}, error) {
	session, col, err := GetCol(conf, "pod")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var ret []interface{}
	var stages []bson.M
	stages = append(stages,
		bson.M{"$match": bson.M{"phase": v1.PodRunning}},
		bson.M{"$match": bson.M{"status": "Running"}},
		bson.M{"$match": bson.M{"ready": bson.M{"$gt": 0}}},
		bson.M{"$group": bson.M{
			"_id": bson.M{
				"$subtract": []bson.M{
					{"$toLong": "$ready"},
					{"$mod": []interface{}{bson.M{"$toLong": "$ready"}, 60 * 1000}},
				},
			},
			"count": bson.M{"$sum": 1},
		}},
		bson.M{"$addFields": bson.M{"timestamp": "$_id"}},
	)
	err = col.Pipe(stages).All(&ret)
	return ret, err
}

func RemovePods(conf *core.Config) error {
	session, col, err := GetCol(conf, "pod")
	if err != nil {
		return err
	}
	defer session.Close()
	_ = col.DropCollection()
	return nil
}
